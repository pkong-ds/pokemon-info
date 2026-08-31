package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// PokemonAPIEntry represents a single Pokémon entry from the PokeAPI "results" array.
type PokemonAPIEntry struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// APIResponse represents the basic structure of the PokeAPI list endpoints.
type APIResponse struct {
	Count    int               `json:"count"`
	Next     *string           `json:"next"` // Use pointer for nullable fields
	Previous *string           `json:"previous"`
	Results  []PokemonAPIEntry `json:"results"`
}

// FormattedPokemon represents the desired output structure for each Pokémon.
type FormattedPokemon struct {
	ID   int    `json:"id" yaml:"id"`
	Slug string `json:"slug" yaml:"slug"`
	Name string `json:"name" yaml:"name"`
	URL  string `json:"url" yaml:"url"`
}

// Global variables for Cobra flags
var outputFormat string
var outputFile string
var resource string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "prepare",
	Short: "Regenerates the pokemon-info Catalog from PokeAPI.",
	Long: `prepare is a maintainer-only tool that regenerates the pokemon-info
Catalog (the name lists embedded in the binary at build time) from the
PokeAPI (pokeapi.co). It outputs the fetched entries as JSON, YAML, or CSV.

It never ships inside the snap; users get Catalog updates via snap refresh.

Examples:
  prepare --output-format yaml --resource pokemon --output-file cmd/pokemon-info/pokemons.yaml
  prepare --output-format yaml --resource move --output-file cmd/pokemon-info/moves.yaml
  prepare --output-format csv > pokemons.csv`,
	// RunE is used for commands that return an error
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate output format
		outputFormat = strings.ToLower(outputFormat)
		if outputFormat != "json" && outputFormat != "yaml" && outputFormat != "csv" {
			return fmt.Errorf("invalid output format: %s. Must be one of json, yaml, csv", outputFormat)
		}

		resource = strings.ToLower(resource)
		if resource != "pokemon" && resource != "move" {
			return fmt.Errorf("invalid resource: %s. Must be one of pokemon, move", resource)
		}

		var allData []FormattedPokemon
		var err error

		switch resource {
		case "pokemon":
			fmt.Fprintln(os.Stderr, "Starting Pokémon data loading process...")
			allData, err = LoadAllPokemonNames()
			if err != nil {
				return fmt.Errorf("error loading Pokémon data: %w", err)
			}
			fmt.Fprintf(os.Stderr, "Successfully loaded %d Pokémon.\n", len(allData))
		case "move":
			fmt.Fprintln(os.Stderr, "Starting Move data loading process...")
			allData, err = LoadAllMoveNames()
			if err != nil {
				return fmt.Errorf("error loading Move data: %w", err)
			}
			fmt.Fprintf(os.Stderr, "Successfully loaded %d Moves.\n", len(allData))
		}

		allPokemon := allData

		// Determine output writer
		var writer io.Writer
		if outputFile != "" {
			file, err := os.Create(outputFile)
			if err != nil {
				return fmt.Errorf("failed to create output file %s: %w", outputFile, err)
			}
			defer file.Close()
			writer = file
			fmt.Fprintf(os.Stderr, "Writing output to file: %s\n", outputFile)
		} else {
			writer = os.Stdout
			fmt.Fprintln(os.Stderr, "Writing output to stdout...")
		}

		// Write output based on format
		switch outputFormat {
		case "json":
			encoder := json.NewEncoder(writer)
			encoder.SetIndent("", "  ") // Pretty print JSON
			if err := encoder.Encode(allPokemon); err != nil {
				return fmt.Errorf("failed to encode JSON: %w", err)
			}
		case "yaml":
			encoder := yaml.NewEncoder(writer)
			// encoder.SetIndent(2) // Default indentation for yaml.v3 is 2 spaces
			if err := encoder.Encode(allPokemon); err != nil {
				return fmt.Errorf("failed to encode YAML: %w", err)
			}
		case "csv":
			csvWriter := csv.NewWriter(writer)
			// Write header
			header := []string{"ID", "Slug", "Name", "URL"}
			if err := csvWriter.Write(header); err != nil {
				return fmt.Errorf("failed to write CSV header: %w", err)
			}
			// Write data rows
			for _, p := range allPokemon {
				row := []string{
					strconv.Itoa(p.ID),
					p.Slug,
					p.Name,
					p.URL,
				}
				if err := csvWriter.Write(row); err != nil {
					return fmt.Errorf("failed to write CSV row for Pokémon ID %d: %w", p.ID, err)
				}
			}
			csvWriter.Flush() // Ensure all data is written
			if err := csvWriter.Error(); err != nil {
				return fmt.Errorf("error flushing CSV writer: %w", err)
			}
		default:
			// This case should ideally be caught by the initial validation
			return fmt.Errorf("unsupported output format: %s", outputFormat)
		}

		fmt.Fprintln(os.Stderr, "Output generation complete.")
		return nil
	},
}

