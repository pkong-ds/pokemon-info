package main

import (
	"fmt"
	"os"
	"strings"
)

func fetchMove(m *MoveEntry) {
	if m == nil {
		fmt.Fprintln(os.Stderr, "Error: MoveEntry is nil")
		os.Exit(1)
	}
	fmt.Printf("Fetching details for %s from %s...\n\n", m.Name, m.URL)

	var moveData MoveAPIData
	if err := fetchAndUnmarshal(m.URL, &moveData); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON response: %v\n", err)
		os.Exit(1)
	}

	// Header
	fmt.Printf("--- %s (ID: %d) ---\n", strings.Title(moveData.Name), moveData.ID)

	// Core stats
	fmt.Printf("Type: %s\n", strings.Title(moveData.Type.Name))
	fmt.Printf("Damage Class: %s\n", strings.Title(moveData.DamageClass.Name))
	fmt.Printf("Power: %s\n", intPtrToString(moveData.Power))
	fmt.Printf("Accuracy: %s\n", intPtrToString(moveData.Accuracy))
	fmt.Printf("PP: %s\n", intPtrToString(moveData.PP))
	fmt.Printf("Priority: %d\n", moveData.Priority)
	fmt.Printf("Generation: %s\n", strings.Title(strings.ReplaceAll(moveData.Generation.Name, "-", " ")))
	fmt.Printf("Target: %s\n", strings.Title(strings.ReplaceAll(moveData.Target.Name, "-", " ")))

	// Effect
	printMoveEffect(&moveData)

	// Meta
	printMoveMeta(&moveData)

	// Stat changes
	printMoveStatChanges(&moveData)

	// Flavor text
	printMoveFlavorText(&moveData)

	// Learned by count
	fmt.Printf("\nLearned By: %d Pokémon\n", len(moveData.LearnedByPokemon))

	// Machines
	printMoveMachines(&moveData)
}

func intPtrToString(p *int) string {
	if p == nil {
		return "-"
	}
	return fmt.Sprintf("%d", *p)
}

func printMoveEffect(m *MoveAPIData) {
	for _, entry := range m.EffectEntries {
		if entry.Language.Name == "en" {
			effect := entry.Effect
			shortEffect := entry.ShortEffect

			// Replace $effect_chance placeholder with actual value
			if m.EffectChance != nil {
				chance := fmt.Sprintf("%d", *m.EffectChance)
				effect = strings.ReplaceAll(effect, "$effect_chance", chance)
				shortEffect = strings.ReplaceAll(shortEffect, "$effect_chance", chance)
			}

			fmt.Printf("\nEffect: %s\n", strings.TrimSpace(effect))
			fmt.Printf("  (short: %s)\n", strings.TrimSpace(shortEffect))
			return
		}
	}
}

func printMoveMeta(m *MoveAPIData) {
	if m.Meta == nil {
		return
	}
	meta := m.Meta

	fmt.Println("\nMeta:")
	fmt.Printf("  Category: %s\n", strings.Title(strings.ReplaceAll(meta.Category.Name, "-", " ")))

	if meta.Ailment.Name != "none" {
		fmt.Printf("  Ailment: %s (%d%% chance)\n", strings.Title(meta.Ailment.Name), meta.AilmentChance)
	}

	if meta.MinHits != nil && meta.MaxHits != nil {
		fmt.Printf("  Hits: %d-%d\n", *meta.MinHits, *meta.MaxHits)
	}
	if meta.MinTurns != nil && meta.MaxTurns != nil {
		fmt.Printf("  Turns: %d-%d\n", *meta.MinTurns, *meta.MaxTurns)
	}

	fmt.Printf("  Crit Rate: %d | Drain: %d%% | Healing: %d%%\n", meta.CritRate, meta.Drain, meta.Healing)
	fmt.Printf("  Flinch Chance: %d%% | Stat Chance: %d%%\n", meta.FlinchChance, meta.StatChance)
}

func printMoveStatChanges(m *MoveAPIData) {
	if len(m.StatChanges) == 0 {
		return
	}
	fmt.Println("\nStat Changes:")
	for _, sc := range m.StatChanges {
		sign := "+"
		if sc.Change < 0 {
			sign = ""
		}
		fmt.Printf("  %s: %s%d\n", strings.Title(strings.ReplaceAll(sc.Stat.Name, "-", " ")), sign, sc.Change)
	}
}

func printMoveFlavorText(m *MoveAPIData) {
	// Find the latest English flavor text (last one in the list)
	var latest *MoveFlavorText
	for i := len(m.FlavorTextEntries) - 1; i >= 0; i-- {
		if m.FlavorTextEntries[i].Language.Name == "en" {
			latest = &m.FlavorTextEntries[i]
			break
		}
	}
	if latest == nil {
		return
	}

	// Clean up flavor text (API sometimes has newlines/form feeds in the text)
	text := strings.ReplaceAll(latest.FlavorText, "\n", " ")
	text = strings.ReplaceAll(text, "\f", " ")
	text = strings.Join(strings.Fields(text), " ")

	fmt.Printf("\nFlavor Text (%s):\n", strings.Title(strings.ReplaceAll(latest.VersionGroup.Name, "-", " ")))
	fmt.Printf("  \"%s\"\n", text)
}

func printMoveMachines(m *MoveAPIData) {
	if len(m.Machines) == 0 {
		return
	}
	fmt.Println("\nMachines:")
	for _, mv := range m.Machines {
		fmt.Printf("  - %s\n", strings.Title(strings.ReplaceAll(mv.VersionGroup.Name, "-", " ")))
	}
}
