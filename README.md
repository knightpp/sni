
# sni / StatusNotifierItem

This library implements `org.kde.StatusNotifierItem` and `com.canonical.dbusmenu` specs.
That enables us to create tray icons with menus.

## Used resources
- https://git.sailfishos.org/mer-core/qtbase/commit/38abd653774aa0b3c5cdfd9a8b78619605230726
- https://www.freedesktop.org/wiki/Specifications/StatusNotifierItem/
## License

[MIT](https://choosealicense.com/licenses/mit/)

  
## Usage/Examples

```go
package main

import (
	"log"
	"os"
	"os/signal"
	"sni"
	"sni/menu"
)

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
	tray, err := sni.NewTray("MyApp", "Descriptive title", tree)
	if err != nil {
		return err
	}
	defer tray.Close()
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
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
```

  
## Acknowledgements

 - [A Rust implementation of the KDE/freedesktop StatusNotifierItem specification ](https://github.com/iovxw/ksni)
  