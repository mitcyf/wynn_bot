package statscard

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"net/http"
	"os"

	"wynn_bot/models"

	"github.com/fogleman/gg"
)

var colorMap = map[string]color.RGBA{
	"BLACK":      {R: 21, G: 21, B: 24, A: 255},
	"RED":        {R: 129, G: 34, B: 28, A: 255},
	"GREEN":      {R: 69, G: 91, B: 16, A: 255},
	"BROWN":      {R: 96, G: 62, B: 37, A: 255},
	"BLUE":       {R: 43, G: 49, B: 123, A: 255},
	"PURPLE":     {R: 100, G: 37, B: 135, A: 255},
	"CYAN":       {R: 16, G: 114, B: 114, A: 255},
	"SILVER":     {R: 115, G: 115, B: 111, A: 255},
	"GRAY":       {R: 52, G: 57, B: 60, A: 255},
	"PINK":       {R: 179, G: 102, B: 125, A: 255},
	"LIME":       {R: 93, G: 145, B: 22, A: 255},
	"YELLOW":     {R: 186, G: 158, B: 45, A: 255},
	"LIGHT_BLUE": {R: 42, G: 131, B: 160, A: 255},
	"MAGENTA":    {R: 145, G: 57, B: 138, A: 255},
	"ORANGE":     {R: 180, G: 93, B: 21, A: 255},
	"WHITE":      {R: 182, G: 187, B: 186, A: 255},
}

const width, height = 562, 952 // magic numbers
const header_width, header_height = width, 118
const image_width, image_height = width / 2, 434
const main_width, main_height = width, height - header_height
const banner_width, banner_height = width - image_width, height - header_height - footer_height
const footer_width, footer_height = width, 262

func RecolorImage(img image.Image, recolor color.RGBA) image.Image {
	bounds := img.Bounds()
	recoloredImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := img.At(x, y)
			_, _, _, alpha := originalColor.RGBA() // Preserve alpha channel
			recolored := color.RGBA{
				R: uint8(recolor.R),
				G: uint8(recolor.G),
				B: uint8(recolor.B),
				A: uint8(alpha >> 8), // Alpha is in the high byte
			}
			recoloredImg.Set(x, y, recolored)
		}
	}
	return recoloredImg
}

func LoadImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image file: %v", err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}
	return img, nil
}

func CreateBanner(data models.PlayerData) (*gg.Context, error) {
	api_URL := "https://api.wynncraft.com/v3/guild/%s"
	guild := data.Guild.Name

	url := fmt.Sprintf(api_URL, guild)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("cannot access url")
	} else {
		defer resp.Body.Close()

		var guildData models.GuildData
		err = json.NewDecoder(resp.Body).Decode(&guildData)

		if err != nil {
			return nil, fmt.Errorf("cannot decode json")
		} else {
			banner := gg.NewContext(banner_width, banner_height)

			banner.SetColor(colorMap[guildData.Banner.Base])
			banner.Clear()
			banner.Scale(banner_width/160.0, banner_height/320.0)

			for _, layer := range guildData.Banner.Layers {
				color := colorMap[layer.Colour]
				image_path := fmt.Sprintf("statscard/banner/%s.png", layer.Pattern)
				layer_base, err := LoadImage(image_path)
				if err != nil {
					return nil, fmt.Errorf("problem loading image")
				}
				recolored_image := RecolorImage(layer_base, color)
				if err != nil {
					return nil, fmt.Errorf("cannot draw layer")
				}
				banner.DrawImage(recolored_image, 0, 0)
			}

			return banner, nil
		}

	}
}

func CreateStatsCard(data models.PlayerData, outputDir string, fileName string) error {

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
	avatar_img, err := fetchAvatar(data.Username)
	if err != nil {
		return fmt.Errorf("failed to fetch avatar: %v", err)
	}
	scaling := min(image_width/512.0, image_height/869.0)
	avatar := gg.NewContext(int(math.Round(512*scaling)), int(math.Round(869*scaling)))
	avatar.Scale(scaling, scaling)
	avatar.DrawImage(avatar_img, 0, 0)

	card.DrawImageAnchored(avatar.Image(), image_width/2, header_height+image_height/2, 0.5, 0.5)

	// guild background
	banner, err := CreateBanner(data)
	if err != nil {
		return fmt.Errorf("error creating banner %v", err)
	}

	card.DrawImage(banner.Image(), image_width, header_height)

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
