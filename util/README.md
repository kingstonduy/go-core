#### Usage

##### Mapping
```go
   	input := map[string]interface{}{
		"name": "John",
		"age":  18,
	}
	err = util.MapStruct(
        input, 
        &balRes,
        util.WithDecodeTimeFormat(time.RFC3339Nano),
        util.WithWeaklyTypedInput(true),    
    )
	if err != nil {
		return nil, err
	}
```