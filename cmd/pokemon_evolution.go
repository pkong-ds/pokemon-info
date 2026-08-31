package main

import (
	"fmt"
	"strings"
)

// PokemonEvolutionChainAPIData represents the top-level structure for evolution chain data
type PokemonEvolutionChainAPIData struct {
	ID              int               `json:"id"`
	BabyTriggerItem *NamedAPIResource `json:"baby_trigger_item"`
	Chain           ChainLink         `json:"chain"`
}

// ChainLink represents a link in the evolution chain
type ChainLink struct {
	IsBaby           bool              `json:"is_baby"`
	Species          NamedAPIResource  `json:"species"`
	EvolutionDetails []EvolutionDetail `json:"evolution_details"`
	EvolvesTo        []ChainLink       `json:"evolves_to"`
}

// EvolutionDetail represents the details of an evolution
type EvolutionDetail struct {
	Item                  *NamedAPIResource `json:"item"`
	Trigger               NamedAPIResource  `json:"trigger"`
	Gender                *int              `json:"gender"`
	HeldItem              *NamedAPIResource `json:"held_item"`
	KnownMove             *NamedAPIResource `json:"known_move"`
	KnownMoveType         *NamedAPIResource `json:"known_move_type"`
	Location              *NamedAPIResource `json:"location"`
	MinLevel              *int              `json:"min_level"`
	MinHappiness          *int              `json:"min_happiness"`
	MinBeauty             *int              `json:"min_beauty"`
	MinAffection          *int              `json:"min_affection"`
	NeedsOverworldRain    bool              `json:"needs_overworld_rain"`
	PartySpecies          *NamedAPIResource `json:"party_species"`
	PartyType             *NamedAPIResource `json:"party_type"`
	RelativePhysicalStats *int              `json:"relative_physical_stats"`
	TimeOfDay             string            `json:"time_of_day"`
	TradeSpecies          *NamedAPIResource `json:"trade_species"`
	TurnUpsideDown        bool              `json:"turn_upside_down"`
}

// PrintNonAPIInfo returns a string with all non-NamedAPIResource information
func (e *EvolutionDetail) PrintNonAPIInfo() string {
	var info strings.Builder

	info.WriteString("Evolution Details (Non-API Info):\n")

	if e.Gender != nil {
		info.WriteString(fmt.Sprintf("- Gender: %d\n", *e.Gender))
	}

	if e.MinLevel != nil {
		info.WriteString(fmt.Sprintf("- Min Level: %d\n", *e.MinLevel))
	}

	if e.MinHappiness != nil {
		info.WriteString(fmt.Sprintf("- Min Happiness: %d\n", *e.MinHappiness))
	}

	if e.MinBeauty != nil {
		info.WriteString(fmt.Sprintf("- Min Beauty: %d\n", *e.MinBeauty))
	}

	if e.MinAffection != nil {
		info.WriteString(fmt.Sprintf("- Min Affection: %d\n", *e.MinAffection))
	}

	if e.NeedsOverworldRain {
		info.WriteString("- Needs Overworld Rain: Yes\n")
	}

	if e.RelativePhysicalStats != nil {
		info.WriteString(fmt.Sprintf("- Relative Physical Stats: %d\n", *e.RelativePhysicalStats))
	}

	if e.TimeOfDay != "" {
		info.WriteString(fmt.Sprintf("- Time of Day: %s\n", e.TimeOfDay))
	}

	if e.TurnUpsideDown {
		info.WriteString("- Turn Upside Down: Yes\n")
	}

	return info.String()
}

// PrintEvolutionChain prints the entire evolution chain as a tree diagram
func (p *PokemonEvolutionChainAPIData) PrintEvolutionChain() {
	fmt.Printf("Evolution Chain #%d\n", p.ID)
	fmt.Println("└── " + p.Chain.Species.Name)

	// Print evolution details for the base species if any
	if len(p.Chain.EvolutionDetails) > 0 {
		fmt.Println("    Evolution Details:")
		for _, detail := range p.Chain.EvolutionDetails {
			printEvolutionDetail(detail, "    ")
		}
	}

	// Recursively print the rest of the chain
	for _, evolution := range p.Chain.EvolvesTo {
		printChainLink(evolution, "    ")
	}
}

