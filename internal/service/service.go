package service

import (
	"bytes"
	"errors"
	"html/template"
	"mstuca_schedule/internal/botErrors"
	"mstuca_schedule/internal/models"
	"mstuca_schedule/internal/service/processor"
	schedulegetter "mstuca_schedule/internal/service/schedule_getter"
	"mstuca_schedule/pkg/cache"
	"mstuca_schedule/pkg/logger"
	"regexp"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	profileKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(editProfileText, editProfileText),
			tgbotapi.NewInlineKeyboardButtonData(goToscheduleText, goToscheduleText),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(githubProject, "http://1.com"),
		),
	)

	editPersonTitleKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(LOH, LOH),
			tgbotapi.NewInlineKeyboardButtonData(Emelya, Emelya),
		),
	)

	pastRegisterKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(editProfileTextV2, editProfileText),
			tgbotapi.NewInlineKeyboardButtonData(goToscheduleText, goToscheduleText),
		),
	)

	groupExpression       = regexp.MustCompile(`[А-Яа-я]+\d+`)
	nameExpression        = regexp.MustCompile(`[А-Яа-я]+`)
	groupAnswerExpression = regexp.MustCompile(`[А-Яа-я]+\d+`)
	nameAnswerExpression  = regexp.MustCompile(`^[А-ЯЁ][а-яё]+\s[А-ЯЁ][а-яё]+\s[А-ЯЁ][а-яё]+$`)
	registerFromTmpl      = template.Must(template.New("register").Parse(registerFromTmplText))
	subgroupExpression    = regexp.MustCompile(`^\d$`)
)

type UpdateProcessor interface {
	Process(update *tgbotapi.Update) tgbotapi.Chattable
}

type updateProcessor struct {
	scheduleGetter schedulegetter.ScheduleGetter
	processor      processor.Processor
	logger         logger.Logger
	cache          cache.Cache
}

func New(logger logger.Logger) (UpdateProcessor, error) {

	scheduleGetter := schedulegetter.New()

	processor, err := processor.New()
	if err != nil {
		return nil, err
	}

	return &updateProcessor{
		scheduleGetter: scheduleGetter,
		processor:      processor,
		logger:         logger,
		cache:          cache.New(),
	}, nil
}

func (up *updateProcessor) Process(update *tgbotapi.Update) tgbotapi.Chattable {

	if update.Message != nil {
		return up.processMessage(update.Message)
	}

	if update.CallbackQuery != nil {
		return up.proccessCallbackQuery(update.CallbackQuery)
	}

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		UnkonwnMessage,
	)

	return msg
}

func (up *updateProcessor) processMessage(message *tgbotapi.Message) tgbotapi.Chattable {
	switch {
	case message.Text == startQuery:
		newMessage := tgbotapi.NewMessage(message.Chat.ID, startMessage)
		newMessage.ReplyMarkup = profileKeyboard
		return newMessage

	case groupExpression.MatchString(message.Text):
		user, err := up.cache.GetUser(message.From.ID)
		if err != nil {

		}

		if user.State != editGroupNameState {
			return tgbotapi.NewMessage(
				message.Chat.ID,
				blockAnswerText,
			)
		}

		groups, err := up.scheduleGetter.GetGroupID(message.Text)
		if err != nil {
			if errors.Is(err, botErrors.ErrNoGroupsFound) {
				return tgbotapi.NewMessage(
					message.Chat.ID,
					botErrors.ErrNoGroupsFound.Error(),
				)
			} else {
				up.logger.Error("error while getting groups", botErrors.Err(err))
				return tgbotapi.NewMessage(
					message.Chat.ID,
					internalServiceError,
				)
			}
		}

		newMessage := tgbotapi.NewMessage(message.Chat.ID, chooseGroupText)
		markup := tgbotapi.NewInlineKeyboardMarkup(
			doKeybordRowsFromResponse(groups),
		)
		newMessage.ReplyMarkup = markup

		user.State = editSubGroupState
		up.cache.SaveUserInfo(user)

		return newMessage

	case nameExpression.MatchString(message.Text):
		user, err := up.cache.GetUser(message.From.ID)
		if err != nil {
			return tgbotapi.NewMessage(
				message.Chat.ID, // !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
				internalServiceError,
			)
		}

		if user.State != editPersonNameState {
			return tgbotapi.NewMessage(
				message.Chat.ID,
				blockAnswerText,
			)
		}

		if user.Title == LOH {
			user.Name = validateAnswer(message.Text)

			user.State = editGroupNameState
			up.cache.SaveUserInfo(user)

			userMessageText, err := registerMessageParse(user)
			if err != nil {
				return tgbotapi.NewMessage(
					message.Chat.ID,
					internalServiceError,
				)
			}

			return tgbotapi.NewMessage(
				message.Chat.ID,
				//message.MessageID,
				userMessageText,
			)
		}

		//TODO

	case subgroupExpression.MatchString(message.Text):
		user, err := up.cache.GetUser(message.From.ID)
		if err != nil {
			return tgbotapi.NewMessage(
				message.Chat.ID, // !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
				internalServiceError,
			)
		}

		if user.State != lastState {
			return tgbotapi.NewMessage(
				message.Chat.ID,
				blockAnswerText,
			)
		}

		validatedAnswerText := validateAnswer(message.Text)

		user.Subgroup, err = strconv.Atoi(validatedAnswerText)
		if err != nil {
			return tgbotapi.NewMessage(
				message.Chat.ID,
				BadNumberMessage,
			)
		}

		userMessageText, err := registerMessageParse(user)
		if err != nil {
			return tgbotapi.NewMessage(
				message.Chat.ID,
				internalServiceError,
			)
		}

		// err = up.processor.SaveProfile(user)
		// if err != nil {
		// 	return tgbotapi.NewMessage(
		// 		message.Chat.ID,
		// 		internalServiceError,
		// 	)
		// }

		up.cache.DeleteUser(user.ID)

		msg := tgbotapi.NewMessage(
			message.Chat.ID,
			userMessageText,
		)

		msg.ReplyMarkup = pastRegisterKeyboard

		return msg
	}

	return tgbotapi.NewMessage(message.Chat.ID, UnkonwnMessage)
}

