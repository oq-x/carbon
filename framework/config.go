package framework

import (
	"os"

	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
	"github.com/joho/godotenv"
)

type colors struct {
	Primary int
	Success int
	Danger  int
	Warning int
}

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load environment variables.")
	} else {
		log.Info("Loaded environment variables.")
		var id, _ = snowflake.Parse(os.Getenv("ID"))
		Config = config{
			Token:    os.Getenv("TOKEN"),
			ID:       id,
			MongoURI: os.Getenv("MONGO_URI"),
			Public:   os.Getenv("PUBLIC"),
			Colors: colors{
				Primary: 0x006CFF,
				Success: 0x00AB0D,
				Danger:  0xC70000,
				Warning: 0xFFC900,
			},
		}

	}
}

type config struct {
	Token    string
	ID       snowflake.ID
	MongoURI string
	Colors   colors
	Public   string
}

var Config config
