# Eventprint ![Analytics](https://ga-beacon.appspot.com/UA-34529482-6/eventprint/readme?pixel)

Eventprint is a check-in solution for hackathons. It has a web interface for
checking in attendees and a desktop client that generates badges and prints
them.

## Getting Started

### Development

Install [Docker](https://www.docker.io/) and
[Fig](http://orchardup.github.io/fig/). Go into the `server` directory and run
`fig up`.

### Production

#### Server

1. Set the correct config vars in `server/config/config.yml` (or export them as
   environment variables). Make sure you replace the comment entirely with the
   value you set.
2. Start the server with `go run server.go` when you're in the `server`
   directory.

#### Client

1. In the `client` directory run `go run main.go`. Run `go run main.go -help`
   for options.

#### Gotchas

Eventprint makes a few assumptions about your environment and will sometimes
silently fail.

* It assumes all of the environment variables are set correctly (including the
  database).
* It assumes that a client is attached when checking attendees in. If one is
  not, it'll silently fail.
