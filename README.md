ipvlan-docker-plugin
=================

IPVlan is a lightweight L2 and L3 network implementation that does not require traditional bridges. The purpose of this is to have references for plugging into the docker networking APIs which are now available as part of libnetwork. Libnetwork is still under development and is considered as experimental at this point.

This example uses Git to bootstrap multi-host environment. I like the idea of using git as a datastore since it offers a good amount of operational oversite and collaboration along with a record of exactly what changed in the infra. If infrastructure as code is the paradigm, then infra code being sourced from Git makes sense to me.

The Docker hosts need to share the same network segment at the moment as I havent tested beyond that. I will add VXLAN tunnels to abstract the underlay and next hop routes for IPVlan L3 mode is all that is required for connectivity. The same concept with Quagga or a BGP daemon could be used with IPVlan for underlay integration as desired.

### Pre-Requisites

1. Install the Docker experimental binary from the instructions at: [Docker Experimental](https://github.com/docker/docker/tree/master/experimental). (stop other docker instances)
	- Quick Experimental Install: `wget -qO- https://experimental.docker.com/ | sh`

### QuickStart Instructions (L2 Mode)

1. Start Docker with the following. **TODO:** How to specify the plugin socket without having to pass a bridge name `foo` since ipvlan/macvlan do not use traditional bridges. This example is running docker in the foreground so you can see the logs realtime.

```
    docker -d --default-network=ipvlan:foo`
```

2. Start a Git server. An easy one to use is [go git service](http://gogs.io/docs/installation/install_from_binary.md).

An example `ini` file to listen on port 80 and respond to `git clone http://username:passwd@172.16.86.1/username/git-overlay.git`

```
APP_NAME = Gogs: Go Git Service
RUN_USER = root
RUN_MODE = prod

[database]
DB_TYPE = sqlite3
HOST = 127.0.0.1:3306
NAME = gogs
USER = root
PASSWD =
SSL_MODE = disable
PATH = data/gogs.db

[repository]
ROOT = /Users/brent/gogs-repositories

[server]
DOMAIN = localhost
HTTP_PORT = 80
ROOT_URL = http://172.16.86.1:3000/
OFFLINE_MODE = false
```

Start with:
```
 sudo ./gogs web -c "custom/conf/app.ini"
```

3. Start the plugin on 2 Docker hosts. Both examples use a debug flag `-d` for lots of extra info.

# Host #1 example:
```
go run main.go -d --host-interface=eth2 --mode=l3 --ipvlan-subnet=10.1.48.0/24 --git-multihost=true --repo="http://username:passwd@172.16.86.1/username/git-overlay.git"
```

# Host #2 example:

```
go run main.go -d --host-interface=eth2 --mode=l3 --ipvlan-subnet=10.1.51.0/24 --git-multihost=true --repo="http://username:passwd@172.16.86.1/username/git-overlay.git"
```

Lastly start up some containers and check reachability:

```
docker run -i -t --rm ubuntu
```

Some example debug output is:

```
DEBU[0033] copying [ data/endpoints.latest ] endpoint cache to [ data/endpoints.old ] endpoint cache dir for diffing
DEBU[0042] Interval time expired
DEBU[0042] Running: git -C data/endpoints.latest pull
remote: Counting objects: 4, done.
remote: Compressing objects: 100% (3/3), done.
remote: Total 4 (delta 0), reused 4 (delta 0)
Unpacking objects: 100% (4/4), done.
From http://172.16.86.1/nerdalert/git-overlay
 + 3418d2a...f7aadc2 master     -> origin/master  (forced update)
Merge made by the 'recursive' strategy.
 endpoints/10.1.1.51.json | 7 +++++++
 1 file changed, 7 insertions(+)
 create mode 100644 endpoints/10.1.1.51.json
20150818023229
DEBU[0043] checking repo for updates 2015-08-18 06:32:29.761949388 +0000 UTC
DEBU[0043] New records learned [10.1.1.51.json]
INFO[0043] Endpoint [ 10.1.1.51.json ] was added
DEBU[0043] Umarshalled New Endpoint: [ 10.1.1.51 ]
DEBU[0043] Umarshalled Network [ 10.1.51.0/24 ]
DEBU[0043] Umarshalled Meta [  ]
DEBU[0043] Umarshalled SegID [ 0 ]
WARN[0043] Remote endpoint [ 10.1.1.51 ] != Local endpoint [ 10.1.1.48 ] adding route to remote endpoint [ 10.1.1.51 ] using interface [ eth2 ]
DEBU[0043] Adding route in the default namespace for IPVlan L3 mode with the following:
DEBU[0043] IP Prefix: [ 10.1.51.0/24 ] - Next Hop: [ 10.1.1.51 ] - Source Interface: [ eth2 ]
DEBU[0043] Current record list: [10.1.1.48.json 10.1.1.51.json]
INFO[0043] Records removed [ [] ]
```