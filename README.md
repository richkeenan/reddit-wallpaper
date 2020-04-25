# Reddit Wallpaper

Grabs the most upvoted, landscape image from /r/EarthPorn within the last 24 hours and sets it as your current  desktop background. You can run this once, or install it as a service or through cron/Task Scheduler.

## Building

`go build main.go`

## Running

To adhere to Reddit's API usage rules you need to run this under an account - a completely separate account from your main Reddit account is recommended. Follow the instructions [here](https://turnage.gitbooks.io/graw/content/chapter1.html) for setting up the account with an Application.

Create a file called `wallpaperbot.agent` with the following contents:

```
user_agent: "windows:wallpaperbot:0.0.1 (by /u/{reddit username})"
client_id: "{client id}"
client_secret: "{client secret}"
username: "{reddit username}"
password: "{reddit password}"
```

Run using `go run main.go` for a one-off invocation. Or run `go build main.go` and point a cron job, Task Scheduler etc at `main.exe` to run it on a schedule.
