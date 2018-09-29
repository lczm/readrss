package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	ui "github.com/gizak/termui"
	"github.com/mmcdole/gofeed"
)

// focusStack
// 0 - rssNames
// 1 - rssContent
// 2 - rssContentExtended
// 3 - helpPage / errorPage

type Configuration struct {
	Rss []string
}

var (
	rssNamesCounter   = -1
	rssContentCounter = -1
	focusStack        = 0
	fullStack         = make(map[string]map[string]string)
)

func getCurrentFocus(stack []string, position int) string {

	focusString := stack[position]
	return focusString
}

func openInBrowser(url string) {
	// uses xdg-open
	shellCommand := "xdg-open"

	exec.Command(shellCommand, url).Start()
}

func errorPage(width int, height int, content []string) {

	page := ui.NewList()
	page.Items = content
	page.Overflow = "wrap"
	page.Width = width
	page.Height = height

	page.ItemFgColor = ui.ColorRed
	page.BorderFg = ui.ColorRed

	ui.Clear()
	ui.Render(page)

	focusStack = 3
}

func getConfig() []string {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}

	return configuration.Rss
}

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

	leftWidth := halfWidth
	rightWidth := termWidth - halfWidth

	addRssHeaderHeight := 3
	rssNamesHeight := termHeight - addRssHeaderHeight
	rssContentHeight := termHeight

	// get the contents of config file for rss links
	configRss := getConfig()

	rssHeader := []string{
		"Press 'a' to Add a RSS Feed",
	}

	rssNamesItems := []string{}

	for i := 0; i < len(configRss); i++ {
		rssNamesItems = append(rssNamesItems, configRss[i])
	}

	rssContentItems := []string{}

	// widget declarations

	addRssHeader := ui.NewList()
	addRssHeader.Items = rssHeader
	addRssHeader.Overflow = "wrap"
	addRssHeader.ItemFgColor = ui.ColorCyan
	addRssHeader.Width = leftWidth
	addRssHeader.Height = addRssHeaderHeight
	addRssHeader.BorderFg = ui.ColorBlue

	rssNames := ui.NewList()
	rssNames.Items = rssNamesItems
	rssNames.Overflow = "wrap"
	rssNames.ItemFgColor = ui.ColorCyan
	rssNames.Width = leftWidth
	rssNames.Height = rssNamesHeight
	// offset of the Y is the height of the top widget
	rssNames.Y = addRssHeaderHeight
	rssNames.BorderFg = ui.ColorMagenta

	rssContent := ui.NewList()
	rssContent.Items = rssContentItems
	rssContent.Overflow = "wrap"
	rssContent.ItemFgColor = ui.ColorCyan
	rssContent.X = halfWidth
	rssContent.Width = rightWidth
	rssContent.Height = rssContentHeight
	rssContent.BorderFg = ui.ColorDefault

	inputString := ""
	inputParagraph := ui.NewPar(inputString)
	inputParagraph.Height = 3
	inputParagraph.Width = 40
	inputParagraph.X = termWidth/2 - (40 / 2)
	inputParagraph.Y = termHeight / 2
	inputParagraph.BorderFg = ui.ColorDefault

	// put stuff in here later on
	helpPageItems := []string{
		"q : Exit",
		"a : Add RSS feed",
		"j : Move down",
		"k : Move up",
		"H : Help page [this]",
		"o : Open in browser",
		"Tab : Move focus between names and content",
		"Enter : Get feed / Expand details",
		"Esc : Escape out of current mode",
		"Press Esc to get out of this screen",
	}

	helpPage := ui.NewList()
	helpPage.Items = helpPageItems

	// positioning
	helpPage.Width = termWidth
	helpPage.Height = termHeight
	// position shouldn't matter as it takes up the whole page
	helpPage.X = 0
	helpPage.Y = 0

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

	ui.Handle("<Enter>", func(ui.Event) {

		if focusStack == 0 {
			if len(rssNamesItems) < 1 {
				// Error
				// ui.Clear()
				// ui.Render(addRssHeader, rssNames, rssContent)

				contents := []string{
					"There is currently no content in the header",
				}

				errorPage(termWidth, termHeight, contents)
			} else {
				focusString := getCurrentFocus(rssNamesItems,
					rssNamesCounter)
				fp := gofeed.NewParser()
				feed, _ := fp.ParseURL(focusString)
				items := feed.Items
				for i := 0; i < len(items); i++ {
					rssContentItems = append(rssContentItems, items[i].Title)
					// store an alternate version to reference back
					tempStack := make(map[string]string)

					tempStack["Description"] = items[i].Description
					tempStack["Published"] = items[i].Published
					tempStack["Link"] = items[i].Link

					fullStack[items[i].Title] = tempStack

				}
				rssContent.Items = rssContentItems
				ui.Clear()
				ui.Render(addRssHeader, rssNames, rssContent)
			}

		} else if focusStack == 1 {
			if len(rssContentItems) < 1 {
				contents := []string{
					"There is currently no content in the header",
				}

				errorPage(termWidth, termHeight, contents)
			} else {
				focusString := getCurrentFocus(rssContentItems,
					rssContentCounter)

				// Extended page items
				rssContentExtendedItems := []string{}

				// Description
				descriptionString := "Description : " + fullStack[focusString]["Description"]
				rssContentExtendedItems = append(rssContentExtendedItems, descriptionString)

				// blank line
				rssContentExtendedItems = append(rssContentExtendedItems, "")

				// Published
				publishedString := "Published : " + fullStack[focusString]["Published"]
				rssContentExtendedItems = append(rssContentExtendedItems, publishedString)

				// blank line
				rssContentExtendedItems = append(rssContentExtendedItems, "")

				// Link to source
				linkString := "Link: " + fullStack[focusString]["Link"]
				rssContentExtendedItems = append(rssContentExtendedItems, linkString)

				// blank line
				rssContentExtendedItems = append(rssContentExtendedItems, "")

				// widget for the new page
				rssContentExtended := ui.NewList()
				rssContentExtended.Overflow = "wrap"
				rssContentExtended.ItemFgColor = ui.ColorCyan
				rssContentExtended.Width = termWidth
				rssContentExtended.Height = termHeight
				rssContentExtended.BorderFg = ui.ColorMagenta

				rssContentExtended.Items = rssContentExtendedItems

				ui.Clear()
				ui.Render(rssContentExtended)

				focusStack = 2
			}
		} else {
			contents := []string{
				"There is currently no content in the header",
			}
			errorPage(termWidth, termHeight, contents)
		}

	})

	ui.Handle("o", func(ui.Event) {

		// mainly will only work
		// when the focus is on rssContent
		if focusStack == 1 {
			// get where the current focus is on

			focusString := getCurrentFocus(rssContentItems, rssContentCounter)

			link := fullStack[focusString]["Link"]

			openInBrowser(link)

			ui.Clear()
			ui.Render(addRssHeader, rssNames, rssContent)
		}

	})

	// TODO : set this to cancel out rssContentExtended
	ui.Handle("<Escape>", func(ui.Event) {
		// if the focus is currently on the rssContentExtended page
		if focusStack == 2 {
			ui.Clear()
			ui.Render(addRssHeader, rssNames, rssContent)
			// push it back to nothing
			focusStack = 1
		} else if focusStack == 3 {
			ui.Clear()
			ui.Render(addRssHeader, rssNames, rssContent)
			// set it back to default
			focusStack = 1
		}
	})

	ui.Handle("j", func(ui.Event) {

		// TODO : make this into a function
		stackClone := []string{}

		if focusStack == 0 {
			// use the rssNames stack
			stackClone = append(stackClone, rssNamesItems...)
			switch {
			case rssNamesCounter == len(stackClone)-1:
			default:
				rssNamesCounter++
			}

			stackClone[rssNamesCounter] = "[" +
				stackClone[rssNamesCounter] + "]" +
				"(fg-red,bg-green)"
		} else {
			// use the rssContent stack
			stackClone = append(stackClone, rssContentItems...)
			switch {
			case rssContentCounter == len(stackClone)-1:
			default:
				rssContentCounter++
			}
			stackClone[rssContentCounter] = "[" +
				stackClone[rssContentCounter] + "]" +
				"(fg-red,bg-green)"
		}

		// TODO : convert this into a regex function
		// add the syntax highlighting format

		// reset the Items with the new highlighting
		if focusStack == 0 {
			rssNames.Items = stackClone
		} else {
			rssContent.Items = stackClone
		}

		// re-render the background
		ui.Render(addRssHeader, rssNames, rssContent)

	})

	ui.Handle("k", func(ui.Event) {

		// TODO : make this into a function
		stackClone := []string{}

		if focusStack == 0 {
			// use the rssNames stack
			stackClone = append(stackClone, rssNamesItems...)
			switch {
			case rssNamesCounter < 1:
			default:
				rssNamesCounter--
			}

			stackClone[rssNamesCounter] = "[" +
				stackClone[rssNamesCounter] + "]" +
				"(fg-red,bg-green)"
		} else {
			// use the rssContent stack
			stackClone = append(stackClone, rssContentItems...)
			switch {
			case rssContentCounter < 1:
			default:
				rssContentCounter--
			}
			stackClone[rssContentCounter] = "[" +
				stackClone[rssContentCounter] + "]" +
				"(fg-red,bg-green)"
		}

		// TODO : convert this into a regex function
		// add the syntax highlighting format

		// reset the Items with the new highlighting
		if focusStack == 0 {
			rssNames.Items = stackClone
		} else {
			rssContent.Items = stackClone
		}

		// re-render the background
		ui.Render(addRssHeader, rssNames, rssContent)

	})

	ui.Handle("H", func(ui.Event) {
		ui.Clear()
		ui.Render(helpPage)
		focusStack = 3
	})

	ui.Handle("<Tab>", func(ui.Event) {
		// TODO : switch focus

		// 0 == rssNames
		// 1 == rssContent
		if focusStack == 0 {
			focusStack = 1
			rssNames.BorderFg = ui.ColorDefault
			rssContent.BorderFg = ui.ColorMagenta
			ui.Clear()
			ui.Render(addRssHeader, rssNames, rssContent)
		} else if focusStack == 1 {
			// switch it back to 0
			focusStack = 0
			rssNames.BorderFg = ui.ColorMagenta
			rssContent.BorderFg = ui.ColorDefault
			ui.Clear()
			ui.Render(addRssHeader, rssNames, rssContent)
		} else {
			focusStack = 0
			ui.Clear()
			ui.Render(addRssHeader, rssNames, rssContent)
		}
	})

	ui.Loop()
}
