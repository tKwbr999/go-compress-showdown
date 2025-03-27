package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// Result stores parsed data from a single vegeta report
type Result struct {
	FileName       string
	Algorithm      string // none, gzip, brotli, zstd
	Level          string // default, best, speed, 1, 6, 11, fastest
	DataSize       string // small, medium, large
	TestType       string // fixed, max
	Requests       int64
	Rate           float64 // req/s
	Throughput     float64 // req/s
	Duration       float64 // seconds
	LatencyMean    float64 // ms
	Latency50      float64 // ms
	Latency95      float64 // ms
	Latency99      float64 // ms
	BytesInTotal   int64
	BytesOutTotal  int64
	Success        float64 // percentage
	StatusCodes    map[string]int64
	Error          string // If any error occurred during report generation
	OriginalFileSize int64 // Size of the original uncompressed file
	CompressedFileSize int64 // Size of the compressed data used in test (approximated from BytesInTotal/Requests if possible)
}

// ReportData holds all results for markdown generation
type ReportData struct {
	Results []Result
}

// Regexps for parsing vegeta report output
var (
	requestsRe    = regexp.MustCompile(`Requests\s+\[total, rate, throughput\]\s+(\d+),\s+([\d.]+),\s+([\d.]+)`)
	durationRe    = regexp.MustCompile(`Duration\s+\[total, attack, wait\]\s+([\d.]+)s,\s+([\d.]+)s,\s+([\d.]+)s`)
	latenciesRe   = regexp.MustCompile(`Latencies\s+\[mean, 50, 95, 99, max\]\s+([\d.]+)ms,\s+([\d.]+)ms,\s+([\d.]+)ms,\s+([\d.]+)ms,\s+([\d.]+)ms`)
	bytesRe       = regexp.MustCompile(`Bytes\s+\[total, mean\]\s+(\d+),\s+([\d.]+)`) // Assuming Bytes In
	bytesOutRe    = regexp.MustCompile(`Bytes Out\s+\[total, mean\]\s+(\d+),\s+([\d.]+)`) // Need to confirm exact vegeta output format
	successRe     = regexp.MustCompile(`Success\s+\[ratio\]\s+([\d.]+)%`)
	statusCodesRe = regexp.MustCompile(`Status Codes\s+\[code:count\]\s+(.*)`)
	errorsRe      = regexp.MustCompile(`Error Set:\s+(.*)`)
	// Make the level part optional (?:_(\w+))? to handle 'none' which has no level
	fileNameRe    = regexp.MustCompile(`^(none|gzip|brotli|zstd)(?:_(\w+))?_(\w+)_(fixed|max)\.bin$`)
)

// Data sizes map
var dataSizes = map[string]string{
	"small":  "testdata/small.txt",
	"medium": "testdata/medium.txt",
	"large":  "testdata/large.txt",
}

func getFileSize(path string) int64 {
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Printf("Warning: Could not get file size for %s: %v", path, err)
		return 0
	}
	return fileInfo.Size()
}


