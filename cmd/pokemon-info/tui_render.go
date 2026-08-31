package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var typeColors = map[string]string{
	"normal":   "240",
	"fire":     "203",
	"water":    "39",
	"electric": "220",
	"grass":    "82",
	"ice":      "159",
	"fighting": "197",
	"poison":   "135",
	"ground":   "179",
	"flying":   "183",
	"psychic":  "212",
	"bug":      "114",
	"rock":     "137",
	"ghost":    "141",
	"dragon":   "99",
	"dark":     "95",
	"steel":    "250",
	"fairy":    "218",
}

func typeBadge(t string) string {
	color, ok := typeColors[strings.ToLower(t)]
	if !ok {
		color = "250"
	}
	return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color(color)).Padding(0, 1).Render(strings.Title(t))
}

func classBadge(c string) string {
	color := "250"
	switch strings.ToLower(c) {
	case "physical":
		color = "203"
	case "special":
		color = "39"
	case "status":
		color = "137"
	}
	return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color(color)).Padding(0, 1).Render(strings.Title(c))
}

func prettyName(s string) string {
	return strings.Title(strings.ReplaceAll(s, "-", " "))
}

func statName(apiName string) string {
	switch apiName {
	case "hp":
		return "HP"
	case "attack":
		return "Attack"
	case "defense":
		return "Defense"
	case "special-attack":
		return "Sp. Atk"
	case "special-defense":
		return "Sp. Def"
	case "speed":
		return "Speed"
	default:
		return prettyName(apiName)
	}
}

func englishGenus(genera []Genus) string {
	for _, g := range genera {
		if g.Language.Name == "en" {
			return g.Genus
		}
	}
	return ""
}

func speciesFlavor(s *PokemonSpeciesAPIData) string {
	if s == nil {
		return ""
	}
	var latest *FlavorTextEntry
	for i := len(s.FlavorTextEntries) - 1; i >= 0; i-- {
		if s.FlavorTextEntries[i].Language.Name == "en" {
			latest = &s.FlavorTextEntries[i]
			break
		}
	}
	if latest == nil {
		return ""
	}
	text := strings.ReplaceAll(latest.FlavorText, "\n", " ")
	text = strings.ReplaceAll(text, "\f", " ")
	return strings.Join(strings.Fields(text), " ")
}

func renderPokemon(d *pokemonDetail, w int) string {
	if d == nil || d.Pokemon == nil {
		return dimStyle.Render("no data")
	}
	p := d.Pokemon
	wrap := lipgloss.NewStyle().Width(w)

	var b strings.Builder
	b.WriteString(nameStyle.Render(strings.Title(p.Name)))
	b.WriteString(dimStyle.Render(fmt.Sprintf("  #%d", p.ID)))
	if d.Species != nil {
		if g := englishGenus(d.Species.Genera); g != "" {
			b.WriteString(dimStyle.Render("  " + g))
		}
		if d.Species.IsLegendary {
			b.WriteString("  " + tagStyle.Background(lipgloss.Color("220")).Render("legendary"))
		}
		if d.Species.IsMythical {
			b.WriteString("  " + tagStyle.Background(lipgloss.Color("212")).Render("mythical"))
		}
	}
	b.WriteString("\n\n")

	var badges []string
	for _, t := range p.Types {
		badges = append(badges, typeBadge(t.TypeInfo.Name))
	}
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, badges...))
	b.WriteString("\n\n")

	b.WriteString(labelStyle.Render("Info") + "\n")
	b.WriteString(fmt.Sprintf("  %-14s %.1f m\n", "Height", float64(p.Height)/10))
	b.WriteString(fmt.Sprintf("  %-14s %.1f kg\n", "Weight", float64(p.Weight)/10))
	b.WriteString(fmt.Sprintf("  %-14s %d\n", "Base Exp", p.BaseExperience))
	if d.Species != nil {
		b.WriteString(fmt.Sprintf("  %-14s %d\n", "Capture Rate", d.Species.CaptureRate))
	}
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Abilities") + "\n")
	for _, a := range p.Abilities {
		marker := ""
		if a.IsHidden {
			marker = dimStyle.Render(" (hidden)")
		}
		b.WriteString("  " + prettyName(a.Ability.Name) + marker + "\n")
	}
	b.WriteString("\n")

	b.WriteString(labelStyle.Render("Base Stats") + "\n")
	total := 0
	for _, s := range p.Stats {
		total += s.BaseStat
		b.WriteString(statBar(s.StatInfo.Name, s.BaseStat))
	}
	b.WriteString(statBar("total", total))
	b.WriteString("\n")

	if chain := evolutionChainString(d.Chain); chain != "" {
		b.WriteString(chain)
		b.WriteString("\n")
	}
	if ft := speciesFlavor(d.Species); ft != "" {
		b.WriteString(labelStyle.Render("Flavor Text") + "\n")
		b.WriteString(wrap.Render(dimStyle.Render(ft)))
		b.WriteString("\n")
	}
	return wrap.Render(b.String())
}

