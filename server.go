package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type GameResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Id          string `json:"id"`
}

type GameDetailResponse struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	Id            string `json:"id"`
	CurrentPrice  int    `json:"currentPrice"`
	SellerName    string `json:"sellerName"`
	DeveloperName string `json:"developerName"`
	PublisherName string `json:"publisherName"`
	ThumbnailUrl  string `json:"thumbnailUrl"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func GameHandler(w http.ResponseWriter, r *http.Request) {
	// Open the JSON file
	file, err := os.Open("games.json")
	if err != nil {
		errorJson := ErrorResponse{Error: "error opening json file"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorJson)
		return
	}
	defer file.Close()

	var jsonData GamesJson
	err = json.NewDecoder(file).Decode(&jsonData)
	if err != nil {
		errorJson := ErrorResponse{Error: "error reading json file"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorJson)
		return
	}
	defer file.Close()
	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")
	// Encode and write JSON data to the response
	gameResponses := []GameResponse{}
	games := jsonData.Data.Catalog.SearchStore.Elements
	for i := 0; i < len(games); i++ {
		game := games[i]
		gameResponse := GameResponse{Title: game.Title, Description: game.Description, Id: game.ID}
		gameResponses = append(gameResponses, gameResponse)
	}
	json.NewEncoder(w).Encode(gameResponses)
	return
}

func GameDetailHandler(w http.ResponseWriter, r *http.Request) {
	// Open the JSON file
	file, err := os.Open("games.json")
	if err != nil {
		errorJson := ErrorResponse{Error: "error opening json file"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorJson)
	}
	defer file.Close()

	var jsonData GamesJson
	err = json.NewDecoder(file).Decode(&jsonData)
	if err != nil {
		errorJson := ErrorResponse{Error: "error reading json file"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorJson)
	}
	// Get the id from the request URL
	id := r.URL.Query().Get("id")
	w.Header().Set("Content-Type", "application/json")

	// Check if the query parameter is present
	if id != "" {
		games := jsonData.Data.Catalog.SearchStore.Elements
		for i := 0; i < len(games); i++ {
			game := games[i]
			if id == game.ID {
				gameDetailResponse := GameDetailResponse{Title: game.Title, Description: game.Description, Id: game.ID, CurrentPrice: game.CurrentPrice, SellerName: game.Seller.Name, DeveloperName: game.DeveloperDisplayName, PublisherName: game.PublisherDisplayName, ThumbnailUrl: game.KeyImages[0].URL}
				json.NewEncoder(w).Encode(gameDetailResponse)
				return
			}
		}

		w.WriteHeader(http.StatusNotFound)
		errorJson := ErrorResponse{Error: "game not found"}
		json.NewEncoder(w).Encode(errorJson)
		return

	} else {
		errorJson := ErrorResponse{Error: "id parameter required"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorJson)
		return
	}
}

func main() {

	// Handle requests to the root path ("/")
	http.HandleFunc("/", GameHandler)

	http.HandleFunc("/game", GameDetailHandler)

	// Start the HTTP server on port 8080
	fmt.Println("Server listening on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
