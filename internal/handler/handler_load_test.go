package handler_test

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

const (
	targetHost               = "http://localhost:8080" // テスト対象ホスト
	targetsDir               = "../../vegeta_targets"  // ターゲットファイルディレクトリ (テストファイルからの相対パス)
	resultsDir               = "../../results"         // 結果出力ディレクトリ (テストファイルからの相対パス)
	fixedRate                = 50                      // 固定レート (リクエスト/秒)
	fixedDuration            = 15 * time.Second        // 固定レートテスト時間
	maxDuration              = 15 * time.Second        // 最大スループットテスト時間
	maxWorkers        uint64 = 50                      // 最大ワーカー数
	serverWaitTimeout        = 60 * time.Second        // サーバー起動待機タイムアウト
)

// setup ensures the results directory exists and the server is running.
func setup(t *testing.T) {
	t.Helper()
	// 結果ディレクトリを作成
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		t.Fatalf("Failed to create results directory %s: %v", resultsDir, err)
	}

	// サーバーが起動するまで待機 (簡易ヘルスチェック)
	startTime := time.Now()
	for {
		resp, err := http.Get(targetHost + "/")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			t.Logf("Server is up and running at %s", targetHost)
			break
		}
		if time.Since(startTime) > serverWaitTimeout {
			t.Fatalf("Server did not start within %v", serverWaitTimeout)
		}
		t.Log("Waiting for server to start...")
		time.Sleep(2 * time.Second)
	}
}

// runVegetaTest executes a single vegeta test (fixed rate or max throughput).
func runVegetaTest(t *testing.T, targetName, testType string, attacker *vegeta.Attacker, rate vegeta.Rate, duration time.Duration, targeter vegeta.Targeter) {
	t.Helper()
	reportFileName := filepath.Join(resultsDir, fmt.Sprintf("vegeta_report_%s_%s.txt", targetName, testType))
	binFileName := filepath.Join(resultsDir, fmt.Sprintf("vegeta_results_%s_%s.bin", targetName, testType))

	t.Logf("[%s] Running %s test (Rate: %v, Duration: %v)...", targetName, testType, rate, duration)

	var metrics vegeta.Metrics
	binFile, err := os.Create(binFileName)
	if err != nil {
		t.Errorf("[%s] Failed to create bin file %s: %v", targetName, binFileName, err)
		// binファイルが作れなくてもレポートは試みる
	} else {
		defer binFile.Close()
	}

	// Attack実行、結果をbinファイルとmetricsに格納
	var results <-chan *vegeta.Result
	if binFile != nil {
		results = attacker.Attack(targeter, rate, duration, targetName+"_"+testType)
	} else {
		// binファイルが作成できなかった場合、ファイル書き込みなしで実行
		for res := range attacker.Attack(targeter, rate, duration, targetName+"_"+testType) {
			metrics.Add(res)
		}
		metrics.Close()
	}

	// 結果をエンコードしてファイルに書き込む (binFileがnilでない場合)
	successCount := uint64(0)
	totalCount := uint64(0)
	if binFile != nil {
		enc := vegeta.NewEncoder(binFile)
		for result := range results {
			totalCount++
			if result.Code >= 200 && result.Code < 300 {
				successCount++
			}
			metrics.Add(result)
			if err := enc.Encode(result); err != nil {
				t.Errorf("[%s] Failed to encode result to %s: %v", targetName, binFileName, err)
				// エンコードエラーが発生しても処理は続行
			}
		}
		metrics.Close()
	} else {
		// binファイルがない場合はmetricsからカウント (概算になる可能性あり)
		successCount = uint64(metrics.StatusCodes["200"]) // 簡易的に200のみカウント
		totalCount = metrics.Requests
	}

	t.Logf("[%s] %s test completed. Total requests: %d, Success: %d", targetName, testType, totalCount, successCount)

	// レポート生成
	t.Logf("[%s] Generating %s report...", targetName, testType)
	reportFile, err := os.Create(reportFileName)
	if err != nil {
		t.Errorf("[%s] Failed to create report file %s: %v", targetName, reportFileName, err)
		return
	}
	defer reportFile.Close()

	reporter := vegeta.NewTextReporter(&metrics) // テキストレポーターを使用
	if err := reporter.Report(reportFile); err != nil {
		t.Errorf("[%s] Failed to write report to %s: %v", targetName, reportFileName, err)
	}

	// ターミナルにも簡易レポートを出力 (エラーが見やすいように)
	var reportBuf bytes.Buffer
	reporterStdOut := vegeta.NewTextReporter(&metrics)
	if err := reporterStdOut.Report(&reportBuf); err != nil {
		t.Errorf("[%s] Failed to generate stdout report: %v", targetName, err)
	} else {
		t.Logf("[%s] %s Report:\n%s", targetName, testType, reportBuf.String())
	}
}

