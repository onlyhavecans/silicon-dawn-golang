# Silicon Dawn

I made this so I can pull random cards on my phone and review the text.

I originally designed this for use in [Pythonista](http://omz-software.com/pythonista/) but I rewrote it in rust to be a web app.
Then I got very tired of maintaining and compiling rust for the app.
I wanted to take it to the NEXT LEVEL and rewrite it in golang and stuff it in a docker container.


## Instructions

### Docker

I publish the docker container at `skwrl/silicon-dawn`.
It is fully self-contained and uses port 3200 internally.

1. Install docker however you choose
1. copy the `docker-compose.yml` from this repo
1. Change the port settings as desired
1. `docker-compose up`

### Bare Go Binary

1. Install golang
1. Check out this repository wherever you choose
1. `go run . get` to hydrate the cards data directory
1. `go run . serve` to start the webserver
1. Browse to http://localhost:3200 to enjoy your pick
1. Refresh the page for a fresh pick

Or! if you are super lazy check out [my hosted copy of this](https://silicon-dawn.cards).