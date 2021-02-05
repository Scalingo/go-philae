# Changelog

## To Be Released

* Various dependencies update following Dependabot activation

## v4.4.2

* Stop using github.com/coreos/etcd, use instead go.etcd.io/etcd/v3

## v4.4.1

* Dependencies: stop using global github.com/Scalingo/go-utils module, use submodule
* Use /v4 suffix in module

## v4.4.0

* [swiftprobe] Make it compatible with API v3 of Openstack Keystone

## v4.3.3

* Add PostgreSQL probe

## v4.3.2

* Add License to open-source project
* Ditch GoConvey for better go 1.12 compat

## v4.3.1

* Hotfix, nil value type is still the type
* Add go.mod/go.sum for go modules

## v4.3.0

* Update Prober to return a correct Error in Prober.Result.Error
* Update Prober to have a default higher timeout (10 seconds) and use context
* TCP probe is now resolving DNS first

## v4.2.0

* Use new github status API
* Fix mongo initializer when there are query parameters

## v4.1.0

* Add `genericprobe` for Pinger (Ping() error)

## v4.0.0

* Migrate to context embedded logger

## v3.1.0

* tcpprobe: add tcpprobe

## v3.0.1

* nsqprobe: use https scheme if tls config is set

## v3.0.0

* nsqprobe: add tls.Config argument to handle https test

## v2.0.1

* Sirupsen -> sirupsen

## v2.0.0

* Add ability to configure logger and verbosity of philaehandler

Breaking Change:

```
philaehandler.NewPhilaeRouter(http.Router, *prober.Prober)
philaehandler.NewHandler(*prober.Prober)

// become

philaehandler.NewPhilaeRouter(http.Router, *prober.Prober, philaehandler.HandlerOpts)
philaehandler.NewHandler(*prober.Prober, philaehandler.HandlerOpts)
```

## v1.0.0

* Initial release (semver version|to use with dep)
