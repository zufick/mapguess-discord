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
	m, _ := game.DiscordSession.ChannelMessageSendComplex(gi.game.ChannelId, gi.GetCurrentRoundResponse())
	gi.game.CurrentRound.SetMessage(m)
	AddCountryReactions(m, gi.game.CurrentRound.CountryOptions)
}
func (gi *GameInteractions) OnRoundEnd() {
	game.DiscordSession.ChannelMessageDelete(gi.game.ChannelId, gi.game.CurrentRound.Message.ID)
	m, _ := game.DiscordSession.ChannelMessageSendComplex(gi.game.ChannelId, gi.GetRoundEndResponse())
	gi.game.CurrentRound.ResultMessage = m
	game.DiscordSession.MessageReactionAdd(m.ChannelID, m.ID, "▶")
}
func (gi *GameInteractions) OnGameEnd() {}

func OnMessageReaction(s *discordgo.Session, mr *discordgo.MessageReactionAdd) {
	g := game.GetGame(mr.ChannelID)

	if g == nil {
		return
	}

	if mr.MessageID == g.CurrentRound.Message.ID {
		g.SetUserAnswer(mr.Member.User.ID, mr.Emoji.Name)
		return
	}

	if mr.MessageID == g.CurrentRound.ResultMessage.ID && mr.Emoji.Name == "▶" {
		reactions, _ := s.MessageReactions(mr.ChannelID, mr.MessageID, "▶", 100, "", "")
		for _, user := range g.Users {
			var reacted bool
			for _, reactionUser := range reactions {
				if user.Profile.ID == reactionUser.ID {
					reacted = true
				}
			}
			if reacted == false {
				return
			}
		}
		s.ChannelMessageDelete(g.ChannelId, g.CurrentRound.ResultMessage.ID)
		g.StartRound()
		return
	}
}

func NewGameInteractions(g *game.Game) *GameInteractions {
	game.DiscordSession.AddHandler(OnMessageReaction)

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

func (gi *GameInteractions) GetCurrentRoundResponse() *discordgo.MessageSend {
	description := "Select one country:"
	for _, c := range gi.game.CurrentRound.CountryOptions {
		symbol, _ := utils.GetCountryCodeSymbol(c)
		name := utils.GetCountryName(c)
		description += fmt.Sprintf("\n %s - %s \n", symbol, name)
	}

	return &discordgo.MessageSend{
		Content: "Guess where it is",
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Guess where it is",
				Description: description,
				Image:       &discordgo.MessageEmbedImage{URL: gi.game.CurrentRound.ImgUrl},
			},
		},
	}
}

func (gi *GameInteractions) GetRoundEndResponse() *discordgo.MessageSend {
	countryEmoji, _ := utils.GetCountryCodeSymbol(gi.game.CurrentRound.CorrectCountry)
	description := fmt.Sprintf("Round ended! Correct country: %s %s.\n", countryEmoji, utils.GetCountryName(gi.game.CurrentRound.CorrectCountry))

	if len(gi.game.CurrentRound.Winners) > 0 {
		description += "\nWinners:\n"
		for _, w := range gi.game.CurrentRound.Winners {
			description += fmt.Sprintf("%s (+1 point, total: %d)\n", w.Profile.Username, w.Score)
		}
	} else {
		description += "There are no winners\n"
	}

	description += "\nPress ▶ emoji to start the next round\n"

	return &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Round ended",
				Description: description,
				Image:       &discordgo.MessageEmbedImage{URL: gi.game.CurrentRound.ImgUrl},
			},
		},
	}
}

func AddCountryReactions(m *discordgo.Message, countryCodes []string) {
	for _, c := range countryCodes {
		cSymbol, _ := utils.GetCountryCodeSymbol(c)
		game.DiscordSession.MessageReactionAdd(m.ChannelID, m.ID, cSymbol)
	}
}
