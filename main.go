package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	width = 96
)

type Glyphs struct {
	glyphs []string
}

type Chars struct {
	basicAccented   []Glyphs
	basicLatin      []Glyphs
	latinSupplement []Glyphs
}

type model struct {
	cursor   int
	sections []string
	selected int
}

func initialModel() model {
	return model{
		sections: []string{"Basic Accented", "Basic Latin", "Latin-1 Supplement"},
		selected: 0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

var (
	uppercase = []string{"Á,À,Â,Ã,Ä,Å", "Ç,Č,Ć", "Đ", "É,È,Ê,Ë", "Í,Ì,Î,Ï", "Ñ", "Ó,Ò,Ô,Õ,Ö", "Ú,Ù,Û,Ü"}
	lowercase = []string{"á,à,â,ã,ä,å", "ç,č,ć", "đ", "é,è,ê,ë", "í,ì,î,ï", "ñ", "ó,ò,ô,õ,ö", "ú,ù,û,ü"}
	symbols   = []string{"$,₱,€,¥,£,¢", "¡,¿", "“", "°", "•", "‰", "©,®", "‹,›,×, «,»", "Æ,Œ,æ,œ,ß,§"}
	latinx    = []string{"Ø,Ý,Ÿ,Š,Ž", "ø,ý,ÿ,š,ž"}

	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	// Tabs.

	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	tabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	tab = lipgloss.NewStyle().
		Border(tabBorder, true).
		BorderForeground(highlight).
		Padding(0, 1)

	activeTab = tab.Copy().Border(activeTabBorder, true)

	tabGap = tab.Copy().
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)

	subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}

	titleBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"}).
			Align(lipgloss.Center)

	titleStyle = lipgloss.NewStyle().
			Inherit(titleBarStyle).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#FF5F87")).
			Padding(0, 1)

	textStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true, true, true, true)

	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(0, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)
	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#888B7E")).
			Padding(0, 3).
			MarginTop(1)

	docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)
)

func renderGlyphs(glyphs []string) []string {
	formattedGlyphs := make([]string, 0)
	for i := 0; i < len(glyphs); i++ {
		formattedGlyphs = append(formattedGlyphs, textStyle.Render(glyphs[i]))
	}
	return formattedGlyphs
}

func (m model) View() string {
	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	doc := strings.Builder{}

	// Tabs
	// get model.selected to determine activeTab
	{
		row := lipgloss.JoinHorizontal(
			lipgloss.Top,
			activeTab.Render("Basic Accented"),
			tab.Render("Basic Latin"),
			tab.Render("Latin-1 Supplement"),
		)
		gap := tabGap.Render(strings.Repeat(" ", max(0, width-lipgloss.Width(row)-2)))
		row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
		doc.WriteString(row + "\n\n")
	}

	{
		title := titleStyle.Width(physicalWidth / 3).Render("GLYPHS")

		upcase := Glyphs{glyphs: uppercase}
		downcase := Glyphs{glyphs: lowercase}
		syms := Glyphs{glyphs: symbols}
		latin := Glyphs{glyphs: latinx}
		basicSet := []Glyphs{upcase, downcase, syms, latin}
		charset := Chars{basicSet, nil, nil}
		rows := make([]string, len(charset.basicAccented))

		for i := 0; i < len(rows); i++ {
			rows[i] = lipgloss.JoinHorizontal(lipgloss.Center, renderGlyphs(charset.basicAccented[i].glyphs)...)
		}
		//fmt.Printf("chars %v %d", charset.basicAccented[1].glyphs, len(rows))

		okButton := buttonStyle.Render("Ok")
		body := lipgloss.JoinVertical(lipgloss.Center, rows...)
		buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton)
		titleUi := lipgloss.JoinHorizontal(lipgloss.Center, title)
		ui := lipgloss.JoinVertical(lipgloss.Center, body, buttons)
		view := lipgloss.JoinVertical(lipgloss.Center, titleUi, ui)

		dialog := lipgloss.Place(width, 9,
			lipgloss.Center, lipgloss.Center,
			dialogBoxStyle.Render(view),
			lipgloss.WithWhitespaceChars("猫咪"),
			lipgloss.WithWhitespaceForeground(subtle),
		)
		doc.WriteString(dialog + "\n\n")
	}

	if physicalWidth > 0 {
		docStyle = docStyle.MaxWidth(physicalWidth)
	}

	return docStyle.Render(doc.String())
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("セバエラー")
		os.Exit(1)
	}
}
