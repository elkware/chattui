# ChatTUI

ChatTUI is a simple chat client for the terminal.

It is written in Go and uses the tview library for the terminal interface.

## Features

- Continuous conversation support
- Custom instructions
- Config and chat history persistence

![chattui.gif](chattui.gif)

## Installation

### From source

Clone the repository and build the binary.

```bash
cd chattui
go build -o chattui cmd/main.go
```

## Usage

To start the client, simply run the `chattui` command.

```bash
chattui
```

## Configuration

The UI provides a configuration interface to set the Open AI API key and the optional custom instructions.
