package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"math/rand"
	"github.com/smythg4/go-pokedex/internal/pokeapi"
)

type cliCommand struct {
	name			string
	description		string
	callback		func(*Config, []string) error
}

type Config struct {
	//strings to store next and previous URLs for navigation
	Next			*string
	Previous		*string
}

var commandRegistry map[string]cliCommand

var userPokeDex map[string]pokeapi.Pokemon

func init() {
    commandRegistry = map[string]cliCommand{
        "exit": {
            name:        "exit",
            description: "Exit the Pokedex",
            callback:    commandExit,
        },
        "help": {
            name:        "help",
            description: "Displays a help message",
            callback:    commandHelp,
        },
		"map": {
			name:		"map",
			description: "Displays locations in the Pokemon world 20 at a time",
			callback:	commandMap,
		},
		"mapb": {
			name:		"mapb",
			description: "Displays the previous page of locations in the Pokemon world",
			callback:	commandMapBack,
		},
		"explore": {
			name:		"explore",
			description: "Given a location area name, returns list of Pokemon found in that area",
			callback:	commandExplore,
		},
		"catch": {
			name:			"catch",
			description:	"Given a Pokemon name it makes an attempt to catch it and add it to your pokedex",
			callback:	commandCatch,
		},
		"pokedex": {
			name:			"pokedex",
			description:	"Displays names of every Pokemon you've caught",
			callback:	commandDex,
		},
		"inspect": {
			name:			"inspect",
			description:	"Displays details about a pokemon you've caught",
			callback:	commandInspect,
		},
    }

	userPokeDex = make(map[string]pokeapi.Pokemon)
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandExit(conf *Config, params []string)  error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf *Config, params []string)  error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage: ")
	fmt.Println()
	for name, cmd := range commandRegistry {
		fmt.Printf("%s: %s\n", name, cmd.description)
	}
	return nil
}

func commandMap(conf *Config, params []string)  error {

	client := pokeapi.NewClient()

	resp, err := client.ListLocationAreas(conf.Next)
	if err != nil {
		return err
	}

	for _, area := range resp.Results {
		fmt.Println(area.Name)
	}

	conf.Next = resp.Next
	conf.Previous = resp.Previous

	return nil
}

func commandMapBack(conf *Config, params []string)  error {

	if conf.Previous == nil {
		fmt.Println("You're on the first page")
		return nil
	}

	client := pokeapi.NewClient()

	resp, err := client.ListLocationAreas(conf.Previous)
	if err != nil {
		return err
	}

	for _, area := range resp.Results {
		fmt.Println(area.Name)
	}

	conf.Next = resp.Next
	conf.Previous = resp.Previous

	return nil

	return nil
}

func commandExplore(conf *Config, params []string) error {

	if len(params) < 1 {
		return fmt.Errorf("No location area provided.")
	}

	// location area name should be the first (and only) argument provided
	name := params[0]

	fmt.Printf("Exporing %s...\n", name)

	client := pokeapi.NewClient()

	resp, err := client.GetLocationAreaDetails(name)
	
	if err != nil {
		return err
	}
	fmt.Println("Found Pokemon:")

	for _, encounter := range resp.PokemonEncounters {
		fmt.Printf(" - %s\n",encounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(conf *Config, params []string)  error {
	if len(params) < 1 {
		return fmt.Errorf("No pokemon name provided.")
	}

	name := params[0]

	client := pokeapi.NewClient()
	this_pokemon, err := client.GetPokemonDetails(name)
	if err != nil {
		return fmt.Errorf("Error finding %s: %w", name, err)
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", this_pokemon.Name)

	die_roll := rand.Intn(400)

	if die_roll > this_pokemon.BaseExperience {
		fmt.Printf("%s was caught!\n", this_pokemon.Name)
		fmt.Println("You may now inspect it with the inspect command.")
		// add to pokedex
		userPokeDex[name] = this_pokemon
	} else {
		fmt.Printf("%s escaped!\n", this_pokemon.Name)
		// don't add to pokedex
	}

	return nil
}

func commandDex(conf *Config, params []string)  error {
	fmt.Println("Your Pokedex:")
	for key, _ := range userPokeDex {
		fmt.Printf(" - %s\n",key)
	}
	return nil
}

func commandInspect(conf *Config, params []string) error {
	if len(params) < 1 {
		return fmt.Errorf("No pokemon name provided.")
	}

	name := params[0]

	my_pokemon, ok := userPokeDex[name]
	if !ok {
		// check if they caught this pokemon yet
		fmt.Println("You have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", my_pokemon.Name)
	fmt.Printf("Height: %d\n", my_pokemon.Height)
	fmt.Printf("Weight: %d\n", my_pokemon.Weight)
	
	fmt.Println("Stats:")
	for _, stats := range my_pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stats.Stat.Name, stats.BaseStat)
	}

	fmt.Println("Types:")
	for _, types := range my_pokemon.Types {
		fmt.Printf(" - %s\n", types.Type.Name)
	}

	return nil
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	config := Config{}

    for  {
		fmt.Print("Pokedex > ")
		scanner.Scan()
        line := scanner.Text()
        clean_line := cleanInput(line)
		if len(clean_line[0]) > 0 {
			command, ok := commandRegistry[clean_line[0]]
			if ok {
				err := command.callback(&config, clean_line[1:])
				if err != nil {
					fmt.Printf("Error with command %s: %s\n", command.name, err)
				}
			} else {
				fmt.Println("Unknown command")
			}
		} else {
			continue
		}
    }
}