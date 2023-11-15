package ui

import (
	"fmt"

	"chattui/internal"
	"chattui/internal/config"
	"github.com/gdamore/tcell/v2"
	tv "github.com/rivo/tview"
)

type MainUI struct {
	NewChatBtn    *tv.Button
	ClearChatsBtn *tv.Button
	ChatList      *tv.List
	ChatBox       *tv.TextView
	ChatInput     *tv.InputField
	ExitBtn       *tv.Button
	ConfigBtn     *tv.Button
	ModelDropdown *tv.DropDown
	ExitModal     *tv.Modal
	root          *tv.Flex
	app           *tv.Application
	configScreen  *ConfigUI
	focusBoxes    []*tv.Box
	currentFocus  int
}

func (c *MainUI) Init() {
	c.ModelDropdown = tv.NewDropDown()
	c.ModelDropdown.SetLabel("Model: ")
	c.ModelDropdown.SetOptions([]string{"gpt-3.5-turbo", "gpt-4"}, nil)
	c.ModelDropdown.SetCurrentOption(0)

	c.ClearChatsBtn = tv.NewButton("Clear all chats")
	c.ClearChatsBtn.SetMouseCapture(func(action tv.MouseAction, event *tcell.EventMouse) (tv.MouseAction, *tcell.EventMouse) {
		if action == tv.MouseLeftClick && c.ClearChatsBtn.InRect(event.Position()) {
			modal := tv.NewModal()
			modal.SetText("Are you sure you want to clear all chats?")
			modal.AddButtons([]string{"Yes", "No"})
			modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "Yes" {
					internal.Chats = make([]*internal.Chat, 0)
					c.ChatList.Clear()
					c.ChatBox.SetText("")
					c.ChatInput.SetText("")
					c.ChatInput.SetPlaceholder(" Ready.")
					c.ChatInput.SetDisabled(false)
					c.ModelDropdown.SetDisabled(false)
					c.ModelDropdown.SetCurrentOption(0)
				}
				c.app.SetRoot(c.root, true)
			})
			c.app.SetRoot(modal, true)
		}
		return action, event
	})

	c.NewChatBtn = tv.NewButton("New Chat")
	c.ChatList = tv.NewList()
	c.ChatList.SetBorder(true)
	c.ChatList.SetHighlightFullLine(true)
	c.ChatList.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		messages := ""
		for _, msg := range internal.Chats[index].Messages {
			if msg.Role == "user" {
				messages += "[blue]" + internal.GetUsernameFromOS() + ":[white]\n" + msg.Content + "\n"
			} else {
				messages += "[green]AI:[white]\n" + msg.Content + "\n"
			}
			messages += "\n"
		}
		c.ChatBox.SetText(messages)
		if internal.Chats[index].Model == "" {
			c.ModelDropdown.SetDisabled(false)
			c.ModelDropdown.SetCurrentOption(0)
		} else {
			c.ModelDropdown.SetDisabled(true)
			if internal.Chats[index].Model == "gpt-3.5-turbo" {
				c.ModelDropdown.SetCurrentOption(0)
			} else {
				c.ModelDropdown.SetCurrentOption(1)
			}
		}
	})
	for _, chat := range internal.Chats {
		c.ChatList.AddItem(chat.Name, "", 0, nil)
	}
	if len(internal.Chats) > 0 {
		c.ChatList.SetCurrentItem(len(internal.Chats) - 1)
		if internal.Chats[len(internal.Chats)-1].Model == "" {
			c.ModelDropdown.SetDisabled(false)
			c.ModelDropdown.SetCurrentOption(0)
		} else {
			c.ModelDropdown.SetDisabled(true)
			if internal.Chats[len(internal.Chats)-1].Model == "gpt-3.5-turbo" {
				c.ModelDropdown.SetCurrentOption(0)
			} else {
				c.ModelDropdown.SetCurrentOption(1)
			}
		}
	}

	c.ChatBox = tv.NewTextView()
	c.ChatBox.SetBorder(true)
	c.ChatBox.SetScrollable(true)
	c.ChatBox.SetDynamicColors(true)
	c.ChatBox.SetChangedFunc(func() {
		c.app.Draw()
	},
	)
	if len(internal.Chats) > 0 {
		messages := ""
		for _, msg := range internal.Chats[len(internal.Chats)-1].Messages {
			if msg.Role == "user" {
				messages += "[blue]" + internal.GetUsernameFromOS() + ":[white]\n" + msg.Content + "\n"
			} else {
				messages += "[green]AI:[white]\n" + msg.Content + "\n"
			}
			messages += "\n"
		}
		c.ChatBox.SetText(messages)
	}

	c.ChatInput = tv.NewInputField()
	c.ChatInput.SetBorder(true)
	c.ChatInput.SetPlaceholder(" Ready.")
	c.ChatInput.SetLabel(internal.GetUsernameFromOS() + "> ")
	c.ChatInput.SetRect(0, 0, 0, 3)

	c.ExitBtn = tv.NewButton("X")
	c.ConfigBtn = tv.NewButton("Config")

	c.ExitModal = tv.NewModal()
	c.ExitModal.SetText("Are you sure you want to quit?")
	c.ExitModal.AddButtons([]string{"Yes", "No"})
	c.ExitModal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Yes" {
			c.app.Stop()
		}
		if buttonLabel == "No" {
			c.app.SetRoot(c.root, true)
		}
	},
	)

	c.focusBoxes = []*tv.Box{
		c.ExitBtn.Box, c.ConfigBtn.Box, c.ModelDropdown.Box,
		c.ChatList.Box, c.NewChatBtn.Box, c.ChatBox.Box, c.ChatInput.Box,
	}
	c.currentFocus = 6

	c.configScreen = NewConfigUI(c.app, c.root)

	// layout for the menu bar
	exitAndCfgLayout := tv.NewFlex().SetDirection(tv.FlexColumn)
	exitAndCfgLayout.AddItem(tv.NewFrame(nil), 1, 1, false)
	exitAndCfgLayout.AddItem(c.ExitBtn, 3, 1, false)
	exitAndCfgLayout.AddItem(tv.NewFrame(nil), 3, 1, false)
	exitAndCfgLayout.AddItem(c.ConfigBtn, 10, 1, false)
	exitAndCfgLayout.AddItem(tv.NewFrame(nil), 3, 1, false)
	exitAndCfgLayout.AddItem(c.ModelDropdown, 0, 1, false)

	// layout for the chat list and chat box
	chatLayout := tv.NewFlex().SetDirection(tv.FlexColumn)
	left := tv.NewFlex().SetDirection(tv.FlexRow)

	ccBtn := tv.NewFlex().SetDirection(tv.FlexRow).AddItem(c.ClearChatsBtn, 1, 1, false)
	ccBtn.SetBorder(true)
	left.AddItem(ccBtn, 3, 2, false)

	left.AddItem(c.ChatList, 0, 1, false)
	btn := tv.NewFlex().SetDirection(tv.FlexRow).AddItem(c.NewChatBtn, 1, 1, false)
	btn.SetBorder(true)
	left.AddItem(btn, 3, 2, false)

	right := tv.NewFlex().SetDirection(tv.FlexRow)
	right.AddItem(c.ChatBox, 0, 1, false)
	right.AddItem(c.ChatInput, 3, 1, false)
	chatLayout.AddItem(left, 0, 1, false)
	chatLayout.AddItem(right, 0, 4, false)

	// add the menu bar and chat layout to the root layout
	c.root.AddItem(exitAndCfgLayout, 1, 1, false)
	c.root.AddItem(chatLayout, 0, 1, true)
	c.app.SetFocus(c.focusBoxes[c.currentFocus])
}

