package Helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"
)

const (
	serverAddress = "lic-auth-usa1.botbuddy.net:7843"
)

var key = []byte("AES_LOLWAT_KEYBB")
var authenticated = false

func License() bool {
	licenseKey, err := readLicenseKey()
	if err != nil {
		fmt.Println("Your license.key file is missing. It should be next to your BotBuddy.exe file.")
	}

	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("License server is offline. Please contact the team on Discord.")
		return false
	}

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	block, err := aes.NewCipher(key)
	if err != nil {
		return false
	}

	// Encrypt license key
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return false
	}

	nonceSize := 12
	nonce := make([]byte, nonceSize)
	_, err = rand.Read(nonce)
	if err != nil {
		return false
	}

	ciphertext := aesgcm.Seal(nil, nonce, []byte(licenseKey), nil)

	// Send encrypted license key
	encryptedMsg := append(nonce, ciphertext...)
	msgLen := uint32(len(encryptedMsg))
	err = binary.Write(conn, binary.LittleEndian, msgLen)
	if err != nil {
		fmt.Println("Problem sending encrypted message to license server:", err)
		return false
	}

	_, err = conn.Write(encryptedMsg)
	if err != nil {
		fmt.Println("Problem sending encrypted message to license server:", err)
		return false
	}

	// Receive encrypted response
	var responseLen uint32
	err = binary.Read(conn, binary.LittleEndian, &responseLen)
	if err != nil {
		fmt.Println("Problem receiving encrypted response from license server:", err)
		return false
	}

	encryptedResponse := make([]byte, responseLen)
	_, err = conn.Read(encryptedResponse)
	if err != nil {
		fmt.Println("Problem receiving encrypted response from license server:", err)
		return false
	}

	// Decrypt response
	nonce, ciphertext = encryptedResponse[:nonceSize], encryptedResponse[nonceSize:]
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println("Problem decrypting response from license server:", err)
		return false
	}

	authenticated = !strings.Contains(string(plaintext), "INVALID")
	return authenticated
}

func readLicenseKey() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}

	licensePath := filepath.Join(filepath.Dir(exe), "license.key")
	b, err := ioutil.ReadFile(licensePath)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
