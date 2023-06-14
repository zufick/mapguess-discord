package commands

import (
	"github.com/bwmarrin/discordgo"
	"mapguess-discord/game"
)

var (
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"start": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := game.StartGame(i.ChannelID)
			if err != nil {
				s.InteractionRespond(i.Interaction, game.GetGameExistsResponse())
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: game.GetGameInvitationResponse(i.ChannelID),
			})
		},
	}
)
