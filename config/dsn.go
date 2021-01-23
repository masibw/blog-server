package config

import (
	"fmt"
	"os"
)

func DSN() string {
	return PureDSN() + "?parseTime=true&loc=Asia%2FTokyo&collation=utf8mb4_bin"
}

func PureDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("DB_HOST"),
		"3306",
		os.Getenv("MYSQL_DATABASE"),
	)
}