// printChainLink is a helper function to recursively print each link in the chain
func printChainLink(link ChainLink, indent string) {
	// Print the species name with the branch symbol
	fmt.Println(indent + "└── " + link.Species.Name)

	// Print evolution details if any
	if len(link.EvolutionDetails) > 0 {
		fmt.Println(indent + "    Evolution Details:")
		for _, detail := range link.EvolutionDetails {
			printEvolutionDetail(detail, indent+"    ")
		}
	}

	// Recursively print the next evolutions
	for _, evolution := range link.EvolvesTo {
		printChainLink(evolution, indent+"    ")
	}
}

// printEvolutionDetail is a helper function to print evolution details
func printEvolutionDetail(detail EvolutionDetail, indent string) {
	// Print trigger
	fmt.Println(indent + "- Trigger: " + detail.Trigger.Name)

	// Print item if present
	if detail.Item != nil {
		fmt.Println(indent + "- Item: " + detail.Item.Name)
	}

	// Print level if present
	if detail.MinLevel != nil {
		fmt.Println(indent + fmt.Sprintf("- Min Level: %d", *detail.MinLevel))
	}

	// Print happiness if present
	if detail.MinHappiness != nil {
		fmt.Println(indent + fmt.Sprintf("- Min Happiness: %d", *detail.MinHappiness))
	}

	// Print beauty if present
	if detail.MinBeauty != nil {
		fmt.Println(indent + fmt.Sprintf("- Min Beauty: %d", *detail.MinBeauty))
	}

	// Print affection if present
	if detail.MinAffection != nil {
		fmt.Println(indent + fmt.Sprintf("- Min Affection: %d", *detail.MinAffection))
	}

	// Print held item if present
	if detail.HeldItem != nil {
		fmt.Println(indent + "- Held Item: " + detail.HeldItem.Name)
	}

	// Print known move if present
	if detail.KnownMove != nil {
		fmt.Println(indent + "- Known Move: " + detail.KnownMove.Name)
	}

	// Print known move type if present
	if detail.KnownMoveType != nil {
		fmt.Println(indent + "- Known Move Type: " + detail.KnownMoveType.Name)
	}

	// Print location if present
	if detail.Location != nil {
		fmt.Println(indent + "- Location: " + detail.Location.Name)
	}

	// Print time of day if present
	if detail.TimeOfDay != "" {
		fmt.Println(indent + "- Time of Day: " + detail.TimeOfDay)
	}

	// Print other boolean conditions
	if detail.NeedsOverworldRain {
		fmt.Println(indent + "- Needs Overworld Rain: Yes")
	}

	if detail.TurnUpsideDown {
		fmt.Println(indent + "- Turn Upside Down: Yes")
	}

	// Print party species if present
	if detail.PartySpecies != nil {
		fmt.Println(indent + "- Party Species: " + detail.PartySpecies.Name)
	}

	// Print party type if present
	if detail.PartyType != nil {
		fmt.Println(indent + "- Party Type: " + detail.PartyType.Name)
	}

	// Print trade species if present
	if detail.TradeSpecies != nil {
		fmt.Println(indent + "- Trade Species: " + detail.TradeSpecies.Name)
	}

	// Print relative physical stats if present
	if detail.RelativePhysicalStats != nil {
		var statsDesc string
		switch *detail.RelativePhysicalStats {
		case -1:
			statsDesc = "Attack < Defense"
		case 0:
			statsDesc = "Attack = Defense"
		case 1:
			statsDesc = "Attack > Defense"
		default:
			statsDesc = fmt.Sprintf("%d", *detail.RelativePhysicalStats)
		}
		fmt.Println(indent + "- Relative Physical Stats: " + statsDesc)
	}

	// Print gender if present
	if detail.Gender != nil {
		var genderDesc string
		switch *detail.Gender {
		case 1:
			genderDesc = "Female"
		case 2:
			genderDesc = "Male"
		default:
			genderDesc = fmt.Sprintf("%d", *detail.Gender)
		}
		fmt.Println(indent + "- Gender: " + genderDesc)
	}
}
