package www

import (
	"errors"
	"fmt"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-sanitize"
	"github.com/whosonfirst/go-whosonfirst-image"
	"net/http"
)

type RasterSize struct {
	Label     string
	MaxHeight int
	MaxWidth  int
}

type RasterOptions struct {
	Format string
	Sizes  map[string]RasterSize
}

func NewDefaultRasterOptions() (*RasterOptions, error) {

	xsm := RasterSize{
		Label:     "xsm",
		MaxHeight: 100,
		MaxWidth:  100,
	}

	sm := RasterSize{
		Label:     "sm",
		MaxHeight: 300,
		MaxWidth:  300,
	}

	med := RasterSize{
		Label:     "med",
		MaxHeight: 640,
		MaxWidth:  640,
	}

	lg := RasterSize{
		Label:     "lg",
		MaxHeight: 1024,
		MaxWidth:  1024,
	}

	sz := map[string]RasterSize{
		"xsm": xsm,
		"sm":  sm,
		"med": med,
		"lg":  lg,
	}

	opts := RasterOptions{
		Format: "png",
		Sizes:  sz,
	}

	return &opts, nil
}

func RasterHandler(r reader.Reader, opts *RasterOptions) (http.Handler, error) {

	if opts.Format != "png" {
		return nil, errors.New("Invalid or unsupported raster format")
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		uri, err, status := ParseURIFromRequest(req, r)

		if err != nil {
			http.Error(rsp, err.Error(), status)
			return
		}

		sn_opts := sanitize.DefaultOptions()

		sz := "lg"

		query := req.URL.Query()
		query_sz := query.Get("size")

		req_sz, err := sanitize.SanitizeString(query_sz, sn_opts)

		if err != nil {
			http.Error(rsp, err.Error(), status)
			return
		}

		if req_sz != "" {
			sz = req_sz
		}

		sz_info, ok := opts.Sizes[sz]

		if !ok {
			http.Error(rsp, "Invalid output size", http.StatusBadRequest)
			return
		}

		img_opts := image.NewDefaultOptions()

		img_opts.Writer = rsp
		img_opts.Height = sz_info.MaxHeight
		img_opts.Width = sz_info.MaxWidth

		content_type := fmt.Sprintf("image/%s", opts.Format)
		rsp.Header().Set("Content-Type", content_type)

		f := uri.Feature
		image.FeatureToPNG(f, img_opts)
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
