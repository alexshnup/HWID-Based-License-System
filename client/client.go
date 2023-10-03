package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"HWID-Based-License-System/client/hwinfo"
)

var (
	licenseServer   string = "http://127.0.0.1:9347/" //Your license server address
	PathLicenseFile string = "license.dat"
	hardwareInfo    *hwinfo.HardwareInfo
	diskSerial      string
	salt            string = "12345salt"
)

func logDirectoryInfo(path string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return fmt.Errorf("error reading directory: %v", err)
	}

	for _, file := range files {
		log.Println(file.Name(), file.IsDir())
	}
	return nil
}

func checkFileExist(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func md5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func init() {
	log.Println("HwinfoInit...")

	// Log directory information
	if err := logDirectoryInfo("./"); err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Initialize hardware information
	hardwareInfo, err := hwinfo.GetHardwareInfo()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("Block: %+v, Disk Serial: %s", hardwareInfo.Block, hardwareInfo.Disk)
	diskSerial = hardwareInfo.Disk
	log.Println("HwinfoInit OK")
}

func LicenseCheck() {
	fmt.Println("LC....")
	key := ""
	hdd := ""
	if len(diskSerial) > 0 {
		hdd = diskSerial
	}

	name, _ := os.Hostname()
	// usr, _ := user.Current()
	// usr := fmt.Sprint(rutoken.SerialRutoken)
	usr := name

	if !checkFileExist(PathLicenseFile) {

		fmt.Println("License file not found.")

		fmt.Print("Try activate license from Env")
		// scan := bufio.NewScanner(os.Stdin)
		// scan.Scan()

		if key == "" {
			// key := scan.Text()
			key = os.Getenv("License")
			if len(key) == 0 {
				fmt.Print("Please set \"License\" env with key")
				os.Exit(0)
			}
		}

		fmt.Println("Key:", key)

		os.WriteFile(PathLicenseFile, []byte(key), 0600)

		// fmt.Println("HWID:", md5Hash(name+usr.Username))
		hwid := md5Hash(name + usr + hdd)
		fmt.Println("HWID:", hwid)

		fmt.Println("Connecting to license server...")

		client := &http.Client{}
		data := url.Values{}
		data.Set("license", key)
		data.Set("hwid", hwid)

		// fmt.Println("0", salt, md5Hash("0"+salt+hwid))
		// fmt.Println("1", salt, md5Hash("1"+salt+hwid))
		// fmt.Println("2", salt, md5Hash("2"+salt+hwid))
		// fmt.Println("data.Encode()", data.Encode())

		u, _ := url.ParseRequestURI(licenseServer)
		urlStr := fmt.Sprintf("%v", u)
		r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp, err := client.Do(r)
		if err != nil {
			fmt.Println("Unable to connect to license server.")
			os.Exit(0)
		} else {
			defer resp.Body.Close()
			resp_body, _ := ioutil.ReadAll(resp.Body)
			if resp.StatusCode == 200 {
				if string(resp_body) != md5Hash("0"+salt+hwid) {
					if string(resp_body) == md5Hash("1"+salt+hwid) {
						fmt.Println("License is Expired.")

						os.Exit(0)
					} else if string(resp_body) == md5Hash("2"+salt+hwid) {
						fmt.Println("Registered!")

						fmt.Println("DO NOT DELETE THE FILE! " + PathLicenseFile)
						fmt.Println(" ")
					} else {
						fmt.Println("Unable to verify to license server.")

						os.Exit(0)
					}
				}
			} else {
				fmt.Println(resp.StatusCode)
			}
		}

		// } else {
		// 	os.Exit(0)
		// }
	} else {

		dat, err := os.ReadFile(PathLicenseFile)
		if err != nil {
			log.Fatalf("failed reading data from file: %s", err)
		}

		key := string(dat)
		hwid := md5Hash(name + usr + hdd)

		fmt.Println("HWID:", hwid)

		client := &http.Client{}
		data := url.Values{}
		data.Set("license", key)
		// data.Set("hwid", md5Hash(name+usr.Username))
		data.Set("hwid", hwid)
		u, _ := url.ParseRequestURI(licenseServer)
		urlStr := fmt.Sprintf("%v", u)
		r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp, err := client.Do(r)
		if err != nil {
			fmt.Println("Unable to connect to license server.")

			os.Exit(0)
		} else {
			defer resp.Body.Close()
			resp_body, _ := ioutil.ReadAll(resp.Body)
			if resp.StatusCode == 200 {
				if string(resp_body) != md5Hash("0"+salt+hwid) {
					if string(resp_body) == md5Hash("1"+salt+hwid) {
						fmt.Println("License is Expired.")

						os.Exit(0)
					} else {
						fmt.Println("Unable to verify to license server.")
						os.Exit(0)
					}
				}
			} else {
				fmt.Println("Unable connect to license server.")
				os.Exit(0)
			}
		}

	}
}

func main() {
	LicenseCheck()
	fmt.Println("License OK")
	time.Sleep(1 * time.Second)
}
