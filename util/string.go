package util

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func StringToFloat32(str string) float32 {
	value, err := strconv.ParseFloat(str, 32)
	if err != nil {
		value = 0.0
	}

	return float32(value)
}

func StringToInt64(str string) int64 {
	i, err := strconv.Atoi(str)
	if err != nil {
		log.Print(err)
	}

	return int64(i)
}

func Float32ToString(f float32) string {
	return strconv.FormatFloat(float64(f), 'f', 6, 64)
}

func FloatDeceimal(f float32) float32 {
	value, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", f), 32)
	return float32(value)
}

/**
 * 大写
 */
func Upper(s string) string  {
	return strings.ToUpper(s)
}

/**
 * 小写
 */
func Lower(s string) string  {
	return strings.ToLower(s)
}
