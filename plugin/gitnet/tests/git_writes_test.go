package tests

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

var jsonBlob = []byte(`[
  {
    "Endpoint": "192.168.1.254",
    "Network": "10.1.100.10/24",
    "Meta": "ipvlan:is:kewl"
  },
  {
    "Endpoint": "192.168.1.254",
    "Network": "10.1.100.10/24",
    "Meta": ""
  }
]`)

type LocalEndpoint struct {
	Endpoint string `json:"Endpoint"`
	Network  string `json:"Network"`
	Meta     string `json:"Meta"`
}

func TestConfigUnmarshall(t *testing.T) {
	//	var out bytes.Buffer
	var endpoints []LocalEndpoint
	json.Unmarshal(jsonBlob, &endpoints)
	for _, endpoint := range endpoints {
		if endpoint.Endpoint == "192.168.1.254" {
			assert.Equal(t, endpoint.Endpoint, "192.168.1.254")
			assert.Equal(t, endpoint.Network, "10.1.100.10/24")
			assert.Equal(t, endpoint.Meta, "ipvlan:is:kewl")
			assert.NotEqual(t, endpoint.Meta, "10.1.100.10/24")
		}
		break
	}
}

var jsonBlob2 = []byte(`[
  {
    "Endpoint": "192.168.1.254",
    "Network": "10.1.100.10/24",
    "Meta": "ipvlan:is:kewl"
  }
]`)

func TestConfigMarshall(t *testing.T) {
	var epConfig = &LocalEndpoint{
		Endpoint: "192.168.1.254",
		Network:  "10.1.100.10/24",
		Meta:     "ipvlan:is:kewl",
	}
	var jsonSlice []LocalEndpoint
	jsonSlice = append(jsonSlice, *epConfig)
	marshall, err := json.MarshalIndent(jsonSlice, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, jsonBlob2, marshall)
}
