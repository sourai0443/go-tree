package main

import (
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
	filesMap := make(map[string]int)
	var desCnt int
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if skipFunc != nil && skipFunc(d) {
			return nil
		}

		p, _ := filepath.Split(path)
		if d.IsDir() {
			filesMap[p] = 0
			des, err := os.ReadDir(path)
			if err != nil {
				return err
			}
			desCnt = len(des)
		}
		filesMap[p]++
		fmt.Println(desCnt, filesMap[p], d.Name(), p)

		cnt := strings.Count(path, string(os.PathSeparator))
		for i := 0; i < cnt; i++ {
			fmt.Fprint(out, "│")
		}

		if desCnt == filesMap[p] {
			fmt.Fprintf(out, "└%s\n", d.Name())
			desCnt = 0
		} else {
			fmt.Fprintf(out, "┝%s\n", d.Name())
		}
		return nil
	})
}
