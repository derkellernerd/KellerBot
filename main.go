package main

import (
	"log"
	"net/http"

	"github.com/derkellernerd/dori/auth"
	"github.com/derkellernerd/dori/chat"
	"github.com/derkellernerd/dori/core"
)

const Version = "0.1.0"

func main() {
	log.Printf("Hallo, ich bin KellerBot der Twitch Bot: %s", Version)

	env := core.NewEnvironment()
	chat, err := chat.NewChat(env)

	_ = auth.NewTwitchAuth(env, func() {
		chat.Start()
		if err != nil {
			panic(err)
		}
	})

	go http.ListenAndServe(":8080", nil)

	select {}
}
