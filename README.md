# Read RSS

## How it looks
![Preview](preview.gif)

## Configuration - example
```
{
    "Rss": [
        "https://xkcd.com/rss.xml"
    ]
}
```

## Installation
```
git clone https://github.com/lczm/readrss
```
Ensure that you have dependencies installed (dependencies below)
Go Build/Run readrss.go and use it

## Hotkeys
"q : Exit"
"a : Add RSS feed"
"j : Move down"
"k : Move up"
"o : Open in browser"
"H : Help page [this]"
"Tab : Move focus between names and content"
"Enter : Get feed / Expand details"
"Esc : Escape out of current mode"
"Press Esc to get out of this screen"

## Dependencies
[gofeed](https://github.com/mmcdole/gofeed)
[termui](https://github.com/gizak/termui)
