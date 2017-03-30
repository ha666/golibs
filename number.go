package golibs

import "fmt"

func ByteToGByte(byteValue float64) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	var mod float64 = 1024.0
	var index int = 0
	for byteValue >= mod {
		byteValue /= mod
		index++
	}
	return fmt.Sprintf("%.4f", byteValue) + units[index]
}
