package table

import (
	"bytes"
	"fmt"
)

const PADDING = 1
const EMPTY_MISSING_VAL = 0

type Alignment int

const (
	LEFT Alignment = iota
	CENTER
	RIGHT
)

const (
	SINGLE_HORIZONTAL_DIV rune = '\u2500' // ─
	SINGLE_VERTICAL_DIV   rune = '\u2502' // │
	SINGLE_CROSS_DIV      rune = '\u253C' // ┼

	DOUBLE_HORIZONTAL_DIV rune = '\u2550' // ═
	DOUBLE_VERTICAL_DIV   rune = '\u2551' // ║
	DOUBLE_CROSS_DIV      rune = '\u256C' // ╬
)

type Table struct {
	columns       map[string]*column
	columnOrder   []string
	VerticalDiv   rune
	HorizontalDiv rune
	CrossDiv      rune
	numEntries    int
}

type column struct {
	header       string
	alignment    Alignment
	missingVal   rune
	entries      map[int]string
	maxEntrySize int
}

/**
 * Adds a column to the table.
 *
 * @note Columns are added from left -> Right
 *
 * @param name The Header name when printed
 * @alignment Column alignment of Values such as CENTER aligned
 * @missingVal Character to print if entry is not supplied. If 0 nothing will be printed.
 *
 * @return An error will be returned if the table already has a column with the header name already present.
 */
func (t *Table) CreateColumn(name string, alignment Alignment, missingVal rune) error {
	if t.columns == nil {
		t.columns = make(map[string]*column)
		t.columnOrder = make([]string, 0)
	}

	if _, hasCol := t.columns[name]; hasCol {
		return fmt.Errorf("column with name '%s' already exists", name)
	}

	col := &column{
		header:       name,
		alignment:    alignment,
		missingVal:   missingVal,
		entries:      make(map[int]string),
		maxEntrySize: len(name),
	}

	t.columns[col.header] = col
	t.columnOrder = append(t.columnOrder, col.header)

	return nil
}

/**
 * Adds a row to the table.
 *
 * @note Rows are added from Top to Bottom
 *
 * @param row Contents of the row. Keys = Header Name, Values = Column Value
 *
 * @return An error is returned if one or more Keys refer to unspecified columns
 */
func (t *Table) AddEntry(row map[string]any) error {
	t.numEntries++

	var err error
	var missingHdrs []string

	for hdr, val := range row {
		if col, has := t.columns[hdr]; has {
			valStr := fmt.Sprint(val)
			col.entries[t.numEntries] = valStr
			if newMaxLen := len(valStr); col.maxEntrySize < newMaxLen {
				col.maxEntrySize = newMaxLen
			}
		} else {
			missingHdrs = append(missingHdrs, hdr)
		}
	}
	if 0 < len(missingHdrs) {
		err = fmt.Errorf("no column(s) with header name(s) exist: %v", missingHdrs)
	}
	return err
}

func (t *Table) String() string {
	buf := bytes.NewBufferString("")

	if t.HorizontalDiv == 0 {
		t.HorizontalDiv = SINGLE_HORIZONTAL_DIV
	}
	if t.VerticalDiv == 0 {
		t.VerticalDiv = SINGLE_VERTICAL_DIV
	}
	if t.CrossDiv == 0 {
		t.CrossDiv = SINGLE_CROSS_DIV
	}

	for i := -1; i <= t.numEntries; i++ {
		for pos, colName := range t.columnOrder {
			col := t.columns[colName]
			switch i {
			case -1:
				header := formatColEntry(col.header, col.alignment, col.maxEntrySize)
				buf.WriteString(pad(header, PADDING))
			case 0:
				var hdrDiv string
				for i := 0; i < col.maxEntrySize+(PADDING*2); i++ {
					hdrDiv += string(t.HorizontalDiv)
				}
				buf.WriteString(string(hdrDiv))
			default:
				entry, has := col.entries[i]
				if !has {
					if col.missingVal == EMPTY_MISSING_VAL {
						entry = ""
					} else {
						entry = string(col.missingVal)
					}
				}

				row := formatColEntry(entry, col.alignment, col.maxEntrySize)
				buf.WriteString(pad(row, PADDING))
			}

			if pos+1 != len(t.columnOrder) {
				if i == 0 {
					buf.WriteString(string(t.CrossDiv))
				} else {
					buf.WriteString(string(t.VerticalDiv))
				}
			} else {
				buf.WriteString("\n\r")
			}
		}
	}

	return buf.String()
}

/**
 * Aligns a value to a certain width.
 *
 * @param val Value to align to width.
 * @param alignment How the value should be aligned
 * @param width Size of the width to align to
 *
 * @example formatColEntry("data", LEFT,   8) = "data    "
 * @example formatColEntry("data", CENTER, 8) = "  data  "
 * @example formatColEntry("data", RIGHT,  8) = "    data"
 *
 * @return The aligned value
 */
func formatColEntry(val string, alignment Alignment, width int) string {
	var valFmt string
	switch alignment {
	case RIGHT:
		valFmt = fmt.Sprintf("%%%ds%%s", width)
	case CENTER:
		valWidth := len(val)
		left := (width - valWidth) / 2
		right := width - valWidth - left
		valFmt = fmt.Sprintf("%%%ds%%%ds", valWidth+left, right)
	default:
		valFmt = fmt.Sprintf("%%-%ds%%s", width)
	}
	return fmt.Sprintf(valFmt, val, "")
}

/**
 * Pads a string value with N number of spaces on the left and right.
 *
 * @param val Value to add padding to
 * @param padding How much padding to add on the left and right
 *
 * @example pad("data", 0) = "data"
 * @example pad("data", 1) = " data "
 * @example pad("data", 2) = "  data  "
 *
 * @return The value with padding
 */
func pad(val string, padding int) string {
	valFmt := fmt.Sprintf("%%%ds%%s%%%ds", padding, padding)
	return fmt.Sprintf(valFmt, "", val, "")
}
