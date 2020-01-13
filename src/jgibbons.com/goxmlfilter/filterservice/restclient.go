package filterservice

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
)

func callRestApi(url string, ignoreCertAuthority bool) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	basicAuth(req)

	resp, err := createHttpClient(ignoreCertAuthority).Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// If using self certified certs, then maybe ignore the authorizer, unless
// you have your own certs all installed properly....
func createHttpClient(ignoreCertAuthority bool) *http.Client {
	var client *http.Client
	if ignoreCertAuthority {
		log.Println("Ignoring HTTPS certificate authority checks.")
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	} else {
		client = http.DefaultClient
	}
	return client
}

func basicAuth(req *http.Request) {
	if userName, exists := os.LookupEnv("xmlfilt_username"); exists {
		if password, pexists := os.LookupEnv("xmlfilt_password"); pexists {
			fmt.Printf("Setting basic auth for user [%s]\n", userName)
			req.SetBasicAuth(userName, password)
		} else {
			fmt.Println("No environment variable 'xmlfilt_password' found, basic auth not applied")
		}
	} else {
		fmt.Println("No environment variable 'xmlfilt_username' found, basic auth not applied")
	}
}
