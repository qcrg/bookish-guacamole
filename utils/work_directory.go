package utils

import (
	"os"
	"path"
	"path/filepath"
	"sync"
)

func get_project_dir_impl() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	for {

		stats, _ := os.Stat(filepath.Join(wd, "go.mod"))
		if stats != nil && !stats.IsDir() {
			return wd
		}
		if wd == "/" {
			panic("Work directory, where is go.mod file, not found")
		}
		wd = path.Dir(wd)
	}
}

var GetProjectDir = sync.OnceValue(get_project_dir_impl)
