package main

import (
	"flag"
	"log"

	tgClient "notekeeper-tgbot/clients/telegram"
	eventConsumer "notekeeper-tgbot/consumer/event-consumer"
	"notekeeper-tgbot/events/telegram"
	"notekeeper-tgbot/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "./files"
	batchSize   = 100
)

func main() {
	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()), files.New(storagePath),
	)

	log.Print("service started")

	consumer := eventConsumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service stopped", err)
	}

}

func mustToken() string {
	token := flag.String(
		"bot-token", "", "Token for Telegram API",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token not specified")
	}

	return *token
}
