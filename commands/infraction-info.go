package commands

import (
	"carbon/framework"
	"fmt"
	"strconv"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"go.mongodb.org/mongo-driver/bson"
)

var InfractionInfo = framework.CommandStruct{
	Data: discord.SlashCommandCreate{
		Name:        "infraction-info",
		Description: "Get info of an infraction",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionInt{
				Name:        "id",
				Description: "The ID of the infraction",
				Required:    true,
			},
		},
	},
	RequiredPermissions: []int64{framework.Permissions.MODERATE_MEMBERS},
	Execute: func(Interaction events.ApplicationCommandInteractionCreate) {
		id, _ := Interaction.SlashCommandInteractionData().OptInt("id")
		infraction := framework.FindDocument("infractions", bson.D{{Key: "ID", Value: id}})
		if len(infraction) < 1 {
			Interaction.CreateMessage(discord.MessageCreate{Content: "Unknown infraction.", Flags: discord.MessageFlagEphemeral})
			return
		}
		embed := discord.Embed{
			Title: strconv.Itoa(id),
			Description: fmt.Sprintf("Type: %s\nModerator: %s\nUser: %s\nReason: %s",
				infraction["Type"],
				"<@"+infraction["ModeratorID"].(string)+">",
				"<@"+infraction["UserID"].(string)+">",
				infraction["Reason"]),
			Color: framework.Config.Colors.Primary,
		}
		row := discord.ActionRowComponent{discord.ButtonComponent{CustomID: fmt.Sprintf("deleteinfraction-%d", infraction["ID"]), Label: "Delete Infraction", Style: 4}}
		Interaction.CreateMessage(discord.MessageCreate{Embeds: []discord.Embed{embed}, Components: []discord.ContainerComponent{row}})
	},
}
