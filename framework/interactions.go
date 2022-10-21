package framework

import "github.com/bwmarrin/discordgo"

type Interaction struct {
	Session *discordgo.Session
	*discordgo.InteractionCreate
}

func (interaction Interaction) GetUser() *discordgo.User {
	if interaction.User == nil {
		return interaction.Member.User
	} else {
		return interaction.User
	}
}

func (interaction Interaction) GetOption(name string) *discordgo.ApplicationCommandInteractionDataOption {
	options := interaction.ApplicationCommandData().Options
	var option *discordgo.ApplicationCommandInteractionDataOption
	if options[0].Type == 1 {
		options = options[0].Options
	} else if options[0].Type == 2 {
		options = options[0].Options[0].Options
	} else {
		for _, i := range options {
			if i.Name == name {
				option = i
			}
		}
	}
	return option
}
func (interaction Interaction) GetSubcommand() string {
	var option string
	if interaction.ApplicationCommandData().Options[0].Type == 1 {
		option = interaction.ApplicationCommandData().Options[0].Name
	} else if interaction.ApplicationCommandData().Options[0].Type == 2 {
		if interaction.ApplicationCommandData().Options[0].Options[0].Type == 1 {
			option = interaction.ApplicationCommandData().Options[0].Options[0].Name
		}
	}
	return option
}
func (interaction Interaction) GetSubcommandGroup() string {
	var option string
	if interaction.ApplicationCommandData().Options[0].Type == 2 {
		option = interaction.ApplicationCommandData().Options[0].Name
	}
	return option
}
func (interaction Interaction) Reply(data *discordgo.InteractionResponseData) error {
	return interaction.Session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: 4,
		Data: data,
	})
}
