package main

import (
	"fmt"
	"io"
	"os"
)

const (
	bufferSize = 8192
)

func simpleCat(filename string) error {
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

	buf := make([]byte, bufferSize)

	for {
		n, err := file.Read(buf)

		if n > 0 {
			written, writeErr := os.Stdout.Write(buf[:n])
			if writeErr != nil {
				return fmt.Errorf("write error: %w", writeErr)
			}
			if written != n {
				return fmt.Errorf("write error: wrote %d bytes, expected %d", written, n)
			}
		}


		if err == io.EOF {
			return nil
		}
		if err != nil {
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
