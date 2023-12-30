package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"
	"unsafe"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"

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
		Id:          time.Now().String(),
		Description: res,
	})

	rss, err := feed.ToRss()
	if err != nil {
		log.Fatal(err)
	}

	return rss
}

func generateHTML(data []Doc) string {
	questions := ""
	for _, item := range data {
		questions += fmt.Sprintf("# %s \n\n", item.Docs)
		for _, question := range item.Qs {
			questions += fmt.Sprintf("- %s \n", question)
		}
	}
	return Md2HTML(questions)
}

func Md2HTML(md string) string {
	if md == "" {
		return ""
	}
	var buf bytes.Buffer
	markdown := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithExtensions(
			extension.GFM,
			extension.Strikethrough,
			extension.TaskList,
			extension.Linkify,
			extension.Table,
			extension.DefinitionList,
			extension.Footnote,
			extension.Typographer,
			extension.NewTypographer(
				extension.WithTypographicSubstitutions(extension.TypographicSubstitutions{
					extension.LeftSingleQuote:  []byte("&sbquo;"),
					extension.RightSingleQuote: nil, // nil disables a substitution
				}),
			),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
	if err := markdown.Convert(StringToBytes(md), &buf); err != nil {
		return ""
	}

	return buf.String()
}

func StringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
