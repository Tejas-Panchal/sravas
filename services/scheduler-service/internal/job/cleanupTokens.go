package job

import "log"

// CleanupTokens removes expired refresh tokens from the store
func CleanupTokens() {
	log.Println("[scheduler] cleanup tokens: removing expired refresh tokens")
	// TODO: query DB for refresh tokens where expires_at < now
	// TODO: delete expired rows
}
