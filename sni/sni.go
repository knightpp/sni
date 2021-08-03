package sni

type Status string
type Category string

const (
	// StatusPassive doesn't convey important information to the user, it can
	// be considered an "idle" status and is likely that visualizations will
	// chose to hide it.
	StatusPassive Status = "Passive"
	// StatusActive tells that item is active, is more important that the
	// item will be shown in some way to the user.
	StatusActive Status = "Active"
	// StatusNeedsAttention carries really important information for the user,
	// such as battery charge running out and is wants to incentive the direct
	// user intervention. Visualizations should emphasize in some way the items
	// with NeedsAttention status.
	StatusNeedsAttention Status = "NeedsAttention"

	// CategoryApplicationStatus describes the status of a generic application,
	// for instance the current state of a media player. In the case where
	// the category of the item can not be known, such as when the item is
	// being proxied from another incompatible or emulated system,
	// ApplicationStatus can be used a sensible default fallback.
	CategoryApplicationStatus Category = "ApplicationStatus"
	// CategoryCommunications describes the status of communication oriented
	// applications, like an instant messenger or an email client.
	CategoryCommunications Category = "Communications"
	// CategorySystemServices describes services of the system not seen as
	// a stand alone application by the user, such as an indicator for the
	// activity of a disk indexing service.
	CategorySystemServices Category = "SystemServices"
	// CategoryHardware describes the state and control of a particular
	// hardware, such as an indicator of the battery charge or sound card
	// volume control.
	CategoryHardware Category = "Hardware"
)
