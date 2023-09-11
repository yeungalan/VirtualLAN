package main

import (
	"crypto/rand"
	"io"
	"net"

	"golang.org/x/crypto/curve25519"
)

func generatePrivateKey() ([32]byte, error) {
	var privateKey [32]byte
	_, err := io.ReadFull(rand.Reader, privateKey[:])
	if err != nil {
		return [32]byte{}, err
	}
	// Ensure the private key is a valid Curve25519 private key
	privateKey[0] &= 248
	privateKey[31] &= 127
	privateKey[31] |= 64
	return privateKey, nil
}

func generatePublicKey(privateKey [32]byte) [32]byte {
	var publicKey [32]byte
	curve25519.ScalarBaseMult(&publicKey, &privateKey)
	return publicKey
}

func generateSharedSecret(privateKey [32]byte, publicKey [32]byte) [32]byte {
	var sharedSecret [32]byte
	curve25519.ScalarMult(&sharedSecret, &privateKey, &publicKey)
	return sharedSecret
}

// GetAllIPsInCIDR returns a slice containing all IP addresses in the given CIDR range.
func GetAllIPsInCIDR(cidr string) ([]string, error) {
	ipRange := []string{}

	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	// Convert the IPNet IP address to a 4-byte slice
	ip := ipNet.IP.To4()

	// Get the network and broadcast addresses
	network := ipNet.IP
	broadcast := make(net.IP, len(network))
	copy(broadcast, network)
	for i := range broadcast {
		broadcast[i] |= ^network[i]
	}

	// Iterate through the IP range and add each IP address to the list
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); incrementIP(ip) {
		ipRange = append(ipRange, ip.String())
	}

	return ipRange, nil
}

// incrementIP increments the IP address by 1.
func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