func NewMainUI(app *tv.Application, root *tv.Flex) *MainUI {
	ui := &MainUI{}
	ui.root = root
	ui.app = app
	ui.Init()
	return ui
}

func (c *MainUI) SetupMouseListeners() {
	c.ExitBtn.SetMouseCapture(func(action tv.MouseAction, event *tcell.EventMouse) (tv.MouseAction, *tcell.EventMouse) {
		if action == tv.MouseLeftClick && c.ExitBtn.InRect(event.Position()) {
			c.app.SetRoot(c.ExitModal, true)
		}
		return action, event
	})

	c.ConfigBtn.SetMouseCapture(func(action tv.MouseAction, event *tcell.EventMouse) (tv.MouseAction, *tcell.EventMouse) {
		if action == tv.MouseLeftClick && c.ConfigBtn.InRect(event.Position()) {
			config.AppConfig.Load()
			(c.configScreen.Form.GetFormItemByLabel("Open AI API Key")).(*tv.InputField).SetText(config.AppConfig.ApiKey)
			(c.configScreen.Form.GetFormItemByLabel("Custom instructions")).(*tv.TextArea).SetText(config.AppConfig.CustomInstructions, true)
			c.app.SetRoot(c.configScreen.FormRoot, true)
		}
		return action, event
	})
	c.NewChatBtn.SetMouseCapture(func(action tv.MouseAction, event *tcell.EventMouse) (tv.MouseAction, *tcell.EventMouse) {
		if action == tv.MouseLeftClick && c.NewChatBtn.InRect(event.Position()) {
			c.newChatBtnClick()
		}
		return action, event
	})

}

