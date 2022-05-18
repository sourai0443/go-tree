package main

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	getDirNames(".", nil, os.Stdout)
}

// 0 => ┝、or │ or └
// 1 => │　+ ┝、or │ or └

func getDirNames(root string, skipFunc func(entry os.DirEntry) bool, out io.Writer) {
	// 現在処理中のファイル数
	filesMap := make(map[string]int)
	// rootに指定したファイルが+1されるため打消し
	filesMap[root] = -1
	// ディレクトリ内のファイル数（ディレクトリを含む）
	dirFiles := make(map[string]int)
	// 空白列を記録するマップ
	blankMap := make(map[int]interface{})

	filepath.WalkDir(root, func(orgPath string, d fs.DirEntry, err error) error {
		if skipFunc != nil && skipFunc(d) {
			return nil
		}

		path := filepath.Clean(orgPath)
		p, _ := filepath.Split(path)
		p = filepath.Clean(p)

		// 処理中のファイルがディレクトリの場合の処理
		// ディレクトリ内のファイル数を記録する（└を打つため）
		if d.IsDir() {
			des, err := os.ReadDir(path)
			if err != nil {
				return err
			}
			dirFiles[path] = len(des)
		}

		// 処理中のファイルが何番目かをインクリメント
		filesMap[p]++
		// インデントの深さを取得
		cnt := strings.Count(path, string(os.PathSeparator))

		buf := &bytes.Buffer{}
		// 空白に指定されていなければ、│を出力し空白の場合は空文字列を出力
		for i := 0; i < cnt; i++ {
			if _, ok := blankMap[i]; ok {
				fmt.Fprint(buf, " ")
			} else {
				fmt.Fprint(buf, "│")
			}
		}

		if dirFiles[p] == filesMap[p] {
			// ディレクトリのファイル数と処理中のファイル数が一致した場合は└を出力。同じ列を空白に指定するマップに追加
			fmt.Fprintf(buf, "└%s\n", d.Name())
			blankMap[cnt] = struct{}{}
		} else if strings.EqualFold(root, orgPath) {
			// rootと同じファイルは罫線を着けずに出力
			fmt.Fprintln(buf, root)
		} else {
			// ファイルの終端でも、rootファイルでもない場合は├を出力し、同じ列を空白に指定するマップから指定列数を削除
			fmt.Fprintf(buf, "┝%s\n", d.Name())
			delete(blankMap, cnt)
		}
		// 文字列として出力
		fmt.Fprint(out, buf.String())
		return nil
	})
}
