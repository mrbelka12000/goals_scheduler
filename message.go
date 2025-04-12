package goals_scheduler

const (
	DateFormat = "2006-01-02"

	SomethingWentWrong = "Что то пошло не так, повторите пожалуйста еще раз"

	MessageTimerFormat = `
		Введите интервал для напоминания
Допустимые единицы времени: "s", "m", "h". Ex: 5h = каждые 5 часов`

	MessageTimeFormat = `
		Укажите время для напоминание
Формат час:минута. Пример: 14:58 (UTC+5)`

	MessageDayFormat = `
		Укажите дeнь:
`

	MessageInputText = "Введите текст цели"

	MessageDeadline = "Введите крайний срок для цели"

	MessageChooseDay = "Выберите дни"

	MessageDone = "Цель сохранилась"

	MessageChooseMethod = "Выберите метод:"

	MessageNeedText = "Пожалуйста введите сообщение"

	MessageDateBelow = "Нельзя ставить цели на прошлое"

	MessageSet = "Сохранено"
)
