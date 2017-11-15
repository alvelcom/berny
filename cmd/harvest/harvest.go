package harvest

import (
	"flag"
	"log"
	"strings"

	"github.com/alvelcom/redoubt/api"
)

var fs = flag.NewFlagSet("redoubt harvest", flag.ExitOnError)
var (
	fServer = fs.String("server", "127.0.0.1:2326", "Server to connect to")
	info    api.MachineInfo
)

func init() {
	info.Extra = map[string]string{}
}

func Main(args []string) {
	prepareFlags(args)
	log.Printf("MachineInfo: %+v", info)

}

func prepareFlags(args []string) {
	ips := GetLocalIPs()
	hostInfo, _ := GetHostInfo(ips)

	fIPs := fs.String("ips", strings.Join(ips, ","),
		`Comma separated list of host IPs`)
	fs.StringVar(&info.FQDN, "fqdn", hostInfo.FQDN,
		`Machine's FQDN`)
	fs.StringVar(&info.Host, "host", hostInfo.Hostname,
		`Machine's hostname`)
	fs.StringVar(&info.Domain, "domain", hostInfo.Domain,
		`Machine's domainname`)
	fs.StringVar(&info.Cluster, "cluster", "",
		`Cluster name, that machine belonds to`)
	fs.StringVar(&info.NodeType, "node-type", "",
		``)
	fs.StringVar(&info.Id, "id", "",
		`Machine ID, usualy seq. number or hash`)
	fs.StringVar(&info.Provider, "provider", "",
		`Cloud provider for the machine`)
	fs.StringVar(&info.Region, "region", "",
		`Cloud provider's region`)
	fs.StringVar(&info.City, "city", "",
		`City name, where that machine is located`)
	fs.StringVar(&info.Country, "country", "",
		`Country name, where that machine is located`)
	fs.StringVar(&info.Geo, "geo", "",
		`Free form geographical info`)
	fExtra := fs.String("extra", "",
		`Comma separated list of key=value pairs for extra server's parameters`)

	fs.Parse(args)

	info.IPs = strings.Split(*fIPs, ",")
	for _, kv := range strings.Split(*fExtra, ",") {
		lst := strings.SplitN(kv, "=", 2)
		if len(lst) == 2 {
			info.Extra[lst[0]] = lst[1]
		}
	}
}
