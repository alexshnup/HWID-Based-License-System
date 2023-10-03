package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/jaypipes/ghw"
)

const (
	// licenseServer   = "http://127.0.0.1:9347/" //Your license server address
	licenseServer   = "http://licserver:9347/"
	PathLicenseFile = "license.dat"
	salt            = "12345salt"
)

// HardwareInfo contains relevant hardware information
type HardwareInfo struct {
	Block *ghw.BlockInfo
	Disk  string
}

var (
	hardwareInfo *HardwareInfo
	diskSerial   string
)

// GetHardwareInfo initializes and retrieves hardware information
func ghwInfo() (*HardwareInfo, error) {
	block, err := ghw.Block()
	if err != nil {
		return nil, fmt.Errorf("error getting block storage info: %v", err)
	}

	diskSerial := ""
	if len(block.Disks) > 0 {
		diskSerial = block.Disks[0].SerialNumber
	}

	return &HardwareInfo{
		Block: block,
		Disk:  diskSerial,
	}, nil
}

func getHardwareInfo() error {
	var err error
	hardwareInfo, err = ghwInfo()
	if err != nil {
		return fmt.Errorf("getting hardware info: %v", err)
	}
	diskSerial = hardwareInfo.Disk
	return nil
}

func checkFileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func hashSHA256(text string) string {
	hasher := sha256.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func sendRequestToLicenseServer(data url.Values) (*http.Response, error) {
	u, err := url.ParseRequestURI(licenseServer)
	if err != nil {
		return nil, fmt.Errorf("parsing license server URI: %v", err)
	}

	r, err := http.NewRequest("POST", u.String(), bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("creating new request: %v", err)
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	return client.Do(r)
}

func getHWID() (string, error) {
	name, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("getting hostname: %v", err)
	}

	if len(diskSerial) == 0 {
		return "", errors.New("disk serial number is empty")
	}

	return hashSHA256(name + diskSerial), nil
}

func getLicenseKey() (string, error) {
	if checkFileExist(PathLicenseFile) {
		data, err := os.ReadFile(PathLicenseFile)
		if err != nil {
			return "", fmt.Errorf("reading license file: %v", err)
		}
		return string(data), nil
	}

	key := os.Getenv("License")
	if len(key) == 0 {
		return "", errors.New("license key is not provided in the environment variable")
	}

	return key, nil
}

func licenseCheck(key string, hwid string) error {
	data := url.Values{}
	data.Set("license", key)
	data.Set("hwid", hwid)

	resp, err := sendRequestToLicenseServer(data)
	if err != nil {
		return fmt.Errorf("sending request to license server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code from license server: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}

	fmt.Printf("Response from license server: %s\n", respBody)

	//expect
	// fmt.Printf("Response from license server expired: %s\n", hashSHA256("1"+salt+hwid))
	// fmt.Printf("Response from license server valid New Register: %s\n", hashSHA256("2"+salt+hwid))
	// fmt.Printf("Response from license server good license: %s\n", hashSHA256("0"+salt+hwid))

	switch string(respBody) {
	case hashSHA256("0" + salt + hwid):
		fmt.Println("License Good")
		return nil
	case hashSHA256("1" + salt + hwid):
		return errors.New("license is expired")
	case hashSHA256("2" + salt + hwid):
		if !checkFileExist(PathLicenseFile) {
			fmt.Println("License file not found.")
			fmt.Print("Try activate license from Env")
			os.WriteFile(PathLicenseFile, []byte(key), 0600)
		}
		return nil // License is valid and registered
	default:
		return errors.New("unable to verify license server response")
	}
}

func main() {
	if err := getHardwareInfo(); err != nil {
		log.Fatalf("Error: %v", err)
	}

	hwid, err := getHWID()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("HWID: %s", hwid)

	key, err := getLicenseKey()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if err := licenseCheck(key, hwid); err != nil {
		log.Fatalf("License check failed: %v", err)
	}

	log.Println("License OK")
	time.Sleep(1 * time.Second)
}
