# Kitty-Slack
A slack bot to annoy your colleagues with cat gifs

![Meow](https://upload.wikimedia.org/wikipedia/commons/thumb/3/3c/Creative-Tail-Animal-cat.svg/240px-Creative-Tail-Animal-cat.svg.png)

## Commands

Write `@kitty` with one following messages:

* `meow` or `gimme more` - Sends you a random cat gif
* `meow!` or `gimme more!` - Sends you 5 random cat gifs! Whoop!


## Installation

Make sure you have an environment variable `SLACK_TOKEN` set containing the oauth token string. Build from source using `go build` or

Use in a docker container

```
# Compile
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

# Build
docker build -t cbrgm/kitty-slack .

# Run
docker run -itd --name kitty-slack -e token=<SLACK_TOKEN> cbrgm/kitty-slack
```
