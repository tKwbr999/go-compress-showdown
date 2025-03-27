package data

import (
	"fmt"
	"log" // log パッケージを追加
	"os"
	"path/filepath"
	"sync"
)

var (
	cache     = make(map[string][]byte)
	cacheOnce sync.Once
	cacheErr  error
)

const testdataDir = "testdata" // testdataディレクトリのパス

// loadAndCache は testdata ディレクトリから全ての .txt ファイルを読み込み、キャッシュします。
// 初回呼び出し時のみ実行されます。
func loadAndCache() {
	cacheOnce.Do(func() {
		log.Println("Attempting to load test data...") // 開始ログ
		files := []string{"small.txt", "medium.txt", "large.txt"}
		successCount := 0
		for _, filename := range files {
			key := filename[:len(filename)-len(filepath.Ext(filename))] // "small", "medium", "large"
			path := filepath.Join(testdataDir, filename)
			log.Printf("Reading file: %s", path) // ファイル読み込みログ

			dataBytes, err := os.ReadFile(path)
			if err != nil {
				cacheErr = fmt.Errorf("failed to read test data file '%s': %w", path, err)
				log.Printf("ERROR reading file %s: %v", path, err) // エラーログ
				// エラーが発生したらキャッシュ処理を中断
				return
			}
			cache[key] = dataBytes
			log.Printf("Successfully loaded %s (%d bytes)", path, len(dataBytes)) // 成功ログ
			successCount++
		}
		if cacheErr == nil {
			log.Printf("Successfully loaded %d test data files.", successCount) // 完了ログ
		}
	})
}

// LoadData は指定されたキーに対応するテストデータをキャッシュから取得します。
// 初回呼び出し時に testdata ディレクトリからファイルを読み込みキャッシュします。
// キーは "small", "medium", "large" を想定しています。
func LoadData(key string) ([]byte, error) {
	loadAndCache() // 必要であればキャッシュを初期化
	if cacheErr != nil {
		return nil, cacheErr // キャッシュ初期化時のエラーを返す
	}

	dataBytes, ok := cache[key]
	if !ok {
		return nil, fmt.Errorf("test data not found for key: %s (expected small, medium, or large)", key)
	}
	return dataBytes, nil
}

// GetDataKeys は利用可能なデータキーのリストを返します。
func GetDataKeys() []string {
	loadAndCache()
	keys := make([]string, 0, len(cache))
	for k := range cache {
		keys = append(keys, k)
	}
	// 必要であればソートする
	// sort.Strings(keys)
	return keys
}