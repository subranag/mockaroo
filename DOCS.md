# Mock-a-(roo)🦘 Docs
This page contains detailed documentation on how to use mockaroo to create complex HTTP(S) mocking solutions, click on each section to see how to proceed

  * [All Examples](#all-examples)
  * [Starting Mockaroo](#starting-mockaroo)
  * [The Server Section](#the-server-section)
  * [The Mock Blocks](#the-mock-blocks)
  * [Matching Query Params](#matching-query-params)
  * [Matching Headers](#matching-headers)
  * [Accessing Request Body](#accessing-request-body)
  * [Capturing Path Variables](#capturing-path-variables)
  * [Template Execution Response](#template-execution-response)
  * [File in Response](#file-in-response)
  * [The Complete Example](#the-complete-example)

## All Examples
all examples that have been listed below have been coded into a single big mock file which can be found [here](https://github.com/subranag/mockaroo/blob/master/sample/uber_example.hcl)

the required information is available in both the places here in the docs and the **uber big file** containing all examples

## Starting Mockaroo
Starting mockaroo is simple [download](https://github.com/subranag/mockaroo/releases) mockaroo binary for your target platform and drop it in you PATH if required, you may have to do `chmod +x` for linux and darwin and rename the file if you like, then simply start mockaroo by pointing it to you mock `hcl` file e.g.
```
mockaroo -conf ./<path_to_your_hcl_file>
```
the mock files are written in HCL https://www.terraform.io/docs/language/syntax/configuration.html , **HCL is a superb configuration language for clear configuration and readability**
> ⚠️**NOTE**: the file extension should be HCL otherwise you might get an error


## The Server Section 
the server section in the mock HCL deals with specifying HTTP(S) server related configuration, see sample file with documentation as well in HCL 
```hcl
server {
  /* 
    the server and port binding typically of the form
    localhost:<port> OR 127.0.0.1:<port> OR 0.0.0.0:<port> 
    you can also bind it to a specific ip and a port that is available 
  */   
  listen_addr = "localhost:5000"

  /* 
    since this will be running locally snake oil is used in the name 
    NOTHING is stopping you from using real cert and key
  */
  // e.g. "/home/subbu/snake_oil_cert/server.crt"
  snake_oil_cert = "/<path>/<cert_file_name>.crt"

  // e.g. "/home/subbu/snake_oil_cert/server.key" 
  snake_oil_key  = "/<path>/<key_file_name>.key" 

  /* 
    you can also provide a request_log_path where requests will be logged
    this is OPTIONAL, but recommended to see what kind of requests are coming 
    along
  */
  request_log_path = "/var/tmp/requests.log"
  ...
```
> ⚠️**NOTE**: the server will start in HTTPS mode if and only if BOTH snake_oil_cert and snake_oil_key are present

## The Mock Blocks
after the server section is declared in the HCL file you need declare *one or more* mock blocks in the mockaroo file 

there can be several mock blocks, typically you can declare all the mocks required for a single use-case or scenario in a single mock HCL file, see example below where two mocks are declared 

```hcl
server {

  listen_addr = "localhost:5000"

  // you can declare several mock sections and give each mock a meaningful name
  mock "get_user" {
    request {
      path = "/user/{userId}"
      verb = "GET"
    }
    response {
      body = <<EOF
            user id from GET {{.PathVariable "userId"}}
            EOF
    }
  }

  // another mock with the same HTTP path but different verb "POST"
  mock "post_user" {
    request {
      path = "/user/{userId}"
      verb = "POST"
    }
    response {
      # NOTE: response status is 201 created
      status = 201

      body = <<EOF
            user id from POST {{.PathVariable "userId"}}
            EOF
    }
  }
}
```
> 🚨 **NOTE**: there is order to matching mocks , mocks should be declared in decreasing order of *specificity in matching* from the *MOST SPECIFIC MATCH* to *LEAST SPECIFIC MATCH* otherwise matching may not work as expected

## Matching Query Params
you can specify query parameters to match query params on incoming request, please look at the example below to see how to do it, you can match query params to value 

```hcl
server {

  listen_addr = "localhost:5000"

  mock "get_beer_lager" {
    request {
      path = "/beer"
      verb = "GET"

      // you can match criteria on query params 
      // you can specify multiple query params
      /*
      queries = {
        a = "b"
        c = "b" 
        ...
      }
      */
      queries = {
        type = "lager"
      }
    }
    response {
      body = <<EOF
            my beer is {{.Form.Get "type"}}, brand is {{.Fake.BeerMalt}}
            EOF
    }
  }

  mock "get_beer_ipa" {
    request {
      path = "/beer"
      verb = "GET"

      // here we match ipa
      queries = {
        type = "ipa"
      }
    }
    response {

      body = <<EOF
            my beer is {{.Form.Get "type"}}, brand is {{.Fake.BeerHop}}
            EOF
    }
  }
}
```
> ℹ️ **NOTE**: currently regex matching is not available for query params but will be added in the future

## Matching Headers
the headers in the HTTP request can be matched as well, take a look at the example below to see how header matching works, documentation is in the mock itself

```hcl
server {

  listen_addr = "localhost:5000"

    mock "get_patient" {
    request {
      path = "/patient/{patientId}"
      verb = "GET"

      // you can match headers as regex as well as plain string 
      // all GET requests that come in with Origin "Clinic .*" will
      // match this request
      // you can specify several headers as well
      headers = {
        "Origin" = "Clinic .*"
      }
    }
    response {
      // you can specify headers for response as well you can specify 
      // any number of headers  
      headers = {
        Content-Type = "application/json"
      }

      body = <<EOF
{
    "patient_id": "{{.PathVariable "patientId"}}",
    "patient_name": "{{.Fake.Name}}",
    "email": "{{.Fake.Email}}",
    "street": "{{.Fake.Street}}",
    "state": "{{.Fake.State}}",
    "ssn": "{{.Fake.SSN}}",
    "phone": "{{.Fake.PhoneFormatted}}",
    "source": "Clinic/{{.Headers.Get "Origin"}}"
}
            EOF
    }
  }

  mock "get_patient_hospital" {
    request {
      path = "/patient/{patientId}"
      verb = "GET"

      // you can match headers as regex as well as plain string 
      // all GET requests that come in with Origin "Clinic .*" will
      // match this request
      // you can specify several headers as well
      headers = {
        "Origin" = "Hospital .*"
      }
    }
    response {
      // you can specify headers for response as well you can specify 
      // any number of headers  
      headers = {
        Content-Type = "application/json"
      }

      body = <<EOF
{
    "patient_id": "{{.PathVariable "patientId"}}",
    "patient_name": "{{.Fake.Name}}",
    "email": "{{.Fake.Email}}",
    "street": "{{.Fake.Street}}",
    "state": "{{.Fake.State}}",
    "ssn": "{{.Fake.SSN}}",
    "phone": "{{.Fake.PhoneFormatted}}",
    "source": "Hospital/{{.Headers.Get "Origin"}}"
}
            EOF
    }
  }
}
```

## Accessing Request Body
if the RAW quest body can be parsed as JSON the entire parsed JSON is available to the template context when sending back response let us look at example below

```hcl
server {

  listen_addr = "localhost:5000"

  mock "request_body_in_template" {
    request {
      path = "/request/body"
      verb = "POST"
    }

    response {
      // if the request body can be parsed as JSON the entire request is available 
      // as dictionary an array references you can access them in templates using 
      // the GOLANG template index function see example below
      body = <<EOF
id is {{index .JsonBody "id"}}
array values are {{index .JsonBody "value" 0}} {{index .JsonBody "value" 1}} {{index .JsonBody "value" 2}}
            EOF
    }
  }

}
```

After starting the server with the above configuration you can verify the output using cURL commands shown below
```
curl -X POST -H "Content-Type: application/json" -d'{"id":"a", "value":[1, 2, 3]}' "http://localhost:5000/request/body"
```
expected output should be 
```
id is a
array values are 1 2 3
```
you can write some very powerful mocks having the entire request JSON in the template context

## Capturing Path Variables
path variables can be captured and are available in context during template execution, if you do not want to specify a path variable you can also specify a `*` which will match any text in that path component please checkout the example below 

```hcl
server {

  listen_addr = "localhost:5000"
  
  mock "capture_path_variables" {
    // you can capture path variables by just naming them 
    // if you do not want them stored in a named variable then just mention *
    // the variable will still be captured and be present in the template context
    // the path variable name will of the form pvarN where N is the 1 based index from 
    // the start of the path components see the repose template for usage
    request {
      path = "/path/{a}/{b}/*/{d}"
      verb = "GET"
    }
    response {
      body = <<EOF
      the request was for path/{{.PathVariable "a"}}/{{.PathVariable "b"}}/{{.PathVariable "pvar4"}}/{{.PathVariable "d"}}
            EOF
    }
  }
}
```
Now fire this cURL request
```
curl "http://localhost:5000/test/this/path/correctly"
```
the response should be 
```
      the request was for path/test/this/path/correctly
```

## Template Execution Response
The power of golang template engine is available for the response rending, the complete reference for how go template engine works is outside the scope of this document, but a very comprehensive tutorial can be found [here](https://learn.hashicorp.com/tutorials/nomad/go-template-syntax?in=nomad/templates)

The custom functions available during the template executions are listed here 

| Template Call | Description |
|---------------| ------------|
| `{{.Method}}` | Request method GET/PUT/POST etc|
| `{{.Protocol}}` | HTTP/1.0 HTTP/1.1 etc|
| `{{.Host}}` | host from which the request came|
| `{{.RemoteAddr}}` | remote address|
| `{{.Headers.Get "key"}}` | get the value of request headers|
| `{{.Form.Get "key"}}` | form contains all url query params and POST form data|
| `{{.PathVars "key"}}` | this template variable contains the key value map of all path variables|
| `{{.Fake.<FakeFunction>}}` | using the Fake context you can call all fake functions on gofakeit list of all functions [here](https://github.com/brianvoe/gofakeit#functions) e.g. `{{.Fake.PhoneFormatted}}`|
| `{{.PathVariable "key"}}` | same as PathVars gets the value of path variable captured|
| `{{.RandomInt min max}}` | generated a random int in the interval [min, max) e.g. `{{.RandomInt 5 10}}`|
| `{{.RandomFloat min max}}` | generated a random int in the interval [min, max) e.g. `{{.RandomFloat 1.1 2.2}}`|

> ⚠️**NOTE**: all random data generation is stable i.e. they will same random values in sequence, the seed prime for generation is `2011`

## File In Response
you can include a file path as response the contents of the file will be sent as response this combined with the right MIME type can allow you to send binary response for a mock see the example below

```hcl
server {

  listen_addr = "localhost:5000"

  mock "file_request" {
    request {
      path = "/etc/passwd"
      verb = "GET"
    }
    response {
      headers = {
        Content-Type = "text/plain"
      }
      file = "/etc/passwd"
    }
  }
}
```
Now if you run the cURL command 
```
curl "http://localhost:5000/etc/passwd"
```
you should see you passwd file

## The Complete Example
all of the above examples have been tested and have been dumped into a single big uber example file with all the relevant documentation please take a look [here](https://github.com/subranag/mockaroo/blob/master/sample/uber_example.hcl)

# Happy Mocking!!!!

