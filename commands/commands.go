package commands

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"mapguess-discord/utils/phrases"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "start",
			Description: phrases.StartCommandDescription,
		},
	}
)

func RegisterCommands(dg *discordgo.Session) {
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionMessageComponent:
			if h, ok := invitationHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}
	})

	for _, command := range commands {
		_, err := dg.ApplicationCommandCreate(dg.State.User.ID, "", command)
		if err != nil {
			log.Panicf("Cannot create command %v", command, err)
		}
	}
}
