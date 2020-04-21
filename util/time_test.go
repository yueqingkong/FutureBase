package util

import (
	"testing"
)

func TestStringYearToTime(t *testing.T) {
	s := "2019-11-05"
	ti:= StringYearToTime(s)
	t.Log(t, "show Time: ", ti.Format("2006-01-02 15:04:05"))
}
