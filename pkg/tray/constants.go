package tray

import (
	"fmt"
	"os"
)

type (
	TextDirection = string
	MenuStatus    = string
)

const (
	MENU_PATH = "/MenuBar"
	SNI_PATH  = "/StatusNotifierItem"
	// SNI_INTERFACE_NAME stands for StatusNotifierItem
	SNI_INTERFACE_NAME = "org.kde.StatusNotifierItem"
	// SNW_INTERFACE_NAME stands for StatusNotifierWatcher
	SNW_INTERFACE_NAME      = "org.kde.StatusNotifierWatcher"
	DBUSMENU_INTERFACE_NAME = "com.canonical.dbusmenu"

	// TextDirectionLTR stands for left to right
	TextDirectionLTR TextDirection = "ltr"
	// TextDirectionRTL stands for right to left
	TextDirectionRTL TextDirection = "rtl"

	// MenuStatusNormal would be right to use in almost all cases
	MenuStatusNormal MenuStatus = "normal"
	// MenuStatusNotice is used to indicate higher priority then normal
	MenuStatusNotice MenuStatus = "notice"
)

// NameBySpec returns name in format defined in the spec. For example:
// org.kdeStatusNotifierItem-<process id>-<instance number>
func NameBySpec(instance uint32) string {
	return fmt.Sprintf("org.kde.StatusNotifierItem-%d-%d", os.Getpid(), instance)
}

type Pixmap struct {
	Width  int32
	Heigth int32
	// Data is represented in ARGB32 format and is in the network byte order
	Data []byte
}

type ToolTip struct {
	First  string
	Second []struct {
		First  int32
		Second int32
		Third  []byte
	}
	Third  string
	Fourth string
}