// init is called by Go before main()
func init() {
	// Define persistent flags, available to this command and all its children
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output-format", "f", "json", "Output format (json, yaml, csv)")
	rootCmd.PersistentFlags().StringVarP(&outputFile, "output-file", "o", "", "Output file path (default: stdout)")
	rootCmd.PersistentFlags().StringVarP(&resource, "resource", "r", "pokemon", "Resource to fetch (pokemon, move)")
}

// titleCase converts a snake_case or kebab-case string to Title Case.
func titleCase(s string) string {
	if s == "" {
		return ""
	}
	parts := strings.Split(s, "-")
	for i, part := range parts {
		if len(part) > 0 {
			runes := []rune(part)
			runes[0] = unicode.ToUpper(runes[0])
			for j := 1; j < len(runes); j++ {
				runes[j] = unicode.ToLower(runes[j])
			}
			parts[i] = string(runes)
		}
	}
	return strings.Join(parts, " ")
}

// extractIDFromURL parses the Pokémon ID from its URL.
func extractIDFromURL(url string) (int, error) {
	if !strings.HasSuffix(url, "/") {
		url += "/" // Ensure trailing slash for consistent parsing
	}
	parts := strings.Split(strings.TrimSuffix(url, "/"), "/")
	if len(parts) < 2 {
		return 0, fmt.Errorf("invalid Pokémon URL format: %s", url)
	}
	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("could not parse ID '%s' from URL %s: %w", idStr, url, err)
	}
	return id, nil
}

// LoadAllPokemonNames fetches all Pokémon names from the PokeAPI and formats them.
func LoadAllPokemonNames() ([]FormattedPokemon, error) {
	client := &http.Client{Timeout: 20 * time.Second} // Increased timeout for potentially large fetches
	baseURL := "https://pokeapi.co/api/v2/pokemon"

	// 1. Make initial request to get the count
	fmt.Fprintln(os.Stderr, "Fetching initial Pokémon list to get count...")
	req1, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create initial request: %w", err)
	}
	req1.Header.Set("User-Agent", "PokemonCLI/1.0") // Good practice to set User-Agent

	resp1, err := client.Do(req1)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch initial Pokémon list: %w", err)
	}
	defer resp1.Body.Close()

	if resp1.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("initial API request failed with status: %s", resp1.Status)
	}

	var initialResponse APIResponse
	if err := json.NewDecoder(resp1.Body).Decode(&initialResponse); err != nil {
		return nil, fmt.Errorf("failed to decode initial API response: %w", err)
	}

	totalCount := initialResponse.Count
	fmt.Fprintf(os.Stderr, "Total Pokémon count from API: %d\n", totalCount)
	if totalCount == 0 {
		return []FormattedPokemon{}, nil // No Pokémon to fetch
	}

	// 2. Make request to get all Pokémon using the count
	fetchAllURL := fmt.Sprintf("%s?limit=%d&offset=0", baseURL, totalCount)
	fmt.Fprintf(os.Stderr, "Fetching all %d Pokémon from: %s\n", totalCount, fetchAllURL)

	req2, err := http.NewRequest("GET", fetchAllURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for all Pokémon: %w", err)
	}
	req2.Header.Set("User-Agent", "PokemonCLI/1.0")

	resp2, err := client.Do(req2)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all Pokémon: %w", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request for all Pokémon failed with status: %s", resp2.Status)
	}

	var allPokemonResponse APIResponse
	if err := json.NewDecoder(resp2.Body).Decode(&allPokemonResponse); err != nil {
		return nil, fmt.Errorf("failed to decode API response for all Pokémon: %w", err)
	}

	// Check if we received the expected number of results
	if len(allPokemonResponse.Results) != totalCount {
		fmt.Fprintf(os.Stderr, "Warning: Expected %d Pokémon, but received %d. The API might have changed or there was an issue.\n", totalCount, len(allPokemonResponse.Results))
		// Proceed with what was received, or handle as a more critical error if desired.
	}

	// 3. Store all Pokémon names and IDs in the desired format
	formattedPokemons := make([]FormattedPokemon, 0, len(allPokemonResponse.Results))
	fmt.Fprintln(os.Stderr, "Processing Pokémon data...")

	for i, p := range allPokemonResponse.Results {
		id, err := extractIDFromURL(p.URL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not extract ID for Pokémon '%s' (URL: %s): %v. Skipping entry.\n", p.Name, p.URL, err)
			continue // Skip this entry if ID cannot be parsed
		}

		formattedPokemon := FormattedPokemon{
			ID:   id,
			Slug: p.Name,            // API provides it in snake-case
			Name: titleCase(p.Name), // Convert to Title Case
			URL:  p.URL,
		}
		formattedPokemons = append(formattedPokemons, formattedPokemon)
		if (i+1)%100 == 0 { // Print progress every 100 Pokémon
			fmt.Fprintf(os.Stderr, "Processed %d/%d Pokémon...\n", i+1, len(allPokemonResponse.Results))
		}
	}

	fmt.Fprintf(os.Stderr, "Successfully processed %d Pokémon.\n", len(formattedPokemons))
	return formattedPokemons, nil
}

