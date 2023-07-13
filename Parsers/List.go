package Parsers

import (
	windowsClient "WindowsClient"

	"fmt"
)

func ListParserHandler(mode string, client *windowsClient.Client) {
	switch mode {
	case "all":
		client.ListAll()
	case "latest":
		client.ListLatest()
	case "installed":
		client.ListInstalled()
	default:
		fmt.Printf("\"%s\" is an invalid mode\n", mode)
	}
}
