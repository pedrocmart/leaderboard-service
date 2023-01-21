package utils

import "os"

func GetEnvOrDefault(key, defvalue string) (value string) {
	value, found := os.LookupEnv(key)
	if !found {
		value = defvalue
	}
	return
}
