package main

import (
	"AppUtils"
	"WindowsClient"

	"commands"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

const PvmVersion = "0.1.0.a1"

func main() {
	cmdInstall := &cobra.Command{
		Use:     "install <version>",
		Aliases: []string{"i"},
		Short:   "Installs the specified python version",
		Long:    `Installs the specified python version. If it is not a valid version, the program will exit`,
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch len(args) {
			case 1:
				commands.Install(args[0], "")
			case 2:
				commands.Install(args[0], args[1])
			default:
				log.Fatalln("Too many parameters. If you're trying to use an alias, please wrap it in double quotes (ex: \"my-custom-alias\")")
			}
		},
	}

	cmdReinstall := &cobra.Command{
		Use:     "reinstall <version>",
		Aliases: []string{"r"},
		Short:   "Reinstalls the specified python version",
		Long:    `Reinstalls the specified python version, that can either be a specific version, "all" or an alias. If not found, the program will exit`,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			commands.Reinstall(args[0])
		},
	}

	cmdUninstall := &cobra.Command{
		Use:     "uninstall <version>",
		Aliases: []string{"u"},
		Short:   "Uninstalls the specified python version",
		Long:    `Uninstalls the specified python version, that can either be a specific version, "all" or an alias. If not found, the program will exit`,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			commands.Uninstall(args[0])
		},
	}

	cmdUse := &cobra.Command{
		Use:   "use <slug>",
		Short: "Switch to use the specified alias",
		Long:  `Activates the specified python version. It can be either the version number or the installation alias. If not found, the program will exit`,
		Args:  cobra.ExactArgs(1),
		Run:   func(cmd *cobra.Command, args []string) { commands.Use(args[0]) },
	}

	cmdList := &cobra.Command{
		Use:     "list <mode>",
		Aliases: []string{"l"},
		Short:   "Lists all the mode-specified versions",
		Long:    `Lists all the mode-specified versions. Valid modes are "all" (lists stable and unstable versions), "installed", "latest" (lists the latest 5 releases for each major python version). If not a valid mode, the program will exit`,
		Args:    cobra.ExactArgs(1),
		Run:     func(cmd *cobra.Command, args []string) { commands.List(args[0]) },
	}

	cmdOn := &cobra.Command{
		Use:   "on",
		Short: "Enables python version management",
		Long:  `Enables python version management by creating a symlink`,
		Args:  cobra.ExactArgs(1),
		Run:   func(cmd *cobra.Command, args []string) { commands.ToggleAppState("enable") },
	}

	cmdOff := &cobra.Command{
		Use:   "off",
		Short: "Disables python version management",
		Long:  `Disables python version management by removing the symlink`,
		Args:  cobra.ExactArgs(1),
		Run:   func(cmd *cobra.Command, args []string) { commands.ToggleAppState("disable") },
	}

	cmdVersion := &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "Disables python version management",
		Long:    `Disables python version management by removing the symlink`,
		Args:    cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			client := WindowsClient.NewClient()

			fmt.Println("PVM (Python Version Manager) for Windows")
			fmt.Println("----------------------------------------")
			fmt.Println("Version: " + PvmVersion)
			fmt.Println("Arch:    " + WindowsClient.NewClient().Arch)
			fmt.Println("AppRoot: " + AppUtils.GetWorkingDir())
			fmt.Println("Client:  " + client.ClientInfo())
		},
	}

	var rootCmd = &cobra.Command{Use: "pvm"}

	rootCmd.AddCommand(
		cmdInstall,
		cmdReinstall,
		cmdUninstall,
		cmdUse,
		cmdList,
		cmdOn,
		cmdOff,
		cmdVersion,
	)

	rootCmd.Execute()
}
