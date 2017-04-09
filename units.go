package pbars

import "fmt"

// UnitFormatFunc is a function that takes in some decimal number of units and converts it
// to the appropriate unit.
type UnitFormatFunc func(v float64) string

// NoUnitFunc simply formats the value as a string, no unit conversion is performed
func NoUnitFunc(v float64) string {
	return fmt.Sprintf("%.2f", v)
}

var byteunits = []string{"B", "KB", "MB", "GB", "TB", "PB"}

// ByteFormatFunc formats bytes into B, KB, MB, GB..
func ByteFormatFunc(v float64) string {
	unit := byteunits[0]
	for i := 0; i < len(byteunits); i++ {
		unit = byteunits[i]
		if v < 1024.0 {
			break
		}
		v /= 1024
	}
	return fmt.Sprintf("%.2f%s", v, unit)
}

var _ UnitFormatFunc = NoUnitFunc
var _ UnitFormatFunc = ByteFormatFunc
