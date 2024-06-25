package web

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
)

// squashSpace turn many spaces into one
func squashSpace(text string) string {
	return strings.TrimSpace(regexp.MustCompile(`\s{2,}`).ReplaceAllString(text, " "))
}

func handleParsingError(err error, meta string) {
	log.Printf("parse error [meta: %s]: %s", meta, err.Error())
}

func parseDate(date string, times string) (time.Time, error) {
	if date != "" && times != "" {
		loc, _ := time.LoadLocation("Europe/Moscow")
		localDate, _ := time.ParseInLocation(time.DateOnly, date, loc)
		fmt.Println(strings.Fields(times))
		t, _ := time.Parse(time.TimeOnly, strings.Fields(times)[3]+":00")

		return localDate.Add(time.Hour*time.Duration(t.Hour()) + time.Minute*time.Duration(t.Minute())), nil
	}

	return time.Time{}, fmt.Errorf("parseDate: empty input")
}
