package routes

import (
	"runtime"
	"strings"
)

var packagePath = getPackagePath()

func getPackagePath() string {
	pc, _, _, _ := runtime.Caller(0)
	name := runtime.FuncForPC(pc).Name()
	return name[:strings.LastIndex(name, ".")]
}
