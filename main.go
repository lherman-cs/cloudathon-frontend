package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

const (
	LogFile   = "log.txt"
	JoinDelay = 10 // in seconds
)

var i int = 0

func uploadFileHandler(api *API, event map[string]interface{}) {
	h := DefaultHelper()
	info, err := h.GetFileInfo(event["file_id"].(string))
	if err != nil {
		Error(err)
		return
	}

	user := info["user"].(string)
	var url string

	if i == 0 {
		url = "https://i.imgur.com/LHqpqYs.jpg"
	} else {
		url = "https://i.imgur.com/uuvS3uI.jpg"
	}
	i = (i + 1) % 2

	result, _ := api.Game.Helper.Request(CheckImageURL, map[string]string{
		"url":     url,
		"user_id": user,
	})

	if result != nil && result["body"] != nil {
		api.SendMessage(result["body"].(string))
	} else {
		api.SendMessage("I'm sorry. I can't analyze your image")
	}
}

func messageHandler(api *API, event map[string]interface{}) {
	if event == nil || event["user"] == nil || event["text"] == nil {
		return
	}

	user := event["user"].(string)
	text := event["text"].(string)
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)

	switch {
	case text == "start":
		api.SendMessage(api.Game.Start())
		if !api.Game.Info.IsCounterStarted {
			go func(api *API) {
				var counter int64 = JoinDelay
				for ; counter >= 0; counter-- {
					if api.Game.Info.IsRunning == false {
						return
					}
					api.SendMessage(strconv.FormatInt(counter, 10))
					time.Sleep(time.Second)
				}

				api.Game.Mutex.Lock()
				api.Game.Info.CanJoin = false
				api.Game.Mutex.Unlock()
				api.SendMessage("The game has begun! Hunt!")

				gameHandler(api)
			}(api)
			api.Game.Info.IsCounterStarted = true
		}
	case text == "join":
		if api.Game.Info.CanJoin {
			msg, err := api.Game.Join(user)
			if err != nil {
				Error(err)
			} else {
				api.SendMessage(msg)
			}
		} else {
			api.SendMessage("Game is running. No new players are allowed")
		}

	case text == "stop":
		api.SendMessage(api.Game.Stop())
		i = 0
		api.Game = DefaultGame()
	}

}

// gameHandler is a handler to keep giving hints
func gameHandler(api *API) {
	data, _ := api.Game.Helper.Request(StartURL, nil)

	if data != nil {
		api.SendMessage(data["body"].(string))
	}
}

func eventHandler(api *API, event map[string]interface{}) {
	eventType := event["type"]

	switch {
	case eventType == "file_shared":
		uploadFileHandler(api, event)
	case eventType == "message":
		messageHandler(api, event)
	}
	Debug(event)
	Debug("=====================================")
}

func main() {
	f, err := os.Create(LogFile)
	if err != nil {
		log.Fatalln("Can't create a log file")
	}
	defer f.Close()
	log.SetOutput(f)

	api, _ := New(DefaultGame())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(api *API, c chan os.Signal) {
		<-c
		i = 0
		api.SendMessage(api.Game.Stop())
		os.Exit(0)
	}(api, c)

	api.Start(eventHandler)
}
