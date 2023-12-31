package api

import (
	"encoding/json"
	"fmt"
	"io"

	"net/http"
)

const HACKERNEWS_TOP_STORIES_API = `https://hacker-news.firebaseio.com/v0/topstories.json?limitToFirst=10&orderBy="$key"`
const HACKERNEWS_NEWS_DETAIL_API = "https://hacker-news.firebaseio.com/v0/item/%s.json"

type NewsId int

type NewsItem struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

func GetNewsIds() ([]NewsId, error) {
	response, err := http.Get(HACKERNEWS_TOP_STORIES_API)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err)
	}

	var newsIds []NewsId
	err = json.Unmarshal(body, &newsIds)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API response: %s", err)
	}

	return newsIds, nil
}

func GetNewsDetails(newsId string) (NewsItem, error) {
	apiURL := fmt.Sprintf(HACKERNEWS_NEWS_DETAIL_API, newsId)

	response, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return NewsItem{}, fmt.Errorf("failed to read detail response body: %s", err)

	}

	var details NewsItem
	err = json.Unmarshal(body, &details)
	if err != nil {
		return NewsItem{}, fmt.Errorf("failed to parse details API response: %s", err)
	}

	return details, nil
}
