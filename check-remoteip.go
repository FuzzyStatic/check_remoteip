/**
 * @Author: Allen Flickinger <FuzzyStatic>
 * @Date:   2017-09-12T19:13:12-04:00
 * @Email:  allen.flickinger@gmail.com
 * @Last modified by:   FuzzyStatic
 * @Last modified time: 2017-09-15T21:35:30-04:00
 */

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"check_remoteip/addr"
	"check_remoteip/gmail"
)

const confFilePath = "./config.json"
const ipFilePath = "REMOTE_IP"
const ipFileMode = 0664

type confJSON struct {
	From string `json:"From"`
	To   string `json:"To"`
}

func main() {
	var (
		conf      confJSON
		confFile  []byte
		ipFile    *os.File
		prevIP    []byte
		prevIPStr string
		newIPStr  string
		err       error
	)

	if confFile, err = ioutil.ReadFile(confFilePath); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}

	json.Unmarshal(confFile, &conf)

	if _, err = os.Stat(ipFilePath); os.IsNotExist(err) {
		if ipFile, err = os.Create(ipFilePath); err != nil {
			log.Printf("Error: %v", err)
			os.Exit(1)
		}
		defer ipFile.Close()
	}

	if prevIP, err = ioutil.ReadFile(ipFilePath); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}

	if newIPStr, err = addr.GetRemoteIP(); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}

	prevIPStr = strings.TrimSpace(string(prevIP))

	if prevIPStr != newIPStr {
		log.Printf("Previous Remote IP: %s", string(prevIP))
		log.Printf("New Remote IP: %s", newIPStr)

		// Write Remote IP to file
		ioutil.WriteFile(ipFilePath, []byte(newIPStr), 0644)

		if string(prevIP) != "" {
			gmail.Send([]byte(
				"From: " + conf.From + "\r\n" +
					"To: " + conf.To + "\r\n" +
					"Subject: New Remote IP: " + newIPStr + "\r\n\r\n" +
					"Remote IP changed from " + prevIPStr + " to " + newIPStr))
		}
	}
}
