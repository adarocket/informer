package common

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	pb "github.com/adarocket/proto/proto-gen/common"
)

// GetServerBasicData -
func (commonStatistic *CommonStatistic) GetServerBasicData() *pb.ServerBasicData {
	var serverBasicData pb.ServerBasicData
	var err error

	serverBasicData.Ipv4, serverBasicData.Ipv6, err = externalIP()
	if err != nil {
		log.Println(err)
	}

	distrib, err := getLinuxData()
	if err != nil {
		log.Println(err)
	}

	serverBasicData.LinuxName = distrib.ID
	serverBasicData.LinuxVersion = distrib.Version

	return &serverBasicData
}

func externalIP() (ipv4, ipv6 string, err error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ipv4, ipv6, err
	}

	for _, iface := range ifaces {
		if iface.Name != "eth0" {
			continue
		}
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}

		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return ipv4, ipv6, err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip4 := ip.To4()
			if ip4 == nil {
				ip = ip.To16()
				if ipv6 == "" {
					ipv6 = ip.String()
				}
				continue
			}
			if ipv4 == "" {
				ipv4 = ip4.String()
			}
		}

	}
	return ipv4, ipv6, err
}

type distribInfo struct {
	ID, Version string
}

func getLinuxData() (distrib *distribInfo, err error) {
	file, err := os.Open("/usr/lib/os-release")
	if err != nil {
		return distrib, err
	}
	defer file.Close()

	distrib = new(distribInfo)
	distribStats := map[string]*string{
		"NAME":    &distrib.ID,
		"VERSION": &distrib.Version,
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		i := strings.IndexRune(line, '=')
		if i < 0 {
			continue
		}
		fld := line[:i]
		if ptr := distribStats[fld]; ptr != nil {
			val := strings.TrimSpace(line[i+1:])
			*ptr = val
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan error for /proc/meminfo: %s", err)
	}

	return distrib, nil
}

// func getLinuxData() (distrib *distribInfo, err error) {
// 	file, err := os.Open("/etc/lsb-release")
// 	if err != nil {
// 		return distrib, err
// 	}
// 	defer file.Close()

// 	distrib = new(distribInfo)
// 	distribStats := map[string]*string{
// 		"DISTRIB_ID":          &distrib.ID,
// 		"DISTRIB_RELEASE":     &distrib.Version,
// 		"DISTRIB_DESCRIPTION": &distrib.Description,
// 	}

// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		i := strings.IndexRune(line, '=')
// 		if i < 0 {
// 			continue
// 		}
// 		fld := line[:i]
// 		if ptr := distribStats[fld]; ptr != nil {
// 			val := strings.TrimSpace(line[i+1:])
// 			*ptr = val
// 		}
// 	}

// 	if err := scanner.Err(); err != nil {
// 		return nil, fmt.Errorf("scan error for /proc/meminfo: %s", err)
// 	}

// 	return distrib, nil
// }
