package globals

import (
	"database/sql"

	zLogger "github.com/iabdulzahid/go-logger/logger"
)

var Logger zLogger.Logger
var DB *sql.DB
