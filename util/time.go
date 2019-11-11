package util

import (
	"errors"
	"strconv"
	"time"
)

func IsoTime(t time.Time) string {
	iso := t.UTC().String()
	isoBytes := []byte(iso)
	iso = string(isoBytes[:10]) + "T" + string(isoBytes[11:23]) + "Z"
	return iso
}

func IsoToTime(iso string) (time.Time, error) {
	nilTime := time.Now()
	if iso == "" {
		return nilTime, errors.New("illegal parameter")
	}
	// "2018-03-18T06:51:05.933Z"
	isoBytes := []byte(iso)
	year, err := strconv.Atoi(string(isoBytes[0:4]))
	if err != nil {
		return nilTime, errors.New("illegal year")
	}
	month, err := strconv.Atoi(string(isoBytes[5:7]))
	if err != nil {
		return nilTime, errors.New("illegal month")
	}
	day, err := strconv.Atoi(string(isoBytes[8:10]))
	if err != nil {
		return nilTime, errors.New("illegal day")
	}
	hour, err := strconv.Atoi(string(isoBytes[11:13]))
	if err != nil {
		return nilTime, errors.New("illegal hour")
	}
	min, err := strconv.Atoi(string(isoBytes[14:16]))
	if err != nil {
		return nilTime, errors.New("illegal min")
	}
	sec, err := strconv.Atoi(string(isoBytes[17:19]))
	if err != nil {
		return nilTime, errors.New("illegal sec")
	}
	nsec, err := strconv.Atoi(string(isoBytes[20 : len(isoBytes)-1]))
	if err != nil {
		return nilTime, errors.New("illegal nsec")
	}

	return time.Date(year, time.Month(month), day, hour, min, sec, nsec, time.UTC), nil
}

func UnixMillis(t time.Time) string {
	return strconv.FormatInt(t.UnixNano()/int64(time.Millisecond), 10)
}

func TimestampToTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

func TimeToString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func StringToTime(str string) time.Time {
	var result time.Time

	local, _ := time.LoadLocation("Local")
	t, err := time.ParseInLocation("2006-01-02 15:04:05", str, local)
	if err != nil {
		result = time.Time{}
	} else {
		result = t
	}
	return result
}

func StringYearToTime(str string) time.Time {
	var result time.Time

	local, _ := time.LoadLocation("Local")
	t, err := time.ParseInLocation("2006-01-02", str, local)
	if err != nil {
		result = time.Time{}
	} else {
		result = t
	}
	return result
}

func SubHours(t1 time.Time, t2 time.Time) float32 {
	return float32(t1.Sub(t2).Hours())
}

/**
 * 添加秒
 */
func AddSecond(t time.Time, second int64) time.Time {
	return t.Add(time.Duration(second) * time.Second)
}
