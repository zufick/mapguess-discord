package commands

import (
	"github.com/bwmarrin/discordgo"
	"mapguess-discord/game"
)

var (
	componentsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"invitation_join": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			game.JoinGame(i.ChannelID, i.Member.User)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: game.GetGameInvitationResponse(i.ChannelID),
			})
		},
	}
)
