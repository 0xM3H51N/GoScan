package scanner

import (
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
}

func read(conn net.Conn) (string, error) {
	var banner strings.Builder
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if n > 0 {
			banner.Write(buf[:n])
		}
		if err != nil {
			break
		}
	}
	return banner.String(), nil
}

func worker(jobs <-chan Info, result chan<- Result, timeout int, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		add := []string{job.IP, strconv.Itoa(job.Port)}
		raddr := strings.Join(add, ":")

		conn, err := net.DialTimeout("tcp", raddr, time.Duration(timeout)*time.Second)
		if err == nil {
			conn.SetReadDeadline(time.Now().Add(time.Second * 2))
			banner, err := read(conn)
			if err != nil {
				continue
			}
			result <- Result{IP: job.IP, Port: job.Port, Status: "open", Banner: banner}
			conn.Close()
		} else {
			result <- Result{IP: job.IP, Port: job.Port, Status: "close", Banner: ""}

		}

	}
}

func Scanner(IPs []string, ports []int, workers int, timeout int) ([]Result, error) {
	var wg sync.WaitGroup

	bufferSize := len(IPs) * len(ports)
	job := make(chan Info, bufferSize)
	chResult := make(chan Result, bufferSize)

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go worker(job, chResult, timeout, &wg)
	}

	go func() {
		for _, ip := range IPs {
			for _, port := range ports {
				job <- Info{IP: ip, Port: port}
			}
		}
		close(job)
	}()

	go func() {
		wg.Wait()
		close(chResult)
	}()

	var result []Result
	for r := range chResult {
		result = append(result, r)
	}

	return result, nil
}
