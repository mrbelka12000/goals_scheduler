package bot

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/yanzay/tbot/v2"
)

const (
	BTN_PREV         = "<"
	BTN_NEXT         = ">"
	CallbackCalendar = "calendar"
)

type (
	calendar struct {
		client *tbot.Client
		// Key is fmt.Sprintf(chatID+:+messageID)
		store map[string]*info
		sync.Mutex
	}
	info struct {
		year  int
		month time.Month
	}
)

func newCalendar(client *tbot.Client) *calendar {
	return &calendar{
		client: client,
		store:  make(map[string]*info),
	}
}

func (c *calendar) calendarHandler(m *tbot.Message) {
	now := time.Now()
	calendar := c.GenerateCalendar(now.Year(), now.Month())

	msg, _ := c.client.SendMessage(m.Chat.ID, "Дедлайн", tbot.OptInlineKeyboardMarkup(calendar))

	c.Lock()
	c.store[fmt.Sprintf("%v:%v", m.Chat.ID, msg.MessageID)] = &info{
		year:  now.Year(),
		month: now.Month(),
	}
	c.Unlock()
}

func (c *calendar) GenerateCalendar(year int, month time.Month) *tbot.InlineKeyboardMarkup {
	val := &tbot.InlineKeyboardMarkup{}

	val.InlineKeyboard = append(val.InlineKeyboard, addMonthYearRow(year, month))
	val.InlineKeyboard = append(val.InlineKeyboard, addDaysNamesRow())
	val.InlineKeyboard = append(val.InlineKeyboard, generateMonth(year, int(month))...)
	val.InlineKeyboard = append(val.InlineKeyboard, addSpecialButtons())

	return val
}

func (c *calendar) handleCallback(cq *tbot.CallbackQuery) string {
	key := fmt.Sprintf("%s:%d", cq.Message.Chat.ID, cq.Message.MessageID)
	c.Lock()
	defer c.Unlock()

	data := cq.Data[len(CallbackCalendar)+1:]
	info, ok := c.store[key]
	if !ok {
		return ""
	}

	if data == "<" {
		now := time.Now()
		if info.year < now.Year() {
			c.client.AnswerCallbackQuery(cq.ID, tbot.OptText("Нельзя ставить цели на прошлое"))
			return ""
		}
		if info.month <= 2 && info.year < now.Year() {
			c.client.AnswerCallbackQuery(cq.ID, tbot.OptText("Нельзя ставить цели на прошлое"))
			return ""
		}

		if info.month != 1 {
			info.month--
		} else {
			info.month = 12
			info.year--
		}
	} else if data == ">" {
		if info.month != 12 {
			info.month++
		} else {
			info.month = 1
			info.year++
		}
	} else {
		c.client.DeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
		c.client.AnswerCallbackQuery(cq.ID, tbot.OptText("OK"))
		return data
	}

	markup := c.GenerateCalendar(info.year, info.month)
	c.client.EditMessageReplyMarkup(cq.Message.Chat.ID, cq.Message.MessageID, tbot.OptInlineKeyboardMarkup(markup))
	return ""
}

func addMonthYearRow(year int, month time.Month) []tbot.InlineKeyboardButton {
	btn := tbot.InlineKeyboardButton{Text: fmt.Sprintf("%s %v", month, year), CallbackData: "-"}
	row := []tbot.InlineKeyboardButton{btn}
	return row
}

func addDaysNamesRow() []tbot.InlineKeyboardButton {
	days := [7]string{"Mo", "Tu", "We", "Th", "Fr", "Sa", "Su"}
	var rowDays []tbot.InlineKeyboardButton
	for _, day := range days {
		btn := tbot.InlineKeyboardButton{
			Text: day, CallbackData: "-",
		}
		rowDays = append(rowDays, btn)
	}
	return rowDays
}

func generateMonth(year int, month int) [][]tbot.InlineKeyboardButton {
	firstDay := date(year, month, 0)
	now := time.Now()
	amountDaysInMonth := date(year, month+1, 0).Day()
	var rows [][]tbot.InlineKeyboardButton

	weekday := int(firstDay.Weekday())
	rowDays := []tbot.InlineKeyboardButton{}
	for i := 1; i <= weekday; i++ {
		btn := tbot.InlineKeyboardButton{Text: " ", CallbackData: "-"}
		rowDays = append(rowDays, btn)
	}

	amountWeek := weekday
	for i := 1; i <= amountDaysInMonth; i++ {
		if amountWeek == 7 {
			rows = append(rows, rowDays)
			amountWeek = 0
			rowDays = []tbot.InlineKeyboardButton{}
		}
		if i < now.Day() && year == now.Year() && month == int(now.Month()) {
			rowDays = append(rowDays, tbot.InlineKeyboardButton{Text: string(i), CallbackData: "-"})
			amountWeek++
			continue
		}
		day := strconv.Itoa(i)
		if len(day) == 1 {
			day = fmt.Sprintf("0%v", day)
		}
		monthStr := strconv.Itoa(month)
		if len(monthStr) == 1 {
			monthStr = fmt.Sprintf("0%v", monthStr)
		}

		btnText := fmt.Sprintf("%v", i)
		if time.Now().Day() == i {
			btnText = fmt.Sprintf("%v", i)
		}
		//btn := tgbotapi.NewInlineKeyboardButtonData(btnText, fmt.Sprintf("%v.%v.%v", year, monthStr, day))
		rowDays = append(rowDays, tbot.InlineKeyboardButton{Text: btnText, CallbackData: fmt.Sprintf("%v %v-%v-%v", CallbackCalendar, year, monthStr, day)})
		amountWeek++
	}

	for len(rowDays) != 7 {
		rowDays = append(rowDays, tbot.InlineKeyboardButton{Text: " ", CallbackData: "-"})
	}
	rows = append(rows, rowDays)

	return rows
}

func date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func addSpecialButtons() []tbot.InlineKeyboardButton {
	var rowDays []tbot.InlineKeyboardButton
	btnPrev := tbot.InlineKeyboardButton{Text: BTN_PREV, CallbackData: fmt.Sprintf("%v %v", CallbackCalendar, BTN_PREV)}
	btnNext := tbot.InlineKeyboardButton{Text: BTN_NEXT, CallbackData: fmt.Sprintf("%v %v", CallbackCalendar, BTN_NEXT)}
	rowDays = append(rowDays, btnPrev, btnNext)
	return rowDays
}
