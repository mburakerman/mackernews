package main

import (
	"fmt"

	_ "embed"
	"strconv"

	"github.com/getlantern/systray"
	"github.com/mburakerman/mackernews/api"
	"github.com/skratchdot/open-golang/open"
)

//go:embed icon.png
var iconByte []byte

const GITHUB_URL = "https://github.com/mburakerman/mackernews/"

var strayNewsItems []*systray.MenuItem
var strayRefreshItem *systray.MenuItem
var strayAboutItem *systray.MenuItem
var strayQuitItem *systray.MenuItem

func listNewsItems() {
	newsIds, err := api.GetNewsIds()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for index, newsId := range newsIds {
		newsIdString := strconv.Itoa(int(newsId))
		newsDetailItem, err := api.GetNewsDetails(newsIdString)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		title := fmt.Sprintf("%d. %s", index+1, newsDetailItem.Title)
		strayNewsItem := systray.AddMenuItem(title, newsDetailItem.URL)
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
	systray.SetIcon(iconByte)
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
