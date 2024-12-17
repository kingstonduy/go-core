#### Usage

``` go

func GetConfigure() config.Configure {
	c, err := NewViperConfig(
		config.WithDefaults(map[string]interface{}{
			"DEFAULT_KEY": "DEFAULT_VALUE",
		}),
		config.WithTagName("config"),
		config.WithAutomaticEnv(true),
		config.WithConfigFile("sample.env"),
	)

	if err != nil {
		panic(err)
	}
	
	return cfg
}

```

```go
cfg, err := GetConfigure()
if err != nil {
    panic(err)
}

username := cfg.GetString("DB_YUGABYTE_USER")
password := cfg.GetString("DB_YUGABYTE_PASSWORD")
host := cfg.GetString("DB_YUGABYTE_HOST")
port := cfg.GetInt("DB_YUGABYTE_PORT")

```

##### Overrides env file by ENV_FILE key 
``` go
export GO_ENV_FILE=env_file && go run main.go
```