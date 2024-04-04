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
