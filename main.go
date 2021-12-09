package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"os"
	"strings"
	//"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/term"
)

const (
	width = 96
)

var (
	subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}

	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 0).
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

func main() {
	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	doc := strings.Builder{}

	{
		okButton := buttonStyle.Render("Ok")

		title := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render("GLYPHS")
		buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton)
		ui := lipgloss.JoinVertical(lipgloss.Center, title, buttons)

		dialog := lipgloss.Place(width, 9,
			lipgloss.Center, lipgloss.Center,
			dialogBoxStyle.Render(ui),
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
