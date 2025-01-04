package statscard

import (
	"fmt"
	"image"
	"image/color"
	"net/http"

	"wynn_bot/models"

	"github.com/fogleman/gg"
)

func CreateStatsCard(data models.PlayerData, outputDir string, fileName string) error {
	const width, height = 1500, 2000 // magic numbers
	const header_width, header_height = width, 150
	const image_width, image_height = 512, 869
	const main_width, main_height = width, height - header_height

	card := gg.NewContext(width, height)

	// background
	card.SetColor(color.RGBA{R: 64, G: 64, B: 64, A: 255})
	card.Clear()

	// boxes
	card.SetColor(color.RGBA{R: 191, G: 166, B: 191, A: 255})
	card.DrawRectangle(0, 0, header_width, header_height)
	card.Fill()

	card.SetColor(color.RGBA{R: 15, G: 15, B: 15, A: 255})
	card.DrawRectangle(0, header_height, image_width, image_height)
	card.Fill()

	// player avatar
	avatar, err := fetchAvatar(data.Username)
	if err != nil {
		return fmt.Errorf("failed to fetch avatar: %v", err)
	}
	card.DrawImage(avatar, 0, header_height)

	// guild background
	banner_test, err := fetchImage("https://beta-cdn.wynncraft.com/nextgen/banners/CIRCLE_MIDDLE.svg")
	if err != nil {
		return fmt.Errorf("whats up %v", err)
	}
	card.DrawImage(banner_test, image_width, header_height)

	saveErr := card.SavePNG(outputDir + "/" + fileName)
	if saveErr != nil {
		return fmt.Errorf("failed to save image: %v", saveErr)
	}

	return nil

}

func fetchImage(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image: %v", err)
	}
	defer resp.Body.Close()

	// Decode the image
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	return img, nil
}

func fetchAvatar(username string) (image.Image, error) {
	// Build the API URL
	url := fmt.Sprintf("https://nmsr.nickac.dev/fullbody/%s", username)

	return fetchImage(url)
}
