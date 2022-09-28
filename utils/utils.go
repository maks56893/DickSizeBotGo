package utils

import (
	models "DickSizeBot/postgres/models/dick_size"
	"context"
	"fmt"
	"math/rand"
	"time"
)

func GenerateDickSize() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(40)
}

func CheckLastMeasureDateIsToday(ctx context.Context, repo models.Repository, userid int64, chatId int64) bool {
	dickAndDate, err := repo.GetLastMeasureByUserInThisChat(ctx, userid, chatId)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	date := dickAndDate.Measure_date
	if date.Year() != 0001 {
		dateDiff := int(time.Now().Sub(date) / time.Hour)
		if dateDiff <= 24 {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func CheckIsTodayMeasure(measure models.DickSize) bool {
	today := time.Now()

	if measure.Measure_date.Year() == today.Year() && measure.Measure_date.Month() == today.Month() && measure.Measure_date.Day() == today.Day() {
		return true
	} else {
		return false
	}
}
