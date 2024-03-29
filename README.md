# Go Philae v5.2.2

Go Philae is the go implementation of our Philae health check protocol.

## Architecture

The go-philae library is based around the idea of probe. Each service has some dependencies. Those dependencies are checked using probes.

The prober is the component that will take all probes, check everyone of them, and return an aggregated result into a single response.

Finally, the handler is a small utility that will take a prober and an existing router and add an `/_health` endpoint to that router. This endpoint is configured to run check every time someone call it.

## Usage

Import:

```
github.com/Scalingo/go-philae/v5
```

To use it in an existing project, you will need to add a prober with some probes, pass it to the handler and generate the route.

```go
router := handlers.NewRouter("http")
// ... configure your routes

probes := prober.NewProber()
probes.AddProbe(redisprobe.NewRedisProbeFromURL("redis", config.E["REDIS_URL"]))
/// ... add as many probes as needed

globalRouter := philaehandler.NewPhilaeRouter(router, prober)

http.ListenAndServe(":8080", globalRouter)
```

## Creating a probe

To create a probe, you will need to implement the following interface:

```go
type Probe interface {
  Name() string // Return the name of the probe
  Check(context.Context) error // Return nil if the probe check was successfull or an error otherwise
}
```

That's all folks.

### Convention

* The check must be as lightweight as possible
* The check should not modify data used by the service
* The check should not depend on the service state
* The check should not take more than 3 second

## Existing probes

### DockerProbe

Check that the docker daemon is running.

#### Usage

```go
dockerprobe.NewDockerProbe(name, endpoint string)
```

* name: The name of the probe
* endpoint: The docker API endpoint

### EtcdProbe

Check that a etcd server is running

#### Usage

```go
etcdprobe.NewEtcdProbe(name string, client etcd.KeysAPI)
```

* name: The name of the probe
* client: An etcd Keys API client correctly configured

### GithubProbe

Check that GitHub isn't reporting any issue.
It use the official GitHub status API to check if there is no "major" problem with the GitHub infrastructure.

#### Usage

```go
githubprobe.NewGithubProbe(name string)
```

* name: The name of the probe

### GitlabProbe

Check that GitLab isn't reporting any issue.
It use the official GitLab Status (StatusIO page) to check if there no "major" problem with the GitLab infrastructure.

#### Usage

```go
gitlabprobe.NewGitLabProbe(name string)
```

* name: The name of the probe

### HTTPProbe

Check that an HTTP service is running fine.
It will send a GET request to an endpoint and check that the response code is in the 2XX or 3XX class.

#### Usage

```go
httpprobe.NewHTTPProbe(name, endpoint sring, opts HTTPOptions)
```

* name: The name of the probe
* endpoint: Endpoint which should be checked (e.g.: http://google.com)
* opts: General Options

HTTPOptions params:

* Username used for basic auth
* Password used for basic auth
* Checker custom checker that will check the response sent by the server
* ExpectedStatusCode will check for a specific status code
* DialTimeout provide a custom timeout to first byte
* ResponseTimeout provide a custom timeout from first byte to the end of the response

### MongoProbe

Check that a MongoDB database is up and running.

#### Usage

```go
mongoprobe.NewMongoProbe(name, url string)
```

* name: The name of the probe
* url: Url used to check the probe

### NsqProbe

Check that a nsq database is up and running.

#### Usage

```go
nsqprobe.NewNSQProbe(name, host string, port int)
```

* name: The name of the probe
* host: the IP address (or FQDN) of the nsq server
* port: The port on which the nsq server is running

### PhilaeProbe

Check that another service using Philae probe is running and healthy.

#### Usage

```go
philaeprobe.NewPhilaeProbe(name, endpoint string, dialTimeout, responseTimeout int)
```

* name: The name of the probe
* endpoint: The philae endpoint (e.g.: "http://example.com/_health")
* dialTimeout, responseTimeout: see HTTPProbe

### RedisProbe

Check that a Redis server is up and running

#### Usage

```go
redisprobe.NewRedisProbe(name, host, password) string
```

* name: The name of the probe
* host: The Redis host
* password: The password needed to access the database


```go
redisprobe.NewRedisProbeFromURL(name, url string)
```

* name: The name of the probe
* url: The url of the Redis server (e.g.: "redis://:password@example.com")

### PostgreSQLProbe

Check that a PostgreSQL server is up and running

#### Usage

```go
pgsqlprobe.NewPostgreSQLProbe(name, host, password) string
```

* name: The name of the probe
* host: The PostgreSQL host
* password: The password needed to access the database


```go
pgsqlprobe.NewPostgreSQLProbeFromURL(name, url string)
```

* name: The name of the probe
* url: The URL of the PostgreSQL server (e.g.: `postgres://username:password@example.com`)

### MySQLProbe

Check that a MySQL server is up and running

#### Usage

```go
mysqlprobe.NewMySQLProbe(name, host, password) string
```

* name: The name of the probe
* host: The MySQL host
* password: The password needed to access the database


```go
mysqlprobe.NewMySQLProbeFromURL(name, url string)
```

* name: The name of the probe
* url: The URL of the MySQL server (e.g.: `mysql://username:password@example.com`)

### SampleProbe

A probe only used for testing. It will always return the same result

#### Usage

```go
sampleprobe.NewSampleProbe(name string, result bool)
```

* name: The name of the probe
* result: Is the check successful or not

```go
sampleprobe.NewTimedSampleProbe(name string, result bool, time time.Duration)
```

* name: The name of the probe
* result: Is the check successful or not
* time: The time the probe will take before returning a result

### StatusIOProbe

This probe will check that a service using StatusIO is healthy

#### Usage

```go
statusioprobe.NewStatusIOProbe(name, id string)
```

* name: The name of the probe
* id: The StatusIO service id

### SwiftProbe

Check that a Swift host is up and healthy.
It creates a new connection and try to authenticate.

#### Usage

```go
swiftprobe.NewSwiftProbe(name, url, region, tenant, username, password string)
```

* name: The name of the probe
* url: Url of the Swift host
* tenant: The tenant name needed to authenticate
* username: The username needed to authenticate
* password: The password needed to authenticate

### TCPProbe

Check that a TCP server accept connection.

#### Usage

```go
tcpprobe.NewTCPProbe(name, endpoint string, opts TCPOptions)
```

* name: The name of the probe
* endpoint: Endpoint to probe

Options:

* Timeout (default 5s)

## Release a New Version

Bump new version number in:

- `CHANGELOG.md`
- `README.md`

Commit, tag and create a new release:

```sh
version="5.2.2"

git switch --create release/${version}
git add CHANGELOG.md README.md
git commit --message="feat: bump v${version}"
git push --set-upstream origin release/${version}
gh pr create --reviewer=john-scalingo --fill-first
```

The make a PR. Once the PR is merged:

```sh
git pull origin master
git tag ${version}
git push origin master ${version}
gh release create v${version} --generate-notes
```

The title of the release should be the version number and the text of the
release should be the generated notes from GitHub.
