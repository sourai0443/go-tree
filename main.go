package main

import (
	"bytes"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	target          string
	out             string
	isDirectoryOnly bool
)

func init() {
	kingpin.Flag("target", "").Short('t').Default(".").StringVar(&target)
	kingpin.Flag("out", "").Short('o').Default("").StringVar(&out)
	kingpin.Flag("dir-only", "").Short('d').Default("false").BoolVar(&isDirectoryOnly)
	kingpin.Parse()
}

const (
	outputFileNamePattern = "200601021504.txt"
	timeStampPattern      = "2006/01/02 15:04:05"
	blank                 = ""
	space                 = " "
	verticalBorder        = "│"
	rightExistBorder      = "├"
	endBorder             = "└"
)

func main() {
	skip := func(entry os.DirEntry) bool {
		if isDirectoryOnly {
			return !entry.IsDir()
		} else {
			return false
		}
	}

	if !strings.EqualFold(out, blank) {
		timeStamp := time.Now()
		o, err := os.OpenFile(timeStamp.Format(outputFileNamePattern), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0766)
		if err != nil {
			panic(err)
		}
		fmt.Fprintln(o, blank)
		fmt.Fprintln(o, timeStamp.Format(timeStampPattern), strings.Join(os.Args, space))
		getDirNames(target, o, skip)
	} else {
		getDirNames(target, os.Stdout, skip)
	}

}

func getDirNames(root string, out io.Writer, skipFunc func(entry os.DirEntry) bool) {
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
			if dirFiles[path], err = getFileCount(path, isDirectoryOnly); err != nil {
				panic(err)
			}
		}

		// 処理中のファイルが何番目かをインクリメント
		filesMap[p]++
		// インデントの深さを取得
		cnt := strings.Count(path, string(os.PathSeparator))
		// バッファを宣言（バッファに書き込んだものを指定した出力先に出力する）。
		buf := &bytes.Buffer{}
		// 空白に指定されていなければ、│を出力し空白の場合は空文字列を出力
		printVerticalBorder(buf, cnt, blankMap)

		if dirFiles[p] == filesMap[p] {
			// ディレクトリのファイル数と処理中のファイル数が一致した場合は└を出力。同じ列を空白に指定するマップに追加
			fmt.Fprintln(buf, endBorder+d.Name())
			blankMap[cnt] = struct{}{}
		} else if strings.EqualFold(root, orgPath) {
			// rootと同じファイルは罫線を着けずに出力
			fmt.Fprintln(buf, root)
		} else {
			// ファイルの終端でも、rootファイルでもない場合は├を出力し、同じ列を空白に指定するマップから指定列数を削除
			fmt.Fprintln(buf, rightExistBorder+d.Name())
			delete(blankMap, cnt)
		}
		// 文字列として出力
		fmt.Fprint(out, buf.String())
		return nil
	})
}

func getFileCount(path string, isDirectoryOnly bool) (int, error) {
	var fileCnt int
	des, err := os.ReadDir(path)
	if err != nil {
		return 0, err
	}

	if isDirectoryOnly {
		// ディレクトリのみの指定の場合は、ファイル数をディレクトリのみをカウントする
		for _, de := range des {
			if de.IsDir() {
				fileCnt++
			}
		}
	} else {
		// ディレクトリ、ファイルを含む場合はターゲット以下全てのファイル数を取得
		fileCnt = len(des)
	}
	return fileCnt, nil
}

func printVerticalBorder(out io.Writer, max int, blankMap map[int]interface{}) {
	for i := 0; i < max; i++ {
		if _, ok := blankMap[i]; ok {
			fmt.Fprint(out, " ")
		} else {
			fmt.Fprint(out, "│")
		}
	}
}