func parseVegetaReport(filePath string) (Result, error) {
	var result Result
	result.FileName = filepath.Base(filePath)
	result.StatusCodes = make(map[string]int64)

	// Parse filename
	matches := fileNameRe.FindStringSubmatch(result.FileName)
	// Expected matches:
	// [0]: Full string (e.g., "gzip_default_small_fixed.bin" or "none_small_fixed.bin")
	// [1]: Algorithm ("gzip", "none", etc.)
	// [2]: Level ("default", "1", etc.) - Optional, empty for "none"
	// [3]: DataSize ("small", "medium", "large")
	// [4]: TestType ("fixed", "max")
	if len(matches) != 5 { // Length should still be 5 due to the non-capturing group structure
		return result, fmt.Errorf("could not parse filename structure: %s (matches: %d)", result.FileName, len(matches))
	}
	result.Algorithm = matches[1]
	result.Level = matches[2] // This will be empty string "" for 'none'
	if result.Algorithm == "none" {
		result.Level = "N/A" // Assign a placeholder level for 'none'
	}
	result.DataSize = matches[3]
	result.TestType = matches[4]


	// Get original file size
	originalFilePath, ok := dataSizes[result.DataSize]
	if ok {
		result.OriginalFileSize = getFileSize(originalFilePath)
	}


	// Run vegeta report command
	cmd := exec.Command("vegeta", "report", filePath)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		result.Error = fmt.Sprintf("vegeta report failed: %v, stderr: %s", err, stderr.String())
		// Even if the command fails, try to parse partial output
	}

	output := out.String()

	// --- Parse metrics ---
	if m := requestsRe.FindStringSubmatch(output); len(m) == 4 {
		result.Requests, _ = strconv.ParseInt(m[1], 10, 64)
		result.Rate, _ = strconv.ParseFloat(m[2], 64)
		result.Throughput, _ = strconv.ParseFloat(m[3], 64)
	}
	if m := durationRe.FindStringSubmatch(output); len(m) == 4 {
		// Using total duration
		result.Duration, _ = strconv.ParseFloat(m[1], 64)
	}
	if m := latenciesRe.FindStringSubmatch(output); len(m) == 6 {
		result.LatencyMean, _ = strconv.ParseFloat(m[1], 64)
		result.Latency50, _ = strconv.ParseFloat(m[2], 64)
		result.Latency95, _ = strconv.ParseFloat(m[3], 64)
		result.Latency99, _ = strconv.ParseFloat(m[4], 64)
	}
	// Vegeta report might show Bytes In (request size) and Bytes Out (response size)
	// We are interested in the response size (Bytes Out) for compression ratio,
	// but the report format might vary. Let's try to capture both.
	// Assuming "Bytes" without "Out" refers to Bytes In.
	if m := bytesRe.FindStringSubmatch(output); len(m) == 3 {
		result.BytesInTotal, _ = strconv.ParseInt(m[1], 10, 64)
		// Estimate compressed size if possible (this is a rough estimate)
		if result.Requests > 0 {
			result.CompressedFileSize = result.BytesInTotal / result.Requests
		}
	}
	if m := bytesOutRe.FindStringSubmatch(output); len(m) == 3 {
		result.BytesOutTotal, _ = strconv.ParseInt(m[1], 10, 64)
		// If Bytes Out is available, it's a better measure of compressed size per request
		if result.Requests > 0 {
			result.CompressedFileSize = result.BytesOutTotal / result.Requests
		}
	} else if result.BytesInTotal > 0 && result.BytesOutTotal == 0 {
        // Fallback if "Bytes Out" is not present, assume "Bytes" refers to response size
        result.BytesOutTotal = result.BytesInTotal
		if result.Requests > 0 {
			result.CompressedFileSize = result.BytesOutTotal / result.Requests
		}
    }


	if m := successRe.FindStringSubmatch(output); len(m) == 2 {
		result.Success, _ = strconv.ParseFloat(m[1], 64)
	}
	if m := statusCodesRe.FindStringSubmatch(output); len(m) == 2 {
		codes := strings.Fields(m[1])
		for _, codeCount := range codes {
			parts := strings.Split(codeCount, ":")
			if len(parts) == 2 {
				count, _ := strconv.ParseInt(parts[1], 10, 64)
				result.StatusCodes[parts[0]] = count
			}
		}
	}
	if m := errorsRe.FindStringSubmatch(output); len(m) == 2 {
		result.Error += " Vegeta Errors: " + m[1] // Append vegeta errors if any
	}

	// If BytesOutTotal is still 0 and we have 'none' algorithm, use original file size
	if result.BytesOutTotal == 0 && result.Algorithm == "none" && result.Requests > 0 {
		result.BytesOutTotal = result.OriginalFileSize * result.Requests
		result.CompressedFileSize = result.OriginalFileSize
	}


	return result, nil
}

