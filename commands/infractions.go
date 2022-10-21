package commands

import (
	"carbon/framework"
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"go.mongodb.org/mongo-driver/bson"
)

var Infractions = framework.CommandStruct{
	Data: discord.SlashCommandCreate{
		Name:        "infractions",
		Description: "Check your or an other user's infractions.",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionUser{
				Name:        "member",
				Description: "The member you want to check the infractions for",
			},
		},
	},
	RequiredPermissions: []int64{},
	Execute: func(Interaction events.ApplicationCommandInteractionCreate) {
		member, exists := Interaction.SlashCommandInteractionData().OptUser("member")
		if !exists {
			member = Interaction.User()
		}
		infractions := framework.FindDocuments("infractions", bson.D{{Key: "UserID", Value: member.ID.String()}})
		var embed discord.Embed
		var row discord.ActionRowComponent
		if len(infractions) < 1 {
			embed = discord.Embed{
				Title:       "Infractions",
				Color:       framework.Config.Colors.Primary,
				Description: fmt.Sprintf("%s has no infractions.", member.Tag()),
			}
			row = discord.ActionRowComponent{discord.ButtonComponent{CustomID: "clearinfractions-disabled", Label: "Clear Infractions", Disabled: true, Style: 4}}
		} else {
			description := fmt.Sprintf("**%s's infractions**\n\n", member.Tag())
			for i, f := range infractions {
				if i > 9 {
					continue
				}
				description += fmt.Sprintf("**%d**. `%d` - %s\n", i+1, f["ID"], f["Type"])
			}
			embed = discord.Embed{
				Title:       "Infractions",
				Description: description,
				Color:       framework.Config.Colors.Primary,
			}
			row = discord.ActionRowComponent{discord.ButtonComponent{CustomID: fmt.Sprintf("clearinfractions-%s", member.ID.String()), Label: "Clear Infractions", Style: 4}}
		}
		Interaction.CreateMessage(discord.MessageCreate{Embeds: []discord.Embed{embed}, Components: []discord.ContainerComponent{row}})
	},
}
