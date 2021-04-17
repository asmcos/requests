package init

import (
	"os"
	"path"
	"runtime"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	println("filename:" + filename)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}
