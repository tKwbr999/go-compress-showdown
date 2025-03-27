package compressor_test

import (
	"bytes"
	"go-compress-showdown/internal/compressor" // compressor パッケージをインポート
	"io"
	"strconv"
	"testing"

	"github.com/andybalholm/brotli" // brotli パッケージをインポート
)

// benchmarkBrotliCompress は指定された品質レベルで圧縮ベンチマークを実行します
func benchmarkBrotliCompress(b *testing.B, quality int, qualityStr string) {
	b.Helper()
	b.ReportAllocs()
	b.SetBytes(int64(len(testData))) // testData は gzip_test.go で定義済み

	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		// compressor パッケージの関数を使用
		writer, err := compressor.GetBrotliWriter(&buf, qualityStr)
		if err != nil {
			b.Fatalf("Failed to get brotli writer (quality: %s): %v", qualityStr, err)
		}

		_, err = writer.Write(testData)
		if err != nil {
			writer.Close()
			b.Fatalf("Failed to write data (quality: %s): %v", qualityStr, err)
		}
		err = writer.Close()
		if err != nil {
			b.Fatalf("Failed to close writer (quality: %s): %v", qualityStr, err)
		}
	}
}

// 代表的な品質レベルでのベンチマーク関数
func BenchmarkBrotliCompressQ1(b *testing.B) {
	benchmarkBrotliCompress(b, 1, "1")
}
func BenchmarkBrotliCompressQ6(b *testing.B) { // デフォルト相当
	benchmarkBrotliCompress(b, 6, "6")
}
func BenchmarkBrotliCompressQ11(b *testing.B) { // 最高圧縮
	benchmarkBrotliCompress(b, 11, "11")
}

// benchmarkBrotliDecompress は解凍ベンチマークを実行します
func benchmarkBrotliDecompress(b *testing.B, quality int, qualityStr string) {
	b.Helper()
	b.ReportAllocs()
	b.SetBytes(int64(len(testData)))

	// 事前にテストデータを圧縮しておく
	var compressedBuf bytes.Buffer
	writer, _ := compressor.GetBrotliWriter(&compressedBuf, qualityStr)
	writer.Write(testData)
	writer.Close()
	compressedData := compressedBuf.Bytes()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// brotli パッケージの Reader を使用
		reader := brotli.NewReader(bytes.NewReader(compressedData))

		_, err := io.Copy(io.Discard, reader)
		if err != nil {
			// brotli.Reader は Close() メソッドを持たない
			b.Fatalf("Failed to decompress data (quality: %s): %v", qualityStr, err)
		}
	}
}

// 各品質レベルで圧縮されたデータの解凍ベンチマーク関数
func BenchmarkBrotliDecompressQ1(b *testing.B) {
	benchmarkBrotliDecompress(b, 1, "1")
}
func BenchmarkBrotliDecompressQ6(b *testing.B) {
	benchmarkBrotliDecompress(b, 6, "6")
}
func BenchmarkBrotliDecompressQ11(b *testing.B) {
	benchmarkBrotliDecompress(b, 11, "11")
}

// --- Helper to avoid redefining testData ---
// generateTestData は gzip_test.go で定義されているため、ここでは再定義しない
// testData も同様

// strconv を使うためインポートを追加
var _ = strconv.Itoa // Use strconv to avoid "imported and not used" error if benchmarks are commented out