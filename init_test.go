package parexec

import (
	"runtime"
)

const TOTALBENCH = 2000000

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
