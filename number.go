package golibs

import "fmt"

func ByteToGByte(byte_value float64) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	var mod float64 = 1024.0
	var index int = 0
	for byte_value >= mod {
		byte_value /= mod
		index++
	}
	return fmt.Sprintf("%.4f", byte_value) + units[index]
}
