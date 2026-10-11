package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
	"time"
)

const (
	endpoint  = "https://api.rawg.io/api/games"
	pageCount = 3
	pageSize  = 5
)

type game struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Released string  `json:"released"`
	Rating   float64 `json:"rating"`
}

type apiResponse struct {
	Results []game `json:"results"`
}

type pageResult struct {
	page  int
	games []game
	err   error
}

func main() {
	key := os.Getenv("RAWG_API_KEY")
	if key == "" {
		fmt.Fprintln(os.Stderr, "set RAWG_API_KEY before running this example")
		os.Exit(2)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client := &http.Client{Timeout: 12 * time.Second}
	results := make(chan pageResult, pageCount)

	for page := 1; page <= pageCount; page++ {
		go func(page int) {
			games, err := fetchPage(ctx, client, key, page)
			results <- pageResult{page: page, games: games, err: err}
		}(page)
	}

	completed := make([]pageResult, 0, pageCount)
	for i := 0; i < pageCount; i++ {
		completed = append(completed, <-results)
	}
	sort.Slice(completed, func(i, j int) bool { return completed[i].page < completed[j].page })

	successful := 0
	for _, result := range completed {
		if result.err != nil {
			fmt.Fprintf(os.Stderr, "page %d: %v\n", result.page, result.err)
			continue
		}
		successful++
		fmt.Printf("\nPage %d (%d games)\n", result.page, len(result.games))
		for _, item := range result.games {
			fmt.Printf("%-7d | %-36s | %-10s | %.2f\n", item.ID, item.Name, item.Released, item.Rating)
		}
	}
	fmt.Printf("\nFetched %d of %d pages.\n", successful, pageCount)
	if successful != pageCount {
		os.Exit(1)
	}
}

func fetchPage(ctx context.Context, client *http.Client, key string, page int) ([]game, error) {
	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	query := requestURL.Query()
	query.Set("key", key)
	query.Set("page", fmt.Sprint(page))
	query.Set("page_size", fmt.Sprint(pageSize))
	requestURL.RawQuery = query.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("API returned %s", response.Status)
	}
	var payload apiResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("could not decode API response: %w", err)
	}
	return payload.Results, nil
}
