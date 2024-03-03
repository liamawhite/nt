package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/subosito/gozaru"
	"sigs.k8s.io/yaml"
)

type inputLater struct {
	Name    string
	Link    string
	Authors string
	Topics  string
	Content string
}

type Later struct {
	Type    string   `json:"type"`
	Name    string   `json:"name"`
	Link    string   `json:"link"`
	Authors []string `json:"authors"`
	Topics  []string `json:"topics"`
	Status  string   `json:"status"`
}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// This is a very quick and dirty implementation of the later command.
// Once I am happy with the workflow, I should build a better implementation
// which is more configurable and has fuzzy search for topics and authors.
var laterCmd = &cobra.Command{
	Use:     "later",
	Aliases: []string{"l"},
	RunE: func(cmd *cobra.Command, args []string) error {
		var input inputLater

		// Read the clipboard and check if it is a URL
		// If it is a URL, use it as the link and try to fetch the Title
		clipboardContent, _ := clipboard.ReadAll()
		if IsUrl(clipboardContent) {
			input.Link = clipboardContent

			// Fetch the Title
			res, err := http.Get(clipboardContent)
			if err == nil && res.StatusCode == 200 {
				defer res.Body.Close()
				doc, err := goquery.NewDocumentFromReader(res.Body)
				if err == nil {
					input.Name = doc.Find("title").Text()
				}
			}
		}

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Title("Name").Description("Name of the resource").Value(&input.Name),
				huh.NewInput().Title("Link").Description("Link to the resource").Value(&input.Link),
				huh.NewInput().Title("Topics").Description("Topics related to the resource (comma separated)").Value(&input.Topics),
				huh.NewInput().Title("Authors").Description("Authors of the resource (comma separated)").Value(&input.Authors),
				huh.NewText().Title("Notes").Description("Any notes you want to add").Value(&input.Content),
			),
		).WithProgramOptions(tea.WithAltScreen())

		if err := form.Run(); err != nil {
			return err
		}

		later := Later{
			Type:   "source",
			Name:   input.Name,
			Link:   input.Link,
			Status: "later",
		}

		topics := make([]string, 0)
		for _, t := range strings.Split(input.Topics, ",") {
			topics = append(topics, strings.TrimSpace(t))
		}
		later.Topics = topics

		authors := make([]string, 0)
		for _, a := range strings.Split(input.Authors, ",") {
			authors = append(authors, fmt.Sprintf("[[%v]]", strings.TrimSpace(a)))
		}
		later.Authors = authors

		buf := bytes.NewBuffer(nil)

		// Write the struct as front matter
		y, err := yaml.Marshal(later)
		if err != nil {
			return err
		}

		buf.WriteString("---\n")
		buf.Write(y)
		buf.WriteString("---\n")
		buf.WriteString(input.Content)

		// Write the buffer to a file
		// Sanitize the name to something that can be used as a file Name
		filename := gozaru.Sanitize(input.Name) + ".md"
		path := filepath.Join(laterDir, filename)
		if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
			return err
		}
		fmt.Println("Resource saved to", path)
		return nil
	},
}

var (
	laterDir string
)

func init() {
	rootCmd.AddCommand(laterCmd)

	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	laterCmd.Flags().StringVarP(&laterDir, "dir", "d", filepath.Join(home, "notes", "sources"), "Directory to store the resource")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// laterCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// laterCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
