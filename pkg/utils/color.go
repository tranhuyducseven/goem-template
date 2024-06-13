package utils

import (
	"fmt"
)

type ForegroundColor int

const (
	Black   ForegroundColor = 30
	Red     ForegroundColor = 31
	Green   ForegroundColor = 32
	Yellow  ForegroundColor = 33
	Blue    ForegroundColor = 34
	Magenta ForegroundColor = 35
	Cyan    ForegroundColor = 36
	White   ForegroundColor = 37
)

func ColorMethod(method string) string {
	switch method {
	case "GET":
		return Logf(Green, method)
	case "POST":
		return Logf(Yellow, method)
	case "PUT":
		return Logf(Blue, method)
	case "DELETE":
		return Logf(Red, method)
	case "PATCH":
		return Logf(Magenta, method)
	default:
		return Logf(Cyan, method)
	}
}

func MagentaMsg(s string) string {
	return Logf(Magenta, s)
}

func BlueMsg(s string) string {
	return Logf(Blue, s)
}

func YellowMsg(s string) string {
	return Logf(Yellow, s)
}

func Logf(foreground ForegroundColor, v ...interface{}) string {
	return fmt.Sprintf("\033[%dm %s\033[0m", foreground, fmt.Sprint(v...))
}
