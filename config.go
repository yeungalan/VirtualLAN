package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
)

type Configuration struct {
	CIDR       string `json:"cidr"`
	ListenPort int    `json:"listen_port"`
	EndPoint   string `json:"endpoint"`
}

func generateConfig(config Configuration, serverIP string) (string, string) {
	serverPrivateKey, err := generatePrivateKey()
	if err != nil {
		log.Fatal(err)
	}
	serverPublicKey := generatePublicKey(serverPrivateKey)
	serverPrivateKeyBase32 := base64.StdEncoding.EncodeToString(serverPrivateKey[:])
	serverPublicKeyBase32 := base64.StdEncoding.EncodeToString(serverPublicKey[:])

	suffix := strings.Split(config.CIDR, "/")[1]

	str := "[Interface]\n"
	str += "Address = " + serverIP + "/" + suffix + "\n"
	str += "ListenPort = " + strconv.FormatInt(int64(config.ListenPort), 10) + "\n"
	str += "PrivateKey = " + serverPrivateKeyBase32 + "\n"
	str += "\n\n"

	return str, serverPublicKeyBase32
}

func readFromFile(fileName string) (Configuration, error) {
	var config Configuration

	file, err := os.Open(fileName)
	if err != nil {
		return Configuration{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return Configuration{}, err
	}

	return config, nil
}
