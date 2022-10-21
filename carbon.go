package main

import (
	"carbon/commands"
	"carbon/framework"
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/log"
)

var (
	CmdHandler *framework.CommandHandler
)

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func main() {
	framework.LoadConfig()
	client, err := disgo.New(framework.Config.Token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				7,
			),
		),
		bot.WithCacheConfigOpts(
			cache.WithCacheFlags(cache.FlagGuilds),
			cache.WithCacheFlags(cache.FlagRoles),
		),
		bot.WithEventListenerFunc(applicationCommandHandler),
		bot.WithEventListenerFunc(messageComponentHandler),
	)

	log.SetLevel(log.LevelDebug)
	if err = client.OpenGateway(context.TODO()); err != nil {
		panic(err)
	}

	CmdHandler = framework.NewCommandHandler()
	log.Debug("Loaded command handler")
	registerCommands(client)
	framework.MongoConnect(framework.Config.MongoURI)

	log.Info("Carbon has started")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
func registerCommands(client bot.Client) {
	CmdHandler.Register([]framework.CommandStruct{
		commands.Ban,
		commands.Info,
		commands.Statistics,
		commands.Timeout,
		commands.Warn,
		commands.Infractions,
		commands.InfractionInfo,
		commands.RemoveTimeout,
		commands.Settings,
		commands.Mute,
		commands.Unmute,
	})
	keys := []discord.ApplicationCommandCreate{}
	for key := range CmdHandler.Cmds {
		cmd := CmdHandler.Cmds[key]
		log.Debug("Loaded command ", cmd.Data.Name)
		keys = append(keys, cmd.Data)
	}
	_, err := client.Rest().SetGlobalCommands(framework.Config.ID, keys)
	if err == nil {
		log.Info("Registered application commands.")
	} else {
		log.Errorf("Failed to register application commands: %s", err)
	}
}