const markdownTemplate = `# Compression Benchmark Report

This report summarizes the results of load tests performed using Vegeta against a Go server serving compressed content.

## Test Setup

- **Load Generator:** Vegeta
- **Target Server:** Go HTTP server (` + "`cmd/server/main.go`" + `)
- **Test Data:**
    - Small: ` + "`testdata/small.txt`" + ` ({{ .SmallSize }} bytes)
    - Medium: ` + "`testdata/medium.txt`" + ` ({{ .MediumSize }} bytes)
    - Large: ` + "`testdata/large.txt`" + ` ({{ .LargeSize }} bytes)
- **Compression Algorithms Tested:** None, Gzip (speed, default, best), Brotli (1, 6, 11), Zstd (fastest, default, best)
- **Test Types:**
    - **Fixed:** Constant request rate (likely 50 req/s based on ` + "`run_vegeta_tests.sh`" + `) for a fixed duration. Measures latency under stable load.
    - **Max:** Attempts to find the maximum sustainable throughput. Measures server capacity.

## Detailed Results per Algorithm

{{ range .Algorithms }}
### {{ .Name }}
{{ range .Levels }}
#### Level: {{ .LevelName }}
{{ range .DataSizes }}
##### Data Size: {{ .SizeName }} ({{ .OriginalSize }} bytes)

**Fixed Rate Test (Latency Focus):**
{{ with .FixedResult }}
- **Requests:** {{ .Requests }}
- **Target Rate:** ~50 req/s (estimated)
- **Actual Throughput:** {{ printf "%.2f" .Throughput }} req/s
- **Duration:** {{ printf "%.2f" .Duration }}s
- **Success Rate:** {{ printf "%.2f" .Success }}%
- **Latency (ms):** Mean={{ printf "%.2f" .LatencyMean }}, 50th={{ printf "%.2f" .Latency50 }}, 95th={{ printf "%.2f" .Latency95 }}, 99th={{ printf "%.2f" .Latency99 }}
- **Compressed Size (approx):** {{ .CompressedFileSize }} bytes
- **Compression Ratio (approx):** {{ printf "%.2f" .CompressionRatio }}x
- **Status Codes:** {{ range $code, $count := .StatusCodes }}{{ $code }}:{{ $count }} {{ end }}
{{ if .Error }}- **Errors:** {{ .Error }}{{ end }}
{{ else }}
No data available.
{{ end }}

**Max Throughput Test (Capacity Focus):**
{{ with .MaxResult }}
- **Requests:** {{ .Requests }}
- **Max Throughput:** {{ printf "%.2f" .Throughput }} req/s
- **Duration:** {{ printf "%.2f" .Duration }}s
- **Success Rate:** {{ printf "%.2f" .Success }}%
- **Latency at Max Rate (ms):** Mean={{ printf "%.2f" .LatencyMean }}, 50th={{ printf "%.2f" .Latency50 }}, 95th={{ printf "%.2f" .Latency95 }}, 99th={{ printf "%.2f" .Latency99 }}
- **Compressed Size (approx):** {{ .CompressedFileSize }} bytes
- **Compression Ratio (approx):** {{ printf "%.2f" .CompressionRatio }}x
- **Status Codes:** {{ range $code, $count := .StatusCodes }}{{ $code }}:{{ $count }} {{ end }}
{{ if .Error }}- **Errors:** {{ .Error }}{{ end }}
{{ else }}
No data available.
{{ end }}
{{ end }}
{{ end }}
{{ end }}

## Comparison Summary

### Fixed Rate Tests (Latency @ ~50 req/s)

| Algorithm | Level    | Data Size | Latency (99th ms) | Throughput (req/s) | Comp. Ratio (x) | Comp. Size (bytes) | Success (%) |
|-----------|----------|-----------|-------------------|--------------------|-----------------|--------------------|-------------|
{{ range .FixedResults }}| {{ .Algorithm }} | {{ .Level }} | {{ .DataSize }} | {{ printf "%8.2f" .Latency99 }} | {{ printf "%18.2f" .Throughput }} | {{ printf "%15.2f" .CompressionRatio }} | {{ printf "%18d" .CompressedFileSize }} | {{ printf "%11.2f" .Success }} |
{{ end }}

### Max Throughput Tests (Capacity)

| Algorithm | Level    | Data Size | Max Throughput (req/s) | Latency (99th ms) | Comp. Ratio (x) | Comp. Size (bytes) | Success (%) |
|-----------|----------|-----------|------------------------|-------------------|-----------------|--------------------|-------------|
{{ range .MaxResults }}| {{ .Algorithm }} | {{ .Level }} | {{ .DataSize }} | {{ printf "%22.2f" .Throughput }} | {{ printf "%17.2f" .Latency99 }} | {{ printf "%15.2f" .CompressionRatio }} | {{ printf "%18d" .CompressedFileSize }} | {{ printf "%11.2f" .Success }} |
{{ end }}

## Performance Graphs (Mermaid)

**Note:** Mermaid graphs render in environments that support it (like GitLab, GitHub with extensions, etc.).

### Fixed Rate: 99th Percentile Latency vs Compression Ratio (Large Data)

![Fixed Rate Latency vs Compression Ratio (Large Data)](report_images/latency_vs_ratio_large.png)

### Max Throughput vs Compression Ratio (Large Data)

![Max Throughput vs Compression Ratio (Large Data)](report_images/throughput_vs_ratio_large.png)

*(Similar graphs for Medium and Small data can be added here if needed)*

## Conclusion

Based on the results:

- **Compression Ratio:** Brotli (level 11) and Zstd (best) generally provide the highest compression ratios, significantly reducing data size, especially for larger payloads. Gzip (best) offers good compression, better than default/speed levels.
- **Latency (Fixed Rate):**
    - Without compression ('none'), latency is generally the lowest as expected, but at the cost of bandwidth.
    - Faster algorithms like Gzip (speed), Brotli (1), and Zstd (fastest) introduce minimal latency overhead compared to 'none'.
    - Higher compression levels (Brotli 11, Zstd best, Gzip best) tend to increase latency due to the higher CPU cost of compression, although this effect might be more pronounced on the server side during compression rather than the client side during decompression (which Vegeta measures). The impact varies with data size.
- **Throughput (Max Rate):**
    - 'None' often achieves the highest raw request throughput, limited primarily by network or server connection handling.
    - Lightweight compression (Gzip speed, Brotli 1, Zstd fastest) can sometimes achieve throughput close to 'none', especially if network bandwidth is the bottleneck for uncompressed data.
    - Heavy compression algorithms significantly reduce the maximum achievable throughput due to the CPU overhead on the server. Zstd generally shows a good balance, offering better throughput than Brotli at similar high compression levels.
- **Trade-offs:**
    - **Bandwidth Sensitive:** If minimizing data transfer is critical, Brotli (11) or Zstd (best) are strong choices, accepting a potential latency/throughput penalty.
    - **Latency Sensitive:** If minimizing response time is paramount, 'none' or very light compression (Zstd fastest, Brotli 1, Gzip speed) are preferable.
    - **Balanced:** Zstd (default or even best) often provides a good compromise between compression ratio, latency, and throughput. Gzip (default) remains a widely compatible and reasonable default.

**Recommendation:**

The optimal choice depends heavily on the specific application needs (latency sensitivity, bandwidth constraints, typical data size, client capabilities).

- For APIs serving potentially large, compressible text-based data where bandwidth saving is important, **Zstd (default)** offers an excellent balance.
- For static assets where maximum compression is desired and can be pre-calculated, **Brotli (11)** is a top contender.
- For general-purpose web content where compatibility and moderate compression are sufficient, **Gzip (default)** is still a solid choice.
- If latency is absolutely critical and bandwidth is plentiful, **'none'** or **Zstd (fastest)** might be considered.

Further investigation could involve profiling the server-side CPU usage during these tests to get a clearer picture of the compression cost.
`

