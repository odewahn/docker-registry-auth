package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/parnurzeal/gorequest"
)

// This function parses the Www-Authenticate header provided in the challenge
// It has the following format
// Bearer realm="https://gitlab.com/jwt/auth",service="container_registry",scope="repository:andrew18/container-test:pull"
func parseBearer(bearer []string) map[string]string {
	out := make(map[string]string)
	for _, b := range bearer {
		for _, s := range strings.Split(b, " ") {
			if s == "Bearer" {
				continue
			}
			for _, params := range strings.Split(s, ",") {
				fields := strings.Split(params, "=")
				key := fields[0]
				val := strings.Replace(fields[1], "\"", "", -1)
				out[key] = val
			}
		}
	}
	return out
}

func main() {

	// Based on
	// http://www.cakesolutions.net/teamblogs/docker-registry-api-calls-as-an-authenticated-user

	request := gorequest.New()

	//url := "https://index.docker.io/v2/odewahn/myalpine/tags/list"
	url := "https://registry.gitlab.com/v2/andrew18/container-test/tags/list"

	// First step is to get the endpoint where we'll be authenticating
	resp, _, _ := request.Get(url).End()

	// This has the various things we'll need to parse and use in the request
	params := parseBearer(resp.Header["Www-Authenticate"])
	paramsJSON, _ := json.Marshal(&params)
	log.Println(string(paramsJSON))

	// Get the token
	challenge := gorequest.New()
	resp, body, _ := challenge.Get(params["realm"]).
		SetBasicAuth(os.Getenv("USER"), os.Getenv("PWD")).
		Query(string(paramsJSON)).
		End()

	token := make(map[string]string)
	json.Unmarshal([]byte(body), &token)

	// Now reissue the challenge with the toekn in the Header
	// curl -IL https://index.docker.io/v2/odewahn/image/tags/list
	authenticatedRequest := gorequest.New()

	resp, body, _ = authenticatedRequest.Get(url).
		Set("Authorization", "Bearer "+token["token"]).
		End()

	fmt.Println(body)

}
