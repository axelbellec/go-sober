package platform

// InitPlatform
// This init will call ours bootstraps functions in platform module
// We do it to avoid circular dependencies and Go default alphabetic order initialisation
func InitPlatform() {
	initConfig()
	initLogger()
}
