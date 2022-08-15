package tray

import (
	"image"

	"github.com/godbus/dbus/v5"
	"github.com/knightpp/sni/generated/status_notifier_item"
	"github.com/knightpp/sni/pkg/sni"
)

// SetId sets an id that should be unique for this application and consistent
// between sessions, such as the application name itself.
func (t *Tray) SetId(id string) *Tray {
	t.propsSni["Id"].Value = id
	return t
}

// SetTitle sets a name that describes the application, it can be more
// descriptive than Id.
func (t *Tray) SetTitle(title string) *Tray {
	t.propsSni["Title"].Value = title
	return t
}

// SetIconName sets StatusNotifierItem IconName property
func (t *Tray) SetIconName(name string) *Tray {
	t.propsSni["IconName"].Value = name
	return t
}

// SetIconPixmap sets StatusNotifierItem IconPixmap property.
// Copies src to a new image and changes pixel format to ARGB32.
//
// Note: see SetIconPixmapRaw
func (t *Tray) SetIconPixmap(src image.Image) *Tray {
	t.propsSni["IconPixmap"].Value = []Pixmap{
		imageToArgb32(src),
	}
	return t
}

// SetIconPixmapRaw sets StatusNotifierItem IconPixmap property.
//
// Note: see SetIconPixmap for higher level abstraction
func (t *Tray) SetIconPixmapRaw(pixmaps []Pixmap) *Tray {
	t.propsSni["IconPixmap"].Value = pixmaps
	return t
}

// GetIconName returns StatusNotifierItem IconName property
func (t *Tray) GetIconName() string {
	v, ok := t.propsSni["IconName"]
	if !ok {
		panic("GetIconName(): no such value in a map")
	}
	s, ok := v.Value.(string)
	if !ok {
		panic("GetIconName(): value is not string")
	}
	return s
}

// SetWindowId sets WindowId property
func (t *Tray) SetWindowId(id int32) *Tray {
	t.propsSni["WindowId"].Value = id
	return t
}

// SetItemIsMenu sets ItemIsMenu property
func (t *Tray) SetItemIsMenu(b bool) *Tray {
	t.propsSni["ItemIsMenu"].Value = b
	return t
}

// SetOverlayIconName is a property of StatusNotifierItem
func (t *Tray) SetOverlayIconName(name string) *Tray {
	t.propsSni["OverlayIconName"].Value = name
	return t
}

// SetOverlayIconPixmap sets StatusNotifierItem OverlayIconPixmap property.
// Copies src to a new image and changes pixel format to ARGB32.
//
// Note: see SetOverlayIconPixmapRaw
func (t *Tray) SetOverlayIconPixmap(src image.Image) *Tray {
	t.propsSni["OverlayIconPixmap"].Value = []Pixmap{
		imageToArgb32(src),
	}
	return t
}

// SetAttentionIconName sets StatusNotifierItem AttentionIconName prop.
func (t *Tray) SetOverlayIconPixmapRaw(pixmaps []Pixmap) *Tray {
	t.propsSni["OverlayIconPixmap"].Value = pixmaps
	return t
}

// SetAttentionIconName sets StatusNotifierItem AttentionIconName prop.
func (t *Tray) SetAttentionIconName(name string) *Tray {
	t.propsSni["AttentionIconName"].Value = name
	return t
}

// SetAttentionIconPixmap sets StatusNotifierItem AttentionIconPixmap property.
// Copies src to a new image and changes pixel format to ARGB32.
//
// Note: see SetOverlayIconPixmapRaw
func (t *Tray) SetAttentionIconPixmap(src image.Image) *Tray {
	t.propsSni["AttentionIconPixmap"].Value = []Pixmap{
		imageToArgb32(src),
	}
	return t
}

// SetAttentionMovieName sets StatusNotifierItem AttentionMovieName prop.
func (t *Tray) SetAttentionMovieName(name string) *Tray {
	t.propsSni["AttentionMovieName"].Value = name
	return t
}

// SetToolTipRaw sets StatusNotifierItem ToolTip prop.
func (t *Tray) SetToolTipRaw(tooltip ToolTip) *Tray {
	t.propsSni["ToolTip"].Value = tooltip
	return t
}

