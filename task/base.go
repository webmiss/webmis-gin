package task

import "fmt"

type Base struct{}

/* 输出到控制台 */
func (Base) Print(content ...interface{}) {
	fmt.Println(content...)
}
