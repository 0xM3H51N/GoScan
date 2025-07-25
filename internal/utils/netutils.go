package utils

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func checkDublicates() {

}

func validateIP(ip string) error {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return fmt.Errorf("IP could not be parsed as an IP")
	}

	return nil
}

func ipParsing() {

}

func validatePortsInRange(port int) error {

	if port < 1 || port > 65535 {
		return fmt.Errorf("Port is not in range 1-65535")
	}
	return nil
}

func getPortsRange(rng string) ([]int, error) {
	splitRange := strings.Split(rng, "-")

	firstPort := splitRange[0]
	lastPort := splitRange[len(splitRange)-1]

	fp, err1 := strconv.Atoi(firstPort)
	lp, err2 := strconv.Atoi(lastPort)

	if err1 != nil || err2 != nil {
		log.Printf("Ports range input error: %v, %v", err1, err2)
		return nil, fmt.Errorf("Ports range input error: %w, %w", err1, err2)
	}

	err1 := validatePortsInRange(fp)
	err2 := validatePortsInRange(lp)
	if err1 != nil || err2 != nil {

	}

	var ports []int
	for i := fp; i <= lp; i++ {
		ports = append(ports, i)
	}
	return ports, nil
}

func portsParsing(portSlice []string) ([]int, error) {
	var ports []int
	for _, p := range portSlice {
		if strings.Contains(p, "-") {
			portsRange, err := getPortsRange(p)
			if err != nil {
				continue
			}

			for _, p := range portsRange {
				ports = append(ports, p)
			}
		} else {

			ports = append(ports, p)
		}
	}

	return ports, nil
}

func extractIPsFromFile(file string) ([]string, error) {
	fd, err := os.Open(file)
	if err != nil {
		return []string{}, fmt.Errorf("Error openning file %s: %w", file, err)
	}

	defer fd.Close()

	scanner := bufio.NewScanner(fd)

	var ipList []string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		line = strings.TrimSpace(line)
		err := validateIP(line)
		if err != nil {
			log.Printf("Found a non valid: %s, %v", line, err)
			continue
		}

		ipList = append(ipList, line)
	}

	if err := scanner.Err(); err != nil {
		return []string{}, fmt.Errorf("Error reading from file: %w", err)
	}

	return ipList, nil
}
