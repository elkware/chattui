package main

import (
	"encoding/json"
	"os"

	"chattui/internal"
	"chattui/internal/ui"
	tv "github.com/rivo/tview"
)

func main() {
	loadChatHistory()
	app := tv.NewApplication()
	defer saveChats()
	defer app.Stop()
	app.EnableMouse(true)

	root := tv.NewFlex()
	root.SetDirection(tv.FlexRow)
	root.SetBorder(true)

	mainUI := ui.NewMainUI(app, root)
	mainUI.SetupMouseListeners()
	mainUI.SetupKeyboardListeners()
	app.SetRoot(root, true)

	defer func() {
		if r := recover(); r != nil {
			if app != nil {
				app.Stop()
			}
		}
	}()

	// Run the application.
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func loadChatHistory() {
	// load chats from file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	cfg, err := os.ReadFile(homeDir + "/.chattui-history.json")
	if err != nil {
		if os.IsNotExist(err) {
			// Return an empty config
			internal.Chats = make([]*internal.Chat, 0)
			return
		} else {
			panic(err)
		}
	}
	if err := json.Unmarshal(cfg, &internal.Chats); err != nil {
		panic(err)
	}
}

func saveChats() {
	// save chats to file
	ret, err := json.Marshal(internal.Chats)
	if err != nil {
		panic(err)
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(homeDir+"/.chattui-history.json", ret, 0644); err != nil {
		panic(err)
	}
}
