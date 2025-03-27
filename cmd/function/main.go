package main

import (
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"go-compress-showdown/internal/data"
	"go-compress-showdown/internal/handler"
)

func main() {
	// テストデータを事前に読み込み、キャッシュを初期化する
	// Cloud Functions のインスタンス起動時に一度だけ実行される想定
	// GetDataKeys() を呼び出すことで内部的に loadAndCache() が実行される
	_ = data.GetDataKeys() // 戻り値は不要なので破棄
	// キャッシュ初期化時にエラーが発生したか確認 (LoadData でも可)
	_, err := data.LoadData("small") // ダミーで small を読み込み、エラーを確認
	if err != nil {
		// 起動時にデータロードに失敗したら Fatal で終了させる
		log.Fatalf("Failed to initialize test data cache: %v", err)
	}

	// HTTP ルーターを作成
	mux := http.NewServeMux()

	// 各エンドポイントにハンドラーを登録
	mux.HandleFunc("/none", handler.NoneHandler)
	mux.HandleFunc("/gzip", handler.GzipHandler)
	mux.HandleFunc("/brotli", handler.BrotliHandler)
	mux.HandleFunc("/zstd", handler.ZstdHandler)

	// Cloud Functions フレームワークに HTTP 関数としてルーターを登録
	// これにより、すべてのリクエストが mux によって処理される
	funcframework.RegisterHTTPFunction("/", mux.ServeHTTP)

	// 環境変数 PORT が設定されていれば、そのポートでリッスン
	// Cloud Functions 環境では自動的に設定される
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	// フレームワークを開始し、リクエストの待受を開始
	log.Printf("Starting function server on port %s...\n", port)
	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}