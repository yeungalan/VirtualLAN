package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

var serverPublicKey string
var serverConfigHeader string
var ipRange []string
var config Configuration

func main() {
	var err error

	config, err = readFromFile("config.json")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	ipRange, err = GetAllIPsInCIDR(config.CIDR)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(ipRange) < 2 {
		log.Fatal("CIDR size == 0")
	}

	ipRange = ipRange[1:]

	serverConfigHeader, serverPublicKey = generateConfig(config, ipRange[0])

	log.Println("Server IP", ipRange[0])
	ipRange = ipRange[1:]

	http.HandleFunc("/getClient", getClient)
	//http.HandleFunc("/getServer", getServer)

	port := 8080
	fmt.Printf("Server is listening on port %d...\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func getClient(w http.ResponseWriter, r *http.Request) {
	//username, _ := mv(r, "username", false)
	//password, _ := mv(r, "password", false)
	//computerId, _ := mv(r, "computerId", false)

	//fmt.Println(ipRange)
	if len(ipRange) == 0 {
		sendErrorResponse(w, "CIDR full")
		return
	}

	clientIP := ipRange[0]
	log.Println("New Client joining...! assigned ip", clientIP)

	clientConfigHeader, clientPublicKey := generateConfig(config, ipRange[0])
	ipRange = ipRange[1:]

	clientConfigHeader += "[Peer]\n"
	clientConfigHeader += "PublicKey = " + serverPublicKey + "\n"
	clientConfigHeader += "AllowedIPs = " + config.CIDR + "\n"
	clientConfigHeader += "Endpoint = " + config.EndPoint + ":" + strconv.FormatInt(int64(config.ListenPort), 10) + "\n"

	serverConfigHeader += "[Peer]\n"
	serverConfigHeader += "PublicKey = " + clientPublicKey + "\n"
	serverConfigHeader += "AllowedIPs = " + clientIP + "/32\n"
	serverConfigHeader += "\n"

	restartServer()

	sendTextResponse(w, clientConfigHeader)
}

func getServer(w http.ResponseWriter, r *http.Request) {
	sendTextResponse(w, serverConfigHeader)
}

func restartServer() {
	log.Println("Server restarting...")

	cmd := exec.Command("wireguard.exe", "/uninstalltunnelservice", "wg2")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		//fmt.Printf("Failed to run WireGuard command: %v\n", err)
		//return
	}
	log.Println("WireGuard service stopped")

	os.Remove("wg2.conf")

	// Create the destination file
	file, err := os.Create("wg2.conf")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	// Copy the contents of the response body to the destination file
	_, err = io.WriteString(file, serverConfigHeader)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	log.Println("Config updated")

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	log.Println("Restarting...")

	// Run the first WireGuard command
	cmd = exec.Command("wireguard.exe", "/installtunnelservice", path+"/wg2.conf")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Failed to run WireGuard command: %v\n", err)
		return
	}

	log.Println("Done...")

}
