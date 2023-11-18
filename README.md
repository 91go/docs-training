# docs-training

```shell

go install github.com/91go/docs-training@latest

```


Used to extract questions from private docusaurus blog. We can use these questions as interview questions or just as a breakthrough-point to familiarize ourselves with the notes.

Often we are faced with problems that we have solved before and then are unable to solve them when they arise again. Much of this is due to lack of familiarity with your notes, and this tool will help you become more familiar with your notes. 


---

You can use this command in github actions, for example


```yaml

jobs:
  docs-training:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/setup-go@v3
        with:
          go-version: '1.20'
      - uses: actions/checkout@v3
      - id: generate
        run: |
          go install github.com/91go/docs-training@latest
          echo -e "$(docs-training action -n ${{ env.NUM }} -w algo)" >> ms-algo.md
          echo -e "$(docs-training action -n ${{ env.NUM }} -w arch)" >> ms-arch.md
          echo -e "$(docs-training action -n ${{ env.NUM }} -w lang)" >> ms-lang.md
          echo -e "$(docs-training action -n ${{ env.NUM }} -w works)" >> works.md
      - name: Concatenate Markdown variables
        run: |
          {
            echo "## algo" && cat ms-algo.md
            echo "## arch" && cat ms-arch.md
            echo "## lang" && cat ms-lang.md
            echo "## works" && cat works.md
          } >> output.md
      - name: Convert Markdown to HTML
        run: |
          go install github.com/kpym/gm@latest
          gm output.md > output.html
      - name: Convert HTML to RSS
        run: |
          go install github.com/91go/feedgen@latest
          feedgen gen --filename output.html --title ${{ env.RSS_TITLE }} --author ${{ env.RSS_AUTHOR }} --description ${{ env.RSS_DESCRIPTION }} --mail ${{ env.RSS_MAIL }} > feed.xml
      - ...git commit and push OR publish to s3 OR other any other service can storage feed.xml file

```
