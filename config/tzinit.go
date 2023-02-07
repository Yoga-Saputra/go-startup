package config

import "os"

func init() {
	os.Setenv("TZ", "Asia/Manila")
}
