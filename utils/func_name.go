package utils

import (
	"runtime"
	"strings"
)

func GetFuncName() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic("")
	}
	ffn := runtime.FuncForPC(pc).Name()
	return ffn[strings.LastIndex(ffn, ".")+1:]
}
