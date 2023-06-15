package commands

import "github.com/bwmarrin/discordgo"

var (
	answerHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
)

func AddAnswerHandler(s string, f func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	answerHandlers[s] = f
}
