# Mock-a-(roo)🦘
comprehensive HTTP/HTTPS interface mocking tool for all your development and testing needs! 
- supports complex request matching (path/query/headers/verb)
- mock response can be fully templated; with request body/query/header params in template context
- config language is *hashicorp HCL*, **human readable with single line and multiline comments**, perfect for documenting your interface
- mock complex HTTP API and send the mockaroo file to client developer for resting, communicate intent clearly, document your API
- run integration tests locally with ease
- add delays to mock real world response times  
- support for HTTPS in case you need it

## Getting Started: Mock HTTP Server in under a minute
- Step 1: download mockaroo binary
- Step 2: add this content to a file `mock.hcl`
```hcl
server {
    listen_addr = "localhost:5000"
    mock "hello_world" {
        request {
            path = "/hello/{userName}"
            verb = "GET"
        }
        response {
            body = <<EOF
            hello world {{.PathVariable "userName"}}
            EOF
        }
    }
}
```
- Step 3: start the mockaroo server from your terminal
```
mockaroo -conf ./mock.hcl
```
- Step 4: curl away 
```
curl "http://localhost:5000/hello/buddy"
```
That's it; read the complete documentation for all the features available for complex mocking

## Why ?
There are several reasons why I always felt the need for something like mockaroo in my local dev box, these are some of them
- I am experimenting with a new API and want to send a live working version to my teammates, simply publish a mock file