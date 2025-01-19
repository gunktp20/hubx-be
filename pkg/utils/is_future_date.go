package utils

import (
	"errors"
	"time"
)

func IsFutureDate(date time.Time) (bool, error) {

	// โหลดเขตเวลาไทย (Thailand Standard Time)
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return false, errors.New("failed to load Thailand time zone")
	}

	// ดึงวันที่ปัจจุบัน (แค่ปี, เดือน, วัน) ในเขตเวลาไทย
	now := time.Now().In(loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	// แปลง date ให้เป็นเขตเวลาไทย
	dateInThai := date.In(loc)

	// ตรวจสอบว่าวันที่เป็นในอดีตหรือวันนี้
	if dateInThai.Before(today) || dateInThai.Equal(today) {
		return false, errors.New("the event date cannot be today or in the past")
	}

	// เป็นวันที่ในอนาคต
	return true, nil
}
