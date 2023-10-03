package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var (
	// PORT int = 9347
	PORT int    = 9347
	salt string = "salt"
)

func checkFileExist(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func CheckFileExist(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func createFile(pathFile string) error {
	file, err := os.Create(pathFile)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func randomString(n int) string {
	var letterRunes = []rune("1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func md5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// LICENSE:EXPDATE:EMAIL:HWID = 0,1,2
func checkHandler(response http.ResponseWriter, request *http.Request) {

	request.ParseForm()
	license := request.FormValue("license")
	hwid := request.FormValue("hwid")

	database, _ := readLines("/app/db/dbfile")

	fmt.Printf("\nrequest.RequestURI=%v\n", request.RequestURI)

	for _, table := range database {

		row := strings.Split(table, ":")

		t, err := time.Parse("2006-01-02", row[1])
		if err != nil {
			fmt.Println("ERROR: Error reading database")
		}

		t2, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))

		if license == row[0] && t.After(t2) {
			if hwid == row[3] {
				log.Printf("%v from %v Registed, Good licnese %s request.URL.Path %s", time.Now(), request.RemoteAddr, license, request.URL.Path)
				// fmt.Fprintf(response, "0") //Registed, Good licnese

				out := ""
				m5 := md5Hash("0" + salt + hwid)

				switch request.URL.Path {
				case "/":
					fmt.Printf("\nm5=%s\n", m5)
					out = m5
					// case "/gat":
					// 	fmt.Printf("\nm5=%s\n", m5)
					// 	out = m5 + DecryptString(gantner)
					// 	fmt.Printf("\nout=%s\n", out)
				}

				// out := md5Hash(md5Hash("0") + md5Hash("solt2022"+stlsolt))
				fmt.Fprintf(response, out) //Registed, Good licnese
				return
			} else if row[3] == "NOTSET" {
				b, err := ioutil.ReadFile("/app/db/dbfile")
				if err != nil {
					fmt.Println("READfromCHECK")
					os.Exit(0)
				}

				str := string(b)
				edit := row[0] + ":" + row[1] + ":" + row[2] + ":" + hwid
				res := strings.Replace(str, table, edit, -1)

				err = ioutil.WriteFile("/app/db/dbfile", []byte(res), 0644)
				if err != nil {
					fmt.Println("WRITEfromCHECK")
					os.Exit(0)
				}

				log.Printf("%v from %v NEW Registed, Good licnese %s", time.Now(), request.RemoteAddr, license)
				// fmt.Fprintf(response, "2") //Registed, Good licnese
				out := md5Hash("2" + salt + hwid)

				fmt.Fprintf(response, out) //Registed, Good licnese
				return

			}
		} else if license == row[0] && !t.After(t2) {

			log.Printf("%v from %v registerd but license %s", time.Now(), request.RemoteAddr, license)
			// fmt.Fprintf(response, "1") //registerd but license experied

			out := md5Hash("1" + salt + hwid)
			fmt.Fprintf(response, out) //Registed, Good licnese
			return
		}

		if license == row[0] && hwid != row[3] {
			fmt.Printf("from %v license %s problem hwid", request.RemoteAddr, license)
		}
	}
}

func serverAPI() {
	router := mux.NewRouter()
	router.HandleFunc("/", checkHandler).Methods("POST")
	http.Handle("/", router)

	http.ListenAndServe(":"+string(strconv.Itoa(PORT)), nil)
}

func GetNewLicNumber() string {
	return randomString(4) + "-" + randomString(4) + "-" + randomString(4)
}

func LicExist(toFind string) bool {
	database, _ := readLines("/app/db/dbfile")
	for _, v := range database {
		if strings.Contains(v, toFind) {
			return true
		}
	}
	return false
}

func main() {
	salt = os.Getenv("SALT")
	if salt == "" {
		fmt.Println("SALT not set")
		os.Exit(0)
	}
	fmt.Println("License Server")
	fmt.Println("Github: https://github.com/alexshnup/easy-license-system")
	fmt.Println("Forked from: https://github.com/SaturnsVoid/HWID-Based-License-System")

	if !checkFileExist("/app/db/dbfile") {
		fmt.Println("Database does not exist, creating new database.")
		_ = createFile("/app/db/dbfile")
	}

	database, _ := readLines("/app/db/dbfile")

	fmt.Println("Total Licenses:", len(database))

	go serverAPI()
	for {
		fmt.Println(" ")
		fmt.Print("$> ")
		scan := bufio.NewScanner(os.Stdin)
		scan.Scan()
		switch scan.Text() {
		case "list":
			database, _ = readLines("/app/db/dbfile")
			for _, table := range database {
				fmt.Println(table)
			}
		case "add":
			var email string
			var experation string
			var license string

			fmt.Print("License Email: ")
			scan = bufio.NewScanner(os.Stdin)
			scan.Scan()
			email = scan.Text()

		exp:
			fmt.Print("License Experation (YYYY-MM-DD): ")
			scan = bufio.NewScanner(os.Stdin)
			scan.Scan()
			_, err := time.Parse("2006-01-02", scan.Text())
			if err != nil {
				fmt.Println("Experation must be in the YYYY-MM-DD Format.")
				goto exp
			}
			experation = scan.Text()

			b, err := ioutil.ReadFile("/app/db/dbfile")
			if err != nil {
				os.Exit(0)
			}

			license = randomString(4) + "-" + randomString(4) + "-" + randomString(4)
			for LicExist(license) {
				fmt.Println("Generating...")
				license = randomString(4) + "-" + randomString(4) + "-" + randomString(4)
			}

			str := string(b)
			str = str + "\r\n" + license + ":" + experation + ":" + email + ":NOTSET"

			re := regexp.MustCompile("(?m)^\\s*$[\r\n]*")
			str2 := strings.Trim(re.ReplaceAllString(str, ""), "\r\n")

			err = ioutil.WriteFile("/app/db/dbfile", []byte(str2), 0644)
			if err != nil {
				os.Exit(0)
			}

			fmt.Println("New License Generated:", license, "for", email)
		case "add bulk":
			var experation string

			fmt.Println("Bulk accounts will be added to database without emails. You can add emails at a later time.")
			fmt.Println(" ")

			reader := bufio.NewReader(os.Stdin)
			fmt.Print("How many keys to generate? (#): ")
			bytes, _, err := reader.ReadLine()
			if err != nil {
				os.Exit(0)
			}

			amount := string(bytes)

			n, err := strconv.Atoi(amount)
			if err != nil {
				os.Exit(0)
			}

		expb:
			fmt.Print("License Experation ( YYYY-MM-DD ): ")
			scan = bufio.NewScanner(os.Stdin)
			scan.Scan()
			if len(scan.Text()) == 0 {

				experation = "2019-12-30"
				fmt.Printf("Default %s", experation)
				goto expbdefault
			}
			_, err = time.Parse("2006-01-02", scan.Text())
			if err != nil {
				fmt.Println("Experation must be in the YYYY-MM-DD Format.")
				goto expb
			}
			experation = scan.Text()

		expbdefault:
			for i := 0; i < n; i++ {
			restart:
				var old string
				license := randomString(4) + "-" + randomString(4) + "-" + randomString(4)
				if license != old {
					b, err := ioutil.ReadFile("/app/db/dbfile")
					if err != nil {
						os.Exit(0)
					}

					str := string(b)
					str = str + "\r\n" + license + ":" + experation + ":null" + ":NOTSET"

					re := regexp.MustCompile("(?m)^\\s*$[\r\n]*")
					str2 := strings.Trim(re.ReplaceAllString(str, ""), "\r\n")

					err = ioutil.WriteFile("/app/db/dbfile", []byte(str2), 0644)
					if err != nil {
						os.Exit(0)
					}
					fmt.Println("New License Generated:", license)
					old = license
				} else {
					goto restart
				}
			}

		case "remove":
			fmt.Print("What licence would you like to remove?: ")
			scan := bufio.NewScanner(os.Stdin)
			scan.Scan()

			for _, table := range database {

				row := strings.Split(table, ":")

				if scan.Text() == row[3] { //Found in DB

					b, err := ioutil.ReadFile("/app/db/dbfile")
					if err != nil {
						os.Exit(0)
					}

					str := string(b)
					res := strings.Replace(str, table, "", -1)

					re := regexp.MustCompile("(?m)^\\s*$[\r\n]*")
					reres := strings.Trim(re.ReplaceAllString(res, ""), "\r\n")

					err = ioutil.WriteFile("/app/db/dbfile", []byte(reres), 0644)
					if err != nil {
						os.Exit(0)
					}
				}
			}

			fmt.Println("Done")
		case "exit":
			os.Exit(0)
			// default:
			// 	fmt.Println("Unknown Command")
		}
		time.Sleep(1 * time.Second)
	}
}
