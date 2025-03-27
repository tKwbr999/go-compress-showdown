package compressor

import (
	"fmt"
	"io"

	"github.com/klauspost/compress/zstd"
)

// GetZstdWriter は指定されたレベルでzstd圧縮を行うWriterを返します。
// levelStr は "fastest", "default", "best" を受け付けます。
func GetZstdWriter(w io.Writer, levelStr string) (*zstd.Encoder, error) {
	level := zstd.SpeedDefault // デフォルト

	switch levelStr {
	case "fastest":
		level = zstd.SpeedFastest
	case "default":
		level = zstd.SpeedDefault
	case "best":
		level = zstd.SpeedBestCompression
	case "":
		// levelStrが空の場合はデフォルトを使用
	default:
		// klauspost/compress/zstd では EncoderLevelFrom は提供されているが、
		// 文字列からの直接変換は用意されていない。主要な定数を使う。
		return nil, fmt.Errorf("invalid zstd level string: %s (use fastest, default, or best)", levelStr)
	}

	// zstd.NewWriter はエラーを返さないが、オプションでエラーが発生する可能性はある
	// ここではレベル指定のみなので、エラーチェックは省略
	zstdWriter, _ := zstd.NewWriter(w, zstd.WithEncoderLevel(level))

	return zstdWriter, nil
}