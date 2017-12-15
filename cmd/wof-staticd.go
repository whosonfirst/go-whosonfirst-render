package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-http-mapzenjs"
	"github.com/whosonfirst/go-whosonfirst-render/http"
	"github.com/whosonfirst/go-whosonfirst-render/reader"
	"log"
	gohttp "net/http"
	"os"
	"path/filepath"
)

func main() {

	var host = flag.String("host", "localhost", "The hostname to listen for requests on")
	var port = flag.Int("port", 8080, "The port number to listen for requests on")

	var source = flag.String("source", "fs", "...")
	var root = flag.String("fs-root", "", "...")

	var s3_bucket = flag.String("s3-bucket", "whosonfirst.mapzen.com", "...")
	var s3_prefix = flag.String("s3-prefix", "", "...")
	var s3_region = flag.String("s3-region", "us-east-1", "...")
	var s3_creds = flag.String("s3-credentials", "", "...")

	var api_key = flag.String("mapzen-apikey", "mapzen-xxxxxxx", "")

	flag.Parse()

	var r reader.Reader
	var err error

	switch *source {
	case "fs":
		r, err = reader.NewFSReader(*root)
	case "s3":

		cfg := reader.S3Config{
			Bucket:      *s3_bucket,
			Prefix:      *s3_prefix,
			Region:      *s3_region,
			Credentials: *s3_creds,
		}

		r, err = reader.NewS3Reader(cfg)
	default:
		err = errors.New("Unknown or invalid source")
	}

	if err != nil {
		log.Fatal(err)
	}

	// my kingdom for a generic middleware handler that I can use
	// to update markup generated by another handler...
	// (20171214/thisisaaronland)

	html_opts := http.NewDefaultHTMLOptions()
	html_opts.MapzenAPIKey = *api_key

	html_handler, err := http.HTMLHandler(r, html_opts)

	if err != nil {
		log.Fatal(err)
	}

	svg_handler, err := http.SVGHandler(r)

	if err != nil {
		log.Fatal(err)
	}

	spr_handler, err := http.SPRHandler(r)

	if err != nil {
		log.Fatal(err)
	}

	ping_handler, err := http.PingHandler()

	if err != nil {
		log.Fatal(err)
	}

	static_handler, err := http.StaticHandler()

	if err != nil {
		log.Fatal(err)
	}

	mapzenjs_handler, err := mapzenjs.MapzenJSHandler()

	if err != nil {
		log.Fatal(err)
	}

	id_func := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		path := req.URL.Path
		ext := filepath.Ext(path)
		
		switch ext {
		case ".svg":
			svg_handler.ServeHTTP(rsp, req)
		case ".spr":
			spr_handler.ServeHTTP(rsp, req)
		default:
			html_handler.ServeHTTP(rsp, req)
		}

		return
	}

	id_handler := gohttp.HandlerFunc(id_func)

	mux := gohttp.NewServeMux()

	mux.Handle("/id/", id_handler)
	mux.Handle("/svg/", svg_handler)
	mux.Handle("/spr/", spr_handler)

	mux.Handle("/ping", ping_handler)

	mux.Handle("/javascript/mapzen.min.js", mapzenjs_handler)
	mux.Handle("/javascript/tangram.min.js", mapzenjs_handler)
	mux.Handle("/javascript/mapzen.js", mapzenjs_handler)
	mux.Handle("/javascript/tangram.js", mapzenjs_handler)
	mux.Handle("/css/mapzen.js.css", mapzenjs_handler)
	mux.Handle("/tangram/refill-style.zip", mapzenjs_handler)

	mux.Handle("/javascript/slippymap.crosshairs.js", static_handler)
	mux.Handle("/javascript/whosonfirst.spr.js", static_handler)
	mux.Handle("/css/whosonfirst.spr.css", static_handler)

	address := fmt.Sprintf("%s:%d", *host, *port)
	log.Printf("listening on %s\n", address)

	err = gohttp.ListenAndServe(address, mux)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
