package utils

import "os"

func GetEnv(key string, default_ string) string {
	val, present := os.LookupEnv(key)
	if present {
		return val
	} else {
		return default_
	}
}
