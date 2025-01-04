package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"wynn_bot/models"
	"wynn_bot/statscard"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

type optionMap = map[string]*discordgo.ApplicationCommandInteractionDataOption

// type APIResponse struct {
// 	Data []models.PlayerData `json:"data"`
// }

func parseOptions(options []*discordgo.ApplicationCommandInteractionDataOption) (om optionMap) {
	om = make(optionMap)
	for _, opt := range options {
		om[opt.Name] = opt
	}
	return
}

// func interactionAuthor(i *discordgo.Interaction) *discordgo.User {
// 	if i.Member != nil {
// 		return i.Member.User
// 	}
// 	return i.User
// }

func loadEnv(key string, verbose bool) (string, error) {
	err := godotenv.Load("secrets.env")
	if err != nil {
		return "", fmt.Errorf("can't load secrets.env")
	}

	token := os.Getenv(key)
	if token == "" {
		return "", fmt.Errorf(("token not in environment"))
	}

	if verbose {
		fmt.Println("token retrieved successfully! debug: ", token)
	}

	return token, nil
}

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "stats",
		Description: "Displays the wynncraft stats for a player.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "username",
				Description: "The player's username or uuid.",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
		},
	},
}

func getPlayerStat(s *discordgo.Session, i *discordgo.InteractionCreate, opts optionMap) {
	builder := new(strings.Builder)
	// builder.WriteString(opts["username"].StringValue() + "\n")

	api_URL := "https://api.wynncraft.com/v3/player/%s?fullResult"
	username := opts["username"].StringValue()

	url := fmt.Sprintf(api_URL, username)
	resp, err := http.Get(url)
	if err != nil {
		builder.WriteString("cannot access url \n")
	} else {
		builder.WriteString("accessing: " + url + "\n")
		defer resp.Body.Close()

		var playerData models.PlayerData
		err = json.NewDecoder(resp.Body).Decode(&playerData)

		if err != nil {
			builder.WriteString("cannot decode json \n")
		} else {
			builder.WriteString("does emojis work :sob: \n")
			builder.WriteString("```CSS\n wait am overthinking this again``` \n")

			statscard.CreateStatsCard(playerData, "statscard/images", "testing.png")

			fmt.Printf("Debug output: %+v\n", &playerData.Guild)

		}

	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: builder.String(),
		},
	})

	if err != nil {
		log.Panicf("could not respond to interaction: %s", err)
	}
}

// messageCreate is a handler function that processes new messages.
// func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
// 	// Ignore messages from the bot itself
// 	if m.Author.ID == s.State.User.ID {
// 		return
// 	}

// 	// Respond to a specific command
// 	if m.Content == "!hello" {
// 		s.ChannelMessageSend(m.ChannelID, "Hello! How can I assist you?")
// 	}
// }

func main() {
	// load .env file
	token, err := loadEnv("DISCORD_TOKEN", true)
	if err != nil {
		fmt.Println("error retrieving token: ", err)
	}
	App, err := loadEnv("APP_ID", true)
	if err != nil {
		fmt.Println("error retrieving token: ", err)
	}
	Guild, err := loadEnv("GUILD_ID", true)
	if err != nil {
		fmt.Println("error retrieving token: ", err)
	}

	// start a discord session
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating discord session:", err)
		return
	}

	// add command handler
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		data := i.ApplicationCommandData()
		if data.Name != "stats" {
			return
		}

		getPlayerStat(s, i, parseOptions(data.Options))
	})

	// Open a connection to Discord
	err = session.Open()
	if err != nil {
		fmt.Println("Error opening connection:", err)
		return
	}

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as %s", r.User.String())
	})

	_, err = session.ApplicationCommandBulkOverwrite(App, Guild, commands)
	if err != nil {
		log.Fatalf("could not register commands: %s", err)
	}

	fmt.Println("Bot is now running. Press CTRL+C to exit.")

	// Wait for a termination signal (e.g., CTRL+C)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	// Cleanly close the Discord session
	session.Close()

}
