package Helpers

import "time"

func TimeString() string {
	t := time.Now()
	return t.Format("15:04:05")
}
