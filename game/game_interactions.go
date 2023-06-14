package game

import "github.com/bwmarrin/discordgo"

func GetGameExistsResponse() *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: ":x: The game has already started",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}
}

func GetGameInvitationResponse(channelId string) *discordgo.InteractionResponseData {
	game := games[channelId]

	description := "New game invitation"
	buttons := []discordgo.MessageComponent{
		discordgo.Button{
			Emoji: discordgo.ComponentEmoji{
				Name: "➕",
			},
			Label:    "Join",
			Style:    discordgo.SecondaryButton,
			CustomID: "invitation_join",
		},
	}

	if len(game.users) > 1 {
		buttons = append(buttons, discordgo.Button{
			Emoji: discordgo.ComponentEmoji{
				Name: "▶",
			},
			Label:    "Start match",
			Style:    discordgo.SecondaryButton,
			CustomID: "invitation_start",
		})
	}

	if len(game.users) > 0 {
		description += "\nJoined players:\n"

		for _, v := range game.users {
			description += "@" + v.profile.Username + "\n"
		}
	}

	return &discordgo.InteractionResponseData{
		Content: "Starting the game...",
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Game invitation",
				Description: description,
				Image:       &discordgo.MessageEmbedImage{URL: "https://cdn.discordapp.com/attachments/630411106506113036/1118465338300964944/map_copy.png"},
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: buttons,
			},
		},
	}
}
