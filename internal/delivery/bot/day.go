package bot

//
//import (
//	"encoding/json"
//	"fmt"
//	"strings"
//	"sync"
//
//	"github.com/rs/zerolog"
//	"github.com/yanzay/tbot/v2"
//
//	gs "github.com/mrbelka12000/goals_scheduler"
//	"github.com/mrbelka12000/goals_scheduler/internal/models"
//)
//
//const (
//	chooseEmoji = "🎯"
//)
//
//type (
//	day struct {
//		client *tbot.Client
//
//		// Key is fmt.Sprintf(chatID+:+messageID)
//		store map[string]DayInfo
//		sync.Mutex
//		log zerolog.Logger
//	}
//
//	DayInfo struct {
//		Mark []bool
//	}
//
//	Day struct {
//		Action  string `json:"action,omitempty"`
//		Weekday gs.Day `json:"weekday"`
//	}
//)
//
//func newDay(client *tbot.Client, log zerolog.Logger) *day {
//	return &day{
//		client: client,
//
//		store: make(map[string]DayInfo),
//		log:   log,
//	}
//}
//
//func (d *day) getBaseKeyboard(info *DayInfo) *tbot.InlineKeyboardMarkup {
//	markup := &tbot.InlineKeyboardMarkup{
//		InlineKeyboard: [][]tbot.InlineKeyboardButton{
//			{
//				{
//					Text: "Пн",
//					CallbackData: callbackDataBuilder(CallbackData{
//						Type: gs.CallbackTypeDay,
//						Day: &Day{
//							Action:  ActionDayMark,
//							Weekday: gs.Monday,
//						},
//					}),
//				},
//			},
//			{
//				{
//					Text: "Вт",
//					CallbackData: callbackDataBuilder(CallbackData{
//						Type: gs.CallbackTypeDay,
//						Day: &Day{
//							Action:  ActionDayMark,
//							Weekday: gs.Tuesday,
//						},
//					}),
//				},
//			},
//			{
//				{
//					Text: "Ср",
//					CallbackData: callbackDataBuilder(CallbackData{
//						Type: gs.CallbackTypeDay,
//						Day: &Day{
//							Action:  ActionDayMark,
//							Weekday: gs.Wednesday,
//						},
//					}),
//				},
//			},
//			{
//				{
//					Text: "Чт",
//					CallbackData: callbackDataBuilder(CallbackData{
//						Type: gs.CallbackTypeDay,
//						Day: &Day{
//							Action:  ActionDayMark,
//							Weekday: gs.Thursday,
//						},
//					}),
//				},
//			},
//			{
//				{
//					Text: "Пт",
//					CallbackData: callbackDataBuilder(CallbackData{
//						Type: gs.CallbackTypeDay,
//						Day: &Day{
//							Action:  ActionDayMark,
//							Weekday: gs.Friday,
//						},
//					}),
//				},
//			},
//			{
//				{
//					Text: "Сб",
//					CallbackData: callbackDataBuilder(CallbackData{
//						Type: gs.CallbackTypeDay,
//						Day: &Day{
//							Action:  ActionDayMark,
//							Weekday: gs.Saturday,
//						},
//					}),
//				},
//			},
//			{
//				{
//					Text: "Вс",
//					CallbackData: callbackDataBuilder(CallbackData{
//						Type: gs.CallbackTypeDay,
//						Day: &Day{
//							Action:  ActionDayMark,
//							Weekday: gs.Sunday,
//						},
//					}),
//				},
//			},
//			{
//				{
//					Text: "Выбрать",
//					CallbackData: callbackDataBuilder(CallbackData{
//						Type: gs.CallbackTypeDay,
//						Day: &Day{
//							Action:  ActionDaySubmit,
//							Weekday: gs.Wednesday,
//						},
//					}),
//				},
//			},
//		},
//	}
//
//	if info == nil {
//		return markup
//	}
//
//	for i, v := range info.Mark {
//		text := markup.InlineKeyboard[i][0].Text
//		enabled := strings.Contains(text, chooseEmoji)
//
//		if !v {
//			if enabled {
//				text = strings.Replace(text, chooseEmoji, "", 1)
//			}
//		} else {
//			if !enabled {
//				text = fmt.Sprintf("%s %s", text, chooseEmoji)
//			}
//		}
//		markup.InlineKeyboard[i][0].Text = text
//	}
//
//	//for i := 0; i < len(markup.InlineKeyboard); i++ {
//	//	for j, v := range markup.InlineKeyboard[i] {
//	//		if strings.Contains(v.CallbackData, fmt.Sprint(day)) {
//	//
//	//
//	//
//	//			if ! {
//	//			} else {
//	//				text = strings.Replace(text, chooseEmoji, "", 1)
//	//			}
//	//
//	//			markup.InlineKeyboard[i][j].Text = text
//	//		}
//	//	}
//	//}
//
//	return markup
//}
//
//func (d *day) handleCallbackDay(cq *tbot.CallbackQuery, data *Day) string {
//	key := fmt.Sprintf("%s:%d", cq.Message.Chat.ID, cq.Message.MessageID)
//	d.Lock()
//	defer d.Unlock()
//
//	switch data.Action {
//	case ActionDayMark:
//		info, ok := d.store[key]
//		if !ok {
//			info = DayInfo{
//				Mark: make([]bool, 7),
//			}
//		}
//
//		if !info.Mark[data.Weekday] {
//			info.Mark[data.Weekday] = true
//		} else {
//			info.Mark[data.Weekday] = false
//		}
//
//		markup := d.getBaseKeyboard(&info)
//		d.client.EditMessageReplyMarkup(cq.Message.Chat.ID, cq.Message.MessageID, tbot.OptInlineKeyboardMarkup(markup))
//
//		d.store[key] = info
//	case ActionDaySubmit:
//		info, ok := d.store[key]
//		if !ok || !anyTrue(info.Mark) {
//			d.client.AnswerCallbackQuery(cq.ID, tbot.OptText(gs.MessageChooseDay))
//			return ""
//		}
//
//		jsonData, err := json.Marshal(info)
//		if err != nil {
//			d.log.Err(err).Msg("failed to marshal info")
//			d.client.AnswerCallbackQuery(cq.ID, tbot.OptText(gs.SomethingWentWrong))
//			return ""
//		}
//
//		delete(d.store, key)
//		d.client.DeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
//		return string(jsonData)
//	}
//	return ""
//}
//
//func anyTrue(arr []bool) bool {
//	for _, v := range arr {
//		if v {
//			return true
//		}
//	}
//	return false
//}
