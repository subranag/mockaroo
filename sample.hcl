server {
    # the server and port binding typically localhost:<port>/127.0.0.1:<port>/0.0.0.0:<port>
    listen_addr = "localhost:5000"
    # the certificate if any typically a self signed cert 
    snake_oil_cert = "/var/tmp/my_snake_oil_cert"

    # the file to which the requests will be teed apart from the console
    request_log_path = "/var/tmp/requests.log"

    mock "user_request" {
        request {
            path = "/users"
            verb = "GET"
        }

        response {
            headers = {
                Content-Type = "application/json"
                Response-Length = "50"
            }

            response_body = <<EOF
            [
                {
                    "name" : "bob",
                    "height" : 5.5,
                    "age" : 40
                },
                {
                    "name" : "Jack",
                    "height" : 5.11,
                    "age" : 42
                },
                {
                    "name" : "Brosnann",
                    "height" : 6.11,
                    "age" : 32
                }
            ]
            EOF
        }
    }

    mock "image_resp" {
        request {
            path = "/"
            verb = "GET"
        }
    }
}