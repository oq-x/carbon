package commands

import (
	"carbon/framework"
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var RemoveTimeout = framework.CommandStruct{
	Data: discord.SlashCommandCreate{
		Name:        "remove-timeout",
		Description: "Remove a member's timeout",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionUser{
				Name:        "member",
				Description: "The member you want to untime-out",
				Required:    true,
			},
		},
	},
	RequiredPermissions: []int64{framework.Permissions.MODERATE_MEMBERS},
	Execute: func(Interaction events.ApplicationCommandInteractionCreate) {
		user, _ := Interaction.SlashCommandInteractionData().OptUser("member")
		me, _ := Interaction.Client().Caches().GetSelfMember(*Interaction.GuildID())
		guild, _ := Interaction.Client().Caches().Guilds().Get(*Interaction.GuildID())
		member, _ := Interaction.Client().Caches().Members().Get(guild.ID, user.ID)
		if framework.HigherMember(Interaction.Client(), guild, &Interaction.Member().Member, &member) != Interaction.Member().User.ID.String() {
			Interaction.CreateMessage(discord.MessageCreate{
				Content: fmt.Sprintf("You cannot untime-out **%s** because they have a higher role than you.", member.User.Tag()),
				Flags:   discord.MessageFlagEphemeral,
			})
			return
		}

		if Interaction.Member().Permissions.Has(discord.PermissionAdministrator) {
			Interaction.CreateMessage(discord.MessageCreate{
				Content: fmt.Sprintf("I cannot untime-out **%s** because they are an administrator.", member.User.Tag()),
				Flags:   discord.MessageFlagEphemeral,
			})
			return
		}

		if framework.HigherMember(Interaction.Client(), guild, &me, &member) != me.User.ID.String() {
			Interaction.CreateMessage(discord.MessageCreate{
				Content: fmt.Sprintf("I cannot untime-out **%s** because they have a higher role than me.", member.User.Tag()),
				Flags:   discord.MessageFlagEphemeral,
			})
			return
		}
		var res string
		_, TimeoutError := Interaction.Client().Rest().UpdateMember(*Interaction.GuildID(), member.User.ID, discord.MemberUpdate{CommunicationDisabledUntil: nil})
		if TimeoutError != nil {
			res = fmt.Sprintf("❌ I was not able to untime-out **%s**.", member.User.Tag())
		} else {
			res = fmt.Sprintf("☑️ I was able to time-out **%s**.", member.User.Tag())
		}

		Interaction.CreateMessage(discord.MessageCreate{
			Content: res,
		})
	},
}
