package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Restaurant struct {
	Name     string    `json:"name"`
	Address  Address   `json:"address"`
	Rating   Rating    `json:"rating"`
	Cuisines []Cuisine `json:"cuisines"`
}

type Address struct {
	FirstLine  string `json:"firstLine"`
	City       string `json:"city"`
	PostalCode string `json:"postalCode"`
}

type Rating struct {
	StarRating float64 `json:"starRating"`
}

type Cuisine struct {
	Name string `json:"name"`
}

func main() {
	restaurants, err := RetrieveRestaurant()
	if err != nil {
		log.Fatalf("failed to retrieve restaurants: %v", err)
	}

	DisplayTopRestaurants(restaurants)
}

func RetrieveRestaurant() ([]Restaurant, error) {
	var postCode string
	fmt.Println("Please enter a postal code")
	fmt.Scan(&postCode)

	resp, err := http.Get(fmt.Sprintf("https://uk.api.just-eat.io/discovery/uk/restaurants/enriched/bypostcode/%s", postCode))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve data from API: %v", err)
	}
	defer resp.Body.Close()

	var data struct {
		Restaurants []Restaurant `json:"restaurants"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %v", err)
	}

	return data.Restaurants, nil
}

func DisplayTopRestaurants(restaurants []Restaurant) {
	for i, v := range restaurants {
		if i >= 10 {
			break
		}
		PrintRestaurant(i+1, v)
	}
}

func PrintRestaurant(index int, restaurant Restaurant) {
	fmt.Printf("Restaurant: %d\n", index)
	fmt.Printf("Name: %s\n", restaurant.Name)
	fmt.Printf("Address: %s, %s, %s\n", restaurant.Address.FirstLine, restaurant.Address.City, restaurant.Address.PostalCode)
	fmt.Printf("Cuisines: %s\n", JoinCuisines(restaurant.Cuisines))
	fmt.Printf("Rating: %.1f\n", restaurant.Rating.StarRating)
	fmt.Println("---------------------------------------------------------------")
}

func JoinCuisines(cuisines []Cuisine) string {
	names := make([]string, len(cuisines))
	for i, c := range cuisines {
		names[i] = c.Name
	}
	return strings.Join(names, ", ")
}

