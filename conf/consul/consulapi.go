package consul

import (
	"github.com/hashicorp/consul/api"
	"log"
)

var client *api.Client
var kv *api.KV
var session *api.Session
var agent *api.Agent
var catalog *api.Catalog

//var acl *api.ACL
//var txn *api.Txn
//
//var connect *api.Connect
//var event *api.Event
//var discoveryChain *api.DiscoveryChain
//var health *api.Health
//var operator *api.Operator

func init() {
	cli, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}
	client = cli

	kv = client.KV()
	session = client.Session()
	agent = client.Agent()
	catalog = client.Catalog()
	//acl = client.ACL()
	//txn = client.Txn()
	//
	//connect = client.Connect()
	//event = client.Event()
	//discoveryChain = client.DiscoveryChain()
	//health = client.Health()
	//operator = client.Operator()
}
