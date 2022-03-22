package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

func syncRepo() {
	cmd := exec.Command("/usr/bin/createrepo_c", ".")
	cmd.Dir = "repo"
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.PError("Error syncing repository", map[string]interface{}{
			"error":  err.Error(),
			"output": string(output),
		})
		return
	}
	log.Info("Synced repository")
}

func makeRepoFile() {
	repoFileDataTemplate := `[%s]
name=%s
baseurl=%s
enabled=1
gpgcheck=0`
	repoFileData := fmt.Sprintf(repoFileDataTemplate, Config.YumRepoID, Config.YumRepoDescription, Config.YumRepoBaseurl)

	err := os.WriteFile(path.Join("repo", Config.YumRepoID+".repo"), []byte(repoFileData), os.ModePerm)
	if err != nil {
		log.PError("Error writing repo file", map[string]interface{}{
			"error": err.Error(),
		})
	}
}
