Feature: process DinoPay webhook event PaymentCreated
  DinoPay sends a webhook event of type PaymentCreated

  Background: the dinopay-gateway is up and running
    Given a running dinopay-gateway

  Scenario: the payment is valid and matches the accountNumber of a walletera user
    Given a PaymentCreated event:
    """json
    {
      "id": "647f9176-466a-4d8c-b027-d53b4da77d4d",
      "type": "PaymentCreated",
      "time": "2023-07-07T19:31:11.123Z",
      "data": {
        "id": "bb17667e-daac-41f6-ada3-2c22f24caf22",
        "amount": 100,
        "currency": "USD",
        "sourceAccount": {
          "accountHolder": "john doe",
          "accountNumber": "IE12BOFI90000112345678"
        },
        "destinationAccount": {
          "accountHolder": "jane doe",
          "accountNumber": "IE12BOFI90000112349876"
        },
        "createdAt": "2023-07-07T19:31:11Z",
        "updatedAt": "2023-07-07T19:31:11Z"
      }
    }
    """
    And  a payments endpoint to create deposits:
    # the json below is a mockserver expectation
    """json
    {
      "id": "postDepositSucceed",
      "httpRequest" : {
        "method": "POST",
        "path": "/deposits",
        "body": {
            "type": "JSON",
            "json": {
              "id": "${json-unit.any-string}",
              "amount": 100,
              "currency": "USD",
              "customerId": "9fd3bc09-99da-4486-950a-11082f5fd966",
              "externalId": "bb17667e-daac-41f6-ada3-2c22f24caf22"
              },
            "matchType": "ONLY_MATCHING_FIELDS"
        }
      },
      "httpResponse" : {
        "statusCode" : 201,
        "headers" : {
          "content-type" : [ "application/json" ]
        }
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
    When the webhook event is received
    Then the dinopay-gateway creates the corresponding payment on the Payments API
    And the dinopay-gateway produces the following log:
    """
    DinoPay event PaymentCreated processed successfully
    """
