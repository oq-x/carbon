package commands

import (
	"carbon/framework"
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/json"
	"github.com/disgoorg/disgo/rest"
	"github.com/google/uuid"
)

var Timeout = framework.CommandStruct{
	Data: discord.SlashCommandCreate{
		Name:        "timeout",
		Description: "Time-out a member.",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionUser{
				Name:        "member",
				Description: "The member you want to time-out",
				Required:    true,
			},
			discord.ApplicationCommandOptionString{
				Name:        "duration",
				Description: "The duration",
				Required:    true,
			},
			discord.ApplicationCommandOptionString{
				Name:        "reason",
				Description: "The reason why you want to time-out this member",
			},
		},
	},
	RequiredPermissions: []int64{framework.Permissions.MODERATE_MEMBERS},
	Execute: func(Interaction events.ApplicationCommandInteractionCreate) {
		user, _ := Interaction.SlashCommandInteractionData().OptUser("member")
		me, _ := Interaction.Client().Caches().GetSelfMember(*Interaction.GuildID())
		guild, _ := Interaction.Client().Caches().Guilds().Get(*Interaction.GuildID())
		member, _ := Interaction.Client().Caches().Members().Get(guild.ID, user.ID)
		rawDuration, _ := Interaction.SlashCommandInteractionData().OptString("duration")
		reason, reasonExists := Interaction.SlashCommandInteractionData().OptString("reason")
		if !reasonExists {
			reason = "None"
		}
		duration, durationError := time.ParseDuration(rawDuration)
		if durationError != nil {
			Interaction.CreateMessage(discord.MessageCreate{
				Content: fmt.Sprintf("I wasn't able to parse **%s**", rawDuration),
				Flags:   discord.MessageFlagEphemeral,
			})
			return
		}
		if framework.HigherMember(Interaction.Client(), guild, &Interaction.Member().Member, &member) != Interaction.Member().User.ID.String() {
			Interaction.CreateMessage(discord.MessageCreate{
				Content: fmt.Sprintf("You cannot time-out **%s** because they have a higher role than you.", member.User.Tag()),
				Flags:   discord.MessageFlagEphemeral,
			})
			return
		}

		if Interaction.Member().Permissions.Has(discord.PermissionAdministrator) {
			Interaction.CreateMessage(discord.MessageCreate{
				Content: fmt.Sprintf("I cannot time-out **%s** because they are an administrator.", member.User.Tag()),
				Flags:   discord.MessageFlagEphemeral,
			})
			return
		}

		if framework.HigherMember(Interaction.Client(), guild, &me, &member) != me.User.ID.String() {
			Interaction.CreateMessage(discord.MessageCreate{
				Content: fmt.Sprintf("I cannot time-out **%s** because they have a higher role than me.", member.User.Tag()),
				Flags:   discord.MessageFlagEphemeral,
			})
			return
		}
		var res string
		Time := json.New(time.Now().Add(duration))
		_, TimeoutError := Interaction.Client().Rest().UpdateMember(*Interaction.GuildID(), member.User.ID, discord.MemberUpdate{CommunicationDisabledUntil: &Time}, rest.WithReason(fmt.Sprintf("Timed out by %s. Reason: %s", Interaction.User().Tag(), reason)))
		if TimeoutError != nil {
			res = fmt.Sprintf("❌ I was not able to time-out **%s**.", member.User.Tag())
		} else {
			res = fmt.Sprintf("☑️ I was able to time-out **%s**. *Reason*: `%s`", member.User.Tag(), reason)
		}
		ID := uuid.New().ID()
		framework.InsertDocument("infractions", framework.Infraction{
			ID:          int(ID),
			Type:        "Timeout",
			UserID:      me.User.ID.String(),
			ModeratorID: Interaction.User().ID.String(),
			GuildID:     Interaction.GuildID().String(),
			Reason:      reason,
		}.Data())

		Interaction.CreateMessage(discord.MessageCreate{
			Content: res,
		})
	},
}
