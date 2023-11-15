# ChatTUI

ChatTUI is a simple chat client for the terminal.

It is written in Go and uses the tview library for the terminal interface.

## Features

- Continuous conversation support
- Custom instructions
- Config and chat history persistence

![chattui.gif](chattui.gif)

## Installation

```bash
go get github.com/elkware/chattui
```

## Usage

To start the client, simply run the `chattui` command.

```bash
chattui
```

## Configuration

The UI provides a configuration interface to set the Open AI API key and the optional custom instructions.
