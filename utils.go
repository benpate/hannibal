package hannibal

import (
	"net/http"
	"time"
)

func TimeFormat(value time.Time) string {
	return value.UTC().Format(http.TimeFormat)
}
