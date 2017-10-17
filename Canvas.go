package main

import (
    "fmt"
    "log"
    "sort"
    "time"
    "strings"
    "net/url"
)

//Read File for Secret and Key
func ReadFile() (key, secret string){
    //path := "Special.txt"
    key = "Hello"
    secret = "world"
    return key, secret
}

//Create Hashed Message
func SignedMessage(address, secret string) string{
    //timestamp := strings.Replace(time.Now().UTC().Format(time.RFC1123), "UTC", "GMT", -1)
    //Testing purposes
    //timestamp := time.Date(2015, time.December, 1, 9, 24, 50, 0, time.UTC)
    u, err := url.Parse(address)
    if err != nil{
        log.Fatal(err)
    }
    address = SortParams(u.RawQuery)
    fmt.Println(address)
    return "hello"
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

func main(){
    key, secret := ReadFile()
    address :=  "https://portal.inshosteddata.com/api/account/self/dump?limit=100&after=45"
    timestamp := strings.Replace(time.Now().UTC().Format(time.RFC1123), "UTC", "GMT", -1)
    msg := SignedMessage(address, secret)
    fmt.Printf("Key: %s \n Secret: %s\n Time: %s\n URL: %s\n Message: %s\n", key, secret, timestamp, address, msg)
}
