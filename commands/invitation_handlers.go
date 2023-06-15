package commands

import (
	"github.com/bwmarrin/discordgo"
	"mapguess-discord/game"
)

var (
	invitationHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"invitation_join": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			game.JoinGame(i.ChannelID, i.Interaction)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: GetGameInvitationResponse(i.ChannelID),
			})
		},
		"invitation_start": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			//game.StartMatch(i.ChannelID)
			//game.GetCurrentRoundResponse(i.ChannelID)
			g := game.GetGame(i.ChannelID)
			gi := NewGameInteractions(g)
			g.RegisterGameListener(gi)
			g.StartMatch()

			s.InteractionResponseDelete(i.Interaction)
		},
	}
)
