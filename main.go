package main

import (
	updatechecker "github.com/Christian1984/go-update-checker"
	"log"
	"os"
	"os/exec"

	"github.com/bomkz/vtolvr-utils/definitions"
	"github.com/bomkz/vtolvr-utils/richpresence"

	"github.com/getlantern/systray"
)

var Version = "0.0"

func main() {
	//hsvrApp := app.New()
	//aircraft := fyne.NewStaticResource("aircraft", definitions.Icon)
	//hsvrApp.SetIcon(aircraft)

	//definitions.FrontendWindow = hsvrApp.NewWindow("HSVR API Frontend")

	filename := "hsvr-utils.log"
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

	uc := updatechecker.New("bomkz", "hsvr-utils", "HSVR Utilities", "https://github.com/bomkz/hsvr-utils/releases/latest", 0, false)
	uc.CheckForUpdate(Version)

	needsUpdate = uc.UpdateAvailable
	systray.Run(onReady, onExit)
	//hsvrApp.Run()

}

var needsUpdate bool

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
	//showFrontend := systray.AddMenuItem("Show Frontend", "Open Frontend GUI.")

	systray.AddSeparator()
	systray.AddMenuItem("Current version: "+Version, "")

	var c *systray.MenuItem

	if needsUpdate {
		c = systray.AddMenuItem("Update Available", "")
	}

	if exists {
		enableStartup.Check()
	}

	go richpresence.HandleInit()

	for {
		select {
		//case <-showFrontend.ClickedCh:
		//apifrontend.BuildFrontend()
		case <-c.ClickedCh:
			err := exec.Command("open", "https://github.com/bomkz/hsvr-utils/releases/latest").Run()
			if err != nil {
				log.Println(err)
			}
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

func onExit() {
	log.Fatal("Shutting down VTOL VR Utils")
}