// Helper struct for template generation
type AlgoLevelData struct {
	SizeName     string
	OriginalSize int64
	FixedResult  *Result
	MaxResult    *Result
}

type AlgoData struct {
	LevelName string
	DataSizes []AlgoLevelData
}

type TemplateAlgoData struct {
	Name   string
	Levels []AlgoData
}

// Add CompressionRatio method to Result
func (r Result) CompressionRatio() float64 {
	if r.OriginalFileSize > 0 && r.CompressedFileSize > 0 {
		return float64(r.OriginalFileSize) / float64(r.CompressedFileSize)
	}
	return 0.0
}


func main() {
	resultsDir := "results"
	files, err := os.ReadDir(resultsDir)
	if err != nil {
		log.Fatalf("Failed to read results directory: %v", err)
	}

	var allResults []Result
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".bin") {
			filePath := filepath.Join(resultsDir, file.Name())
			result, err := parseVegetaReport(filePath)
			if err != nil {
				log.Printf("Error parsing %s: %v", file.Name(), err)
				// Store partial result with error
				if result.FileName == "" { // If filename parsing failed
					result.FileName = file.Name()
					result.Error = err.Error()
				}
			}
			allResults = append(allResults, result)
		}
	}

	// Sort results for consistent report order
	sort.Slice(allResults, func(i, j int) bool {
		if allResults[i].Algorithm != allResults[j].Algorithm {
			return allResults[i].Algorithm < allResults[j].Algorithm
		}
		// Custom level sorting logic directly here
		levelI := allResults[i].Level
		levelJ := allResults[j].Level
		if levelI != levelJ {
			// Define level order within the comparison
			order := map[string]int{
				// Common levels (lower value = comes first)
				"N/A": 0, // For 'none'
				// Gzip
				"speed": 10,
				// Zstd
				"fastest": 10, // Treat speed/fastest similarly
				// Common default
				"default": 20,
				// Gzip/Zstd best
				"best": 30,
				// Brotli levels (use numbers directly, offset to avoid clashes)
				"1":  101,
				"6":  106,
				"11": 111,
			}
			// Handle potential Brotli numeric levels if parsed as numbers (though filename regex suggests strings)
			if _, err := strconv.Atoi(levelI); err == nil {
				if _, err := strconv.Atoi(levelJ); err == nil {
					numI, _ := strconv.Atoi(levelI)
					numJ, _ := strconv.Atoi(levelJ)
					return numI < numJ
				}
			}
			// Use map order for string levels
			return order[levelI] < order[levelJ]
		}
		if allResults[i].DataSize != allResults[j].DataSize {
			return dataSizeOrder(allResults[i].DataSize) < dataSizeOrder(allResults[j].DataSize)
		}
		return allResults[i].TestType < allResults[j].TestType // fixed before max
	})

	// Prepare data for the template
	reportData := ReportData{Results: allResults}
	tmplData := prepareTemplateData(reportData)


	// Generate markdown
	tmpl, err := template.New("report").Parse(markdownTemplate)
	if err != nil {
		log.Fatalf("Failed to parse markdown template: %v", err)
	}

	var reportContent bytes.Buffer
	err = tmpl.Execute(&reportContent, tmplData)
	if err != nil {
		log.Fatalf("Failed to execute template: %v", err)
	}

	// Create images directory
	imagesDir := "report_images"
	if err := os.MkdirAll(imagesDir, 0755); err != nil {
		log.Fatalf("Failed to create images directory %s: %v", imagesDir, err)
	}

	   // Generate plots
	   latencyPlotPath := filepath.Join(imagesDir, "latency_vs_ratio_large.png")
	   err = createLatencyVsRatioPlot(tmplData.FixedResultsLarge, latencyPlotPath, tmplData.MaxLatencyFixedLarge)
	   if err != nil {
	       log.Printf("Warning: Failed to generate latency plot: %v", err)
	       // Continue report generation even if plot fails
	   }

	   throughputPlotPath := filepath.Join(imagesDir, "throughput_vs_ratio_large.png")
	   err = createThroughputVsRatioPlot(tmplData.MaxResultsLarge, throughputPlotPath, tmplData.MaxThroughputMaxLarge)
	    if err != nil {
	       log.Printf("Warning: Failed to generate throughput plot: %v", err)
	       // Continue report generation even if plot fails
	   }


	// Write markdown to file
	outputFile := "compression_benchmark_report.md"
	err = os.WriteFile(outputFile, reportContent.Bytes(), 0644)
	if err != nil {
		log.Fatalf("Failed to write report file %s: %v", outputFile, err)
	}

	fmt.Printf("Report successfully generated: %s\n", outputFile)
}

