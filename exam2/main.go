package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

//โครงสร้างข้อมูล
type Stat struct {
	BaseStat int `json:"base_stat"`
	Effort   int `json:"effort"`
	Stat     struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"stat"`
}

type Sprites struct {
	BackDefault     string  `json:"back_default"`
	BackFemale      *string `json:"back_female"`
	BackShiny       string  `json:"back_shiny"`
	BackShinyFemale *string `json:"back_shiny_female"`
	FrontDefault    string  `json:"front_default"`
	FrontFemale     *string `json:"front_female"`
	FrontShiny      string  `json:"front_shiny"`
	FrontShinyFemale *string `json:"front_shiny_female"`
}

type PokemonResponse struct {
	Stats   []Stat  `json:"stats"`
	Name    string  `json:"name"`
	Sprites Sprites `json:"sprites"`
}

func fetchPokemonData(id string) (*PokemonResponse, error) {
	client := resty.New()

	//ดึง stats จาก /pokemon/{id}
	pokemonURL := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", id)
	resp1, err := client.R().Get(pokemonURL)
	if err != nil {
		return nil, err
	}

	var pokemonData struct {
		Stats []Stat `json:"stats"`
	}
	if err := json.Unmarshal(resp1.Body(), &pokemonData); err != nil {
		return nil, err
	}

	//ดึง name และ sprites จาก /pokemon-form/{id}
	formURL := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon-form/%s", id)
	resp2, err := client.R().Get(formURL)
	if err != nil {
		return nil, err
	}

	var formData struct {
		Name    string  `json:"name"`
		Sprites Sprites `json:"sprites"`
	}
	if err := json.Unmarshal(resp2.Body(), &formData); err != nil {
		return nil, err
	}

	//คืนค่า JSON ตามรูปแบบ
	return &PokemonResponse{
		Stats:   pokemonData.Stats,
		Name:    formData.Name,
		Sprites: formData.Sprites,
	}, nil
}

func main() {
	r := gin.Default()

	r.POST("/pokemon", func(c *gin.Context) {
		var request struct {
			ID string `json:"id"`
		}

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		pokemonData, err := fetchPokemonData(request.ID)
		if err != nil {
			log.Println("Error fetching Pokémon data:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Pokémon data"})
			return
		}

		c.JSON(http.StatusOK, pokemonData)
	})

	r.Run(":8080") //Server on port 8080
}
