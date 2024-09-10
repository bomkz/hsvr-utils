package main

import (
	"fmt"
	"github.com/bomkz/hsvr-utils/definitions"
	"github.com/bomkz/hsvr-utils/richpresence"
	"github.com/getlantern/systray"
	"log"
	"os"
)

func main() {
	//hsvrApp := app.New()
	//aircraft := fyne.NewStaticResource("aircraft", definitions.Icon)
	//hsvrApp.SetIcon(aircraft)

	//definitions.FrontendWindow = hsvrApp.NewWindow("HSVR API Frontend")

	filename := "hsvr-utils.log"
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return
	}
	createFileIfNotExists()
	file, err := openLogFile(homedir + "\\" + filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	log.Println("log file created")

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

	if exists {
		enableStartup.Check()
	}

	go richpresence.HandleInit()

	for {
		select {
		//case <-showFrontend.ClickedCh:
		//apifrontend.BuildFrontend()

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
