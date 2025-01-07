package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wynn_bot/chartings"
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
	{
		Name:        "charttest",
		Description: "Testing command for the charting function",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "test",
				Description: "not used",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
		},
	},
}

func getPlayerStat(s *discordgo.Session, i *discordgo.InteractionCreate, opts optionMap) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Generating stats card, please wait...",
		},
	})

	if err != nil {
		log.Panicf("could not respond to interaction: %s", err)
		return
	}

	api_URL := "https://api.wynncraft.com/v3/player/%s?fullResult"
	username := opts["username"].StringValue()

	url := fmt.Sprintf(api_URL, username)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to access URL: %s", err)
		_, _ = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: stringPointer("Failed to access the player data URL."),
		})
		return
	}
	defer resp.Body.Close()

	var playerData models.PlayerData
	err = json.NewDecoder(resp.Body).Decode(&playerData)
	if err != nil {
		log.Printf("Failed to decode JSON: %s", err)
		_, _ = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: stringPointer("Failed to decode player data JSON."),
		})
		return
	}

	// Generate the stats card
	imagePath := "statcard.png"
	err = statscard.CreateStatsCard(playerData, "statscard/images", imagePath)
	if err != nil {
		log.Printf("Failed to generate stats card: %s", err)
		_, _ = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: stringPointer("Failed to generate stats card."),
		})
		return
	}

	// Open the generated image
	file, err := os.Open("statscard/images/statcard.png")
	if err != nil {
		log.Printf("Failed to open stats card: %s", err)
		_, _ = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: stringPointer("Failed to open the generated stats card."),
		})
		return
	}
	defer file.Close()

	// Retry logic for editing the response
	maxRetries := 5
	for attempts := 0; attempts < maxRetries; attempts++ {
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: stringPointer(""),
			Files: []*discordgo.File{
				{
					Name:   "statcard.png",
					Reader: file,
				},
			},
		})
		if err == nil {
			break // Success
		}
		log.Printf("Retry %d: Failed to edit interaction response with image: %s", attempts+1, err)
		time.Sleep(time.Second) // Wait before retrying
	}

	if err != nil {
		log.Printf("Failed to edit interaction response after retries: %s", err)
		_, _ = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: stringPointer("Failed to edit interaction response after multiple attempts."),
		})
	}
}

func ChartTest(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Chart generating, please wait...",
		},
	})
	if err != nil {
		log.Printf("Could not respond to interaction: %s", err)
		return
	}

	// Example chart data
	data := &chartings.ChartData{
		X:        []float64{1, 2, 3, 4, 5},
		Y:        []float64{10, 20, 15, 25, 30},
		XLegends: []string{"A", "B", "C", "D", "E"},
		YLegends: []string{"Low", "Medium", "High", "Very High"},
		XLabel:   "Categories",
		YLabel:   "Values",
		Title:    "Test Chart",
		Desc:     "This is a test chart",
	}

	buffer, err := chartings.Render(data)
	if err != nil {
		log.Printf("Failed to render chart: %s", err)
		_, _ = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: stringPointer("Failed to generate chart."),
		})
		return
	}

	// Retry logic for editing the response
	maxRetries := 5
	for attempts := 0; attempts < maxRetries; attempts++ {
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: stringPointer("Chart generated!"),
			Files: []*discordgo.File{
				{
					Name:   "chart.png",
					Reader: buffer,
				},
			},
		})
		if err == nil {
			break // Success
		}
		log.Printf("Retry %d: Failed to edit interaction response with image: %s", attempts+1, err)
		time.Sleep(time.Second) // Wait before retrying
	}

	if err != nil {
		log.Printf("Failed to edit interaction response after retries: %s", err)
		_, _ = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: stringPointer("Failed to edit interaction response after multiple attempts."),
		})
	}
}

func stringPointer(s string) *string {
	return &s
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
		if data.Name == "stats" {
			getPlayerStat(s, i, parseOptions(data.Options))
		} else if data.Name == "charttest" {

		}

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
