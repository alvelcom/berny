package harvest

import (
	"errors"
	"net"
	"strings"
)

var ErrNoAddrFound = errors.New("LookupAddr returns empty list")

type HostInfo struct {
	Hostname string
	Domain   string
	FQDN     string
}

func GetLocalIPs() []string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	var list []string
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if !ip.IsGlobalUnicast() {
				continue
			}
			list = append(list, ip.String())
		}
	}
	return list
}

func GetHostInfo(ips []string) (hi HostInfo, err error) {
	for _, ip := range ips {
		hi, err = checkIP(ip)
		if err == nil {
			return
		}
	}
	return
}

func checkIP(ip string) (HostInfo, error) {
	info := HostInfo{}
	addrs, err := net.LookupAddr(ip)
	if err != nil {
		return info, err
	}
	if len(addrs) == 0 {
		return info, ErrNoAddrFound
	}

	info.FQDN = strings.TrimSuffix(addrs[0], ".")

	labels := strings.SplitN(info.FQDN, ".", 2)
	if len(labels) == 2 {
		info.Hostname = labels[0]
		info.Domain = labels[1]
	}
	return info, nil
}
