package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/knightpp/sni/pkg/menu"
	"github.com/knightpp/sni/pkg/tray"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	tree := menu.NewItem().Submenu(
		menu.NewItem().Label("Button 1").IconName("emblem-default").
			OnClick(func() {
				log.Print("Button 1 clicked!!!")
			}),
		menu.NewItem().Separator(true),
		menu.NewItem().Label("Button 2").IconName("help-about").
			OnClick(func() {
				log.Print("Button 2 clicked!!!")
			}),
	).Build()

	tray, err := tray.NewTray("MyApp", "Descriptive title", tree)
	if err != nil {
		return err
	}
	defer tray.Close()

	err = tray.Setup()
	if err != nil {
		return err
	}

	// go func() {
	// 	var err error
	// 	for {
	// 		tray.SetIconName("emblem-mail")
	// 		log.Print("Changed to: ", tray.GetIconName())
	// 		err = tray.SignalNewIcon()
	// 		if err != nil {
	// 			log.Print("Error: ", err)
	// 			break
	// 		}
	// 		time.Sleep(2 * time.Second)
	// 		tray.SetIconName("emblem-default")
	// 		log.Print("Changed to: ", tray.GetIconName())
	// 		err = tray.SignalNewIcon()
	// 		if err != nil {
	// 			log.Print("Error: ", err)
	// 			break
	// 		}
	// 		time.Sleep(2 * time.Second)
	// 	}
	// }()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Print("SIG INTERRUPT (Ctrl + C): exitting")
	return nil
}
