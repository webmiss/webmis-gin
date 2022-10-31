package task

import (
	"webmis/library"
	"webmis/service"
	"webmis/util"
)

/* 日志 */
type Logs struct{ Base }

/* 消费者 */
func (r Logs) Log() {
	for {
		redis := (&library.Redis{}).New("")
		data := redis.BLPop("logs", 10)
		redis.Close()
		if data == nil {
			continue
		}
		// 保存
		msg := (&util.Type{}).Strval(data[1])
		res := r._logsWrite(msg)
		if !res {
			(&service.Logs{}).File("upload/erp/Logs.json", string(util.JsonEncode(msg)))
		}
	}
}

/* 日志-写入 */
func (r Logs) _logsWrite(msg string) bool {
	// 数据
	data := map[string]interface{}{}
	util.JsonDecode(msg, &data)
	r.Print(data)
	return true
}
