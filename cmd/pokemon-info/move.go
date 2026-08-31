package main

// MoveAPIData is the top-level structure for the PokeAPI Move response.
type MoveAPIData struct {
	ID                int                    `json:"id"`
	Name              string                 `json:"name"`
	Accuracy          *int                   `json:"accuracy"`
	EffectChance      *int                   `json:"effect_chance"`
	PP                *int                   `json:"pp"`
	Priority          int                    `json:"priority"`
	Power             *int                   `json:"power"`
	DamageClass       NamedAPIResource       `json:"damage_class"`
	EffectEntries     []VerboseEffect        `json:"effect_entries"`
	FlavorTextEntries []MoveFlavorText       `json:"flavor_text_entries"`
	Generation        NamedAPIResource       `json:"generation"`
	Meta              *MoveMetaData          `json:"meta"`
	Names             []MoveName             `json:"names"`
	StatChanges       []MoveStatChange       `json:"stat_changes"`
	Target            NamedAPIResource       `json:"target"`
	Type              NamedAPIResource       `json:"type"`
	ContestType       *NamedAPIResource      `json:"contest_type"`
	Machines          []MachineVersionDetail `json:"machines"`
	LearnedByPokemon  []NamedAPIResource     `json:"learned_by_pokemon"`
	PastValues        []PastMoveStatValues   `json:"past_values"`
}

// VerboseEffect holds the effect description for a move.
type VerboseEffect struct {
	Effect      string           `json:"effect"`
	ShortEffect string           `json:"short_effect"`
	Language    NamedAPIResource `json:"language"`
}

// MoveFlavorText holds localized flavor text for a move.
type MoveFlavorText struct {
	FlavorText   string           `json:"flavor_text"`
	Language     NamedAPIResource `json:"language"`
	VersionGroup NamedAPIResource `json:"version_group"`
}

// MoveMetaData holds metadata about a move.
type MoveMetaData struct {
	Ailment       NamedAPIResource `json:"ailment"`
	Category      NamedAPIResource `json:"category"`
	MinHits       *int             `json:"min_hits"`
	MaxHits       *int             `json:"max_hits"`
	MinTurns      *int             `json:"min_turns"`
	MaxTurns      *int             `json:"max_turns"`
	Drain         int              `json:"drain"`
	Healing       int              `json:"healing"`
	CritRate      int              `json:"crit_rate"`
	AilmentChance int              `json:"ailment_chance"`
	FlinchChance  int              `json:"flinch_chance"`
	StatChance    int              `json:"stat_chance"`
}

// MoveStatChange holds information about a stat change caused by a move.
type MoveStatChange struct {
	Change int              `json:"change"`
	Stat   NamedAPIResource `json:"stat"`
}

// MoveName holds the localized name for a move.
type MoveName struct {
	Name     string           `json:"name"`
	Language NamedAPIResource `json:"language"`
}

// MachineVersionDetail holds information about a TM/HM that teaches a move.
type MachineVersionDetail struct {
	Machine      APIResource      `json:"machine"`
	VersionGroup NamedAPIResource `json:"version_group"`
}

// APIResource is a reference with only a URL (no name).
type APIResource struct {
	URL string `json:"url"`
}

// PastMoveStatValues holds historical stat values for a move across version groups.
type PastMoveStatValues struct {
	Accuracy      *int              `json:"accuracy"`
	EffectChance  *int              `json:"effect_chance"`
	Power         *int              `json:"power"`
	PP            *int              `json:"pp"`
	EffectEntries []VerboseEffect   `json:"effect_entries"`
	Type          *NamedAPIResource `json:"type"`
	VersionGroup  NamedAPIResource  `json:"version_group"`
}
