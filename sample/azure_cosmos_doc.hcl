server {
  listen_addr = "localhost:5002"

  mock "query_document" {

    request {
      path = "/dbs/{databaseId}/colls/{collectionId}/docs"
      verb = "POST"

      headers = {
        x-ms-documentdb-isquery = ".*"
      }
    }

    response {

      # NOTE: the response from the server is actually 201 for create
      status = 201

      delay {
        min_millis = 200
        max_millis = 400
      }

      headers = {
        Cache-Control                 = "no-store, no-cache"
        Pragma                        = "no-cache"
        Transfer-Encoding             = "chunked"
        Content-Type                  = "application/json"
        Server                        = "Microsoft-HTTPAPI/2.0"
        Strict-Transport-Security     = "max-age=31536000"
        x-ms-last-state-change-utc    = "Fri, 25 Mar 2016 22:39:02.501 GMT"
        etag                          = "00003200-0000-0000-0000-56f9e84d0000"
        x-ms-resource-quota           = "documentSize=10240;documentsSize=10485760;collectionSize=10485760;"
        x-ms-resource-usage           = "documentSize=0;documentsSize=1;collectionSize=1;"
        x-ms-schemaversion            = "1.1"
        x-ms-alt-content-path         = "dbs/testdb/colls/testcoll"
        x-ms-quorum-acked-lsn         = "602"
        x-ms-current-write-quorum     = "3"
        x-ms-current-replica-set-size = "4"
        x-ms-request-charge           = "12.38"
        x-ms-serviceversion           = "version=1.6.52.5"
        x-ms-activity-id              = "856acd38-320d-47df-ab6f-9761bb987668"
        x-ms-session-token            = "0:603"
        Set-Cookie                    = "x-ms-session-token#0=603; Domain=querydemo.documents.azure.com; Path=/dbs/1KtjAA==/colls/1KtjAImkcgw="
        Set-Cookie                    = "x-ms-session-token=603; Domain=querydemo.documents.azure.com; Path=/dbs/1KtjAA==/colls/1KtjAImkcgw="
        x-ms-gatewayversion           = "version=1.6.52.5"
        Date                          = "Tue, 29 Mar 2016 02:28:30 GMT"
      }

      body = <<EOF
{  
  "id": "AndersenFamily",  
  "LastName": "Andersen",  
  "Parents": [  
    {  
      "FamilyName": null,  
      "FirstName": "Thomas"  
    },  
    {  
      "FamilyName": null,  
      "FirstName": "Mary Kay"  
    }  
  ],  
  "Children": [  
    {  
      "FamilyName": null,  
      "FirstName": "Henriette Thaulow",  
      "Gender": "female",  
      "Grade": 5,  
      "Pets": [  
        {  
          "GivenName": "Fluffy"  
        }  
      ]  
    }  
  ],  
  "Address": {  
    "State": "WA",  
    "County": "King",  
    "City": "Seattle"  
  },  
  "IsRegistered": true,  
  "_rid": "1KtjAImkcgwBAAAAAAAAAA==",  
  "_self": "dbs/1KtjAA==/colls/1KtjAImkcgw=/docs/1KtjAImkcgwBAAAAAAAAAA==/",  
  "_etag": "\"00003200-0000-0000-0000-56f9e84d0000\"",  
  "_ts": 1459218509,  
  "_attachments": "attachments/"  
}
EOF
    }

  }


  mock "create_document" {
    request {
      path = "/dbs/{databaseId}/colls/{collectionId}/docs"
      verb = "POST"
    }
    response {

      # NOTE: the response from the server is actually 201 for create
      status = 201

      delay {
        min_millis = 800
        max_millis = 1200
      }

      headers = {
        Cache-Control                 = "no-store, no-cache"
        Pragma                        = "no-cache"
        Transfer-Encoding             = "chunked"
        Content-Type                  = "application/json"
        Server                        = "Microsoft-HTTPAPI/2.0"
        Strict-Transport-Security     = "max-age=31536000"
        x-ms-last-state-change-utc    = "Fri, 25 Mar 2016 22:39:02.501 GMT"
        etag                          = "00003200-0000-0000-0000-56f9e84d0000"
        x-ms-resource-quota           = "documentSize=10240;documentsSize=10485760;collectionSize=10485760;"
        x-ms-resource-usage           = "documentSize=0;documentsSize=1;collectionSize=1;"
        x-ms-schemaversion            = "1.1"
        x-ms-alt-content-path         = "dbs/testdb/colls/testcoll"
        x-ms-quorum-acked-lsn         = "602"
        x-ms-current-write-quorum     = "3"
        x-ms-current-replica-set-size = "4"
        x-ms-request-charge           = "12.38"
        x-ms-serviceversion           = "version=1.6.52.5"
        x-ms-activity-id              = "856acd38-320d-47df-ab6f-9761bb987668"
        x-ms-session-token            = "0:603"
      }

      body = <<EOF
{  
  "id": "AndersenFamily",  
  "LastName": "Andersen",  
  "Parents": [  
    {  
      "FamilyName": null,  
      "FirstName": "Thomas"  
    },  
    {  
      "FamilyName": null,  
      "FirstName": "Mary Kay"  
    }  
  ],  
  "Children": [  
    {  
      "FamilyName": null,  
      "FirstName": "Henriette Thaulow",  
      "Gender": "female",  
      "Grade": 5,  
      "Pets": [  
        {  
          "GivenName": "Fluffy"  
        }  
      ]  
    }  
  ],  
  "Address": {  
    "State": "WA",  
    "County": "King",  
    "City": "Seattle"  
  },  
  "IsRegistered": true,  
  "_rid": "1KtjAImkcgwBAAAAAAAAAA==",  
  "_self": "dbs/1KtjAA==/colls/1KtjAImkcgw=/docs/1KtjAImkcgwBAAAAAAAAAA==/",  
  "_etag": "\"00003200-0000-0000-0000-56f9e84d0000\"",  
  "_ts": 1459218509,  
  "_attachments": "attachments/"  
}
EOF
    }

  }

  mock "list_docs" {
    request {
      path = "/dbs/{databaseId}/colls/{collectionId}/docs"
      verb = "GET"
    }

    response {

      delay {
        min_millis = 200
        max_millis = 500
      }

      headers = {
        Cache-Control              = "no-store, no-cache"
        Pragma                     = "no-cache"
        Transfer-Encoding          = "chunked"
        Content-Type               = "application/json"
        Content-Location           = "https://querydemo.documents.azure.com/dbs/testdb/colls/testcoll/docs"
        Server                     = "Microsoft-HTTPAPI/2.0"
        Strict-Transport-Security  = "max-age=31536000"
        x-ms-last-state-change-utc = "Sun, 27 Mar 2016 22:39:13.369 GMT"
        x-ms-resource-quota        = "documentSize=10240;documentsSize=10485760;collectionSize=10485760;"
        x-ms-resource-usage        = "documentSize=0;documentsSize=2;collectionSize=2;"
        x-ms-item-count            = "2"
        x-ms-schemaversion         = "1.1"
        x-ms-alt-content-path      = "dbs/testdb/colls/testcoll"
        x-ms-content-path          = "d9RzAJRFKgw="
        x-ms-request-charge        = "1"
        x-ms-serviceversion        = "version=1.6.52.5"
        x-ms-activity-id           = "46e2e9a5-4917-4ff6-9be5-6f206c38bb6b"
        x-ms-session-token         = "0:772"
        Set-Cookie                 = "x-ms-session-token#0=772; Domain=querydemo.documents.azure.com; Path=/dbs/testdb/colls/testcoll"
        Set-Cookie                 = "x-ms-session-token=772; Domain=querydemo.documents.azure.com; Path=/dbs/testdb/colls/testcoll"
        x-ms-gatewayversion        = "version=1.6.52.5"
        Date                       = "Tue, 29 Mar 2016 02:03:07 GMT"
      }

      body = <<EOF
{  
  "_rid": "d9RzAJRFKgw=",  
  "Documents": [  
    {  
      "id": "SalesOrder1",  
      "ponumber": "PO18009186470",  
      "OrderDate": "2005-07-01T00:00:00",  
      "ShippedDate": "0001-01-01T00:00:00",  
      "AccountNumber": "Account1",  
      "SubTotal": 419.4589,  
      "TaxAmount": 12.5838,  
      "Freight": 472.3108,  
      "TotalDue": 985.018,  
      "Items": [  
        {  
          "OrderQty": 1,  
          "ProductId": 760,  
          "UnitPrice": 419.4589,  
          "LineTotal": 419.4589  
        }  
      ],  
      "_rid": "d9RzAJRFKgwBAAAAAAAAAA==",  
      "_self": "dbs/d9RzAA==/colls/d9RzAJRFKgw=/docs/d9RzAJRFKgwBAAAAAAAAAA==/",  
      "_etag": "\"0000d986-0000-0000-0000-56f9e25b0000\"",  
      "_ts": 1459216987,  
      "_attachments": "attachments/"  
    },  
    {  
      "id": "SalesOrder2",  
      "ponumber": "PO15428132599",  
      "OrderDate": "2005-07-01T00:00:00",  
      "DueDate": "2005-07-13T00:00:00",  
      "ShippedDate": "2005-07-08T00:00:00",  
      "AccountNumber": "Account2",  
      "SubTotal": 6107.0820,  
      "TaxAmt": 586.1203,  
      "Freight": 183.1626,  
      "TotalDue": 4893.3929,  
      "DiscountAmt": 1982.872,  
      "Items": [  
        {  
          "OrderQty": 3,  
          "ProductCode": "A-123",  
          "ProductName": "Product 1",  
          "CurrencySymbol": "$",  
          "CurrencyCode": "USD",  
          "UnitPrice": 17.1,  
          "LineTotal": 5.7  
        }  
      ],  
      "_rid": "d9RzAJRFKgwCAAAAAAAAAA==",  
      "_self": "dbs/d9RzAA==/colls/d9RzAJRFKgw=/docs/d9RzAJRFKgwCAAAAAAAAAA==/",  
      "_etag": "\"0000da86-0000-0000-0000-56f9e25b0000\"",  
      "_ts": 1459216987,  
      "_attachments": "attachments/"  
    }  
  ],  
  "_count": 2  
} 
EOF
    }
  }

  mock "get_document" {
    request {
      path = "/dbs/{databaseId}/colls/{collectionId}/docs/{docId}"
      verb = "GET"
    }
    response {

      delay {
        min_millis = 200
        max_millis = 500
      }

      headers = {
        Cache-Control              = "no-store, no-cache"
        Pragma                     = "no-cache"
        Transfer-Encoding          = "chunked"
        Content-Type               = "application/json"
        Content-Location           = "https://querydemo.documents.azure.com/dbs/testdb/colls/testcoll/docs/SalesOrder1"
        Server                     = "Microsoft-HTTPAPI/2.0"
        Strict-Transport-Security  = "max-age=31536000"
        x-ms-last-state-change-utc = "Mon, 28 Mar 2016 14:47:03.949 GMT"
        etag                       = "0000d986-0000-0000-0000-56f9e25b0000"
        x-ms-resource-quota        = "documentSize=10240;documentsSize=10485760;collectionSize=10485760;"
        x-ms-resource-usage        = "documentSize=0;documentsSize=2;collectionSize=2;"
        x-ms-schemaversion         = "1.1"
        x-ms-alt-content-path      = "dbs/testdb/colls/testcoll"
        x-ms-content-path          = "d9RzAJRFKgw="
        x-ms-request-charge        = "1"
        x-ms-serviceversion        = "version=1.6.52.5"
        x-ms-activity-id           = "c22bc349-2c02-4b80-81b9-a2d758c92902"
        x-ms-session-token         = "0:772"
        Set-Cookie                 = "x-ms-session-token#0=772; Domain=querydemo.documents.azure.com; Path=/dbs/testdb/colls/testcoll"
        Set-Cookie                 = "x-ms-session-token=772; Domain=querydemo.documents.azure.com; Path=/dbs/testdb/colls/testcoll"
        x-ms-gatewayversion        = "version=1.6.52.5"
        Date                       = "Tue, 29 Mar 2016 02:03:06 GMT"
      }

      body = <<EOF
{  
  "id": "{{.PathVariable "docId"}}",  
  "ponumber": "PO18009186470",  
  "OrderDate": "2005-07-01T00:00:00",  
  "ShippedDate": "0001-01-01T00:00:00",  
  "AccountNumber": "Account1",  
  "SubTotal": 419.4589,  
  "TaxAmount": 12.5838,  
  "Freight": 472.3108,  
  "TotalDue": 985.018,  
  "Items": [  
    {  
      "OrderQty": 1,  
      "ProductId": 760,  
      "UnitPrice": 419.4589,  
      "LineTotal": 419.4589  
    }  
  ],  
  "_rid": "d9RzAJRFKgwBAAAAAAAAAA==",  
  "_self": "dbs/d9RzAA==/colls/d9RzAJRFKgw=/docs/d9RzAJRFKgwBAAAAAAAAAA==/",  
  "_etag": "\"0000d986-0000-0000-0000-56f9e25b0000\"",  
  "_ts": 1459216987,  
  "_attachments": "attachments/"  
}
EOF
    }

  }
}
