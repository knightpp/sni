package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	codegen := func(p string, args ...string) error {
		filenameExt := path.Base(p)
		filenameSnake := SnakeCase(strings.TrimSuffix(filenameExt, ".xml"))
		baseDir := path.Join("../generated/", filenameSnake)
		outputFile := path.Join(baseDir, filenameSnake+".go")
		_ = os.MkdirAll(baseDir, 0o750)
		cmd := exec.Command("./dbus-codegen-go")
		cmd.Args = append(cmd.Args,
			"--output="+outputFile,
			"--package="+filenameSnake)
		cmd.Args = append(cmd.Args, args...)
		cmd.Args = append(cmd.Args, p)
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	params := map[string][]string{
		"./dbus_interfaces/DBusMenu.xml": {
			"--prefix=com.canonical",
		},
		"./dbus_interfaces/StatusNotifierItem.xml": {
			"--prefix=org.kde",
		},
		"./dbus_interfaces/StatusNotifierWatcher.xml": {
			"--prefix=org.kde",
		},
		"./dbus_interfaces/DBus.xml": {
			"--prefix=org.freedesktop",
		},
	}
	for path, args := range params {
		log.Printf("Executing with path: %s", path)
		if err := codegen(path, args...); err != nil {
			return err
		}
	}

	return nil
}

func SnakeCase(camel string) string {
	var buf bytes.Buffer
	for _, c := range camel {
		if 'A' <= c && c <= 'Z' {
			// just convert [A-Z] to _[a-z]
			if buf.Len() > 0 {
				buf.WriteRune('_')
			}
			buf.WriteRune(c - 'A' + 'a')
		} else {
			buf.WriteRune(c)
		}
	}
	return buf.String()
}
