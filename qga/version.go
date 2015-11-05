package qga

var (
	// Version string (git descrive --long)
	Version string
	// BuildTime build time
	BuildTime string
)

// GetVersion display current qemu-ga version
func GetVersion() string {
	return Version
}
