package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"

	mapset "github.com/deckarep/golang-set/v2"
)

var suffixList = []string{
	"Cargo.toml",
	".git",
	"package.json",
	"node_modules",
	".terraform",
	".terragrunt-cache",
	".direnv",
	"go.mod",
}
var projectPaths = mapset.NewSet[string]()

func main() {
	var pathname string
	if len(os.Args[1:]) == 1 {
		pathname = os.Args[1]
		if pathname == "-h" || pathname == "--help" {
			usage()
			os.Exit(0)
		}

		var err error
		pathname, err = filepath.Abs(pathname)
		if err != nil {
			log.Panicln(err)
		}
		fileInfo, err := os.Stat(pathname)
		if err != nil {
			log.Panicln(err)
		}
		if !fileInfo.IsDir() {
			usage()
			os.Exit(1)
		}
	} else {
		usage()
		os.Exit(1)
	}

	fileSystem := os.DirFS(pathname)

	fs.WalkDir(fileSystem, ".", func(subPath string, d fs.DirEntry, err error) error {
		for _, suffix := range suffixList {
			if filepath.Base(subPath) == suffix {
				parentPath := path.Join(pathname, path.Dir(subPath))
				if !projectPaths.Contains(parentPath) {
					fmt.Println(parentPath)
					projectPaths.Add(parentPath)
				}
				return fs.SkipDir
			}
		}
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})

	// paths := projectPaths.ToSlice()
	// slices.SortFunc(paths, func(a, b string) int {
	// 	return cmp.Compare(strings.Count(a, "/"), strings.Count(b, "/"))
	// })
	// for _, path := range paths {
	// 	fmt.Println(path)
	// }
}

func usage() {
	fmt.Printf(`
%s <path>
	<path>		path that should be searched for potential project directorys
`, os.Args[0])
}
