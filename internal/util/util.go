package util

import (
	"strconv"
	"time"
)

func GetExpireTime(expDate string) time.Duration {
	// if exp date contains 'h' then return time.Hour
	if expDate[len(expDate)-1] == 'h' {
		// remove 'h' from expDate and convert to int
		expDate = expDate[:len(expDate)-1]
		expDateInt, _ := strconv.Atoi(expDate)
		return time.Hour * time.Duration(expDateInt)
	} else if expDate[len(expDate)-1] == 'd' {
		expDate = expDate[:len(expDate)-1]
		expDateInt, _ := strconv.Atoi(expDate)
		return time.Hour * time.Duration(expDateInt*24)
	} else if expDate[len(expDate)-1] == 's' {
		expDate = expDate[:len(expDate)-1]
		expDateInt, _ := strconv.Atoi(expDate)
		return time.Second * time.Duration(expDateInt)
	} else {
		return 0
	}
}
