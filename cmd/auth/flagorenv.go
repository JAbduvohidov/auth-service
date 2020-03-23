package main

import "os"

func FlagOrEnv(flag, key string) (string, bool) {
	if flag == "" {
		return os.LookupEnv(key)
	}
	return flag, true
}