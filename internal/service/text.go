package service

const (
	editProfileText      = "Настроить профиль"
	editProfileTextV2    = "Настроить профиль заново"
	blockAnswerText      = "Не то сообщение"
	internalServiceError = "Мне стало плохо в моменте, попробуйте еще раз"
	BadNumberMessage     = "Введите нормально число"
	goToscheduleText     = "Перейти к расписанию"
	githubProject        = "Github))"
	UnkonwnMessage       = `Данный функционал пока не поддерживается. Не понимаю вас(, попробуйте использовать копки под сообщениями.
	Если их нет. Попробуйте ввести /start`
	startQuery           = "/start"
	editPersonTitleText  = `Вы ученик или преподаватель?`
	editPersonTitleState = 0
	LOH                  = "Студент"
	Emelya               = "Преподаватель"
	editPersonNameText   = "Введите свое ФИО"
	editPersonNameState  = 1
	editGroupNameText    = "Введите название группы в виде БИС201 или МАГ241"
	editGroupNameState   = 2
	chooseGroupText      = `Выберите вашу группу из списка`
	editSubGroupText     = "Введите номер своей подгруппы"
	editSubGroupState    = 3
	registerFromTmplText = `Ваш профиль:
	Имя: {{.Name}}
	Титул: {{.Title}}{{if eq .Title "Студент"}}
	Группа: {{.GroupName}}
	Подгруппа №: {{.Subgroup}}
	{{end}}`
	lastStateText = "Ваш профиль готов. Для показа расписания нажмите на соответствующую кнопку"
	lastState     = 4
)

var (
	editPersonNameTextByte = []byte(`Введите свое ФИО`)
	editGroupNameTextByte  = []byte(`Введите название группы в виде БИС201 или МАГ241`)
	//chooseGroupTextByte    = []byte(`\nВыберите вашу группу из списка`)
	editSubGroupTextByte = []byte(`Введите номер своей подгруппы`)
	lastStateTextByte    = []byte(`Ваш профиль готов. Для показа расписания нажмите на соответствующую кнопку`)
)

const startMessage = `Здравствуйте, я бот-помощник получения расписания МГТУ ГА. Проверьте данные своего пользователя. 
Если вы используете меня впервые - нажмите на кнопку "Настроить профиль".
Если ваш профиль настроен - можете посмотреть свое расписание`
