package main

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/jxeng/shortcut"
)

func main() {
	filename := "vtolvrutil.log"
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	createFileIfNotExists()
	file, err := openLogFile(homedir + "\\" + filename)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	log.Println("log file created")
	systray.Run(onReady, onExit)
}

func createFileIfNotExists() {
	filename := "vtolvrutil.log"
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
		file.Close()

	}
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

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("VTOL VR Utilities")
	systray.SetTooltip("VTOLVR Utils")
	quit := systray.AddMenuItem("Quit", "Quit App")

	exists, err := checkIfStartupExists()

	if err != nil {
		log.Fatal(err)
		return
	}

	enableStartup := systray.AddMenuItemCheckbox("Start on boot", "Start the app when you log in.", false)

	if exists {
		enableStartup.Check()
	}
	steamID32 := findCurrentUID()

	int64SteamID64 := convertID3ToID64(steamID32)

	steamID64 = strconv.Itoa(int(int64SteamID64))

	go connectWS()

	go richPresenceHandler()

	quit.SetIcon(icon.Data)

	for {
		select {
		case <-quit.ClickedCh:
			onExit()
		case <-enableStartup.ClickedCh:

			exists, err = checkIfStartupExists()

			if err != nil {
				log.Fatal(err)
			}

			if !exists {
				err := makeLink()

				if err != nil {
					log.Fatal(err)
				}
				enableStartup.Check()
			} else if exists {
				err := deleteLink()

				if err != nil {
					log.Fatal(err)
				}
				enableStartup.Uncheck()
			}

		}

	}

}

func deleteLink() error {
	appdata, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	startupDir := appdata + "\\Microsoft\\Windows\\Start Menu\\Programs\\Startup\\vtolvr.lnk"

	os.Remove(startupDir)

	return nil

}

func onExit() {
	log.Fatal("Shutting down VTOL VR Utils")
}
