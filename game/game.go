package game

import (
	"errors"
	"github.com/bwmarrin/discordgo"
)

type User struct {
	profile *discordgo.User
	score   int
}

type Game struct {
	channelId string
	users     map[string]*User // user id - user
}

var (
	games = map[string]*Game{} // channel id - game
)

func StartGame(channelId string) error {
	if _, exists := games[channelId]; exists == true {
		return errors.New("Game already exists")
	}

	games[channelId] = &Game{
		channelId: channelId,
		users:     map[string]*User{},
	}
	return nil
}

func JoinGame(channelId string, user *discordgo.User) {
	games[channelId].users[user.ID] = &User{
		profile: user,
		score:   0,
	}
}
