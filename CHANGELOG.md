# Changelog

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
