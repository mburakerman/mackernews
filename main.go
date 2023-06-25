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

func listNewsItems() {
	newsIds, err := getNewsIds()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for i := 0; i < NEWS_LIMIT && i < len(newsIds); i++ {
		newsId := newsIds[i]
		newsDetailItem, err := getNewsDetails(strconv.Itoa(int(newsId)))
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		item := systray.AddMenuItem(newsDetailItem.Title, newsDetailItem.URL)
		item.Disabled()
		go func() {
			for {
				<-item.ClickedCh
				open.Run(newsDetailItem.URL)
			}
		}()
	}
}

func onReady() {
	pngPath := "./icon.png"
	iconBytes, err := os.ReadFile(pngPath)
	if err != nil {
		panic(err)
	}

	systray.SetIcon(iconBytes)
	systray.SetTooltip("Mackernews")

	listNewsItems()

	systray.AddSeparator()
	aboutItem := systray.AddMenuItem("About Mackernews", "")
	quitItem := systray.AddMenuItem("Quit", "Quit Mackernews")

	for {
		select {
		case <-aboutItem.ClickedCh:
			open.Run("https://github.com/mburakerman/mackernews/")

		case <-quitItem.ClickedCh:
			systray.Quit()
			return
		}
	}
}

func main() {
	systray.Run(onReady, nil)
}
