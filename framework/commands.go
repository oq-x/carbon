package framework

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type (
	execute       func(events.ApplicationCommandInteractionCreate)
	CommandStruct struct {
		Execute             execute
		Data                discord.SlashCommandCreate
		RequiredPermissions []int64
	}
	CmdMap map[string]CommandStruct

	CommandHandler struct {
		Cmds CmdMap
	}
)

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{make(CmdMap)}
}

func (handler CommandHandler) GetCmds() CmdMap {
	return handler.Cmds
}

func (handler CommandHandler) Get(name string) (*CommandStruct, bool) {
	cmd, found := handler.Cmds[name]
	return &cmd, found
}

func (handler CommandHandler) Register(commands []CommandStruct) {
	for _, command := range commands {
		handler.Cmds[command.Data.Name] = command
	}
}
