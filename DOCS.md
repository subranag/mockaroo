# Mock-a-(roo)🦘 Docs
This page contains detailed documentation on how to use mockaroo to create complex HTTP(S) mocking solutions, click on each section to see how to proceed

## Starting Mockaroo
Starting mockaroo is simple download mockaroo binary for your target platform and drop it in you PATH if required, then simply start mockaroo by pointing it to you mock `hcl` file e.g.
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

```
server {
    ...
    // you can declare several mock sections and give each mock a meaningful name
    mock "get_user" {
        request {
            path = "/user/{userId}"
            verb = "GET"
        }
        response {
            body = <<EOF
            user id {{.PathVariable "userId"}}
            EOF
        }
    }

    // another mock with the same HTTP path but different verb "GET"
    mock "post_user" {
        request {
            path = "/user/{userId}"
            verb = "GET"
        }
        response {
            # NOTE: response code is 201 created
            code = 201

            body = <<EOF
            user id {{.PathVariable "userId"}}
            EOF
        }
    }
    ...
}
```
> 🚨 **NOTE**: there is order to matching mocks , mock should be declared in order from *MOST SPECIFIC MATCH* to *LEAST SPECIFIC MATCH* in descending order: otherwise you might have wrong matching (please see relevant section)