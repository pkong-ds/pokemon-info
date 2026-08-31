package main

import (
	// Required for embedding files
	_ "embed"
	"fmt"
	"log" // For fatal errors during initialization
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3" // For YAML parsing
)

//go:embed pokemons.yaml
var pokemonYAMLFile []byte

//go:embed moves.yaml
var moveYAMLFile []byte

// PokemonEntry defines the structure for each entry in the YAML file.
type PokemonEntry struct {
	ID   int    `yaml:"id"`
	Slug string `yaml:"slug"`
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

// MoveEntry defines the structure for each entry in the moves YAML file.
type MoveEntry struct {
	ID   int    `yaml:"id"`
	Slug string `yaml:"slug"`
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

var allPokemonEntries []PokemonEntry
var allMoveEntries []MoveEntry

// This init function will run once when the package is initialized.
func init() {
	// Load and parse Pokémon names from the embedded YAML file.
	err := yaml.Unmarshal(pokemonYAMLFile, &allPokemonEntries)
	if err != nil {
		log.Fatalf("Failed to parse pokemons.yaml: %v", err)
	}

	// If no Pokémon names are loaded, it's a critical issue for autocompletion.
	if len(allPokemonEntries) == 0 {
		log.Fatalln("No Pokémon names loaded from pokemons.yaml. Autocompletion will not work.")
	}

	// Load and parse move names from the embedded YAML file.
	err = yaml.Unmarshal(moveYAMLFile, &allMoveEntries)
	if err != nil {
		log.Fatalf("Failed to parse moves.yaml: %v", err)
	}

	// If no moves are loaded, it's a critical issue for autocompletion.
	if len(allMoveEntries) == 0 {
		log.Fatalln("No move names loaded from moves.yaml. Autocompletion will not work.")
	}

	// Register subcommands
	rootCmd.AddCommand(movesCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "pokemon-info <pokemon-name>",
	Short: "Fetches info for a given Pokémon.",
	Long: `pokemon-info is a CLI tool to display the base statistics
of a specified Pokémon. It features autocompletion for Pokémon names loaded from a YAML file.`,
	Example: `  pokemon-info Pikachu
  pokemon-info Char`, // User types Char then TAB
	Args: cobra.ExactArgs(1), // Expects exactly one argument
	Run: func(cmd *cobra.Command, args []string) {

		targetPokemonEntry := resolveWithSuggestions(allPokemonEntries, args[0], "Pokémon")
		if targetPokemonEntry == nil {
			os.Exit(1)
		}

		data, err := fetchPokemon(targetPokemonEntry)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching Pokémon data: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Fetched data for %s (ID: %d)\n", data.Name, data.ID)
	},
	// This is where the magic for autocompletion happens for the positional argument
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// ShellCompDirectiveNoFileComp tells the shell to not perform file completion
		// if our suggestions don't match.
		return suggestNames(allPokemonEntries, toComplete), cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveNoSpace
	},
}

var movesCmd = &cobra.Command{
	Use:   "moves <move-name>",
	Short: "Fetches info for a given Pokémon move.",
	Long:  `Fetches and displays detailed information about a Pokémon move from PokeAPI.`,
	Example: `  pokemon-info moves thunderbolt
  pokemon-info moves thunder-punch`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetMoveEntry := resolveWithSuggestions(allMoveEntries, args[0], "Move")
		if targetMoveEntry == nil {
			os.Exit(1)
		}
		fetchMove(targetMoveEntry)
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return suggestNames(allMoveEntries, toComplete), cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveNoSpace
	},
}
