package telegram

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"read-adviser-bot/lib/e"
	"read-adviser-bot/storage"
	"strconv"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	p.logger.Info(
		"got new command",
		slog.String("from", username),
		slog.String("chat_id", strconv.Itoa(chatID)),
		slog.String("text", text),
	)

	// add page: https://...
	// rnd page: /rnd
	// help: /help
	// start: /start: hi + help

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() {
		err = e.WrapIfErr("can't do command: save page", err)
	}()

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	exists, err := p.storage.Exists(page)
	if err != nil {
		return err
	}
	if exists {
		if err = p.tg.SendMessage(chatID, msgAlreadyExists); err != nil {
			return err
		}
		return fmt.Errorf("page is already saved")
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() {
		err = e.WrapIfErr("can't do command: send random", err)
	}()

	page, err := p.storage.PickRandom(username)
	if errors.Is(err, storage.ErrNoSavedPages) {
		p.logger.Error("there are no pages saved")
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}
	if err != nil {
		return err
	}

	if err = p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