// LoadAllMoveNames fetches all move names from the PokeAPI and formats them.
func LoadAllMoveNames() ([]FormattedPokemon, error) {
	client := &http.Client{Timeout: 20 * time.Second}
	baseURL := "https://pokeapi.co/api/v2/move"

	// 1. Make initial request to get the count
	fmt.Fprintln(os.Stderr, "Fetching initial move list to get count...")
	req1, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create initial request: %w", err)
	}
	req1.Header.Set("User-Agent", "PokemonCLI/1.0")

	resp1, err := client.Do(req1)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch initial move list: %w", err)
	}
	defer resp1.Body.Close()

	if resp1.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("initial API request failed with status: %s", resp1.Status)
	}

	var initialResponse APIResponse
	if err := json.NewDecoder(resp1.Body).Decode(&initialResponse); err != nil {
		return nil, fmt.Errorf("failed to decode initial API response: %w", err)
	}

	totalCount := initialResponse.Count
	fmt.Fprintf(os.Stderr, "Total move count from API: %d\n", totalCount)
	if totalCount == 0 {
		return []FormattedPokemon{}, nil
	}

	// 2. Make request to get all moves using the count
	fetchAllURL := fmt.Sprintf("%s?limit=%d&offset=0", baseURL, totalCount)
	fmt.Fprintf(os.Stderr, "Fetching all %d moves from: %s\n", totalCount, fetchAllURL)

	req2, err := http.NewRequest("GET", fetchAllURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for all moves: %w", err)
	}
	req2.Header.Set("User-Agent", "PokemonCLI/1.0")

	resp2, err := client.Do(req2)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all moves: %w", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request for all moves failed with status: %s", resp2.Status)
	}

	var allMovesResponse APIResponse
	if err := json.NewDecoder(resp2.Body).Decode(&allMovesResponse); err != nil {
		return nil, fmt.Errorf("failed to decode API response for all moves: %w", err)
	}

	if len(allMovesResponse.Results) != totalCount {
		fmt.Fprintf(os.Stderr, "Warning: Expected %d moves, but received %d.\n", totalCount, len(allMovesResponse.Results))
	}

	// 3. Store all move names and IDs in the desired format
	formattedMoves := make([]FormattedPokemon, 0, len(allMovesResponse.Results))
	fmt.Fprintln(os.Stderr, "Processing move data...")

	for i, m := range allMovesResponse.Results {
		id, err := extractIDFromURL(m.URL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not extract ID for move '%s' (URL: %s): %v. Skipping entry.\n", m.Name, m.URL, err)
			continue
		}

		formattedMove := FormattedPokemon{
			ID:   id,
			Slug: m.Name,
			Name: titleCase(m.Name),
			URL:  m.URL,
		}
		formattedMoves = append(formattedMoves, formattedMove)
		if (i+1)%100 == 0 {
			fmt.Fprintf(os.Stderr, "Processed %d/%d moves...\n", i+1, len(allMovesResponse.Results))
		}
	}

	fmt.Fprintf(os.Stderr, "Successfully processed %d moves.\n", len(formattedMoves))
	return formattedMoves, nil
}

// main is the entry point of the application
func main() {
	// Execute the root command. Cobra handles parsing flags and calling the appropriate Run function.
	if err := rootCmd.Execute(); err != nil {
		// Cobra's default error handling already prints the error,
		// but we explicitly print to os.Stderr for clarity if needed.
		// fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1) // Exit with a non-zero status code to indicate failure
	}
}
