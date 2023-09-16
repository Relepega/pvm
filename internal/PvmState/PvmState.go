package PvmState

import (
	"AppUtils"

	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func Install() {
	envVars, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.ALL_ACCESS) // Enviroment Variable for <username>
	if err != nil {
		log.Fatalln(err)
	}
	defer envVars.Close()

	// get path
	pathContent, _, err := envVars.GetStringValue("PATH")
	if err != nil {
		log.Fatalln(err)
	}

	// chech if content is already there
	symlinkTarget := filepath.Join(os.Getenv("localappdata"), "Python")
	append := fmt.Sprintf("%s;%s;%s\\Scripts", AppUtils.GetWorkingDir(), symlinkTarget, symlinkTarget)

	if strings.Contains(pathContent, append) {
		fmt.Println("PVM is already installed")
		return
	}

	// append it
	pathContent += ";" + append

	err = envVars.SetStringValue("PATH", pathContent)
	if err != nil {
		log.Fatalln(err)
	}

	// 🎉
	fmt.Println("PVM is now installed, restart your shell to apply the changes!")
}

func Uninstall() {
	envVars, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.ALL_ACCESS) // Enviroment Variable for <username>
	if err != nil {
		log.Fatalln(err)
	}
	defer envVars.Close()

	// get path
	pathContent, _, err := envVars.GetStringValue("PATH")
	if err != nil {
		log.Fatalln(err)
	}

	// chech if content is already there
	symlinkTarget := filepath.Join(os.Getenv("localappdata"), "Python")
	substr := fmt.Sprintf("%s;%s;%s\\Scripts", AppUtils.GetWorkingDir(), symlinkTarget, symlinkTarget)

	if !strings.Contains(pathContent, substr) {
		fmt.Println("PVM is not installed")
		return
	}

	pathContent = strings.Replace(pathContent, substr, "", 1)

	// append it
	err = envVars.SetStringValue("PATH", pathContent)
	if err != nil {
		log.Fatalln(err)
	}

	// 🎉
	fmt.Println("PVM is now uninstalled, restart your shell to apply the changes!")
}
