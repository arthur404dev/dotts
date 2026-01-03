package theme

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/rivo/uniseg"
)

// ApplyGradient renders a string with a horizontal color gradient
func ApplyGradient(input string, colors ...lipgloss.TerminalColor) string {
	if input == "" || len(colors) < 2 {
		return input
	}

	var b strings.Builder
	clusters := toGraphemeClusters(input)
	if len(clusters) == 0 {
		return input
	}

	ramp := blendColors(len(clusters), colors...)
	if len(ramp) == 0 {
		return input
	}

	for i, cluster := range clusters {
		style := lipgloss.NewStyle().Foreground(ramp[i])
		b.WriteString(style.Render(cluster))
	}

	return b.String()
}

// ApplyBoldGradient renders a bold string with a horizontal color gradient
func ApplyBoldGradient(input string, colors ...lipgloss.TerminalColor) string {
	if input == "" || len(colors) < 2 {
		return input
	}

	var b strings.Builder
	clusters := toGraphemeClusters(input)
	if len(clusters) == 0 {
		return input
	}

	ramp := blendColors(len(clusters), colors...)
	if len(ramp) == 0 {
		return input
	}

	for i, cluster := range clusters {
		style := lipgloss.NewStyle().Foreground(ramp[i]).Bold(true)
		b.WriteString(style.Render(cluster))
	}

	return b.String()
}

// toGraphemeClusters splits a string into grapheme clusters
func toGraphemeClusters(input string) []string {
	var clusters []string
	gr := uniseg.NewGraphemes(input)
	for gr.Next() {
		clusters = append(clusters, gr.Str())
	}
	return clusters
}

// termColorToColorful converts a lipgloss.TerminalColor to colorful.Color
func termColorToColorful(tc lipgloss.TerminalColor) colorful.Color {
	// For lipgloss.Color (hex strings), we can extract the hex
	if c, ok := tc.(lipgloss.Color); ok {
		cf, err := colorful.Hex(string(c))
		if err == nil {
			return cf
		}
	}

	// For ANSI colors or fallback, use a default approach
	// We'll map common ANSI colors to approximate hex values
	if c, ok := tc.(lipgloss.ANSIColor); ok {
		return ansiToColorful(int(c))
	}

	// Fallback to white
	return colorful.Color{R: 1, G: 1, B: 1}
}

// ansiToColorful maps ANSI color codes to approximate colorful.Color values
func ansiToColorful(code int) colorful.Color {
	// Standard ANSI color approximations
	ansiColors := map[int]string{
		0:  "#000000", // Black
		1:  "#cc0000", // Red
		2:  "#00cc00", // Green
		3:  "#cccc00", // Yellow
		4:  "#0000cc", // Blue
		5:  "#cc00cc", // Magenta
		6:  "#00cccc", // Cyan
		7:  "#cccccc", // White
		8:  "#666666", // Bright Black
		9:  "#ff0000", // Bright Red
		10: "#00ff00", // Bright Green
		11: "#ffff00", // Bright Yellow
		12: "#0000ff", // Bright Blue
		13: "#ff00ff", // Bright Magenta
		14: "#00ffff", // Bright Cyan
		15: "#ffffff", // Bright White
	}

	if hex, ok := ansiColors[code]; ok {
		cf, _ := colorful.Hex(hex)
		return cf
	}

	return colorful.Color{R: 1, G: 1, B: 1}
}

// blendColors creates a color ramp between stops using HCL interpolation
func blendColors(size int, stops ...lipgloss.TerminalColor) []lipgloss.Color {
	if len(stops) < 2 || size == 0 {
		return nil
	}

	// Convert to colorful.Color for HCL blending
	cfStops := make([]colorful.Color, len(stops))
	for i, stop := range stops {
		cfStops[i] = termColorToColorful(stop)
	}

	// Handle single character case
	if size == 1 {
		return []lipgloss.Color{lipgloss.Color(cfStops[0].Hex())}
	}

	numSegments := len(cfStops) - 1
	result := make([]lipgloss.Color, 0, size)

	// Distribute sizes across segments
	segmentSizes := make([]int, numSegments)
	baseSize := size / numSegments
	remainder := size % numSegments

	for i := range numSegments {
		segmentSizes[i] = baseSize
		if i < remainder {
			segmentSizes[i]++
		}
	}

	// Generate colors for each segment
	for i := range numSegments {
		c1 := cfStops[i]
		c2 := cfStops[i+1]
		segSize := segmentSizes[i]

		for j := range segSize {
			t := 0.0
			if segSize > 1 {
				t = float64(j) / float64(segSize-1)
			}
			// HCL blending stays in gamut and looks smooth
			blended := c1.BlendHcl(c2, t)
			result = append(result, lipgloss.Color(blended.Hex()))
		}
	}

	return result
}

// FadeToBackground creates a gradient that fades from a color to the background
func FadeToBackground(input string, startColor lipgloss.TerminalColor) string {
	t := Current()
	return ApplyGradient(input, startColor, t.BgSubtle)
}
