package PvmState

import (
	"AppUtils"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func Install() {
	isAdmin, err := AppUtils.IsAppRunningAsAdmin()

	if err != nil {
		panic(err)
	}

	if !isAdmin {
		fmt.Println("This command must be run as Admin")
		os.Exit(0)
	}

	// env := os.Environ()

	// for _, s := range env {
	// 	println(s)
	// }

	path, exists := os.LookupEnv("PATH")

	if !exists {
		log.Fatalln("PATH Enviroment variable doesn't exist")
	}

	pvm_home := AppUtils.GetWorkingDir()
	symlinkTarget := filepath.Join(os.Getenv("localappdata"), "Python")

	newPath := fmt.Sprintf("%s;%s;%s;%s\\Scripts", path, pvm_home, symlinkTarget, symlinkTarget)

	AppUtils.RunCmd(fmt.Sprintf(`powershell -Command "[Environment]::SetEnvironmentVariable("PATH", %s, [EnvironmentVariableTarget]::Machine)`, newPath))

}
