package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"mapguess-discord/api"
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
		"invitation_start": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			test := api.GetRandomPhoto()
			fmt.Println("\n")
			fmt.Println(test)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: test.Result.Data[0].FileUrl,
				},
			})
		},
	}
)
