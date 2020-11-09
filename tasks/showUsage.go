package tasks

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/astaxie/beego"
)

// printMemUsage printMemUsage
func printMemUsage() error {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	msg := "Alloc = " + bToMb(m.Alloc) + " MiB, TotalAlloc = " + bToMb(m.TotalAlloc) + " MiB, Sys = " + bToMb(m.Sys) + " MiB, NumGC = " + fmt.Sprint(m.NumGC)
	beego.Informational(msg)
	return nil
}

func bToMb(b uint64) string {
	toString := strconv.FormatUint(b/1024/1024, 10)
	return toString
}
