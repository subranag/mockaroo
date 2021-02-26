server {
  # the server and port binding typically localhost:<port>/127.0.0.1:<port>/0.0.0.0:<port>
  listen_addr = "localhost:5000"

  # provide snake oil cert and key path (both are required to start in https mode)
  snake_oil_cert = "/home/subbu/snake_oil_cert/server.crt"
  snake_oil_key  = "/home/subbu/snake_oil_cert/server.key"

  # the file to which the requests will be teed apart from the console
  request_log_path = "/var/tmp/requests.log"

  mock "test_api" {
    request {
      path = "/"
      verb = "GET"
      queries = {
        list = ""
      }
    }

    response {
      headers = {
        Transfer-Encoding = "chunked"
        Content-Type      = "application/xml"
        Date              = "Wed, 26 Oct 2016 22:08:54 GMT"
        x-ms-version      = "2016-05-31"
        Server            = "Windows-Azure-Blob/1.0 Microsoft-HTTPAPI/2.0"
      }

      delay {
        min_millis = 500
        max_millis = 1000
      }

      body = <<EOF
<?xml version="1.0" encoding="utf-8"?>  
<EnumerationResults ServiceEndpoint="https://myaccount.blob.core.windows.net/">  
  <MaxResults>3</MaxResults>  
  <Containers>  
    <Container>  
      <Name>audio</Name>  
      <Properties>  
        <Last-Modified>Wed, 26 Oct 2016 20:39:39 GMT</Last-Modified>  
        <Etag>0x8CACB9BD7C6B1B2</Etag> 
        <PublicAccess>container</PublicAccess> 
      </Properties>  
    </Container>  
    <Container>  
      <Name>images</Name>  
      <Properties>  
        <Last-Modified>Wed, 26 Oct 2016 20:39:39 GMT</Last-Modified>  
        <Etag>0x8CACB9BD7C1EEEC</Etag>  
      </Properties>  
    </Container>  
    <Container>  
      <Name>textfiles</Name>  
      <Properties>  
        <Last-Modified>Wed, 26 Oct 2016 20:39:39 GMT</Last-Modified>  
        <Etag>0x8CACB9BD7BACAC3</Etag>  
      </Properties>  
    </Container>  
  </Containers>  
  <NextMarker>video</NextMarker>  
</EnumerationResults>   
            EOF
    }
  }

  mock "user_request" {
    request {
      path = "/users"
      verb = "GET"
    }

    response {
      headers = {
        Content-Type = "application/json"
      }

      body = <<EOF
            [
                {
                    "name" : "{{.NewUUID}}",
                    "height" : {{.RandomFloat 4.0 8.0}},
                    "age" : {{.RandomInt 21 90}},
                    "user-agent": "{{.Headers.Get "User-Agent"}}",
                    "Method": "{{.Method}}",
                    "RemoteAddr": "{{.RemoteAddr}}"
                },
                {
                    "name" : "{{.NewUUID}}",
                    "height" : {{.RandomFloat 4.0 8.0}},
                    "age" : {{.RandomInt 21 90}}
                },
                {
                    "name" : "{{.NewUUID}}",
                    "height" : {{.RandomFloat 4.0 8.0}},
                    "age" : {{.RandomInt 21 90}}
                }
            ]
            EOF
    }
  }

  mock "test_path" {
    request {
      path = "/test/{operation}"
      verb = "GET"

      // request headers fully support regexp for matching
      headers = {
        "Origin" = ".*"
      }
    }

    response {
      status = 409
      body   = <<EOF
Hello World {{.PathVariable "operation"}}
            EOF
    }
  }

  mock "image" {
    request {
      path = "/asset/image"
      verb = "GET"
    }

    response {
      headers = {
        Content-Type = "image/png"
      }
      file = "/home/subbu/development/workspace/game-assets/fairy-tale-backgrounds/_PNG/1/background.png"
    }
  }
}
