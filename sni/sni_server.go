package sni

import (
	"log"

	"github.com/knightpp/sni/generated/status_notifier_item"

	"github.com/godbus/dbus/v5"
)

type SniServer struct {
	*status_notifier_item.StatusNotifierItem
}

func NewSniServer() *SniServer {
	return &SniServer{}
}

// ContextMenu is org.kde.StatusNotifierItem.ContextMenu method.
func (s *SniServer) ContextMenu(x, y int32) (err *dbus.Error) {
	log.Printf("ContextMenu(x = %d, y = %d)", x, y)
	return nil
}

// Activate is org.kde.StatusNotifierItem.Activate method.
func (s *SniServer) Activate(x, y int32) (err *dbus.Error) {
	log.Printf("Activate(x = %d, y = %d)", x, y)
	return nil
}

// SecondaryActivate is org.kde.StatusNotifierItem.SecondaryActivate method.
func (s *SniServer) SecondaryActivate(x, y int32) (err *dbus.Error) {
	log.Printf("SecondaryActivate(x = %d, y = %d)", x, y)
	return nil
}

// Scroll is org.kde.StatusNotifierItem.Scroll method.
func (s *SniServer) Scroll(delta int32, orientation string) (err *dbus.Error) {
	log.Printf("Scroll(delta = %d, orientation = %s)", delta, orientation)
	return nil
}
