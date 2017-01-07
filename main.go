package main

import (
	"bufio"
	"log"
	"os"
	"time"

	"github.com/sogko/go-wordpress"
	"net/http"
	"path/filepath"
	"strings"
	"flag"
)

type LJPost struct {
	Date     time.Time
	Subject  string
	Picture  string
	ItemID   uint
	Security string
	Friends  string
}

var username, password, baseURL *string

func main() {
	dirPath := flag.String("d", "", "Directory with LJ archive")
	username = flag.String("username", "", "Wordpress blog user name")
	password = flag.String("password", "", "Wordpress blog password")
	baseURL = flag.String("url", "", "Wordpress API URL")


	err := filepath.Walk(*dirPath, readFromFile)
	if err != nil {
		log.Println(err)
	}

}

func readFromFile(path string, fileInfo os.FileInfo, err error) error {
	if !fileInfo.IsDir() {
		file, err := os.Open(path)
		if err != nil {
			log.Printf("Error in opening file %s for reading. %s", path, err)
			return err
		}
		defer file.Close()

		var lines []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		postToWP(lines)
		// Processing into LJPost type variable
		return scanner.Err()
	}
	return nil
}

func postToWP(fileContents []string) error {

	var p wordpress.Post
	p.Status = wordpress.PostStatusPublish
	// fmt.Println(fileContents)

	//Обработочка файлов
	var text []string
	for _, cStr := range fileContents {
		log.Printf("Processing string \"%s\"", cStr)
		if strings.HasPrefix(cStr, "Date:") {

			//Обработка Date
			layoutStr := "Date:      2006-01-02 15:04"
			dTime, err := time.Parse(layoutStr, cStr)
			if err != nil {
				dTime = time.Now()
			}
			// Layout:  Wed, 29 Jun 2005 20:01 +00
			layoutString := "2006-01-02 15:04:05"
			p.Date = dTime.Format(layoutString)
			// p.Date = dTime.Format(time.RFC1123Z)

		} else if strings.HasPrefix(cStr, "Subject:") {
			p.Title.Raw = strings.Split(cStr, "Subject:   ")[1]
		} else if strings.HasPrefix(cStr, "Mood:") {

		} else if strings.HasPrefix(cStr, "Music:") {

		} else if strings.HasPrefix(cStr, "ItemID:") {

		} else if strings.HasPrefix(cStr, "Tags:") {

		} else if strings.HasPrefix(cStr, "Picture:") {

		} else if strings.HasPrefix(cStr, "Security:") {
			p.Status = wordpress.PostStatusPrivate
			p.Password = "12321"


		} else if strings.HasPrefix(cStr, "Friends:") {

		} else {
			// Это просто текст
			text = append(text, cStr)
		}

	}
	if text[0] == "" {
		text = text[1:]
	}

	p.Excerpt.Raw = p.Title.Raw
	p.CommentStatus = wordpress.CommentStatusOpen
	p.Format = wordpress.PostFormatStandard
	p.Type =   wordpress.PostTypePost
	p.Author = 1

	p.Content.Raw = strings.Join(text, "\n")
	log.Println("\nDate: ", p.Date, "\nSubject: ", p.Title.Raw, "\nContent: ", p.Content.Raw)

	// create wp-api client
	//log.Println("Creating client")
	client := wordpress.NewClient(&wordpress.Options{
		BaseAPIURL: *baseURL, // example: `http://192.168.99.100:32777/wp-json/wp/v2`
		Username:   *username,
		Password:   *password,
	})
	//
	//// ===============
	//
	//log.Println("Posting")
	newPost, resp, body, err := client.Posts().Create(&p)
	//log.Println("Posting finished")
	if err != nil {
		log.Println("Error: ", err)
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		log.Println("Received Status code: %v", resp.StatusCode)
	}
	if body == nil {
		log.Println("body should not be nil")
	}
	if newPost == nil {
		log.Println("newPost should not be nil")
	}
	log.Println("Post posted\n",newPost, "\n=========\n")
	return nil
}
