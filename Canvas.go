package main

import (
    "fmt"
    "log"
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
func SignedMessage(URL, secret string) string{
    //timestamp := strings.Replace(time.Now().UTC().Format(time.RFC1123), "UTC", "GMT", -1)
    //Testing purposes
    //timestamp := time.Date(2015, time.December, 1, 9, 24, 50, 0, time.UTC)
    URL = SortParams(URL)
    fmt.Println(URL)
    return "hello"
}

//Sort URL query parameters
func SortParams(URL string) string{
    u, err := url.Parse(URL)
    if err != nil{
        log.Fatal(err)
    }
    sQuery, _ := url.ParseQuery(u.RawQuery)
    fmt.Println(sQuery)
    fmt.Println(sQuery["limit"][0])
    return URL
}

func main(){
    key, secret := ReadFile()
    URL :=  "https://portal.inshosteddata.com/api/account/self/dump?limit=100&after=45"
    timestamp := strings.Replace(time.Now().UTC().Format(time.RFC1123), "UTC", "GMT", -1)
    msg := SignedMessage(URL, secret)
    fmt.Printf("Key: %s \n Secret: %s\n Time: %s\n URL: %s\n Message: %s\n", key, secret, timestamp, URL, msg)
}
