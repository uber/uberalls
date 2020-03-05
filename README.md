# uberalls [![Build Status](https://travis-ci.org/uber/uberalls.svg?branch=master)](https://travis-ci.org/uber/uberalls) [![Coverage Status](https://coveralls.io/repos/uber/uberalls/badge.svg)](https://coveralls.io/r/uber/uberalls)

Code coverage metric storage service. Provide coverage metrics on differentials
with [Phabricator][] and [Jenkins][], just like [Coveralls][] does for GitHub
and TravisCI.

[Phabricator]: http://phabricator.org/
[Jenkins]: https://jenkins-ci.org/
[Coveralls]: https://coveralls.io/

## Docker
```console
docker run -p 3000:3000 haxandmat/uberalls:latest
```
## ... via [`docker-compose`](https://github.com/docker/compose)

Example `docker-compose.yml` for `uberalls`:

```yaml
version: '3'

services:
  server:
    restart: unless-stopped
    image: haxandmat/uberalls:latest
    ports:
      - "3000:3000"

```

Run `docker-compose up -d`, wait for it to initialize completely, and visit `http://localhost:3000`, or `http://host-ip:3000` (as appropriate).

## Running

Configure your database by editing `config/default.json` or specify a different
file by passing an `UBERALLS_CONFIG` environment variable. The development and
test runs use a SQLite database.

If you'd like to use MySQL, you could use the following configuration:

```json
{
  "dbType": "mysql",
  "dbLocation": "user:password@/dbname?charset=utf8",
  "listenPort": 8080,
  "listenAddress": "0.0.0.0"
}
```

## Jenkins integration

Uberalls works best when paired with our [Phabricator Jenkins Plugin][], which
will record Cobertura data on master runs, and compare coverage to the base
revision on differentials.

![Jenkins Integration](https://raw.githubusercontent.com/Haxandmat/uberalls/master/docs/jenkins-integration.png)

[Phabricator Jenkins Plugin]: https://github.com/uber/phabricator-jenkins-plugin

On differentials, the delta in coverage is calculated by taking the base commit
(from conduit metadata) and comparing that to the current coverage amount.

In order to have a baseline to compare against, you must also have jenkins build
your project on your mainline branch ("master" by default). You can either
create a separate job, or enable SCM polling under Build Triggers:

![scm polling](https://raw.githubusercontent.com/Haxandmat/uberalls/master/docs/scm-polling.png)

## Development

Get the source

```bash
go get github.com/uber/uberalls
```

Install Glide and dependencies

```bash
cd $GOPATH/src/github.com/uber/uberalls
go get github.com/Masterminds/glide
glide install
```

Run the thing

```bash
go build && ./uberalls
```

Run the tests

```bash
go test .
```

## License

MIT Licensed
