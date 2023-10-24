# SSE Example

Example for https://github.com/bigskysoftware/htmx/pull/1794

# Start

run:

```
go build && ./example
```

Afterwards open http://localhost:8080/


You see a changing message like:

```
hello 669. 
```

But if you open the sse-url with curl:

```
curl http://localhost:8080/sse
```

You see, that the full message looks like this:

```
data: hello 669. <a href="https://github.com/bigskysoftware/htmx/pull/1794">Here</a> is a fantastic pr.
```


If you change the `main.go` file in line 36 to 41 to use the last stable version of htmx, then you can see, that it works in the browser as expected.