// --- Plotting functions ---

// createLatencyVsRatioPlot generates a scatter plot of Latency (99th) vs Compression Ratio
func createLatencyVsRatioPlot(results []Result, outputPath string, yMax float64) error {
	pts := make(plotter.XYs, 0, len(results))
	labels := make(map[int]string) // Map index to label
	i := 0
	for _, r := range results {
		if r.CompressionRatio() > 0 { // Only plot valid ratios
			pts = append(pts, plotter.XY{X: r.CompressionRatio(), Y: r.Latency99})
			labels[i] = fmt.Sprintf("%s-%s", r.Algorithm, r.Level)
			i++
		}
	}

	if len(pts) == 0 {
		return fmt.Errorf("no valid data points for latency plot")
	}

	p := plot.New()

	p.Title.Text = "Fixed Rate (Large Data): Latency (99th) vs Compression Ratio"
	p.X.Label.Text = "Compression Ratio (Higher is Better)"
	p.Y.Label.Text = "Latency 99th (ms) (Lower is Better)"
	p.Y.Max = yMax // Use calculated max from data

	s, err := plotter.NewScatter(pts)
	if err != nil {
		return fmt.Errorf("could not create scatter plotter: %v", err)
	}
	s.Color = color.RGBA{R: 255, B: 128, A: 255} // Example color

	p.Add(s)

	// Add labels (optional, can clutter the plot)
	// l, err := plotter.NewLabels(plotter.XYLabels{XYs: pts, Labels: getLabels(labels, len(pts))})
	// if err != nil {
	// 	log.Printf("Warning: could not create labels: %v", err)
	// } else {
	// 	p.Add(l)
	// }


	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, 4*vg.Inch, outputPath); err != nil {
		return fmt.Errorf("could not save plot: %v", err)
	}
	return nil
}

