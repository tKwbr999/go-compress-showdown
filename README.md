# go-compress-showdown
GZIP、Brotli、Zstandard実際どれが速いの？って話

## 概要

このプロジェクトは、Go言語における主要な圧縮アルゴリズム（Gzip, Brotli, Zstandard）のパフォーマンスを比較するためのベンチマークツールです。
生成AIのレスポンスのようなテキストデータを効率的に転送するための最適な圧縮方式を見つけることを目的としています。
`Makefile` を使用して、サーバーの起動、負荷テストの実行、レポート生成などを簡単に行うことができます。

## 使い方 (Makefile コマンド)

Go環境がセットアップされていることを前提とします。

1.  **負荷テストの実行とレポート生成:**
    以下のコマンドで、サーバーの起動、`vegeta` を利用した負荷テスト (内部的には `go test` 経由で実行)、サーバーの停止、そして結果レポート (`compression_benchmark_report.md` と `compression_benchmark_report.ja.md`) の生成までを一括で行います。
    ```bash
    make clean-results # 既存の結果を削除 (任意)
    make test-load     # サーバー起動、負荷テスト実行、サーバー停止
    make report        # 結果からレポートを生成
    ```
    *注意:* `test-load` は完了までに時間がかかる場合があります。

2.  **個別の操作:**
    *   **サーバーの起動:**
        ```bash
        make run-server
        make wait-for-server # サーバーが応答可能になるまで待機
        ```
    *   **サーバーの停止:**
        ```bash
        make stop-server
        ```
    *   **負荷テストのみ実行 (サーバーは別途起動):**
        `make run-server` でサーバーを起動後、以下のコマンドを実行します。
        ```bash
        go test -v -timeout 30m ./internal/handler -run ^TestLoad$
        ```
        テスト完了後、`make stop-server` でサーバーを停止してください。
    *   **レポート生成のみ:**
        `results` ディレクトリに `vegeta` の結果ファイル (`.bin`) が存在する場合に実行します。
        ```bash
        make report
        ```
    *   **結果ファイルのクリーンアップ:**
        ```bash
        make clean-results
        ```

## Go 標準ベンチマーク

圧縮・解凍処理自体の純粋なパフォーマンスは、以下のコマンドで測定できます。
```bash
go test -bench=. -benchmem ./internal/compressor
```

## GitHub Actions による自動実行

このリポジトリでは、コードがプッシュされるたびにGitHub Actionsによって自動的にベンチマーク (`make test-load` と `make report`) が実行されます。

**最新のベンチマーク結果:**

最新のワークフロー実行結果は、リポジトリの **Actions** タブから確認できます。
各ワークフロー実行ページの下部にある **Artifacts** セクションから、以下の結果ファイルをダウンロードできます。

*   `go-benchmark-results`: Go標準ベンチマーク (`go test -bench`) の結果 (`go_benchmark_results.txt`)
*   `vegeta-results`: `vegeta` によるHTTP負荷テストの結果 (各ターゲットごとの `.txt` レポートファイルを含む `results` ディレクトリ)
*   `benchmark-report`: `make report` で生成されたマークダウンレポート (`compression_benchmark_report.md`, `compression_benchmark_report.ja.md`)

[![Go Compression Benchmark](https://github.com/tKwbr999/go-compress-showdown/actions/workflows/benchmark.yml/badge.svg)](https://github.com/tKwbr999/go-compress-showdown/actions/workflows/benchmark.yml)

## 関連ドキュメント

*   [生成AIレスポンス高速化総合レポート (ai-response-optimization.md)](ai-response-optimization.md): ストリーミング、圧縮、LLM推論最適化、キャッシュ、UX改善など、生成AIサービスのレスポンスを高速化するための技術を包括的に解説したレポート。
*   [圧縮ベンチマーク計画 (compression_benchmark_plan.md)](compression_benchmark_plan.md): Go言語における各種圧縮アルゴリズムのパフォーマンス比較計画書。
*   [圧縮ベンチマークレポート (compression_benchmark_report.md)](compression_benchmark_report.md): 圧縮アルゴリズムのベンチマーク結果レポート（英語版）。
*   [圧縮ベンチマークレポート (日本語版) (compression_benchmark_report.ja.md)](compression_benchmark_report.ja.md): 圧縮アルゴリズムのベンチマーク結果レポート（日本語版）。
*   [レポート生成プログラム修正計画 (report_generator_fix_plan.md)](report_generator_fix_plan.md): ベンチマークレポート生成プログラムの不具合修正計画書。
