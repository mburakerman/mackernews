package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"net/http"
	"strconv"

	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
)

const GITHUB_URL = "https://github.com/mburakerman/mackernews/"
const HACKERNEWS_TOP_STORIES_API = "https://hacker-news.firebaseio.com/v0/topstories.json"
const HACKERNEWS_NEWS_DETAIL_API = "https://hacker-news.firebaseio.com/v0/item/%s.json"
const NEWS_LIMIT = 5

type NewsId int

type NewsItem struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

func getNewsIds() ([]NewsId, error) {
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

	return newsIds[:NEWS_LIMIT], nil
}

func getNewsDetails(newsId string) (NewsItem, error) {
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

var strayNewsItems []*systray.MenuItem
var strayRefreshItem *systray.MenuItem
var strayAboutItem *systray.MenuItem
var strayQuitItem *systray.MenuItem

func listNewsItems() {
	newsIds, err := getNewsIds()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, newsId := range newsIds {
		newsIdString := strconv.Itoa(int(newsId))
		newsDetailItem, err := getNewsDetails(newsIdString)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		strayNewsItem := systray.AddMenuItem(newsDetailItem.Title, newsDetailItem.URL)
		strayNewsItems = append(strayNewsItems, strayNewsItem)

		go func(item *systray.MenuItem) {
			for {
				<-item.ClickedCh
				open.Run(newsDetailItem.URL)
			}
		}(strayNewsItem)
	}
}

func listAllItems() {
	listNewsItems()

	systray.AddSeparator()
	strayRefreshItem = systray.AddMenuItem("Refresh", "")
	strayAboutItem = systray.AddMenuItem("About Mackernews", "")
	strayQuitItem = systray.AddMenuItem("Quit", "Quit Mackernews")
}

func onReady() {
	pngPath := "./icon.png"
	iconBytes, err := os.ReadFile(pngPath)
	if err != nil {
		panic(err)
	}

	systray.SetIcon(iconBytes)
	systray.SetTooltip("Mackernews")

	listAllItems()

	for {
		select {
		case <-strayRefreshItem.ClickedCh:
			for _, strayNewsItem := range strayNewsItems {
				strayNewsItem.Hide()
			}
			strayRefreshItem.Hide()
			strayAboutItem.Hide()
			strayQuitItem.Hide()

			listAllItems()

		case <-strayAboutItem.ClickedCh:
			open.Run(GITHUB_URL)

		case <-strayQuitItem.ClickedCh:
			systray.Quit()
			return
		}
	}
}

func main() {
	systray.Run(onReady, nil)
}
