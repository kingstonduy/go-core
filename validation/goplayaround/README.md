#### Usage

```go
    func GetValidator() validation.Validator {
	    return goplayaround.NewGpValidator()
    }
```

```go
    var request interface{}
    err := h.validation.Validate(request)
    if err != nil {
        // handle errors
    }
```