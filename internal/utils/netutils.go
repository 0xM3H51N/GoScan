package utils

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"slices"
	"strconv"
	"strings"
)

func checkDuPlicatesIPs(l []string) []string {

	slices.Sort(l)
	list := slices.Compact(l)

	return list
}

func validateIP(ip string) error {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return fmt.Errorf("invalid IP address")
	}

	return nil
}

func validateRange(i int, s string) error {
	switch {
	case s == "port":
		if i < 1 || i > 65535 {
			return fmt.Errorf("port is not in range 1-65535")
		}
		return nil
	case s == "ip":
		if i < 1 || i > 255 {
			return fmt.Errorf("IP is not in range 1-255")
		}
		return nil
	default:
		return fmt.Errorf("range validation error")
	}

}

func IP2Int(ip net.IP) []int {
	var intIP [4]int
	for i, b := range ip {
		intIP[i] = int(b)
	}
	return intIP[:]
}

func getIPsByCIDR(ip string) ([]string, error) {

	_, IPNetwork, err := net.ParseCIDR(ip)
	if err != nil {
		return nil, fmt.Errorf("error parsing CIDR: %w", err)
	}

	IPv4 := IPNetwork.IP.To4()
	mask := IPNetwork.Mask
	broadcast := make(net.IP, len(IPv4))

	for i := 0; i < len(IPv4); i++ {
		broadcast[i] = IPv4[i] | ^mask[i]
	}

	ipCounts := make(net.IP, len(IPv4))
	for i := 0; i < len(IPv4); i++ {
		if IPv4[i] == 0 {
			IPv4[i] = 1
		}
		ipCounts[i] = (broadcast[i] - IPv4[i]) + 1
	}

	baseIP := IP2Int(IPv4)
	ipRng := IP2Int(ipCounts)

	var list []string
	for i := baseIP[0]; i < (baseIP[0] + ipRng[0]); i++ {
		for j := baseIP[1]; j < (baseIP[1] + ipRng[1]); j++ {
			for k := baseIP[2]; k < (baseIP[2] + ipRng[2]); k++ {
				for l := baseIP[3]; l < (baseIP[3] + ipRng[3]); l++ {
					c := []int{i, j, k, l}
					var arr [4]string
					for x, s := range c {
						str := strconv.Itoa(s)
						arr[x] = str
					}
					strarr := strings.Join(arr[:], ".")
					list = append(list, strarr)
				}
			}
		}
	}
	return list, nil
}

func getNumberOfIPs(ip []string) (int, []int, []int, error) {
	var ipCounts [4]int
	var startRng [4]int
	for i, str := range ip {
		if strings.Contains(str, "-") {
			strSplited := strings.Split(str, "-")
			start, err1 := strconv.Atoi(strSplited[0])
			end, err2 := strconv.Atoi(strSplited[1])
			if err1 != nil || err2 != nil {
				return 0, nil, nil, fmt.Errorf("could convert IP range to int: %w, %w", err1, err2)
			}
			ipCounts[i] = (end - start) + 1
			startRng[i] = start
		} else {
			ipCounts[i] = 1
			n, _ := strconv.Atoi(str)
			startRng[i] = n
		}
	}
	numberOfIPs := (((ipCounts[3] * ipCounts[2]) * ipCounts[1]) * ipCounts[0])
	return numberOfIPs, ipCounts[:], startRng[:], nil
}

func getIPsByRange(ip string) ([]string, error) {

	ipSplit := strings.Split(ip, ".")
	if len(ipSplit) != 4 {
		return nil, fmt.Errorf("invalid IP format: %s", ip)
	}

	_, ipCounts, startIP, err := getNumberOfIPs(ipSplit)
	if err != nil {
		return nil, fmt.Errorf("error getting IPs range: %w", err)
	}

	var ipList []string

	for i := startIP[0]; i < (startIP[0] + ipCounts[0]); i++ {
		for j := startIP[1]; j < (startIP[1] + ipCounts[1]); j++ {
			for k := startIP[2]; k < (startIP[2] + ipCounts[2]); k++ {
				for l := startIP[3]; l < (startIP[3] + ipCounts[3]); l++ {
					c := []int{i, j, k, l}
					var arr [4]string
					for x, s := range c {
						str := strconv.Itoa(s)
						arr[x] = str
					}
					strarr := strings.Join(arr[:], ".")
					ipList = append(ipList, strarr)

				}
			}
		}
	}

	return ipList, nil
}

func IpParsing(ipSlice []string) ([]string, error) {
	var fList []string
	for _, ip := range ipSlice {
		if strings.Contains(ip, "-") {

			IPsRange, err := getIPsByRange(ip)
			if err != nil {
				continue
			}

			fList = append(fList, IPsRange...)

		} else if strings.Contains(ip, "/") {

			list, err := getIPsByCIDR(ip)
			if err != nil {
				continue
			}

			fList = append(fList, list...)

		} else {
			err := validateIP(ip)
			if err != nil {
				continue
			}
			fList = append(fList, ip)

		}
	}

	ipList := checkDuPlicatesIPs(fList)

	return ipList, nil
}

func getPortsRange(rng string) ([]int, error) {
	splitRange := strings.Split(rng, "-")

	firstPort := splitRange[0]
	lastPort := splitRange[len(splitRange)-1]

	fp, err1 := strconv.Atoi(firstPort)
	lp, err2 := strconv.Atoi(lastPort)

	if err1 != nil || err2 != nil {
		return nil, fmt.Errorf("ports range input error: %w, %w", err1, err2)
	}

	err1 = validateRange(fp, "port")
	err2 = validateRange(lp, "port")
	if err1 != nil || err2 != nil {
		return nil, fmt.Errorf("range input error: %w, %w", err1, err2)
	}

	var ports []int
	for i := fp; i <= lp; i++ {
		ports = append(ports, i)
	}
	return ports, nil
}

func PortsParsing(portSlice []string) ([]int, error) {
	var ports []int
	for _, p := range portSlice {
		if strings.Contains(p, "-") {
			portsRange, err := getPortsRange(p)
			if err != nil {
				continue
			}

			ports = append(ports, portsRange...)

		} else {
			port, err := strconv.Atoi(p)
			if err != nil {
				continue
			}
			err = validateRange(port, "port")
			if err != nil {
				continue
			}
			ports = append(ports, port)
		}
	}

	if len(ports) == 0 {
		return nil, fmt.Errorf("port input Error: check your input")
	}

	return ports, nil
}

func ExtractIPsFromFile(file string) ([]string, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("error openning file %s: %w", file, err)
	}

	defer fd.Close()

	scanner := bufio.NewScanner(fd)

	var list []string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		line = strings.TrimSpace(line)
		err := validateIP(line)
		if err != nil {
			log.Printf("Invalid IP skipped: %s, (%v)", line, err)
			continue
		}

		list = append(list, line)
	}

	if err := scanner.Err(); err != nil {
		return []string{}, fmt.Errorf("error reading from file: %w", err)
	}

	ipList := checkDuPlicatesIPs(list)

	return ipList, nil
}
