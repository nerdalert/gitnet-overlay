package control

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os/exec"
	"time"
)

const (
	timeFmt = "20060102150405"
)

func (g *ControlConfig) backupConf() error {

	//	tar -cvzf tarballname.tar.gz itemtocompress
	tarBackup := fmt.Sprintf(DefaultBackupPath + "/" + "endpoint-snapshot-" + getTime() + ".tar.gz")
	cmd := exec.Command("tar", "cvzf", tarBackup, EndpointPushRoot)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Error creating config backup for [%s]", DefaultBackupPath, out, err)
	}
	if err != nil {
		fmt.Printf("Error creating config backup: [%s]\n", err)
		return err
	}
	return nil
}

func getTime() string {
	t := time.Now().UTC().Local()
	fmt.Println(t.Format(timeFmt))
	return t.Format(timeFmt)
}
