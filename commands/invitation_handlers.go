package commands

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"mapguess-discord/game"
)

var (
	invitationHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"invitation_join": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := game.JoinGame(i.ChannelID, i.Interaction)
			if err != nil {
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				if err != nil {
					log.Fatal(err)
				}

				return
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: GetGameInvitationResponse(i.ChannelID),
			})
			if err != nil {
				log.Fatal(err)
			}
		},
		"invitation_start": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			g := game.GetGame(i.ChannelID)
			if g.MatchStarted {
				return
			}

			gi := NewGameInteractions(g)
			g.RegisterGameListener(gi)
			g.StartMatch()

			err := s.ChannelMessageDelete(i.ChannelID, i.Interaction.Message.ID)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
)
