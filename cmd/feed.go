package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/gorilla/feeds"
	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
)

type Doc struct {
	Docs string   `yaml:"docs"`
	Qs   []string `yaml:"qs"`
}

// feedCmd represents the feed command
var feedCmd = &cobra.Command{
	Use:   "feed",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		outputFile := args[1]
		yamlData, err := os.ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}

		var data []Doc
		if err := yaml.Unmarshal(yamlData, &data); err != nil {
			log.Fatal(err)
		}

		rssFeed := generateRSS(data)
		if err := os.WriteFile(outputFile, []byte(rssFeed), 0o644); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Successfully wrote RSS feed to file")
	},
}

func init() {
	rootCmd.AddCommand(feedCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// feedCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// feedCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func generateRSS(data []Doc) string {
	feed := &feeds.Feed{
		Title:       "docs",
		Link:        &feeds.Link{Href: "https://github.com/hxhacking/docs"},
		Description: "just for fun",
	}

	res := generateHTML(data)
	feed.Items = append(feed.Items, &feeds.Item{
		Title:       "docs",
		Link:        &feeds.Link{Href: "https://github.com/hxhacking/docs"},
		Description: res,
	})

	rss, err := feed.ToRss()
	if err != nil {
		log.Fatal(err)
	}

	return rss
}

func generateHTML(data []Doc) string {
	// Generate HTML markup using the data
	questions := ""

	for _, item := range data {
		questions += fmt.Sprintf("<h2>%s</h2>", item.Docs)
		questions += "<ul>"
		for _, question := range item.Qs {
			// questions += "<!--<li><input disabled=\"\" type=\"checkbox\">" + question + "</li>-->"
			questions += fmt.Sprintf("<li>%s</li>", question)
		}
		questions += "</ul>"
	}

	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		  <head>
		  </head>
		  <body>
		        %s
		  </body>
		</html>
	  `, questions)
}
