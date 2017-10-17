package main

import (
    "bufio"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "fmt"
    "github.com/tidwall/gjson"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "os"
    "sort"
    "strings"
    "time"
)

//Define File Message Structure
type FileMessage struct{
    Filename    string  `json:"filename"`
    Table       string  `json:"table"`
    FileURL     string  `json:"url"`
    Partial     bool    `json:"partial"`
}

//Define Response Message Structure
type Message struct{
    SchemaVersion   string          `json:"schemaVersion"`
    Incomplete      bool            `json:"incomplete"`
    Files           []FileMessage   `json:"files"`
}

//Read File for Secret and Key
func ReadFile() (key, secret string){
    path := "Special.txt"
    file, err := os.Open(path)
    if err !=nil {
        log.Fatal(err)
    }
    defer file.Close()
    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan(){
        lines = append(lines, scanner.Text())
    }

    key = lines[0]
    secret = lines[1]
    return key, secret
}

//Generate message for hashing
func GenerateMessage(secret, timestamp, canvasURL string) string{
    u, err := url.Parse(canvasURL)
    if err != nil{
        log.Fatal(err)
    }
    queryStr := SortParams(u.RawQuery)
    msgSlice := []string {"GET", u.Host, "", "", u.Path, queryStr, timestamp, secret}
    msg := strings.Join(msgSlice, "\n")

    return msg
}

//HMAC hash creation
func ComputeHash(timestamp, canvasURL, secret string) string{
    msg := GenerateMessage(secret, timestamp, canvasURL)
    byteSecret := []byte(secret)
    hash := hmac.New(sha256.New, byteSecret)
    hash.Write([]byte(msg))
    return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

//Sort URL query parameters alphabetically
func SortParams(rawQuery string) string{
    sQuery, _ := url.ParseQuery(rawQuery)
    sQueryK := make([]string, len(sQuery))
    i := 0
    for k, _ := range sQuery{
        sQueryK[i] = k
        i++
    }
    sort.Strings(sQueryK)
    rawQuery = ""
    for x:=0; x<len(sQueryK); x++{
        rawQuery = rawQuery + sQueryK[x] + "=" + sQuery[sQueryK[x]][0]
        if x<(len(sQueryK)-1){
            rawQuery = rawQuery + "&"
        }
    }
    return rawQuery
}

//Generate Headers and parse response --Current WIP--
func httpRequest(msg, key, canvasURL, timestamp string){

    client := &http.Client{}
    req, err := http.NewRequest("GET", canvasURL, nil)
    req.Header.Add("Authorization", "HMACAuth " + key + ":" + msg)
    req.Header.Add("Date", timestamp)
    resp, err := client.Do(req)
    if err !=nil {
        fmt.Println("ERROR", err)
    }
    if resp !=nil{
        defer resp.Body.Close()
    }

    //body response returned in bytes
    body, readErr := ioutil.ReadAll(resp.Body)
    if readErr !=nil{
        fmt.Println("body error:", readErr)
    }

    var respMsg Message
    gjson.Unmarshal(body, &respMsg)
    fmt.Printf("Schema Version: %s, Incomplete: %t \n", respMsg.SchemaVersion, respMsg.Incomplete)
    fmt.Println("File names:")
    for i:=0; i<5; i++{
        fmt.Println(respMsg.Files[i].Filename)
    }
}

func main(){
    key, secret := ReadFile()
    canvasURL :=  "https://portal.inshosteddata.com/api/account/self/file/sync"
    //Timestamp, replace UTC with GMT and convert to string in standard RFC1123 format
    timestamp := strings.Replace(time.Now().UTC().Format(time.RFC1123), "UTC", "GMT", -1)
    msg := ComputeHash(timestamp, canvasURL, secret)
    httpRequest(msg, key, canvasURL, timestamp)
}
