#!/bin/bash

# テストパラメータ (必要に応じて調整してください)
FIXED_RATE=50       # 固定レートテストの秒間リクエスト数 (負荷軽減のため調整)
FIXED_DURATION="15s" # 固定レートテストの実行時間 (負荷軽減のため調整)
MAX_DURATION="15s"   # 最大スループットテストの実行時間 (負荷軽減のため調整)

# テスト対象ホスト (環境変数 TARGET_HOST が設定されていればそれを使用)
TARGET_HOST="${TARGET_HOST:-http://localhost:8080}"

# ディレクトリ設定
RESULTS_DIR="results"
TARGETS_DIR="vegeta_targets"
TEMP_DIR=$(mktemp -d) # 一時ファイル用ディレクトリ

# vegeta コマンドのパス (環境によってはフルパス指定が必要な場合あり)
VEGETA_CMD="vegeta"

# 結果ディレクトリを作成 (存在しない場合)
mkdir -p "$RESULTS_DIR"

# 一時ディレクトリ削除のトラップ
trap 'rm -rf "$TEMP_DIR"' EXIT

# ターゲットファイル一覧を取得
# shellcheck disable=SC2012
targets=$(ls "${TARGETS_DIR}"/*.txt | sort) # ファイル名順にソート

echo "Starting Vegeta tests..."
echo "Target Host: ${TARGET_HOST}"
echo "Fixed Rate: ${FIXED_RATE}/s, Duration: ${FIXED_DURATION}"
echo "Max Throughput Duration: ${MAX_DURATION}"
echo "Results will be saved in '${RESULTS_DIR}' directory."
echo "=================================================="

# サーバーが起動しているか確認 (簡易チェック)
if ! curl -s -o /dev/null "${TARGET_HOST}/"; then
    echo "ERROR: HTTP server is not running or not accessible at ${TARGET_HOST}/"
    echo "Please start the server first (e.g., ./server_app)"
    exit 1
fi


for target_file in $targets; do
    # ターゲット名を取得 (例: none_small)
    target_name=$(basename "$target_file" .txt)
    echo ""
    echo "--- Testing target: $target_name ---"

    # 一時的なターゲットファイルを作成し、ホスト名を置換
    temp_target_file="${TEMP_DIR}/${target_name}.txt"
    sed "s|http://localhost:8080|${TARGET_HOST}|g" "$target_file" > "$temp_target_file"

    # --- 固定レートテスト ---
    fixed_bin_output="${RESULTS_DIR}/${target_name}_fixed.bin"
    fixed_txt_output="${RESULTS_DIR}/${target_name}_fixed.txt"
    echo "[${target_name}] Running fixed rate test..."
    if ${VEGETA_CMD} attack -targets="$temp_target_file" -rate="$FIXED_RATE" -duration="$FIXED_DURATION" -output="$fixed_bin_output"; then
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
    # rate=0 で最大スループットを試みる (-max-workers を調整)
    if ${VEGETA_CMD} attack -targets="$temp_target_file" -rate=0 -max-workers=50 -duration="$MAX_DURATION" -output="$max_bin_output"; then
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