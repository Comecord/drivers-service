package middlewares

import (
	"context"
	"drivers-service/config"
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/telegram"
)

func SendTelegram(title string, message string, config *config.Config) {

	telegramService, _ := telegram.New(config.Notify.TelegramToken)

	telegramService.AddReceivers(config.Notify.TelegramChatId)
	notify.UseServices(telegramService)

	// Send a test message.
	_ = notify.Send(
		context.Background(),
		"Subject/Title",
		"The actual message - Hello, you awesome gophers! :)",
	)

}
