package control

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"os"
)

// Useful or not, on the fence
type ControlConfig struct {
	Repo         string // URL for git repo.
	TimeInterval int    // time between checking repo for endpoint updates
	BaseDir      string // base cache directory
	MasterIface  string // master netlink iface "ethX"
}

type toString interface {
	String() string
}

// GetParam accessors
func (g *ControlConfig) GetRepo() string {
	return g.Repo
}
func (g *ControlConfig) GetTimeInterval() int {
	return g.TimeInterval
}
func (g *ControlConfig) GetBasePath() string {
	return g.BaseDir
}
func (g *ControlConfig) String() string {
	s := fmt.Sprintf("Repository: [%s] \n"+
		"Polling interval: [%d] \n"+
		"Base directory path: [%s] \n",
		g.GetRepo(), g.GetTimeInterval(), g.GetBasePath())
	return s
}

func Run(gitURL, masterIface string, timeInterval int, daemon bool) {
	g := gitNet(gitURL, masterIface, timeInterval)
	log.Debugf("connecting to [ %s ] with the following paramters: \n%s ", g.GetRepo(), g)
	MasterEthIface = masterIface
	// Locally generated configs get pushed from this directory
	git := newGit(EndpointPushRoot)
	gitExists := fmt.Sprintf("%s/.git", EndpointPushRoot)
	if _, err := os.Stat(gitExists); err != nil {
		c := git.clone(g.GetRepo())
		err := c.Run()
		if err != nil {
			log.Debugf("Error cloning ensure the git server is reachable [ %s ] Git returned: %s", g.GetRepo(), err)
		}
	}
	if daemon == true {
		g.Read() // run as a daemon for every (n) seconds
	} else {
		g.Read() // TODO: if false, run one time and exit
	}

}

func gitNet(gitURL, masterIface string, timeInterval int) *ControlConfig {
	err := initCachePath()
	if err != nil {
		log.Warnf("Encountered an error while creating cache directories: %s", err)
	}

	return &ControlConfig{
		Repo:         gitURL,
		TimeInterval: timeInterval,
		BaseDir:      BaseDirectory,
		MasterIface:  masterIface,
	}
}

func initCachePath() error {
	log.Debugf("Initializing endpoint cache...")
	var err error
	if pathExists(BaseDirectory) {
		log.Debugf("Existing cache dir found, attempting to remove")
		os.RemoveAll(BaseDirectory)
		if err != nil {
			log.Warnf("Could not delete the old Git cache path [ %s ]: %s", BaseDirectory, err)
		} else {
			log.Warnf("Succesfully removed the old cache dir [ %s ]", BaseDirectory)
		}
	}
	// Create the cache subdirectories
	time.Sleep(1 * time.Second)
	log.Warnf("Creating the directory [ %s ]", EndpointStoreOldRoot)
	if err = CreatePaths(EndpointStoreOldRoot); err != nil {
		log.Fatalf("Could not create the directory [ %s ]: %s", EndpointStoreOldRoot, err)
		return err
	}
	if err = CreatePaths(EndpointStoreLatestRoot); err != nil {
		log.Fatalf("Could not create the directory [ %s ]: %s", EndpointStoreLatestRoot, err)
		return err
	}
	if err = CreatePaths(EndpointPushRoot); err != nil {
		log.Fatalf("Could not create the directory [ %s ]: %s", EndpointPushRoot, err)
		return err
	}
	if err = CreatePaths(DefaultBackupPath); err != nil {
		log.Warnf("Could not create the directory [ %s ]: %s", DefaultBackupPath, err)
		return err
	}
	time.Sleep(1 * time.Second)
	return nil
}
