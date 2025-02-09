package discord

import (
	"log"

	"github.com/Sanjar0126/wiki_change_stream/config"
	"github.com/bwmarrin/discordgo"
)

func NewDiscord(cfg *config.Config, handler *Handler) *discordgo.Session {
	dg, err := discordgo.New("Bot " + cfg.DiscordBotToken)
	if err != nil {
		log.Fatal("Error creating Discord session: ", err)
	}

	dg.AddHandler(handler.MessageHandle)
	dg.AddHandler(handler.Ready)

	return dg
}
