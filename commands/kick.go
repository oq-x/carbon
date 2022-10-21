package commands

import (
	"carbon/framework"
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/rest"
	"github.com/google/uuid"
)

var Kick = framework.CommandStruct{
	Data: discord.SlashCommandCreate{
		Name:        "kick",
		Description: "Kick a member",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionUser{
				Name:        "member",
				Description: "The member you want to kick",
				Required:    true,
			},
			discord.ApplicationCommandOptionString{
				Name:        "reason",
				Description: "The reason why you want to kick this member",
			},
		},
	},
	RequiredPermissions: []int64{framework.Permissions.KICK_MEMBERS},
	Execute: func(Interaction events.ApplicationCommandInteractionCreate) {
		user, exists := Interaction.SlashCommandInteractionData().OptUser("member")
		guild, _ := Interaction.Guild()
		member, _ := Interaction.Client().Caches().Members().Get(guild.ID, user.ID)
		me, _ := Interaction.Client().Caches().GetSelfMember(*Interaction.GuildID())
		reason, reasonExists := Interaction.SlashCommandInteractionData().OptString("reason")
		if !reasonExists {
			reason = "None"
		}
		if !exists {
			Interaction.CreateMessage(discord.MessageCreate{
				Content: "Please choose a member that is in this server!",
				Flags:   discord.MessageFlagEphemeral,
			})
			return
		}
		if framework.HigherMember(Interaction.Client(), guild, &Interaction.Member().Member, &member) != Interaction.User().ID.String() {
			Interaction.CreateMessage(discord.MessageCreate{
				Content: fmt.Sprintf("You cannot kick **%s#%s** because they have a higher role than you.", member.User.Username, member.User.Discriminator),
				Flags:   discord.MessageFlagEphemeral,
			})
			return
		}
		if framework.HigherMember(Interaction.Client(), guild, &me, &member) != me.User.ID.String() {
			Interaction.CreateMessage(discord.MessageCreate{
				Content: fmt.Sprintf("I cannot kick **%s#%s** because they have a higher role than me.", member.User.Username, member.User.Discriminator),
				Flags:   64,
			})
			return
		}
		var res string
		banError := Interaction.Client().Rest().RemoveMember(*Interaction.GuildID(), member.User.ID, rest.WithReason(fmt.Sprintf("Kicked by %s. Reason: %s", Interaction.User().Tag(), reason)))
		if banError != nil {
			res = fmt.Sprintf("❌ I was not able to kick **%s#%s**.", member.User.Username, member.User.Discriminator)
		} else {
			res = fmt.Sprintf("☑️ I was able to kick **%s#%s**.\n*Reason*: `%s`", member.User.Username, member.User.Discriminator, reason)
		}
		ID := uuid.New().ID()
		framework.InsertDocument("infractions", framework.Infraction{
			ID:          int(ID),
			Type:        "Kick",
			UserID:      member.User.ID.String(),
			ModeratorID: Interaction.User().ID.String(),
			Reason:      reason,
			GuildID:     Interaction.GuildID().String(),
		}.Data())
		Interaction.CreateMessage(discord.MessageCreate{
			Content: res,
		})
	},
}
