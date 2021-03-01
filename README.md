# Mock-a-(roo)🦘
comprehensive HTTP/HTTPS interface mocking tool for all your development and testing needs! 
- supports complex request matching (path/query/headers/verb)
- mock response can be fully templated; with request body/query/header params in template context
- config language is *hashicorp HCL*, **human readable with single line and multiline comments**, perfect for documenting your interface
- mock complex HTTP API and send the mockaroo file to client developer for testing; communicate intent clearly; document your API
- run integration tests locally with ease
- add delays to mock real world response times
- full support for *generating fake data* using https://github.com/brianvoe/gofakeit
- support for HTTPS in case you need it

## Getting Started: Mock HTTP Server in under a minute
- Step 1: download mockaroo latest binary [here](https://github.com/subranag/mockaroo/releases) for your target platform and rename the file to `mockaroo` you wil have to do `chmod +x mockaroo` for darwin and linux
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
That's it; read the [complete documentation](https://github.com/subranag/mockaroo/blob/master/DOCS.md) for all the features available for complex mocking

## Why ?
There are several reasons why I always felt the need for something like mockaroo in my local dev box, these are some of them
- I am experimenting with a new API and want to send a *live working version* with use-cases to my teammates, simply publish a mock file
- I am a UX developer who wants to get started on UI work and wants to communicate *my API needs* to the back end developer, with clear documentation
- I wand to publish API documentation that is "ALIVE" simply download mock file and code to it 
- Publish integration test contract that can be run and verified locally in the dev box
- Test my client code against exhaustive fake data that almost looks real 
- I want to test my service locally but *the service I am using does not run locally*
- I can check-in my mock HCL file with code so that other developers can contribute to the use cases and enhance the tests

... and many more, if any of these reasons resonate with you then mockaroo might be a good fit for you 

## Show me the goods 
Just to show that mockaroo is a great fit for mocking needs I have created a samples folder https://github.com/subranag/mockaroo/tree/master/sample which contains several detailed examples I will keep adding to the list in the future, please check it out

## Acknowledgments
This project would have not been possible without these two awesome projects 
- **Gorilla Mux** : awesome router and route matching library https://github.com/gorilla/mux
- **Gofakeit**    : awesome library in golang to create fake data https://github.com/brianvoe/gofakeit
- **Hashicorp Config Language**: really awesome config language which is the basis of terraform https://github.com/hashicorp/hcl

## Documentation 
detailed documentation is available [here](https://github.com/subranag/mockaroo/blob/master/DOCS.md), the [sample](https://github.com/subranag/mockaroo/tree/master/sample) folder contains som really cool examples from top service providers