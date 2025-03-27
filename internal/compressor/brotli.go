package compressor

import (
	"fmt"
	"io"
	"strconv"

	"github.com/andybalholm/brotli"
)

// GetBrotliWriter は指定された品質レベルでbrotli圧縮を行うWriterを返します。
// qualityStr は数値文字列 (0-11) を受け付けます。デフォルトは 6 です。
func GetBrotliWriter(w io.Writer, qualityStr string) (*brotli.Writer, error) {
	quality := brotli.DefaultCompression // デフォルトはレベル6

	if qualityStr != "" {
		parsedQuality, err := strconv.Atoi(qualityStr)
		if err != nil {
			return nil, fmt.Errorf("invalid brotli quality string: %s", qualityStr)
		}
		// Brotliの品質レベルは 0 から 11
		if parsedQuality < 0 {
			quality = 0
		} else if parsedQuality > 11 {
			quality = 11
		} else {
			quality = parsedQuality
		}
	}

	// andybalholm/brotli では NewWriterLevel ではなく WriterOptions を使う
	options := brotli.WriterOptions{
		Quality: quality,
	}
	brotliWriter := brotli.NewWriterOptions(w, options)

	return brotliWriter, nil // NewWriterOptions はエラーを返さない
}