package parexec

import (
	"runtime"
)

const TOTALBENCH = 200000

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
