package commands

import (
	"context"
	"strconv"
	"time"

	cash2 "DickSizeBot/cash"
	. "DickSizeBot/logger"
	"DickSizeBot/pagination"
	models "DickSizeBot/postgres/models/dick_size"
	"DickSizeBot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DuelCommandObj struct {
	ctx  context.Context
	repo models.Repository
	cash cash2.Cache
	Bot  *tgbotapi.BotAPI
}

func NewDuelCommandObj(ctx context.Context, repo models.Repository, bot *tgbotapi.BotAPI) DuelCommandObj {
	return DuelCommandObj{
		ctx:  ctx,
		repo: repo,
		cash: cash2.NewCache(),
		Bot:  bot,
	}
}

func (cmd *DuelCommandObj) Execute(update tgbotapi.Update) (msg tgbotapi.MessageConfig) {
	if !(update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup()) {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Работает только в группе")
		msg.ParseMode = "HTML"
		msg.ReplyToMessageID = update.Message.MessageID
		// _, err := cmd.Bot.Send(msg)
		// if err != nil {
		// 	Log.Printf(err.Error())
		// }
		return
	}

	if !utils.CheckLastUsersDuelIsToday(cmd.ctx, cmd.repo, update.Message.From.ID, update.Message.Chat.ID) {
		userData := cmd.repo.GetAllCredentials(cmd.ctx, update.Message.Chat.ID)

		var usersKeyboardButtons = tgbotapi.NewInlineKeyboardMarkup()
		for _, userCred := range userData {
			if userCred.UserId == update.Message.From.ID {
				continue
			}

			userId := "duel_user_id#" + strconv.Itoa(int(userCred.UserId))

			buttonText := userCred.Fname + " @" + userCred.Username + " " + userCred.Lname

			row := tgbotapi.NewInlineKeyboardRow(
				tgbotapi.InlineKeyboardButton{
					Text:         buttonText,
					CallbackData: &userId,
				},
			)
			usersKeyboardButtons = utils.AddRowToInlineKeyboard(&usersKeyboardButtons, row)
		}

		if _, ok := cmd.cash.Get("duelKeyboard"); ok {
			Log.Infof("cash for duel keyboard already exists, deleting it...")
			_ = cmd.cash.Delete("duelKeyboard")
		}
		cmd.cash.Set("duelKeyboard", usersKeyboardButtons.InlineKeyboard, 10*time.Minute)

		keyboardWithUsers := pagination.NewInlineKeyboardPaginator(1, "page#1", usersKeyboardButtons.InlineKeyboard)

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "С кем хочешь помериться?")
		msg.ReplyToMessageID = update.Message.MessageID
		msg.ParseMode = "HTML"
		msg.ReplyMarkup = keyboardWithUsers
	}

	return
}

func (cmd *DuelCommandObj) ExecuteFromCallback(update tgbotapi.Update) (msg tgbotapi.MessageConfig) {
	if !(update.CallbackQuery.Message.Chat.IsGroup() || update.CallbackQuery.Message.Chat.IsSuperGroup()) {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Работает только в группе")
		msg.ParseMode = "HTML"
		msg.ReplyToMessageID = update.Message.MessageID
		// _, err := cmd.Bot.Send(msg)
		// if err != nil {
		// 	Log.Printf(err.Error())
		// }
		return
	}

	if !utils.CheckLastUsersDuelIsToday(cmd.ctx, cmd.repo, update.CallbackQuery.Message.From.ID, update.CallbackQuery.Message.Chat.ID) {
		userData := cmd.repo.GetAllCredentials(cmd.ctx, update.CallbackQuery.Message.Chat.ID)

		var usersKeyboardButtons = tgbotapi.NewInlineKeyboardMarkup()
		for _, userCred := range userData {
			if userCred.UserId == update.CallbackQuery.Message.From.ID {
				continue
			}

			userId := "duel_user_id#" + strconv.Itoa(int(userCred.UserId)) + "#caller_user_id#" + strconv.Itoa(int(update.CallbackQuery.Message.From.ID))

			buttonText := userCred.Fname + " @" + userCred.Username + " " + userCred.Lname

			row := tgbotapi.NewInlineKeyboardRow(
				tgbotapi.InlineKeyboardButton{
					Text:         buttonText,
					CallbackData: &userId,
				},
			)
			usersKeyboardButtons = utils.AddRowToInlineKeyboard(&usersKeyboardButtons, row)
		}

		// if _, ok := cmd.cash.Get("duelKeyboard"); ok {
		// 	Log.Infof("cash for duel keyboard already exists, deleting it...")
		// 	_ = cmd.cash.Del("duelKeyboard")
		// }
		// cmd.cash.Set("duelKeyboard", usersKeyboardButtons.InlineKeyboard, 10*time.Minute)

		test := pagination.NewInlineKeyboardPaginator(1, "page#1", usersKeyboardButtons.InlineKeyboard)

		msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "С кем хочешь помериться?")
		msg.ReplyMarkup = test
		msg.ParseMode = "HTML"
		//		msg.ReplyToMessageID = update.CallbackQuery.Message.MessageID
	}

	return
}
