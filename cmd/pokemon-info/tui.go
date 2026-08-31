package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type page int

const (
	pageLanding page = iota
	pagePokemon
	pageMoves
	pageHelp
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffd75e"))
	dimStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	selStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffd75e"))
	nameStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffffff"))
	labelStyle = lipgloss.NewStyle().Bold(true)
	errStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
	barStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffd75e"))
	tagStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("0")).Padding(0, 1)
	padStyle   = lipgloss.NewStyle().Padding(1, 2)
)

type slashCommand struct {
	name string
	desc string
	page page
	quit bool
}

var slashCommands = []slashCommand{
	{name: "/pokemons", desc: "browse every Pokémon", page: pagePokemon},
	{name: "/moves", desc: "browse every move", page: pageMoves},
	{name: "/help", desc: "keybindings and about", page: pageHelp},
	{name: "/quit", desc: "exit", quit: true},
}

type rootModel struct {
	page    page
	width   int
	height  int
	input   textinput.Model
	sel     int
	pokemon *browserModel
	moves   *browserModel
}

func runTUI() error {
	ti := textinput.New()
	ti.Placeholder = "/"
	ti.Prompt = "› "
	ti.PromptStyle = selStyle
	ti.CharLimit = 32

	m := rootModel{input: ti, page: pagePokemon}
	m.pokemon = newBrowser(kindPokemon, m.width, m.height)
	_, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	return err
}

func (m rootModel) Init() tea.Cmd {
	if m.page == pagePokemon {
		return tea.Batch(m.input.Focus(), m.pokemon.enter())
	}
	return m.input.Focus()
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.pokemon != nil {
			m.pokemon.resize(msg.Width, msg.Height)
		}
		if m.moves != nil {
			m.moves.resize(msg.Width, msg.Height)
		}
		return m, nil
	case tea.KeyMsg:
		switch m.page {
		case pageLanding:
			return m.updateLanding(msg)
		case pageHelp:
			if msg.Type == tea.KeyCtrlC {
				return m, tea.Quit
			}
			m.page = pageLanding
			return m, m.input.Focus()
		default:
			cmd, act := m.browser().Update(msg)
			switch act {
			case actQuit:
				return m, tea.Quit
			case actBack:
				m.page = pageLanding
				return m, m.input.Focus()
			}
			return m, cmd
		}
	default:
		if m.page == pagePokemon || m.page == pageMoves {
			cmd, _ := m.browser().Update(msg)
			return m, cmd
		}
		return m, nil
	}
}

func (m *rootModel) browser() *browserModel {
	if m.page == pageMoves {
		return m.moves
	}
	return m.pokemon
}

func (m rootModel) updateLanding(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case msg.Type == tea.KeyCtrlC:
		return m, tea.Quit
	case m.input.Value() == "" && (msg.String() == "q" || msg.String() == "esc"):
		return m, tea.Quit
	case msg.String() == "esc":
		m.input.SetValue("")
		m.sel = 0
		return m, nil
	case msg.String() == "up":
		if m.sel > 0 {
			m.sel--
		}
		return m, nil
	case msg.String() == "down":
		if m.sel < len(m.filteredCommands())-1 {
			m.sel++
		}
		return m, nil
	case msg.String() == "enter":
		cmds := m.filteredCommands()
		if len(cmds) == 0 {
			return m, nil
		}
		if m.sel >= len(cmds) {
			m.sel = len(cmds) - 1
		}
		return m.execCommand(cmds[m.sel])
	default:
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		m.sel = 0
		return m, cmd
	}
}

// filteredCommands returns the slash commands matching the current input.
// The leading "/" is optional when typing.
func (m rootModel) filteredCommands() []slashCommand {
	value := strings.TrimPrefix(strings.ToLower(m.input.Value()), "/")
	if value == "" {
		return slashCommands
	}
	var cmds []slashCommand
	for _, c := range slashCommands {
		if strings.HasPrefix(strings.TrimPrefix(c.name, "/"), value) {
			cmds = append(cmds, c)
		}
	}
	return cmds
}

func (m rootModel) execCommand(c slashCommand) (tea.Model, tea.Cmd) {
	if c.quit {
		return m, tea.Quit
	}
	m.input.SetValue("")
	m.sel = 0

	switch c.page {
	case pageHelp:
		m.page = pageHelp
		return m, nil
	case pagePokemon:
		if m.pokemon == nil {
			m.pokemon = newBrowser(kindPokemon, m.width, m.height)
		}
		m.page = pagePokemon
		return m, m.pokemon.enter()
	case pageMoves:
		if m.moves == nil {
			m.moves = newBrowser(kindMoves, m.width, m.height)
		}
		m.page = pageMoves
		return m, m.moves.enter()
	}
	return m, nil
}

func (m rootModel) View() string {
	switch m.page {
	case pageHelp:
		return m.viewHelp()
	case pagePokemon:
		return m.pokemon.View()
	case pageMoves:
		return m.moves.View()
	default:
		return m.viewLanding()
	}
}

func (m rootModel) viewLanding() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("pokemon-info"))
	b.WriteString(dimStyle.Render("  · Pokédex in your terminal"))
	b.WriteString("\n\n")

	cmds := m.filteredCommands()
	for i, c := range cmds {
		style := dimStyle
		if i == m.sel {
			style = selStyle
		}
		b.WriteString("  ")
		b.WriteString(style.Render(c.name))
		b.WriteString(dimStyle.Render("  " + c.desc))
		b.WriteString("\n")
	}

	b.WriteString("\n  ")
	b.WriteString(m.input.View())
	b.WriteString("\n\n")
	b.WriteString(dimStyle.Render(fmt.Sprintf("  %d Pokémon · %d moves embedded · details fetched live from PokeAPI",
		len(allPokemonEntries), len(allMoveEntries))))
	b.WriteString("\n")
	b.WriteString(dimStyle.Render("  enter run · ↑↓ pick · esc clear · q quit"))
	return padStyle.Render(b.String())
}

func (m rootModel) viewHelp() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Keybindings"))
	b.WriteString("\n\n")
	b.WriteString("  ↑↓ / j k    move selection\n")
	b.WriteString("  /           filter the list\n")
	b.WriteString("  esc         clear filter, or back to this menu\n")
	b.WriteString("  enter       reload details for selection\n")
	b.WriteString("  f           cycle Pokémon art: small → large → off\n")
	b.WriteString("  q           quit\n\n")
	b.WriteString(titleStyle.Render("About"))
	b.WriteString("\n\n")
	b.WriteString("  Data comes live from PokeAPI; the name catalog is\n")
	b.WriteString("  embedded so filtering works offline. TUI inspired by\n")
	b.WriteString("  poketex (github.com/ckaznable/poketex); Pokémon art from\n")
	b.WriteString("  pokemon-colorscripts (phoneybadber) and poketex's gen 9\n")
	b.WriteString("  pack (Caruban).\n\n")
	b.WriteString(dimStyle.Render("  esc back"))
	return padStyle.Render(b.String())
}
