.PHONY: report clean-results run-server stop-server wait-for-server test-load

SERVER_PID_FILE := server.pid
SERVER_URL := http://localhost:8080
WAIT_TIMEOUT := 60 # seconds

# レポートを生成する
report:
	@echo "Generating benchmark report..."
	@go run report_generator.go
	@echo "Report generation finished."

# 結果ファイルを削除する
clean-results:
	@echo "Cleaning up results directory..."
	@rm -rf results/*.bin results/*.txt
	@echo "Results cleaned."

# サーバーをバックグラウンドで起動する
run-server:
	@echo "Starting server in background..."
	@go run cmd/server/main.go & echo $$! > $(SERVER_PID_FILE)
	@echo "Server process started with PID $$(cat $(SERVER_PID_FILE))."

# サーバーを停止する
stop-server:
	@if [ -f $(SERVER_PID_FILE) ]; then \
		PID=$$(cat $(SERVER_PID_FILE)); \
		echo "Stopping server process with PID $$PID..."; \
		if kill $$PID; then \
			echo "Server stopped."; \
		else \
			echo "Failed to stop server process $$PID. It might have already stopped."; \
		fi; \
		rm -f $(SERVER_PID_FILE); \
	else \
		echo "Server PID file not found. Server might not be running or was stopped manually."; \
	fi

# サーバーが起動するまで待機する
wait-for-server:
	@echo "Waiting for server at $(SERVER_URL) to start..."
	@start_time=$$(date +%s); \
	while ! curl -s -o /dev/null -w "%{http_code}" $(SERVER_URL)/ | grep -q 200; do \
		current_time=$$(date +%s); \
		elapsed_time=$$((current_time - start_time)); \
		if [ $$elapsed_time -ge $(WAIT_TIMEOUT) ]; then \
			echo "Server did not start within $(WAIT_TIMEOUT) seconds."; \
			make stop-server; \
			exit 1; \
		fi; \
		echo -n "."; \
		sleep 2; \
	done; \
	echo " Server started."

# 負荷テストを実行する (サーバーの起動と停止を含む)
test-load: run-server wait-for-server
	@echo "Running load tests..."
	@trap 'echo "Load test interrupted. Stopping server..."; make stop-server' INT TERM EXIT; \
	go test -v -timeout 30m ./internal/handler -run ^TestLoad$; \
	TEST_EXIT_CODE=$$?; \
	trap - INT TERM EXIT; \
	echo "Load tests finished. Stopping server..."; \
	make stop-server; \
	exit $$TEST_EXIT_CODE

# (オプション) ベンチマークテスト全体を実行するターゲット
# run-benchmark: clean-results test-load report