package compressor

import (
	"compress/gzip"
	"fmt"
	"io"
	"strconv"
)

// GetGzipWriter は指定されたレベルでgzip圧縮を行うWriterを返します。
// levelStr は "speed", "default", "best" または数値文字列を受け付けます。
func GetGzipWriter(w io.Writer, levelStr string) (*gzip.Writer, error) {
	level := gzip.DefaultCompression // デフォルトは標準圧縮

	switch levelStr {
	case "speed":
		level = gzip.BestSpeed
	case "default":
		level = gzip.DefaultCompression
	case "best":
		level = gzip.BestCompression
	case "":
		// levelStrが空の場合はデフォルトを使用
	default:
		// 数値指定の場合
		parsedLevel, err := strconv.Atoi(levelStr)
		if err != nil {
			return nil, fmt.Errorf("invalid gzip level string: %s", levelStr)
		}
		// gzipパッケージのレベル範囲に調整 (ただし、ドキュメント上は-1から9)
		// ここでは BestSpeed, DefaultCompression, BestCompression の定数を使うことを推奨
		// 数値指定は参考程度とし、主要な定数にマッピングする方が安全かもしれない
		if parsedLevel < gzip.BestSpeed { // -2 より小さい場合
			level = gzip.BestSpeed
		} else if parsedLevel > gzip.BestCompression { // 9 より大きい場合
			level = gzip.BestCompression
		} else {
			level = parsedLevel
		}
	}

	gzipWriter, err := gzip.NewWriterLevel(w, level)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip writer: %w", err)
	}
	return gzipWriter, nil
}