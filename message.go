package goals_scheduler

const (
	DateFormat = "2006-01-02"

	SomethingWentWrong = "Что то пошло не так, повторите пожалуйста еще раз"
	MessageTimerFormat = `
		Введите интервал для напоминания
Допустимые единицы времени: "s", "m", "h". Ex: 5h = каждые 5 часов`
	MessageTimeFormat = `
		Укажите время для напоминание
Формат час:минута. Пример: 14:58 (UTC)`
	MessageDayFormat = `
		Укажите дeнь:
`

	MessageStart     = "Введите текст цели"
	MessageDeadline  = "Введите крайний срок для цели"
	MessageChooseDay = "Выберите день для напоминания"
	MessageDone      = "Цель сохранилась"
)
