package main

import "testing"

func TestParseShifts(t *testing.T) {
	parseShifts("1-2=3-4,3-4=1-2", 4)
}
