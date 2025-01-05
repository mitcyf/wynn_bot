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
	"time"

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

const width, height = 562, 952 // magic numbers
const headerWidth, headerHeight = width, 102
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

func ParseTime(rawTime string) string {
	parsedTime, err := time.Parse(time.RFC3339Nano, rawTime)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return ""
	}

	return parsedTime.Format("Jan 02, 2006")

}

func TimeAgo(isoTimestamp string) string {
	// Parse the ISO 8601 timestamp
	parsedTime, err := time.Parse(time.RFC3339Nano, isoTimestamp)
	if err != nil {
		return "Invalid timestamp"
	}

	// Calculate the difference
	now := time.Now()
	duration := now.Sub(parsedTime)

	// Helper function for pluralization
	pluralize := func(value int, unit string) string {
		if value == 1 {
			return fmt.Sprintf("%d %s ago", value, unit) // Singular
		}
		return fmt.Sprintf("%d %ss ago", value, unit) // Plural
	}

	// Determine the appropriate time unit
	if duration < time.Minute {
		return pluralize(int(duration.Seconds()), "second")
	} else if duration < time.Hour {
		return pluralize(int(duration.Minutes()), "minute")
	} else if duration < 24*time.Hour {
		return pluralize(int(duration.Hours()), "hour")
	} else if duration < 30*24*time.Hour {
		return pluralize(int(duration.Hours()/24), "day")
	} else if duration < 12*30*24*time.Hour {
		return pluralize(int(duration.Hours()/(24*30)), "month")
	} else {
		return pluralize(int(duration.Hours()/(24*365)), "year")
	}
}

func CreateBanner(data models.PlayerData) (*gg.Context, error) {
	apiURL := "https://api.wynncraft.com/v3/guild/%s"

	var bannerBase string
	var bannerLayers []models.BannerLayer

	if data.Guild != nil {
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
				banner.SetColor(darkenedMap["SILVER"])
				banner.Clear()
			} else {
				// banner := gg.NewContext(bannerWidth, bannerHeight)
				bannerBase = guildData.Banner.Base
				bannerLayers = guildData.Banner.Layers
			}

		}
	} else {
		bannerBase = "SILVER"
		bannerLayers = []models.BannerLayer{}

		// modify to customize banner if no guild
		bannerLayers = append(bannerLayers,
			models.BannerLayer{Colour: "GRAY", Pattern: "BORDER"},
			models.BannerLayer{Colour: "GRAY", Pattern: "MOJANG"},
		)
	}
	banner := gg.NewContext(bannerWidth, bannerHeight)

	banner.SetColor(darkenedMap[bannerBase])
	banner.Clear()
	banner.Scale(bannerWidth/160.0, bannerHeight/320.0)

	for _, layer := range bannerLayers {
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

	return banner, nil
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

	card.DrawImageAnchored(avatar.Image(), imageWidth/2, headerHeight+imageHeight/2, 0.5, 0.5)

	// guild background
	banner, err := CreateBanner(data)
	if err != nil {
		return fmt.Errorf("error creating banner %v", err)
	}

	card.DrawImage(banner.Image(), imageWidth, headerHeight)

	// header content

	// rank badge

	rankImg, err := LoadImage(fmt.Sprintf("statscard/ranks_upscale/rank_none.png"))

	if data.RankBadge != nil {
		rankBadge := (*data.RankBadge)[15 : len(*data.RankBadge)-4]
		rankImg, err = LoadImage(fmt.Sprintf("statscard/ranks_upscale/%s.png", rankBadge))
		if err != nil {
			return fmt.Errorf("error loading rank badge: %v", err)
		}
	} else {

	}
	badge := gg.NewContext(int(math.Round(headerWidth/2.0)), int(math.Round(headerHeight/3.0)))
	// badge.Scale(math.Round(headerHeight/30.0), math.Round(headerHeight/30.0))
	badge.DrawImage(rankImg, 0, 0)

	card.DrawImageAnchored(badge.Image(), 15, 30, 0, 0.5)

	// text

	if data.LegacyRankColour != nil {
		card.SetHexColor(data.LegacyRankColour.Sub)
	} else {
		card.SetColor(color.RGBA{R: 221, G: 225, B: 218, A: 255})
	}

	if err := card.LoadFontFace("statscard/fonts/minecraft.ttf", 42); err != nil {
		panic(err)
	}
	card.DrawStringAnchored(data.Username, float64(rankImg.Bounds().Max.X)+30, 30, 0, 0.4)

	subtitle1 := ""
	subtitle2 := "" // line breaks don't work with gg for some reason
	subtitle1 += "first joined " + ParseTime(data.FirstJoin)
	if data.Online {
		subtitle2 += "currently online on world " + *data.Server
	} else {
		subtitle2 += "last seen " + TimeAgo(data.LastJoin)
		if data.Server != nil {
			subtitle2 += " on world " + *data.Server
		}
	}
	card.SetColor(color.RGBA{R: 221, G: 225, B: 218, A: 255})
	if err := card.LoadFontFace("statscard/fonts/comfortaa.ttf", 16); err != nil {
		panic(err)
	}
	card.DrawStringAnchored(subtitle1, 20, float64(rankImg.Bounds().Max.Y)+32, 0, 0)
	card.DrawStringAnchored(subtitle2, 20, float64(rankImg.Bounds().Max.Y)+55, 0, 0)

	// saving the image

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
