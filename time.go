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

func UnixNano() int64 {
	return time.Now().UnixNano()
}

func Unix() int64 {
	return time.Now().Unix()
}
