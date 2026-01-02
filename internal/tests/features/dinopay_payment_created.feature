Feature: process DinoPay webhook event PaymentCreated
  DinoPay sends a webhook event of type PaymentCreated

  Background: the dinopay-gateway is up and running
    Given a running dinopay-gateway

  Scenario: the payment is valid and matches the accountNumber of a walletera user
    Given a DinoPay PaymentCreated event:
    """
    data/dinopay_payment_created_event.json
    """
    And  an accounts endpoint to get accounts:
    # the json below is a mockserver expectation
    """json
    {
      "id": "getAccount",
      "httpRequest" : {
        "method": "GET",
        "path": "/accounts",
        "queryStringParameters": {
            "dinopayAccountNumber": "IE12BOFI90000112349876"
        }
      },
      "httpResponse" : {
        "statusCode" : 200,
        "headers" : {
          "content-type" : [ "application/json" ]
        },
        "body": [
          {
            "id": "11111111-2222-3333-4444-555555555555",
            "customerId": "9fd3bc09-99da-4486-950a-11082f5fd966",
            "customerAccountId": "DINO-ACC-001",
            "institutionName": "DinoPay Bank",
            "institutionId": "DNPY-001",
            "currency": "ARS",
            "accountDetails": {
              "accountType": "dinopay",
              "accountHolder": "jane doe",
              "accountNumber": "IE12BOFI90000112349876"
            }
          }
        ]
      },
      "priority" : 0,
      "timeToLive" : {
        "unlimited" : true
      },
      "times" : {
        "unlimited" : true
      }
    }
    """
    And  a payments endpoint to create payments:
    """
    data/payments_create_payments_endpoint_expectation.json
    """
    When the webhook event is received
    Then the dinopay-gateway creates the corresponding payment on the Payments API
    And the dinopay-gateway produces the following log:
    """
    DinoPay event PaymentCreated processed successfully
    """
    And the dinopay-gateway produces the following log:
    """
    Gateway event InboundPaymentReceived processed successfully
    """
