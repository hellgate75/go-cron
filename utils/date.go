package utils

import "time"


// Calculates the next year since Now
func NextYear() time.Time {
	return time.Now().AddDate(1, 0, 0)
}

// Calculates new date since now adding given number of years
func AddYears(m int) time.Time {
	return time.Now().AddDate(m, 0, 0)
}

// Calculates the next year since given date
func NextYearSince(date time.Time) time.Time {
	return date.AddDate(1, 0, 0)
}

// Calculates new date since given date adding given number of years
func AddYearsSince(date time.Time, m int) time.Time {
	return date.AddDate(m, 0, 0)
}


// Calculates the next month since Now
func NextMonth() time.Time {
	return time.Now().AddDate(0, 1, 0)
}

// Calculates new date since now adding given number of months
func AddMonths(m int) time.Time {
	return time.Now().AddDate(0, m, 0)
}

// Calculates the next month since given date
func NextMonthSince(date time.Time) time.Time {
	return date.AddDate(0, 1, 0)
}

// Calculates new date since given date adding given number of months
func AddMonthsSince(date time.Time, m int) time.Time {
	return date.AddDate(0, m, 0)
}

