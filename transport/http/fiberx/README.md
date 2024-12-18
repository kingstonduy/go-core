```go
	fConfig := fiberx.DefaultFiberConfig
	fConfig.BodyLimit = 1

	fiberApp := fiberx.NewFiberApp(
		fiberx.WithBasePath(cfg.ServerConfig.BasePath),
		fiberx.WithSwaggerPath("/swagger/*"),
		fiberx.WithMetricsPath("/metrics"),
		fiberx.WithServiceName(cfg.ServerConfig.ApplicationName),
		fiberx.WithFiberConfig(fConfig),
	)
```
