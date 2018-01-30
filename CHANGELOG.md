# Changelog

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
