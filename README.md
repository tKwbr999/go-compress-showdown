# go-compress-showdown
GZIP、Brotli、Zstandard実際どれが速いの？って話

## 概要

このプロジェクトは、Go言語における主要な圧縮アルゴリズム（Gzip, Brotli, Zstandard）のパフォーマンスを比較するためのベンチマークツールです。
生成AIのレスポンスのようなテキストデータを効率的に転送するための最適な圧縮方式を見つけることを目的としています。

## ベンチマーク

### 実行方法 (ローカル)

1.  Go環境と `vegeta` をインストールします。
    ```bash
    go install github.com/tsenart/vegeta/v12@latest
    ```
2.  サーバーをビルドして起動します。
    ```bash
    go build -o server_app ./cmd/server
    ./server_app
    ```
3.  別のターミナルでベンチマークスクリプトを実行します。
    ```bash
    chmod +x ./run_vegeta_tests.sh
    ./run_vegeta_tests.sh
    ```
4.  `results` ディレクトリにテスト結果が出力されます。

### Go 標準ベンチマーク

純粋な圧縮・解凍性能は以下のコマンドで測定できます。
```bash
go test -bench=. -benchmem ./internal/compressor
```

### GitHub Actions による自動実行

このリポジトリでは、コードがプッシュされるたびにGitHub Actionsによって自動的にベンチマークが実行されます。

**最新のベンチマーク結果:**

最新のワークフロー実行結果は、リポジトリの **Actions** タブから確認できます。
各ワークフロー実行ページの下部にある **Artifacts** セクションから、以下の結果ファイルをダウンロードできます。

*   `go-benchmark-results`: Go標準ベンチマーク (`go test -bench`) の結果 (`go_benchmark_results.txt`)
*   `vegeta-results`: `vegeta` によるHTTP負荷テストの結果 (各ターゲットごとの `.txt` レポートファイルを含む `results` ディレクトリ)

[![Go Compression Benchmark](https://github.com/tKwbr999/go-compress-showdown/actions/workflows/benchmark.yml/badge.svg)](https://github.com/tKwbr999/go-compress-showdown/actions/workflows/benchmark.yml)
