package main

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"io/ioutil"
	"encoding/json"
	"regexp"
	"log"
	"flag"
)

type user struct {
	Username string `json:"username"`
	Password string `json:"password"` 
}

type config struct {
	User user `json:"user"`
}

const baseURL = "https://jira.alterdata.com.br/"
const authFile = "config.json"

func removeHTMLtags(str string) string {
	r := regexp.MustCompile("<br>|<p>|</p>")
	return r.ReplaceAllString(str,"")
}

func logIssue(issue jira.Issue) {
	fmt.Printf("Key: %s\n", issue.Key)

	if issue.Fields.Description != "" {
		fmt.Printf("Description: %s\n", issue.Fields.Description)
	}

	for _, comment := range issue.Fields.Comments.Comments {
		fmt.Printf("%s: %s\n", comment.Author.Name, removeHTMLtags(comment.Body))
	} 
}

func main() {	
	c := config{}	

	issueKey := flag.String("i", "", "issue")
	flag.StringVar(&c.User.Username, "u", "", "user")
	flag.StringVar(&c.User.Password, "p", "", "pwd")

	flag.Parse()		

	if c.User.Username == "" && c.User.Password == "" {
		bAuth, _ := ioutil.ReadFile(authFile)	
		json.Unmarshal(bAuth, &c)
	}

	if c.User.Username == "" {
		log.Fatal("ERR_NOUSR")
	}

	if c.User.Password == "" {
		log.Fatal("ERR_NOPWD")
	}

	tp := jira.BasicAuthTransport{
		Username: c.User.Username,
		Password: c.User.Password,
	}

	jiraClient, err := jira.NewClient(tp.Client(), baseURL)

	if err != nil {
		log.Fatalf("ERR_AUTH: %v", err)
	}

	jiraIssue, _, err := jiraClient.Issue.Get(*issueKey, nil)

	if err != nil {
		log.Fatal(err)
	}

	logIssue(*jiraIssue)
}