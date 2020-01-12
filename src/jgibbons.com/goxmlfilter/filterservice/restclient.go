package filterservice

import (
	"fmt"
	"net/http"
	"os"
)

func callRestApi(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	basicAuth(req)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
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
