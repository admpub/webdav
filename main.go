package main

import (
	"runtime"

	"github.com/admpub/webdav/v4/cmd"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cmd.Execute()
}
