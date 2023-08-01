package utils

import (
	// "context"
	"sync"
	"time"
)

const NanosecondsPerSecond = uint64(time.Second) / uint64(time.Nanosecond)
const NanosecondsPerMillisecond = uint64(time.Millisecond) / uint64(time.Nanosecond)
const TimeSecFormat string = "2006-01-02 15:04:05"
const TimeMinFormat string = "2006-01-02 15:04"
const TimeHourFormat string = "2006-01-02 15"
const TimeDayFormat string = "2006-01-02"
const TimeSecFormatWithT string = "2006-01-02T15:04:05"

//const TimeZoneChina string = "Asia/Shanghai"

// UnixNow 返回当前的 Unix 时间，使用 Second 作为单位
func UnixSecondNow() uint64 {
	var now int64 = time.Now().UnixNano() / int64(NanosecondsPerSecond)
	return uint64(now)
}

// UnixMilliNow 返回当前的 Unix 时间，使用 Millisecond 作为单位
func UnixMilliNow() uint64 {
	var now int64 = time.Now().UnixNano() / int64(NanosecondsPerMillisecond)
	return uint64(now)
}

// ParseUnixMilliTimestamp 将一个 Unix Millisecond 时间转换为 time.Time 结构体
func ParseUnixMilliTimestamp(timestamp uint64) time.Time {
	var t = time.Unix(0, 0)
	return t.Add(time.Duration(timestamp * NanosecondsPerMillisecond))
}

// UnixMilli 将一个 time.Time 结构体 转换为 Unix Millisecond Timestamp
func UnixMilli(t time.Time) uint64 {
	return uint64(t.UnixNano()) / NanosecondsPerMillisecond
}

// 获取时间戳所在当地日期的当地时间0点的时间戳
func GetTimeZoneDateUnixMilli(timestamp uint64, timeZoneStr string) uint64 {
	tz, _ := LoadLocation(timeZoneStr)
	t := ParseUnixMilliTimestamp(timestamp)
	t = t.In(tz)
	y, m, d := t.Date()
	return UnixMilli(time.Date(y, m, d, 0, 0, 0, 0, tz))
}

func Time2DateStrInTimeZone(timestamp uint64, timeZone string) string {
	location, _ := LoadLocation(timeZone)
	return ParseUnixMilliTimestamp(timestamp).In(location).Format("2006-01-02")
}
func Time2DateStrInTimeZoneFormat(timestamp uint64, timeZone string, format string) string {
	location, _ := LoadLocation(timeZone)
	return ParseUnixMilliTimestamp(timestamp).In(location).Format(format)
}

func Time2DateStrInLocal(timestamp uint64) string {
	return Time2DateStrInTimeZone(timestamp, "Local")
}

func TimeStr2TimeUint64InTimeZoneFormat(timestamp string, timeZone string, format string) (uint64, error) {
	tz, _ := LoadLocation(timeZone)
	res, err := time.ParseInLocation(format, timestamp, tz)
	if err != nil {
		return 0, err
	}
	return UnixMilli(res), nil
}

func Time2TimeUint64(timestamp string, timeZone string) (uint64, error) {
	tz, _ := LoadLocation(timeZone)
	res, err := time.ParseInLocation("2006-01-02 15:04:05", timestamp, tz)
	if err != nil {
		return 0, err
	}
	return UnixMilli(res), nil
}

func Time2TimeStrInTimeZone(timestamp uint64, timeZone string) string {
	if timestamp == 0 {
		return ""
	}
	location, _ := LoadLocation(timeZone)
	return ParseUnixMilliTimestamp(timestamp).In(location).Format("2006-01-02 15:04:05")
}

func Time2TimeStrInTimeZoneFormat(timestamp uint64, timeZone string, format string) string {
	location, _ := LoadLocation(timeZone)
	return ParseUnixMilliTimestamp(timestamp).In(location).Format(format)
}

// 导出文件名称
func Time2TimeStrInTimeZone2(timestamp uint64, timeZone string) string {
	location, _ := LoadLocation(timeZone)
	return ParseUnixMilliTimestamp(timestamp).In(location).Format("20060102150405")
}

func Time2TimeStrInLocal(timestamp uint64) string {
	return Time2TimeStrInTimeZone(timestamp, "Local")
}

var loadLocationResultMap sync.Map
var loadLocationMutex sync.Mutex

type loadLocationResult struct {
	location *time.Location
	err      error
}

// LoadLocation 实现了带有 cache 的 time.LoadLocation。尽量使用本函数而不是 time.LoadLocation 从而提升系统性能。
//
//
// Go 内置的 time.LoadLocation 会在每次被调用时读取本地的目录，速度比较慢，同时容易造成 too many open files 的错误。
// 可以参考下面的 issue 获取更多信息：
// https://github.com/golang/go/issues/24844
// https://github.com/golang/go/issues/26106
//
func LoadLocation(name string) (*time.Location, error) {
	result, ok := loadLocationResultMap.Load(name)
	if !ok {
		loadLocationMutex.Lock()
		defer loadLocationMutex.Unlock()

		result, ok = loadLocationResultMap.Load(name)
		if ok {
			result := result.(*loadLocationResult)
			return result.location, result.err
		}

		location, err := time.LoadLocation(name)
		result := &loadLocationResult{location: location, err: err}
		loadLocationResultMap.Store(name, result)
		if err != nil {
			return nil, err
		}
		return location, err
	} else {
		result := result.(*loadLocationResult)
		return result.location, result.err
	}
}