func (c *MainUI) newChatBtnClick() {
	internal.Chats = append(internal.Chats, &internal.Chat{
		Name:     fmt.Sprint("Chat ", len(internal.Chats)+1),
		Model:    "",
		Messages: make([]internal.ChatMessage, 0),
	})
	c.ChatList.AddItem(fmt.Sprint("Chat ", len(internal.Chats)), "", 0, nil)
	c.ChatList.SetCurrentItem(len(internal.Chats) - 1)
	c.ChatBox.SetText("")
	c.ChatInput.SetPlaceholder(" Ready.")
	c.ChatInput.SetText("")
	c.ModelDropdown.SetDisabled(false)
	c.ModelDropdown.SetCurrentOption(0)
}

func (c *MainUI) SetupKeyboardListeners() {
	c.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			if c.currentFocus == 7 {
				c.currentFocus = 0
			}
			c.app.SetFocus(c.focusBoxes[c.currentFocus])
			c.currentFocus++
			return nil
		}
		if event.Key() == tcell.KeyBacktab {
			if c.currentFocus < 0 {
				c.currentFocus = 6
			}
			c.app.SetFocus(c.focusBoxes[c.currentFocus])
			c.currentFocus--
			return nil
		}
		return event
	})

	c.NewChatBtn.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			c.newChatBtnClick()
		}
		return event
	})

	c.ExitBtn.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			c.app.SetRoot(c.ExitModal, true)
		}
		return event
	})

	c.setupChatFunc()

}

func (c *MainUI) setupChatFunc() {
	c.ChatInput.SetFinishedFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			input := c.ChatInput.GetText()
			if input == "" {
				return
			}
			c.ChatInput.SetDisabled(true)
			if len(internal.Chats) == 0 {
				internal.Chats = append(internal.Chats, &internal.Chat{
					Name:     fmt.Sprint("Chat ", len(internal.Chats)+1),
					Model:    "",
					Messages: make([]internal.ChatMessage, 0),
				})
				c.ChatList.AddItem(fmt.Sprint("Chat ", len(internal.Chats)), "", 0, nil)
				c.ChatList.SetCurrentItem(0)
			}
			chat := internal.Chats[c.ChatList.GetCurrentItem()]
			if chat.Model == "" {
				_, chat.Model = c.ModelDropdown.GetCurrentOption()
			}
			c.ModelDropdown.SetDisabled(true)

			chat.Messages = append(chat.Messages, internal.ChatMessage{
				Role:    "user",
				Content: input,
			})

			c.ChatBox.SetText(c.ChatBox.GetText(false) + "[blue]" + internal.GetUsernameFromOS() + ":[white]\n" + input + "\n")
			go func() {
				c.ChatInput.SetText("")
				c.ChatInput.SetPlaceholder(" Thinking...")
				resp, err := internal.CallChatGPT(chat)
				if err != nil {
					c.ChatBox.SetText(c.ChatBox.GetText(false) + "\n[red]AI:\n[white]An error occurred while calling the API.\n\n")
					c.ChatBox.ScrollToEnd()
					c.ChatInput.SetPlaceholder(" Ready.")
					c.ChatInput.SetDisabled(false)
					return
				}
				chat.Messages = append(chat.Messages, internal.ChatMessage{
					Role:    "assistant",
					Content: resp,
				})
				response := "\n[green]AI:\n[white]" + resp + "\n\n"
				c.ChatBox.SetText(c.ChatBox.GetText(false) + response)

				c.ChatBox.ScrollToEnd()
				c.ChatInput.SetPlaceholder(" Ready.")
				c.ChatInput.SetDisabled(false)
			}()
		}
	},
	)
}
