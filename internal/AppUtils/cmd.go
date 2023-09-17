package AppUtils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func RunCmd(command string) (cmd *exec.Cmd, err error) {
	osShell := []string{"bash", "-c"}

	if runtime.GOOS == "windows" {
		osShell = []string{"cmd", "/C"}
	}

	cmd = exec.Command(osShell[0], osShell[1], command)
	// cmd.Stdout = os.Stdout
	cmd.Stdout = nil
	cmd.Stderr = os.Stdout
	err = cmd.Run()

	if err != nil {
		return cmd, err
	}

	return cmd, nil
}

func CreateSymlink(src string, dest string) error {
	// fmt.Println("\nsrc : " + src)
	// fmt.Println("dest: " + dest)

	// run command
	symlinkCommand := fmt.Sprintf("New-Item -Force -ItemType SymbolicLink -Path '%s' -Target '%s'", dest, src)
	command := []string{"powershell.exe", "-noprofile", `Start-Process -WindowStyle hidden -Verb RunAs -Wait powershell.exe -Args "` + symlinkCommand + `"`}

	_, err := RunCmd(strings.Join(command, " "))

	if err != nil {
		return err
	}

	return nil
}

func GetPythonVersionInUse() (string, error) {
	cmd := exec.Command("python", "-V")
	out, err := cmd.CombinedOutput()

	if err != nil {
		return "", err
	}

	output := strings.TrimSpace(string(out))
	res := strings.Split(output, " ")[1]

	return res, nil
}
