package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"mapguess-discord/game"
	"mapguess-discord/utils"
)

type GameInteractions struct {
	game *game.Game
}

func (gi *GameInteractions) OnRoundStart() {
	for _, u := range gi.game.Users {
		fmt.Println("\nSending interaction for " + u.Profile.Username + "\n")
		//r := gi.GetCurrentRoundResponse()
		game.DiscordSession.InteractionRespond(u.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "test for " + u.Profile.Username,
			},
		})
	}
}
func (gi *GameInteractions) OnRoundEnd() {}
func (gi *GameInteractions) OnGameEnd()  {}

func NewGameInteractions(g *game.Game) *GameInteractions {
	return &GameInteractions{
		game: g,
	}
}

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
	game := game.GetGame(channelId)

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

	if len(game.Users) > 1 {
		buttons = append(buttons, discordgo.Button{
			Emoji: discordgo.ComponentEmoji{
				Name: "▶",
			},
			Label:    "Start match",
			Style:    discordgo.SecondaryButton,
			CustomID: "invitation_start",
		})
	}

	if len(game.Users) > 0 {
		description += "\nJoined players:\n"

		for _, v := range game.Users {
			description += "@" + v.Profile.Username + "\n"
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

func (gi *GameInteractions) GetCurrentRoundResponse() *discordgo.InteractionResponseData {
	buttons := []discordgo.MessageComponent{}

	for _, c := range gi.game.Round.CountryOptions {
		AddAnswerHandler("selectanswer_"+c, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "ANSWERED " + c,
				},
			})
		})

		cc, _ := utils.GetCountryCodeSymbol(c)
		buttons = append(buttons, discordgo.Button{
			Emoji: discordgo.ComponentEmoji{
				Name: cc,
			},
			Label:    utils.GetCountryName(c),
			Style:    discordgo.SecondaryButton,
			CustomID: "selectanswer_" + c,
		})
	}

	return &discordgo.InteractionResponseData{
		Content: "Guess where it is",
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Guess where it is",
				Image: &discordgo.MessageEmbedImage{URL: gi.game.Round.ImgUrl},
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: buttons,
			},
		},
		Flags: discordgo.MessageFlagsEphemeral,
	}
}
