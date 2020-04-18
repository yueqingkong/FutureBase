package plat

import (
	"errors"
	"log"
	"strconv"
	"time"
)

var (
	format = "2006-01-02 15:04:05"
)

func IsoTime(t time.Time) string {
	iso := t.UTC().String()
	isoBytes := []byte(iso)
	iso = string(isoBytes[:10]) + "T" + string(isoBytes[11:19]) + "Z"
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

func StringToTime(t string) time.Time {
	result, err := time.Parse(format, t)
	if err != nil {
		log.Print(err)
	}

	return result
}
