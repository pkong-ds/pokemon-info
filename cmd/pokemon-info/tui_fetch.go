package main

import (
	"encoding/json"
	"fmt"
)

// getJSON fetches url and unmarshals into target without printing or exiting,
// so the TUI can handle errors interactively.
func getJSON(url string, target interface{}) error {
	body, err := cachedGet(url)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, target)
}

// pokemonDetail bundles everything needed to render one Pokémon in the TUI.
type pokemonDetail struct {
	Pokemon *PokemonAPIData
	Species *PokemonSpeciesAPIData
	Chain   *PokemonEvolutionChainAPIData
}

func fetchPokemonDetail(p *PokemonEntry) (*pokemonDetail, error) {
	if p == nil {
		return nil, fmt.Errorf("no Pokémon selected")
	}

	var apiData PokemonAPIData
	if err := getJSON(p.URL, &apiData); err != nil {
		return nil, err
	}

	d := &pokemonDetail{Pokemon: &apiData}

	var species PokemonSpeciesAPIData
	if err := getJSON(apiData.Species.URL, &species); err != nil {
		return nil, err
	}
	d.Species = &species

	var chain PokemonEvolutionChainAPIData
	if err := getJSON(species.EvolutionChain.URL, &chain); err != nil {
		return nil, err
	}
	d.Chain = &chain

	return d, nil
}

func fetchMoveDetail(m *MoveEntry) (*MoveAPIData, error) {
	if m == nil {
		return nil, fmt.Errorf("no move selected")
	}

	var moveData MoveAPIData
	if err := getJSON(m.URL, &moveData); err != nil {
		return nil, err
	}
	return &moveData, nil
}
