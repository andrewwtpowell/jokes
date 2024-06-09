// Query joke API for a quick joke
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const jokeUrl = "https://v2.jokeapi.dev/"

type categoryResult struct {
    Error           bool
    Categories      []string
    CategoryAliases []catAlias
    Timestamp       int
}

type catAlias struct {
    Alias       string
    Resolved    string
}

type singleResult struct {
    Error       bool
    Category    string
    Type        string
    Joke        string
    Flags       flags
    Id          int
    Safe        bool
    Lang        string
}

type twopartResult struct {
    Error       bool
    Category    string
    Type        string
    Setup       string
    Delivery    string
    Flags       flags
    Id          int
    Safe        bool
    Lang        string
}

type flags struct {
    Nsfw        bool
    Religious   bool
    Political   bool
    Racist      bool
    Sexist      bool
    Explicit    bool 
}

func main() {

    //fmt.Printf("Checking server availability...\n")
    if _, err := getEndpoint("ping"); err != nil {
        log.Fatal(err)
    }

    categories, err := getAvailCategories()
    if err != nil {
        log.Fatal(err)
    }

    catPtr := flag.String("cat", "Any", strings.Join(categories, " "))
    typePtr := flag.String("type", "single", "joke type (single or twopart)")

    flag.Parse()

    found := false
    for _, cat := range categories {
        
        if *catPtr == cat {
            found = true
            break
        }
    }

    if !found {
        log.Fatal("invalid category entered")
    }

    if *typePtr == "single" {
        joke, err := getSingleJoke(*catPtr)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("%s\n", joke)
    } else if *typePtr == "twopart" {
        setup, delivery, err := getTwopartJoke(*catPtr)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("%s\n", setup)
        time.Sleep(2)
        fmt.Printf("%s\n", delivery)
    } else {
        log.Fatal("invalid type entered")
    }
}

func getTwopartJoke(category string) (string, string, error) {
    
    resp, err := getEndpoint("joke/" + category + "?type=twopart")
    if err != nil {
        panic(err)
    }

    defer resp.Body.Close()

    var result twopartResult
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        panic(err)
    }

    if result.Error == true {
        return "", "", fmt.Errorf("received error response")
    }

    if result.Setup == "" {
        return "", "", fmt.Errorf("empty joke field received")
    }

    return result.Setup, result.Delivery, nil
}

func getSingleJoke(category string) (string, error) {
    
    resp, err := getEndpoint("joke/" + category + "?type=single")
    if err != nil {
        panic(err)
    }

    defer resp.Body.Close()

    var result singleResult
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        panic(err)
    }

    if result.Error == true {
        return "", fmt.Errorf("received error response")
    }

    if result.Joke == "" {
        return "", fmt.Errorf("empty joke field received")
    }

    return result.Joke, nil
}

func getAvailCategories() ([]string, error) {

    resp, err := getEndpoint("categories")
    if err != nil {
        panic(err)
    }

    defer resp.Body.Close()

    var result categoryResult
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        panic(err)
    }

    return result.Categories, nil
}

func getEndpoint(endpoint string) (*http.Response, error) {

    dest := jokeUrl + endpoint 
    //fmt.Printf("Sending request to %s\n", dest)
    resp, err := http.Get(dest)
    if err != nil {
        panic(err)
    }

    if resp.StatusCode != http.StatusOK {
        return resp, fmt.Errorf("accessing endpoint %s failed: %s\n", endpoint, resp.Status)
    }

    return resp, nil
}
