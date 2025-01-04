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

const darkenFactor = 0.5

var darkenedMap = func() map[string]color.RGBA {
	darkened := make(map[string]color.RGBA)

	for name, col := range colorMap {
		darkened[name] = color.RGBA{
			R: uint8(float64(col.R) * darkenFactor),
			G: uint8(float64(col.G) * darkenFactor),
			B: uint8(float64(col.B) * darkenFactor),
			A: col.A, // Keep alpha unchanged
		}
	}

	return darkened
}()

const darken = 0.5

const width, height = 562, 952 // magic numbers
const headerWidth, headerHeight = width, 118
const imageWidth, imageHeight = width / 2, 434
const mainWidth, mainHeight = width, height - headerHeight
const bannerWidth, bannerHeight = width - imageWidth, height - headerHeight - footerHeight
const footerHidth, footerHeight = width, 262

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
	apiURL := "https://api.wynncraft.com/v3/guild/%s"
	guild := data.Guild.Name

	url := fmt.Sprintf(apiURL, guild)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("cannot access url")
	} else {
		defer resp.Body.Close()

		var guildData models.GuildData
		err = json.NewDecoder(resp.Body).Decode(&guildData)
		banner := gg.NewContext(bannerWidth, bannerHeight)

		if err != nil {
			banner.SetColor(darkenedMap["GRAY"])
			banner.Clear()
		} else {
			// banner := gg.NewContext(bannerWidth, bannerHeight)

			banner.SetColor(darkenedMap[guildData.Banner.Base])
			banner.Clear()
			banner.Scale(bannerWidth/160.0, bannerHeight/320.0)

			for _, layer := range guildData.Banner.Layers {
				color := darkenedMap[layer.Colour]
				imagePath := fmt.Sprintf("statscard/banner/%s.png", layer.Pattern)
				layerBase, err := LoadImage(imagePath)
				if err != nil {
					return nil, fmt.Errorf("problem loading image")
				}
				recoloredImage := RecolorImage(layerBase, color)
				if err != nil {
					return nil, fmt.Errorf("cannot draw layer")
				}
				banner.DrawImage(recoloredImage, 0, 0)
			}
		}
		return banner, nil

	}
}

func CreateStatsCard(data models.PlayerData, outputDir string, fileName string) error {

	card := gg.NewContext(width, height)

	// background
	card.SetColor(color.RGBA{R: 64, G: 64, B: 64, A: 255})
	card.Clear() // this only ends up in the footer tbh

	background, err := LoadImage("statscard/images/background.png")
	if err != nil {
		return fmt.Errorf("failed to load background: %v", err)
	}
	card.DrawImage(background, 0, 0)

	// boxes
	card.SetColor(color.RGBA{R: 0, G: 0, B: 0, A: 158})
	card.DrawRectangle(0, 0, headerWidth, headerHeight)

	card.DrawRectangle(0, headerHeight+imageHeight, imageWidth, bannerHeight-imageHeight)
	card.Fill()

	// player avatar

	avatarImg, err := fetchAvatar(data.Username)
	if err != nil {
		return fmt.Errorf("failed to fetch avatar: %v", err)
	}

	scaling := min(imageWidth/512.0, imageHeight/869.0) * 0.9
	avatar := gg.NewContext(int(math.Round(512*scaling)), int(math.Round(869*scaling)))
	avatar.Scale(scaling, scaling)
	avatar.DrawImage(avatarImg, 0, 0)

	card.DrawImageAnchored(avatar.Image(), imageWidth/2, headerHeight+imageHeight/2, 0.5, 0.4)

	// guild background
	banner, err := CreateBanner(data)
	if err != nil {
		return fmt.Errorf("error creating banner %v", err)
	}

	card.DrawImage(banner.Image(), imageWidth, headerHeight)

	// header content
	rankBadge := (*data.RankBadge)[15 : len(*data.RankBadge)-4]
	rankImg, err := LoadImage(fmt.Sprintf("statscard/ranks/%s.png", rankBadge))
	if err != nil {
		return fmt.Errorf("error loading rank badge: %v", err)
	}
	badge := gg.NewContext(int(math.Round(headerWidth/3.0)), int(math.Round(headerHeight/3.0)))
	badge.Scale(math.Round(headerHeight/45.0), math.Round(headerHeight/45.0))
	badge.DrawImage(rankImg, 0, 0)

	card.DrawImage(badge.Image(), 15, 10)

	if err := card.LoadFontFace("/Library/Fonts/Impact.ttf", 96); err != nil {
		panic(err)
	}

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
