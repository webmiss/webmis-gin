package service

import (
	"webmis/config"
	"webmis/library"
	"webmis/util"
)

/* 日志 */
type Logs struct{}

/* 文件 */
func (Logs) File(file string, content string) {
	(&library.FileEo{}).New(config.Env().RootDir)
	(&library.FileEo{}).WriterEnd(file, content+"\n")
}

/* 生产者 */
func (Logs) Log(data interface{}) {
	redis := (&library.Redis{}).New("")
	redis.RPush("logs", util.JsonEncode(data))
	redis.Close()
}
