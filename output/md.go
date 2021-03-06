package output

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ma6254/FictionDown/store"
	pb "gopkg.in/cheggaaa/pb.v1"
)

type Markdown struct {
}

func (t *Markdown) Conv(src store.Store, outpath string) (err error) {

	f, err := os.Create(outpath)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "---\n")
	fmt.Fprintf(f, "title: %#v\n", src.BookName)
	fmt.Fprintf(f, "description: %#v\n", src.Description)
	fmt.Fprintf(f, "creator: %#v\n", src.Author)
	fmt.Fprintf(f, "lang: %s\n", "zh-CN")

	if src.CoverURL != "" {

		client := &http.Client{}
		req, err := http.NewRequest("GET", src.CoverURL, nil)
		if err != nil {
			return err
		}
		req.Header.Add(
			"user-agent",
			"Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Mobile Safari/537.36",
		)
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		coverBuf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		tempfile, err := ioutil.TempFile("", "book_cover_*.jpg")
		if err != nil {
			return err
		}

		ioutil.WriteFile(tempfile.Name(), coverBuf, 0775)

		log.Printf("Save Cover Image: %#v", tempfile.Name())

		fmt.Fprintf(f, "cover-image: %#v\n", tempfile.Name())
	}
	fmt.Fprintf(f, "---\n\n")

	fmt.Fprintf(f, "# 简介\n\n")
	dlist := strings.Split(src.Description, "\n")

	for _, cc := range dlist {
		fmt.Fprintf(f, "<p style=\"text-indent:2em\">%s</p>\n",
			strings.Replace(cc, "*", "□", -1),
		)
	}
	fmt.Fprintf(f, "\n")

	for _, v1 := range src.Volumes {
		var VIP string
		if v1.IsVIP {
			VIP = "收费"
		} else {
			VIP = "免费"
		}
		fmt.Fprintf(f, "# %#v_%s\n\n", v1.Name, VIP)
		log.Printf("正在转换卷: %s", v1.Name)
		bar := pb.StartNew(len(v1.Chapters))
		for _, v2 := range v1.Chapters {
			// s += fmt.Sprintf(`<h1><a href=%#v>%s</a></h1>`, v2.URL, v2.Name)
			fmt.Fprintf(f, "## [%s](%s)\n\n", v2.Name, v2.URL)
			for _, cc := range v2.Text {
				fmt.Fprintf(f, "<p style=\"text-indent:2em\">%s</p>\n",
					strings.Replace(cc, "*", "□", -1),
				)
			}
			bar.Increment()
			fmt.Fprintf(f, "\n")
		}
		bar.Finish()
	}
	return nil
}
