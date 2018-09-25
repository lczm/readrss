package main

import (
	"fmt"
	"os"
	// "strings"

	ui "github.com/gizak/termui"
	"github.com/mmcdole/gofeed"
)

// TODO
// -start making functions

var (
	rssNamesCounter = 0
)

func main() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	// gets the width of the terminal
	termWidth := ui.TermWidth()
	// gets the height of the terminal
	termHeight := ui.TermHeight()

	halfWidth := termWidth / 2

	// leftWidth := utf8.RuneCountInString("Press 'a' to Add a RSS Feed")

	leftWidth := halfWidth
	rightWidth := termWidth - halfWidth

	addRssHeaderHeight := 3
	rssNamesHeight := termHeight - addRssHeaderHeight
	rssContentHeight := termHeight

	rssHeader := []string{
		"Press 'a' to Add a RSS Feed",
	}

	rssNamesItems := []string{}

	rssContentItems := []string{}

	addRssHeader := ui.NewList()
	addRssHeader.Items = rssHeader
	addRssHeader.Overflow = "wrap"
	addRssHeader.ItemFgColor = ui.ColorCyan
	addRssHeader.Width = leftWidth
	addRssHeader.Height = addRssHeaderHeight

	rssNames := ui.NewList()
	rssNames.Items = rssNamesItems
	rssNames.Overflow = "wrap"
	rssNames.ItemFgColor = ui.ColorCyan
	rssNames.Width = leftWidth
	rssNames.Height = rssNamesHeight
	// offset of the Y is the height of the top widget
	rssNames.Y = addRssHeaderHeight

	rssContent := ui.NewList()
	rssContent.Items = rssContentItems
	rssContent.Overflow = "wrap"
	rssContent.ItemFgColor = ui.ColorCyan
	rssContent.X = halfWidth
	rssContent.Width = rightWidth
	rssContent.Height = rssContentHeight

	inputString := ""
	inputParagraph := ui.NewPar(inputString)
	inputParagraph.Height = 3
	inputParagraph.Width = 40
	inputParagraph.X = termWidth/2 - (40 / 2)
	inputParagraph.Y = termHeight / 2

	ui.Render(addRssHeader, rssNames, rssContent)

	ui.Handle("q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("a", func(ui.Event) {
		ui.Render(inputParagraph)
	inputloop:
		for {

			character := make([]byte, 1, 1)
			// read from the standard input and do some processing
			input, _ := os.Stdin.Read(character)

			// TODO : check if whatever is entered is a valid url
			if input == 1 {
				// if its an escape - break out of the loop
				if character[0] == 27 {
					break inputloop
				} else if character[0] == 8 && len(inputString) > 0 ||
					// back spaces
					character[0] == 127 && len(inputString) > 0 {
					// remove the last character of the string
					inputString = inputString[:len(inputString)-1]

					inputParagraph := ui.NewPar(inputString)
					inputParagraph.Height = 3
					inputParagraph.Width = 40
					inputParagraph.X = termWidth/2 - (40 / 2)
					inputParagraph.Y = termHeight / 2

					// clears the entire screen
					ui.Clear()
					// render the background
					ui.Render(addRssHeader, rssNames, rssContent)
					// render the new textbox
					ui.Render(inputParagraph)
				} else if character[0] == 13 {
					// if its an enter - save & submit
					if len(inputString) < 1 {
						ui.Clear()
						break inputloop
					}
					rssNamesItems = append(rssNamesItems, inputString)
					rssNames.Items = rssNamesItems

					// reset the string so the previous
					// input is wiped
					inputString = ""

					break inputloop
				} else {
					// convert b[0] to a letter
					convertedAscii := string(character[0])
					inputString += convertedAscii

					inputParagraph := ui.NewPar(inputString)
					inputParagraph.Height = 3
					inputParagraph.Width = 40
					inputParagraph.X = termWidth/2 - (40 / 2)
					inputParagraph.Y = termHeight / 2

					// clears the entire screen
					ui.Clear()
					// render the background
					ui.Render(addRssHeader, rssNames, rssContent)
					// render the new textbox
					ui.Render(inputParagraph)
				}
			} else {
				// if its not a valid character
				break inputloop
			}
		}

		ui.Clear()
		ui.Render(addRssHeader, rssNames, rssContent)
	})

	ui.Handle("r", func(ui.Event) {
		// find the input
		rssInputFeed := rssNamesItems[len(rssNamesItems)-1]

		fp := gofeed.NewParser()
		feed, _ := fp.ParseURL(rssInputFeed)
		// fmt.Println(feed.Title)
		// feed the title to the content stack
		// TODO : have another stack that serves this while storing
		// the full json
		rssContentItems = append(rssContentItems, feed.Title)
		rssContent.Items = rssContentItems
		ui.Render(addRssHeader, rssNames, rssContent)
	})

	ui.Handle("j", func(ui.Event) {
		// get all the strings from the rssNames array
		// rssNamesCounter will be the 'height'

		// rssNamesItemsClone := rssNamesItems
		rssNamesItemsClone := []string{}
		rssNamesItemsClone = append(rssNamesItemsClone, rssNamesItems...)

		rssNamesItemsClone[rssNamesCounter] = "[" +
			rssNamesItemsClone[rssNamesCounter] +
			"]" + "(fg-red,bg-green)"

		// reset the Items with the new highlighting
		rssNames.Items = rssNamesItemsClone
		// re-render the background
		ui.Render(addRssHeader, rssNames, rssContent)

		// add one when going down
		rssNamesCounter++
	})

	ui.Handle("k", func(ui.Event) {
		fmt.Println("hello this is k")

		// minus one when going up
		rssNamesCounter--
	})

	ui.Handle("tab", func(ui.Event) {
		fmt.Println("hello this is tab")
	})

	ui.Loop()

}
