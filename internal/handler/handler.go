package handler

import (
	"fmt"
	"go-compress-showdown/internal/compressor"
	"go-compress-showdown/internal/data" // data パッケージをインポート
	"log"
	"net/http"
)

// NoneHandler は圧縮なしでデータを返します。
func NoneHandler(w http.ResponseWriter, r *http.Request) {
	dataKey := r.URL.Query().Get("data")
	if dataKey == "" {
		http.Error(w, "Query parameter 'data' is required (e.g., ?data=small)", http.StatusBadRequest)
		return
	}

	// data パッケージからデータを読み込む
	dataBytes, err := data.LoadData(dataKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load data for key '%s': %v", dataKey, err), http.StatusNotFound) // データが見つからない場合は 404
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(dataBytes)
}

// GzipHandler はgzip圧縮してデータを返します。
func GzipHandler(w http.ResponseWriter, r *http.Request) {
	dataKey := r.URL.Query().Get("data")
	levelStr := r.URL.Query().Get("level")

	if dataKey == "" {
		http.Error(w, "Query parameter 'data' is required (e.g., ?data=small)", http.StatusBadRequest)
		return
	}

	// data パッケージからデータを読み込む
	originalData, err := data.LoadData(dataKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load data for key '%s': %v", dataKey, err), http.StatusNotFound)
		return
	}

	// Content-Encoding ヘッダーを設定
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", "text/plain")

	// gzip Writer を取得
	gzipWriter, err := compressor.GetGzipWriter(w, levelStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to initialize gzip writer: %v", err), http.StatusInternalServerError)
		return
	}
	defer gzipWriter.Close()

	// データを書き込み (圧縮される)
	_, err = gzipWriter.Write(originalData)
	if err != nil {
		log.Printf("Error writing gzipped data: %v", err)
	}

	// Closeが呼ばれることで最終的な圧縮データがフラッシュされる
}


// BrotliHandler はbrotli圧縮してデータを返します。
func BrotliHandler(w http.ResponseWriter, r *http.Request) {
	dataKey := r.URL.Query().Get("data")
	qualityStr := r.URL.Query().Get("level") // Brotliでは quality と呼ばれることが多い

	if dataKey == "" {
		http.Error(w, "Query parameter 'data' is required (e.g., ?data=small)", http.StatusBadRequest)
		return
	}

	// data パッケージからデータを読み込む
	originalData, err := data.LoadData(dataKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load data for key '%s': %v", dataKey, err), http.StatusNotFound)
		return
	}

	// Content-Encoding ヘッダーを設定
	w.Header().Set("Content-Encoding", "br")
	w.Header().Set("Content-Type", "text/plain")

	// brotli Writer を取得
	brotliWriter, err := compressor.GetBrotliWriter(w, qualityStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to initialize brotli writer: %v", err), http.StatusInternalServerError)
		return
	}
	defer brotliWriter.Close()

	// データを書き込み (圧縮される)
	_, err = brotliWriter.Write(originalData)
	if err != nil {
		log.Printf("Error writing brotli data: %v", err)
	}

	// Closeが呼ばれることで最終的な圧縮データがフラッシュされる
}

// ZstdHandler はzstd圧縮してデータを返します。
func ZstdHandler(w http.ResponseWriter, r *http.Request) {
	dataKey := r.URL.Query().Get("data")
	levelStr := r.URL.Query().Get("level")

	if dataKey == "" {
		http.Error(w, "Query parameter 'data' is required (e.g., ?data=small)", http.StatusBadRequest)
		return
	}

	// data パッケージからデータを読み込む
	originalData, err := data.LoadData(dataKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load data for key '%s': %v", dataKey, err), http.StatusNotFound)
		return
	}

	// Content-Encoding ヘッダーを設定
	w.Header().Set("Content-Encoding", "zstd")
	w.Header().Set("Content-Type", "text/plain")

	// zstd Writer を取得
	zstdWriter, err := compressor.GetZstdWriter(w, levelStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to initialize zstd writer: %v", err), http.StatusInternalServerError)
		return
	}
	defer zstdWriter.Close()

	// データを書き込み (圧縮される)
	_, err = zstdWriter.Write(originalData)
	if err != nil {
		log.Printf("Error writing zstd data: %v", err)
	}

	// Closeが呼ばれることで最終的な圧縮データがフラッシュされる
}
