### Transport

To inject request information into the request context, using **transport.MonitorRequest()**
``` go
    monitoredCtx := transport.MonitorRequest(
            ctx.UserContext(),
            transport.MonitorRequestData{
                ClientIP:       ctx.IP(),
                Protocol:       metadata.ProtocolHTTP,
                Method:         ctx.Method(),
                RequestPath:    app.getPath(ctx.Path(), true),
                Hostname:       ctx.Hostname(),
                UserAgent:      string(ctx.Context().UserAgent()),
                ClientTime:     trace.Cts,
                ContentLength:  len(ctx.Request().Body()),
                Request:        req,
                ClientID:       trace.Cid,
                From:           trace.From,
                To:             trace.To,
                Username:       trace.Username,
                RequestHeaders: reqHeaders,
                RemoteHost:     ctx.Context().RemoteAddr().String(),
                XForwardedFor:  strings.Join(ctx.IPs(), ","),
            },
            transport.WithLogger(app.getLogger()), // optional
            transport.WithTracer(app.getTracer()), // optional
        )

```

To get the response with injected request information, using **transport.GetResponse[T]**
**Note**: In golang, context can only pass from parent function to child function, not reverse. So if we wan to get injected request information in response, make sure that we are using the monitored request  
```go
    monitoredCtx := transport.MonitorRequest(
            ctx.UserContext(),
            transport.MonitorRequestData{
                ClientIP:       ctx.IP(),
                Protocol:       metadata.ProtocolHTTP,
                Method:         ctx.Method(),
                RequestPath:    app.getPath(ctx.Path(), true),
                Hostname:       ctx.Hostname(),
                UserAgent:      string(ctx.Context().UserAgent()),
                ClientTime:     trace.Cts,
                ContentLength:  len(ctx.Request().Body()),
                Request:        req,
                ClientID:       trace.Cid,
                From:           trace.From,
                To:             trace.To,
                Username:       trace.Username,
                RequestHeaders: reqHeaders,
                RemoteHost:     ctx.Context().RemoteAddr().String(),
                XForwardedFor:  strings.Join(ctx.IPs(), ","),
            },
            transport.WithLogger(app.getLogger()), // optional
            transport.WithTracer(app.getTracer()), // optional
        )


    // This function will also log the response 
    transport.GetResponse[TResp](
		monitoredCtx,
		transport.WithData(resp), // optional
		transport.WithError(err), // optional
	)
```



