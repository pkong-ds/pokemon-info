package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type browseKind int

const (
	kindPokemon browseKind = iota
	kindMoves
)

type browserAction int

const (
	actNone browserAction = iota
	actBack
	actQuit
)

type artMode int

const (
	artSmall artMode = iota
	artLarge
	artOff
)

type browserItem struct {
	slug string
	name string
	id   int
}

func (i browserItem) Title() string       { return i.name }
func (i browserItem) Description() string { return fmt.Sprintf("#%04d", i.id) }
func (i browserItem) FilterValue() string { return i.name }

type detailResult struct {
	slug string
	val  any
	err  error
}

type selectionTick struct{ idx int }

type browserModel struct {
	kind     browseKind
	list     list.Model
	spinner  spinner.Model
	entryOf  map[string]any
	cache    map[string]any
	fetching map[string]bool
	lastIdx  int
	lastSlug string
	width    int
	height   int
	artMode  artMode
}

func listWidth(w int) int {
	lw := w * 38 / 100
	if lw < 26 {
		lw = 26
	}
	if lw > 46 {
		lw = 46
	}
	return lw
}

func listHeight(h int) int {
	return h - 1
}

func newBrowser(kind browseKind, width, height int) *browserModel {
	if width <= 0 {
		width = 80
	}
	if height <= 0 {
		height = 24
	}

	b := &browserModel{
		kind:     kind,
		spinner:  spinner.New(spinner.WithSpinner(spinner.Dot), spinner.WithStyle(dimStyle)),
		entryOf:  map[string]any{},
		cache:    map[string]any{},
		fetching: map[string]bool{},
		width:    width,
		height:   height,
	}

	var items []list.Item
	if kind == kindPokemon {
		for i := range allPokemonEntries {
			e := &allPokemonEntries[i]
			item := browserItem{slug: e.Slug, name: e.Name, id: e.ID}
			b.entryOf[e.Slug] = e
			items = append(items, item)
		}
		b.list = list.New(items, list.NewDefaultDelegate(), listWidth(width), listHeight(height))
		b.list.Title = fmt.Sprintf("Pokémon (%d)", len(items))
	} else {
		for i := range allMoveEntries {
			e := &allMoveEntries[i]
			item := browserItem{slug: e.Slug, name: e.Name, id: e.ID}
			b.entryOf[e.Slug] = e
			items = append(items, item)
		}
		b.list = list.New(items, list.NewDefaultDelegate(), listWidth(width), listHeight(height))
		b.list.Title = fmt.Sprintf("Moves (%d)", len(items))
	}

	b.list.SetShowStatusBar(false)
	b.list.SetShowHelp(false)
	return b
}

func (b *browserModel) resize(width, height int) {
	if width <= 0 || height <= 0 {
		return
	}
	b.width = width
	b.height = height
	b.list.SetSize(listWidth(width), listHeight(height))
}

func (b *browserModel) selectedItem() browserItem {
	it, ok := b.list.SelectedItem().(browserItem)
	if !ok {
		return browserItem{}
	}
	return it
}

func (b *browserModel) enter() tea.Cmd {
	item := b.selectedItem()
	if item.slug == "" {
		return nil
	}
	return tea.Batch(b.fetchCmd(item), spinner.Tick)
}

func (b *browserModel) fetchCmd(item browserItem) tea.Cmd {
	b.fetching[item.slug] = true
	return func() tea.Msg {
		var val any
		var err error
		if b.kind == kindPokemon {
			val, err = fetchPokemonDetail(b.entryOf[item.slug].(*PokemonEntry))
		} else {
			val, err = fetchMoveDetail(b.entryOf[item.slug].(*MoveEntry))
		}
		return detailResult{slug: item.slug, val: val, err: err}
	}
}

