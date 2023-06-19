package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"mapguess-discord/game"
	"mapguess-discord/utils/countries"
	"mapguess-discord/utils/phrases"
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
	err := game.DiscordSession.ChannelMessageDelete(gi.game.ChannelId, gi.game.CurrentRound.Message.ID)
	if err != nil {
		log.Fatal(err)
	}

	m, _ := game.DiscordSession.ChannelMessageSendComplex(gi.game.ChannelId, gi.GetRoundEndResponse())
	gi.game.CurrentRound.ResultMessage = m

	if gi.game.CurrentRoundNumber < game.MaxRounds {
		err = game.DiscordSession.MessageReactionAdd(m.ChannelID, m.ID, "▶")
		if err != nil {
			log.Fatal(err)
		}
	}
}
func (gi *GameInteractions) OnMatchEnd() {
	_, err := game.DiscordSession.ChannelMessageSendComplex(gi.game.ChannelId, gi.GetMatchEndResponse())
	if err != nil {
		log.Fatal(err)
	}
}

func OnMessageReaction(s *discordgo.Session, mr *discordgo.MessageReactionAdd) {
	g := game.GetGame(mr.ChannelID)

	if g == nil {
		return
	}

	if mr.MessageID == g.CurrentRound.Message.ID {
		if g.HasUser(mr.Member.User.ID) {
			user := g.GetUser(mr.Member.User.ID)
			if g.GetUser(mr.Member.User.ID).CurrentRoundAnswer != "" {
				sym, _ := countries.GetCountryCodeSymbol(user.CurrentRoundAnswer)
				err := s.MessageReactionRemove(mr.ChannelID, g.CurrentRound.Message.ID, sym, mr.Member.User.ID)
				if err != nil {
					fmt.Println(err)
				}
			}
			g.SetUserAnswer(mr.Member.User.ID, mr.Emoji.Name)
		}
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
		err := s.ChannelMessageDelete(g.ChannelId, g.CurrentRound.ResultMessage.ID)
		if err != nil {
			fmt.Println(err)
		}
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
			Content: phrases.GameExistsResponseContent,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}
}

func GetGameInvitationResponse(channelId string) *discordgo.InteractionResponseData {
	g := game.GetGame(channelId)

	description := phrases.GameInvitationDescription
	buttons := []discordgo.MessageComponent{
		discordgo.Button{
			Emoji: discordgo.ComponentEmoji{
				Name: "➕",
			},
			Label:    phrases.GameInvitationJoinLabel,
			Style:    discordgo.SecondaryButton,
			CustomID: "invitation_join",
		},
	}

	if len(g.Users) > 1 {
		buttons = append(buttons, discordgo.Button{
			Emoji: discordgo.ComponentEmoji{
				Name: "▶",
			},
			Label:    phrases.GameInvitationStartLabel,
			Style:    discordgo.SecondaryButton,
			CustomID: "invitation_start",
		})
	}

	if len(g.Users) > 0 {
		description += "\n" + phrases.GameInvitationJoinedPlayers + ":\n"

		for _, v := range g.Users {
			description += "@" + v.Profile.Username + "\n"
		}
	}

	return &discordgo.InteractionResponseData{
		Content: phrases.GameInvitationStarting,
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       phrases.GameInvitationTitle,
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
		symbol, _ := countries.GetCountryCodeSymbol(c)
		name := countries.GetCountryName(c)
		description += fmt.Sprintf("\n %s - %s \n", symbol, name)
	}

	return &discordgo.MessageSend{
		Content: phrases.CurrentRoundMessageContent,
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       phrases.CurrentRoundMessageContent,
				Description: description,
				Image:       &discordgo.MessageEmbedImage{URL: gi.game.CurrentRound.ImgUrl},
			},
		},
	}
}

func (gi *GameInteractions) GetRoundEndResponse() *discordgo.MessageSend {
	countryEmoji, _ := countries.GetCountryCodeSymbol(gi.game.CurrentRound.CorrectCountry)
	description := fmt.Sprintf(phrases.RoundEndedText+" %s %s.\n", countryEmoji, countries.GetCountryName(gi.game.CurrentRound.CorrectCountry))

	if len(gi.game.CurrentRound.Winners) > 0 {
		description += "\n" + phrases.RoundEndedWinners + "\n"
		for _, w := range gi.game.CurrentRound.Winners {
			description += fmt.Sprintf("%s "+phrases.RoundEndedUserScore+"\n", w.Profile.Username, w.Score)
		}
	} else {
		description += phrases.RoundEndedNoWinners + "\n"
	}

	if gi.game.CurrentRoundNumber < game.MaxRounds {
		description += "\n" + phrases.RoundEndedRestartText + "\n"
	}

	return &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       phrases.RoundEndedTitle,
				Description: description,
				Image:       &discordgo.MessageEmbedImage{URL: gi.game.CurrentRound.ImgUrl},
			},
		},
	}
}

func (gi *GameInteractions) GetMatchEndResponse() *discordgo.MessageSend {
	description := phrases.PlayerRating + "\n"

	for _, u := range gi.game.GetUsersSortedByScore() {
		description += fmt.Sprintf("%s "+phrases.PlayerPoints+"\n", u.Profile.Username, u.Score)
	}

	return &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       phrases.MatchEndedTitle,
				Description: description,
			},
		},
	}
}

func AddCountryReactions(m *discordgo.Message, countryCodes []string) {
	for _, c := range countryCodes {
		cSymbol, _ := countries.GetCountryCodeSymbol(c)
		err := game.DiscordSession.MessageReactionAdd(m.ChannelID, m.ID, cSymbol)
		if err != nil {
			log.Fatal(err)
		}
	}
}
