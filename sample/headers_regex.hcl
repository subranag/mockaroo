server {
  # the server and port binding typically localhost:<port>/127.0.0.1:<port>/0.0.0.0:<port>
  listen_addr = "localhost:7000"

  mock "origin_header_known" {
    request {
      // NOTE: path is same 
      path = "/test/{operation}"
      verb = "POST"

      // NOTE: please include more specific matches in ORDER abive
      // both request paths are for /test/{operation}
      // but this request expectes to be matched with Origin header matching
      headers = {
        "Origin" = ".*"
      }
    }

    response {
      body = <<EOF
Hello World {{.PathVariable "operation"}} I know my Origin {{.Headers.Get "Origin"}}
            EOF
    }
  }

  mock "origin_header_unknown" {
    request {
      // NOTE: path is same 
      path = "/test/{operation}"
      verb = "POST"
    }

    response {
      body = <<EOF
Hello World {{.PathVariable "operation"}} I do not know my origin
            EOF
    }
  }
}
