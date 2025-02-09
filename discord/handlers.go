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
	"github.com/spf13/cast"
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

func (h *Handler) MemberJoin(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	log.Printf("New member joined: %s#%s", m.User.Username, m.User.Discriminator)

	welcomeEmbed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("Welcome to %s!", m.GuildID),
		Description: "Thanks for joining! Here's some information to help you get started:",
		Color:       config.BotInfoColor,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "SetLang command",
				Value:  "type !setLang <lang code>, e.g. !setLang en (e.g., ru, fr, es, etc.). Default is en",
				Inline: false,
			},
			{
				Name:   "Recent command",
				Value:  "type !recent for getting recent wiki changes for the current language, you can enter offset and limit !recent <offset> <limit>",
				Inline: false,
			},
			{
				Name:   "Daily stats command",
				Value:  "type !stats [yyyy-mm-dd] to display how many changes occurred on that date for the chosen language",
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "We hope you enjoy your stay!",
		},
	}

	dmChannel, err := s.UserChannelCreate(m.User.ID)
	if err != nil {
		log.Printf("Error creating DM channel: %v", err)
		return
	}

	_, err = s.ChannelMessageSendEmbed(dmChannel.ID, welcomeEmbed)
	if err != nil {
		log.Printf("Error sending welcome DM: %v", err)
		return
	}

	log.Printf("Sent welcome DM to %s", m.User.Username)
}

func (h *Handler) MessageHandle(s *discordgo.Session, m *discordgo.MessageCreate) {
	var (
		commandArg string
	)

	if m.Author.ID == s.State.User.ID {
		return
	}

	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		channel, err = s.Channel(m.ChannelID)
		if err != nil {
			log.Printf("Error getting channel: %v", err)
			return
		}
	}

	if channel.Type != discordgo.ChannelTypeDM {
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

		_, err = s.ChannelMessageSend(m.ChannelID, response)

	case "stats":
		date := commandArg

		discordUser, err := h.db.DiscordUser().GetOrCreate(
			context.Background(), m.Author.ID)
		if err != nil {
			log.Println("error while getting user from db", err)
			return
		}

		count, err := h.db.WikiChanges().GetCountDate(date, discordUser.Lang)
		if err != nil {
			log.Println("error while getting changes count from db", err)
			return
		}

		response := fmt.Sprintf("Changes for %s %s lang: %d", date, discordUser.Lang, count)

		_, err = s.ChannelMessageSend(m.ChannelID, response)
	case "recent":
		var (
			offset int64 = 0
			limit  int64 = 10
		)

		discordUser, err := h.db.DiscordUser().GetOrCreate(
			context.Background(), m.Author.ID)
		if err != nil {
			log.Println("error while getting user from db", err)
			return
		}

		if len(args) > 1 {
			offset = cast.ToInt64(args[1])
		}

		if len(args) > 2 {
			limit = cast.ToInt64(args[2])
		}

		wikiChanges, _, err := h.db.WikiChanges().GetAll(
			context.Background(), offset, limit, discordUser.Lang)
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
