package PvmState

import (
	"AppUtils"
	"WindowsClient"
	"os/exec"
	"path/filepath"

	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func Handler(mode string) {
	registry, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.ALL_ACCESS) // Enviroment Variable for <username>
	if err != nil {
		log.Fatalln(err)
	}
	defer registry.Close()

	pathContent, _, err := registry.GetStringValue("PATH")
	if err != nil {
		log.Fatalln(err)
	}

	symlinkTarget := strings.ReplaceAll(WindowsClient.SymlinkDest, "%localappdata%", os.Getenv("localappdata"))
	pvmInstallPaths := fmt.Sprintf("%s;%s;%s\\Scripts", AppUtils.GetWorkingDir(), symlinkTarget, symlinkTarget)

	switch mode {
	case "install":
		install(&registry, &pathContent, &pvmInstallPaths)
	case "uninstall":
		uninstall(&registry, &pathContent, &pvmInstallPaths)
	case "enable":
		enable(&symlinkTarget)
	case "disable":
		disable(&symlinkTarget)
	default:
		log.Fatalf(`"%s" is not a valid mode`, mode)
	}
}

func install(registry *registry.Key, pathContent *string, pvmInstallPaths *string) {
	if strings.Contains(*pathContent, *pvmInstallPaths) {
		fmt.Println("PVM is already installed")
		return
	}

	splitted := strings.Split(*pathContent, ";")

	// do some work for fixing msys2 conflicts (issue #2)
	// maybe it's inefficient af, but works :)
	newPathContent := make([]string, len(splitted)+1)

	newPathContent[0] = splitted[0]
	newPathContent[1] = *pvmInstallPaths

	for i, s := range splitted {
		if i == 0 {
			continue
		}
		newPathContent[i+1] = s
	}

	err := registry.SetStringValue("PATH", strings.Join(newPathContent, ";"))
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("PVM is now installed, restart your shell to apply the changes!")
}

func uninstall(registry *registry.Key, pathContent *string, pvmInstallPaths *string) {
	if !strings.Contains(*pathContent, *pvmInstallPaths) {
		fmt.Println("PVM is not installed")
		return
	}

	err := registry.SetStringValue("PATH", strings.Replace(*pathContent, *pvmInstallPaths, "", 1))
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("PVM is now uninstalled, restart your shell to apply the changes!")
}

func enable(symlinkTarget *string) {
	f := filepath.Clean(*symlinkTarget)

	_, err := os.Lstat(f + ".pvmbak")
	if err != nil {
		if _, err := os.Lstat(f); err == nil {
			fmt.Println("PVM-managed python installations are already enabled!")
			os.Exit(0)
		}

		fmt.Println(`No disable state has been found. Please use "pvm use <version>" or "pvm install <version>" instead`)
		os.Exit(0)
	}

	cmd := exec.Command("cmd.exe", "/C", "move", f+".pvmbak", f)
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}

	// re-check to confirm it has been renamed
	// if not, that means the symlink is invalid
	// and thus, we delete it
	_, err = os.Lstat(f + ".pvmbak")
	if err == nil {
		fmt.Println("A PVM disable state has been found, but the link to the python installation is invalid")
		fmt.Println("Now removing it for security reasons")
		fmt.Println(`Please use "pvm use <version>" or "pvm install <version>" instead`)
		os.Remove(f + ".pvmbak")
		os.Exit(0)
	}

	fmt.Println("PVM-managed python installations are now enabled!")
}

func disable(symlinkTarget *string) {
	f := filepath.Clean(*symlinkTarget)

	_, err := os.Lstat(f + ".pvmbak")
	if err == nil {
		fmt.Println("PVM-managed python installations are already disabled!")
		os.Exit(0)
	}

	_, err = os.Lstat(f)
	if err == nil {
		cmd := exec.Command("cmd.exe", "/C", "move", f, f+".pvmbak")
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		err = cmd.Wait()
		if err != nil {
			log.Fatal(err)
		}
	}

	// re-check to confirm it has been renamed
	// if not, that means the symlink is invalid
	// and thus, we delete it
	_, err = os.Lstat(f)
	if err == nil {
		os.Remove(f)
	}

	fmt.Println("PVM-managed python installations are now disabled")
	fmt.Println(`To re-enable them either use "pvm on", "pvm use <version>" or "pvm install <version>"`)
}
