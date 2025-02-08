package discord

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Sanjar0126/wiki_change_stream/config"
	"github.com/Sanjar0126/wiki_change_stream/models"
	"github.com/Sanjar0126/wiki_change_stream/storage"
	"github.com/bwmarrin/discordgo"
)

type Handler struct {
	config *config.Config
	db     storage.StorageI
}

type HandlerOptions struct {
	Config *config.Config
	DB     storage.StorageI
}

func NewHandler(opts *HandlerOptions) *Handler {
	return &Handler{
		config: opts.Config,
		db:     opts.DB,
	}
}

func (h *Handler) Ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Bot is ready! Logged in as: %v#%v\n", event.User.Username, event.User.Discriminator)
}

func (h *Handler) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	var (
		err        error
		commandArg string
	)

	if m.ChannelID == h.config.DiscordChannelId {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if !strings.HasPrefix(m.Content, config.BotCommandPrefix) {
			return
		}

		args := strings.Split(strings.TrimPrefix(m.Content, config.BotCommandPrefix), " ")
		if len(args) == 0 {
			return
		}

		command := strings.ToLower(args[0])

		if len(args) > 1 {
			commandArg = args[1]
		}

		switch command {
		case "ping":
			response := fmt.Sprintf("Pong! User ID: %s", m.Author.ID)
			_, err = s.ChannelMessageSend(m.ChannelID, response)
		case "setlang":
			lang := commandArg

			discordUser, err := h.db.DiscordUser().GetOrCreate(
				context.Background(), m.Author.ID)
			if err != nil {
				log.Println("error while getting user from db", err)
				return
			}

			_, err = h.db.DiscordUser().Update(context.Background(), models.DiscordUser{
				BId:      discordUser.BId,
				AuthorId: discordUser.AuthorId,
				Lang:     lang,
			})

			if err != nil {
				log.Println("error while updating user from db", err)
				return
			}

			response := fmt.Sprintf("Language set to %s", lang)
			log.Println(response)

			_, err = s.ChannelMessageSend(m.ChannelID, response)
		case "recent":
			discordUser, err := h.db.DiscordUser().GetOrCreate(
				context.Background(), m.Author.ID)
			if err != nil {
				log.Println("error while getting user from db", err)
				return
			}

			wikiChanges, _, err := h.db.WikiChanges().GetAll(
				context.Background(), 0, 10, discordUser.Lang)
			if err != nil {
				log.Println("error while getting wiki changes from db", err)
				return
			}

			embeds := []*discordgo.MessageEmbed{}

			for _, wikiChange := range wikiChanges {
				embed := &discordgo.MessageEmbed{
					Title: "Wiki recent changes",
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "ID",
							Value:  fmt.Sprintf("%d", wikiChange.ID),
							Inline: true,
						},
						{
							Name:   "Title",
							Value:  wikiChange.Title,
							Inline: true,
						},
						{
							Name:   "Title url",
							Value:  wikiChange.TitleURL,
							Inline: true,
						},
						{
							Name:   "Type",
							Value:  wikiChange.Type,
							Inline: true,
						},
						{
							Name:   "Comment",
							Value:  wikiChange.Comment,
							Inline: true,
						},
						{
							Name:   "User",
							Value:  wikiChange.User,
							Inline: true,
						},
						{
							Name:   "Is bot",
							Value:  fmt.Sprintf("%t", wikiChange.Bot),
							Inline: true,
						},
						{
							Name:   "Server name",
							Value:  wikiChange.ServerName,
							Inline: true,
						},
						{
							Name:   "Wiki type",
							Value:  wikiChange.Wiki,
							Inline: true,
						},
						{
							Name: "Timestamp",
							Value: time.Unix(
								int64(wikiChange.Timestamp), 0).Format("2006-01-02 15:04:05"),
							Inline: true,
						},
					},
					Color: config.BotInfoColor,
				}

				embeds = append(embeds, embed)
			}

			_, err = s.ChannelMessageSendEmbeds(m.ChannelID, embeds)
		}

		if err != nil {
			log.Printf("error while sending message %s, %v", command, err)
		}
	}
}
