package compressor_test

import (
	"bytes"
	"compress/gzip"
	"go-compress-showdown/internal/compressor" // compressor パッケージをインポート
	"io"
	"testing"
)

// ベンチマーク用のテストデータ (ある程度のサイズを持たせる)
var testData = []byte(generateTestData(1024 * 10)) // 10KBのテストデータ

// generateTestData は指定されたサイズの繰り返し文字列データを生成します
func generateTestData(size int) string {
	base := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	data := ""
	for len(data) < size {
		data += base
	}
	return data[:size]
}

// benchmarkGzipCompress は指定されたレベルで圧縮ベンチマークを実行します
func benchmarkGzipCompress(b *testing.B, level int, levelStr string) {
	b.Helper()
	b.ReportAllocs() // メモリアロケーションをレポート
	b.SetBytes(int64(len(testData))) // 処理バイト数を設定

	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		// compressor パッケージの関数を使用
		writer, err := compressor.GetGzipWriter(&buf, levelStr)
		if err != nil {
			b.Fatalf("Failed to get gzip writer (level: %s): %v", levelStr, err)
		}

		_, err = writer.Write(testData)
		if err != nil {
			writer.Close() // エラー時もCloseを試みる
			b.Fatalf("Failed to write data (level: %s): %v", levelStr, err)
		}
		err = writer.Close() // Closeで圧縮が完了する
		if err != nil {
			b.Fatalf("Failed to close writer (level: %s): %v", levelStr, err)
		}
		// buf.Bytes() を使って圧縮結果にアクセスできるが、ベンチマークでは不要
	}
}

// 各圧縮レベルでのベンチマーク関数
func BenchmarkGzipCompressSpeed(b *testing.B) {
	benchmarkGzipCompress(b, gzip.BestSpeed, "speed")
}
func BenchmarkGzipCompressDefault(b *testing.B) {
	benchmarkGzipCompress(b, gzip.DefaultCompression, "default")
}
func BenchmarkGzipCompressBest(b *testing.B) {
	benchmarkGzipCompress(b, gzip.BestCompression, "best")
}

// benchmarkGzipDecompress は解凍ベンチマークを実行します
func benchmarkGzipDecompress(b *testing.B, level int, levelStr string) {
	b.Helper()
	b.ReportAllocs()
	b.SetBytes(int64(len(testData)))

	// 事前にテストデータを圧縮しておく
	var compressedBuf bytes.Buffer
	writer, _ := compressor.GetGzipWriter(&compressedBuf, levelStr)
	writer.Write(testData)
	writer.Close()
	compressedData := compressedBuf.Bytes()

	b.ResetTimer() // 圧縮にかかった時間をリセット

	for i := 0; i < b.N; i++ {
		reader, err := gzip.NewReader(bytes.NewReader(compressedData))
		if err != nil {
			b.Fatalf("Failed to create gzip reader (level: %s): %v", levelStr, err)
		}
		// 解凍結果を読み捨てる
		_, err = io.Copy(io.Discard, reader)
		if err != nil {
			reader.Close()
			b.Fatalf("Failed to decompress data (level: %s): %v", levelStr, err)
		}
		err = reader.Close()
		if err != nil {
			b.Fatalf("Failed to close reader (level: %s): %v", levelStr, err)
		}
	}
}

// 各圧縮レベルで圧縮されたデータの解凍ベンチマーク関数
func BenchmarkGzipDecompressSpeed(b *testing.B) {
	benchmarkGzipDecompress(b, gzip.BestSpeed, "speed")
}
func BenchmarkGzipDecompressDefault(b *testing.B) {
	benchmarkGzipDecompress(b, gzip.DefaultCompression, "default")
}
func BenchmarkGzipDecompressBest(b *testing.B) {
	benchmarkGzipDecompress(b, gzip.BestCompression, "best")
}