package main

import (
	"carbon/framework"
	"strconv"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/log"
	"go.mongodb.org/mongo-driver/bson"
)

func applicationCommandHandler(i *events.ApplicationCommandInteractionCreate) {
	if i.GuildID() == nil {
		i.CreateMessage(discord.MessageCreate{
			Content: "Commands must be ran in guilds.",
			Flags:   discord.MessageFlagEphemeral,
		})
		return
	}
	name := i.Data.CommandName()
	command, found := CmdHandler.Get(name)
	if !found {
		return
	}
	missingPermissions := []int64{}
	if len(command.RequiredPermissions) > 0 {
		for _, x := range command.RequiredPermissions {
			if !i.Member().Permissions.Has(discord.Permissions(x)) {
				missingPermissions = append(missingPermissions, x)
			}
		}
	}
	if len(missingPermissions) > 0 {
		i.CreateMessage(discord.MessageCreate{
			Content: "You aren't permitted to run this command.",
			Flags:   discord.MessageFlagEphemeral,
		})
		return
	}
	command.Execute(*i)
	log.Debugf("Command %s has been ran", name)

}

func messageComponentHandler(i *events.ComponentInteractionCreate) {
	if i.GuildID() == nil {
		i.CreateMessage(discord.MessageCreate{
			Content: "Commands must be ran in guilds.",
			Flags:   discord.MessageFlagEphemeral,
		})
		return
	}
	switch i.ButtonInteractionData().Type() {
	case 2:
		{
			if strings.Split(i.ButtonInteractionData().CustomID(), "-")[0] == "clearinfractions" {
				if !i.Member().Permissions.Has(discord.PermissionModerateMembers) {
					i.CreateMessage(discord.MessageCreate{Content: "You are not permitted to use this component.", Flags: discord.MessageFlagEphemeral})
					return
				}
				framework.DeleteDocuments("infractions", bson.D{{Key: "UserID", Value: strings.Split(i.ButtonInteractionData().CustomID(), "-")[1]}})

				i.CreateMessage(discord.MessageCreate{Content: "Successfully cleared infractions."})
			} else if strings.Split(i.ButtonInteractionData().CustomID(), "-")[0] == "deleteinfraction" {
				if !i.Member().Permissions.Has(discord.PermissionModerateMembers) {
					i.CreateMessage(discord.MessageCreate{Content: "You are not permitted to use this component.", Flags: discord.MessageFlagEphemeral})
					return
				}
				id, _ := strconv.Atoi(strings.Split(i.ButtonInteractionData().CustomID(), "-")[1])
				framework.DeleteDocuments("infractions", bson.D{{Key: "ID", Value: id}})
				row := discord.ActionRowComponent{discord.ButtonComponent{CustomID: "deleteinfraction-disabled", Label: "Delete Infraction", Style: 4, Disabled: true}}
				i.UpdateMessage(discord.MessageUpdate{Components: &[]discord.ContainerComponent{row}})
				i.Client().Rest().CreateFollowupMessage(i.Client().ID(), i.Token(), discord.MessageCreate{Content: "Successfully deleted infraction."})
			}

			break
		}
	}
}
