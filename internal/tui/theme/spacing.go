// Package theme provides the visual theming system for the TUI.
package theme

import "strings"

// Spacing constants define standard spacing values in cells (horizontal) or lines (vertical).
// These provide consistent visual rhythm across all TUI components.
const (
	SpaceNone = 0 // No spacing
	SpaceXS   = 1 // Extra small: tight spacing for compact layouts
	SpaceSM   = 2 // Small: default spacing between related elements
	SpaceMD   = 3 // Medium: spacing between distinct sections
	SpaceLG   = 4 // Large: major section separation
	SpaceXL   = 6 // Extra large: page-level separation
)

// Gap returns a horizontal gap (string of spaces) of the specified width.
// Returns an empty string if size is zero or negative.
func Gap(size int) string {
	if size <= 0 {
		return ""
	}
	return strings.Repeat(" ", size)
}

// VGap returns a vertical gap (empty lines) of the specified height.
// Returns an empty string if size is zero or negative.
func VGap(size int) string {
	if size <= 0 {
		return ""
	}
	return strings.Repeat("\n", size)
}

// Padding represents padding values for all four sides of a component.
// Values are in cells (horizontal) or lines (vertical).
type Padding struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

// PaddingAll creates uniform padding on all sides.
func PaddingAll(size int) Padding {
	return Padding{
		Top:    size,
		Right:  size,
		Bottom: size,
		Left:   size,
	}
}

// PaddingXY creates padding with horizontal (x) and vertical (y) values.
// The x value applies to Left and Right, y applies to Top and Bottom.
func PaddingXY(x, y int) Padding {
	return Padding{
		Top:    y,
		Right:  x,
		Bottom: y,
		Left:   x,
	}
}

// PaddingTRBL creates padding with explicit values for each side,
// following CSS convention: Top, Right, Bottom, Left (clockwise from top).
func PaddingTRBL(top, right, bottom, left int) Padding {
	return Padding{
		Top:    top,
		Right:  right,
		Bottom: bottom,
		Left:   left,
	}
}

// PaddingNone returns zero padding on all sides.
func PaddingNone() Padding {
	return Padding{}
}

// Horizontal returns the total horizontal padding (Left + Right).
func (p Padding) Horizontal() int {
	return p.Left + p.Right
}

// Vertical returns the total vertical padding (Top + Bottom).
func (p Padding) Vertical() int {
	return p.Top + p.Bottom
}
