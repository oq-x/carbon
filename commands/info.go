package commands

import (
	"carbon/framework"
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var Info = framework.CommandStruct{
	Data: discord.SlashCommandCreate{
		Name:        "info",
		Description: "Information about Carbon",
		Options:     []discord.ApplicationCommandOption{},
	},
	RequiredPermissions: []int64{},
	Execute: func(Interaction events.ApplicationCommandInteractionCreate) {
		me, _ := Interaction.Client().Caches().GetSelfUser()
		var components []discord.ContainerComponent
		if framework.Config.Public == "true" {
			components = []discord.ContainerComponent{
				discord.ActionRowComponent{discord.ButtonComponent{
					Style: discord.ButtonStyleLink,
					Label: "Add me to your server",
					URL:   fmt.Sprintf("https://discord.com/api/oauth2/authorize?client_id=%s&permissions=1099780063238&scope=bot%%20applications.commands", framework.Config.ID),
				}},
			}
		} else {
			components = []discord.ContainerComponent{}
		}
		Interaction.CreateMessage(discord.MessageCreate{
			Embeds: []discord.Embed{
				{
					Title:       "Carbon",
					Description: "Carbon is a Discord bot made to help you better moderate your server.\n\n[GitHub](https://github.com/oq-x/carbon)",
					Thumbnail:   &discord.EmbedResource{URL: *me.AvatarURL()},
					Color:       framework.Config.Colors.Primary,
				},
			},
			Components: components,
		})
	},
}
