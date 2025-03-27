#!/bin/bash

# テストパラメータ (必要に応じて調整してください)
FIXED_RATE=100      # 固定レートテストの秒間リクエスト数
FIXED_DURATION="30s" # 固定レートテストの実行時間
MAX_DURATION="30s"   # 最大スループットテストの実行時間

# ディレクトリ設定
RESULTS_DIR="results"
TARGETS_DIR="vegeta_targets"

# vegeta コマンドのパス (環境によってはフルパス指定が必要な場合あり)
VEGETA_CMD="vegeta"

# ターゲットファイル一覧を取得
# shellcheck disable=SC2012
targets=$(ls "${TARGETS_DIR}"/*.txt | sort) # ファイル名順にソート

echo "Starting Vegeta tests..."
echo "Fixed Rate: ${FIXED_RATE}/s, Duration: ${FIXED_DURATION}"
echo "Max Throughput Duration: ${MAX_DURATION}"
echo "Results will be saved in '${RESULTS_DIR}' directory."
echo "=================================================="

# サーバーが起動しているか確認 (簡易チェック)
if ! curl -s -o /dev/null http://localhost:8080/; then
    echo "ERROR: HTTP server is not running or not accessible at http://localhost:8080/"
    echo "Please start the server first (e.g., ./server_app)"
    exit 1
fi


for target_file in $targets; do
    # ターゲット名を取得 (例: none_small)
    target_name=$(basename "$target_file" .txt)
    echo ""
    echo "--- Testing target: $target_name ---"

    # --- 固定レートテスト ---
    fixed_bin_output="${RESULTS_DIR}/${target_name}_fixed.bin"
    fixed_txt_output="${RESULTS_DIR}/${target_name}_fixed.txt"
    echo "[${target_name}] Running fixed rate test..."
    if ${VEGETA_CMD} attack -targets="$target_file" -rate="$FIXED_RATE" -duration="$FIXED_DURATION" -output="$fixed_bin_output"; then
        echo "[${target_name}] Generating fixed rate report..."
        ${VEGETA_CMD} report "$fixed_bin_output" > "$fixed_txt_output"
        echo "[${target_name}] Fixed rate test completed."
    else
        echo "[${target_name}] ERROR: Fixed rate test failed."
    fi
    sleep 1 # サーバー負荷軽減のため少し待機

    # --- 最大スループットテスト ---
    max_bin_output="${RESULTS_DIR}/${target_name}_max.bin"
    max_txt_output="${RESULTS_DIR}/${target_name}_max.txt"
    echo "[${target_name}] Running max throughput test..."
    # rate=0 で最大スループットを試みる (-max-workers が必要)
    if ${VEGETA_CMD} attack -targets="$target_file" -rate=0 -max-workers=100 -duration="$MAX_DURATION" -output="$max_bin_output"; then
        echo "[${target_name}] Generating max throughput report..."
        ${VEGETA_CMD} report "$max_bin_output" > "$max_txt_output"
        echo "[${target_name}] Max throughput test completed."
    else
        echo "[${target_name}] ERROR: Max throughput test failed."
    fi
    sleep 1 # サーバー負荷軽減のため少し待機

done

echo "=================================================="
echo "All Vegeta tests completed."
echo "Check the '${RESULTS_DIR}' directory for results."