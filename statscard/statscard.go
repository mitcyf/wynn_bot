package statscard

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
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

// const mainWidth, mainHeight = width, height - headerHeight
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

func CreateBanner(data models.PlayerData) (*gg.Context, *models.GuildData, error) {
	apiURL := "https://api.wynncraft.com/v3/guild/%s"

	var bannerBase string
	var bannerLayers []models.BannerLayer
	var guildData *models.GuildData

	if data.Guild != nil {
		guild := data.Guild.Name
		url := fmt.Sprintf(apiURL, guild)
		resp, err := http.Get(url)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot access url")
		} else {
			defer resp.Body.Close()

			// var guildData models.GuildData
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
			return nil, nil, fmt.Errorf("problem loading image")
		}
		recoloredImage := RecolorImage(layerBase, color)
		banner.DrawImage(recoloredImage, 0, 0)
	}
	if data.Guild != nil {
		return banner, guildData, nil
	} else {
		return banner, nil, nil
	}
}

func formatNumber(n float64) string {
	if n >= 1e9 {
		return fmt.Sprintf("%.1fB", n/1e9)
	} else if n >= 1e6 {
		return fmt.Sprintf("%.1fM", n/1e6)
	} else if n >= 1e3 {
		return fmt.Sprintf("%.1fK", n/1e3)
	} else {
		return fmt.Sprintf("%.1f", n) // For numbers less than 1, use default formatting
	}
}

func drawPieChart(card *gg.Context, x, y, outerRadius, innerRadius float64, data map[string]int, colors map[string]string) {
	total := 0
	for _, value := range data {
		total += value
	}

	if total == 0 {
		return
	}

	startAngle := rand.Float64() * 2 * math.Pi
	for key, value := range data {
		if value > 0 {
			percentage := float64(value) / float64(total)
			endAngle := startAngle + (percentage * 2 * math.Pi)

			card.SetHexColor(colors[key])
			card.MoveTo(x, y)
			card.DrawArc(x, y, outerRadius, startAngle, endAngle)
			card.LineTo(x+innerRadius*math.Cos(endAngle), y+innerRadius*math.Sin(endAngle))
			card.DrawArc(x, y, innerRadius, endAngle, startAngle)
			card.ClosePath()
			card.Fill()

			startAngle = endAngle
		}
	}
}

