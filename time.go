package golibs

import "time"

const (
	Time_TIMEyyyyMMdd       string = "20060102"
	Time_TIMEStandard       string = "2006-01-02 15:04:05"
	Time_TIMEyyyyMMddHHmmss string = "20060102150405"
)

// Since返回从t到现在经过的毫秒数
func Since(t time.Time) int64 {
	return time.Since(t).Nanoseconds() / 1000000
}

// 返回当前时间戳（纳秒）
func UnixNano() int64 {
	return time.Now().UnixNano()
}

// 返回当前时间戳（秒）
func Unix() int64 {
	return time.Now().Unix()
}

// 返回当前时间字符串
func StandardTime() string {
	return time.Now().Format(Time_TIMEStandard)
}

// 返回从2000-01-01 00:00:00到现在经过的纳秒数
func From2000Nano() int64 {
	timeA, _ := time.Parse(Time_TIMEStandard, "2000-01-01 00:00:00")
	return time.Since(timeA).Nanoseconds()
}
