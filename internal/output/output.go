package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/0xM3H51N/goscan/internal/scanner"
)

type PortResult struct {
	Port   int    `json:"port"`
	Status string `json:"status"`
	Banner string `json:"banner"`
}

func formatAsJson(results map[string][]PortResult, fileName string) error {
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}

	if fileName == "" {
		fmt.Println(string(jsonData))
	} else {
		fd, err := os.Create(fileName)
		if err != nil {
			return err
		}

		defer fd.Close()

		_, err = fd.WriteString(string(jsonData))
		if err != nil {
			return err
		}
	}

	return nil
}

func formatAsPlain(results map[string][]scanner.Result, fileName string) error {

	if fileName == "" {
		for ip, r := range results {
			fmt.Println("=================")
			fmt.Printf("%s\n", ip)
			fmt.Println("=================")
			fmt.Printf("port\tstatus\tversion\n")
			for _, s := range r {
				fmt.Printf("%d\t%s\t%s\n", s.Port, s.Status, s.Banner)
			}
		}
	} else {
		fd, err := os.Create(fileName)
		if err != nil {
			return err
		}

		defer fd.Close()

		for ip, r := range results {
			fmt.Fprintf(fd, "=================\n")
			fmt.Fprintf(fd, "%s\n", ip)
			fmt.Fprintf(fd, "=================\n")
			fmt.Fprintf(fd, "port\tstatus\tversion\n")
			for _, s := range r {
				fmt.Fprintf(fd, "%d\t%s\t%s\n", s.Port, s.Status, s.Banner)
			}
		}
	}
	return nil
}

func WriteOutput(res []scanner.Result, fileName string, format string) error {

	f := strings.ToLower(format)
	switch f {
	case "json":
		results := make(map[string][]PortResult)
		for _, r := range res {
			pr := PortResult{Port: r.Port, Status: r.Status, Banner: r.Banner}
			results[r.IP] = append(results[r.IP], pr)
		}
		err := formatAsJson(results, fileName)
		if err != nil {
			return err
		}
	case "plaintext":
		results := make(map[string][]scanner.Result)
		for _, r := range res {
			results[r.IP] = append(results[r.IP], r)
		}
		err := formatAsPlain(results, fileName)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown output format: %s", format)
	}
	return nil
}
