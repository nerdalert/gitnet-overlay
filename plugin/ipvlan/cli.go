package ipvlan

import "github.com/codegangsta/cli"

var (
	//  Exported user CLI flag config options
	FlagIPVlanMode     = cli.StringFlag{Name: "mode", Value: ipVlanMode, Usage: "name of the ipvlan mode [l2|l3]. (default: l2)"}
	FlagGateway        = cli.StringFlag{Name: "gateway", Value: "", Usage: "IP of the default gateway (defaultL2 mode: first usable address of a subnet. Subnet 192.168.1.0/24 would mean the container gateway to 192.168.1.1)"}
	FlagSubnet         = cli.StringFlag{Name: "ipvlan-subnet", Value: defaultSubnet, Usage: "subnet for the containers (l2 mode: 192.168.1.0/24)"}
	FlagMtu            = cli.IntFlag{Name: "mtu", Value: cliMTU, Usage: "MTU of the container interface (default: 1500)"}
	FlagIpvlanEthIface = cli.StringFlag{Name: "host-interface", Value: IpVlanEthIface, Usage: "(required) interface that the container will be communicating outside of the docker host with"}
	// Git bootstrap variables
	FlagGitPollingInt = cli.IntFlag{Name: "poll-interval", Value: 10, Usage: "polling interval in seconds (default: 15 seconds)"}
	FlagGitRepo       = cli.StringFlag{Name: "repo", Value: "", Usage: "(required) URL for the Git boostrap repo. Example: http://username:password@172.16.86.1/nerdalert/git-overlay.git"}
	FlagGitBootstrap  = cli.BoolFlag{Name: "git-multihost", Usage: "(default [false]) run an experimental multi-host endoint bootstrap."}
)

var (
	// These are the default values that are overwritten if flags are used at runtime
	ipVlanMode     = "l2"             // ipvlan l2 is the default
	IpVlanEthIface = ""               // default to eth0?
	defaultSubnet  = "192.168.1.0/24" // Should this just be the eth0 IP subnet?
	gatewayIP      = ""               // GW required for L2. increment network addr+1 if not defined
	cliMTU         = 1500

	BaseDir                = "data" // configurable base directory where records are cached
	gitBootstrapURL        = ""
	gitBoostrapBool        = false
	gitPollInterval        = 10 // time in seconds
	defaultGitPollInterval = 10 // time in seconds
//	gitRepo   = "http://172.16.86.1/nerdalert/git-overlay.git"
)
