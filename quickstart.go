package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

const credentials_file = "credentials.json"
const token_file = "token.json"

func main() {
	service := service_wrapper(os.Args[0])
	if len(os.Args) < 2 {
		list_labels(service)
	} else {
		c := configuration(os.Args, "from")
		message := message(c)
		_, err := service.Users.Messages.Send("me", message).Do()
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			log.Printf("Message sent to %v", c.receiver)
		}
	}
}

func list_labels(srv *gmail.Service) {
	user := "me"
	r, err := srv.Users.Labels.List(user).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}
	if len(r.Labels) == 0 {
		fmt.Println("No labels found.")
		return
	}
	fmt.Println("Labels:")
	for _, l := range r.Labels {
		fmt.Printf("- %s\n", l.Name)
	}
}

func message(c *config) *gmail.Message {
	header := make(map[string]string)
	header["From"] = c.from
	header["To"] = c.receiver
	header["Subject"] = c.subject
	header["Content-Type"] = "text/plain"
	//	Does not work reliably.
	//	Sometimes received content is only a few, unprintable, characters.
	//	header["Content-Transfer-Encoding"] = "base64"
	var message string
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + message_content(c.infile)
	return &gmail.Message{Raw: base64.URLEncoding.EncodeToString([]byte(message))}
}

func message_content(infile string) string {
	bytes, err := os.ReadFile(infile)
	if err != nil {
		log.Fatalf("Unable to retrieve content: %v", err)
	}
	return string(bytes)
}

func service_wrapper(program string) *gmail.Service {
	directory := path.Dir(program)
	credentials := service_file("credentials", path.Join(directory, credentials_file))
	token := service_file("token", path.Join(directory, token_file))
	return service(credentials, token)
}

func service_file(name, fallback string) (result string) {
	result = os.Getenv(name)
	if result == "" {
		result = fallback
	}
	return result
}

func service(credentials, token string) *gmail.Service {
	ctx := context.Background()
	b, err := ioutil.ReadFile(credentials)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope, gmail.GmailComposeScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config, token)
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}
	return srv
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config, tokFile string) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	log.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
