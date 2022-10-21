package commands

import (
	"carbon/framework"
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

var Unmute = framework.CommandStruct{
	Data: discord.SlashCommandCreate{
		Name:        "unmute",
		Description: "Unmute a member.",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionUser{
				Name:        "member",
				Description: "The member you want to unmute",
				Required:    true,
			},
		},
	},
	RequiredPermissions: []int64{framework.Permissions.MODERATE_MEMBERS},
	Execute: func(Interaction events.ApplicationCommandInteractionCreate) {
		user, _ := Interaction.SlashCommandInteractionData().OptUser("member")
		reason, reasonExists := Interaction.SlashCommandInteractionData().OptString("reason")
		guild, _ := Interaction.Guild()
		member, _ := Interaction.Client().Caches().Members().Get(guild.ID, user.ID)
		guildData := framework.FindDocument("guilds", framework.Guild{ID: guild.ID.String()}.Data())
		if len(guildData) < 1 || guildData["MuteRole"] == nil {
			Interaction.CreateMessage(discord.MessageCreate{
				Content: "This server does not have a mute role set. Please use `/settings mute-role`.",
				Flags:   discord.MessageFlagEphemeral,
			})
			return
		}
		id, _ := snowflake.Parse(guildData["MuteRole"].(string))
		role, exists := Interaction.Client().Caches().Roles().Get(guild.ID, id)
		if !exists {
			Interaction.CreateMessage(discord.MessageCreate{
				Content: "This server's mute role has been deleted. Please use `/settings mute-role`.",
				Flags:   discord.MessageFlagEphemeral,
			})
			return
		}
		if framework.HigherMember(Interaction.Client(), guild, &Interaction.Member().Member, &member) != Interaction.User().ID.String() {
			Interaction.CreateMessage(discord.MessageCreate{
				Content: fmt.Sprintf("You cannot unmute **%s** because they have a higher role than you.", member.User.Tag()),
				Flags:   discord.MessageFlagEphemeral,
			})
			return
		}
		if !reasonExists {
			reason = "None"
		}
		var res string
		UnmuteError := Interaction.Client().Rest().RemoveMemberRole(*Interaction.GuildID(), member.User.ID, role.ID)
		if UnmuteError != nil {
			res = fmt.Sprintf("❌ I was not able to unmute **%s**.", member.User.Tag())
		} else {
			res = fmt.Sprintf("☑️ I was able to mute **%s**. *Reason*: `%s`", member.User.Tag(), reason)
		}

		Interaction.CreateMessage(discord.MessageCreate{
			Content: res,
		})
	},
}
