package main

// PokemonSpeciesAPIData represents the data structure for Pokemon species information
type PokemonSpeciesAPIData struct {
	ID                   int                    `json:"id"`
	Name                 string                 `json:"name"`
	Order                int                    `json:"order"`
	GenderRate           int                    `json:"gender_rate"`
	CaptureRate          int                    `json:"capture_rate"`
	BaseHappiness        int                    `json:"base_happiness"`
	IsBaby               bool                   `json:"is_baby"`
	IsLegendary          bool                   `json:"is_legendary"`
	IsMythical           bool                   `json:"is_mythical"`
	HatchCounter         int                    `json:"hatch_counter"`
	HasGenderDifferences bool                   `json:"has_gender_differences"`
	FormsSwitchable      bool                   `json:"forms_switchable"`
	Color                NamedAPIResource       `json:"color"`
	Shape                NamedAPIResource       `json:"shape"`
	EvolvesFromSpecies   *NamedAPIResource      `json:"evolves_from_species"`
	EvolutionChain       EvolutionChainResource `json:"evolution_chain"`
	Habitat              NamedAPIResource       `json:"habitat"`
	Generation           NamedAPIResource       `json:"generation"`
	GrowthRate           NamedAPIResource       `json:"growth_rate"`
	EggGroups            []NamedAPIResource     `json:"egg_groups"`
	Names                []PokemonSpeciesName   `json:"names"`
	PalParkEncounters    []PalParkEncounter     `json:"pal_park_encounters"`
	FlavorTextEntries    []FlavorTextEntry      `json:"flavor_text_entries"`
	FormDescriptions     []FormDescription      `json:"form_descriptions"`
	Genera               []Genus                `json:"genera"`
	Varieties            []PokemonVariety       `json:"varieties"`
	PokedexNumbers       []PokedexNumber        `json:"pokedex_numbers"`
}

// EvolutionChainResource represents a reference to an evolution chain
type EvolutionChainResource struct {
	URL string `json:"url"`
}

// PokemonSpeciesName represents a name for a Pokemon species in a specific language
type PokemonSpeciesName struct {
	Name     string           `json:"name"`
	Language NamedAPIResource `json:"language"`
}

// FlavorTextEntry represents a flavor text entry for a Pokemon species
type FlavorTextEntry struct {
	FlavorText string           `json:"flavor_text"`
	Language   NamedAPIResource `json:"language"`
	Version    NamedAPIResource `json:"version"`
}

// FormDescription represents a form description for a Pokemon species
type FormDescription struct {
	Description string           `json:"description"`
	Language    NamedAPIResource `json:"language"`
}

// Genus represents a genus for a Pokemon species in a specific language
type Genus struct {
	Genus    string           `json:"genus"`
	Language NamedAPIResource `json:"language"`
}

// PalParkEncounter represents a Pal Park encounter for a Pokemon species
type PalParkEncounter struct {
	BaseScore int              `json:"base_score"`
	Rate      int              `json:"rate"`
	Area      NamedAPIResource `json:"area"`
}

// PokemonVariety represents a variety of a Pokemon species
type PokemonVariety struct {
	IsDefault bool             `json:"is_default"`
	Pokemon   NamedAPIResource `json:"pokemon"`
}

// PokedexNumber represents a Pokedex number for a Pokemon species
type PokedexNumber struct {
	EntryNumber int              `json:"entry_number"`
	Pokedex     NamedAPIResource `json:"pokedex"`
}
