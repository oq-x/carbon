package commands

import (
	"carbon/framework"
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/google/uuid"
)

var Warn = framework.CommandStruct{
	Data: discord.SlashCommandCreate{
		Name:        "warn",
		Description: "Warn a member.",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionUser{
				Name:        "member",
				Description: "The member you want to warn",
				Required:    true,
			},
			discord.ApplicationCommandOptionString{
				Name:        "reason",
				Description: "The reason why you want to warn this member",
			},
		},
	},
	RequiredPermissions: []int64{framework.Permissions.MODERATE_MEMBERS},
	Execute: func(Interaction events.ApplicationCommandInteractionCreate) {
		user, _ := Interaction.SlashCommandInteractionData().OptUser("member")
		reason, reasonExists := Interaction.SlashCommandInteractionData().OptString("reason")
		guild, _ := Interaction.Guild()
		member, _ := Interaction.Client().Caches().Members().Get(guild.ID, user.ID)
		if framework.HigherMember(Interaction.Client(), guild, &Interaction.Member().Member, &member) != Interaction.User().ID.String() {
			Interaction.CreateMessage(discord.MessageCreate{
				Content: fmt.Sprintf("You cannot warn **%s** because they have a higher role than you.", member.User.Tag()),
				Flags:   discord.MessageFlagEphemeral,
			})
			return
		}
		if !reasonExists {
			reason = "None"
		}
		ID := uuid.New().ID()
		framework.InsertDocument("infractions", framework.Infraction{
			ID:          int(ID),
			Type:        "Warning",
			UserID:      user.ID.String(),
			ModeratorID: Interaction.User().ID.String(),
			GuildID:     Interaction.GuildID().String(),
			Reason:      reason,
		}.Data())

		Interaction.CreateMessage(discord.MessageCreate{
			Content: fmt.Sprintf("☑️ Successfully added infraction **%d** to user **%s**.\n*Reason*: `%s`", ID, user.Tag(), reason),
		})
	},
}
