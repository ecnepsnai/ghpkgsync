package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
)

func syncRPMRepo() {
	cmd := exec.Command("/usr/bin/createrepo_c", ".")
	cmd.Dir = "repo/rpm"
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.PError("Error syncing rpm repository", map[string]interface{}{
			"error":  err.Error(),
			"output": string(output),
		})
		return
	}
	log.Info("Synced rpm repository")
	runtime.GC()
}

func makeRPMRepoFile() {
	repoFileDataTemplate := `[%s]
name=%s
baseurl=%s
enabled=1
gpgcheck=0`
	repoFileData := fmt.Sprintf(repoFileDataTemplate, Config.YumRepoID, Config.YumRepoDescription, Config.YumRepoBaseurl)

	err := os.WriteFile(path.Join("repo", "rpm", Config.YumRepoID+".repo"), []byte(repoFileData), os.ModePerm)
	if err != nil {
		log.PError("Error writing rpm repo file", map[string]interface{}{
			"error": err.Error(),
		})
	}
}
