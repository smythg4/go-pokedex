package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"github.com/smythg4/go-pokedex/internal/pokecache"
)

var pokeCache = pokecache.NewCache(5 * time.Minute)

const BaseURL = "https://pokeapi.co/api/v2"

// Client represents a PokeAPI client
type Client struct {
	// You could add fields like HTTP client, base URL, etc.
}

// LocationAreaResponse represents the response from the location-area endpoint
type LocationAreaResponse struct {
	Next     *string             `json:"next"`
	Previous *string             `json:"previous"`
	Results  []LocationAreaResult `json:"results"`
}

// LocationAreaResult represents a single location area
type LocationAreaResult struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationAreaDetailResponse struct {
	Name 				string 				`json:"name"`
	PokemonEncounters 	[]PokemonEncounters	`json:"pokemon_encounters"`
}

type PokemonEncounters struct {
	Pokemon	Pokemon	`json:"pokemon"`
}

type Pokemon struct {
	Name			string		`json:"name"`
	URL 			string		`json:"url"`
	BaseExperience	int 		`json:"base_experience"`
	Height			int 		`json:"height"`
	Weight			int 		`json:"weight"`
	Stats			[]Stats 	`json:"stats"`
	Types 			[]Types 	`json:"types"`
}

type Stat struct {
	Name 	string 	`json:"name"`
	URL 	string	`json:"url"`
}

type Stats struct {
	BaseStat 	int		`json:"base_stat"`
	Effort		int 	`json:"effort"`
	Stat 		Stat 	`json:"stat"`
}

type Type struct {
	Name 	string 	`json:"name"`
	URL 	string 	`json:"url"`
}

type Types struct {
	Slot 	int 	`json:"slot"`
	Type 	Type 	`json:"type"`
}


// NewClient creates a new PokeAPI client
func NewClient() *Client {
	return &Client{}
}

// ListLocationAreas fetches a page of location areas
func (c *Client) ListLocationAreas(pageURL *string) (LocationAreaResponse, error) {
	url := fmt.Sprintf("%s/location-area", BaseURL)
	if pageURL != nil {
		url = *pageURL
	}

	//check the cache first
	if data, found := pokeCache.Get(url); found {
		var result LocationAreaResponse
		err := json.Unmarshal(data, &result)
		if err != nil {
			return LocationAreaResponse{}, fmt.Errorf("Error unmarshaling cached data: %w", err)
		}
		
		return result, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return LocationAreaResponse{}, fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return LocationAreaResponse{}, fmt.Errorf("unsuccessful status code: %d", resp.StatusCode)
	}

    data, err := io.ReadAll(resp.Body)
    if err != nil {
        return LocationAreaResponse{}, fmt.Errorf("error reading byte data off response body: %w",err)
    }
    
	// add results to cache
	pokeCache.Add(url, data)

    var result LocationAreaResponse
    err = json.Unmarshal(data, &result)
    if err != nil {
        return LocationAreaResponse{}, fmt.Errorf("error unmarshaling JSON: %w", err)
    }
    
    return result, nil
}

func (c *Client) GetLocationAreaDetails(name string) (LocationAreaDetailResponse, error) {
	url := fmt.Sprintf("%s/location-area/%s", BaseURL, name)

	//check the cache first
	if data, found := pokeCache.Get(url); found {
		var details LocationAreaDetailResponse
		err := json.Unmarshal(data, &details)
		if err != nil {
			return LocationAreaDetailResponse{}, fmt.Errorf("Error unmarshaling cached data: %w", err)
		}
		
		return details, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return LocationAreaDetailResponse{}, fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return LocationAreaDetailResponse{}, fmt.Errorf("unsuccessful status code: %d", resp.StatusCode)
	}

    data, err := io.ReadAll(resp.Body)
    if err != nil {
        return LocationAreaDetailResponse{}, fmt.Errorf("error reading byte data off response body: %w",err)
    }

	// add results to cache
	pokeCache.Add(url, data)

    var details LocationAreaDetailResponse
    err = json.Unmarshal(data, &details)
    if err != nil {
        return LocationAreaDetailResponse{}, fmt.Errorf("error unmarshaling JSON: %w", err)
    }
    
    return details, nil
}

func (c *Client) GetPokemonDetails(name string) (Pokemon, error) {
	url := fmt.Sprintf("%s/pokemon/%s", BaseURL, name)

	// check the cache
	if data, found := pokeCache.Get(url); found {
		var result Pokemon 
		err := json.Unmarshal(data, &result)
		if err != nil {
			return Pokemon{}, fmt.Errorf("Error unmarshaling cached data: %w", err)
		}
		return result, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return Pokemon{}, fmt.Errorf("Error getting response: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return Pokemon{}, fmt.Errorf("Unsuccessful status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
    if err != nil {
        return Pokemon{}, fmt.Errorf("error reading byte data off response body: %w",err)
    }

	// add results to cache
	pokeCache.Add(url, data)

	var result Pokemon 
	err = json.Unmarshal(data, &result)
	if err != nil {
		return Pokemon{}, fmt.Errorf("Error unmarshaling fetched data: %w", err)
	}
	return result, nil

}