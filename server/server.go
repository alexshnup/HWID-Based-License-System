package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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
	PORT   int    = 9347
	salt   string = "12345salt"
	token  string = "1234TOKEN"
	dbfile string = "/app/db/dbfile"
)

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
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

func sha256Hash(text string) string {
	hasher := sha256.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// LICENSE:EXPDATE:EMAIL:HWID = 0,1,2
func checkHandler(response http.ResponseWriter, request *http.Request) {

	request.ParseForm()
	license := request.FormValue("license")
	hwid := request.FormValue("hwid")

	database, _ := readLines(dbfile)

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

				out := sha256Hash("0" + salt + hwid)

				// out := md5Hash(md5Hash("0") + md5Hash("solt2022"+stlsolt))
				fmt.Fprint(response, out) //Registed, Good licnese
				return
			} else if row[3] == "NOTSET" {
				b, err := os.ReadFile(dbfile)
				if err != nil {
					fmt.Println("READfromCHECK")
					os.Exit(0)
				}

				str := string(b)
				edit := row[0] + ":" + row[1] + ":" + row[2] + ":" + hwid
				res := strings.Replace(str, table, edit, -1)

				err = os.WriteFile(dbfile, []byte(res), 0644)
				if err != nil {
					fmt.Println("WRITEfromCHECK")
					os.Exit(0)
				}

				log.Printf("%v from %v NEW Registed, Good licnese %s", time.Now(), request.RemoteAddr, license)
				// fmt.Fprintf(response, "2") //Registed, Good licnese
				out := sha256Hash("2" + salt + hwid)

				fmt.Fprint(response, out) //Registed, Good licnese
				return

			}
		} else if license == row[0] && !t.After(t2) {

			log.Printf("%v from %v registerd but license %s", time.Now(), request.RemoteAddr, license)
			// fmt.Fprintf(response, "1") //registerd but license experied

			out := sha256Hash("1" + salt + hwid)
			fmt.Fprint(response, out) //Registed, Good licnese
			return
		}

		if license == row[0] && hwid != row[3] {
			fmt.Printf("from %v license %s problem hwid", request.RemoteAddr, license)
		}
	}
}

// LicenseRequest is a struct that represents the JSON payload from the client for adding a license.
type LicenseRequest struct {
	Email      string `json:"email"`
	Expiration string `json:"expiration"`
}
type ResetRequest struct {
	Key string `json:"key"`
}
type RemoveRequest struct {
	Email string `json:"email"`
}

func addKeyHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON request
	var request LicenseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate email
	if request.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Validate and parse expiration
	_, err := time.Parse("2006-01-02", request.Expiration)
	if err != nil {
		http.Error(w, "Expiration must be in the YYYY-MM-DD format", http.StatusBadRequest)
		return
	}

	// // Read the dbfile
	b, _ := os.ReadFile(dbfile)
	// if err != nil {
	// 	http.Error(w, "Failed to read data: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// Generate the license
	license := randomString(4) + "-" + randomString(4) + "-" + randomString(4)
	for LicExist(license) {
		license = randomString(4) + "-" + randomString(4) + "-" + randomString(4)
	}

	// Append the new license to the existing data and cleanup
	str := string(b)
	str = str + "\r\n" + license + ":" + request.Expiration + ":" + request.Email + ":NOTSET"

	re := regexp.MustCompile("(?m)^\\s*$[\r\n]*")
	str2 := strings.Trim(re.ReplaceAllString(str, ""), "\r\n")

	// Write back to the dbfile
	err = os.WriteFile(dbfile, []byte(str2), 0644)
	if err != nil {
		http.Error(w, "Failed to write data", http.StatusInternalServerError)
		return
	}

	// Respond to the client
	response := map[string]string{
		"message":  "New license generated",
		"license":  license,
		"email":    request.Email,
		"exp_date": request.Expiration,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func listKeysHandler(w http.ResponseWriter, r *http.Request) {
	// Open the dbfile
	file, err := os.Open(dbfile)
	if err != nil {
		http.Error(w, "Failed to read data", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Read the lines from the file
	var keys []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		keys = append(keys, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		http.Error(w, "Error reading data", http.StatusInternalServerError)
		return
	}

	// Respond to the client with the list of keys
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(keys)
}

func resetKeyHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON request
	var req ResetRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	fmt.Printf("Resetting key: %s\n", req.Key)

	// Open the dbfile
	file, err := os.Open(dbfile)
	if err != nil {
		http.Error(w, "Failed to read data", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Read the lines from the file and try to find and reset the key
	var lines []string
	var found bool
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) > 0 && parts[0] == req.Key {
			found = true
			parts[len(parts)-1] = "NOTSET"  // Reset the last column
			line = strings.Join(parts, ":") // Recombine the line
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		http.Error(w, "Error reading data", http.StatusInternalServerError)
		return
	}

	if !found {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	// Write the modified data back to the file
	err = os.WriteFile(dbfile, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		http.Error(w, "Failed to write data", http.StatusInternalServerError)
		return
	}

	// Respond to the client with success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func removeKeyHandler(w http.ResponseWriter, r *http.Request) {
	var req RemoveRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Read the dbfile
	b, err := ioutil.ReadFile(dbfile)
	if err != nil {
		http.Error(w, "Failed to read data", http.StatusInternalServerError)
		return
	}

	str := string(b)
	var removed bool

	// Go through each line in the database, keeping those that don't match
	// the provided email.
	var newDB []string
	lines := strings.Split(str, "\n")
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) > 3 && parts[2] == req.Email {
			removed = true
			continue // skip this line
		}
		newDB = append(newDB, line)
	}

	if !removed {
		http.Error(w, "Email not found", http.StatusNotFound)
		return
	}

	// Clean up any empty lines
	re := regexp.MustCompile("(?m)^\\s*$[\r\n]*")
	newStr := strings.Trim(re.ReplaceAllString(strings.Join(newDB, "\n"), ""), "\r\n")

	// Write the modified data back to the file
	err = os.WriteFile(dbfile, []byte(newStr), 0644)
	if err != nil {
		http.Error(w, "Failed to write data", http.StatusInternalServerError)
		return
	}

	// Respond to the client with success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func validateTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip token validation for open routes
		if r.RequestURI == "/" && r.Method == "POST" {
			next.ServeHTTP(w, r)
			return
		}
		// Get the token from the Authorization header
		t := r.Header.Get("Authorization")

		// Validate the token
		if t != token {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Call the next handler if token is valid
		next.ServeHTTP(w, r)
	})
}

func serverAPI() {
	router := mux.NewRouter()

	// Apply the middleware
	router.Use(validateTokenMiddleware)

	router.HandleFunc("/", checkHandler).Methods("POST")
	router.HandleFunc("/add", addKeyHandler).Methods("POST")
	router.HandleFunc("/list", listKeysHandler).Methods("GET")
	router.HandleFunc("/reset-key", resetKeyHandler).Methods("POST")
	router.HandleFunc("/remove", removeKeyHandler).Methods("DELETE")
	http.Handle("/", router)

	http.ListenAndServe(":"+string(strconv.Itoa(PORT)), nil)
}

func GetNewLicNumber() string {
	return randomString(4) + "-" + randomString(4) + "-" + randomString(4)
}

func LicExist(toFind string) bool {
	database, _ := readLines(dbfile)
	for _, v := range database {
		if strings.Contains(v, toFind) {
			return true
		}
	}
	return false
}

func getEnvVar(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	} else {
		fmt.Printf("%s not set, will use default %s\n", key, defaultValue)
		return defaultValue
	}
}

func main() {
	salt = getEnvVar("SALT", salt)
	token = getEnvVar("TOKEN", token)
	fmt.Println("Using salt:", salt)
	fmt.Println("Using token:", token)

	fmt.Println("License Server")
	fmt.Println("Github: https://github.com/alexshnup/easy-license-system")
	fmt.Println("Forked from: https://github.com/SaturnsVoid/HWID-Based-License-System")

	if !fileExists(dbfile) {
		log.Println("Database does not exist, creating new database.")
		if err := createFile(dbfile); err != nil {
			log.Fatalf("Failed to create file: %v", err)
		}
	}

	database, err := readLines(dbfile)
	if err != nil {
		log.Fatalf("Error reading database: %v", err)
	}

	fmt.Println("Total Licenses:", len(database))

	serverAPI()
}
