package dateUtil

import (
	"strconv"
	"time"
)

const DefaultDate =  "1970-01-01 00:00:00"
const DefaultLocationLayout = "2006-01-02T15:04:05+08:00"
const DefaultLayout = "2006-01-02 15:04:05"

func ParseInLocationDefault(date string) time.Time {
	if date == "" {
		date = DefaultDate
	}
	loc, _ := time.LoadLocation("Local")
	theTime, err := time.ParseInLocation(DefaultLocationLayout, date, loc)
	if err != nil{
		date = DefaultDate
		theTime, _ = time.ParseInLocation(DefaultLayout, date, loc)
	}
	return theTime
}

func ToDateString(time time.Time) string {
	return time.Format(DefaultLayout)
}

func Now() string {
	return ToDateString(time.Now())
}

const BirthdayLayout = "20060102"

func IsKidForIdcard(idcard string) bool{
	if len(idcard) != 18 {
		return false
	}

	birthday := idcard[6:14]
	yearStr := birthday[0:4]
	year,err := strconv.Atoi(yearStr)
	if err != nil{
		return false
	}

	newTime, err := time.Parse(BirthdayLayout, strconv.Itoa(year + 18) + birthday[4:])
	if err != nil{
		return false
	}

	if newTime.Before(time.Now()){
		return true
	}

	return false
}