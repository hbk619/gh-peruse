package internal

import "time"

func TimeMustParse(layout string, str string) time.Time {
	parse, err := time.Parse(layout, str)
	if err != nil {
		panic(err)
	}
	return parse
}
