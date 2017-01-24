package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var resultsLimit = 20
var mavenCentralUrl = "http://search.maven.org/solrsearch/select?q=%s&rows=%d&wt=json"
var MakeRequest = makeRequest
var out io.Writer = os.Stdout

type Spellcheck struct {
	Suggestions []interface{}
}

type Doc struct {
	Id            string
	A             string
	G             string
	LatestVersion string
}

type Response struct {
	NumFound int
	Docs     []Doc
}

type Dependencies struct {
	Response   Response
	Spellcheck Spellcheck
}

func parseDependencies(body []byte) (Dependencies, error) {
	var m Dependencies
	err := json.Unmarshal(body, &m)
	if err != nil {
		fmt.Fprint(out, "Parse error: ", err)
	}
	return m, err
}

func makeRequest(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}

func outputResultsFound(dependencies Dependencies) {
	fmt.Fprintln(out, "")
	resultsFound := fmt.Sprintf("Results found: %d", dependencies.Response.NumFound)
	fmt.Fprintln(out, resultsFound)
}

func outputSpellCheckSuggestions(dependencies Dependencies) {
	suggestions := dependencies.Spellcheck.Suggestions[1]
	fmt.Println(suggestions)
}

func outputGradleResults(dependencies Dependencies) {
	docs := dependencies.Response.Docs
	for i := range docs {
		line := fmt.Sprintf("%s:%s", docs[i].Id, docs[i].LatestVersion)
		fmt.Fprintln(out, line)
	}

	outputResultsFound(dependencies)
}

func outputMavenResults(dependencies Dependencies) {
	docs := dependencies.Response.Docs
	for i := range docs {
		groupId := fmt.Sprintf("<groupId>%s</groupId>", docs[i].G)
		artifactId := fmt.Sprintf("<artifactId>%s</artifactId>", docs[i].A)
		version := fmt.Sprintf("<version>%s</version>", docs[i].LatestVersion)

		fmt.Fprintln(out, "")
		fmt.Fprintln(out, "<dependency>")
		fmt.Fprintln(out, groupId)
		fmt.Fprintln(out, artifactId)
		fmt.Fprintln(out, version)
		fmt.Fprintln(out, "</dependency>")
	}

	outputResultsFound(dependencies)
}

func search(url string) Dependencies {
	body := MakeRequest(url)
	dependencies, err := parseDependencies([]byte(body))
	if err != nil {
		log.Fatal(err)
	}
	return dependencies
}

var repositoryTypes = []string{"gradle", "maven"}

func getFullUrl(searchText string) string {
	fullUrl := fmt.Sprintf(mavenCentralUrl, searchText, resultsLimit)
	return fullUrl
}

func GradleSearchAction(c *cli.Context) error {
	searchText := c.Args().Get(0)
	if searchText == "" {
		log.Fatal("\nERROR: Missing gradle search text\n")
		return nil
	}

	result := search(getFullUrl(searchText))
	outputGradleResults(result)

	return nil
}

func MavenSearchAction(c *cli.Context) error {
	searchText := c.Args().Get(0)
	if searchText == "" {
		log.Fatal("\nERROR: Missing maven search text\n")
		return nil
	}
	result := search(getFullUrl(searchText))
	outputMavenResults(result)

	return nil
}

func main() {
	app := cli.NewApp()
	app.Version = "1.0.0"
	app.EnableBashCompletion = true
	app.Name = "jarsearch"
	app.Usage = "commandline tool for retrieving build dependency information"
	app.Commands = []cli.Command{
		cli.Command{
			Aliases: []string{repositoryTypes[0]},
			Usage:   fmt.Sprintf("Example: jarsearch %s <search_text>", repositoryTypes[0]),
			Action:  GradleSearchAction,
		},
		cli.Command{
			Aliases: []string{repositoryTypes[1]},
			Usage:   fmt.Sprintf("Example: jarsearch %s <search_text>", repositoryTypes[1]),
			Action:  MavenSearchAction,
		},
	}
	app.CommandNotFound = func(c *cli.Context, command string) {
		fmt.Fprintf(
			c.App.Writer,
			"ERROR: Repository type %q not supported. Must be one of: %s \n",
			command,
			repositoryTypes)
	}

	app.Run(os.Args)
}
