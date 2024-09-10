package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/jxeng/shortcut"
)

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

func deleteLink() error {
	appdata, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	startupDir := appdata + "\\Microsoft\\Windows\\Start Menu\\Programs\\Startup\\vtolvr.lnk"

	err = os.Remove(startupDir)
	if err != nil {
		return err
	}

	return nil

}

func openLogFile(path string) (*os.File, error) {
	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return logFile, nil
}

func makeLink() error {

	ex, err := os.Executable()

	if err != nil {
		log.Fatal(err)
	}

	appdata, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}

	startupDir := appdata + "\\Microsoft\\Windows\\Start Menu\\Programs\\Startup\\vtolvr.lnk"

	sc := shortcut.Shortcut{
		ShortcutPath:     startupDir,
		Target:           ex,
		IconLocation:     "%SystemRoot%\\System32\\SHELL32.dll,0",
		Arguments:        "",
		Description:      "",
		Hotkey:           "",
		WindowStyle:      "1",
		WorkingDirectory: "",
	}

	err = shortcut.Create(sc)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func createFileIfNotExists() {
	filename := "hsvr-util.log"
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
		return
	}
	file, err := os.OpenFile(homedir+"\\"+filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		if os.IsExist(err) {
			return
		} else {
			log.Fatal(err)
		}
	} else {
		err := file.Close()
		if err != nil {
			return
		}

	}
}

func checkIfStartupExists() (bool, error) {

	appdata, err := os.UserConfigDir()
	if err != nil {
		return false, err
	}

	startupDir := appdata + "\\Microsoft\\Windows\\Start Menu\\Programs\\Startup\\vtolvr.lnk"

	if _, err := os.Stat(startupDir); err == nil {
		return true, nil

	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil

	} else {
		return false, err

	}
}
