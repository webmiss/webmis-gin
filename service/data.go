package service

import (
	"crypto/rand"
	"math/big"
	"sort"
	"time"
	"webmis/config"
	"webmis/library"
	"webmis/model"
	"webmis/util"
)

/* 数据类 */
type Data struct{ Base }

// 机器标识
var machineId int64 = config.Env().MachineId

const (
	max8bit  = uint(8)  //随机数位数
	max10bit = uint(10) //机器位数
	max12bit = uint(12) //序列数位数
)

// 分区时间
var partition = map[string]int{
	"p2208": 1661961600,
	"p2209": 1664553600,
	"plast": 1664553600,
}

/* 薄雾算法 */
func (Data) Mist(redisName string) int64 {
	// 自增ID
	redis := (&library.Redis{}).New("")
	autoId := redis.Incr(redisName)
	redis.Close()
	// 随机数
	randA, _ := rand.Int(rand.Reader, big.NewInt(255))
	randB, _ := rand.Int(rand.Reader, big.NewInt(255))
	// 位运算
	mist := int64((autoId << (max8bit + max8bit)) | (randA.Int64() << max8bit) | randB.Int64())
	return mist
}

/* 雪花算法 */
func (Data) Snowflake() int64 {
	// 时间戳
	now := time.Now()
	t := now.UnixNano() / 1e6
	// 随机数
	rand, _ := rand.Int(rand.Reader, big.NewInt(4095))
	// 位运算
	mist := int64((t << (max10bit + max12bit)) | (machineId << max12bit) | rand.Int64())
	return mist
}

/* 图片地址 */
func (Data) Img(img interface{}) string {
	str := (&util.Type{}).Strval(img)
	if str == "" {
		return ""
	}
	return config.Env().BaseURL + str
}

/*
* 分区-获取ID
* p2209 = (&service.Data{}).PartitionID("2022-10-01 00:00:00", "logs", "")
 */
func (r Data) PartitionID(date string, table string, column string) map[string]interface{} {
	if column == "" {
		column = "ctime"
	}
	t := util.Time()
	m := (&model.Model{})
	m.Table(table)
	m.Columns("id", column)
	m.Where(column+" < ?", t)
	m.Order(column + " DESC, id DESC")
	one := m.FindFirst()
	one["date"] = date
	one["time"] = t
	one["table"] = table
	return one
}

/*
* 分区-获取名称
* p = (&service.Data{}).PartitionName(1661961600, 1664553600)
 */
func (r Data) PartitionName(stime int, etime int) string {
	p1 := r.__getPartitionTime(stime)
	p2 := r.__getPartitionTime(etime)
	arr := []string{}
	start := false
	all := r.__getPartitionKeys()
	for _, k := range all {
		if k == p1 {
			start = true
		}
		if start {
			arr = append(arr, k)
		}
		if k == p2 {
			break
		}
	}
	return util.Implode(",", arr)
}
func (r Data) __getPartitionTime(time int) string {
	name := ""
	all := r.__getPartitionKeys()
	for _, k := range all {
		if time < partition[k] {
			return k
		}
		name = k
	}
	return name
}
func (r Data) __getPartitionKeys() []string {
	data := []string{}
	for k := range partition {
		data = append(data, k)
	}
	sort.Strings(data)
	return data
}
