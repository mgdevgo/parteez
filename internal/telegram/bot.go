package telegram

import (
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	"gopkg.in/telebot.v4"
)

type Bot struct {
	bot        *telebot.Bot
	token      string
	appName    string
	webAppURL  string
	publicHost string
}

func NewBot(token, appName, webAppURL, publicHost string) (*Bot, error) {
	return &Bot{
		token:     token,
		appName:   appName,
		webAppURL: webAppURL,
	}, nil
}

func (b *Bot) Run() error {
	const op = "telegram.bot.Run"

	settings := telebot.Settings{
		Token: b.token,
		// Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	tbot, err := telebot.NewBot(settings)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	b.bot = tbot

	b.bot.Handle("/start", func(ctx telebot.Context) error {
		return ctx.Send(
			"Привет! Добро пожаловать в бот бронирования отелей! Надеюсь, вам понравится мое приложение",
			&telebot.ReplyMarkup{
				InlineKeyboard: [][]telebot.InlineButton{
					{
						telebot.InlineButton{
							Text: "Открыть приложение",
							WebApp: &telebot.WebApp{
								URL: b.webAppURL,
							},
						},
					},
				},
			})
	})

	b.bot.Handle("/help", func(ctx telebot.Context) error {
		return ctx.Send(
			"Я всего лишь простой бот, поэтому все, что я могу сделать, это отправить ссылку на мини-приложение.",
			&telebot.ReplyMarkup{
				InlineKeyboard: [][]telebot.InlineButton{
					{
						telebot.InlineButton{
							Text: "Открыть приложение",
							WebApp: &telebot.WebApp{
								URL: b.webAppURL,
							},
						},
					},
				},
			})
	})

	b.bot.Start()

	webhook, err := b.bot.Webhook()
	if err != nil {
		return fmt.Errorf("%s: get webhook status error: %w", op, err)
	}

	if webhook.Endpoint.PublicURL != "" {
		log.Info("🤖 webhook info: ", webhook)
	} else {
		b.setWebhook()
	}


	return nil
}

func (b *Bot) setWebhook() error {
	const op = "telegram.bot.setWebhook"
	log.Infof("🤖 setting a webhook to @BotFather: %s/bot", b.publicHost)

	webhook := &telebot.Webhook{
		Endpoint: &telebot.WebhookEndpoint{
			PublicURL: b.publicHost + "/bot",
		},
	}

	if err := b.bot.SetWebhook(webhook); err != nil {
		log.Warn("🤖 webhook not set", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("🤖 webhook set")
	return nil
}

func (b *Bot) onMessage(ctx telebot.Context) error {
	panic("not implemented")
}
