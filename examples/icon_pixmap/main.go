package main

import (
	"encoding/base64"
	"image/png"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/knightpp/sni"
	"github.com/knightpp/sni/menu"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
func run() error {
	tree := menu.NewItem().Build()
	tray, err := sni.NewTray("MyApp", "Descriptive title", tree)
	if err != nil {
		return err
	}
	defer tray.Close()
	imgBase64 := "iVBORw0KGgoAAAANSUhEUgAAAEAAAABACAYAAACqaXHeAAAAZUlEQVR42" +
		"u3QQREAAAQAMDIJqLeXHM4WYRlTHY+lAAECBAgQIECAAAECBAgQIECAAAECBAgQIECAAAE" +
		"CBAgQIECAAAECBAgQIECAAAECBAgQIECAAAECBAgQIECAAAECBAgQcN8CKKVrQf7diS4AA" +
		"AAASUVORK5CYII="

	img, err := png.Decode(
		base64.NewDecoder(base64.StdEncoding, strings.NewReader(imgBase64)))
	if err != nil {
		return err
	}
	tray.SetIconName("")
	// tray.SetIconName("help-about")
	tray.SetIconPixmap(img)
	err = tray.Setup()
	if err != nil {
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Print("SIG INTERRUPT (Ctrl + C): exitting")
	return nil
}
