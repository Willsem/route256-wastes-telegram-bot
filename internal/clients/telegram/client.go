package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/pkg/log"
)

type Config struct {
	Token         string `yaml:"token"`
	Timeout       int    `yaml:"timeout"`
	MessageBuffer int    `yaml:"message_buffer"`
}

type Client struct {
	client *tgbotapi.BotAPI
	logger log.Logger

	timeout int

	messageUpdates chan *models.Message
}

func NewClient(config Config, logger log.Logger) (*Client, error) {
	client, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to connecting to telegrag bot: %v", err)
	}

	c := &Client{
		client:         client,
		logger:         logger,
		timeout:        config.Timeout,
		messageUpdates: make(chan *models.Message, config.MessageBuffer),
	}

	go c.listenUpdates()

	return c, nil
}

func (c *Client) SendMessage(userID int64, text string) error {
	msg := tgbotapi.NewMessage(userID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown

	_, err := c.client.Send(msg)
	if err != nil {
		return fmt.Errorf("sending message to telegram: %v", err)
	}
	return nil
}

func (c *Client) GetUpdatesChan() chan *models.Message {
	return c.messageUpdates
}

func (c *Client) listenUpdates() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = c.timeout

	updates := c.client.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			msg := update.Message
			usr := msg.From
			c.logger.Debugf("[%s] %s", usr.UserName, msg.Text)

			c.messageUpdates <- models.NewMessage(
				msg.MessageID,
				models.NewUser(usr.ID, usr.FirstName, usr.LastName, usr.UserName),
				msg.Date, msg.Text,
			)
		}
	}
}
