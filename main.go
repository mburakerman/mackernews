package main

import (
	"fmt"
	"sort"
	"sync"

	"github.com/getlantern/systray"
	"github.com/mburakerman/mackernews/api"
	"github.com/mburakerman/mackernews/icon"
	"github.com/skratchdot/open-golang/open"
)

const GITHUB_URL = "https://github.com/mburakerman/mackernews/"

var strayNewsItems []*systray.MenuItem
var strayRefreshItem *systray.MenuItem
var strayAboutItem *systray.MenuItem
var strayQuitItem *systray.MenuItem

type NewsResult struct {
	idx  int
	news api.NewsItem
}

func listNewsItems() {
	newsIds, err := api.GetNewsIds()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	newsResultsCh := make(chan NewsResult, len(newsIds))

	var wg sync.WaitGroup
	wg.Add(len(newsIds))

	for index, newsId := range newsIds {
		go func(newsId api.NewsId, index int) {
			defer wg.Done()

			newsDetailItem, err := api.GetNewsDetails(int(newsId))
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			newsResultsCh <- NewsResult{idx: index, news: newsDetailItem}
		}(newsId, index)
	}
	wg.Wait()
	close(newsResultsCh)

	var sortedNews []NewsResult
	for res := range newsResultsCh {
		sortedNews = append(sortedNews, res)
	}

	sort.Slice(sortedNews, func(i, j int) bool {
		return sortedNews[i].idx < sortedNews[j].idx
	})

	for _, res := range sortedNews {
		title := fmt.Sprintf("%d. %s", res.idx+1, res.news.Title)
		strayNewsItem := systray.AddMenuItem(title, res.news.URL)
		strayNewsItems = append(strayNewsItems, strayNewsItem)

		go func(item *systray.MenuItem) {
			for {
				<-item.ClickedCh
				open.Run(res.news.URL)
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
	systray.SetIcon(icon.Data)
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
