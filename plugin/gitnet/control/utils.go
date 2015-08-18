package control

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

func CreatePaths(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil && !os.IsExist(err) {
		log.Warnf("Could not create net plugin path directory [ %s ]: %s", path, err)
		return err
	}
	return nil
}

func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

func getIfaceAddrStr(name string) (string, error) {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return "", err
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return "", err
	}
	var addrs4 []net.Addr
	for _, addr := range addrs {
		ip := (addr.(*net.IPNet)).IP
		if ip4 := ip.To4(); len(ip4) == net.IPv4len {
			addrs4 = append(addrs4, addr)
		}
	}
	switch {
	case len(addrs4) == 0:
		return "", fmt.Errorf("Interface [ %s ] has no IP addresses", name)
	case len(addrs4) > 1:
		log.Warnf("Interface %v has more than 1 IPv4 address. Defaulting to IP %v\n", name, (addrs4[0].(*net.IPNet)).IP)
	}
	s := strings.Split(addrs4[0].String(), "/")
	ip, _ := s[0], s[1]
	return ip, err
}

func hostLookup(ipstr string) (error, string) {
	// if the string is not an ip address, try to resolve an ip from hostname
	ipstr = strings.TrimSpace(ipstr)
	var err error
	if resolvedIP, err := dnsLookup(ipstr); err == nil {
		log.Debug("Endpoint hostname received, resolved host the name to the IP: [ %s ]:", resolvedIP)
		return nil, resolvedIP
	}
	return err, ""
}

func dnsLookup(s string) (string, error) {
	ipAddr, err := net.ResolveIPAddr("ip", s)
	return ipAddr.IP.String(), err
}

func getFilenames(dir string) ([]string, error) {
	file, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	res := make([]string, len(file))
	for i, fi := range file {
		res[i] = fi.Name()
	}
	return res, nil
}

func fileExists(s string, endpoints []string) bool {
	for _, endpoint := range endpoints {
		if endpoint == s {
			return true
		}
	}
	return false
}

// Print formatted JSON or YAML for easier debugging
func printPretty(data interface{}, format string) {
	var p []byte
	var err error
	switch format {
	case "json":
		p, err = json.MarshalIndent(data, "", "\t")
	case "yaml":
		p, err = yaml.Marshal(data)
	default:
		fmt.Printf("unsupported format: %s", format)
		return
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s", p)
}

func removeDups(endpoints []string) []string {
	m := map[string]struct{}{}
	for _, v := range endpoints {
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
		}
	}
	list := make([]string, len(m))

	i := 0
	for v := range m {
		list[i] = v
		i++
	}
	return list
}

// AddHostRoute adds a host-scoped route to a device.
func AddRoute(neighborNetwork *net.IPNet, nextHop net.IP, iface netlink.Link) error {
	log.Debugf("Adding route in the default namespace for IPVlan L3 mode with the following:")
	log.Debugf("IP Prefix: [ %s ] - Next Hop: [ %s ] - Source Interface: [ %s ]",
		neighborNetwork, nextHop, iface.Attrs().Name)

	return netlink.RouteAdd(&netlink.Route{
		Scope:     netlink.SCOPE_UNIVERSE,
		LinkIndex: iface.Attrs().Index,
		Dst:       neighborNetwork,
		Gw:        nextHop,
	})
}
