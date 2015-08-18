package main

// Example usage:
// go run *.go \
// -t 30 \
// -l debug \
// -r http://localhost:3000/nerdalert/git-overlay.git nerdalert/git-overlay
// # or 1-line
// go run *.go  -t 10 -l debug -r http://localhost:3000/nerdalert/git-overlay.git nerdalert/git-overlay

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jessevdk/go-flags"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/nerdalert/gitnet-overlay/plugin/gitnet/control"
)

// only daemon mode supported atm. The -d flag is ignored
var opts struct {
	GitRepoFlag       string `short:"r" long:"repo" description:"(required) target repository url - example format: https://github.com/nerdalert/plugin-watch.git"`
	TimeIntervalFlag  int    `short:"t" long:"time" description:"(requiredl) time in seconds between Git repository update checks."`
	BaseDirectoryFlag string `short:"b" long:"backup-path" description:"(default: [ data/ ]) path to the base directory where the git repo and json config resides."`
	Daemon            bool   `short:"d" long:"daemon" description:"(optional:default [true]) run as a daemon. Alternatively could be run via a cron job."`
	LogLevelFlag      string `short:"l" long:"loglevel" description:"(optional:default [info]) set the logging level. Options are [debug, info, warn, error]."`
	Help              bool   `short:"h" long:"help" description:"show app help."`
}

func init() {
	runtime.GOMAXPROCS(1)
	ch := make(chan os.Signal, 1)
	go sigHandler(ch)
}

func sigHandler(ch chan os.Signal) {
	signal.Notify(ch, os.Interrupt, syscall.SIGKILL, syscall.SIGTERM)
	go func() {
		for _ = range ch {
			os.Exit(0)
		}
	}()
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}
	if opts.Help {
		showUsage()
		os.Exit(1)
	}
	// Null check required fields
	if opts.GitRepoFlag == "" {
		showUsage()
		log.Fatal("Required repo name is missing")
		os.Exit(1)
	} else {
		control.GitDatastoreURL = opts.GitRepoFlag // Bind to a global var
	}
	if opts.TimeIntervalFlag < control.DefaultIntervalMin {
		showUsage()
		log.Fatal("The minimum polling interval is 10 seconds.")
		os.Exit(1)
	}
	// Bind opts to a couple global vars for convenience
	if opts.BaseDirectoryFlag != "" {
		control.BaseDirectory = opts.BaseDirectoryFlag
	}

	//	var timeInterval int
	//	timeInterval = opts.TimeIntervalFlag
	//	if opts.TimeIntervalFlag == 0 {
	//		timeInterval = control.DefaultInterval
	//		log.Debug("Polling interval not specified, setting it to 90 seconds")
	//	}
	// Set logrus logging level, default is Info
	switch opts.LogLevelFlag {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
		log.Debug("Logging level is set to : ", log.GetLevel())
	}

	//	control.Run(opts.GitRepoFlag, timeInterval, opts.Daemon)

	//	if opts.Daemon == true {
	//		control.RunGit() // run as a daemon for every (n) seconds
	//	} else {
	//		control.RunGit() // TODO: if false, run one time and exit
	//	}
	//	g := control.gitNet(opts.GitRepoFlag, timeInterval)

}

func showUsage() {
	var usage string
	usage = `
Usage:
  main

Application Options:
    -r, --repo=         (required) target repository url - example format: https://github.com/nerdalert/plugin-watch.git
    -t, --time=         (requiredl) time in seconds between Git repository update checks.
    -c, --config-path=  (recommended: default: [ ./tmp/conf/ ]) path to config files.
    -b, --backup-path=  (recommended: default: [ data/snapshots ]) path to the backup endpoint config files.
    -s, --server=       (optional: default: [ data/endpoints ]) path to config files.
    -d, --daemon=       (optional:default [ true ]) run as a daemon. Alternatively could be run via a cron job.
    -l, --loglevel=     (optional:default [ info ]) set the logging level. Default is 'info'. options are [debug, info, warn, error].
    -h, --help    show app help.

Example daemon mode processing flows every 2 minutes:
	git-control -r github.com/plugin-watch -t 120 -l debug -r https://github.com/nerdalert/plugin-watch.git

Example run-once export:
    TODO:

Help Options:
  -h, --help    Show this help message
  `
	log.Print(usage)
}

//func initSigs() {
//	c := make(chan os.Signal, 1)
//	signal.Notify(c, os.Interrupt, os.Kill)
//	go func() {
//		s := <-c
//		log.Warnln("Got signal: ", s)
//		cleanUp()
//		os.Exit(1)
//	}()
//}
//
//func sigHandler(ch chan os.Signal) {
//	signal.Notify(ch, os.Interrupt, syscall.SIGKILL, syscall.SIGTERM)
//	go func() {
//		for _ = range ch {
//			os.Exit(0)
//		}
//	}()
//}
//
//func cleanUp() {
//	log.Infoln("Exiting process..")
//}
