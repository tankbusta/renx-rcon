package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	rcon "github.com/tankbusta/renx-rcon"
)

func main() {
	server := os.Getenv("GAME_SERVER_ADDRESS")
	rconPassword := os.Getenv("GAME_SERVER_RCON_PASSWORD")

	svr := rcon.NewServer(rconPassword, server)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := svr.Start(ctx); err != nil {
		log.Fatal(err)
	}
}
