package date_utils

const Iso8601DateFormat = "%Y-%m-%d"

func MonthNames() struct {
	return {
		"01": "jan",
		"02": "feb",
		"03": "mar",
		"04": "apr",
		"05": "may",
		"06": "jun",
		"07": "jul",
		"08": "aug",
		"09": "sep",
		"10": "oct",
		"11": "nov",
		"12": "dec",
	}
}

func MonthShortNames() slice {
	return [
		"jan",
		"feb",
		"mar",
		"apr",
		"may",
		"jun",
		"jul",
		"aug",
		"sep",
		"oct",
		"nov",
		"dec",
	]
}

// Get date parts
func GetDateParts(iso8601Time string) map[string][]int {
	if iso8601Time == nil {
		return map[string][]int{"date-parts": []int{}}
	}
	if len(iso8601Time) < 10 {
		iso8601Time = strings.Repeat(iso8601Time, 10, "0")
	}
	year, _ := strconv.Atoi(iso8601Time[0:4])
	month, _ := strconv.Atoi(iso8601Time[5:7])
	day, _ := strconv.Atoi(iso8601Time[8:10])
	dateParts := []int{year, month, day}
	return map[string][]int{"date-parts": dateParts}
}

func GetDateFromUnixTimestamp(timestamp int) string {
	if timestamp == nil {
		return nil
	}
	return time.Unix(timestamp, 0).Format(Iso8601DateFormat)
}
