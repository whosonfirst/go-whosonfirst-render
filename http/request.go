package http

import (
	"errors"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"github.com/whosonfirst/warning"
	_ "log"
	gohttp "net/http"
	"path/filepath"
	"regexp"
	"strconv"
)

var re_wofid *regexp.Regexp

func init() {
	re_wofid = regexp.MustCompile(`^(\d+)(?:(?:\-alt\-.*)?\.[^\.]+)?$`)
}

func IdFromPath(path string) (int64, error) {

	abs_path, err := filepath.Abs(path)

	if err != nil {
		return -1, err
	}

	fname := filepath.Base(abs_path)

	match := re_wofid.FindAllStringSubmatch(fname, -1)

	if len(match) == 0 {
		return -1, errors.New("Unable to parse WOF ID")
	}

	if len(match[0]) != 2 {
		return -1, errors.New("Unable to parse WOF ID")
	}

	wofid, err := strconv.ParseInt(match[0][1], 10, 64)

	if err != nil {
		return -1, err
	}

	return wofid, nil
}

func FeatureFromRequest(req *gohttp.Request, r reader.Reader) (geojson.Feature, error, int) {

	path := req.URL.Path

	wofid, err := IdFromPath(path)

	if err != nil {

		q := req.URL.Query()
		str_id := q.Get("id")

		if str_id == "" {
			return nil, err, gohttp.StatusNotFound
		}

		id, err := strconv.ParseInt(str_id, 10, 64)

		if err != nil {
			return nil, err, gohttp.StatusBadRequest
		}

		wofid = id
	}

	rel_path, err := uri.Id2RelPath(wofid)

	if err != nil {
		return nil, err, gohttp.StatusBadRequest // StatusInternalServerError
	}

	ctx := req.Context()

	fh, err := r.Read(ctx, rel_path)

	if err != nil {
		return nil, err, gohttp.StatusBadRequest // StatusInternalServerError
	}

	f, err := feature.LoadFeatureFromReader(fh)

	if err != nil && !warning.IsWarning(err) {
		return nil, err, gohttp.StatusInternalServerError
	}

	return f, nil, 0
}