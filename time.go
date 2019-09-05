package golibs

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const (
	Time_TIMEyyyyMMdd           string = "20060102"
	Time_TIMEStandard           string = "2006-01-02 15:04:05"
	Time_TIME_HH_mm_ss          string = "15:04:05"
	Time_TIMEMSSQL              string = "2006-01-02T15:04:05.999Z"
	Time_TIMEMYSQL              string = "2006-01-02T15:04:05+08:00"
	Time_TIMEyyyyMMddHHmmss     string = "20060102150405"
	Time_TIMEyyyyMMddHHmmssffff string = "200601021504059999"
	Time_TIMEJavaUtilDate       string = "20060102150405000-0700"
	Time_TIMEISO8601            string = "2006-01-02T15:04:05.999-0700"
)

// 获取当前日期
func GetDate(t time.Time) int {
	return t.Year()*10000 + int(t.Month())*100 + t.Day()
}

// Since返回从t到现在经过的毫秒数
func Since(t time.Time) int64 {
	return time.Since(t).Nanoseconds() / 1000000
}

// 延时delay毫秒，从t开始计时
func Sleep(t time.Time, delay int64) {
	if end := delay - Since(t); end > 0 {
		time.Sleep(time.Duration(end) * time.Millisecond)
	}
}

// 返回当前时间戳（纳秒）
func UnixNano() int64 {
	return time.Now().UnixNano()
}

// 返回当前时间戳（秒）
func Unix() int64 {
	return time.Now().Unix()
}

// 返回当前时间戳（毫秒）
func UnixMilliSecond() int64 {
	return time.Now().UnixNano() / 1000000
}

// 返回当前时间字符串
func StandardTime() string {
	return time.Now().Format(Time_TIMEStandard)
}

// 获取当前时间的int64格式，yyyyMMddHHmmss
func GetTimeInt64() int64 {
	s := time.Now().Format(Time_TIMEyyyyMMddHHmmss)
	i64, _ := strconv.ParseInt(s, 10, 64)
	return i64
}

// 返回从2000-01-01 00:00:00到现在经过的纳秒数
func From2000Nano() int64 {
	time.Sleep(time.Nanosecond)
	timeA, _ := time.Parse(Time_TIMEStandard, "2000-01-01 00:00:00")
	return time.Since(timeA).Nanoseconds()
}

// 把时间字符串转成本地时间
func TimeStringToTime(sourceTime string) (time.Time, error) {
	return time.ParseInLocation(Time_TIMEStandard, sourceTime, time.Local)
}

//版本号转时间,格式:2019.905.1052
func Version2Time(version string) (t time.Time, err error) {

	//region 验证格式
	if Length(version) < 10 {
		return t, errors.New("错误的版本")
	}
	if !IsVersion(version) {
		return t, errors.New("错误的版本号")
	}
	if !strings.Contains(version, ".") {
		return t, errors.New("错误的版本号格式")
	}
	tmpVer := strings.Split(version, ".")
	if len(tmpVer) < 3 {
		return t, errors.New("错误的版本号格式")
	}
	//endregion

	//region 解析年份
	if !IsNumber(tmpVer[0]) {
		return t, errors.New("错误的版本号格式yyyy")
	}
	year, err := strconv.Atoi(tmpVer[0])
	if err != nil {
		return t, err
	}
	if year < 2018 || year > 2036 {
		return t, errors.New("错误的年份")
	}
	//endregion

	//region 解析月日
	if !IsNumber(tmpVer[1]) {
		return t, errors.New("错误的版本号格式MMdd")
	}
	monthDay, err := strconv.Atoi(tmpVer[1])
	if err != nil {
		return t, err
	}
	if monthDay < 101 || monthDay > 1231 {
		return t, errors.New("错误的月日")
	}
	month := monthDay / 100
	day := monthDay % 100
	if month < 1 || month > 12 {
		return t, errors.New("错误的月")
	}
	if day < 1 || day > 31 {
		return t, errors.New("错误的日")
	}
	if month == 2 && day > 29 {
		return t, errors.New("错误的日")
	}
	//endregion

	//region 解析时分
	if !IsNumber(tmpVer[2]) {
		return t, errors.New("错误的版本号格式HHmm")
	}
	hourMinute, err := strconv.Atoi(tmpVer[2])
	if err != nil {
		return t, err
	}
	if hourMinute < 0 || hourMinute > 2359 {
		return t, errors.New("错误的时分")
	}
	hour := hourMinute / 100
	minute := hourMinute % 100
	if hour < 0 || hour > 23 {
		return t, errors.New("错误的时")
	}
	if minute < 0 || minute > 59 {
		return t, errors.New("错误的分")
	}
	//endregion

	return time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.Local), nil

}
