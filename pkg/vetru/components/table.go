package components

import (
	"strings"

	"github.com/arthur404dev/dotts/pkg/vetru/theme"
	"github.com/charmbracelet/lipgloss"
)

// Table is a data table component with optional headers, borders, and striping.
type Table struct {
	theme       *theme.Theme
	headers     []string
	rows        [][]string
	colWidths   []int
	showBorders bool
	striped     bool
}

// NewTable creates a new Table component with borders enabled and striping disabled.
func NewTable(t *theme.Theme) *Table {
	return &Table{
		theme:       t,
		headers:     []string{},
		rows:        [][]string{},
		showBorders: true,
		striped:     false,
	}
}

// SetHeaders sets the table column headers.
func (t *Table) SetHeaders(headers ...string) *Table {
	t.headers = headers
	return t
}

// AddRow appends a row of cells to the table.
func (t *Table) AddRow(cells ...string) *Table {
	t.rows = append(t.rows, cells)
	return t
}

// SetRows replaces all existing rows with the provided ones.
func (t *Table) SetRows(rows [][]string) *Table {
	t.rows = rows
	return t
}

// Clear removes all rows from the table (headers are preserved).
func (t *Table) Clear() *Table {
	t.rows = [][]string{}
	return t
}

// SetColWidths sets explicit column widths (0 for auto-calculation).
func (t *Table) SetColWidths(widths ...int) *Table {
	t.colWidths = widths
	return t
}

// SetShowBorders enables or disables the separator line between header and rows.
func (t *Table) SetShowBorders(show bool) *Table {
	t.showBorders = show
	return t
}

// SetStriped enables or disables alternating row backgrounds.
func (t *Table) SetStriped(striped bool) *Table {
	t.striped = striped
	return t
}

// View renders the table as a formatted string.
func (tb *Table) View() string {
	if len(tb.headers) == 0 && len(tb.rows) == 0 {
		return ""
	}

	t := tb.theme

	numCols := len(tb.headers)
	for _, row := range tb.rows {
		if len(row) > numCols {
			numCols = len(row)
		}
	}

	widths := make([]int, numCols)

	for i, w := range tb.colWidths {
		if i < numCols {
			widths[i] = w
		}
	}

	for i, h := range tb.headers {
		if widths[i] == 0 && len(h) > widths[i] {
			widths[i] = len(h)
		}
	}

	for _, row := range tb.rows {
		for i, cell := range row {
			if i < numCols && widths[i] == 0 && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	for i := range widths {
		if widths[i] == 0 {
			widths[i] = 10
		}
		widths[i] += 2
	}

	var output []string

	if len(tb.headers) > 0 {
		var headerCells []string
		headerStyle := lipgloss.NewStyle().
			Foreground(t.Primary).
			Bold(true).
			Padding(0, 1)

		for i, h := range tb.headers {
			w := widths[i]
			cell := headerStyle.Width(w).Render(h)
			headerCells = append(headerCells, cell)
		}
		output = append(output, lipgloss.JoinHorizontal(lipgloss.Top, headerCells...))

		if tb.showBorders {
			sepStyle := lipgloss.NewStyle().Foreground(t.FgSubtle)
			var totalWidth int
			for _, w := range widths {
				totalWidth += w
			}
			output = append(output, sepStyle.Render(strings.Repeat("─", totalWidth)))
		}
	}

	for rowIdx, row := range tb.rows {
		var rowCells []string
		rowStyle := lipgloss.NewStyle().
			Foreground(t.FgBase).
			Padding(0, 1)

		if tb.striped && rowIdx%2 == 1 {
			rowStyle = rowStyle.Background(t.BgSubtle)
		}

		for i := 0; i < numCols; i++ {
			w := widths[i]
			cell := ""
			if i < len(row) {
				cell = row[i]
			}
			cellView := rowStyle.Width(w).Render(cell)
			rowCells = append(rowCells, cellView)
		}
		output = append(output, lipgloss.JoinHorizontal(lipgloss.Top, rowCells...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, output...)
}
