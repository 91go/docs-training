package cmd

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"

	"github.com/gorilla/feeds"
	"github.com/spf13/cobra"
)

type qs struct {
	FileName string
	Doc      []doc
}

type doc struct {
	Docs string   `yaml:"docs"`
	Qs   []string `yaml:"qs"`
}

// feedCmd represents the feed command
var feedCmd = &cobra.Command{
	Use:   "feed",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		folderName := args[0]
		outputFile := args[1]

		type fi struct {
			FilePath string
			FileName string
		}

		fileName := make(chan fi)

		_, err := os.Stat(folderName)
		if err != nil {
			panic(err)
		}

		go func() {
			err = filepath.WalkDir(folderName, func(fp string, d fs.DirEntry, err error) error {
				if !d.IsDir() && (strings.HasSuffix(fp, "yml") || strings.HasSuffix(fp, "yaml")) {
					fileName <- fi{FilePath: fp, FileName: d.Name()}
				}
				return nil
			})
			close(fileName)
		}()

		// rssFeed := make(chan string)
		rssFeed := make(chan qs)

		go func(fileName chan fi) {
			var wg sync.WaitGroup
			for f := range fileName {
				wg.Add(1)
				go func(f fi) {
					defer wg.Done()

					yamlData, err := os.ReadFile(f.FilePath)
					if err != nil {
						log.Fatal(err)
					}

					var data []doc
					if err := yaml.Unmarshal(yamlData, &data); err != nil {
						log.Fatal(err)
					}
					rssFeed <- qs{FileName: f.FileName, Doc: data}
					// select {
					// case rssFeed <- qs{FileName: f.FileName, Doc: data}:
					// default:
					// 	// Handle the case when the rssFeed channel is blocked
					// 	log.Println("rssFeed channel is blocked. Skipping data.")
					// }
				}(f)
			}
			wg.Wait()
			close(rssFeed)
		}(fileName)

		// var res string
		var resBuilder strings.Builder
		for rf := range rssFeed {
			// res += generateHTML(rf)
			resBuilder.WriteString(generateHTML(rf))
		}

		// for {
		// 	select {
		// 	case rf, ok := <-rssFeed:
		// 		if !ok {
		// 			// rssFeed channel has been closed
		// 			rssFeed = nil
		// 			break
		// 		}
		// 		resBuilder.WriteString(generateHTML(rf))
		// 	default:
		// 		// Handle the case when the rssFeed channel is blocked
		// 		log.Println("rssFeed channel is blocked. Exiting the loop.")
		// 		rssFeed = nil
		// 	}
		// 	if rssFeed == nil {
		// 		break
		// 	}
		// }

		rss := generateRSS(resBuilder.String())

		if err := os.WriteFile(outputFile, []byte(rss), 0o644); err != nil {
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

func generateRSS(res string) string {
	feed := &feeds.Feed{
		Title:       "docs",
		Link:        &feeds.Link{Href: "https://github.com/hxhacking/docs"},
		Description: "just for fun",
	}

	// res := generateHTML(qs)
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

func generateHTML(qs qs) string {
	questions := fmt.Sprintf("\n\n\n# %s \n\n\n", qs.FileName)
	for _, item := range qs.Doc {
		questions += fmt.Sprintf("## %s \n\n", item.Docs)
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
	if err := markdown.Convert([]byte(md), &buf); err != nil {
		return ""
	}

	return buf.String()
}

// func StringToBytes(s string) []byte {
// 	return unsafe.Slice(unsafe.StringData(s), len(s))
// }
