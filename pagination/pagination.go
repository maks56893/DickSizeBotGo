package pagination

import (
	"math"
	"strconv"
	"strings"

	. "DickSizeBot/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PageLabel string

const (
	FirstPageLabel    PageLabel = `« {}`
	PreviousPageLabel PageLabel = `‹ {}`
	NextPageLabel     PageLabel = `{} ›`
	LastPageLabel     PageLabel = `{} »`
	CurrentPageLabel  PageLabel = `·{}·`
	maxUsersPerPage   int       = 6
)

var cancel = "cancel"
var cancelButton = []tgbotapi.InlineKeyboardButton{
	tgbotapi.InlineKeyboardButton{
		Text:         "В пее",
		CallbackData: &cancel,
	},
}

func (l PageLabel) Page(page int) string {
	return strings.Replace(string(l), "{}", strconv.Itoa(page), 1)
}

type InlineKeyboardPaginator struct {
	page             int
	totalPages       int
	data             string
	allUsersKeyboard [][]tgbotapi.InlineKeyboardButton
}

func NewInlineKeyboardPaginator(page int, data string, keyboard [][]tgbotapi.InlineKeyboardButton) *tgbotapi.InlineKeyboardMarkup {
	totalPages := int(math.Ceil(float64(len(keyboard)) / float64(maxUsersPerPage)))

	if page < 1 {
		page = 1
	}
	if totalPages < 1 {
		totalPages = 1
	}
	if len(data) == 0 {
		data = "page#1"
	}

	buttons := (&InlineKeyboardPaginator{
		page:             page,
		totalPages:       totalPages,
		data:             data,
		allUsersKeyboard: keyboard,
	}).buttons()

	if len(buttons) == 0 {
		return nil
	}

	return &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: buttons}
}

func (p *InlineKeyboardPaginator) buttons() [][]tgbotapi.InlineKeyboardButton {
	var resultKeyboard [][]tgbotapi.InlineKeyboardButton

	if p.totalPages == 1 {
		return nil
	} else if p.totalPages <= 5 {
		resultKeyboard = p.lessKeyboard()
		return append(resultKeyboard, cancelButton)
	} else if p.page <= 3 {
		resultKeyboard = p.startKeyboard()
		return append(resultKeyboard, cancelButton)
	} else if p.page > p.totalPages-3 {
		resultKeyboard = p.finishKeyboard()
		return append(resultKeyboard, cancelButton)
	} else {
		resultKeyboard = p.middleKeyboard()
		return append(resultKeyboard, cancelButton)
	}
}

func (p *InlineKeyboardPaginator) listKeyboardForCurrentPage() [][]tgbotapi.InlineKeyboardButton {
	var currentKeyboard [][]tgbotapi.InlineKeyboardButton

	if p.page == 1 {
		for _, row := range p.allUsersKeyboard[:(p.page * maxUsersPerPage)] {
			currentKeyboard = append(currentKeyboard, row)
		}
	} else if p.page > 1 && p.page != p.totalPages {
		for _, row := range p.allUsersKeyboard[(p.page-1)*maxUsersPerPage : p.page*maxUsersPerPage] {
			currentKeyboard = append(currentKeyboard, row)
		}
	} else if p.page == p.totalPages {
		for _, row := range p.allUsersKeyboard[(p.page-1)*maxUsersPerPage : len(p.allUsersKeyboard)] {
			currentKeyboard = append(currentKeyboard, row)
		}
	} else {
		Log.Errorf("Can't create inline keyboard for page: %d", p.page)
	}

	return currentKeyboard
}

func (p *InlineKeyboardPaginator) lessKeyboard() [][]tgbotapi.InlineKeyboardButton {
	keyboardDict := make([]tgbotapi.InlineKeyboardButton, 0, p.totalPages)
	for page := 1; page <= p.totalPages; page++ {
		keyboardDict = append(keyboardDict, p.isCurrentKeyboard(page))
	}

	currentKeyboard := make([][]tgbotapi.InlineKeyboardButton, 0, 12)
	currentKeyboard = p.listKeyboardForCurrentPage()

	currentKeyboard = append(currentKeyboard, keyboardDict)
	return currentKeyboard
}

func (p *InlineKeyboardPaginator) startKeyboard() [][]tgbotapi.InlineKeyboardButton {
	keyboardDict := make([]tgbotapi.InlineKeyboardButton, 0, 5)
	for page := 1; page <= 3; page++ {
		keyboardDict = append(keyboardDict, p.isCurrentKeyboard(page))
	}
	keyboardDict = append(keyboardDict, p.btnText(NextPageLabel.Page(4), 4))
	keyboardDict = append(keyboardDict, p.btnText(LastPageLabel.Page(p.totalPages), p.totalPages))

	currentKeyboard := p.listKeyboardForCurrentPage()

	currentKeyboard = append(currentKeyboard, keyboardDict)
	return currentKeyboard
}

func (p *InlineKeyboardPaginator) middleKeyboard() [][]tgbotapi.InlineKeyboardButton {
	keyboardDict := make([]tgbotapi.InlineKeyboardButton, 0, 5)

	keyboardDict = []tgbotapi.InlineKeyboardButton{
		p.btnText(FirstPageLabel.Page(1), 1),
		p.btnText(PreviousPageLabel.Page(p.page-1), p.page-1),
		p.btnText(CurrentPageLabel.Page(p.page), p.page),
		p.btnText(NextPageLabel.Page(p.page+1), p.page+1),
		p.btnText(LastPageLabel.Page(p.totalPages), p.totalPages),
	}

	currentKeyboard := p.listKeyboardForCurrentPage()

	currentKeyboard = append(currentKeyboard, keyboardDict)
	return currentKeyboard
}

func (p *InlineKeyboardPaginator) finishKeyboard() [][]tgbotapi.InlineKeyboardButton {
	keyboardDict := make([]tgbotapi.InlineKeyboardButton, 0, 5)

	keyboardDict = append(keyboardDict,
		p.btnText(FirstPageLabel.Page(1), 1),
		p.btnText(PreviousPageLabel.Page(p.totalPages-3), p.totalPages-3))

	for i := 3; i <= 5; i++ {
		keyboardDict = append(keyboardDict, p.isCurrentKeyboard(p.totalPages-5+i))
	}

	currentKeyboard := p.listKeyboardForCurrentPage()

	currentKeyboard = append(currentKeyboard, keyboardDict)
	return currentKeyboard
}

func (p *InlineKeyboardPaginator) isCurrentKeyboard(page int) tgbotapi.InlineKeyboardButton {
	if page == p.page {
		return p.btnText(CurrentPageLabel.Page(page), page)
	}
	return p.btn(page)
}

func (p *InlineKeyboardPaginator) btn(page int) tgbotapi.InlineKeyboardButton {
	return p.btnText(strconv.Itoa(page), page)
}

func (p *InlineKeyboardPaginator) btnText(text string, page int) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(text, strings.ReplaceAll(p.data, p.data, "page#"+strconv.Itoa(page)))
}
