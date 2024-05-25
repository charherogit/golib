package caller

import (
	"bytes"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

// 调用方的简短信息
// 获取 eg: [game.go:feedObject:1190] skip 当前行数传递1，上一层2 FFL=file function line
func BriefInfo(skip int) *bytes.Buffer {
	buffer := bytes.NewBuffer(make([]byte, 0, 32))
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		funcName := runtime.FuncForPC(pc).Name()
		fileName := filepath.Base(file)
		funcNames := strings.Split(funcName, ".")
		funcName = funcNames[len(funcNames)-1]
		buffer.WriteString("[")
		buffer.WriteString(fileName)
		buffer.WriteString(":")
		buffer.WriteString(funcName)
		buffer.WriteString(":")
		buffer.WriteString(strconv.Itoa(line))
		buffer.WriteString("]")
	}
	return buffer
}

func BriefInfoStr(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		return fmt.Sprintf("[%s:%s:%d]",
			filepath.Base(file), TrimFn(runtime.FuncForPC(pc).Name()), line)
	}
	return "jump to long"
}

func Fn(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if ok {
		return TrimFn(runtime.FuncForPC(pc).Name())
	}
	return "jump to long"
}

func TrimFn(fn string) string {
	for i := len(fn) - 1; i >= 0; i-- {
		if fn[i] == '.' {
			return fn[i+1:]
		}
	}
	return fn
}

func FuncName(fp any) string {
	v := reflect.ValueOf(fp)
	if v.Kind() == reflect.Func {
		return TrimFn(runtime.FuncForPC(v.Pointer()).Name())
	}
	return "not a function"
}
