package discordclient

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func InitConnection() (error, *discordgo.Session) {
	token := os.Getenv("REGINALDTOKEN")
	discord, err := discordgo.New("Bot " + token)

	if err != nil {
		fmt.Println("LMAO YOU GOON")
		return err, nil
	}

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return err, nil
	}
	return nil, discord
}

func CloseConnection(discord *discordgo.Session) {
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}