func statBar(apiName string, value int) string {
	const width = 18
	n := int(float64(value) / 255.0 * width)
	if n > width {
		n = width
	}
	bar := strings.Repeat("█", n) + strings.Repeat("░", width-n)
	return fmt.Sprintf("  %-9s %s %4d\n", statName(apiName), barStyle.Render(bar), value)
}

func evolutionChainString(c *PokemonEvolutionChainAPIData) string {
	if c == nil {
		return ""
	}
	var b strings.Builder
	b.WriteString(labelStyle.Render("Evolution Chain") + "\n")
	writeChainLink(&c.Chain, "  ", &b)
	return b.String()
}

func writeChainLink(l *ChainLink, indent string, b *strings.Builder) {
	b.WriteString(indent + "└─ " + strings.Title(l.Species.Name))
	conds := evolutionConditions(l.EvolutionDetails)
	if len(conds) > 0 {
		b.WriteString(dimStyle.Render("  (" + strings.Join(conds, ", ") + ")"))
	}
	b.WriteString("\n")
	for _, e := range l.EvolvesTo {
		writeChainLink(&e, indent+"  ", b)
	}
}

func evolutionConditions(details []EvolutionDetail) []string {
	var out []string
	for i := range details {
		out = append(out, evolutionConditionLines(&details[i])...)
	}
	return out
}

func evolutionConditionLines(e *EvolutionDetail) []string {
	var out []string
	add := func(s string) { out = append(out, s) }

	if e.MinLevel != nil {
		add(fmt.Sprintf("Lv. %d", *e.MinLevel))
	} else {
		add(prettyName(e.Trigger.Name))
	}
	if e.Item != nil {
		add("item: " + prettyName(e.Item.Name))
	}
	if e.HeldItem != nil {
		add("holding " + prettyName(e.HeldItem.Name))
	}
	if e.KnownMove != nil {
		add("knows " + prettyName(e.KnownMove.Name))
	}
	if e.KnownMoveType != nil {
		add("knows " + e.KnownMoveType.Name + " move")
	}
	if e.Location != nil {
		add("at " + prettyName(e.Location.Name))
	}
	if e.MinHappiness != nil {
		add(fmt.Sprintf("happiness %d", *e.MinHappiness))
	}
	if e.MinAffection != nil {
		add(fmt.Sprintf("affection %d", *e.MinAffection))
	}
	if e.MinBeauty != nil {
		add(fmt.Sprintf("beauty %d", *e.MinBeauty))
	}
	if e.TimeOfDay != "" {
		add(e.TimeOfDay)
	}
	if e.NeedsOverworldRain {
		add("in rain")
	}
	if e.TurnUpsideDown {
		add("upside down")
	}
	if e.Gender != nil {
		switch *e.Gender {
		case 1:
			add("female")
		case 2:
			add("male")
		}
	}
	if e.RelativePhysicalStats != nil {
		switch *e.RelativePhysicalStats {
		case -1:
			add("Atk < Def")
		case 0:
			add("Atk = Def")
		case 1:
			add("Atk > Def")
		}
	}
	if e.PartySpecies != nil {
		add(prettyName(e.PartySpecies.Name) + " in party")
	}
	if e.PartyType != nil {
		add(e.PartyType.Name + " type in party")
	}
	if e.TradeSpecies != nil {
		add("traded for " + prettyName(e.TradeSpecies.Name))
	}
	return out
}

