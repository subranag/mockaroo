server {
	listen_addr = "localhost:5001"

    
	mock "create_charge" {
		request {
			path = "/v1/charges"
			verb = "POST"
	    } 
        response {

            headers = {
                Content-Type = "application/json"
            }

		    body = <<EOF
{
  "id": "ch_1IOmoN2eZvKYlo2CVukCCtcR",
  "object": "charge",
  "amount": {{.Form.Get "amount"}},
  "amount_captured": 0,
  "amount_refunded": 0,
  "application": null,
  "application_fee": null,
  "application_fee_amount": null,
  "balance_transaction": "txn_1032HU2eZvKYlo2CEPtcnUvl",
  "billing_details": {
    "address": {
      "city": null,
      "country": null,
      "line1": null,
      "line2": null,
      "postal_code": null,
      "state": null
    },
    "email": null,
    "name": "Jenny Rosen",
    "phone": null
  },
  "calculated_statement_descriptor": null,
  "captured": false,
  "created": 1614270935,
  "currency": "{{.Form.Get "currency"}}",
  "customer": null,
  "description": "{{.Form.Get "description"}}",
  "disputed": false,
  "failure_code": null,
  "failure_message": null,
  "fraud_details": {},
  "invoice": null,
  "livemode": false,
  "metadata": {},
  "on_behalf_of": null,
  "order": null,
  "outcome": null,
  "paid": true,
  "payment_intent": null,
  "payment_method": "card_19yUNL2eZvKYlo2CNGsN6EWH",
  "payment_method_details": {
    "card": {
      "brand": "visa",
      "checks": {
        "address_line1_check": null,
        "address_postal_code_check": null,
        "cvc_check": "unchecked"
      },
      "country": "US",
      "exp_month": 12,
      "exp_year": 2020,
      "fingerprint": "Xt5EWLLDS7FJjR1c",
      "funding": "credit",
      "installments": null,
      "last4": "4242",
      "network": "visa",
      "three_d_secure": null,
      "wallet": null
    },
    "type": "card"
  },
  "receipt_email": null,
  "receipt_number": null,
  "receipt_url": "https://pay.stripe.com/receipts/acct_1032D82eZvKYlo2C/ch_1IOmoN2eZvKYlo2CVukCCtcR/rcpt_J0oR1uwQjlm4YLcYF54frnBgVBO5oC0",
  "refunded": false,
  "refunds": {
    "object": "list",
    "data": [],
    "has_more": false,
    "url": "/v1/charges/ch_1IOmoN2eZvKYlo2CVukCCtcR/refunds"
  },
  "review": null,
  "shipping": null,
  "source_transfer": null,
  "statement_descriptor": null,
  "statement_descriptor_suffix": null,
  "status": "succeeded",
  "transfer_data": null,
  "transfer_group": null,
  "source": "{{.Form.Get "source"}}"
}
EOF
		}
	
    }

    mock "capture_charge" {
        request {
			path = "/v1/charges/{chargeId}/capture"
			verb = "POST"
	    }

        response {

            headers = {
                Content-Type = "application/json"
            }

		    body = <<EOF
{
  "id": "{{.PathVariable "chargeId"}}",
  "object": "charge",
  "amount": 100,
  "amount_captured": 0,
  "amount_refunded": 0,
  "application": null,
  "application_fee": null,
  "application_fee_amount": null,
  "balance_transaction": "txn_1032HU2eZvKYlo2CEPtcnUvl",
  "billing_details": {
    "address": {
      "city": null,
      "country": null,
      "line1": null,
      "line2": null,
      "postal_code": null,
      "state": null
    },
    "email": null,
    "name": "Jenny Rosen",
    "phone": null
  },
  "calculated_statement_descriptor": null,
  "captured": false,
  "created": 1614270935,
  "currency": "usd",
  "customer": null,
  "description": "My First Test Charge (created for API docs)",
  "disputed": false,
  "failure_code": null,
  "failure_message": null,
  "fraud_details": {},
  "invoice": null,
  "livemode": false,
  "metadata": {},
  "on_behalf_of": null,
  "order": null,
  "outcome": null,
  "paid": true,
  "payment_intent": null,
  "payment_method": "card_19yUNL2eZvKYlo2CNGsN6EWH",
  "payment_method_details": {
    "card": {
      "brand": "visa",
      "checks": {
        "address_line1_check": null,
        "address_postal_code_check": null,
        "cvc_check": "unchecked"
      },
      "country": "US",
      "exp_month": 12,
      "exp_year": 2020,
      "fingerprint": "Xt5EWLLDS7FJjR1c",
      "funding": "credit",
      "installments": null,
      "last4": "4242",
      "network": "visa",
      "three_d_secure": null,
      "wallet": null
    },
    "type": "card"
  },
  "receipt_email": null,
  "receipt_number": null,
  "receipt_url": "https://pay.stripe.com/receipts/acct_1032D82eZvKYlo2C/ch_1IOmoN2eZvKYlo2CVukCCtcR/rcpt_J0oR1uwQjlm4YLcYF54frnBgVBO5oC0",
  "refunded": false,
  "refunds": {
    "object": "list",
    "data": [],
    "has_more": false,
    "url": "/v1/charges/ch_1IOmoN2eZvKYlo2CVukCCtcR/refunds"
  },
  "review": null,
  "shipping": null,
  "source_transfer": null,
  "statement_descriptor": null,
  "statement_descriptor_suffix": null,
  "status": "succeeded",
  "transfer_data": null,
  "transfer_group": null
}
EOF
		}
    }

    mock "get_charge" {
		request {
			path = "/v1/charges/{chargeId}"
			verb = "GET"
	    } 
        response {

            headers = {
                Content-Type = "application/json"
            }

		    body = <<EOF
{
  "id": "{{.PathVariable "chargeId"}}",
  "object": "charge",
  "amount": 100,
  "amount_captured": 0,
  "amount_refunded": 0,
  "application": null,
  "application_fee": null,
  "application_fee_amount": null,
  "balance_transaction": "txn_1032HU2eZvKYlo2CEPtcnUvl",
  "billing_details": {
    "address": {
      "city": null,
      "country": null,
      "line1": null,
      "line2": null,
      "postal_code": null,
      "state": null
    },
    "email": null,
    "name": "Jenny Rosen",
    "phone": null
  },
  "calculated_statement_descriptor": null,
  "captured": false,
  "created": 1614270935,
  "currency": "usd",
  "customer": null,
  "description": "My First Test Charge (created for API docs)",
  "disputed": false,
  "failure_code": null,
  "failure_message": null,
  "fraud_details": {},
  "invoice": null,
  "livemode": false,
  "metadata": {},
  "on_behalf_of": null,
  "order": null,
  "outcome": null,
  "paid": true,
  "payment_intent": null,
  "payment_method": "card_19yUNL2eZvKYlo2CNGsN6EWH",
  "payment_method_details": {
    "card": {
      "brand": "visa",
      "checks": {
        "address_line1_check": null,
        "address_postal_code_check": null,
        "cvc_check": "unchecked"
      },
      "country": "US",
      "exp_month": 12,
      "exp_year": 2020,
      "fingerprint": "Xt5EWLLDS7FJjR1c",
      "funding": "credit",
      "installments": null,
      "last4": "4242",
      "network": "visa",
      "three_d_secure": null,
      "wallet": null
    },
    "type": "card"
  },
  "receipt_email": null,
  "receipt_number": null,
  "receipt_url": "https://pay.stripe.com/receipts/acct_1032D82eZvKYlo2C/ch_1IOmoN2eZvKYlo2CVukCCtcR/rcpt_J0oR1uwQjlm4YLcYF54frnBgVBO5oC0",
  "refunded": false,
  "refunds": {
    "object": "list",
    "data": [],
    "has_more": false,
    "url": "/v1/charges/ch_1IOmoN2eZvKYlo2CVukCCtcR/refunds"
  },
  "review": null,
  "shipping": null,
  "source_transfer": null,
  "statement_descriptor": null,
  "statement_descriptor_suffix": null,
  "status": "succeeded",
  "transfer_data": null,
  "transfer_group": null
}
EOF
		}
	
    }

    
    mock "list_charge" {
		request {
			path = "/v1/charges"
			verb = "GET"
	    } 
        response {

            headers = {
                Content-Type = "application/json"
            }

		    body = <<EOF
{
  "object": "list",
  "url": "/v1/charges",
  "has_more": false,
  "data": [
    {
      "id": "ch_1IOmoN2eZvKYlo2CVukCCtcR",
      "object": "charge",
      "amount": 100,
      "amount_captured": 0,
      "amount_refunded": 0,
      "application": null,
      "application_fee": null,
      "application_fee_amount": null,
      "balance_transaction": "txn_1032HU2eZvKYlo2CEPtcnUvl",
      "billing_details": {
        "address": {
          "city": null,
          "country": null,
          "line1": null,
          "line2": null,
          "postal_code": null,
          "state": null
        },
        "email": null,
        "name": "Jenny Rosen",
        "phone": null
      },
      "calculated_statement_descriptor": null,
      "captured": false,
      "created": 1614270935,
      "currency": "usd",
      "customer": null,
      "description": "My First Test Charge (created for API docs)",
      "disputed": false,
      "failure_code": null,
      "failure_message": null,
      "fraud_details": {},
      "invoice": null,
      "livemode": false,
      "metadata": {},
      "on_behalf_of": null,
      "order": null,
      "outcome": null,
      "paid": true,
      "payment_intent": null,
      "payment_method": "card_19yUNL2eZvKYlo2CNGsN6EWH",
      "payment_method_details": {
        "card": {
          "brand": "visa",
          "checks": {
            "address_line1_check": null,
            "address_postal_code_check": null,
            "cvc_check": "unchecked"
          },
          "country": "US",
          "exp_month": 12,
          "exp_year": 2020,
          "fingerprint": "Xt5EWLLDS7FJjR1c",
          "funding": "credit",
          "installments": null,
          "last4": "4242",
          "network": "visa",
          "three_d_secure": null,
          "wallet": null
        },
        "type": "card"
      },
      "receipt_email": null,
      "receipt_number": null,
      "receipt_url": "https://pay.stripe.com/receipts/acct_1032D82eZvKYlo2C/ch_1IOmoN2eZvKYlo2CVukCCtcR/rcpt_J0oR1uwQjlm4YLcYF54frnBgVBO5oC0",
      "refunded": false,
      "refunds": {
        "object": "list",
        "data": [],
        "has_more": false,
        "url": "/v1/charges/ch_1IOmoN2eZvKYlo2CVukCCtcR/refunds"
      },
      "review": null,
      "shipping": null,
      "source_transfer": null,
      "statement_descriptor": null,
      "statement_descriptor_suffix": null,
      "status": "succeeded",
      "transfer_data": null,
      "transfer_group": null
    }
  ]
}
EOF
		}
	
    }
}