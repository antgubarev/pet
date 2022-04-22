package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const ContentTypeJson = "application/json"

type Update struct {
	UpdateId int64 `json:"update_id"`
	Message  struct {
		Text string `json:"text"`
		Chat struct {
			Id        int64  `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		} `json:"chat"`
	} `json:"message"`
}

type SetWebHookRequest struct {
	Url            string   `json:"url"`
	AllowedUpdates []string `json:"allowed_updates"`
}

type sendMessageRequest struct {
	ChatId int64  `json:"chat_id"`
	Text   string `json:"text"`
}

type Bot struct {
	apiUrl     string
	lastUpdate map[int64]*Update
}

func NewBot(token string) *Bot {
	return &Bot{
		apiUrl:     "https://api.telegram.org/bot" + token,
		lastUpdate: make(map[int64]*Update),
	}
}

func (b *Bot) InitWebHook(url string) error {
	req := struct {
		Url            string   `json:"url"`
		AllowedUpdates []string `json:"allowed_updates"`
	}{
		Url:            url,
		AllowedUpdates: []string{"message"},
	}

	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal setWebhook request: %v", err)
	}

	resp, err := http.Post(b.apiUrl+"/setWebHook", ContentTypeJson, bytes.NewBuffer(data))
	if checkErr := checkTelegramApiResp(err, *resp); checkErr != nil {
		return fmt.Errorf("setWebHook: %v", checkErr)
	}

	return nil
}

func (b *Bot) InitCmds() error {
	type Command struct {
		Command     string `json:"command"`
		Description string `json:"description"`
	}
	setMyCommandsRequest := struct {
		Commands []Command `json:"commands"`
	}{
		Commands: []Command{
			{
				Command:     "count",
				Description: "Counts words in text",
			},
		},
	}

	data, err := json.Marshal(setMyCommandsRequest)
	if err != nil {
		return fmt.Errorf("marshal setMyCommands request: %v", err)
	}

	resp, err := http.Post(b.apiUrl+"/setMyCommands", ContentTypeJson, bytes.NewBuffer(data))
	if checkErr := checkTelegramApiResp(err, *resp); checkErr != nil {
		return fmt.Errorf("setWebHook: %v", checkErr)
	}

	return nil
}

func (b *Bot) sendMessage(chatId int64, message string) error {
	sendMessageRequest := &sendMessageRequest{
		ChatId: chatId,
		Text:   message,
	}

	data, err := json.Marshal(sendMessageRequest)
	if err != nil {
		return fmt.Errorf("marshal sendMessageRequest: %v", err)
	}

	resp, err := http.Post(b.apiUrl+"/sendMessage", ContentTypeJson, bytes.NewBuffer(data))
	if checkErr := checkTelegramApiResp(err, *resp); checkErr != nil {
		return fmt.Errorf("sendMessageRequest: %v", checkErr)
	}

	return nil
}

func (b *Bot) StartCmdEcho(update *Update) error {
	if err := b.sendMessage(
		update.Message.Chat.Id,
		fmt.Sprintf("Hello, %s %s! Let's start. Choose a command.",
			update.Message.Chat.FirstName,
			update.Message.Chat.LastName)); err != nil {
		return err
	}

	b.setChatLastUpdate(update)

	return nil
}

func (b *Bot) CountCmd(update *Update) error {
	if err := b.sendMessage(
		update.Message.Chat.Id,
		"Input your text"); err != nil {
		return err
	}

	b.setChatLastUpdate(update)

	return nil
}

func (b *Bot) CountInput(update *Update) error {
	count := CountWordsInString(update.Message.Text)
	if err := b.sendMessage(
		update.Message.Chat.Id,
		fmt.Sprintf("Your text has %d words", count)); err != nil {
		return err
	}
	return nil
}

func (b *Bot) setChatLastUpdate(update *Update) {
	b.lastUpdate[update.Message.Chat.Id] = update
}

func (b *Bot) WebHook(update *Update) error {
	if update.Message.Text == "/start" {
		b.StartCmdEcho(update)
		return nil
	}
	if update.Message.Text == "/count" {
		b.CountCmd(update)
		return nil
	}

	//check last
	last, ok := b.lastUpdate[update.Message.Chat.Id]
	if ok && last.Message.Text == "/count" {
		b.CountInput(update)
		return nil
	}

	if err := b.sendMessage(
		update.Message.Chat.Id,
		"Please choose a command before"); err != nil {
		return err
	}

	return nil
}

func checkTelegramApiResp(err error, resp http.Response) error {
	if err != nil {
		return fmt.Errorf("send request to telegram api: %v", err)
	}
	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("wrong resonse code: %d, body: %s", resp.StatusCode, body)
	}
	return nil
}
