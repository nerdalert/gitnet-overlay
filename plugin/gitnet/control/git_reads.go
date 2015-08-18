package control

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

// Run the daemon
func (g *ControlConfig) Read() {
	// set a timer
	timer := time.NewTicker(time.Duration(g.GetTimeInterval()) * time.Second)
	g.Process()

	for {
		// fork a go routine to watch for config updates
		doneChan := make(chan bool)
		go func(doneChan chan bool) {
			// node state ignored for now
			minus, plus, err := DiffDirectories(EndpointStoreOld, EndpointStoreLatest)
			log.Infof("Records removed [ %s ]", minus)
			log.Infof("Records added [ %s ]", plus)
			if err != nil {
				log.Errorf("Error diffing the latest records")
			}
			log.Debugf("Cleaning old endpoint store dir [ %s ] to ensure consistency", EndpointStoreOldRoot)
			// reset the oldcache for diff against the latest update
			err = os.RemoveAll(EndpointStoreOldRoot)

			if err != nil {
				log.Errorf("Error deleting old endpoints dir [ %s ]: %s", EndpointStoreOldRoot, err)
			} else {
				log.Debugf("Succesfully removed old endpoints dir [ %s ]", EndpointStoreOldRoot)
			}

			log.Debugf("copying [ %s ] endpoint cache to [ %s ] endpoint cache dir for diffing", EndpointStoreLatestRoot, EndpointStoreOldRoot)
			err = copyDir(EndpointStoreLatestRoot, EndpointStoreOldRoot)
			if err != nil {
				log.Errorf("Error copying latest endpoints [ %s ] to old endpoints dir [ %s ]", EndpointStoreLatestRoot, EndpointStoreOldRoot)
			}
		}(doneChan)
		// wait for timer to expire or insert a channel event
		for b := true; b; {
			select {
			case <-timer.C:
				log.Debug("Interval time expired")
				// process git operations
				g.Process()
				b = false
			}
		}
		log.Debugf("checking repo for updates %s ", time.Now().UTC())
	}
}

func (g *ControlConfig) Process() {
	git := newGit(EndpointStoreLatestRoot)
	gitExists := fmt.Sprintf("%s/.git", EndpointStoreLatestRoot)

	if _, err := os.Stat(gitExists); err != nil {
		c := git.clone(g.GetRepo())
		err := c.Run()
		if err != nil {
			log.Debugf("Error cloning ensure the git server is reachable [ %s ] Git returned: %s", g.GetRepo(), err)
		}
	}

	cmdOutput := &bytes.Buffer{}

	gitPull := []string{"-C", EndpointStoreLatestRoot, "pull"}
	log.Debugf("Running: git -C %s pull", EndpointStoreLatestRoot)
	err := gitCmd(gitPull).Run()
	if err != nil {
		log.Errorf("Error updating the data store via git: %s", err)
	}
	time.Sleep(1 * time.Second)

	debugGitCmd(cmdOutput.Bytes())
	if !strings.Contains(string(cmdOutput.Bytes()), "Already up-to-date") {
		err := g.backupConf()
		if err != nil {
			log.Fatalf("Error backing up the endpoint state: %s", err)
		}
	}
}