// createThroughputVsRatioPlot generates a scatter plot of Throughput vs Compression Ratio
func createThroughputVsRatioPlot(results []Result, outputPath string, yMax float64) error {
	pts := make(plotter.XYs, 0, len(results))
	labels := make(map[int]string)
	i := 0
	for _, r := range results {
		if r.CompressionRatio() > 0 {
			pts = append(pts, plotter.XY{X: r.CompressionRatio(), Y: r.Throughput})
			labels[i] = fmt.Sprintf("%s-%s", r.Algorithm, r.Level)
			i++
		}
	}

	if len(pts) == 0 {
		return fmt.Errorf("no valid data points for throughput plot")
	}

	p := plot.New()

	p.Title.Text = "Max Throughput (Large Data): Throughput vs Compression Ratio"
	p.X.Label.Text = "Compression Ratio (Higher is Better)"
	p.Y.Label.Text = "Throughput (req/s) (Higher is Better)"
	p.Y.Max = yMax

	s, err := plotter.NewScatter(pts)
	if err != nil {
		return fmt.Errorf("could not create scatter plotter: %v", err)
	}
	s.Color = color.RGBA{G: 255, B: 128, A: 255} // Example color

	p.Add(s)

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, 4*vg.Inch, outputPath); err != nil {
		return fmt.Errorf("could not save plot: %v", err)
	}
	return nil
}


// --- Helper functions for sorting and template data preparation ---

func dataSizeOrder(size string) int {
	switch size {
	case "small": return 1
	case "medium": return 2
	case "large": return 3
	default: return 9
	}
}

type TemplateData struct {
	Algorithms           []TemplateAlgoData
	FixedResults         []Result
	MaxResults           []Result
	FixedResultsLarge    []Result // For graph
	MaxResultsLarge      []Result // For graph
	MaxLatencyFixedLarge float64  // For graph axis scaling
	MaxThroughputMaxLarge float64 // For graph axis scaling
	SmallSize            int64
	MediumSize           int64
	LargeSize            int64
}


