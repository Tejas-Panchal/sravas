package job

import "log"

// ProcessUploads scans for videos stuck in "uploading" status and retries or marks them failed
func ProcessUploads() {
	log.Println("[scheduler] process uploads: scanning for stuck uploads")
	// TODO: query DB for videos with status=uploading older than 30 min
	// TODO: for each, call processing-service to retry, or set status=failed
}
