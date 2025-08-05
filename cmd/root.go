/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/0xM3H51N/goscan/internal/output"
	"github.com/0xM3H51N/goscan/internal/scanner"
	"github.com/0xM3H51N/goscan/internal/utils"
	"github.com/spf13/cobra"
)

const appVersion = "GoScan v0.9.0"

func printBanner() {
	fmt.Println(`
   _____       _____                 
  / ____|     / ____|                
 | |  __  ___| (___   ___ __ _ _ __  
 | | |_ |/ _ \\___ \ / __/ _  |  _ \ 
 | |__| | (_) |___) | (_| (_| | | | |
  \_____|\___/_____/ \___\__,_|_| |_|
 Fast and Minimalist TCP Port Scanner                              
                                     `)
}

func startScanner(IPsList []string, portsList []int, workers int, timeout int, outputFile string, format string) error {

	result, err := scanner.Scanner(IPsList, portsList, workers, timeout)
	if err != nil {
		return err
	}

	err = output.WriteOutput(result, outputFile, format)
	if err != nil {
		return err
	}
	return nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goscan",
	Short: "GoScan — A fast, concurrent TCP port scanner written in Go.",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {
		ipSlice, _ := cmd.Flags().GetStringSlice("ip")
		ipsFile, _ := cmd.Flags().GetString("file")
		portSlice, _ := cmd.Flags().GetStringSlice("port")
		numWorkers, _ := cmd.Flags().GetInt("workers")
		timeout, _ := cmd.Flags().GetInt("timeout")
		format, _ := cmd.Flags().GetString("format")
		outputFile, _ := cmd.Flags().GetString("output")
		version, _ := cmd.Flags().GetBool("version")

		printBanner()

		if version {
			fmt.Println(appVersion)
			os.Exit(0)
		}

		if len(ipSlice) == 0 && ipsFile == "" {
			fmt.Println("Provide either --ip/-i or --file/-f one is required")
			cmd.Help()
			os.Exit(0)
		}

		if len(ipSlice) > 0 && ipsFile != "" {
			fmt.Println("Please provide --ip/-i or --file/-f, not both")
			cmd.Help()
			os.Exit(0)
		}

		if len(ipSlice) > 0 && len(portSlice) > 0 {

			IPsList, err := utils.IpParsing(ipSlice)
			if err != nil {
				fmt.Fprintln(os.Stderr, "IP parsing error: ", err)
				os.Exit(1)
			}

			portsList, err := utils.PortsParsing(portSlice)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Port parsing error:", err)
				os.Exit(1)
			}

			err = startScanner(IPsList, portsList, numWorkers, timeout, outputFile, format)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Scan error: ", err)
				os.Exit(1)
			}

		} else if ipsFile != "" && len(portSlice) > 0 {
			IPsList, err := utils.ExtractIPsFromFile(ipsFile)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Extracting IPs error: ", err)
				os.Exit(1)
			}

			portsList, err := utils.PortsParsing(portSlice)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error parsiing ports: ", err)
				os.Exit(1)
			}

			err = startScanner(IPsList, portsList, numWorkers, timeout, outputFile, format)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error scanning: ", err)
			}
		} else {
			cmd.Help()
			os.Exit(1)
		}

	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.Flags().StringSliceP("ip", "i", []string{},
		"Target IP addresses to scan. Supports:\n"+
			"- Single IP (e.g. 192.168.1.1)\n"+
			"- Comma-separated list (e.g. 1.1.1.1,8.8.8.8)\n"+
			"- Ranges (e.g. 192.168.1.10-20)\n"+
			"- CIDR notation (e.g. 192.168.1.0/24)")

	rootCmd.Flags().StringP("file", "f", "",
		"Path to a file containing IP addresses (one per line)")

	rootCmd.Flags().StringSliceP("port", "p", []string{},
		"Ports to scan. Supports:\n"+
			"- Single port (e.g. 22)\n"+
			"- Comma-separated list (e.g. 22,80,443)\n"+
			"- Range (e.g. 1-1024)")

	rootCmd.Flags().IntP("workers", "w", 5,
		"Number of concurrent workers to scan with.")

	rootCmd.Flags().Int("timeout", 5,
		"Connection timeout in seconds.")

	rootCmd.Flags().StringP("format", "x", "text",
		"Output format: JSON or TEXT (default: TEXT)")

	rootCmd.Flags().StringP("output", "o", "",
		"Save scan results to a file. If not set, output is printed to the terminal")

	rootCmd.Flags().BoolP("version", "v", false,
		"Show tool version and exit")
}
