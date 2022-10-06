package main

import (
	. "DickSizeBot/logger"
	"DickSizeBot/postgres"
	"DickSizeBot/postgres/models/dick_size/db"
	"DickSizeBot/utils"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

const CommandToBot = "@FatBigDickBot"
const CommandToTestingBot = "@TestingDickSizeBot"

const MeasureCommand = "/check_size"
const AverageCommangd = "/get_average"
const TodayCommand = "/last_measures"

//var numericKeyboard = tgbotapi.NewReplyKeyboard(
//	tgbotapi.NewKeyboardButtonRow(
//		tgbotapi.NewKeyboardButton(MeasureCommand),
//		tgbotapi.NewKeyboardButton(AverageCommangd),
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

	bot, err := tgbotapi.NewBotAPI("5445796005:AAHQLY5pFGMOZ_uVbEzel0tK0dRReIVC7bw") //main bot
	// bot, err := tgbotapi.NewBotAPI("5681105337:AAHNnD0p6XcXo7biy9U7F7P-ctSkk-TrWGA") //test bot
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

	repo := db.NewRepo(client)

	bot.Debug = true
	Log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	//	ticker := time.NewTicker(1 * time.Minute)
	//	Log.Printf("Ticker started at: %v", time.Now())

	//go func() {
	//	for {
	//		<-ticker.C
	//		repo.DeleteSizesByTime(ctx)
	//		Log.Printf("Database was cleared at: %v", time.Now())
	//	}
	//
	//}()

	for update := range updates {

		var msg tgbotapi.MessageConfig

		if update.Message != nil { // If we got a message

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
			case AverageCommangd, AverageCommangd + CommandToBot, AverageCommangd + CommandToTestingBot:
				if update.Message.Chat.IsGroup() {
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
				if update.Message.Chat.IsGroup() {
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

				callbackData := "test"
				userData := repo.GetAllCredentials(ctx, update.Message.Chat.ID)
				for data := range userData {
					fmt.Println(data)
				}
				/*				var buttons = []tgbotapi.InlineKeyboardButton{
								{
									Text: "button 1",
									//URL:                          nil,
									//LoginURL:                     nil,
									CallbackData: &callbackData,
								},
								{
									Text: "button 2",
									//URL:                          nil,
									//LoginURL:                     nil,
									CallbackData: &callbackData,
								},
								{
									Text: "button 3",
									//URL:                          nil,
									//LoginURL:                     nil,
									CallbackData: &callbackData,
								},
								{
									Text: "button 4",
									//URL:                          nil,
									//LoginURL:                     nil,
									CallbackData: &callbackData,
								},
							}*/
				var usersKeyboardButtons = tgbotapi.NewInlineKeyboardMarkup()
				//				var testInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(buttons)
				//tgbotapi.NewInlineKeyboardRow(
				//	tgbotapi.InlineKeyboardButton{
				//		Text: "Выбери соперника",
				//		//URL:                          nil,
				//		//LoginURL:                     nil,
				//		CallbackData: &callbackData,
				//	},
				//),
				//)
				for _, userCred := range userData {
					buttonText := userCred.Fname + " @" + userCred.Username + " " + userCred.Lname

					row := tgbotapi.NewInlineKeyboardRow(
						tgbotapi.InlineKeyboardButton{
							Text:         buttonText,
							CallbackData: &callbackData,
						},
					)
					usersKeyboardButtons = utils.AddRowToInlineKeyboard(&usersKeyboardButtons, row)
				}

				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "test inline keyboard")
				msg.ReplyMarkup = usersKeyboardButtons
				msg.ParseMode = "HTML"
				msg.ReplyToMessageID = update.Message.MessageID

			}
			_, err := bot.Send(msg)
			if err != nil {
				Log.Printf(err.Error())
			}

		} /*else if update.InlineQuery != nil {
			time.Sleep(1 * time.Second)
			if update.InlineQuery.Query == "" {
				log.Println(update.InlineQuery)
				if utils.CheckLastMeasureDateIsToday(ctx, repo, update.Message.From.ID, update.Message.Chat.ID) {
					dickModel, _ := repo.GetLastMeasureByUserInThisChat(ctx, update.Message.From.ID, update.Message.Chat.ID)

					replyText := GetRandMeasureReplyPattern(int(dickModel.Dick_size))

					article := tgbotapi.NewInlineQueryResultArticle(update.InlineQuery.ID, "Узнай свой размер", replyText)

					inlineConf := tgbotapi.InlineConfig{
						InlineQueryID: update.InlineQuery.ID,
						IsPersonal:    false,
						CacheTime:     0,
						Results:       []interface{}{article},
					}

					if _, err := bot.Request(inlineConf); err != nil {
						Log.Println(err)
					}
				} else {
					dickSize := utils.GenerateDickSize()

					_, err := repo.InsertSize(ctx, update.InlineQuery.From.ID, update.InlineQuery.From.FirstName, update.InlineQuery.From.LastName, update.InlineQuery.From.UserName, dickSize, 0, false)
					if err != nil {
						Log.Printf(err.Error())
					}

					replyText := GetRandMeasureReplyPattern(dickSize)

					//article := tgbotapi.NewInlineQueryResultArticle(update.InlineQuery.ID, "Достать линейку", replyText)
					//
					//
					//inlineConf := tgbotapi.InlineConfig{
					//	InlineQueryID: update.InlineQuery.ID,
					//	IsPersonal:    false,
					//	CacheTime:     0,
					//	Results:       []interface{}{article},
					//}

					params := make(tgbotapi.Params)
					//params["inline_query_id"] = update.InlineQuery.ID
					//params["is_personal"] = "True"
					//params["cache_time"] = "0"
					params.AddNonEmpty("inline_query_id", update.InlineQuery.ID)
					params.AddBool("is_personal", true)
					params.AddNonEmpty("cache_time", "0")

					//var results = []interface{}
					//
					//var first

					resulVar := fmt.Sprintf("[{\"type\": \"article\", \"id\":\"%d\", \"title\": \"Достать линейку\", \"input_message_content\": \"%s\"}]", update.InlineQuery.ID, replyText)

					err = params.AddInterface("results", resulVar)
					if err != nil {
						Log.Println(err)
					}

					_, err = bot.MakeRequest("answerInlineQuery", params)
					//					_, err = bot.Request(inlineConf)

					if err != nil {
						Log.Println(err)
					}
				}

			}

		} else if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID /*update.CallbackQuery.Data*/ /*, "")
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
		}*/

	}

}
