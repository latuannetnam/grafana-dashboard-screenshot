package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type ProgOpts struct {
	grafanaHost        string
	grafanaPort        int
	grafanaProtocol    string
	grafanaApiToken    string
	grafanaPrefix      string
	grafanaDashboardDd string
	grafanaVariables   string
	outputFile         string
	waitTime           int
}

func dashboardScreenshot(opts ProgOpts) error {
	var buf []byte
	chromeOpts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		// chromedp.Flag("headless", false),
		// chromedp.WindowSize(1920, 4000),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), chromeOpts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	url := opts.grafanaProtocol + "://" + opts.grafanaHost + ":" + fmt.Sprintf("%d", opts.grafanaPort)
	// url := opts.grafanaProtocol + "://" + opts.grafanaHost
	if opts.grafanaPrefix != "" {
		url += "/" + opts.grafanaPrefix
	}
	url += "/d/" + opts.grafanaDashboardDd + "/?orgId=1&kiosk"
	if opts.grafanaVariables != "" {
		url += "&" + opts.grafanaVariables
	}
	log.Printf("url:%s\n", url)

	err := chromedp.Run(ctx,
		network.Enable(),
		network.SetExtraHTTPHeaders(network.Headers(map[string]interface{}{"Authorization": "Bearer " + opts.grafanaApiToken})),
		chromedp.EmulateViewport(1920, 1080),
		chromedp.Navigate(url),
		// chromedp.ActionFunc(func(ctx context.Context) error {
		// 	_, exp, err := runtime.Evaluate(`window.scrollTo(0,document.body.scrollHeight);`).Do(ctx)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	if exp != nil {
		// 		return exp
		// 	}
		// 	return nil
		// }),
		// chromedp.Sleep(time.Duration(opts.waitTime)*time.Second),
		// chromedp.KeyEvent(kb.End),
		chromedp.Sleep(time.Duration(opts.waitTime)*time.Second),
		// chromedp.FullScreenshot(&buf, 90),
		printToPDF(ctx, &buf),
	)

	if err != nil {
		log.Fatal(err)
	} else {
		// if err := os.WriteFile("fullScreenshot.png", buf, 0o644); err != nil {
		// 	log.Fatal(err)
		// }

		if err := os.WriteFile(opts.outputFile, buf, 0o644); err != nil {
			log.Fatal(err)
		}
	}
	return err
}

func printToPDF(ctx context.Context, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithLandscape(true).
				WithPaperWidth(16.5).
				WithPaperHeight(23.4).
				WithPageRanges("1-2").
				WithPrintBackground(true).Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}

func main() {
	startTime := time.Now()
	progOpts := ProgOpts{}
	flag.StringVar(&progOpts.grafanaHost, "grafana_host", "localhost", "Grafana host")
	flag.IntVar(&progOpts.grafanaPort, "grafana_port", 3000, "Grafana port")
	flag.StringVar(&progOpts.grafanaProtocol, "grafana_protocol", "http", "Grafana protocol (http|https)")
	flag.StringVar(&progOpts.grafanaApiToken, "grafana_api_token", "", "Grafana API token")
	flag.StringVar(&progOpts.grafanaPrefix, "grafana_prefix", "", "Grafana prefix (e.g. /grafana)")
	flag.StringVar(&progOpts.grafanaDashboardDd, "grafana_dashboard_id", "", "Grafana dashboard ID")
	flag.StringVar(&progOpts.grafanaVariables, "grafana_variables", "", "Grafana variables")
	flag.StringVar(&progOpts.outputFile, "output_file", "output.pdf", "Output file")
	flag.IntVar(&progOpts.waitTime, "wait_time", 10, "Wait time")
	flag.Parse()
	log.Printf("progOpts:%v\n", progOpts)
	dashboardScreenshot(progOpts)
	elapsedTime := time.Since(startTime)
	log.Printf("Done!. Took:%s", elapsedTime)
}
