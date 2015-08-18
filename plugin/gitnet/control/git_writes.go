package control

import (
	"fmt"
	"net"
	"os"
	"time"

	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type LocalEndpoint struct {
	Endpoint string `json:"Endpoint"`
	Network  string `json:"Network"`
	Meta     string `json:"Meta"`
}

// AnnounceEndpoint takes libnetwork endpoint data writes it
// to a json file and pushes it to a Git repo for distribution
func AnnounceEndpoint(networkCidr *net.IPNet, hostIface, gitRepoURL string) {
	log.Debugf("Cloning [ %s ] into [ %s ]", GitDatastoreURL, EndpointPushRoot)
	if gitRepoURL == "" {
		log.Fatalf("No Git repo was specified")
	}
	initConfigWrite(networkCidr, hostIface, gitRepoURL)
}

// Initialize the repo to be used to announce/write config.
// A seperate repo is initialized to read incoming announcements
func initConfigWrite(networkCidr *net.IPNet, hostIface, gitRepoURL string) {
	var err error
	if !pathExists(EndpointPushSubDir) {
		log.Debugf("[ %s ] dir not found, creating it..", EndpointPushSubDir)
		if err = CreatePaths(EndpointPushSubDir); err != nil {
			log.Fatalf("Could not create the directory [ %s ]: %s", EndpointPushSubDir, err)
		} else {
			log.Warnf("Succesfully created the config path [ %s ]", EndpointPushSubDir)
		}
	}
	// Create the cache subdirectories
	time.Sleep(1 * time.Second)
	localEndpointIP, _ := getIfaceAddrStr(hostIface)
	// Fun Go fact: using a + with sprintf is faster then %s since it uses reflection
	endpointFile := fmt.Sprintf(localEndpointIP + dotjson)
	log.Debugf("The endpoint file name is [ %s ] ", endpointFile)
	log.Debugf("Anouncing this endpoint using the source [ %s ] and advertsing network [ %s ] to datastore file [ %s ]", networkCidr, localEndpointIP, endpointFile)
	endpointConfig := &LocalEndpoint{
		Endpoint: localEndpointIP,
		Network:  networkCidr.String(),
		Meta:     "",
	}
	var configAnnounce []LocalEndpoint
	configAnnounce = append(configAnnounce, *endpointConfig)
	marshallConfig(configAnnounce, configFormat, endpointFile)
	if log.GetLevel().String() == "debug" {
		printPretty(configAnnounce, "json")
	}
	// Parse the repo name
	defer gitPushConfig()
}

func gitPushConfig() {

	log.Debugf("Committing and pushing the new endpoint from [ %s ] ", EndpointPushRoot)
	// Get the ip addr of the endpoint address
	ipaddress, _ := getIfaceAddrStr("en0")
	// Add the IP to the commit msg
	commitMsg := fmt.Sprintf("added endoint %s", ipaddress)
	// name of the endpointconfig file e.g. <ip>.json TODO: uncomment below
	// endpointFile, _ := getEndpointName(ifname)
	// relativeEndpointDir := fmt.Sprintf(EndpointPushSubDir+"/%s", endpointFile)
	// TODO: specify exact endpoint to update rather then the entire dir
	// that needs to reconcile with a first time run when endpoints dir !exist
	// endpoints directory in the git repo
	relativeEndpointDir := fmt.Sprintf("endpoints/")
	// Git commands to push local endpoint
	gitAdd := []string{"-C", EndpointPushRoot, "add", relativeEndpointDir}
	gitCommit := []string{"-C", EndpointPushRoot, "commit", "-m", commitMsg}
	// gitPush := []string{"-C", EndpointPushRepo, "push", "-u", GitDatastoreURL, branch}
	gitPush := []string{"-C", EndpointPushRoot, "push", "-f"}
	// Commit and push the endpoint config
	log.Debugf("Running: git %s", gitAdd)
	err := gitCmd(gitAdd).Run()
	if err != nil {
		log.Errorf("Error in git operation: %s", err)
	}
	time.Sleep(1 * time.Second)
	log.Debugf("Running: git %s", gitCommit)
	err = gitCmd(gitCommit).Run()
	if err != nil {
		log.Errorf("Error in git operation: %s", err)
	}
	time.Sleep(1 * time.Second)
	log.Debugf("Running: git %s", gitPush)
	err = gitCmd(gitPush).Run()
	if err != nil {
		log.Errorf("Error in git operation: %s", err)
	}
}

// get the endpoint json filname
func getEndpointName(ifname string) (string, error) {
	ipaddress, err := getIfaceAddrStr(ifname)
	if err != nil {
		return "", err
	}
	endpointFile := fmt.Sprintf("%s%s", ipaddress, dotjson)
	return endpointFile, nil
}

// MarshallConfig marshall the network cidr, endpoint ip, etc into a json file to anounce in Git
func marshallConfig(data interface{}, format, endpointFile string) error {
	var marshall []byte
	var err error
	switch format {
	case "json":
		marshall, err = json.MarshalIndent(data, "", "  ")
	case "yaml":
		marshall, err = yaml.Marshal(data) // todo yaml not yet supported
	default:
		log.Fatalf("Unsupported encoding format: %s", format)
		return err
	}
	if err != nil {
		log.Warn(err)
		return err
	}
	log.Debugf("The local config record written to [ %s ] and anounced to it's peers via Git is:\n %s", endpointFile, marshall)
	return writeEndpointRecord(endpointFile, marshall)
}

// writeEndpointRecord writes the json local json config to file
func writeEndpointRecord(endpointFile string, marshall []byte) error {
	relativeConfigFile := fmt.Sprintf(EndpointPushSubDir + "/" + endpointFile)
	// Remove existing config if it exists
	os.Remove(relativeConfigFile)
	file, err := os.Create(relativeConfigFile)
	if err != nil {
		log.Debugf("Error creating local endpoint record named [ %s ] to anounce via git: %s", endpointFile, err)
		return err
	}
	file.WriteString(string(marshall))
	if err != nil {
		log.Debugf("Error writing endpoint config to [ %s ]: %s", relativeConfigFile, err)
		return err
	}
	return nil
}
