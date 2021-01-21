package config

import "os"

func IsLocal() bool {
	return os.Getenv("ENV") == "local"
}
