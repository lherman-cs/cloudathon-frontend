package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
)

const token string = "xoxb-342968164674-HsSe5A2b0hRutAZExCkBMpsZ"
const slackGatewayURL string = "https://slack.com/api/rtm.connect?token=" + token
const channelID string = "GA2N2HTFE"

type API struct {
	Conn *websocket.Conn
	Data map[string]interface{}
	Game *Game
}

type EventHandler func(api *API, event map[string]interface{})

func New(game *Game) (api *API, err error) {
	api = &API{}
	api.Game = game

	resp, err := http.Get(slackGatewayURL)
	if err != nil {
		return api, err
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&api.Data)
	if err != nil {
		return api, err
	}

	api.Conn, _, _ = websocket.DefaultDialer.Dial(api.Data["url"].(string), nil)
	return api, err
}

func (api *API) SendMessage(message string) error {
	return api.Conn.WriteJSON(map[string]interface{}{
		"id":      1,
		"type":    "message",
		"channel": channelID,
		"text":    message,
	})
}

func (api *API) Start(eventHandler EventHandler) {
	for {
		eventData := make(map[string]interface{})
		api.Conn.ReadJSON(&eventData)
		go eventHandler(api, eventData)
	}
}
