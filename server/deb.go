package main

import "os/exec"

func syncDEBRepo() {
	cmd := exec.Command("/usr/bin/dpkg-scanpackages", ".")
	cmd.Dir = "repo/deb"
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.PError("Error syncing deb repository", map[string]interface{}{
			"error":  err.Error(),
			"output": string(output),
		})
		return
	}
	log.Info("Synced deb repository")
}
