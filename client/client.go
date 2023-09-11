package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	log.Println("Starting wireguard client...")
	var url, destFileName string

	url = "http://10.0.0.29:8080/getClient"
	destFileName = "wg2.conf"

	// Fetch the INI file from the chosen URL
	log.Println("Obtaining config...")

	err := downloadINIFile(url, destFileName)
	if err != nil {
		fmt.Printf("Failed to fetch INI file: %v\n", err)
		return
	}
	log.Println("Obtained...")

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	log.Println("Starting wireguard...")

	// Run the first WireGuard command
	cmd := exec.Command("wireguard.exe", "/installtunnelservice", path+"/"+destFileName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Failed to run WireGuard command: %v\n", err)
		return
	}
	fmt.Println("Press any key to stop...")

	// Wait for user input
	var input string
	fmt.Scanln(&input)

	// Run the second WireGuard command
	cmd = exec.Command("wireguard.exe", "/uninstalltunnelservice", "wg2")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Failed to run WireGuard command: %v\n", err)
		return
	}
	fmt.Println("WireGuard service uninstalled, press any key to quit")

	os.Remove("wg2.conf")
	fmt.Scanln(&input)
}

func downloadINIFile(url, destFileName string) error {
	// Fetch the file from the URL
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the destination file
	file, err := os.Create(destFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy the contents of the response body to the destination file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
