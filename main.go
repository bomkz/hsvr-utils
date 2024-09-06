package main

import (
	"errors"
	"log"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/bomkz/vtolvr-utils/apifrontend"
	"github.com/bomkz/vtolvr-utils/definitions"
	"github.com/bomkz/vtolvr-utils/richpresence"

	"github.com/getlantern/systray"

	"github.com/jxeng/shortcut"
)

func main() {
	hsvrApp := app.New()
	aircraft := fyne.NewStaticResource("aircraft", definitions.Icon)
	hsvrApp.SetIcon(aircraft)

	definitions.FrontendWindow = hsvrApp.NewWindow("HSVR API Frontend")

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
	go systray.Run(onReady, onExit)
	hsvrApp.Run()

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
		err := file.Close()
		if err != nil {
			return
		}

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
	systray.SetIcon(definitions.Icon)
	systray.SetTitle("VTOL VR Utilities")
	systray.SetTooltip("VTOLVR Utils")

	quit := systray.AddMenuItem("Quit", "Quit App")

	exists, err := checkIfStartupExists()

	if err != nil {
		log.Fatal(err)
		return
	}

	enableStartup := systray.AddMenuItemCheckbox("Start on boot", "Start the app when you log in.", false)
	showFrontend := systray.AddMenuItem("Show Frontend", "Open Frontend GUI.")

	if exists {
		enableStartup.Check()
	}
	steamID32 := findCurrentUID()

	int64SteamID64 := convertID3ToID64(steamID32)

	definitions.SteamID64 = strconv.Itoa(int(int64SteamID64))

	go richpresence.InitRP()

	for {
		select {
		case <-showFrontend.ClickedCh:
			apifrontend.BuildFrontend()
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

	err = os.Remove(startupDir)
	if err != nil {
		return err
	}

	return nil

}

func onExit() {
	log.Fatal("Shutting down VTOL VR Utils")
}