type WeekDate struct {
	ThisWeekStartDate time.Time
	ThisWeekEndDate   time.Time
	NextWeekStartDate time.Time
	NextWeekEndDate   time.Time
}

func GetWeekDate(now time.Time, timeZoneStr string) (weekDate *WeekDate) {
	// 获取这个星期的开始结束时间
	tz, _ := LoadLocation(timeZoneStr)
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	weekDate = new(WeekDate)
	weekDate.ThisWeekStartDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, tz).AddDate(0, 0, offset)
	weekDate.ThisWeekEndDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, tz).AddDate(0, 0, offset+7)

	weekDate.NextWeekStartDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, tz).AddDate(0, 0, offset+7)
	weekDate.NextWeekEndDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, tz).AddDate(0, 0, offset+7+7)

	return
}

// 获取所在时区今天的开始与结束
// 没有传时区默认为所在系统的时区
func GetTodayStartAndEndTimeStamp(timeZone string) (uint64, uint64) {
	loc, err := LoadLocation(timeZone)
	if err != nil {
		loc = time.Local
	}
	now := time.Now().In(loc)
	startTime := UnixMilli(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc))
	after := now.Add(time.Hour * 24)
	endTime := UnixMilli(time.Date(after.Year(), after.Month(), after.Day(), 0, 0, 0, 0, loc))
	return startTime, endTime - 1
}

// GetDateStartAndEndTimeStamp 获取时间戳所在当地日期的开始与结束
func GetDateStartAndEndTimeStamp(timestamp uint64, timeZoneStr string) (uint64, uint64) {
	tz, _ := LoadLocation(timeZoneStr)
	t := ParseUnixMilliTimestamp(timestamp)
	t = t.In(tz)
	y, m, d := t.Date()
	startTime := UnixMilli(time.Date(y, m, d, 0, 0, 0, 0, tz))
	after := t.AddDate(0, 0, 1)
	y, m, d = after.Date()
	endTime := UnixMilli(time.Date(y, m, d, 0, 0, 0, 0, tz))
	return startTime, endTime - 1
}

// 获取指定的开始与结束
func GetStartAndEndTimeStamp(timestamp uint64, timeZoneStr string) (uint64, uint64) {
	tz, _ := LoadLocation(timeZoneStr)
	t := ParseUnixMilliTimestamp(timestamp)
	t = t.In(tz)
	startTime := UnixMilli(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, tz))
	after := t.Add(time.Hour * 24)
	endTime := UnixMilli(time.Date(after.Year(), after.Month(), after.Day(), 0, 0, 0, 0, tz))
	return startTime, endTime - 1
}

// GetMonthStartAndEndTimeStamp 获取时间戳所在当月日期的开始与结束
func GetMonthStartAndEndTimeStamp(timestamp uint64, timeZone string) (uint64, uint64) {
	loc, err := LoadLocation(timeZone)
	if err != nil {
		loc = time.Local
	}
	t := ParseUnixMilliTimestamp(timestamp)
	t = t.In(loc)
	y, m, _ := t.Date()
	startTime := UnixMilli(time.Date(y, m, 1, 0, 0, 0, 0, loc))
	after := time.Date(y, m, 1, 0, 0, 0, 0, loc).AddDate(0, 1, 0)
	y, m, d := after.Date()
	endTime := UnixMilli(time.Date(y, m, d, 0, 0, 0, 0, loc))
	return startTime, endTime - 1
}

// 获取指定日期的0点时间
func GetZeroTime(t time.Time, timeZone string) time.Time {
	loc, err := LoadLocation(timeZone)
	if err != nil {
		loc = time.Local
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
}

//获取指定日期所属月份的第一天0点时间
func GetFirstDayOfMonth(t time.Time, timeZone string) time.Time {
	d := t.AddDate(0, 0, -t.Day()+1)
	return GetZeroTime(d, timeZone)
}

//获取指定日期所属月份的最后一天0点时间
func GetLastDayOfMonth(t time.Time, timeZone string) time.Time {
	return GetFirstDayOfMonth(t, timeZone).AddDate(0, 1, -1)
}

//获取当前周的周一
func GetMondayOfCurrentWeek(t time.Time) time.Time {
	var offset int
	if t.Weekday() == time.Sunday {
		offset = 7
	} else {
		offset = int(t.Weekday())
	}
	return t.AddDate(0, 0, -offset+1)
}

//获取指定日期的下周一和下周日
func GetNextWeek(t time.Time) (time.Time, time.Time) {
	mondayOfWeek := GetMondayOfCurrentWeek(t)      //当前时间的周一
	nextStartWeek := mondayOfWeek.AddDate(0, 0, 7) //下周一
	nextEndWeek := nextStartWeek.AddDate(0, 0, 6)  //下周日
	return nextStartWeek, nextEndWeek
}