// SetSniStatus sets StatusNotifierItem Status property.
//
// It describes the status of this item or of the associated application.
//
// The allowed values for the Status property are:
//
// Passive: The item doesn't convey important information to the user, it
// can be considered an "idle" status and is likely that visualizations
// will chose to hide it.
//
// Active: The item is active, is more important that the item will be shown
// in some way to the user.
//
// NeedsAttention: The item carries really important information for the user,
// such as battery charge running out and is wants to incentive the direct user
// intervention. Visualizations should emphasize in some way the items with
// NeedsAttention status.
func (t *Tray) SetSniStatus(status sni.Status) *Tray {
	t.propsSni["Status"].Value = status
	return t
}

// Represents the way the text direction of the application. This
// allows the server to handle mismatches intelligently. For left-
// to-right the string is "ltr" for right-to-left it is "rtl".
func (t *Tray) SetMenuTextDirection(dir TextDirection) *Tray {
	t.propsMenu["TextDirection"].Value = dir
	return t
}

// Tells if the menus are in a normal state or they believe that they
// could use some attention. Cases for showing them would be if help
// were referring to them or they accessors were being highlighted.
// This property can have two values:
//
// - "normal" in almost all cases and
//
// - "notice" when they should have a higher priority to be shown.
func (t *Tray) SetMenuStatus(status MenuStatus) *Tray {
	t.propsMenu["Status"].Value = status
	return t
}

// A list of directories that should be used for finding icons using
// the icon naming spec. Ideally there should only be one for the icon
// theme, but additional ones are often added by applications for
// app specific icons.
func (t *Tray) SetMenuIconThemePath(path []string) *Tray {
	t.propsMenu["IconThemePath"].Value = path
	return t
}

// SetCategory sets a category of the StatusNotifierItem.
// Default value is ApplicationStatus.
func (t *Tray) SetCategory(cat sni.Category) *Tray {
	t.propsSni["Category"].Value = cat
	return t
}

// SetMenuPath sets a path to a menu.
//
// You shouldn't use this function if you don't know what it is.
func (t *Tray) SetMenuPath(path dbus.ObjectPath) *Tray {
	t.propsSni["Menu"].Value = path
	return t
}

/*-----------------------SIGNALS------------------------*/

// SignalNewIcon emits signal on dbus thus requesting re-rendering of its icon.
// You should emit this signal to reflect change of the icon visually.
func (t *Tray) SignalNewIcon() error {
	err := status_notifier_item.Emit(t.conn, &status_notifier_item.StatusNotifierItem_NewIconSignal{
		Path: SNI_PATH,
		Body: &status_notifier_item.StatusNotifierItem_NewIconSignalBody{},
	})
	return err
}

// SignalNewTitle emits signal on dbus notifying system that title was changed
func (t *Tray) SignalNewTitle() error {
	return status_notifier_item.Emit(t.conn,
		&status_notifier_item.StatusNotifierItem_NewTitleSignal{
			Path: SNI_PATH,
			Body: &status_notifier_item.StatusNotifierItem_NewTitleSignalBody{},
		})
}

// SignalNewAttentionIcon emits signal on dbus notifying system that
// attention icon was changed.
func (t *Tray) SignalNewAttentionIcon() error {
	return status_notifier_item.Emit(t.conn,
		&status_notifier_item.StatusNotifierItem_NewAttentionIconSignal{
			Path: SNI_PATH,
			Body: &status_notifier_item.StatusNotifierItem_NewAttentionIconSignalBody{},
		})
}

// SignalNewOverlayIcon emits signal on dbus notifying system that
// overlay icon was changed.
func (t *Tray) SignalNewOverlayIcon() error {
	return status_notifier_item.Emit(t.conn,
		&status_notifier_item.StatusNotifierItem_NewOverlayIconSignal{
			Path: SNI_PATH,
			Body: &status_notifier_item.StatusNotifierItem_NewOverlayIconSignalBody{},
		})
}

// SignalNewToolTip emits signal on dbus notifying system that
// tooltip was changed.
func (t *Tray) SignalNewToolTip() error {
	return status_notifier_item.Emit(t.conn,
		&status_notifier_item.StatusNotifierItem_NewToolTipSignal{
			Path: SNI_PATH,
			Body: &status_notifier_item.StatusNotifierItem_NewToolTipSignalBody{},
		})
}

// SignalNewStatus emits signal on dbus notifying system that
// status was changed.
func (t *Tray) SignalNewStatus() error {
	return status_notifier_item.Emit(t.conn,
		&status_notifier_item.StatusNotifierItem_NewStatusSignal{
			Path: SNI_PATH,
			Body: &status_notifier_item.StatusNotifierItem_NewStatusSignalBody{},
		})
}
