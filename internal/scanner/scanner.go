package scanner

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Info struct {
	IP   string
	Port int
}

type Result struct {
	IP     string
	Port   int
	Status string
	Banner string
	err    error
}

type ErrorMessages struct {
	TimeoutDuration    string
	PortConnection     string
	UnreachableHost    string
	UnreachableNetwork string
}

func read(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	banner, err := reader.ReadString('\n')

	if err != nil {
		return "", fmt.Errorf("error reader banner: %w", err)
	}

	return banner, nil
}

func worker(jobs <-chan Info, result chan<- Result, timeout int, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		add := []string{job.IP, strconv.Itoa(job.Port)}
		raddr := strings.Join(add, ":")

		conn, err := net.DialTimeout("tcp", raddr, time.Duration(timeout))
		if err == nil {
			banner, err := read(conn)
			if err != nil {
				fmt.Errorf("could not grab banner: %v", err)
			}
			result <- Result{IP: job.IP, Port: job.Port, Status: "open", Banner: banner}
			conn.Close()
		} else {
			result <- Result{IP: job.IP, Port: job.Port, Status: "close", err: err}

		}

	}
}

func Scanner(IPs []string, ports []int, workers int, timeout int) ([]Result, error) {
	var wg sync.WaitGroup

	job := make(chan Info)
	chResult := make(chan Result)

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go worker(job, chResult, timeout, &wg)
	}

	for _, ip := range IPs {
		for _, port := range ports {
			job <- Info{IP: ip, Port: port}
		}
	}
	close(job)

	total := (len(IPs) * len(ports))
	var result []Result
	for i := 0; i < total; i++ {
		res := <-chResult
		result = append(result, res)
	}
	close(chResult)

	wg.Wait()

	return result, nil
}
