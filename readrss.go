package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	ui "github.com/gizak/termui"
	"github.com/mmcdole/gofeed"
)

// struct for loading json
type Configuration struct {
	Rss []string
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

type widgetMaker struct {
	List *ui.List
	Par  *ui.Par
}

func getTermSize() (int, int) {
	width := ui.TermWidth()
	height := ui.TermHeight()
	return width, height
}

func makeListWidget(content []string, label string, width int, height int,
	x int, y int, colour string) *widgetMaker {

	widgetList := ui.NewList()
	widgetList.Items = content
	widgetList.Overflow = "wrap"
	widgetList.ItemFgColor = ui.ColorWhite
	widgetList.Width = width
	widgetList.Height = height

	widgetList.X = x // x location
	widgetList.Y = y // y location

	widgetList.BorderLabel = label

	// TODO : tidy this up

	switch {
	case colour == "Red":
		widgetList.BorderFg = ui.ColorRed
	case colour == "Magenta":
		widgetList.BorderFg = ui.ColorMagenta
	case colour == "Cyan":
		widgetList.BorderFg = ui.ColorCyan
	case colour == "Black":
		widgetList.BorderFg = ui.ColorBlack
	case colour == "White":
		widgetList.BorderFg = ui.ColorWhite
	case colour == "Default":
		widgetList.BorderFg = ui.ColorDefault
	case colour == "Yellow":
		widgetList.BorderFg = ui.ColorYellow
	case colour == "Green":
		widgetList.BorderFg = ui.ColorGreen
	default:
		widgetList.BorderFg = ui.ColorDefault
	}

	return &widgetMaker{List: widgetList}
}

func makeParWidget(content string, width int, height int,
	x int, y int, colour string) *widgetMaker {

	widgetPar := ui.NewPar(content)
	widgetPar.Width = width
	widgetPar.Height = height
	widgetPar.X = x // x location
	widgetPar.Y = y // y location

	switch {
	case colour == "Red":
		widgetPar.BorderFg = ui.ColorRed
	case colour == "Magenta":
		widgetPar.BorderFg = ui.ColorMagenta
	case colour == "Cyan":
		widgetPar.BorderFg = ui.ColorCyan
	case colour == "Black":
		widgetPar.BorderFg = ui.ColorBlack
	case colour == "White":
		widgetPar.BorderFg = ui.ColorWhite
	case colour == "Default":
		widgetPar.BorderFg = ui.ColorDefault
	case colour == "Yellow":
		widgetPar.BorderFg = ui.ColorYellow
	case colour == "Green":
		widgetPar.BorderFg = ui.ColorGreen
	default:
		widgetPar.BorderFg = ui.ColorDefault
	}

	return &widgetMaker{Par: widgetPar}
}

func openInBrowser(url string) {
	// uses xdg-open
	shellCommand := "xdg-open"
	exec.Command(shellCommand, url).Start()
}

// rssHeaders.Items = highlightFocus(rssHeaders.Items, 0) - example usage
func highlightFocus(array []string, position int) []string {
	array[position] = "[" + array[position] + "]" + "(fg-red,bg-green)"
	return array
}

func characterFunction(starterString string, character []byte) string {

	if character[0] == 127 && len(starterString) > 0 {
		starterString = starterString[:len(starterString)-1]
	} else {
		starterString += string(character[0])
	}
	return starterString
}

func changeWidthHeight(widget *widgetMaker, width int, height int) {
	widget.List.Width = width
	widget.List.Height = height
}

func main() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	// state declaration
	state := make(map[string]int)
	state["rssHeaderCounter"] = -1
	state["rssContentCounter"] = -1
	state["focusStack"] = 0

	fullStack := make(map[string]map[string]string)

	// originalWidth, originalHeight := getTermSize()
	originalWidth := ui.TermWidth()
	originalHeight := ui.TermHeight()

	termWidth := &originalWidth
	termHeight := &originalHeight

	// contents of config file
	rssHeadersItems := getConfig()
	rssContentsItems := []string{}

	// rssHeaders := makeListWidget(rssHeadersItems, "Feed", termWidth/2,
	// 	termHeight, 0, 0, "Magenta")
	// rssContents := makeListWidget(rssContentsItems, "Content", termWidth/2,
	// 	termHeight, termWidth/2, 0, "White")
	rssHeaders := makeListWidget(rssHeadersItems, "Feed", *termWidth/2, *termHeight, 0, 0, "Magenta")
	rssContents := makeListWidget(rssContentsItems, "Content", *termWidth/2, *termHeight, *termWidth/2, 0, "White")

	ui.Render(rssHeaders.List, rssContents.List)

	ui.Handle("a", func(ui.Event) {
		contentString := ""
		inputParagraph := makeParWidget(contentString, 40, 3,
			*termWidth/2-(40/2), *termHeight/2, "Default")
		ui.Render(inputParagraph.Par)

	inputloop:
		for {
			character := make([]byte, 1, 1)
			input, _ := os.Stdin.Read(character)
			if input == 1 {
				switch {
				case character[0] == 27: // escape
					break inputloop
				case character[0] == 13: // enter
					// save the inputString to somewhere
					rssHeadersItems = append(rssHeadersItems, contentString)
					rssHeaders.List.Items = rssHeadersItems

					break inputloop
				}
				// feed character into a function
				contentString = characterFunction(contentString, character)
				inputParagraph := makeParWidget(contentString, 40, 3,
					*termWidth/2-(40/2), *termHeight/2, "Default")

				ui.Clear()
				ui.Render(rssHeaders.List, rssContents.List)
				ui.Render(inputParagraph.Par)
			} else {
				break inputloop
			}
		}
		ui.Clear()
		ui.Render(rssHeaders.List, rssContents.List)
	})

	ui.Handle("q", func(ui.Event) { ui.StopLoop() })
	ui.Handle("j", func(ui.Event) { // go down in terms of focus
		stackClone := []string{}

		rssHeaders := makeListWidget(rssHeadersItems, "Feed", *termWidth/2, *termHeight, 0, 0, "Magenta")
		rssContents := makeListWidget(rssContentsItems, "Content", *termWidth/2, *termHeight, *termWidth/2, 0, "White")

	higherSwitchJ:
		switch {
		case state["focusStack"] == 0:
			switch {
			case state["rssHeaderCounter"] < 0:
				state["rssHeaderCounter"] = 0
			case state["rssHeaderCounter"] >= len(rssHeadersItems)-1:
			default:
				state["rssHeaderCounter"]++
			}
			stackClone = append(stackClone, rssHeadersItems...)
			stackClone = highlightFocus(stackClone, state["rssHeaderCounter"])
			rssHeaders.List.Items = stackClone
			ui.Render(rssHeaders.List, rssContents.List)
		case state["focusStack"] == 1:
			switch {
			case len(rssContents.List.Items) == 0:
				break higherSwitchJ
			case state["rssContentCounter"] < 0:
				state["rssContentCounter"] = 0
			case state["rssContentCounter"] >= len(rssContentsItems)-1:
			default:
				state["rssContentCounter"]++
			}
			stackClone = append(stackClone, rssContentsItems...)
			stackClone = highlightFocus(stackClone, state["rssContentCounter"])
			rssContents.List.Items = stackClone
			ui.Render(rssHeaders.List, rssContents.List)
		default:
		}

	})
	ui.Handle("k", func(ui.Event) { // go up in terms of focus
		stackClone := []string{}

		rssHeaders := makeListWidget(rssHeadersItems, "Feed", *termWidth/2, *termHeight, 0, 0, "Magenta")
		rssContents := makeListWidget(rssContentsItems, "Content", *termWidth/2, *termHeight, *termWidth/2, 0, "White")

	higherSwitchK:
		switch {
		case state["focusStack"] == 0:
			switch {
			case state["rssHeaderCounter"] < 1:
				state["rssHeaderCounter"] = 0
			default:
				state["rssHeaderCounter"]--
			}
			stackClone = append(stackClone, rssHeadersItems...)
			stackClone = highlightFocus(stackClone, state["rssHeaderCounter"])
			rssHeaders.List.Items = stackClone
			ui.Render(rssHeaders.List, rssContents.List)
		case state["focusStack"] == 1:
			switch {
			case len(rssContents.List.Items) == 0:
				break higherSwitchK
			case state["rssContentCounter"] < 1:
			case len(rssContentsItems) < 1:
				break
			default:
				state["rssContentCounter"]--
			}
			stackClone = append(stackClone, rssContentsItems...)
			stackClone = highlightFocus(stackClone, state["rssContentCounter"])
			rssContents.List.Items = stackClone
			ui.Render(rssHeaders.List, rssContents.List)
		default:
		}

	})
	ui.Handle("o", func(ui.Event) {
		// open the current link in the browser
		switch {
		case state["focusStack"] == 1:
			focusString := rssContentsItems[state["rssContentCounter"]]
			link := fullStack[focusString]["Link"]
			openInBrowser(link)
		}
	})
	ui.Handle("H", func(ui.Event) {
		helpPageContent := []string{
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
		helpPage := makeListWidget(helpPageContent, "Help", *termWidth, *termHeight, 0, 0, "")
		ui.Render(helpPage.List)
		state["focusStack"] = 3 // set focusStack to 3
	})
	ui.Handle("<Tab>", func(ui.Event) {
		rssHeaders := makeListWidget(rssHeadersItems, "Feed", *termWidth/2, *termHeight, 0, 0, "Magenta")
		rssContents := makeListWidget(rssContentsItems, "Content", *termWidth/2, *termHeight, *termWidth/2, 0, "White")
		switch {
		case state["focusStack"] == 0:
			state["focusStack"] = 1
			rssHeaders.List.BorderFg = ui.ColorWhite
			rssContents.List.BorderFg = ui.ColorMagenta
			ui.Render(rssHeaders.List, rssContents.List)
		case state["focusStack"] == 1:
			state["focusStack"] = 0
			rssHeaders.List.BorderFg = ui.ColorMagenta
			rssContents.List.BorderFg = ui.ColorWhite
			ui.Render(rssHeaders.List, rssContents.List)
		}
	})

	ui.Handle("<Enter>", func(ui.Event) {
		rssHeaders := makeListWidget(rssHeadersItems, "Feed", *termWidth/2, *termHeight, 0, 0, "Magenta")
		rssContents := makeListWidget(rssContentsItems, "Content", *termWidth/2, *termHeight, *termWidth/2, 0, "White")
		switch {
		case state["focusStack"] == 0:
			if len(rssHeadersItems) < 1 {
				errorPage := makeParWidget(
					"There is currently no links",
					40, 3, *termWidth/2-(40/2), *termHeight, "",
				)
				ui.Render(errorPage.Par)
				time.Sleep(3 * time.Second) // show for three seconds
				ui.Clear()
				ui.Render(rssHeaders.List, rssContents.List)
				break
			}
			focusString := rssHeadersItems[state["rssHeaderCounter"]]
			// check if http is contained in side the string
			if strings.Contains(focusString, "http") == false {
				notValid := makeParWidget("Not a valid rss link",
					40, 3, *termWidth/2-(40/2), *termHeight/2, "Red")
				ui.Render(notValid.Par)
				time.Sleep(3 * time.Second)
				ui.Clear()
				ui.Render(rssHeaders.List, rssContents.List)
				break
			}
			fp := gofeed.NewParser()
			feed, _ := fp.ParseURL(focusString)
			items := feed.Items
			rssContentsItems = []string{} // clear out the current array
			for i := 0; i < len(items); i++ {
				rssContentsItems = append(rssContentsItems, items[i].Title)
				tempStack := make(map[string]string)
				tempStack["Description"] = items[i].Description
				tempStack["Published"] = items[i].Published
				tempStack["Link"] = items[i].Link
				fullStack[items[i].Title] = tempStack
			}
			rssContents.List.Items = rssContentsItems
			ui.Render(rssHeaders.List, rssContents.List)

		case state["focusStack"] == 1:
			// go to rssExtended
			extendedItems := []string{
				"Description : " +
					fullStack[rssContentsItems[state["rssContentCounter"]]]["Description"],
				"Published: " +
					fullStack[rssContentsItems[state["rssContentCounter"]]]["Published"],
				"Link: " +
					fullStack[rssContentsItems[state["rssContentCounter"]]]["Link"],
			}
			rssContentExtended := makeListWidget(extendedItems, "Content Extended",
				*termWidth, *termHeight, 0, 0, "default")
			ui.Render(rssContentExtended.List)
			state["focusStack"] = 2
		}
	})
	ui.Handle("<Escape>", func(ui.Event) {
		rssHeaders := makeListWidget(rssHeadersItems, "Feed", *termWidth/2, *termHeight, 0, 0, "Magenta")
		rssContents := makeListWidget(rssContentsItems, "Content", *termWidth/2, *termHeight, *termWidth/2, 0, "White")
		switch {
		case state["focusStack"] == 2: // rssContentExtended
			ui.Clear()
			ui.Render(rssHeaders.List, rssContents.List)
			state["focusStack"] = 1
		case state["focusStack"] == 3: // help page
			ui.Clear()
			ui.Render(rssHeaders.List, rssContents.List)
			state["focusStack"] = 0
		default: // dont do anything
		}

	})
	ui.Handle("<Resize>", func(e ui.Event) {
		*termWidth = ui.TermWidth()
		*termHeight = ui.TermHeight()

		rssHeaders := makeListWidget(rssHeadersItems, "Feed", *termWidth/2, *termHeight, 0, 0, "Magenta")
		rssContents := makeListWidget(rssContentsItems, "Content", *termWidth/2, *termHeight, *termWidth/2, 0, "White")

		ui.Render(rssHeaders.List, rssContents.List)

	})
	ui.Loop()
}
