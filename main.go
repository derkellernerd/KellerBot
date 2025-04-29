package main

import (
	"log"

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

	databaseManager := database.NewDatabaseManager()
	env := core.NewEnvironment()
	env.DatabaseManager = databaseManager

	userRepo := repository.NewUser(env)
	err := userRepo.Migrate()
	if err != nil {
		panic(err)
	}

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

	twitchAuthHandler := handler.NewTwitchAuth(env, userRepo)
	actionWorker := worker.NewAction(env, actionRepo)
	actionHandler := handler.NewAction(env, actionRepo)
	twitchEventHandler := handler.NewTwitchEvent(env, twitchEventRepo, actionWorker)
	chatCommandHandler := handler.NewChatCommand(env, chatCommandRepo, actionWorker)
	eventHandler := handler.NewEvent(env, actionWorker, eventRepo)

	r := gin.Default()
	r.Use(middleware.AcceptCors)

	apiV1 := r.Group("api/v1")
	{
		twitch := apiV1.Group("twitch")
		{
			twitch.GET("login", twitchAuthHandler.Login)
			twitch.GET("callback", twitchAuthHandler.Callback)
		}

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
