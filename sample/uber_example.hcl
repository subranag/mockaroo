server {

  listen_addr = "localhost:5000"

  // you can declare several mock sections and give each mock a meaningful name
  mock "get_user" {
    request {
      path = "/user/{userId}"
      verb = "GET"
    }
    response {
      /*
      all captured path variables can be accessed in the template 
      with the expression {{.PathVariable "<varible_key>"}}
      */
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

  mock "get_patient" {
    request {
      path = "/patient/{patientId}"
      verb = "GET"

      // you can match headers as regex as well as plain string 
      // all GET requests that come in with Origin "Clinic .*" will
      // match this reuquest
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
      // match this reuquest
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

  mock "request_body_in_template" {
    request {
      path = "/request/body"
      verb = "POST"
    }
    response {
      body = <<EOF
id is {{index .JsonBody "id"}}
array values are {{index .JsonBody "value" 0}} {{index .JsonBody "value" 1}} {{index .JsonBody "value" 2}}
            EOF
    }
  }

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