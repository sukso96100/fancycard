package render

import (
	"context"
	"sync"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type RenderOptions struct {
	Width   int
	Height  int
	Quality int
}

var DefaultRenderOptions RenderOptions = RenderOptions{1200, 600, 90}

func RenderImage(templateHTML string, externalTemplateURL string, options RenderOptions) ([]byte, error) {
	// create context
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()
	var buf []byte
	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(ctx, fullScreenshot(ctx, templateHTML, externalTemplateURL, options, &buf)); err != nil {
		return nil, err
	}

	return buf, nil
}

func fullScreenshot(ctx context.Context, templateHTML string, externalTemplateURL string, options RenderOptions, res *[]byte) chromedp.Tasks {
	var navURL string
	if externalTemplateURL != "" {
		navURL = externalTemplateURL
	} else {
		navURL = "data:text/html,"
	}
	var wg sync.WaitGroup
	return chromedp.Tasks{
		chromedp.Navigate(navURL),
		chromedp.ActionFunc(func(ctx context.Context) error {

			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			chromedp.ListenTarget(ctx, func(ev interface{}) {
				switch ev.(type) {
				case *page.EventLoadEventFired:
					wg.Done()
				}
			})
			wg.Add(1)
			return page.SetDocumentContent(frameTree.Frame.ID, templateHTML).Do(ctx)
		}),
		chromedp.ActionFunc(func(c context.Context) error {
			wg.Wait()
			return nil
		}),
		emulation.SetDeviceMetricsOverride(int64(options.Width), int64(options.Height), 1.0, false),
		chromedp.CaptureScreenshot(res),
	}
}