func doKeybordRowsFromResponse(groups []*models.Group) []tgbotapi.InlineKeyboardButton {
	bottons := make([]tgbotapi.InlineKeyboardButton, 0)

	for _, group := range groups {
		bottons = append(bottons,
			tgbotapi.NewInlineKeyboardButtonData(
				group.Label,
				group.Label,
			))
	}

	return bottons
}

func (up *updateProcessor) proccessCallbackQuery(callback *tgbotapi.CallbackQuery) tgbotapi.Chattable {
	switch callback.Data {
	case editProfileText:
		user := &models.User{
			ID: callback.From.ID,
		}

		up.cache.SaveUserInfo(user)

		return tgbotapi.NewEditMessageTextAndMarkup(
			callback.Message.Chat.ID,
			callback.Message.MessageID,
			editPersonTitleText,
			editPersonTitleKeyboard,
		)
	case LOH:
		user, err := up.cache.GetUser(callback.From.ID)
		if err != nil {
			return tgbotapi.NewMessage(
				callback.Message.Chat.ID, // !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
				internalServiceError,
			)
		}

		if user.State != editPersonTitleState {
			return tgbotapi.NewMessage(
				callback.Message.Chat.ID,
				blockAnswerText,
			)
		}

		user.Title = LOH
		user.State = editPersonNameState
		up.cache.SaveUserInfo(user)

		messageText, err := registerMessageParse(user)
		if err != nil {
			return tgbotapi.NewMessage(
				callback.Message.Chat.ID,
				internalServiceError,
			)
		}

		return tgbotapi.NewEditMessageText(
			callback.Message.Chat.ID,
			callback.Message.MessageID,
			messageText,
		)
	default:
		switch {
		case groupAnswerExpression.MatchString(callback.Data):
			user, err := up.cache.GetUser(callback.From.ID)
			if err != nil {

			}

			if user.State != editSubGroupState {
				return tgbotapi.NewMessage(
					callback.Message.Chat.ID,
					blockAnswerText,
				)
			}

			user.GroupName = validateAnswer(callback.Data)

			userMessageText, err := registerMessageParse(user)
			if err != nil {
				return tgbotapi.NewMessage(
					callback.Message.Chat.ID,
					internalServiceError,
				)
			}

			user.State = lastState
			up.cache.SaveUserInfo(user)

			return tgbotapi.NewEditMessageText(
				callback.Message.Chat.ID,
				callback.Message.MessageID,
				userMessageText,
			)

		case nameAnswerExpression.MatchString(callback.Data):

		default:
			return nil
		}
		return nil
	}
}

func validateAnswer(answer string) string {
	return strings.TrimSpace(answer)
}

func registerMessageParse(user *models.User) (string, error) {
	buf := make([]byte, 0, len(registerFromTmplText))
	buffer := bytes.NewBuffer(buf)

	if err := registerFromTmpl.Execute(buffer, user); err != nil {
		return "", err
	}

	switch user.State {
	case 1:
		buffer.Grow(len(editPersonNameTextByte))
		_, _ = buffer.Write(editPersonNameTextByte)
		return buffer.String(), nil
	case 2:
		buffer.Grow(len(editGroupNameTextByte))
		_, _ = buffer.Write(editGroupNameTextByte)
		return buffer.String(), nil
	case 3:
		buffer.Grow(len(editSubGroupTextByte))
		_, _ = buffer.Write(editSubGroupTextByte)
		return buffer.String(), nil
	default:
		buffer.Grow(len(lastStateTextByte))
		_, _ = buffer.Write(lastStateTextByte)
		return buffer.String(), nil
	}
}
