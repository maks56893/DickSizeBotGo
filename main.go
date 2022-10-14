package main

import (
	cash2 "DickSizeBot/cash"
	"DickSizeBot/cash_domain"
	. "DickSizeBot/logger"
	"DickSizeBot/postgres"
	models "DickSizeBot/postgres/models/dick_size"
	"DickSizeBot/postgres/models/dick_size/db"
	"DickSizeBot/utils"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const CommandToBot = "@FatBigDickBot"
const CommandToTestingBot = "@TestingDickSizeBot"

const MeasureCommand = "/check_size"
const AverageCommand = "/get_average"
const TodayCommand = "/last_measures"

//var numericKeyboard = tgbotapi.NewReplyKeyboard(
//	tgbotapi.NewKeyboardButtonRow(
//		tgbotapi.NewKeyboardButton(MeasureCommand),
//		tgbotapi.NewKeyboardButton(AverageCommand),
//		//		tgbotapi.NewKeyboardButton("3"),
//	),
//)
var removeKeyboard = tgbotapi.NewRemoveKeyboard(true)

func main() {
	LoggerInit("trace", "log/bot-log.log", true)
	err := tgbotapi.SetLogger(Log)
	if err != nil {
		return
	}

	//bot, err := tgbotapi.NewBotAPI("5445796005:AAHQLY5pFGMOZ_uVbEzel0tK0dRReIVC7bw") //main bot
	bot, err := tgbotapi.NewBotAPI("5681105337:AAHNnD0p6XcXo7biy9U7F7P-ctSkk-TrWGA") //test bot
	if err != nil {
		log.Panic(err)
	}

	ctx := context.Background()

	user := "postgres"
	pass := "56893"
	host := "localhost"
	database := "postgres"

	client, err := postgres.NewClient(ctx, 2, user, pass, host, "5432", database)
	if err != nil {
		Log.Println(err.Error())
	}

	cash := cash2.NewCash().(cash_domain.ICash)

	repo := db.NewRepo(client)

	bot.Debug = true
	Log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	go func() {
		for {
			time.Sleep(1 * time.Minute)
			weekCount := 0
			if time.Now().Weekday().String() == "Monday" {
				if time.Now().Hour() == 4 && time.Now().Minute() == 00 {
					weekCount++
					if weekCount == 2 {
						repo.DeleteSizesByTime(ctx)
						log.Printf("Database was cleared at: %v", time.Now())
						weekCount = 0
					}

				}
			}
		}

	}()

	//	var msgToDel int
	inputAttempts := 0

	for update := range updates {

		var msg tgbotapi.MessageConfig

		var duelDataFromCallback interface{}
		var hasCash bool

		if update.Message != nil {
			key := "duel_" + strconv.Itoa(int(update.Message.From.ID))
			duelDataFromCallback, hasCash = cash.Get(key)
		}

		if update.Message != nil && !hasCash { // If we got a message

			switch update.Message.Text {
			case MeasureCommand, MeasureCommand + CommandToBot, MeasureCommand + CommandToTestingBot:

				if utils.CheckLastMeasureDateIsToday(ctx, repo, update.Message.From.ID, update.Message.Chat.ID) {
					_ = repo.CreateOrUpdateUser(ctx, update.Message.From.ID, update.Message.From.FirstName, update.Message.From.LastName, update.Message.From.UserName, update.Message.Chat.ID)
					dickModel, _ := repo.GetLastMeasureByUserInThisChat(ctx, update.Message.From.ID, update.Message.Chat.ID)

					msg = tgbotapi.NewMessage(update.Message.Chat.ID, GetRandMeasureReplyPattern(int(dickModel.Dick_size)))
					msg.ReplyMarkup = removeKeyboard
					msg.ReplyToMessageID = update.Message.MessageID
				} else {
					dickSize := utils.GenerateDickSize()

					_, err := repo.InsertSize(ctx, update.Message.From.ID, update.Message.From.FirstName, update.Message.From.LastName, update.Message.From.UserName, dickSize, update.Message.Chat.ID, update.Message.Chat.IsGroup())
					if err != nil {
						Log.Printf(err.Error())
					}

					msg = tgbotapi.NewMessage(update.Message.Chat.ID, GetRandMeasureReplyPattern(dickSize))

					msg.ReplyMarkup = removeKeyboard
					msg.ReplyToMessageID = update.Message.MessageID
				}
			case AverageCommand, AverageCommand + CommandToBot, AverageCommand + CommandToTestingBot:
				if update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup() {
					chatAverages, _ := repo.GetUserAllSizesByChatId(ctx, update.Message.Chat.ID)

					msgText := GetRandAverageRepltText()

					for _, userData := range chatAverages {
						if fname, ok := userData["fname"]; ok {
							msgText += "· "
							msgText += fname
							msgText += " "
						}
						if username, ok := userData["username"]; ok && userData["username"] != "" {
							msgText += "@"
							msgText += username
							msgText += " "
						}
						if lname, ok := userData["lname"]; ok {
							msgText += lname
							msgText += " "
						}
						if average, ok := userData["average"]; ok {
							msgText += average
							msgText += " см"
							msgText += "\n"
						}
					}

					msg = tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
					msg.ReplyMarkup = removeKeyboard
					msg.ParseMode = "HTML"
					msg.ReplyToMessageID = update.Message.MessageID

				} else {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Используй в группе")
					msg.ReplyMarkup = removeKeyboard
					msg.ParseMode = "HTML"

					msg.ReplyToMessageID = update.Message.MessageID
				}
			case TodayCommand, TodayCommand + CommandToBot, TodayCommand + CommandToTestingBot:
				if update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup() {
					todayMeasures, err := repo.SelectOnlyTodaysMeasures(ctx, update.Message.Chat.ID)
					if todayMeasures != nil {
						Log.Errorf("Error while getting today measures: %s", err)

						msgText := GetRandTodayReplyText()

						for _, measure := range todayMeasures {
							if utils.CheckIsTodayMeasure(measure) {
								msgText += "✅ "
							} else {
								msgText += "❗ "
							}
							if measure.Fname != "" {
								msgText += measure.Fname + " "
							}
							if measure.Username != "" {
								msgText += "@" + measure.Username + " "
							}
							if measure.Lname != "" {
								msgText += measure.Lname + " "
							}
							if measure.Dick_size != 0 {
								msgText += strconv.Itoa(int(measure.Dick_size)) + "см"
							}
							if !utils.CheckIsTodayMeasure(measure) {
								measureDay := measure.Measure_date.Day()
								measureMonth := int(measure.Measure_date.Month())
								measureYear := measure.Measure_date.Year()
								msgText += " <i>(отмерено " + strconv.Itoa(measureDay) + "." + strconv.Itoa(measureMonth) + "." + strconv.Itoa(measureYear%100/10) + strconv.Itoa(measureYear%10) + ")" + "</i>\n"
							} else {
								msgText += "\n"
							}
						}

						msg = tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
						msg.ReplyMarkup = removeKeyboard
						msg.ParseMode = "HTML"
						msg.ReplyToMessageID = update.Message.MessageID
					}
				} else {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Используй в группе")
					msg.ReplyMarkup = removeKeyboard
					msg.ParseMode = "HTML"

					msg.ReplyToMessageID = update.Message.MessageID
				}
			case "/test":

				//				callbackData := ""
				userData := repo.GetAllCredentials(ctx, update.Message.Chat.ID)
				for data := range userData {
					fmt.Println(data)
				}
				var usersKeyboardButtons = tgbotapi.NewInlineKeyboardMarkup()
				for _, userCred := range userData {
					userId := "duel_user_id " + strconv.Itoa(int(userCred.UserId))

					buttonText := userCred.Fname + " @" + userCred.Username + " " + userCred.Lname

					row := tgbotapi.NewInlineKeyboardRow(
						tgbotapi.InlineKeyboardButton{
							Text:         buttonText,
							CallbackData: &userId,
						},
					)
					usersKeyboardButtons = utils.AddRowToInlineKeyboard(&usersKeyboardButtons, row)
				}

				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "С кем хочешь помериться?")
				msg.ReplyMarkup = usersKeyboardButtons
				msg.ParseMode = "HTML"
				msg.ReplyToMessageID = update.Message.MessageID

			}
			message, err := bot.Send(msg)
			if err != nil {
				Log.Printf(err.Error())
			}
			//			msgToDel = message.MessageID
			Log.Debugf("Sended message: %v", message)

			//Обработка ввода ставки после callback и с данными к кеше
		} else if update.Message != nil && hasCash {

			duelData := models.Duel{
				CallerUserId: duelDataFromCallback.(models.Duel).CallerUserId,
				CalledUserId: duelDataFromCallback.(models.Duel).CalledUserId,
				ChatID:       duelDataFromCallback.(models.Duel).ChatID,
			}

			if inputAttempts < 2 {
				//Поиск числа
				findCorrectNum, err := regexp.MatchString("[1-5]", update.Message.Text)
				if err != nil {
					Log.Errorf("Error while trying find number in text for duel??????????: %v", err)
				}
				//Поиск отмены
				findCancel, err := regexp.MatchString("саси", update.Message.Text)
				if err != nil {
					Log.Errorf("Error while trying find number in text for duel??????????: %v", err)
				}

				if findCorrectNum {
					duelData.Bet, _ = strconv.Atoi(update.Message.Text)
					duelData.Winner = utils.GetDuelWinner(duelData.CalledUserId, duelData.CallerUserId)

					_ = repo.InsertDuelData(ctx, duelData)
					//					msgText := fmt.Sprintf("Вызвал на дуэль %v, кого вызвали на дуэль: %v, ставка: %d, подебитель: %d", duelData.CallerUserId, duelData.CalledUserId, duelData.Bet, duelData.Winner)
					//					msg = tgbotapi.NewMessage(update.Message.Chat.ID, msgText)

					//Данные победителя и проигравшего
					var winner, loser models.UserCredentials

					if duelData.Winner == duelData.CallerUserId {
						winner = repo.GetUserData(ctx, duelData.CallerUserId)
						loser = repo.GetUserData(ctx, duelData.CalledUserId)
					} else {
						winner = repo.GetUserData(ctx, duelData.CalledUserId)
						loser = repo.GetUserData(ctx, duelData.CallerUserId)
					}

					//Обновляем последний размер пользователя
					winnerDickSize, err := repo.GetLastMeasureByUserInThisChat(ctx, winner.UserId, duelData.ChatID)
					if err != nil {
						Log.Errorf("Error while get last measure by user: %v", err)
					}
					loserDickSize, err := repo.GetLastMeasureByUserInThisChat(ctx, loser.UserId, duelData.ChatID)
					if err != nil {
						Log.Errorf("Error while get last measure by user: %v", err)
					}

					if &winnerDickSize != nil && &loserDickSize != nil {
						repo.IncreaceLastDickSize(ctx, winnerDickSize.Id, duelData.Bet)
						repo.IncreaceLastDickSize(ctx, loserDickSize.Id, duelData.Bet-duelData.Bet*2)
					}

					//Формируем сообщение о дуэли
					msgText := fmt.Sprintf("Победил %s @%s %s, кок сакер %s @%s %s", winner.Fname, winner.Username, winner.Lname, loser.Fname, loser.Username, loser.Lname)
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, msgText)

					key := "duel_" + strconv.Itoa(int(update.Message.From.ID))
					cash.Del(key)
				} else if findCancel {

					//Просто удаляем кеш, так как получили команду отмены
					key := "duel_" + strconv.Itoa(int(update.Message.From.ID))
					cash.Del(key)
				} else {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Давай попробуем ввести число от 1 до 5 ещё раз :)")
					inputAttempts++
				}

			}

			if inputAttempts == 2 {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Поздравляю, кок сакер, ты не смог с двух попыток ввести число. Начинай все заново")
				inputAttempts = 0
			}
			message, err := bot.Send(msg)
			if err != nil {
				Log.Printf(err.Error())
			}

			//			msgToDel = message.MessageID
			Log.Debugf("Sended message: %v", message)

		} else if update.CallbackQuery != nil {
			callbackData := strings.Split(update.CallbackQuery.Data, " ")
			if callbackData[0] == "duel_user_id" {
				callback := tgbotapi.NewCallback(update.CallbackQuery.ID /*update.CallbackQuery.Data*/, "duel callback")

				if _, err := bot.Request(callback); err != nil {
					Log.Errorf("Error while exec remove msg request: %v", err)
				}

				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID /*update.CallbackQuery.Data*/, "Сколько сантиметров хочешь поставить? Если передумал, то введи \"саси\"")
				//				msg.ReplyMarkup = betKeyboard
				if _, err := bot.Send(msg); err != nil {
					Log.Errorf("Error while sending msg after callback: %v", err)
				}

				calledId, _ := strconv.Atoi(callbackData[1])

				dataForDuel := models.Duel{
					CallerUserId: update.CallbackQuery.From.ID,
					CalledUserId: int64(calledId),
					ChatID:       update.CallbackQuery.Message.Chat.ID,
					//					Bet:          "",
				}

				key := "duel_" + strconv.Itoa(int(update.CallbackQuery.From.ID))
				cash.Set(key, dataForDuel, 1*time.Minute)

			}
		}
	}

}
