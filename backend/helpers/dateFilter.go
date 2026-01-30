package helpers

import (
	"time"

	"github.com/gin-gonic/gin"
)

type DateFilter struct {
	Start time.Time
	End   time.Time
	Label string
}

func ParseDateFilter(c *gin.Context) DateFilter {
	now := time.Now()

	startStr := c.Query("start_date")
	endStr := c.Query("end_date")
	rangeParam := c.Query("range")

	if startStr != "" && endStr != "" {
		start, _ := time.Parse("2006-01-02", startStr)
		end, _ := time.Parse("2006-01-02", endStr)

		return DateFilter{
			Start: start,
			End:   end.Add(23*time.Hour + 59*time.Minute + 59*time.Second),
			Label: "custom",
		}
	}

	switch rangeParam {
	case "today":
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		return DateFilter{Start: start, End: now, Label: "today"}

	case "this_month":
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		return DateFilter{Start: start, End: now, Label: "this_month"}

	case "this_year":
		start := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		return DateFilter{Start: start, End: now, Label: "this_year"}

	default:
		start := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		return DateFilter{Start: start, End: now, Label: "ytd"}
	}
}