package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"net/http"
	"strconv"

	"github.com/getlantern/systray"
)

const HACKERNEWS_TOP_STORIES_API = "https://hacker-news.firebaseio.com/v0/topstories.json"
const HACKERNEWS_NEWS_DETAIL_API = "https://hacker-news.firebaseio.com/v0/item/%s.json"
const NEWS_LIMIT = 5

type NewsId int

type NewsItem struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

func getHackernewsIds() []NewsId {
	response, err := http.Get(HACKERNEWS_TOP_STORIES_API)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("failed to read response body: %s", err)
		return nil
	}

	var newsIds []NewsId
	err = json.Unmarshal(body, &newsIds)
	if err != nil {
		fmt.Printf("failed to parse API response: %s", err)
		panic(err)
	}

	return newsIds[:NEWS_LIMIT]
}

func getHackernewsDetails(newsId string) NewsItem {
	apiURL := fmt.Sprintf(HACKERNEWS_NEWS_DETAIL_API, newsId)

	response, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("failed to read detail response body: %s", err)
		panic(err)
	}

	var details NewsItem
	err = json.Unmarshal(body, &details)
	if err != nil {
		fmt.Printf("failed to parse details API response: %s", err)
		panic(err)
	}

	return details
}

func listNewsItems() {
	var newsIds = getHackernewsIds()
	for i := 0; i < NEWS_LIMIT && i < len(newsIds); i++ {
		newsId := newsIds[i]
		newsDetailItem := getHackernewsDetails(strconv.Itoa(int(newsId)))
		fmt.Println(newsDetailItem.Title)
		systray.AddMenuItem(newsDetailItem.Title, newsDetailItem.URL)
	}
}

func onReady() {
	pngPath := "./icon.png"
	iconBytes, err := ioutil.ReadFile(pngPath)
	if err != nil {
		panic(err)
	}

	systray.SetIcon(iconBytes)
	systray.SetTooltip("Hacker News")

	listNewsItems()

	refreshItem := systray.AddMenuItem("Refresh", "")
	mQuit := systray.AddMenuItem("Quit", "Quit the app")

	for {
		select {
		case <-refreshItem.ClickedCh:
			listNewsItems()

		case <-mQuit.ClickedCh:
			systray.Quit()
			return
		}
	}
}

func main() {
	systray.Run(onReady, nil)
}
