# fancycard
Easily customizable Social image (or Open graph image) generator   
Built with Go, Gin, GoQuery and Chromedp

## Build & Run
Simply, Clone this repo, build and run with Go Toolchain.  
You can get Go Toolchain from [here](https://golang.org/dl/).  
You will also need Chrome or Chromium installed on your system to make this program work properly.
```bash
git clone https://github.com/sukso96100/fancycard.git
cd fancycard
go build -o fancycard .
./fancycard
```

## Render Social image

### Pass all required data via URL

Build a URL using `/url` API with following format, Then place URL in your webpage as [Open Graph Protocol Image meta tag](https://ogp.me/) or [Twitter card Image meta tag](https://developer.twitter.com/en/docs/twitter-for-websites/cards/guides/getting-started)
```
https://<Host>/url?
    template=<Built in template name or Remote template URL>
    &<Template param key>=<Template param value>
    &<Template param key>=<Template param value>
    & ...
```

- `<Host>`: Address of fancycard server where fancycard instance is hosted (e.g. `fancycard.app.com`)
- `<Built in template name or Remote template URL>`: Path or URL to tempalte writtn in HTML and Go Template syntax
    - Built in templates: `simple.html`
    - Remote template URL: e.g. `https://raw.githubusercontent.com/sukso96100/fancycard/main/tmpl/templates/simple.html`
- `<Template param key>`: Key of the parameter that the template requires
    - e.g. `Title`, `Date`, `Author`, `Img`...
    - Check template file for what template requires
- `<Template param value>`: Value for template parameter
    - e.g. `Hello world`

#### Example usage
```html
<html>
    <head>
        <!-- Open Graph Protocol - Image Meta Tag -->
        <meta property="og:image" content="https://<Host>/url?template=<Built in template name or Remote template URL>&<Template param key>=<Template param value>&<Template param key>=<Template param value>& ..." />
        <!-- Twitter Card - Image Meta Tag -->
        <meta name="twitter:image" content="https://<Host>/url?template=<Built in template name or Remote template URL>&<Template param key>=<Template param value>&<Template param key>=<Template param value>& ..." />
    </head>
...
</html>
```
### Let fancycard to scrap required data from your website

Instead of building super-long URL, Use `/meta` API, And put required data as meta tags with `name="fancycard:<key>"` attributes.

#### Example usage
```html
<html>
    <head>
        <!-- Open Graph Protocol - Image Meta Tag -->
        <meta property="og:image" content="https://<Host>/meta?url=<Your webpage URL on internet>" />
        <!-- Twitter Card - Image Meta Tag -->
        <meta name="twitter:image" content="https://<Host>/meta?url=<Your webpage URL on internet>" />

        <!-- Fancycard meta tags -->

        <!-- Built-in Template name or Remote Template URL (Required)  -->
        <meta name="fancycard:template" content="simple.html" />

        <!-- Other parameters that the template requires -->
        <meta name="fancycard:<Key>" content="<Value>" />
        <meta name="fancycard:Title" content="Hello world!" />
    </head>
...
</html>
```

## License
MIT License
