package chat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/inetspa/golib/requests"
	"github.com/inetspa/golib/web"
	"net/http"
)

const (
	apiEndpoint = "https://chat-api.one.th/message/api/v1"
)

func NewClient(botId string, token string, tokenType string) Client {
	return Client{
		botId:       botId,
		token:       token,
		tokenType:   tokenType,
		apiEndpoint: apiEndpoint,
	}
}

func (c *Client) FindChatFriend(keyword string) (Friend, error) {
	var friend Friend
	msg := struct {
		BotId   string `json:"bot_id"`
		Keyword string `json:"key_search"`
	}{
		BotId:   c.botId,
		Keyword: keyword,
	}
	body, _ := json.Marshal(&msg)
	r, err := c.send(http.MethodPost, c.url("/searchfriend"), body)
	if err != nil {
		return friend, nil
	}
	if r.Code != http.StatusOK {
		return friend, errors.New(fmt.Sprintf("client return error with code %d (%s)", r.Code, string(r.Body)))
	}
	chatFriendResult := struct {
		Status string `json:"status"`
		Friend Friend `json:"friend"`
	}{}
	if err := json.Unmarshal(r.Body, &chatFriendResult); err != nil {
		return friend, err
	}
	friend = chatFriendResult.Friend
	return friend, nil
}

func (c *Client) PushTextMessage(to string, msg string, customNotify *string) error {
	pushMessage := struct {
		To           string `json:"to"`
		BotId        string `json:"bot_id"`
		Type         string `json:"type"`
		Message      string `json:"message"`
		CustomNotify string `json:"custom_notification,omitempty"`
	}{
		To:      to,
		BotId:   c.botId,
		Type:    "text",
		Message: msg,
	}
	if customNotify != nil {
		pushMessage.CustomNotify = *customNotify
	}
	body, _ := json.Marshal(&pushMessage)
	r, err := c.send(http.MethodPost, c.url("/push_message"), body)
	if err != nil {
		return err
	}
	if r.Code != http.StatusOK {
		return errors.New(fmt.Sprintf("client return error with code %d (%s)", r.Code, string(r.Body)))
	}
	return nil
}

func (c *Client) PushWebView(to string, label string, path string, img string, title string, detail string, customNotify *string) error {
	pushMessage := struct {
		To           string     `json:"to"`
		BotId        string     `json:"bot_id"`
		Type         string     `json:"type"`
		CustomNotify string     `json:"custom_notification,omitempty"`
		Elements     []Elements `json:"elements"`
	}{
		To:    to,
		BotId: c.botId,
		Type:  "template",
		Elements: []Elements{
			{
				Image:  img,
				Title:  title,
				Detail: detail,
				Choices: []Choice{
					{
						Label: label,
						Type:  "webview",
						Url:   path,
						Size:  "full",
					},
				},
			},
		},
	}

	if customNotify != nil {
		pushMessage.CustomNotify = *customNotify
	}
	body, _ := json.Marshal(&pushMessage)
	r, err := c.send(http.MethodPost, c.url("/push_message"), body)
	if err != nil {
		return err
	}
	if r.Code != http.StatusOK {
		return errors.New(fmt.Sprintf("server return error with http code %d : %s", r.Code, string(r.Body)))
	}
	return nil
}

func (c *Client) PushLink(to string, label string, path string, img string, title string, detail string, customNotify *string) error {
	pushMessage := struct {
		To           string     `json:"to"`
		BotId        string     `json:"bot_id"`
		Type         string     `json:"type"`
		CustomNotify string     `json:"custom_notification,omitempty"`
		Elements     []Elements `json:"elements"`
	}{
		To:    to,
		BotId: c.botId,
		Type:  "template",
		Elements: []Elements{
			{
				Image:  img,
				Title:  title,
				Detail: detail,
				Choices: []Choice{
					{
						Label: label,
						Type:  "link",
						Url:   path,
					},
				},
			},
		},
	}

	if customNotify != nil {
		pushMessage.CustomNotify = *customNotify
	}
	body, _ := json.Marshal(&pushMessage)
	r, err := c.send(http.MethodPost, c.url("/push_message"), body)
	if err != nil {
		return err
	}
	if r.Code != http.StatusOK {
		return errors.New(fmt.Sprintf("server return error with http code %d : %s", r.Code, string(r.Body)))
	}
	return nil
}

func (c *Client) PushQuickReply(to string, message string, quickReply []QuickReply) error {
	pushQuickReply := struct {
		To         string       `json:"to"`
		BotId      string       `json:"bot_id"`
		Message    string       `json:"message"`
		QuickReply []QuickReply `json:"quick_reply"`
	}{
		To:         to,
		BotId:      c.botId,
		Message:    message,
		QuickReply: quickReply,
	}
	body, _ := json.Marshal(&pushQuickReply)
	r, err := c.send(http.MethodPost, c.url("/push_quickreply"), body)
	if err != nil {
		return err
	}
	if r.Code != http.StatusOK {
		return errors.New(fmt.Sprintf("server return error with http code %d : %s", r.Code, string(r.Body)))
	}
	return nil
}

func (c *Client) GetChatProfile(oneChatToken string) (Profile, error) {
	var chatProfile Profile
	msg := struct {
		BotId        string `json:"bot_id"`
		OneChatToken string `json:"source"`
	}{
		BotId:        c.botId,
		OneChatToken: oneChatToken,
	}
	body, _ := json.Marshal(&msg)
	r, err := c.send(http.MethodPost, "https://chat-api.one.th/manage/api/v1/getprofile", body)
	if err != nil {
		return chatProfile, err
	}
	chatProfileResult := struct {
		Data   Profile `json:"data"`
		Status string  `json:"status"`
	}{}
	if err := json.Unmarshal(r.Body, &chatProfileResult); err != nil {
		return chatProfile, err
	}
	return chatProfileResult.Data, nil
}

func (c *Client) SetEndpoint(ep string) {
	c.apiEndpoint = ep
}

func (c *Client) send(method string, url string, body []byte) (requests.Response, error) {
	headers := map[string]string{
		web.HeaderContentType:   web.MIMEApplicationJSON,
		web.HeaderAuthorization: fmt.Sprintf("%s %s", c.tokenType, c.token),
	}
	r, err := requests.Request(method, url, headers, bytes.NewBuffer(body), 0)
	if err != nil {
		return r, err
	}
	return r, nil
}

func (c *Client) url(path string) string {
	return fmt.Sprintf("%s%s", c.apiEndpoint, path)
}