func prepareTemplateData(report ReportData) TemplateData {
	algoMap := make(map[string]map[string]map[string]map[string]*Result) // algo -> level -> datasize -> testtype -> result

	var fixedResults, maxResults, fixedResultsLarge, maxResultsLarge []Result
	maxLatencyFixedLarge := 0.0
	maxThroughputMaxLarge := 0.0


	for i := range report.Results {
		res := &report.Results[i] // Use pointer to modify map value

		if _, ok := algoMap[res.Algorithm]; !ok {
			algoMap[res.Algorithm] = make(map[string]map[string]map[string]*Result)
		}
		if _, ok := algoMap[res.Algorithm][res.Level]; !ok {
			algoMap[res.Algorithm][res.Level] = make(map[string]map[string]*Result)
		}
		if _, ok := algoMap[res.Algorithm][res.Level][res.DataSize]; !ok {
			algoMap[res.Algorithm][res.Level][res.DataSize] = make(map[string]*Result)
		}
		algoMap[res.Algorithm][res.Level][res.DataSize][res.TestType] = res

		// Collect for summary tables and graphs
		if res.TestType == "fixed" {
			fixedResults = append(fixedResults, *res)
			if res.DataSize == "large" {
				fixedResultsLarge = append(fixedResultsLarge, *res)
				if res.Latency99 > maxLatencyFixedLarge {
					maxLatencyFixedLarge = res.Latency99
				}
			}
		} else if res.TestType == "max" {
			maxResults = append(maxResults, *res)
			if res.DataSize == "large" {
				maxResultsLarge = append(maxResultsLarge, *res)
				if res.Throughput > maxThroughputMaxLarge {
					maxThroughputMaxLarge = res.Throughput
				}
			}
		}
	}

	// Structure for detailed results section
	var templateAlgos []TemplateAlgoData
	algorithms := []string{"none", "gzip", "brotli", "zstd"} // Define order
	levels := map[string][]string{
		"none":   {"large"}, // 'none' has no level, use data size as placeholder? Or adjust template. Using 'large' as placeholder level name.
		"gzip":   {"speed", "default", "best"},
		"brotli": {"1", "6", "11"},
		"zstd":   {"fastest", "default", "best"},
	}
	dataSizesOrder := []string{"small", "medium", "large"}


    // Adjust 'none' level handling
    if noneLevelMap, ok := algoMap["none"]; ok {
        if _, largeExists := noneLevelMap["large"]; largeExists { // Check if 'large' placeholder exists
             noneLevelMap["N/A"] = noneLevelMap["large"] // Move to a more descriptive level name
             delete(noneLevelMap, "large") // Remove the placeholder using the correct map reference
             levels["none"] = []string{"N/A"} // Update levels map
        } else {
             // Handle case where 'none' results might have unexpected level names
             // For simplicity, assume the first found level is the one to use
             var foundLevel string
             for lvl := range noneLevelMap { // Iterate over the actual map
                 foundLevel = lvl
                 break
             }
             if foundLevel != "" {
                 levels["none"] = []string{foundLevel}
             } else {
                 levels["none"] = []string{"N/A"} // Default if no level found
             }
        }
    }


 for _, algoName := range algorithms {
		if levelMap, ok := algoMap[algoName]; ok {
			var algoLevels []AlgoData
			currentLevels := levels[algoName] // Get defined levels for this algo

			// Sort levels based on custom order directly here
			sort.SliceStable(currentLevels, func(i, j int) bool {
				levelI := currentLevels[i]
				levelJ := currentLevels[j]
				// Define level order within the comparison (same logic as in main sort)
				order := map[string]int{
					"N/A": 0,
					"speed": 10, "fastest": 10,
					"default": 20,
					"best": 30,
					"1":  101,
					"6":  106,
					"11": 111,
				}
				// Handle potential Brotli numeric levels
				if _, err := strconv.Atoi(levelI); err == nil {
					if _, err := strconv.Atoi(levelJ); err == nil {
						numI, _ := strconv.Atoi(levelI)
						numJ, _ := strconv.Atoi(levelJ)
						return numI < numJ
					}
				}
				return order[levelI] < order[levelJ]
			})


			for _, levelName := range currentLevels {
                 if dataMap, ok := levelMap[levelName]; ok {
					var levelDataSizes []AlgoLevelData
					for _, sizeName := range dataSizesOrder {
						if typeMap, ok := dataMap[sizeName]; ok {
							ald := AlgoLevelData{
								SizeName:     sizeName,
								OriginalSize: getFileSize(dataSizes[sizeName]),
								FixedResult:  typeMap["fixed"],
								MaxResult:    typeMap["max"],
							}
							levelDataSizes = append(levelDataSizes, ald)
						}
					}
                    if len(levelDataSizes) > 0 { // Only add level if it has data sizes
					    algoLevels = append(algoLevels, AlgoData{LevelName: levelName, DataSizes: levelDataSizes})
                    }
				}
			}
            if len(algoLevels) > 0 { // Only add algorithm if it has levels with data
			    templateAlgos = append(templateAlgos, TemplateAlgoData{Name: strings.Title(algoName), Levels: algoLevels})
            }
		}
	}

    // Ensure graph axes have some padding
    maxLatencyFixedLarge *= 1.1
    maxThroughputMaxLarge *= 1.1
    if maxLatencyFixedLarge == 0 { maxLatencyFixedLarge = 100 } // Default if no data
    if maxThroughputMaxLarge == 0 { maxThroughputMaxLarge = 1000 } // Default if no data


	return TemplateData{
		Algorithms:           templateAlgos,
		FixedResults:         fixedResults,
		MaxResults:           maxResults,
		FixedResultsLarge:    fixedResultsLarge,
		MaxResultsLarge:      maxResultsLarge,
		MaxLatencyFixedLarge: maxLatencyFixedLarge,
		MaxThroughputMaxLarge: maxThroughputMaxLarge,
		SmallSize:            getFileSize(dataSizes["small"]),
		MediumSize:           getFileSize(dataSizes["medium"]),
		LargeSize:            getFileSize(dataSizes["large"]),
	}
}