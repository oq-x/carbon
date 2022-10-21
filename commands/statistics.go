package commands

import (
	"carbon/framework"
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/rest"
)

var Statistics = framework.CommandStruct{
	Data: discord.SlashCommandCreate{
		Name:        "statistics",
		Description: "Carbon's statistics",
		Options:     []discord.ApplicationCommandOption{},
	},
	RequiredPermissions: []int64{},
	Execute: func(Interaction events.ApplicationCommandInteractionCreate) {
		ping := Interaction.Client().Gateway().Latency()
		True := true
		me, _ := Interaction.Client().Caches().GetSelfUser()
		membercount := 0
		for _, i := range Interaction.Client().Caches().Guilds().All() {
			membercount += i.MemberCount
		}
		Interaction.CreateMessage(discord.MessageCreate{
			Embeds: []discord.Embed{
				{
					Title: "Statistics",
					Fields: []discord.EmbedField{
						{
							Name:   "Guilds",
							Value:  fmt.Sprintf("%d", Interaction.Client().Caches().Guilds().Len()),
							Inline: &True,
						},
						{
							Name:   "Total Members",
							Value:  strconv.Itoa(membercount),
							Inline: &True,
						},
						{
							Name:  "Go Version",
							Value: runtime.Version(),
						},
						{
							Name:   "Disgo Version",
							Value:  disgo.Version,
							Inline: &True,
						},
						{
							Name:   "Gateway Version",
							Value:  fmt.Sprintf("v%d", gateway.Version),
							Inline: &True,
						},
						{
							Name:   "API Version",
							Value:  fmt.Sprintf("v%d", rest.APIVersion),
							Inline: &True,
						},
						{
							Name:   "Websocket Ping",
							Value:  fmt.Sprintf("%s", ping.Round(time.Millisecond)),
							Inline: &True,
						},
					},
					Thumbnail: &discord.EmbedResource{URL: *me.AvatarURL()},
					Color:     framework.Config.Colors.Primary,
				},
			},
		})
	},
}