func (b *browserModel) Update(msg tea.Msg) (tea.Cmd, browserAction) {
	switch msg := msg.(type) {
	case detailResult:
		b.fetching[msg.slug] = false
		if msg.err != nil {
			b.cache[msg.slug] = msg.err
		} else {
			b.cache[msg.slug] = msg.val
		}
		var cmd tea.Cmd
		if len(b.fetching) > 0 {
			cmd = spinner.Tick
		}
		return cmd, actNone

	case spinner.TickMsg:
		newSpinner, tick := b.spinner.Update(msg)
		b.spinner = newSpinner
		var cmd tea.Cmd
		if len(b.fetching) > 0 {
			cmd = tick
		}
		return cmd, actNone

	case selectionTick:
		if msg.idx != b.list.Index() {
			return nil, actNone
		}
		item := b.selectedItem()
		if item.slug == "" || b.fetching[item.slug] || b.cache[item.slug] != nil {
			return nil, actNone
		}
		return tea.Batch(b.fetchCmd(item), spinner.Tick), actNone

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return nil, actQuit
		}
		if b.list.FilterState() == list.Filtering {
			return b.forward(msg), actNone
		}
		switch msg.String() {
		case "esc":
			if b.list.FilterState() == list.FilterApplied {
				return b.forward(msg), actNone
			}
			return nil, actBack
		case "q":
			return nil, actQuit
		case "f":
			if b.kind == kindPokemon {
				b.artMode = (b.artMode + 1) % 3
			}
			return nil, actNone
		case "enter":
			item := b.selectedItem()
			if item.slug == "" {
				return nil, actNone
			}
			return tea.Batch(b.fetchCmd(item), spinner.Tick), actNone
		default:
			return b.forward(msg), actNone
		}
	}
	// Forward everything else to the list: FilterMatchesMsg and friends
	// only take effect when the list's Update sees them.
	return b.forward(msg), actNone
}

// forward passes msg to the list and schedules a debounced detail fetch
// when the selection changed. Selection is compared by index AND slug:
// filtering can narrow the list around the same index while replacing
// the selected item.
func (b *browserModel) forward(msg tea.Msg) tea.Cmd {
	newList, cmd := b.list.Update(msg)
	b.list = newList
	if idx := b.list.Index(); idx != b.lastIdx || b.selectedItem().slug != b.lastSlug {
		b.lastIdx = idx
		b.lastSlug = b.selectedItem().slug
		wait := tea.Tick(250*time.Millisecond, func(time.Time) tea.Msg {
			return selectionTick{idx}
		})
		cmd = tea.Batch(cmd, wait)
	}
	return cmd
}

func (b *browserModel) View() string {
	wList := listWidth(b.width)
	wDetail := b.width - wList - 2
	if wDetail < 30 {
		wDetail = 30
	}
	left := b.list.View()
	right := b.detailView(wDetail)
	keys := "↑↓ move · / filter · enter reload · esc menu · q quit"
	if b.kind == kindPokemon {
		keys = "↑↓ move · / filter · enter reload · f art · esc menu · q quit"
	}
	footer := dimStyle.Render("  " + keys)
	return lipgloss.JoinHorizontal(lipgloss.Top, left, right) + "\n" + footer
}

func (b *browserModel) detailView(w int) string {
	item := b.selectedItem()
	var title, body string
	if item.slug == "" {
		title = dimStyle.Render("No selection")
		body = dimStyle.Render("nothing selected")
	} else {
		title = nameStyle.Render(item.name) + dimStyle.Render(fmt.Sprintf("  #%d", item.id))
		body = b.detailBody(item, w-4)
	}

	content := title + "\n\n" + clip(body, b.height-5)
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("238")).
		Padding(0, 1).
		Width(w).
		Height(b.height - 3)
	return box.Render(content)
}

func (b *browserModel) detailBody(item browserItem, w int) string {
	if v, ok := b.cache[item.slug]; ok {
		if err, isErr := v.(error); isErr {
			return errStyle.Render(err.Error()) + "\n\n" + dimStyle.Render("press enter to retry")
		}
		if b.kind == kindPokemon {
			if d, ok := v.(*pokemonDetail); ok {
				body := renderPokemon(d, w)
				if art := b.art(item.slug, w); art != "" {
					body = art + "\n\n" + body
				}
				return body
			}
		} else {
			if d, ok := v.(*MoveAPIData); ok {
				return renderMove(d, w)
			}
		}
	}
	if b.fetching[item.slug] {
		return b.spinner.View() + " " + dimStyle.Render("loading…")
	}
	return dimStyle.Render("loading details…")
}

// art returns the colorscript for the selected Pokémon, honoring the current
// art mode and falling back to the small set when the large art would not
// fit the detail pane.
func (b *browserModel) art(slug string, w int) string {
	if b.artMode == artOff {
		return ""
	}
	set := "small"
	if b.artMode == artLarge {
		set = "large"
	}
	a, ok := findArt(slug, set)
	if !ok {
		return ""
	}
	if a.width+2 > w {
		if set == "large" {
			if s, ok := findArt(slug, "small"); ok && s.width+2 <= w {
				return s.content
			}
		}
		return ""
	}
	return a.content
}

func clip(s string, maxLines int) string {
	if maxLines < 3 {
		maxLines = 3
	}
	lines := strings.Split(s, "\n")
	if len(lines) <= maxLines {
		return s
	}
	return strings.Join(lines[:maxLines], "\n")
}
