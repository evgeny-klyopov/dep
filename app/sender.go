package app

import (
	"context"
	"golang.org/x/net/proxy"
	"gopkg.in/telegram-bot-api.v4"
	"net"
	"net/http"
	"strconv"
	"strings"

)

type ConfigSender struct {
	configNotifications Notifications
	messageProperties messageProperties
}

type messageProperties struct {
	Event string
	logTasks []LogTask
	Path string
	Stage string
	Name string
	Host string
	Port int
}

type Sender interface {
	Send() error
}

func NewSender(config Notifications, messageProperties messageProperties) Sender {
	var sender Sender
	sender = &ConfigSender{config, messageProperties}

	return sender
}

func (s *ConfigSender) Send() error {
	if s.configNotifications.Telegram != nil {
		telegram := telegram{}
		message := telegram.buildMessage(s.messageProperties)

		for _, configTelegram := range *s.configNotifications.Telegram {
			telegram.Config = configTelegram
			if err := telegram.Send(message);err != nil {
				return err
			}
		}
	}
	return nil
}

type telegram struct {
	Config Telegram
}
func(t *telegram) Send(message string) error {
	var bot *tgbotapi.BotAPI

	if t.Config.UseProxy {
		setting := strings.Split(t.Config.Proxy, "@")
		authData := strings.Split(setting[0], ":")

		dialer, err := proxy.SOCKS5(
			"tcp",
			setting[1],
			&proxy.Auth{User: authData[0], Password:  authData[1]},
			proxy.Direct,
		)
		if err != nil {
			return err
		}

		client := &http.Client{Transport: &http.Transport{DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		}}}

		bot, err = tgbotapi.NewBotAPIWithClient(t.Config.Token, client)
		if err != nil {
			return err
		}
	} else {
		var err error
		bot, err = tgbotapi.NewBotAPI(t.Config.Token)
		if err != nil {
			return err
		}
	}

	msg := tgbotapi.NewMessage(t.Config.ChatId, message)
	msg.ParseMode = "markdown"
	_, err := bot.Send(msg)

	if err != nil {
		return err
	}

	return nil
}
func(t *telegram) buildMessage(messageProperties messageProperties) string {
	messages := []string{
		"*" + messageProperties.Event + "*",
		"-------------------------------------",
		"*Host:* " + messageProperties.Host + ":" + strconv.Itoa(messageProperties.Port),
		"*Stage:* " + messageProperties.Stage,
		"*Release Number:* " + messageProperties.Name,
		"*Release Path:* " + messageProperties.Path,
		"",
	}

	messages = append(messages, "*Log commands:*")
	for _, info := range messageProperties.logTasks {
		messages = append(messages, "_" + info.Task + "_ - " + info.Time)
	}

	return strings.Join(messages, "\n")
}