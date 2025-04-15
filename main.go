package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/derkellernerd/kellerbot/auth"
	"github.com/derkellernerd/kellerbot/chat"
	"github.com/derkellernerd/kellerbot/core"
	"github.com/derkellernerd/kellerbot/database"
	"github.com/derkellernerd/kellerbot/handler"
	"github.com/derkellernerd/kellerbot/middleware"
	"github.com/derkellernerd/kellerbot/repository"
	"github.com/derkellernerd/kellerbot/worker"
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

	actionRepo := repository.NewAction(env)
	err = actionRepo.Migrate()
	if err != nil {
		panic(err)
	}

	twitchEventRepo := repository.NewTwitchEvent(env)
	err = twitchEventRepo.Migrate()
	if err != nil {
		panic(err)
	}

	chatCommandRepo := repository.NewChatCommand(env)
	err = chatCommandRepo.Migrate()
	if err != nil {
		panic(err)
	}

	eventRepo := repository.NewEvent(env)
	err = eventRepo.Migrate()
	if err != nil {
		panic(err)
	}

	actionWorker := worker.NewAction(env, actionRepo)
	actionHandler := handler.NewAction(env, actionRepo)
	twitchEventHandler := handler.NewTwitchEvent(env, twitchEventRepo, actionWorker)
	chatCommandHandler := handler.NewChatCommand(env, chatCommandRepo, actionWorker)
	eventHandler := handler.NewEvent(env, actionWorker, eventRepo)

	r := gin.Default()
	r.Use(middleware.AcceptCors)

	chat, err := chat.NewChat(env, twitchEventRepo, chatCommandRepo, actionWorker, eventRepo)

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
		event := apiV1.Group("event")
		{
			event.GET("", eventHandler.EventGetAll)
			event.POST(":eventId/action/replay", eventHandler.EventActionReplay)

			twitch := event.Group("twitch")
			{
				twitch.GET("", twitchEventHandler.TwitchEventGetAll)
				twitch.POST("", twitchEventHandler.TwitchEventCreate)
				twitch.PUT(":twitchEventId", twitchEventHandler.TwitchEventUpdate)
				twitch.DELETE(":twitchEventId", twitchEventHandler.TwitchEventDelete)
				twitch.POST(":twitchEventId/action/test", twitchEventHandler.TwitchEventTest)
			}

			event.GET("alert", middleware.HeadersMiddleware(), eventHandler.AlertEventHandler)
		}

		action := apiV1.Group("action")
		{
			action.GET("", actionHandler.ActionGetAll)
			action.GET(":actionId", actionHandler.ActionGetById)
			action.POST("", actionHandler.ActionCreate)
			action.PUT(":actionId", actionHandler.ActionUpdate)
			action.POST(":actionId/upload", actionHandler.ActionUploadFile)
			action.DELETE(":actionId", actionHandler.ActionDelete)
		}

		chatCommand := apiV1.Group("chat_command")
		{
			chatCommand.GET("", chatCommandHandler.ChatCommandGetAll)
			chatCommand.POST("", chatCommandHandler.ChatCommandCreate)
			chatCommand.POST(":chatCommandId/action/test", chatCommandHandler.ChatCommandTest)
			chatCommand.PUT(":chatCommandId", chatCommandHandler.ChatCommandUpdate)
			chatCommand.DELETE(":chatCommandId", chatCommandHandler.ChatCommandDelete)
		}
	}

	r.GET("action/:actionId", middleware.HeadersNoCache(), actionHandler.ActionGetFile)

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
