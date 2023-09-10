package commands

import (
	windowsClient "WindowsClient"

	"fmt"
)

func List(mode string) {
	client := windowsClient.NewClient()

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
