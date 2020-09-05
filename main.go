package main

import (
	"bytes"
	b64 "encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/h2non/bimg"
	"github.com/melbahja/got"
)

func index(rw http.ResponseWriter, r *http.Request) {
	var options bimg.Options
	imageSrc, ok := r.URL.Query()["src"]
	if !ok {
		fmt.Fprintf(rw, "%q", "No image")
		return
	}

	var b bytes.Buffer
	cacheFolder := b64.StdEncoding.EncodeToString([]byte(r.URL.String()))
	b.WriteString("cache/")
	b.WriteString(cacheFolder)
	b.WriteString("/")
	b.WriteString(strings.Split(path.Base(imageSrc[0]), ".")[0])

	w, hasWidth := r.URL.Query()["w"]
	h, hasHeight := r.URL.Query()["h"]
	o, hasOptions := r.URL.Query()["o"]
	t, hasType := r.URL.Query()["t"]
	c, hasCompression := r.URL.Query()["c"]

	if hasType {
		switch t[0] {
		case "jpg":
			options.Type = bimg.JPEG
			b.WriteString(".jpg")
		case "webp":
			options.Type = bimg.WEBP
			b.WriteString(".webp")
		case "png":
			options.Type = bimg.PNG
			b.WriteString(".png")
		}
	} else if strings.Contains(r.Header.Get("Accept"), "webp") {
		options.Type = bimg.WEBP
		b.WriteString(".webp")
	} else {
		b.WriteString(path.Ext(imageSrc[0]))
	}

	var opts []string
	if hasOptions {
		opts = strings.Split(o[0], ",")
	}

	name := b.String()
	rw.Header().Set("Content-Disposition", fmt.Sprintf(`filename="%s"`, path.Base(name)))
	
	if fileExists(name) && !Contains(opts, "fresh") {
		http.ServeFile(rw, r, name)
		return
	}
	
	err := os.MkdirAll("cache/" + cacheFolder, 0755)
	if err != nil {
		fmt.Fprintf(rw, "%q", "Could not create folder")
	}

	g := got.New()
	err = g.Download(imageSrc[0], name)
	if err != nil {
		fmt.Fprintf(rw, "%q", "Could not download the image")
		return
	}

	if hasOptions {
		if Contains(opts, "crop") {
			options.Crop = true
		}

		if Contains(opts, "smart") {
			options.Gravity = bimg.GravitySmart
		}

		if Contains(opts, "flip") {
			options.Flip = true
		}

		if Contains(opts, "flop") {
			options.Flop = true
		}

		if Contains(opts, "bw") {
			options.Interpretation = bimg.InterpretationBW
		}
	}

	if hasWidth && hasHeight && !options.Crop {
		options.Enlarge = true
	}

	if hasWidth {
		options.Width, _ = strconv.Atoi(w[0])
	}

	if hasHeight {
		options.Height, _ = strconv.Atoi(h[0])
	}

	if hasCompression {
		options.Compression, _ = strconv.Atoi(c[0])
	}

	options.StripMetadata = true

	read, _ := bimg.Read(name)
	img, err := bimg.NewImage(read).Process(options)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	bimg.Write(name, img)

	rw.Write(img)
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if strings.ToLower(x) == n {
			return true
		}
	}
	return false
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", index)

	log.Fatal(http.ListenAndServe(":1985", mux))
}
