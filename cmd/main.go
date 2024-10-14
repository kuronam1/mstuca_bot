package main

import (
	"context"
	"io"
	"log"
	"log/slog"
	"mstuca_schedule/internal/botErrors"
	"mstuca_schedule/internal/service"
	"mstuca_schedule/pkg/logger"
	"os"
	"os/signal"
	"strings"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kelseyhightower/envconfig"
)

// var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
// 	tgbotapi.NewInlineKeyboardRow(
// 		tgbotapi.NewInlineKeyboardButtonURL("1.com", "http://1.com"),
// 		tgbotapi.NewInlineKeyboardButtonData("2", "2"),
// 		tgbotapi.NewInlineKeyboardButtonData("3", "3"),
// 	),
// 	tgbotapi.NewInlineKeyboardRow(
// 		tgbotapi.NewInlineKeyboardButtonData("4", "4"),
// 		tgbotapi.NewInlineKeyboardButtonData("5", "5"),
// 		tgbotapi.NewInlineKeyboardButtonData("6", "6"),
// 	),
// 	tgbotapi.NewInlineKeyboardRow(
// 		tgbotapi.NewInlineKeyboardButtonData("7", "7"),
// 		tgbotapi.NewInlineKeyboardButtonData("8", "8"),
// 		tgbotapi.NewInlineKeyboardButtonData("9", "9"),
// 	),
// )

type BotConfig struct {
	UpdateTimeout int    `envconfig:"TG_BOT_API_UPDATE_CFG_TIMEOUT"`
	BotDebug      bool   `envconfig:"TG_BOT_API_BOT_DEBUG_LVL"`
	Token         string `envconfig:"TG_BOT_API_TOKEN"`
	LogLevel      string `envconfig:"LOGGER_LVL"`
	IsTest        bool   `envconfig:"LOGGER_ISTEST"`
}

type Bot struct {
	telegramAPI   *tgbotapi.BotAPI
	updateHandler service.UpdateProcessor
	config        *BotConfig
	logger        logger.Logger
}

func main() {
	bot := new(Bot)

	bot.ReadConfig()

	bot.SetupLogger()

	bot.SetupTGapi()

	bot.logger.Info("Authorized on account", "name", bot.telegramAPI.Self.UserName)

	bot.SetupService()

	bot.listenForUpdates()
}

func (b *Bot) SetupLogger() {
	level := strings.ToTitle(b.config.LogLevel)

	var (
		targetWriteFile io.Writer
		err             error
	)

	if b.config.IsTest {
		targetWriteFile = os.Stdout
	} else {
		targetWriteFile, err = os.Open("file")
		log.Fatalf("cannot open the file, error: %v", err)
	}

	slogger := new(slog.Logger)

	switch level {

	case slog.LevelDebug.String():
		slogger = slog.New(
			slog.NewJSONHandler(targetWriteFile,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				}))

	case slog.LevelInfo.String():
		slogger = slog.New(
			slog.NewJSONHandler(targetWriteFile,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				}))
	}

	b.logger = logger.New(slogger)
}

func (b *Bot) ReadConfig() {
	config := new(BotConfig)

	err := envconfig.Process("", config)
	if err != nil {
		log.Fatalf("error while reading config: %v", err)
	}

	b.config = config
}

func (b *Bot) SetupTGapi() {
	b.logger.Debug("setting up api", "token", b.config.Token)

	api, err := tgbotapi.NewBotAPI(b.config.Token)
	if err != nil {
		b.logger.Fatal("[SetupTGapi] error while setting up api", botErrors.Err(err))
	}

	b.logger.Debug("setting up bot level", "is test", b.config.IsTest)

	api.Debug = b.config.BotDebug

	b.telegramAPI = api
}

func (b *Bot) SetupService() {
	var err error
	b.updateHandler, err = service.New(b.logger)
	if err != nil {
		b.logger.Fatal("cannot setup service", botErrors.Err(err))
	}
}

func (b *Bot) listenForUpdates() {
	b.logger.Debug("setting up update config")

	updatesConfig := tgbotapi.NewUpdate(0)
	updatesConfig.Timeout = b.config.UpdateTimeout

	updates := b.telegramAPI.GetUpdatesChan(updatesConfig)

	ctx, cancel := context.WithCancel(context.Background())

	b.logger.Debug("starting listening", "config", updatesConfig)

	go func() {
		for update := range updates {
			select {
			case <-ctx.Done():
				b.logger.Warn("ctx done, stopping reciever")
				b.telegramAPI.StopReceivingUpdates()
				return
			default:
				go b.Handle(update)
			}
		}
	}()

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGTERM)

	<-osSignal
	cancel()

	b.logger.Warn("got os signal, shutdown...")
	b.logger.Info("Bye, bye...")
}

func (b *Bot) Handle(update tgbotapi.Update) {
	b.logger.Debug("got update")

	msg := b.updateHandler.Process(&update)

	if update.CallbackQuery != nil {
		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
		if _, err := b.telegramAPI.Request(callback); err != nil {
			panic(err)
		}
	}

	b.telegramAPI.Send(msg)
}

// func (b *Bot) HandleCallback(update tgbotapi.Update) {
// 	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
// 	if _, err := b.telegramAPI.Request(callback); err != nil {
// 		panic(err)
// 	}

// 	// And finally, send a message containing the data received.
// 	msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "Привет, друг", numericKeyboard)
// 	if _, err := b.telegramAPI.Send(msg); err != nil {
// 		panic(err)
// 	}
// }
