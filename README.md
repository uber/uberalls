# uberalls

Code coverage metric storage service. Provide coverage metrics on differentials
with [Phabricator][] and [Jenkins][], just like [Coveralls][] does for GitHub
and TravisCI.

[Phabricator]: http://phabricator.org/
[Jenkins]: https://jenkins-ci.org/
[Coveralls]: https://coveralls.io/

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

![Jenkins Integration](/docs/jenkins-integration.png)

[Phabricator Jenkins Plugin]: https://github.com/uber/phabricator-jenkins-plugin

## Development

Get the source
```bash
go get github.com/uber/uberalls
```

Install godeps
```bash
cd $GOPATH/src/github.com/uber/uberalls
go get github.com/tools/godep
godep restore
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
