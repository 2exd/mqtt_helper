package utils

import "time"

func GetHhmmss() string {
	currentTime := time.Now()
	return currentTime.Format("15:04:05")
}

func GetHh_mm_ss() string {
	currentTime := time.Now()
	return currentTime.Format("15_04_05")
}
