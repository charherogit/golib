package helper

import (
	"strconv"
	"time"
)

const (
	RequestTimeout  = 30 * time.Second          // 需要加长的请求超时时间
	DaySeconds      = 24 * 3600                 // 每天的总秒数
	HourSeconds     = 3600                      // 每小时的总秒数
	MinuteSeconds   = 60                        // 每分钟的总秒数
	ReadDeadline    = 60 * time.Second          // XXX 客户端心跳30s，这块设置为60s俩个心跳时间
	TimeFormatMonth = "2006-01"                 // 年月
	TimeFormatDay   = "2006-01-02"              // 年月日
	TimeFormatStamp = "2006-01-02 15:04:05"     // 年月日时分秒
	TimeFormatMilli = "2006-01-02 15:04:05.000" // 年月日时分秒毫秒
)

// ResolveTime 秒数转化为天小时分钟
func ResolveTime(seconds int) (day int, hour int, minute int) {
	day = seconds / DaySeconds
	hour = (seconds - day*DaySeconds) / HourSeconds
	minute = (seconds - day*DaySeconds - hour*HourSeconds) / MinuteSeconds
	return
}

// GetTodayBeginTimestamp 获取今天的开始时间戳
func GetTodayBeginTimestamp() int64 {
	t := time.Now()
	return t.Unix() - int64(t.Hour())*HourSeconds - int64(t.Minute())*MinuteSeconds - int64(t.Second())
}

// GetTimeDiffFromNow 获取与当前时间的差
func GetTimeDiffFromNow(timestamp int64) int64 {
	return time.Now().Unix() - timestamp
}

// GetTodayBeginTime 获得今天0点0分0秒时间
func GetTodayBeginTime() time.Time {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// GetDayBeginTimestamp 获得某天的开始时间戳
func GetDayBeginTimestamp(day int64) int64 {
	t := time.Unix(day, 0)
	return t.Unix() - int64(t.Hour())*HourSeconds - int64(t.Minute())*MinuteSeconds - int64(t.Second())
}

// GetDayEndTimestamp 获得某天的结束时间戳
func GetDayEndTimestamp(day int64) int64 {
	return GetDayBeginTimestamp(day) + DaySeconds
}

// GetTodayEndTimestamp 获取今天的结束时间戳
func GetTodayEndTimestamp() int64 {
	return GetTodayBeginTimestamp() + DaySeconds
}

// GetNextDayBeginTimestamp 获取明天的开始时间戳
func GetNextDayBeginTimestamp() int64 {
	return GetTodayBeginTimestamp() + DaySeconds
}

func GetCurWeekBeginTime() time.Time {
	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}

	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
}

func GetNextWeekBeginTimestamp() int64 {
	return GetCurWeekBeginTime().AddDate(0, 0, 7).Unix()
}

func GetCureMothBeginTime() time.Time {
	now := time.Now()
	d := now.AddDate(0, 0, -now.Day()+1)
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

func GetNextMothBeginTimestamp() int64 {
	return GetCureMothBeginTime().AddDate(0, 1, 0).Unix()
}

// 获取n个月之后的开始时间戳
func GetLaterNMonthBeginTimestamp(count int) int64 {
	return GetCureMothBeginTime().AddDate(0, count, 0).Unix()
}

// IsSameDayCompareCurTime 和当前时间对比是否是同一天
func IsSameDayCompareCurTime(timeStamp int64) bool {
	if timeStamp <= 0 {
		return false
	}
	compareTime := time.Unix(timeStamp, 0)
	curTime := time.Now()
	return compareTime.YearDay() == curTime.YearDay() && compareTime.Year() == curTime.Year()
}

// 对比两个时间戳是否是同一天
func IsSameDayCompareTime(timeStamp int64, compareTime int64) bool {
	if timeStamp <= 0 || compareTime <= 0 {
		return false
	}
	time1 := time.Unix(timeStamp, 0)
	time2 := time.Unix(compareTime, 0)
	return time1.YearDay() == time2.YearDay() && time1.Year() == time2.Year()
}

// StringToTime 2006-01-02
func StringToTime(t string) time.Time {
	y, _ := strconv.Atoi(t[:4])
	m, _ := strconv.Atoi(t[5:7])
	d, _ := strconv.Atoi(t[8:])
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
}

func TimeToMomentString(ti time.Time) string {
	s := ti.Format("2006-01-02 15")
	if ti.Minute() < 15 {
		s += ":00"
	} else if ti.Minute() >= 15 && ti.Minute() < 30 {
		s += ":15"
	} else if ti.Minute() >= 30 && ti.Minute() < 45 {
		s += ":30"
	} else if ti.Minute() >= 45 && ti.Minute() < 60 {
		s += ":45"
	}
	return s
}

// 修改时区为本地时区
func TimeToString(t time.Time) string {
	return t.In(time.Local).Format(TimeFormatStamp)
}

func GetTimeFormatDayByTimestamp(t int64) string {
	return time.Unix(t, 0).Format(TimeFormatDay)
}

func GetTimeFormatDayByTime(t time.Time) string {
	return t.In(time.Local).Format(TimeFormatDay)
}

func GetTimeByTimestamp(t int64) time.Time {
	return time.Unix(t, 0)
}

// 根据天数返回其秒数
func GetSecondsByDays(days uint32) int64 {
	return DaySeconds * int64(days)
}

func GetGatGap(x, y int64) int {
	if x < y {
		return 0
	}
	return int((x - y) / int64(DaySeconds))
}

// GetMinuteSecondsGap 时间误差
func GetMinuteSecondsGap(src, dst int64) int64 {
	return src / dst * dst
}

func GetTimeDurationBySecond(t int64) time.Duration {
	return time.Duration(t) * time.Second
}