func renderMove(m *MoveAPIData, w int) string {
	if m == nil {
		return dimStyle.Render("no data")
	}
	wrap := lipgloss.NewStyle().Width(w)

	var b strings.Builder
	b.WriteString(nameStyle.Render(strings.Title(m.Name)))
	b.WriteString(dimStyle.Render(fmt.Sprintf("  #%d", m.ID)))
	b.WriteString("\n\n")
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, typeBadge(m.Type.Name), classBadge(m.DamageClass.Name)))
	b.WriteString("\n\n")

	b.WriteString(labelStyle.Render("Stats") + "\n")
	b.WriteString(fmt.Sprintf("  %-10s %s\n", "Power", intPtrToString(m.Power)))
	b.WriteString(fmt.Sprintf("  %-10s %s\n", "Accuracy", intPtrToString(m.Accuracy)))
	b.WriteString(fmt.Sprintf("  %-10s %s\n", "PP", intPtrToString(m.PP)))
	b.WriteString(fmt.Sprintf("  %-10s %d\n", "Priority", m.Priority))
	b.WriteString(fmt.Sprintf("  %-10s %s\n", "Gen", prettyName(m.Generation.Name)))
	b.WriteString(fmt.Sprintf("  %-10s %s\n", "Target", prettyName(m.Target.Name)))
	b.WriteString("\n")

	if effect := moveEffect(m); effect != "" {
		b.WriteString(labelStyle.Render("Effect") + "\n")
		b.WriteString(wrap.Render(effect))
		b.WriteString("\n\n")
	}

	if meta := moveMeta(m); meta != "" {
		b.WriteString(labelStyle.Render("Meta") + "\n")
		b.WriteString(wrap.Render(meta))
		b.WriteString("\n\n")
	}

	if len(m.StatChanges) > 0 {
		b.WriteString(labelStyle.Render("Stat Changes") + "\n")
		for _, sc := range m.StatChanges {
			sign := "+"
			if sc.Change < 0 {
				sign = ""
			}
			b.WriteString(fmt.Sprintf("  %s: %s%d\n", prettyName(sc.Stat.Name), sign, sc.Change))
		}
		b.WriteString("\n")
	}

	if flavor := moveFlavor(m); flavor != "" {
		b.WriteString(labelStyle.Render("Flavor Text") + "\n")
		b.WriteString(wrap.Render(dimStyle.Render(flavor)))
		b.WriteString("\n\n")
	}

	b.WriteString(labelStyle.Render("Learned By") + "\n")
	b.WriteString(fmt.Sprintf("  %d Pokémon\n", len(m.LearnedByPokemon)))

	if len(m.Machines) > 0 {
		b.WriteString("\n" + labelStyle.Render("Machines") + "\n")
		for _, mv := range m.Machines {
			b.WriteString("  " + prettyName(mv.VersionGroup.Name) + "\n")
		}
	}
	return wrap.Render(b.String())
}

func moveEffect(m *MoveAPIData) string {
	for _, entry := range m.EffectEntries {
		if entry.Language.Name != "en" {
			continue
		}
		effect := strings.TrimSpace(entry.Effect)
		short := strings.TrimSpace(entry.ShortEffect)
		if m.EffectChance != nil {
			chance := fmt.Sprintf("%d", *m.EffectChance)
			effect = strings.ReplaceAll(effect, "$effect_chance", chance)
			short = strings.ReplaceAll(short, "$effect_chance", chance)
		}
		return effect + "\n" + dimStyle.Render("("+short+")")
	}
	return ""
}

func moveMeta(m *MoveAPIData) string {
	if m.Meta == nil {
		return ""
	}
	meta := m.Meta
	var b strings.Builder
	b.WriteString(fmt.Sprintf("  %-14s %s\n", "Category", prettyName(meta.Category.Name)))
	if meta.Ailment.Name != "none" {
		b.WriteString(fmt.Sprintf("  %-14s %s (%d%%)\n", "Ailment", prettyName(meta.Ailment.Name), meta.AilmentChance))
	}
	if meta.MinHits != nil && meta.MaxHits != nil {
		b.WriteString(fmt.Sprintf("  %-14s %d-%d\n", "Hits", *meta.MinHits, *meta.MaxHits))
	}
	if meta.MinTurns != nil && meta.MaxTurns != nil {
		b.WriteString(fmt.Sprintf("  %-14s %d-%d\n", "Turns", *meta.MinTurns, *meta.MaxTurns))
	}
	b.WriteString(fmt.Sprintf("  %-14s %d%%\n", "Drain", meta.Drain))
	b.WriteString(fmt.Sprintf("  %-14s %d%%\n", "Healing", meta.Healing))
	b.WriteString(fmt.Sprintf("  %-14s %d%%\n", "Crit Rate", meta.CritRate))
	b.WriteString(fmt.Sprintf("  %-14s %d%%\n", "Flinch", meta.FlinchChance))
	return b.String()
}

func moveFlavor(m *MoveAPIData) string {
	var latest *MoveFlavorText
	for i := len(m.FlavorTextEntries) - 1; i >= 0; i-- {
		if m.FlavorTextEntries[i].Language.Name == "en" {
			latest = &m.FlavorTextEntries[i]
			break
		}
	}
	if latest == nil {
		return ""
	}
	text := strings.ReplaceAll(latest.FlavorText, "\n", " ")
	text = strings.ReplaceAll(text, "\f", " ")
	return strings.Join(strings.Fields(text), " ")
}
