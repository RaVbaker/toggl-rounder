# toggl-rounder
Small script in Go to round [Toggl](https://www.toggl.com/app/timer) entries up to full half hours and that starts always at full hour.


## How to build?

Run standard `go build` command.

## How to use it?

```
Usage of ./toggl-rounder:
  -dry
    	Should it update toggl? (default true)
  -rounding
    	Should it round last entry?

```

Plus please ensure you've setup `TOGGL_API_KEY` environment variable so it would know how to connect to the [Toggl API](https://github.com/toggl/toggl_api_docs).

Full command to run update:

```TOGGL_API_KEY=toggl-s3cret ./toggl-rounder -dry=false -rounding=true```

--- 
Enjoy!

(c) Rafal "RaVbaker" Piekarski 2019

License? as specified in LICENSE file.
