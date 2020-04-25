package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/reujab/wallpaper"
	"github.com/turnage/graw/reddit"
)

var (
	landscapeRegex = regexp.MustCompile(`\[(\d+)x(\d+)\]`) // e.g. [1024x768]
	subreddit      = "/r/EarthPorn"
)

func main() {
	bot, _ := reddit.NewBotFromAgentFile("./wallpaperbot.agent", 0)
	url := getWallpaperImage(bot)
	file := getImage(url)

	wallpaper.SetFromFile(file)
}

func getWallpaperImage(bot reddit.Bot) string {
	listing, _ := bot.Listing(subreddit, "")

	// Only include images
	n := 0
	for _, x := range listing.Posts {
		if isImage(x) {
			listing.Posts[n] = x
			n++
		}
	}
	listing.Posts = listing.Posts[:n]

	// Sort by upvotes
	sort.SliceStable(listing.Posts, func(i, j int) bool {
		return listing.Posts[i].Ups > listing.Posts[j].Ups
	})

	// Pick appropriate post
	var post *reddit.Post
	for _, p := range listing.Posts {
		if isLandscape(p) && isWithin24Hrs(p) {
			post = p
			fmt.Println("a")
			break
		}
	}

	if post == nil {
		// Just return the most upvoted image
		post = listing.Posts[0]
	}

	fmt.Println(fmt.Sprintf("%s, https://old.reddit.com%s , %s", post.Title, post.Permalink, post.URL))
	return post.URL
}

func isImage(p *reddit.Post) bool {
	return strings.HasSuffix(p.URL, "png") || strings.HasSuffix(p.URL, "jpg")
}

func isWithin24Hrs(p *reddit.Post) bool {
	created := time.Unix(int64(p.CreatedUTC), 0)
	diff := time.Now().UTC().Sub(created)

	return diff < time.Hour*24
}

func isLandscape(p *reddit.Post) bool {
	s := landscapeRegex.FindStringSubmatch(p.Title)
	if len(s) != 3 {
		return false
	}

	w, err := strconv.Atoi(s[1])
	if err != nil {
		return false
	}
	h, err := strconv.Atoi(s[2])
	if err != nil {
		return false
	}

	return w > h
}

func getImage(url string) string {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	tmp := os.TempDir()
	img := "wallpaper"
	filename := path.Join(tmp, img)

	outFile, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer outFile.Close()
	_, err = io.Copy(outFile, res.Body)

	return filename
}
