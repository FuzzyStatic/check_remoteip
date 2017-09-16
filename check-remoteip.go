/**
 * @Author: Allen Flickinger <FuzzyStatic>
 * @Date:   2017-09-12T19:13:12-04:00
 * @Email:  allen.flickinger@gmail.com
 * @Last modified by:   FuzzyStatic
 * @Last modified time: 2017-09-15T22:09:40-04:00
 */

package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gmail "google.golang.org/api/gmail/v1"
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

	if newIPStr, err = getRemoteIP(); err != nil {
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
			send([]byte(
				"From: " + conf.From + "\r\n" +
					"To: " + conf.To + "\r\n" +
					"Subject: New Remote IP: " + newIPStr + "\r\n\r\n" +
					"Remote IP changed from " + prevIPStr + " to " + newIPStr))
		}
	}
}

// GetRemoteIP gets the remote IP as reported by ipecho.net
func getRemoteIP() (string, error) {
	var (
		resp *http.Response
		err  error
		body []byte
	)

	if resp, err = http.Get("http://ipecho.net/plain"); err != nil {
		return "", err
	}

	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return "", err
	}

	resp.Body.Close()
	return string(body), err
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("check_remoteip.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// Send sends email
func send(msgStr []byte) {
	var (
		ctx    context.Context
		secret []byte
		conf   *oauth2.Config
		client *http.Client
		err    error
		gs     *gmail.Service
		msg    gmail.Message
	)

	ctx = context.Background()

	// Reads in our credentials
	if secret, err = ioutil.ReadFile("client_secret.json"); err != nil {
		log.Printf("Error: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/check_remoteip.json
	if conf, err = google.ConfigFromJSON(secret, gmail.GmailSendScope); err != nil {
		log.Printf("Error: %v", err)
	}

	client = getClient(ctx, conf)

	// Create a new gmail service using the client
	if gs, err = gmail.New(client); err != nil {
		log.Printf("Error: %v", err)
	}

	// Place messageStr into message.Raw in base64 encoded format
	msg.Raw = base64.URLEncoding.EncodeToString(msgStr)

	// Send the message
	if _, err = gs.Users.Messages.Send("me", &msg).Do(); err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Println("Message sent!")
	}
}
