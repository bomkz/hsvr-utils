//go:generate goversioninfo -icon=definitions/aircraft.ico -manifest=definitions/hsvr-utils.exe.manifest
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bomkz/hsvr-utils/definitions"
	"github.com/bomkz/hsvr-utils/richpresence"
	"github.com/getlantern/systray"
)

func main() {

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

	if exists {
		enableStartup.Check()
	}

	go richpresence.HandleInit()

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

func onExit() {
	log.Fatal("Shutting down VTOL VR Utils")
}
