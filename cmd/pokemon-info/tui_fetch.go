package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// getJSON fetches url and unmarshals into target without printing or exiting,
// so the TUI can handle errors interactively.
func getJSON(url string, target interface{}) error {
	resp, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request to %s failed with status %s", url, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body from %s: %w", url, err)
	}

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("failed to parse JSON from %s: %w", url, err)
	}
	return nil
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
