# Report Generator 修正計画

## 1. 背景

`compression_benchmark_report.md` に記載されているベンチマーク結果において、特に Latency (レイテンシ) と Duration (テスト時間) の数値が不正確である問題が確認されました。この計画は、`report_generator.go` を修正し、これらの数値を正確にレポートに反映させることを目的とします。

## 2. 原因分析

*   **Latency (ずれ/0.00ms 問題):**
    *   `vegeta report` コマンドのテキスト出力における Latencies の項目順序 (`min, mean, 50, 90, 95, 99, max`) と、`report_generator.go` 内の正規表現が期待する順序 (`mean, 50, 95, 99, max`) が一致していません。
    *   これにより、値がずれてパースされ、`strconv.ParseFloat` がエラーとなり、結果として 0.00ms と表示されている可能性が高いです。
*   **Duration (0.00s 問題):**
    *   `vegeta report` の出力自体は正しい時間 (例: `14.984s`) を示していますが、`report_generator.go` でのパース (`strconv.ParseFloat`) 時のエラーハンドリングが不足しているため、エラーが発生しても検知できず 0.00s となっている可能性があります。
    *   また、特定の `.bin` ファイルに対する `vegeta report` コマンド実行自体が失敗している可能性も考えられますが、現状のエラーハンドリングでは詳細が不明です。
*   **Compressed Size / Ratio (0 bytes / 0.00x 問題):**
    *   `vegeta report` のテキスト形式では、レスポンスボディの合計サイズ (`Bytes Out`) が 0 と報告されます。これは `vegeta` の仕様です。
    *   `report_generator.go` はこの `Bytes Out` が 0 であるため、圧縮後のサイズを計算できず、結果として 0 と表示しています (`none` アルゴリズムを除く)。

## 3. 修正計画

以下の修正を `report_generator.go` に適用します。

1.  **Latency の修正:**
    *   `latenciesRe` (L55付近) の正規表現を、`vegeta report` の実際の出力形式 (`min, mean, 50, 90, 95, 99, max`) に合わせて修正します。
    *   `parseVegetaReport` 関数内 (L138-143付近) で、修正後の正規表現に合わせて正しいインデックス (`m[2]` for mean, `m[3]` for 50th, `m[5]` for 95th, `m[6]` for 99th) から値を抽出するように修正します。
    *   Latency 値のパース (`strconv.ParseFloat`) にエラーチェックを追加し、エラー発生時にはログを出力します。
2.  **Duration の修正:**
    *   Duration 値のパース (`strconv.ParseFloat`, L136付近) にエラーチェックを追加し、エラー発生時にはログを出力します。
    *   `vegeta report` コマンド実行部分 (L115-124付近) のエラーハンドリングを強化し、コマンド失敗時や標準エラー出力があった場合に、対象ファイル名を含む詳細なエラーログを出力します。
3.  **Compressed Size / Ratio の修正:**
    *   **今回は修正を見送ります。** `none` 以外のアルゴリズムでは、現状通り 0 bytes / 0.00x と表示されます。サイズの正確な計測方法については、別途検討します。

## 4. 計画図 (Mermaid)

```mermaid
graph TD
    A[問題分析: レポート数値の不正確さ] --> B{原因特定};
    B --> C1[Duration: 0.00s 問題 (Parse/Exec Error?)] ;
    B --> C2[Comp. Size/Ratio: 0 問題 (Bytes Out 不足)];
    B --> C3[Latency: ずれ/0.00ms 問題 (Regex 不一致)];

    C1 --> F1[修正案: エラーハンドリング強化];
    C2 --> F2[修正案: 今回は見送り];
    C3 --> F3[修正案: Regex とパースロジック修正 (必須)];

    F1 & F2 & F3 --> H[修正計画策定];
    H --> I{ユーザー承認};
    I -- Yes --> J[実装依頼 (別モード)];
    I -- No --> H;
```

## 5. 次のステップ

この承認された計画に基づき、`code` モードに切り替えて `report_generator.go` のコード修正を依頼します。