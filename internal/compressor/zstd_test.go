package compressor_test

import (
	"bytes"
	"go-compress-showdown/internal/compressor" // compressor パッケージをインポート
	"io"
	"testing"

	"github.com/klauspost/compress/zstd" // zstd パッケージをインポート
)

// benchmarkZstdCompress は指定されたレベルで圧縮ベンチマークを実行します
func benchmarkZstdCompress(b *testing.B, level zstd.EncoderLevel, levelStr string) {
	b.Helper()
	b.ReportAllocs()
	b.SetBytes(int64(len(testData))) // testData は gzip_test.go で定義済み

	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		// compressor パッケージの関数を使用
		writer, err := compressor.GetZstdWriter(&buf, levelStr)
		if err != nil {
			b.Fatalf("Failed to get zstd writer (level: %s): %v", levelStr, err)
		}

		_, err = writer.Write(testData)
		if err != nil {
			writer.Close()
			b.Fatalf("Failed to write data (level: %s): %v", levelStr, err)
		}
		err = writer.Close()
		if err != nil {
			b.Fatalf("Failed to close writer (level: %s): %v", levelStr, err)
		}
	}
}

// 各圧縮レベルでのベンチマーク関数
func BenchmarkZstdCompressFastest(b *testing.B) {
	benchmarkZstdCompress(b, zstd.SpeedFastest, "fastest")
}
func BenchmarkZstdCompressDefault(b *testing.B) {
	benchmarkZstdCompress(b, zstd.SpeedDefault, "default")
}
func BenchmarkZstdCompressBest(b *testing.B) {
	benchmarkZstdCompress(b, zstd.SpeedBestCompression, "best")
}

// benchmarkZstdDecompress は解凍ベンチマークを実行します
func benchmarkZstdDecompress(b *testing.B, level zstd.EncoderLevel, levelStr string) {
	b.Helper()
	b.ReportAllocs()
	b.SetBytes(int64(len(testData)))

	// 事前にテストデータを圧縮しておく
	var compressedBuf bytes.Buffer
	writer, _ := compressor.GetZstdWriter(&compressedBuf, levelStr)
	writer.Write(testData)
	writer.Close()
	compressedData := compressedBuf.Bytes()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// zstd パッケージの Decoder を使用
		reader, err := zstd.NewReader(bytes.NewReader(compressedData))
		if err != nil {
			b.Fatalf("Failed to create zstd reader (level: %s): %v", levelStr, err)
		}

		_, err = io.Copy(io.Discard, reader)
		if err != nil {
			reader.Close() // zstd.Decoder は Close() を持つ
			b.Fatalf("Failed to decompress data (level: %s): %v", levelStr, err)
		}
		reader.Close()
	}
}

// 各圧縮レベルで圧縮されたデータの解凍ベンチマーク関数
func BenchmarkZstdDecompressFastest(b *testing.B) {
	benchmarkZstdDecompress(b, zstd.SpeedFastest, "fastest")
}
func BenchmarkZstdDecompressDefault(b *testing.B) {
	benchmarkZstdDecompress(b, zstd.SpeedDefault, "default")
}
func BenchmarkZstdDecompressBest(b *testing.B) {
	benchmarkZstdDecompress(b, zstd.SpeedBestCompression, "best")
}

// --- Helper to avoid redefining testData ---
// testData は gzip_test.go で定義済み