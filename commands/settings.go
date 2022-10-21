package commands

import (
	"carbon/framework"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"go.mongodb.org/mongo-driver/bson"
)

var Settings = framework.CommandStruct{
	Data: discord.SlashCommandCreate{
		Name:        "settings",
		Description: "Modify the server settings",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionSubCommand{
				Name:        "mute-role",
				Description: "The mute role for this server.",
				Options: []discord.ApplicationCommandOption{
					discord.ApplicationCommandOptionRole{
						Name:        "role",
						Description: "The mute role you want (Leave blank to create automatically)",
					},
				},
			},
		},
	},
	RequiredPermissions: []int64{framework.Permissions.MANAGE_GUILD},
	Execute: func(Interaction events.ApplicationCommandInteractionCreate) {
		switch *Interaction.SlashCommandInteractionData().SubCommandName {
		case "mute-role":
			{
				role, exists := Interaction.SlashCommandInteractionData().OptRole("role")
				if !exists {
					permissions := discord.Permissions.Remove(discord.PermissionSendMessages)
					newRole, err := Interaction.Client().Rest().CreateRole(*Interaction.GuildID(), discord.RoleCreate{Name: "Muted", Permissions: &permissions})
					if err != nil {
						Interaction.CreateMessage(discord.MessageCreate{
							Content: "I am missing permissions to create the muted role.",
							Flags:   discord.MessageFlagEphemeral,
						})
						return
					}
					role = *newRole
				}
				guildFilter := framework.Guild{ID: Interaction.GuildID().String()}.Data()
				guild := framework.FindDocument("guilds", guildFilter)
				if len(guild) >= 1 {
					framework.UpdateDocument("guilds", guildFilter, bson.D{{Key: "MuteRole", Value: role.ID.String()}})
				} else {
					framework.InsertDocument("guilds", bson.D{{Key: "ID", Value: Interaction.GuildID().String()}, {Key: "MuteRole", Value: role.ID.String()}})
				}
			}
		}
	},
}
