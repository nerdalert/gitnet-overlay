package control

import (
	"encoding/json"
	"io/ioutil"
	"net"

	log "github.com/Sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

type NodeEvent int

const (
	_ NodeEvent = iota
	EVENT_PREFIX_BEST
	EVENT_PREFIX_WITHDRAWN
	EVENT_NODE_JOIN
	EVENT_NODE_REMOVE
	EVENT_NODE_MODIFIED
)

// todo: unused implement
type Node struct {
	Endpoint string
	Network  string
	Meta     string
	SegID    int
}

// todo: unused implement
func nodeEvant(event NodeEvent, record []string) {
	switch event {
	case EVENT_PREFIX_BEST:
		log.Debugf("Event: prefix withdrawn not yet implemented")
	case EVENT_PREFIX_WITHDRAWN:
		log.Debugf("Event: prefix withdrawn not yet implemented")
	case EVENT_NODE_JOIN:
		log.Debugf("Event: node join occoured")
		nodeAdded(record)
	case EVENT_NODE_REMOVE:
		log.Debugf("Event: node added occoured")
		nodeRemoved(record)
	case EVENT_NODE_MODIFIED:
		log.Debugf("Event: node modified occoured")
	default:
		log.Errorf("unknown event  [ %v ]", event)
	}
}

// New remote endpoint added to the git datastore
func nodeAdded(records []string) {
	// Copy for the first diff at init time
	err := copyDir(EndpointStoreLatestRoot, EndpointStoreOldRoot)
	if err != nil {
		log.Error(err)
	}
	for _, record := range records {
		fdata, err := ioutil.ReadFile(record)

		if err != nil {
			log.Error(err)
		}
		log.Infof("Endpoint [ %s ] was added", stripSingleFilepath(record))

		var state []Node
		err = json.Unmarshal(fdata, &state)
		if err != nil {
			log.Errorf("Error unmarshalling json data, verify the record format:", err)
		}
		for _, s := range state {
			log.Debugf("Umarshalled New Endpoint: [ %s ]", s.Endpoint)
			log.Debugf("Umarshalled Network [ %s ]", s.Network)
			log.Debugf("Umarshalled Meta [ %s ]", s.Meta) // Meta will be empty since there arent any values in it. Add some.
			log.Debugf("Umarshalled SegID [ %d ]", s.SegID)
			localEndpointIP, _ := getIfaceAddrStr(MasterEthIface)
			ipvlanParent, err := netlink.LinkByName(MasterEthIface)
			if err != nil {
				log.Debugf("a problem occurred adding the container subnet default namespace route", err)
			}
			_, cidr, err := net.ParseCIDR(s.Network)
			nextHop := net.ParseIP(s.Endpoint)
			if err != nil {
				log.Debugf("a problem occurred adding the container subnet default namespace route", err)
			}
			if s.Endpoint != localEndpointIP {
				log.Warnf("Remote endpoint [ %s ] != Local endpoint [ %s ] adding route to remote endpoint [ %s ] using interface [ %s ]", s.Endpoint, localEndpointIP, s.Endpoint, ipvlanParent.Attrs().Name)
				if err = AddRoute(cidr, nextHop, ipvlanParent); err != nil {
					log.Errorf("An error occoured adding a netlink route: %s", err)
				}
			}
		}
	}
}

// node removed from the git datastore
func nodeRemoved(records []string) {
	log.Warnf("Event: node removed event not implemented yet")
	for _, record := range records {
		log.Infof("Endpoint [ %s ] was removed", stripSingleFilepath(record))
	}
}

// node file modified in the git datastore
func nodeModified(records string) {
	log.Warnf("Event: node modified [ %s ] event not implemented yet", records)
}
