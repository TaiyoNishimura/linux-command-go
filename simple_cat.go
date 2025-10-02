package main

import (
	"fmt"
	"io"
	"os"
)

const (
	// バッファサイズ（C版と同じく最適なI/Oサイズを使用）
	// 通常は4KB-8KBが適切
	bufferSize = 8192
)

// simpleCat は指定されたファイルを標準出力にコピーする
// C版のsimple_cat関数に相当
func simpleCat(filename string) error {
	// ファイルを開く（"-"の場合は標準入力）
	var file *os.File
	var err error

	if filename == "-" {
		file = os.Stdin
	} else {
		file, err = os.Open(filename)
		if err != nil {
			return fmt.Errorf("%s: %w", filename, err)
		}
		defer file.Close()
	}

	// バッファを確保
	buf := make([]byte, bufferSize)

	// ファイル終端まで読み込みと書き込みを繰り返す
	for {
		// ① ファイルから最大bufferSizeバイト読み込む
		n, err := file.Read(buf)

		// ② 読み込んだデータがあれば書き出す
		if n > 0 {
			// 標準出力に書き込む（full_writeに相当）
			written, writeErr := os.Stdout.Write(buf[:n])
			if writeErr != nil {
				return fmt.Errorf("write error: %w", writeErr)
			}
			if written != n {
				return fmt.Errorf("write error: wrote %d bytes, expected %d", written, n)
			}
		}

		// ③ エラーチェック
		if err == io.EOF {
			// ファイル終端に到達（正常終了）
			return nil
		}
		if err != nil {
			// 読み込みエラー
			return fmt.Errorf("read error: %w", err)
		}
	}
}

func main() {
	// 引数がない場合は標準入力から読む
	if len(os.Args) == 1 {
		if err := simpleCat("-"); err != nil {
			fmt.Fprintf(os.Stderr, "cat: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// 各ファイルを順番に処理
	exitCode := 0
	for _, filename := range os.Args[1:] {
		if err := simpleCat(filename); err != nil {
			fmt.Fprintf(os.Stderr, "cat: %v\n", err)
			exitCode = 1
		}
	}

	os.Exit(exitCode)
}
