package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"slices"

	"github.com/derkellernerd/dori/auth"
	"github.com/derkellernerd/dori/chat"
	"github.com/derkellernerd/dori/core"
	"github.com/derkellernerd/dori/database"
	"github.com/derkellernerd/dori/handler"
	"github.com/derkellernerd/dori/middleware"
	"github.com/derkellernerd/dori/model"
	"github.com/derkellernerd/dori/repository"
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

	chatChannel := make(chan model.ChatEvent)

	commandHandler := handler.NewCommand(env, commandRepo)
	eventHandler := handler.NewEvent(env, chatChannel)

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

	chat, err := chat.NewChat(env, commandRepo, chatChannel)

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
			event.GET("chat", eventHandler.ChatEventHandler)
			event.POST("chat", eventHandler.ChatEventTest)
		}
	}

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
