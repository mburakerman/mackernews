package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"net/http"
	"reflect"
	"strconv"

	"github.com/getlantern/systray"
)

const HACKERNEWS_TOP_STORIES_API = "https://hacker-news.firebaseio.com/v0/topstories.json"
const HACKERNEWS_NEWS_DETAIL_API = "https://hacker-news.firebaseio.com/v0/item/%s.json"
const NEWS_LIMIT = 5

type NewsItem struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

func getHackernewsIds() []int {
	response, err := http.Get(HACKERNEWS_TOP_STORIES_API)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %s", err)
		return nil
	}

	var newsIds []int
	err = json.Unmarshal(body, &newsIds)
	if err != nil {
		fmt.Printf("Failed to parse API response: %s", err)
		return nil
	}

	fmt.Printf("type of a is %v\n", reflect.TypeOf(body))

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
		fmt.Println("Error:", err)
		panic(err)
	}

	var details NewsItem
	err = json.Unmarshal(body, &details)
	if err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}

	return details
}

func onReady() {
	pngPath := "./icon.png"
	iconBytes, err := ioutil.ReadFile(pngPath)
	if err != nil {
		panic(err)
	}

	systray.SetIcon(iconBytes)
	systray.SetTooltip("Hacker News")

	var newsIds = getHackernewsIds()
	for i := 0; i < NEWS_LIMIT && i < len(newsIds); i++ {
		newsId := newsIds[i]
		newsDetailItem := getHackernewsDetails(strconv.Itoa(newsId))
		fmt.Println(newsDetailItem.Title)
		systray.AddMenuItem(newsDetailItem.Title, newsDetailItem.URL)
	}

	mToggle := systray.AddMenuItem("Toggle", "bla bla")
	mQuit := systray.AddMenuItem("Quit", "Quit the app")

	for {
		select {
		case <-mToggle.ClickedCh:
			getHackernewsIds()

		case <-mQuit.ClickedCh:
			systray.Quit()
			return
		}
	}

}

func main() {
	systray.Run(onReady, nil)
}
