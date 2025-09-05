# Changelog

## To Be Released

## v5.2.3

* Various Dependabot upgrades

## v5.2.2

* Various Dependabot upgrades

## v5.2.1

* fix(elasticsearch) Fix unable to use System CA certificate pool

## v5.2.0

* feat(probes): add support for elasticsearch

## v5.1.0

* feat(sampleprobe): enable setting a custom error

## v5.0.0

* feat(probe-interface): add support for context in the Check method
* chore(deps): bump go.etcd.io/etcd/client/v2 from 2.305.4 to 2.305.6
* chore(deps): bump github.com/stretchr/testify from 1.8.0 to 1.8.1
* chore(deps): bump github.com/lib/pq from 1.10.6 to 1.10.7
* chore(deps): bump github.com/go-sql-driver/mysql from 1.6.0 to 1.7.0

## v4.4.7

* chore(deps): bump github.com/Scalingo/go-utils/logger from 1.1.1 to 1.2.0
* chore(deps): bump github.com/stretchr/testify from 1.7.1 to 1.8.0
* chore(deps): bump github.com/sirupsen/logrus from 1.8.1 to 1.9.0
* chore(deps): bump github.com/fsouza/go-dockerclient from 1.7.11 to 1.8.3
* chore(deps): bump github.com/ncw/swift 1.0.53 to github.com/ncw/swift/v2 2.0.1

## v4.4.6

* feat(mysqlprobe): implement MySQL probe

## v4.4.5

* feat(prober): change the exposed CheckOneProbe to make it usable outside of this package
* feat(probe check): implement ErrProbeNotFound

## v4.4.4

* chore(go): use go 1.17
* Bump github.com/lib/pq from 1.10.3 to 1.10.6
* Bump github.com/fsouza/go-dockerclient from 1.7.3 to 1.7.11
* Bump go.etcd.io/etcd/client/v2 from 2.305.0 to 2.305.1
* Bump github.com/Scalingo/go-utils/logger from 1.0.0 to 1.1.0
* chore(deps): bump github.com/stretchr/testify from 1.7.0 to 1.7.1

## v4.4.3

* Various dependencies update following Dependabot activation
* Bump github.com/sirupsen/logrus from 1.7.0 to 1.8.1
* Bump github.com/lib/pq from 1.10.0 to 1.10.2
* Bump github.com/fsouza/go-dockerclient from 1.7.2 to 1.7.3
* Bump go.etcd.io/etcd/client to 2.305.0 (Go modules version)

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
