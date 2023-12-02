package tWXYY

import (
	"fmt"
	"net/http"
	"time"
)

func pingGoogle(proxyStr string) error {
	fmt.Println("Google is ")
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Get("https://www.google.com")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Google is reachable")
	} else {
		fmt.Println("Google is not reachable, Status Code:", resp.StatusCode)
	}
	//fmt.Println(resp)
	return nil
}
