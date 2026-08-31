package main

// PokemonAPIData is the top-level structure for the PokeAPI response.
type PokemonAPIData struct {
	ID             int               `json:"id"`
	Name           string            `json:"name"`
	BaseExperience int               `json:"base_experience"`
	Height         int               `json:"height"`
	Weight         int               `json:"weight"`
	Abilities      []AbilitySlot     `json:"abilities"`
	Stats          []PokemonStat     `json:"stats"`
	Types          []PokemonTypeSlot `json:"types"`
	Sprites        PokemonSprites    `json:"sprites"`
	Species        NamedAPIResource  `json:"species"`
}

// AbilitySlot holds information about a Pokémon's ability.
type AbilitySlot struct {
	Ability  NamedAPIResource `json:"ability"`
	IsHidden bool             `json:"is_hidden"`
	Slot     int              `json:"slot"`
}

// NamedAPIResource contains the name and URL of an ability.
type NamedAPIResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// PokemonStat holds information about a Pokémon's stat.
type PokemonStat struct {
	BaseStat int              `json:"base_stat"`
	Effort   int              `json:"effort"`
	StatInfo NamedAPIResource `json:"stat"`
}

// PokemonTypeSlot holds information about a Pokémon's type.
type PokemonTypeSlot struct {
	Slot     int              `json:"slot"`
	TypeInfo NamedAPIResource `json:"type"`
}

// PokemonSprites holds URLs for various Pokémon sprites.
type PokemonSprites struct {
	FrontDefault string `json:"front_default"`
	Other        struct {
		OfficialArtwork struct {
			FrontDefault string `json:"front_default"`
		} `json:"official-artwork"`
	} `json:"other"`
}
