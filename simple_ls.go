package main

import (
	"fmt"
	"os"
	"sort"
	"syscall"
	"unsafe"
)

func getTerminalWidth() int {
	type winsize struct {
		Row    uint16
		Col    uint16
		Xpixel uint16
		Ypixel uint16
	}

	ws := &winsize{}
	ret, _, _ := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdout),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(ret) == -1 {
		return 80
	}
	return int(ws.Col)
}

// ファイル名のリストを複数カラムで表示
func printMultiColumn(names []string, termWidth int) {
	if len(names) == 0 {
		return
	}

	// 最大ファイル名長を計算
	maxLen := 0
	for _, name := range names {
		if len(name) > maxLen {
			maxLen = len(name)
		}
	}

	// カラム幅（スペース2文字分の余白を追加）
	colWidth := maxLen + 2
	if colWidth > termWidth {
		colWidth = termWidth
	}

	// カラム数を計算
	numCols := termWidth / colWidth
	if numCols < 1 {
		numCols = 1
	}

	// 行数を計算
	numRows := (len(names) + numCols - 1) / numCols

	// 縦方向に並べて表示（GNU lsと同じ）
	for row := 0; row < numRows; row++ {
		for col := 0; col < numCols; col++ {
			idx := col*numRows + row
			if idx >= len(names) {
				break
			}

			name := names[idx]
			fmt.Print(name)

			// 最後のカラムでなければスペースを追加
			if col < numCols-1 && idx+numRows < len(names) {
				padding := colWidth - len(name)
				for i := 0; i < padding; i++ {
					fmt.Print(" ")
				}
			}
		}
		fmt.Println()
	}
}

// ファイル名のリストを1行に1つずつ表示
func printOnePerLine(names []string) {
	for _, name := range names {
		fmt.Println(name)
	}
}

func main() {
	// 引数がなければカレントディレクトリ、あれば指定されたディレクトリ
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	// ディレクトリを開く
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ls: cannot access '%s': %v\n", dir, err)
		os.Exit(2)
	}

	// ファイル名を収集（隠しファイルを除く）
	var names []string
	for _, entry := range entries {
		name := entry.Name()
		// "." で始まるファイルはスキップ（隠しファイル）
		if len(name) > 0 && name[0] != '.' {
			names = append(names, name)
		}
	}

	// アルファベット順にソート
	sort.Strings(names)

	// 標準出力が端末かパイプかを判定
	fileInfo, _ := os.Stdout.Stat()
	isTerminal := (fileInfo.Mode() & os.ModeCharDevice) != 0

	if isTerminal {
		// 端末の場合：複数カラムで表示
		termWidth := getTerminalWidth()
		printMultiColumn(names, termWidth)
	} else {
		// パイプの場合：1行に1つずつ表示
		printOnePerLine(names)
	}
}
