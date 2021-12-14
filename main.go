package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"os"
	"strings"
	"golang.org/x/term"
)

const (
	width = 96
)

var (
	uppercase = []string{"Á,À,Â,Ã,Ä,Å", "Ç", "É,È,Ê,Ë", "Í,Ì,Î,Ï", "Ñ", "Ó,Ò,Ô,Õ,Ö", "Ú,Ù,Û,Ü"}
	lowercase = []string{"á,à,â,ã,ä,å", "ç", "é,è,ê,ë", "í,ì,î,ï", "ñ", "ó,ò,ô,õ,ö", "ú,ù,û,ü"}
	symbols   = []string{"$,₱,€,¥,£,¢", "¡,¿", "“", "°", "•", "‰", "©,®", "‹,›,×, «,»", "Æ,Œ,æ,œ,ß,§"}
	latinx    = []string{"Ø,Ý,Ÿ,Š", "ø,ý,ÿ,š,ž"}

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

func main() {
	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	doc := strings.Builder{}

	{
		title := titleStyle.Width(physicalWidth / 3).Render("GLYPHS")

		allChars := [][]string{uppercase, lowercase, symbols, latinx}
		rows := make([]string, len(allChars))

		for i := 0; i < len(allChars); i++ {
			rows[i] = lipgloss.JoinHorizontal(lipgloss.Center, renderGlyphs(allChars[i])...)
		}

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

	fmt.Println(docStyle.Render(doc.String()))
}
