package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"mapguess-discord/game"
)

var (
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"start": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			gameErr := game.StartGame(i.ChannelID)
			if gameErr != nil {
				err := s.InteractionRespond(i.Interaction, GetGameExistsResponse())
				if err != nil {
					fmt.Println(err)
				}
			}

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: GetGameInvitationResponse(i.ChannelID),
			})
			if err != nil {
				fmt.Println(err)
			}
		},
	}
)
