package utils

import (
	models "DickSizeBot/postgres/models/dick_size"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
	"strconv"
	"time"
)

func CheckLastUsersDuelIsToday(ctx context.Context, repo models.Repository, userid int64, chatId int64) bool {
	duelDate, err := repo.GetLastDuelByUserId(ctx, userid, chatId)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	date := duelDate
	measureMidnight := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999, date.Location())
	if date.Year() != 0001 {
		if time.Now().Before(measureMidnight) {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func GenerateMsgTextTwoUsers(winner models.UserCredentials, loser models.UserCredentials, duel models.Duel) (msgText string) {
	msgText = "Победил в схватке диками — "

	if winner.Fname != "" {
		msgText += winner.Fname + " "
	}
	if winner.Username != "" {
		msgText += "@" + winner.Username + " "
	}
	if winner.Lname != "" {
		msgText += winner.Lname + " "
	}

	if winner.UserId == duel.CallerUserId {
		msgText += " выбросив " + strconv.Itoa(int(duel.CallerRoll)) + " "
	} else {
		msgText += " выбросив " + strconv.Itoa(int(duel.CalledRoll)) + " "
	}

	msgText += ", кок сакер — "

	if loser.Fname != "" {
		msgText += loser.Fname + " "
	}
	if loser.Username != "" {
		msgText += "@" + loser.Username + " "
	}
	if loser.Lname != "" {
		msgText += loser.Lname + " "
	}

	if winner.UserId != duel.CallerUserId {
		msgText += " его ролл - " + strconv.Itoa(int(duel.CallerRoll)) + " "
	} else {
		msgText += " выбросив " + strconv.Itoa(int(duel.CalledRoll)) + " "
	}

	return msgText
}

func GenerateDickSize() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(39) + 1
}

func GetDuelWinner(callerId int64, calledId int64) (winner int64, callerRoll int, calledRoll int) {
	rand.Seed(time.Now().UnixNano())

	callerRoll = rand.Intn(100)
	time.Sleep(10 * time.Millisecond)
	calledRoll = rand.Intn(100)

	if callerRoll >= calledRoll {
		winner = callerId
		return winner, callerRoll, calledRoll
	} else {
		winner = calledId
		return winner, callerRoll, calledRoll
	}
}

func CheckLastMeasureDateIsToday(ctx context.Context, repo models.Repository, userid int64, chatId int64) bool {
	dickAndDate, err := repo.GetLastMeasureByUserInThisChat(ctx, userid, chatId)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	date := dickAndDate.Measure_date
	measureMidnight := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999, date.Location())
	if date.Year() != 0001 {
		if time.Now().Before(measureMidnight) {
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

func AddRowToInlineKeyboard(keyboardMarkup *tgbotapi.InlineKeyboardMarkup, row []tgbotapi.InlineKeyboardButton) tgbotapi.InlineKeyboardMarkup {
	newKeyboard := append(keyboardMarkup.InlineKeyboard, row)
	return tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: newKeyboard,
	}
}
