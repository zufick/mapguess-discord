package game

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"mapguess-discord/api"
	"mapguess-discord/utils"
	"math/rand"
	"time"
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
	MatchStarted  bool
	ChannelId     string
	Users         map[string]*User // user id - user
	CurrentRound  *Round
	gameListeners []GameListener
}

type GameListener interface {
	OnRoundStart()
	OnRoundEnd()
	OnGameEnd()
}

var (
	DiscordSession *discordgo.Session
	games          = map[string]*Game{} // channel id - game
)

func (g *Game) RegisterGameListener(l GameListener) {
	g.gameListeners = append(g.gameListeners, l)
}

func (r *Round) SetMessage(m *discordgo.Message) {
	r.Message = m
}

func StartGame(channelId string) error {
	if _, exists := games[channelId]; exists == true {
		return errors.New("Game already exists")
	}

	games[channelId] = &Game{
		ChannelId: channelId,
		Users:     map[string]*User{},
	}
	return nil
}

func JoinGame(channelId string, i *discordgo.Interaction) error {
	if _, ok := games[channelId]; !ok {
		return errors.New("Error while joining. Cannot find this game.")
	}

	user := i.Member.User

	_, ok := games[channelId].Users[user.ID]
	if ok {
		return errors.New("User has already joined.")
	}

	games[channelId].Users[user.ID] = &User{
		Profile:     user,
		Score:       0,
		Interaction: i,
	}
	return nil
}

func (game *Game) StartMatch() {
	game.MatchStarted = true
	game.StartRound()
}

func (game *Game) StartRound() {
	photos := api.GetRandomPhoto()

	rand.Seed(time.Now().Unix())
	photo := photos.Result.Data[rand.Intn(len(photos.Result.Data))]

	countries := utils.GetRandomCountriesExcept(photos.CountryCode)[:3:4]
	countries = append(countries, photos.CountryCode)

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(countries), func(i, j int) {
		countries[i], countries[j] = countries[j], countries[i]
	})

	game.CurrentRound = &Round{
		ImgUrl:         photo.FileUrl,
		CorrectCountry: photos.CountryCode,
		CountryOptions: countries,
	}

	for _, gl := range game.gameListeners {
		gl.OnRoundStart()
	}
}

func (game *Game) endRound() {
	game.CurrentRound.Ended = true
	for _, u := range game.Users {
		if u.CurrentRoundAnswer == game.CurrentRound.CorrectCountry {
			u.Score += 1
			game.CurrentRound.Winners = append(game.CurrentRound.Winners, u)
		}
		u.CurrentRoundAnswer = ""
	}

	for _, gl := range game.gameListeners {
		gl.OnRoundEnd()
	}
}

func (game *Game) SetUserAnswer(userId string, answer string) {
	u := game.GetUser(userId)
	countryCode := utils.GetStringFromCountrySymbol(answer)

	if u == nil || !utils.HasCountryCode(countryCode) || game.CurrentRound.Ended {
		return
	}

	u.CurrentRoundAnswer = countryCode

	game.CheckRoundEnd()
}

func (game *Game) HasUser(userId string) bool {
	for _, u := range game.Users {
		if u.Profile.ID == userId {
			return true
		}
	}
	return false
}

func (game *Game) GetUser(userId string) *User {
	for _, u := range game.Users {
		if u.Profile.ID == userId {
			return u
		}
	}
	return nil
}

func (game *Game) CheckRoundEnd() {
	for _, u := range game.Users {
		if u.CurrentRoundAnswer == "" {
			return
		}
	}

	game.endRound()
}

func GetGame(channelId string) *Game {
	return games[channelId]
}

func SetSession(s *discordgo.Session) {
	DiscordSession = s
}
