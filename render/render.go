package render

import (
	"context"

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

func RenderImage(templateHTML string, options RenderOptions) ([]byte, error) {
	// create context
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()
	var buf []byte
	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(ctx, fullScreenshot(ctx, templateHTML, options, &buf)); err != nil {
		return nil, err
	}

	return buf, nil
}

func fullScreenshot(ctx context.Context, templateHTML string, options RenderOptions, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		emulation.SetDeviceMetricsOverride(int64(options.Width), int64(options.Height), 1.0, false),
		chromedp.Navigate("data:text/html,"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, templateHTML).Do(ctx)
		}),
		chromedp.FullScreenshot(res, options.Quality),
	}
}