func CreateStatsCard(data models.PlayerData, outputDir string, fileName string) error {

	card := gg.NewContext(width, height)

	// background
	card.SetColor(color.RGBA{R: 19, G: 0, B: 25, A: 255})
	card.Clear() // this only ends up in the footer tbh

	background, err := LoadImage("statscard/images/background.png")
	if err != nil {
		return fmt.Errorf("failed to load background: %v", err)
	}
	card.DrawImage(background, 0, 0)

	// boxes
	card.SetColor(color.RGBA{R: 0, G: 0, B: 0, A: 120})
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
	banner, guild, err := CreateBanner(data)
	if err != nil {
		return fmt.Errorf("error creating banner %v", err)
	}

	card.DrawImage(banner.Image(), imageWidth, headerHeight)

	// header content

	// rank badge

	rankImg, err := LoadImage("statscard/ranks_upscale/rank_none.png")
	if err != nil {
		return fmt.Errorf("error loading none badge: %v", err)
	}

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
	card.SetColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	if err := card.LoadFontFace("statscard/fonts/comfortaa_bold.ttf", 16); err != nil {
		panic(err)
	}
	card.DrawStringAnchored(subtitle1, 20, float64(rankImg.Bounds().Max.Y)+32, 0, 0)
	card.DrawStringAnchored(subtitle2, 20, float64(rankImg.Bounds().Max.Y)+55, 0, 0)

	// guild content

	if guild != nil {

		guild1 := strings.ToLower(data.Guild.Rank) + " of " + guild.Prefix

		// this is such a stupid data structure but idk if fixing it to solve the problem once is worth it
		memberInfo := models.MemberInfo{}
		if data.Guild.Rank == "OWNER" {
			memberInfo = guild.Members.Owner[data.Username]
		} else if data.Guild.Rank == "CHIEF" {
			memberInfo = guild.Members.Chief[data.Username]
		} else if data.Guild.Rank == "STRATEGIST" {
			memberInfo = guild.Members.Strategist[data.Username]
		} else if data.Guild.Rank == "CAPTAIN" {
			memberInfo = guild.Members.Captain[data.Username]
		} else if data.Guild.Rank == "RECRUITER" {
			memberInfo = guild.Members.Recruiter[data.Username]
		} else if data.Guild.Rank == "RECRUIT" {
			memberInfo = guild.Members.Recruit[data.Username]
		}
		fmt.Println(memberInfo.Joined)

		guild2 := "since " + ParseTime(memberInfo.Joined)
		guild3 := guild.Name + ", lv " + strconv.Itoa(guild.Level)
		guild4 := formatNumber(float64(memberInfo.Contributed)) + " xp contributed (#" + strconv.Itoa(*memberInfo.ContributionRank) + ")"

		if err := card.LoadFontFace("statscard/fonts/comfortaa_bold.ttf", 20); err != nil {
			panic(err)
		}
		card.DrawStringAnchored(guild1, 20, imageHeight+headerHeight+30, 0, 0.5)

		if err := card.LoadFontFace("statscard/fonts/comfortaa_bold.ttf", 15); err != nil {
			panic(err)
		}
		card.DrawStringAnchored(guild2, 20, imageHeight+headerHeight+60, 0, 0.5)
		card.DrawStringAnchored(guild3, 20, imageHeight+headerHeight+80, 0, 0.5)

		card.DrawStringAnchored(guild4, 20, imageHeight+headerHeight+120, 0, 0.5)
	}

	// main content

	// could this have been cleaner? yes. did i know struct when I started this? no
	// should've cleaned this up instead of just slapping it into new renderer
	// whatever, ugly code still does the job

	headers := make(map[string]string)
	headersY := make(map[string]int)
	labels := make(map[string]string)
	labelsY := make(map[string]int)
	raids := make(map[string]string)
	raidsY := make(map[string]int)
	raidsColor := make(map[string]string)
	valuesRight := make(map[string]string)
	valuesRightY := make(map[string]int)
	valuesLeft := make(map[string]string)
	valuesLeftY := make(map[string]int)

	nog := data.GlobalData.Raids.List["Nest of the Grootslangs"]
	nol := data.GlobalData.Raids.List["Orphion's Nexus of Light"]
	tcc := data.GlobalData.Raids.List["The Canyon Colossus"]
	tna := data.GlobalData.Raids.List["The Nameless Anomaly"]

	headers["player stats"] = "player stats"
	headersY["player stats"] = 40

	spacing := 22
	statsY := headersY["player stats"] + 5

	labels["playtime"] = "playtime"
	labelsY["playtime"] = statsY + spacing
	valuesRight["playtime"] = strconv.Itoa(int(math.Round(data.Playtime))) + " hr"
	valuesRightY["playtime"] = statsY + spacing

	labels["total levels"] = "total levels"
	labelsY["total levels"] = statsY + spacing*2
	valuesRight["total levels"] = strconv.Itoa(data.GlobalData.TotalLevel)
	valuesRightY["total levels"] = statsY + spacing*2

	labels["kills"] = "kills"
	labelsY["kills"] = statsY + spacing*4
	valuesRight["kills"] = strconv.Itoa(data.GlobalData.KilledMobs)
	valuesRightY["kills"] = statsY + spacing*4

	labels["chests"] = "chests"
	labelsY["chests"] = statsY + spacing*5
	valuesRight["chests"] = strconv.Itoa(data.GlobalData.ChestsFound)
	valuesRightY["chests"] = statsY + spacing*5

	labels["dungeons"] = "dungeons"
	labelsY["dungeons"] = statsY + spacing*6
	valuesRight["dungeons"] = strconv.Itoa(data.GlobalData.Dungeons.Total)
	valuesRightY["dungeons"] = statsY + spacing*6

	labels["quests"] = "quests"
	labelsY["quests"] = statsY + spacing*7
	valuesRight["quests"] = strconv.Itoa(data.GlobalData.CompletedQuests)
	valuesRightY["quests"] = statsY + spacing*7

	labels["wars"] = "wars"
	labelsY["wars"] = statsY + spacing*9
	valuesRight["wars"] = strconv.Itoa(data.GlobalData.Wars)
	valuesRightY["wars"] = statsY + spacing*9

	headers["raid completions"] = "raids completions"
	headersY["raid completions"] = statsY + spacing*11

	raidsYCoord := headersY["raid completions"] + 5

	labels["total"] = "total"
	labelsY["total"] = raidsYCoord + spacing
	valuesLeft["total"] = strconv.Itoa(data.GlobalData.Raids.Total)
	valuesLeftY["total"] = raidsYCoord + spacing

	labels["nog"] = "nog"
	labelsY["nog"] = raidsYCoord + spacing*3
	raids["nog"] = strconv.Itoa(nog)
	raidsY["nog"] = raidsYCoord + spacing*3
	raidsColor["nog"] = "#93c47d"

	labels["nol"] = "nol"
	labelsY["nol"] = raidsYCoord + spacing*4
	raids["nol"] = strconv.Itoa(nol)
	raidsY["nol"] = raidsYCoord + spacing*4
	raidsColor["nol"] = "#ffd966"

	labels["tcc"] = "tcc"
	labelsY["tcc"] = raidsYCoord + spacing*5
	raids["tcc"] = strconv.Itoa(tcc)
	raidsY["tcc"] = raidsYCoord + spacing*5
	raidsColor["tcc"] = "#e06666"

	labels["tna"] = "tna"
	labelsY["tna"] = raidsYCoord + spacing*6
	raids["tna"] = strconv.Itoa(tna)
	raidsY["tna"] = raidsYCoord + spacing*6
	raidsColor["tna"] = "#8e7cc3"

	headers["leaderboards"] = "leaderboards"
	headersY["leaderboards"] = raidsYCoord + spacing*8

	leaderboardsY := headersY["leaderboards"] + 5

	labels["completion"] = "completion"
	labelsY["completion"] = leaderboardsY + spacing
	valuesRight["completion"] = "#" + strconv.Itoa(data.Ranking.GlobalPlayerContent)
	valuesRightY["completion"] = leaderboardsY + spacing

	labels["professions"] = "professions"
	labelsY["professions"] = leaderboardsY + spacing*2
	valuesRight["professions"] = "#" + strconv.Itoa(data.Ranking.ProfessionsGlobalLevel)
	valuesRightY["professions"] = leaderboardsY + spacing*2

	labels["wars won"] = "wars won"
	labelsY["wars won"] = leaderboardsY + spacing*4
	valuesRight["wars won"] = "#" + strconv.Itoa(data.Ranking.WarsCompletion)
	valuesRightY["wars won"] = leaderboardsY + spacing*4

	card.SetColor(color.RGBA{R: 0, G: 0, B: 0, A: 91})
	card.DrawRoundedRectangle(10+imageWidth, float64(statsY)+float64(spacing)*3.5+2, width-imageWidth-20, float64(spacing)*11-5, 15)
	card.DrawRoundedRectangle(10+imageWidth, float64(raidsYCoord)+float64(spacing)*3.5+2, width-imageWidth-20, float64(spacing)*8-5, 15)
	card.DrawRoundedRectangle(10+imageWidth, float64(leaderboardsY)+float64(spacing)*3.5+2, width-imageWidth-20, float64(spacing)*6-5, 15)
	card.Fill()

	card.SetHexColor("#FFFFFF")

	if err := card.LoadFontFace("statscard/fonts/comfortaa_bold.ttf", 24); err != nil {
		panic(err)
	}
	for key, y := range headersY {
		card.DrawStringAnchored(headers[key], 20+imageWidth, headerHeight+float64(y), 0, 0.5)
	}

	if err := card.LoadFontFace("statscard/fonts/comfortaa_bold.ttf", 16); err != nil {
		panic(err)
	}
	for key, y := range labelsY {
		card.DrawStringAnchored(labels[key], 20+imageWidth, headerHeight+float64(y), 0, 0.5)
	}
	for key, y := range valuesRightY {
		card.DrawStringAnchored(valuesRight[key], 20+imageWidth+160, headerHeight+float64(y), 0, 0.5)
	}
	for key, y := range valuesLeftY {
		card.DrawStringAnchored(valuesLeft[key], 20+imageWidth+90, headerHeight+float64(y), 0, 0.5)
	}
	for key, y := range raidsY {
		card.SetHexColor(raidsColor[key])
		card.DrawStringAnchored(raids[key], 20+imageWidth+90, headerHeight+float64(y), 0, 0.5)
	}

	drawPieChart(card, 490, float64(raidsYCoord+180), 45, 35, map[string]int{
		"nog": nog,
		"nol": nol,
		"tcc": tcc,
		"tna": tna,
	}, raidsColor)

	footerImg, err := LoadImage("statscard/images/footer.png")
	if err != nil {
		return fmt.Errorf("cannot load footer image")
	}
	card.DrawImage(footerImg, 0, height-footerHeight)

	if err := card.LoadFontFace("statscard/fonts/comfortaa_bold.ttf", 24); err != nil {
		panic(err)
	}
	card.SetHexColor("#ffffff")
	card.DrawStringAnchored("completion", width/2, height-footerHeight+40, 0.5, 0)

	classes := []string{
		"ARCHER",
		"WARRIOR",
		"ASSASSIN",
		"MAGE",
		"SHAMAN",
	}

	levels := map[string]int{
		"ARCHER":   0,
		"WARRIOR":  0,
		"ASSASSIN": 0,
		"MAGE":     0,
		"SHAMAN":   0,
	}

	classColors := map[string]string{ // hsv sat 70 val 70, 50
		"ARCHER":   "#8936b3",
		"WARRIOR":  "#b34036",
		"ASSASSIN": "#36b3b3",
		"MAGE":     "#b38936",
		"SHAMAN":   "#5fb336",
	}
	classPerfectionColors := map[string]string{
		"ARCHER":   "#622680",
		"WARRIOR":  "#802e26",
		"ASSASSIN": "#268080",
		"MAGE":     "#806226",
		"SHAMAN":   "#448026",
	}
	for _, char := range data.Characters {
		if char.Level > levels[char.Type] {
			levels[char.Type] = char.Level
		}
	}

	card.LoadFontFace("statscard/fonts/minecraft.ttf", 22)

	for index, class := range classes {
		x := float64(index)*width/5.0 + width/10.0
		y := height - footerHeight + 120.0
		drawPieChart(card, x, y, 45, 35, map[string]int{
			"lvl": levels[class],
			"rem": 105 - levels[class],
		}, map[string]string{
			"lvl": classColors[class],
			"rem": "#000000",
		})
		if levels[class] == 106 {
			drawPieChart(card, x, y, 35, 30, map[string]int{
				"": 1,
			}, map[string]string{
				"": classPerfectionColors[class],
			})
		}
		classImg, err := LoadImage(fmt.Sprintf("statscard/classes/%s.png", class))
		if err != nil {
			// return fmt.Errorf("cannot load class images for " + class)
		}
		card.DrawImageAnchored(classImg, int(math.Round(x)), int(math.Round(y)-11), 0.5, 0.5)
		card.SetHexColor(classColors[class])
		card.DrawStringAnchored(strconv.Itoa(levels[class]), x, y+11, 0.5, 0.5)
	}

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
