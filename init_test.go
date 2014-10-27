package parexec

import (
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
}
