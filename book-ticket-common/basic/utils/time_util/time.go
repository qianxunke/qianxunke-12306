package time_util

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const Layout_Standard = "2006-01-02 15:04:05"
const Layout_For_Number = "20060102150405"

//时间戳字符串转标准时间
func TimestampStringToTimeString(cuostring string) (bzString string, err error) {
	if len(cuostring) > 10 {
		cuostring = cuostring[0:9]
	}
	value, err := strconv.ParseInt(cuostring, 10, 64)
	if err != nil {
		return
	}
	bzString = time.Unix(value, 0).Format("2006-01-02 15:04:05")
	return
}

//时间转时间戳
func TimeStringToTimestamp(datetime string) (timestamp int64, err error) {
	if len(datetime) <= 0 {
		err = errors.New("datetime is nil")
		return
	}
	timeLayout := "2006-01-02 15:04:05"    //转化所需模板
	loc, err := time.LoadLocation("Local") //获取时区
	if err != nil {
		return
	}
	tmp, err := time.ParseInLocation(timeLayout, datetime, loc)
	if err != nil {
		return
	}
	timestamp = tmp.Unix() //转化为时间戳 类型是int64
	return
}

//时间蹉转时间
func TimestampToTimeString(timestamp int64) (timeString string) {
	timeString = time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
	return
}

//带时区的时间转标准时间
func TimeZoneStringToTimeString(timeZoneString string) (timeString string, err error) {
	if len(timeZoneString) <= 0 {
		err = errors.New("timeString is nil")
		return
	}
	timeStamp, err := TimeStringToTimestamp(timeZoneString)
	if err != nil {
		return
	}
	timeString = TimestampToTimeString(timeStamp)
	return
}

//带时区的时间转标准时间
func TimeTtringToTimeString(timeZoneString string) (timeString string, err error) {
	if len(timeZoneString) == 0 {
		err = errors.New("timeString is nil")
	}
	if !strings.Contains(timeZoneString, "T") {
		timeString = timeZoneString
		return
	}
	//去掉T
	timeZoneString = strings.Replace(timeZoneString, "T", " ", -1)
	//
	timeZoneString = strings.Replace(timeZoneString, "+08:00", "", -1)

	timeString = timeZoneString
	return
}

//获取当前时间
func GetCurrentTime(layout string) (timeString string) {
	timeString = time.Unix(time.Now().Unix(), 0).Format(layout)
	return
}
