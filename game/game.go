package game

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"mapguess-discord/api"
	"mapguess-discord/utils/countries"
	"mapguess-discord/utils/phrases"
	"math/rand"
	"sort"
	"time"
)

const (
	MaxRounds = 10
)

type Round struct {
	ImgUrl         string
	CorrectCountry string
	CountryOptions []string
	Winners        []*User
	Message        *discordgo.Message
	ResultMessage  *discordgo.Message
	Ended          bool
}

type User struct {
	Profile            *discordgo.User
	Score              int
	Interaction        *discordgo.Interaction
	CurrentRoundAnswer string
}

type Game struct {
	MatchStarted       bool
	ChannelId          string
	Users              map[string]*User // user id - user
	CurrentRound       *Round
	CurrentRoundNumber int
	eventListeners     []EventListener
}

type EventListener interface {
	OnRoundStart()
	OnRoundEnd()
	OnMatchEnd()
}

var (
	DiscordSession *discordgo.Session
	games          = map[string]*Game{} // channel id - game
)

func (g *Game) RegisterGameListener(l EventListener) {
	g.eventListeners = append(g.eventListeners, l)
}

func (r *Round) SetMessage(m *discordgo.Message) {
	r.Message = m
}

func StartGame(channelId string) error {
	if _, exists := games[channelId]; exists == true {
		return errors.New(phrases.GameInvitationErrExists)
	}

	games[channelId] = &Game{
		ChannelId: channelId,
		Users:     map[string]*User{},
	}
	return nil
}

func JoinGame(channelId string, i *discordgo.Interaction) error {
	if _, ok := games[channelId]; !ok {
		return errors.New(phrases.GameInvitationErrNotFound)
	}

	user := i.Member.User

	_, ok := games[channelId].Users[user.ID]
	if ok {
		return errors.New(phrases.GameInvitationErrUserExists)
	}

	games[channelId].Users[user.ID] = &User{
		Profile:     user,
		Score:       0,
		Interaction: i,
	}
	return nil
}

func (g *Game) StartMatch() {
	g.MatchStarted = true
	g.StartRound()
}

func (g *Game) StartRound() {
	g.CurrentRoundNumber++
	photos := api.GetRandomPhoto()

	rand.Seed(time.Now().Unix())
	photo := photos.Result.Data[rand.Intn(len(photos.Result.Data))]

	c := countries.GetRandomCountriesExcept(photos.CountryCode)[:3:4]
	c = append(c, photos.CountryCode)

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(c), func(i, j int) {
		c[i], c[j] = c[j], c[i]
	})

	g.CurrentRound = &Round{
		ImgUrl:         photo.FileUrl,
		CorrectCountry: photos.CountryCode,
		CountryOptions: c,
	}

	for _, gl := range g.eventListeners {
		gl.OnRoundStart()
	}
}

func (g *Game) endRound() {
	g.CurrentRound.Ended = true
	for _, u := range g.Users {
		if u.CurrentRoundAnswer == g.CurrentRound.CorrectCountry {
			u.Score += 1
			g.CurrentRound.Winners = append(g.CurrentRound.Winners, u)
		}
		u.CurrentRoundAnswer = ""
	}

	for _, gl := range g.eventListeners {
		gl.OnRoundEnd()
	}

	if g.CurrentRoundNumber >= MaxRounds {
		g.endMatch()
	}
}

func (g *Game) endMatch() {
	g.MatchStarted = false

	for _, gl := range g.eventListeners {
		gl.OnMatchEnd()
	}
	delete(games, g.ChannelId)
}

func (g *Game) SetUserAnswer(userId string, answer string) {
	u := g.GetUser(userId)
	countryCode := countries.GetStringFromCountrySymbol(answer)

	if u == nil || !countries.HasCountryCode(countryCode) || g.CurrentRound.Ended {
		return
	}

	u.CurrentRoundAnswer = countryCode

	g.CheckRoundEnd()
}

func (g *Game) HasUser(userId string) bool {
	for _, u := range g.Users {
		if u.Profile.ID == userId {
			return true
		}
	}
	return false
}

func (g *Game) GetUser(userId string) *User {
	for _, u := range g.Users {
		if u.Profile.ID == userId {
			return u
		}
	}
	return nil
}

func (g *Game) CheckRoundEnd() {
	for _, u := range g.Users {
		if u.CurrentRoundAnswer == "" {
			return
		}
	}

	g.endRound()
}

func GetGame(channelId string) *Game {
	return games[channelId]
}

func SetSession(s *discordgo.Session) {
	DiscordSession = s
}

func (g *Game) GetUsersSortedByScore() []*User {
	var users []*User

	for _, u := range g.Users {
		users = append(users, u)
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].Score > users[j].Score
	})

	return users
}
