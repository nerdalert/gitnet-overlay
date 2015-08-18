package control

const (
	ifname       = "en0"
	dotjson      = ".json"
	branch       = "master"
	configFormat = "json"
	chunk        = 64000

//	url             = "http://nerdalert:spam99@172.16.86.1/nerdalert/git-overlay.git"
)

// Exported default configurations values
var (
	BaseDirectory           = "data" // configurable base directory where records are cached
	GitDatastoreURL         = ""
	MasterEthIface          = ""
	DefaultInterval         = 20                                  // time in seconds
	DefaultIntervalMin      = 10                                  // time in seconds
	EndpointBase            = BaseDirectory + "/endpoints"        // .old and .new are appended to this
	EndpointPushRoot        = BaseDirectory + "/config"           // new endpoints push their conf from here
	EndpointPushSubDir      = BaseDirectory + "/config/endpoints" // new endpoint sources here
	DefaultBackupPath       = BaseDirectory + "/snapshots"        // backup path for configs
	EndpointStoreLatestRoot = EndpointBase + ".latest"            // latest pull from the git repository
	EndpointStoreOldRoot    = EndpointBase + ".old"               // latest pull from the git repository
	EndpointStoreLatest     = EndpointStoreLatestRoot + "/endpoints"
	EndpointStoreOld        = EndpointStoreOldRoot + "/endpoints"
	timef                   = "2006-01-02 15:04:05.99 -0700 MST" // Time format for logging
)
