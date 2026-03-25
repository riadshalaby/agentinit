package prereq

import (
	"os/exec"
	"runtime"
)

type OS string

const (
	Darwin  OS = "darwin"
	Linux   OS = "linux"
	Windows OS = "windows"
)

type PackageManager struct {
	Name           string
	Installed      bool
	SelfInstallCmd string
}

var runtimeGOOS = runtime.GOOS

func DetectOS() OS {
	return detectOS(runtimeGOOS)
}

func DetectPackageManager(o OS) PackageManager {
	return detectPackageManager(o, exec.LookPath)
}

func detectOS(goos string) OS {
	switch goos {
	case string(Darwin):
		return Darwin
	case string(Windows):
		return Windows
	default:
		return Linux
	}
}

func detectPackageManager(o OS, lookPath func(string) (string, error)) PackageManager {
	switch o {
	case Darwin:
		_, err := lookPath("brew")
		return PackageManager{
			Name:           "brew",
			Installed:      err == nil,
			SelfInstallCmd: `/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`,
		}
	case Windows:
		_, err := lookPath("choco")
		return PackageManager{
			Name:           "choco",
			Installed:      err == nil,
			SelfInstallCmd: "powershell -NoProfile -ExecutionPolicy Bypass -Command \"Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))\"",
		}
	default:
		return PackageManager{}
	}
}