// newTargeterFromVegetaFile reads a vegeta target file and creates a Targeter.
// It expects the format: METHOD URL (e.g., GET http://placeholder/path)
// It replaces "http://localhost:8080" with the actual targetHost.
func newTargeterFromVegetaFile(filePath, host string) (vegeta.Targeter, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open target file %s: %w", filePath, err)
	}
	defer file.Close()

	var targets []vegeta.Target
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") { // Skip empty lines and comments
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid format in %s:%d: '%s'", filePath, lineNumber, line)
		}

		method := parts[0]
		url := parts[1]

		// Replace placeholder host with the actual target host
		url = strings.Replace(url, "http://localhost:8080", host, 1)

		target := vegeta.Target{
			Method: method,
			URL:    url,
		}
		// Add header or body if specified in the file format (needs adjustment)
		// Example: if len(parts) > 2 && strings.HasPrefix(parts[2], "Header:") ...
		// Example: if len(parts) > 2 && strings.HasPrefix(parts[2], "Body:") ...

		targets = append(targets, target)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read target file %s: %w", filePath, err)
	}

	if len(targets) == 0 {
		return nil, fmt.Errorf("no valid targets found in %s", filePath)
	}

	return vegeta.NewStaticTargeter(targets...), nil
}

// TestLoad performs load testing for all targets defined in the targets directory.
func TestLoad(t *testing.T) {
	// Skip load tests in short mode (-short flag)
	if testing.Short() {
		t.Skip("Skipping load tests in short mode")
	}

	setup(t) // Setup results dir and wait for server

	files, err := os.ReadDir(targetsDir)
	if err != nil {
		t.Fatalf("Failed to read targets directory %s: %v", targetsDir, err)
	}

	// Create attackers once
	attackerFixed := vegeta.NewAttacker()
	attackerMax := vegeta.NewAttacker(vegeta.Workers(maxWorkers))

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".txt") {
			continue
		}

		targetFileName := file.Name()
		targetName := strings.TrimSuffix(targetFileName, ".txt")
		targetFilePath := filepath.Join(targetsDir, targetFileName)

		// サブテストとして各ターゲットを実行
		t.Run(targetName, func(t *testing.T) {
			targeter, err := newTargeterFromVegetaFile(targetFilePath, targetHost)
			if err != nil {
				t.Fatalf("Failed to create targeter for %s: %v", targetFilePath, err)
			}

			// --- 固定レートテスト ---
			runVegetaTest(t, targetName, "fixed", attackerFixed, vegeta.Rate{Freq: fixedRate, Per: time.Second}, fixedDuration, targeter)

			time.Sleep(1 * time.Second) // サーバー負荷軽減

			// --- 最大スループットテスト ---
			runVegetaTest(t, targetName, "max", attackerMax, vegeta.Rate{Freq: 0}, maxDuration, targeter)

			time.Sleep(1 * time.Second) // サーバー負荷軽減
		})
	}
}
