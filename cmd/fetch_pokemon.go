package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// Helper function to fetch and unmarshal JSON
func fetchAndUnmarshal(url string, target interface{}) error {
	resp, err := httpClient.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching data from API: %v\n", err)
		os.Exit(1)
		return fmt.Errorf("failed to fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(os.Stderr, "Error: API request failed with status %s: %s\n", resp.Status, string(bodyBytes))
		os.Exit(1)
		return fmt.Errorf("API request to %s failed with status %s: %s", url, resp.Status, string(bodyBytes))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading API response body: %v\n", err)
		os.Exit(1)
		return fmt.Errorf("failed to read response body from %s: %w", url, err)
	}

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("failed to parse JSON from %s: %w", url, err)
	}
	return nil
}

func fetchPokemon(p *PokemonEntry) (*PokemonAPIData, error) {
	if p == nil {
		return nil, fmt.Errorf("PokemonEntry is nil")
	}
	fmt.Printf("Fetching details for %s from %s...\n\n", p.Name, p.URL)

	// Parse the JSON response
	var apiData PokemonAPIData
	if err := fetchAndUnmarshal(p.URL, &apiData); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON response: %v\n", err)
		os.Exit(1)
	}

	var speciesAPIData PokemonSpeciesAPIData
	if err := fetchAndUnmarshal(apiData.Species.URL, &speciesAPIData); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON response: %v\n", err)
		os.Exit(1)
	}

	var evolutionChainAPIData PokemonEvolutionChainAPIData
	if err := fetchAndUnmarshal(speciesAPIData.EvolutionChain.URL, &evolutionChainAPIData); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON response: %v\n", err)
		os.Exit(1)
	}

	// Display the fetched Pokémon information
	fmt.Printf("--- %s (ID: %d) ---\n", strings.Title(apiData.Name), apiData.ID)
	if apiData.Sprites.Other.OfficialArtwork.FrontDefault != "" {
		fmt.Printf("Sprite: %s\n", apiData.Sprites.Other.OfficialArtwork.FrontDefault)
	} else if apiData.Sprites.FrontDefault != "" {
		fmt.Printf("Sprite: %s\n", apiData.Sprites.FrontDefault)
	}

	fmt.Printf("Base Experience: %d\n", apiData.BaseExperience)
	fmt.Printf("Height: %.1f m\n", float64(apiData.Height)/10.0)  // Height is in decimetres
	fmt.Printf("Weight: %.1f kg\n", float64(apiData.Weight)/10.0) // Weight is in hectograms

	fmt.Print("Types: ")
	var types []string
	for _, t := range apiData.Types {
		types = append(types, strings.Title(t.TypeInfo.Name))
	}
	fmt.Println(strings.Join(types, ", "))

	fmt.Println("Abilities:")
	for _, a := range apiData.Abilities {
		hidden := ""
		if a.IsHidden {
			hidden = " (Hidden)"
		}
		fmt.Printf("  - %s%s\n", strings.Title(a.Ability.Name), hidden)
	}

	fmt.Println("Evolution Chain:")
	evolutionChainAPIData.PrintEvolutionChain()

	fmt.Println("Base Stats:")

	// Determine the maximum length of formatted stat names for alignment
	maxStatNameLength := 0
	totalBaseStats := 0 // Initialize totalBaseStats
	formattedStatNames := make([]string, len(apiData.Stats))
	for i, s := range apiData.Stats {
		name := strings.Title(strings.ReplaceAll(s.StatInfo.Name, "-", " "))
		formattedStatNames[i] = name
		if len(name) > maxStatNameLength {
			maxStatNameLength = len(name)
		}
		totalBaseStats += s.BaseStat // Accumulate total base stats
	}

	// Print table header
	tableWidth := maxStatNameLength + 10
	separator := "|" + strings.Repeat("-", tableWidth) + "|"
	topBorder := " " + strings.Repeat("-", tableWidth) + " "
	bottomBorder := topBorder
	fmt.Println(topBorder)
	header := fmt.Sprintf("| %-*s | %-5s |", maxStatNameLength, "Stat", "Value")
	fmt.Println(header)
	fmt.Printf("|%s|%s|\n", strings.Repeat("-", maxStatNameLength+2), strings.Repeat("-", 7))

	// Print each stat in a table row
	for i, s := range apiData.Stats {
		fmt.Printf("| %-*s | %5d |\n", maxStatNameLength, formattedStatNames[i], s.BaseStat)
	}

	// Print table footer
	fmt.Println(separator)
	fmt.Printf("| %-*s | %5d |\n", maxStatNameLength, "Total", totalBaseStats)
	fmt.Println(bottomBorder)

	return &apiData, nil
}
