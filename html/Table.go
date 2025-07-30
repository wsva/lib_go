package html

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Table comment
type Table struct {
	Title       string     `json:"Title"`
	Head        []string   `json:"Head"`
	RowList     [][]string `json:"Rows"`
	ColumnWidth []int      `json:"ColumnWidth"`
}

// SetTitle comment
func (t *Table) SetTitle(title string) {
	t.Title = title
}

// SetHead comment
func (t *Table) SetHead(head []string) {
	t.Head = head
}

// AddRow comment
func (t *Table) AddRow(row []string) {
	t.RowList = append(t.RowList, row)
}

func (t *Table) GetSize() (int, int) {
	maxColumn := len(t.Head)
	for _, v := range t.RowList {
		if maxColumn < len(v) {
			maxColumn = len(v)
		}
	}
	return len(t.RowList) + 1, maxColumn
}

func (t *Table) SetColumnWidth() {
	var width []int
	for _, v := range t.Head {
		stringWidth := (len(v) + utf8.RuneCountInString(v)) / 2
		width = append(width, stringWidth)
	}
	for _, v1 := range t.RowList {
		for k2, v2 := range v1 {
			stringWidth := (len(v2) + utf8.RuneCountInString(v2)) / 2
			if k2 > len(width) {
				width = append(width, stringWidth)
			}
			if stringWidth > width[k2] {
				width[k2] = stringWidth
			}
		}
	}
	t.ColumnWidth = width
}

func (t *Table) GetHTML(notitle bool) string {
	var html strings.Builder
	if !notitle {
		if t.Title != "" {
			html.WriteString(fmt.Sprintf("<p>%s</p>\n", t.Title))
		}
	}
	html.WriteString("<table>\n")
	html.WriteString("<tr>")
	for _, v1 := range t.Head {
		html.WriteString(fmt.Sprintf("<th>%s</th>", v1))
	}
	html.WriteString("</tr>\n")
	for _, v1 := range t.RowList {
		html.WriteString("<tr>")
		for _, v2 := range v1 {
			html.WriteString(fmt.Sprintf("<td>%s</td>", v2))
		}
		html.WriteString("</tr>\n")
	}
	html.WriteString("</table>\n")
	return html.String()
}
