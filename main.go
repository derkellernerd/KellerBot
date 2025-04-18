package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"slices"

	"github.com/derkellernerd/kellerbot/auth"
	"github.com/derkellernerd/kellerbot/chat"
	"github.com/derkellernerd/kellerbot/core"
	"github.com/derkellernerd/kellerbot/database"
	"github.com/derkellernerd/kellerbot/handler"
	"github.com/derkellernerd/kellerbot/middleware"
	"github.com/derkellernerd/kellerbot/model"
	"github.com/derkellernerd/kellerbot/repository"
	"github.com/gin-gonic/gin"
)

const Version = "0.1.0"
const cacheFileName = "/home/sebastian/.cache/kellerbot.json"

func main() {
	log.Printf("Hallo, ich bin KellerBot der Twitch Bot: %s", Version)

	twitchSession, err := loadTwitchSession()
	if err != nil {
		panic(err)
	}

	databaseManager := database.NewDatabaseManager()
	env := core.NewEnvironment(twitchSession)
	env.DatabaseManager = databaseManager

	commandRepo := repository.NewCommand(env)
	err = commandRepo.Migrate()
	if err != nil {
		panic(err)
	}

	alertRepo := repository.NewAlert(env)
	err = alertRepo.Migrate()
	if err != nil {
		panic(err)
	}

	twitchEventRepo := repository.NewTwitchEvent(env)
	err = twitchEventRepo.Migrate()
	if err != nil {
		panic(err)
	}

	commandHandler := handler.NewCommand(env, commandRepo)
	eventHandler := handler.NewEvent(env, alertRepo)
	alertHandler := handler.NewAlert(env, alertRepo)
	twitchEventHandler := handler.NewTwitchEvent(env, twitchEventRepo)

	commands, err := commandRepo.CommandFindAll()

	wantedCommands := []model.Command{
		model.NewCommandMessage("ping", "pong"),
		model.NewCommandMessage("stats", "https://twitchtrends.tv/c/DerKellerNerd"),
		model.NewCommandMessage("futev", "https://futev.de"),
	}

	for _, command := range wantedCommands {
		exists := slices.ContainsFunc(commands, func(n model.Command) bool {
			return n.Command == command.Command
		})

		if !exists {
			err := commandRepo.CommandInsert(&command)
			if err != nil {
				panic(err)
			}
			log.Printf("Command created: %s", command.Command)
		}
	}

	r := gin.Default()
	r.Use(middleware.AcceptCors)

	chat, err := chat.NewChat(env, commandRepo, alertRepo, eventHandler, twitchEventRepo)

	_ = auth.NewTwitchAuth(env, r, func() {
		saveTwitchSession(env.TwitchSession)
		go chat.Start()
		if err != nil {
			panic(err)
		}
	})

	if twitchSession.IsAuthenticated() {
		log.Println("no login needed")
		go chat.Start()
		if err != nil {
			panic(err)
		}
	}

	apiV1 := r.Group("api/v1")
	{
		command := apiV1.Group("command")
		{
			command.GET("", commandHandler.CommandGetAll)
			command.POST("", commandHandler.CommandCreate)
			command.PUT(":commandId", commandHandler.CommandUpdate)
			command.DELETE(":commandId", commandHandler.CommandDelete)
		}

		event := apiV1.Group("event")
		{

			twitch := event.Group("twitch")
			{
				twitch.GET("", twitchEventHandler.TwitchEventGetAll)
				twitch.POST("", twitchEventHandler.TwitchEventCreate)
				twitch.PUT(":twitchEventId", twitchEventHandler.TwitchEventUpdate)
				twitch.DELETE(":twitchEventId", twitchEventHandler.TwitchEventDelete)
			}

			event.GET("status", eventHandler.Status)
			event.POST("chat", eventHandler.ChatEventTest)
			event.POST("alert", eventHandler.AlertEventTest)

			event.GET("chat", middleware.HeadersMiddleware(), eventHandler.ChatEventHandler)
			event.GET("alert", middleware.HeadersMiddleware(), eventHandler.AlertEventHandler)
		}

		alert := apiV1.Group("alert")
		{
			alert.GET("", alertHandler.AlertGetAll)
			alert.POST("", alertHandler.AlertCreate)
			alert.POST(":alertId/upload", alertHandler.AlertUploadFile)
		}
	}

	r.GET("alert/:alertId", middleware.HeadersNoCache(), alertHandler.AlertGetFile)

	r.Run()
}

func loadTwitchSession() (*core.TwitchSession, error) {
	twitchSession := &core.TwitchSession{}
	if _, err := os.Stat(cacheFileName); errors.Is(err, os.ErrNotExist) {
		return twitchSession, nil
	}

	fileContent, err := os.ReadFile(cacheFileName)
	if err != nil {
		return twitchSession, err
	}

	err = json.Unmarshal(fileContent, twitchSession)
	if err != nil {
		return twitchSession, err
	}

	log.Printf("Loaded Session from: %s", cacheFileName)

	return twitchSession, nil
}

func saveTwitchSession(twitchSession *core.TwitchSession) error {
	jsonBytes, err := json.Marshal(twitchSession)
	if err != nil {
		return err
	}

	err = os.WriteFile(cacheFileName, jsonBytes, 0600)

	log.Printf("Saved Session to: %s", cacheFileName)
	return err
}
