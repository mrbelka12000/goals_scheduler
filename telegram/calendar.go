package telegram

import (
	"fmt"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	gs "github.com/mrbelka12000/goals_scheduler"
	"github.com/mrbelka12000/goals_scheduler/pkg/ptr"
)

const (
	ButtonPrev = "<"
	ButtonNext = ">"
)

type (
	calendarInfo struct {
		year  int
		month time.Month
	}
)

func (b *Bot) generateCalendar(chatID int64) {
	now := time.Now()
	inlineCalendar := b.generateCalendarMarkup(now.Year(), now.Month())

	_, err := b.bot.Send(tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:      chatID,
			ReplyMarkup: inlineCalendar,
		},
		Text: gs.MessageDeadline,
	})
	if err != nil {
		b.log.Error().Err(err).Msg("failed to send calendar")
		return
	}

	b.calendarMx.Lock()
	b.calendarStore[fmt.Sprintf("%v", chatID)] = calendarInfo{
		year:  now.Year(),
		month: now.Month(),
	}
	b.calendarMx.Unlock()
}

func (b *Bot) generateCalendarMarkup(year int, month time.Month) tgbotapi.InlineKeyboardMarkup {
	val := tgbotapi.InlineKeyboardMarkup{}

	val.InlineKeyboard = append(val.InlineKeyboard, addMonthYearRow(year, month))
	val.InlineKeyboard = append(val.InlineKeyboard, addDaysNamesRow())
	val.InlineKeyboard = append(val.InlineKeyboard, generateMonth(year, int(month))...)
	val.InlineKeyboard = append(val.InlineKeyboard, addSpecialButtons())

	return val
}

func addMonthYearRow(year int, month time.Month) []tgbotapi.InlineKeyboardButton {
	btn := tgbotapi.InlineKeyboardButton{Text: fmt.Sprintf("%s %v", month, year), CallbackData: ptr.Ref("-")}
	return []tgbotapi.InlineKeyboardButton{btn}
}

func addDaysNamesRow() []tgbotapi.InlineKeyboardButton {
	days := [7]string{"Mo", "Tu", "We", "Th", "Fr", "Sa", "Su"}
	var rowDays []tgbotapi.InlineKeyboardButton
	for _, day := range days {
		btn := tgbotapi.InlineKeyboardButton{
			Text:         day,
			CallbackData: ptr.Ref("-"),
		}
		rowDays = append(rowDays, btn)
	}
	return rowDays
}

func generateMonth(year int, month int) [][]tgbotapi.InlineKeyboardButton {

	var (
		firstDay          = date(year, month, 0)
		amountDaysInMonth = date(year, month+1, 0).Day()
		rowDays           []tgbotapi.InlineKeyboardButton
		rows              [][]tgbotapi.InlineKeyboardButton
		weekday           = int(firstDay.Weekday())
		row               []tgbotapi.InlineKeyboardButton
	)

	for i := 0; i < weekday; i++ {
		row = append(row, tgbotapi.InlineKeyboardButton{Text: " ", CallbackData: ptr.Ref("-")})
	}
	rowDays = append(rowDays, row...)

	amountWeek := weekday
	for i := 1; i <= amountDaysInMonth; i++ {
		if amountWeek == 7 {
			rows = append(rows, rowDays)
			amountWeek = 0
			rowDays = []tgbotapi.InlineKeyboardButton{}
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
		rowDays = append(rowDays, tgbotapi.InlineKeyboardButton{
			Text: btnText,
			CallbackData: ptr.Ref(callbackDataBuilder(CallbackData{
				Type: gs.CallbackTypeCalendar,
				Calendar: &CalendarData{
					Data: fmt.Sprintf("%v-%v-%v", year, monthStr, day),
				},
			})),
		})

		amountWeek++
	}

	for len(rowDays) != 7 {
		rowDays = append(rowDays, tgbotapi.InlineKeyboardButton{Text: " ", CallbackData: ptr.Ref("-")})
	}

	rows = append(rows, rowDays)

	return rows
}

func date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func addSpecialButtons() []tgbotapi.InlineKeyboardButton {
	var rowDays []tgbotapi.InlineKeyboardButton
	btnPrev := tgbotapi.InlineKeyboardButton{
		Text: ButtonPrev,
		CallbackData: ptr.Ref(callbackDataBuilder(
			CallbackData{
				Type: gs.CallbackTypeCalendar,
				Calendar: &CalendarData{
					Action: ButtonPrev,
				},
			})),
	}
	btnNext := tgbotapi.InlineKeyboardButton{
		Text: ButtonNext,
		CallbackData: ptr.Ref(callbackDataBuilder(CallbackData{
			Type: gs.CallbackTypeCalendar,
			Calendar: &CalendarData{
				Action: ButtonNext,
			},
		})),
	}
	rowDays = append(rowDays, btnPrev, btnNext)
	return rowDays
}

func (b *Bot) handleCallbackCalendar(cq *tgbotapi.CallbackQuery, data *CalendarData) {
	key := fmt.Sprintf("%v", cq.Message.Chat.ID)
	b.calendarMx.Lock()
	defer b.calendarMx.Unlock()

	info, ok := b.calendarStore[key]
	if !ok {
		b.log.Error().Msg("calendar not found")
		return
	}
	fmt.Printf("%+v\n%+v\n", data, info)
	if data.Action == ButtonPrev {
		now := time.Now()
		if info.year < now.Year() {
			callback := tgbotapi.NewCallbackWithAlert(cq.ID, gs.MessageDateBelow)
			b.handleSendMessageError(b.bot.Send(callback))
			return
		}
		if info.month <= now.Month() && info.year <= now.Year() {
			callback := tgbotapi.NewCallbackWithAlert(cq.ID, gs.MessageDateBelow)
			b.handleSendMessageError(b.bot.Send(callback))
			return
		}

		if info.month != 1 {
			info.month--
		} else {
			info.month = 12
			info.year--
		}
	} else if data.Action == ButtonNext {
		if info.month != 12 {
			info.month++
		} else {
			info.month = 1
			info.year++
		}
	} else {
		err := b.messageSvc.SetValue(gs.KeyDate, data.Data, cq.Message.From.ID)
		if err != nil {
			b.log.Err(err).Msg("failed to set value")
			callback := tgbotapi.NewCallbackWithAlert(cq.ID, gs.SomethingWentWrong)
			b.handleSendMessageError(b.bot.Send(callback))
			return
		}

		deleteMsg := tgbotapi.NewDeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
		b.handleSendMessageError(b.bot.Send(deleteMsg))

		callback := tgbotapi.NewCallbackWithAlert(cq.ID, gs.MessageSet)
		b.handleSendMessageError(b.bot.Send(callback))
		return
	}

	b.calendarStore[key] = info
	markup := b.generateCalendarMarkup(info.year, info.month)
	edit := tgbotapi.NewEditMessageReplyMarkup(cq.Message.Chat.ID, cq.Message.MessageID, markup)
	b.handleSendMessageError(b.bot.Send(edit))
	return
}
