# toggl-rounder
Small script in Go to round [Toggl](https://www.toggl.com/app/timer) entries up to full half hours and that starts always at full hour.


## How to build?

Run standard `go build` command.

To build all versions for release please run `./build.sh`.

## How to use it?

```
Usage of ./toggl-rounder:
  -api-key secret-key
    	Toggl API KEY secret-key, can also be provided via $TOGGL_API_KEY environment variable (default "")
  -debug
    	Print debugging output of API calls
  -dry
    	Unless set to false it doesn't update records in Toggl (default true)
  -remaining string
    	Decides on what to do with remaining time. Possible options: ["keep" "round_half" "round_up"] (default "keep")
  -version
    	Print the version
    	
```

Plus please ensure you've setup `TOGGL_API_KEY` environment variable or provided the `-api-key` parameter so it would know how to connect to the [Toggl API](https://github.com/toggl/toggl_api_docs).

Full command to run update:

```TOGGL_API_KEY=toggl-s3cret ./toggl-rounder -dry=false -remaining round_up```

--- 
Enjoy!

(c) Rafal "RaVbaker" Piekarski 2019

License? as specified in LICENSE file.
