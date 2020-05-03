# Silicon Dawn

I made this so I can pull random cards on my phone and review the text.

I originally designed this for use in [Pythonista](http://omz-software.com/pythonista/)
Then I rewrote it in rust to be a web app. Then I got very tired of compiling rust.
I wanted to take it to the NEXT LEVEL and rewrite it in golang and stuff it in a docker container.

## Instructions

### Docker

I publish the docker container at [DockerHub](https://hub.docker.com/r/skwrl/silicon-dawn).
It is fully self-contained and uses port 3200 internally.

You can spin up a copy however you choose to do a docker or use my compose files in `/compose`

### Go Binary

1. Install golang
1. Check out this repository wherever you choose
1. `go build -o bin/silicon-dawn`
1. `./bin/silicon-dawn get` to hydrate the cards data directory
1. `./bin/silicon-dawn serve` to start the webserver
1. Browse to http://localhost:3200 to enjoy your pick
1. Refresh the page for a fresh pick

Or! if you are super lazy check out [my hosted copy of this](https://silicon-dawn.cards).
