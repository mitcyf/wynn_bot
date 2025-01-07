package chartings

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/fogleman/gg"
)

// "fmt"
// "image"
// "os"
// "bytes"

// "github.com/fogleman/gg"

type ChartData struct {
	X        []float64 // actual X position
	Y        []float64 // actual Y position
	XLegends []string  // i'm assuming this is what is displayed
	YLegends []string  // and same here
	XLabel   string
	YLabel   string
	Title    string
	Desc     string
	Width    int
	Height   int
}

// Returns image buffer and error
func Render(d *ChartData) (*bytes.Buffer, error) {
	fmt.Println("test")
	chart := gg.NewContext(d.Width, d.Height)
	chart.SetRGB(1, 1, 1) // Background color
	chart.Clear()

	width := d.Width
	height := d.Height

	// Draw title
	chart.SetRGB(0, 0, 0)
	chart.LoadFontFace("statscard/fonts/comfortaa.ttf", 24)
	chart.DrawStringAnchored(d.Title, float64(width)/2, 30, 0.5, 0.5)

	// Calculate margins
	margin := 50.0
	chartWidth := float64(width) - 2*margin
	chartHeight := float64(height) - 2*margin - 40 // Adjust for title

	// Find min/max for axes
	xMin, xMax := math.Inf(1), math.Inf(-1)
	yMin, yMax := math.Inf(1), math.Inf(-1)
	for i := range d.X {
		xMin = math.Min(xMin, d.X[i])
		xMax = math.Max(xMax, d.X[i])
		yMin = math.Min(yMin, d.Y[i])
		yMax = math.Max(yMax, d.Y[i])
	}

	// Draw axes
	chart.SetLineWidth(2)
	chart.DrawLine(margin, float64(height)-margin, margin, margin)
	chart.DrawLine(margin, float64(height)-margin, float64(width)-margin, float64(height)-margin)
	chart.Stroke()

	// Draw labels
	chart.LoadFontFace("statscard/fonts/comfortaa.ttf", 14)
	chart.DrawStringAnchored(d.XLabel, float64(width)/2, float64(height)-20, 0.5, 0.5)
	chart.DrawStringAnchored(d.YLabel, 20, float64(height)/2, 0.5, 0.5)

	// Plot points
	for i := range d.X {
		xNorm := (d.X[i] - xMin) / (xMax - xMin)
		yNorm := (d.Y[i] - yMin) / (yMax - yMin)

		x := margin + xNorm*chartWidth
		y := float64(height) - margin - yNorm*chartHeight

		chart.DrawCircle(x, y, 4)
		chart.Fill()
	}

	// Save to buffer
	buffer := new(bytes.Buffer)
	if err := chart.EncodePNG(buffer); err != nil {
		return nil, err
	}

	return buffer, nil
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
	data := &ChartData{
		X:        []float64{1, 2, 3, 4, 5},
		Y:        []float64{10, 20, 15, 25, 30},
		XLegends: []string{"A", "B", "C", "D", "E"},
		YLegends: []string{"Low", "Medium", "High", "Very High"},
		XLabel:   "Categories",
		YLabel:   "Values",
		Title:    "Test Chart",
		Desc:     "This is a test chart",
	}

	buffer, err := Render(data)
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
