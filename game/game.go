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
}

type User struct {
	Profile     *discordgo.User
	Score       int
	Interaction *discordgo.Interaction
}

type Game struct {
	ChannelId     string
	Users         map[string]*User // user id - user
	Round         *Round
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

func JoinGame(channelId string, i *discordgo.Interaction) {
	user := i.Member.User

	_, ok := games[channelId].Users[user.ID]
	if ok {
		return // User has already joined
	}

	games[channelId].Users[user.ID] = &User{
		Profile:     user,
		Score:       0,
		Interaction: i,
	}
}

func (game *Game) StartMatch() {
	game.startRound()
}

func (game *Game) startRound() {
	photos := api.GetRandomPhoto()

	rand.Seed(time.Now().Unix())
	photo := photos.Result.Data[rand.Intn(len(photos.Result.Data))]

	countries := utils.GetRandomCountriesExcept(photos.CountryCode)[:3:4]
	countries = append(countries, photos.CountryCode)

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(countries), func(i, j int) {
		countries[i], countries[j] = countries[j], countries[i]
	})

	game.Round = &Round{
		ImgUrl:         photo.FileUrl,
		CorrectCountry: photos.CountryCode,
		CountryOptions: countries,
	}

	for _, gl := range game.gameListeners {
		gl.OnRoundStart()
	}
}

func GetGame(channelId string) *Game {
	return games[channelId]
}

func SetSession(s *discordgo.Session) {
	DiscordSession = s
}
