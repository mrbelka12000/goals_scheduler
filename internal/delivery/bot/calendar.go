package bot

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/rs/zerolog"
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
		log zerolog.Logger
	}
	info struct {
		year  int
		month time.Month
	}
)

func newCalendar(client *tbot.Client, log zerolog.Logger) *calendar {
	return &calendar{
		client: client,
		store:  make(map[string]*info),
		log:    log,
	}
}

func (c *calendar) calendarHandler(m *tbot.Message) error {
	now := time.Now()
	inlineCalendar := c.GenerateCalendar(now.Year(), now.Month())

	msg, err := c.client.SendMessage(m.Chat.ID, "Дедлайн", tbot.OptInlineKeyboardMarkup(inlineCalendar))
	if err != nil {
		return fmt.Errorf("send calendar: %w", err)
	}

	c.Lock()
	c.store[fmt.Sprintf("%v:%v", m.Chat.ID, msg.MessageID)] = &info{
		year:  now.Year(),
		month: now.Month(),
	}
	c.Unlock()

	return nil
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
	return []tbot.InlineKeyboardButton{btn}
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

	var (
		firstDay          = date(year, month, 0)
		amountDaysInMonth = date(year, month+1, 0).Day()
		rowDays           []tbot.InlineKeyboardButton
		rows              [][]tbot.InlineKeyboardButton
		weekday           = int(firstDay.Weekday())
		row               []tbot.InlineKeyboardButton
	)

	for i := 0; i < weekday; i++ {
		row = append(row, tbot.InlineKeyboardButton{Text: " ", CallbackData: "-"})
	}
	rowDays = append(rowDays, row...)

	amountWeek := weekday
	for i := 1; i <= amountDaysInMonth; i++ {
		if amountWeek == 7 {
			rows = append(rows, rowDays)
			amountWeek = 0
			rowDays = []tbot.InlineKeyboardButton{}
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
