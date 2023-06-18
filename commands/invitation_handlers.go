package commands

import (
	"github.com/bwmarrin/discordgo"
	"mapguess-discord/game"
)

var (
	invitationHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"invitation_join": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := game.JoinGame(i.ChannelID, i.Interaction)
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
				return
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: GetGameInvitationResponse(i.ChannelID),
			})
		},
		"invitation_start": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			g := game.GetGame(i.ChannelID)
			if g.MatchStarted {
				return
			}

			gi := NewGameInteractions(g)
			g.RegisterGameListener(gi)
			g.StartMatch()

			s.ChannelMessageDelete(i.ChannelID, i.Interaction.Message.ID)
		},
	}
)
