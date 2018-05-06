package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/alvelcom/redoubt/pkg/api"
)

var (
	fServer = flag.String("server", "http://127.0.0.1:2326",
		`Server to connect to`)
	fDir = flag.String("dir", "/var/run/schloss",
		`Directory for products`)
	info api.MachineInfo
)

func init() {
	info.Extra = map[string]string{}
}

func main() {
	prepareFlags()
	log.Printf("MachineInfo: %+v", info)

	c, err := api.NewClient(http.DefaultClient, *fServer, info)
	if err != nil {
		log.Printf("Can't initialize: %s", err)
		return
	}

	prods, tasks, errs, err := c.Harvest(nil)
	if err != nil {
		log.Printf("Can't harvest: %s", err)
		return
	}

	if len(errs) > 0 {
		log.Printf("Errors:")
		for _, err := range errs {
			log.Printf("%7s: %s", err.Type, err.Message)
		}
		return
	}

	log.Printf("Tasks:")
	for _, task := range tasks {
		log.Printf("- %#v", task)
	}
	log.Printf("Saving products:")
	if err := saveProducts(*fDir, prods); err != nil {
		log.Printf("saving error: %s", err)
		return
	}
}

func saveProducts(dir string, ps []api.Product) error {
	for _, p := range ps {
		name := path.Join(p.Name...)
		log.Printf("- %s (%04o)", name, p.Mask)

		name = path.Join(dir, name)
		fd, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
			os.FileMode(p.Mask))
		if err != nil {
			return err
		}

		if _, err := fd.Write(p.Body); err != nil {
			return err
		}
	}

	return nil
}

func prepareFlags() {
	ips := GetLocalIPs()
	hostInfo, _ := GetHostInfo(ips)

	fIPs := flag.String("ips", strings.Join(ips, ","),
		`Comma separated list of host IPs`)
	flag.StringVar(&info.FQDN, "fqdn", hostInfo.FQDN,
		`Machine's FQDN`)
	flag.StringVar(&info.Host, "host", hostInfo.Hostname,
		`Machine's hostname`)
	flag.StringVar(&info.Domain, "domain", hostInfo.Domain,
		`Machine's domainname`)
	flag.StringVar(&info.Cluster, "cluster", "",
		`Cluster name, that machine belonds to`)
	flag.StringVar(&info.NodeType, "node-type", "",
		``)
	flag.StringVar(&info.Id, "id", "",
		`Machine ID, usualy seq. number or hash`)
	flag.StringVar(&info.Provider, "provider", "",
		`Cloud provider for the machine`)
	flag.StringVar(&info.Region, "region", "",
		`Cloud provider's region`)
	flag.StringVar(&info.City, "city", "",
		`City name, where that machine is located`)
	flag.StringVar(&info.Country, "country", "",
		`Country name, where that machine is located`)
	flag.StringVar(&info.Geo, "geo", "",
		`Free form geographical info`)
	fExtra := flag.String("extra", "",
		`Comma separated list of key=value pairs for extra server's parameters`)

	flag.Parse()

	info.IPs = strings.Split(*fIPs, ",")
	for _, kv := range strings.Split(*fExtra, ",") {
		lst := strings.SplitN(kv, "=", 2)
		if len(lst) == 2 {
			info.Extra[lst[0]] = lst[1]
		}
	}
}