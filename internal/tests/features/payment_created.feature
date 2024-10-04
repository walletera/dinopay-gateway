Feature: process PaymentCreated event
  Walletera users with funds on their accounts want to be able to Withdraw those funds (totally or partially)
  to DinoPay accounts.
  - When a Withdrawal is created, A PaymentCreated event is published by the Payments Service.
  - The DinoPay Gateway must listen for PaymentCreated events.
  - Whenever a PaymentCreated event arrives the dinopay-gateway must create a corresponding Payment on DinoPay API.

  Background: the dinopay-gateway is up and running
    Given a running dinopay-gateway

  Scenario: payment created event is processed successfully
    Given a PaymentCreated event:
    """json
    {
      "id": "0fg1833e-3438-4908-b90a-5721670cb067",
      "type": "PaymentCreated",
      "data": {
        "id": "0ae1733e-7538-4908-b90a-5721670cb093",
        "amount": 100,
        "currency": "USD",
        "direction": "outbound",
        "customerId": "2432318c-4ff3-4ac0-b734-9b61779e2e46",
        "status": "pending",
        "beneficiary": {
          "bankName": "dinopay",
          "bankId": "dinopay",
          "accountHolder": "Richard Roe",
          "routingKey": "123456789123456",
          "accountNumber": "1200079635"
        },
        "createdAt": "2024-10-04T00:00:00Z"
      }
    }
    """
    And  a dinopay endpoint to create payments:
    # the json below is a mockserver expectation
    """json
    {
      "id": "createPaymentSucceed",
      "httpRequest" : {
        "method": "POST",
        "path" : "/payments",
        "body": {
            "type": "JSON",
            "json": {
              "customerTransactionId": "0ae1733e-7538-4908-b90a-5721670cb093",
              "amount": 100,
              "currency": "USD",
              "destinationAccount": {
                "accountHolder": "Richard Roe",
                "accountNumber": "1200079635"
              }
            },
            "matchType": "ONLY_MATCHING_FIELDS"
        }
      },
      "httpResponse" : {
        "statusCode" : 201,
        "headers" : {
          "content-type" : [ "application/json" ]
        },
        "body" : {
          "id" : "bb17667e-daac-41f6-ada3-2c22f24caf22",
          "amount" : 100,
          "currency" : "USD",
          "sourceAccount" : {
            "accountHolder" : "john doe",
            "accountNumber" : "IE12BOFI90000112345678"
          },
          "destinationAccount" : {
            "accountHolder" : "jane doe",
            "accountNumber" : "IE12BOFI90000112349876"
          },
          "status" : "pending",
          "customerTransactionId" : "9713ec22-cf8d-4a21-affb-719db00d7388",
          "createdAt" : "2023-07-07",
          "updatedAt" : "2023-07-07"
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
    And  a payments endpoint to update payments:
    # the json below is a mockserver expectation
    """json
    {
      "id": "updatePaymentSucceed",
      "httpRequest" : {
        "method": "PATCH",
        "path": "/payments/0ae1733e-7538-4908-b90a-5721670cb093",
        "body": {
            "type": "JSON",
            "json": {
              "externalId": "bb17667e-daac-41f6-ada3-2c22f24caf22",
              "status": "pending"
              },
            "matchType": "ONLY_MATCHING_FIELDS"
        }
      },
      "httpResponse" : {
        "statusCode" : 200,
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
    When the event is published
    Then the dinopay-gateway creates the corresponding payment on the DinoPay API
    And the dinopay-gateway updates the payment on payments service
    And the dinopay-gateway produces the following log:
    """
    PaymentCreated event processed successfully
    """
    And the dinopay-gateway produces the following log:
    """
    OutboundPaymentCreated event processed successfully
    """

  Scenario: payment created event processing failed when trying to create payment on Dinopay
    Given a PaymentCreated event:
    """json
    {
      "id": "0fg1833e-3438-4908-b90a-5721670cb067",
      "type": "PaymentCreated",
      "data": {
        "id": "0ae1733e-7538-4908-b90a-5721670cb093",
        "amount": 100,
        "currency": "USD",
        "direction": "outbound",
        "customerId": "2432318c-4ff3-4ac0-b734-9b61779e2e46",
        "status": "pending",
        "beneficiary": {
          "bankName": "dinopay",
          "bankId": "dinopay",
          "accountHolder": "Richard Roe",
          "routingKey": "123456789123456",
          "accountNumber": "1200079635"
        },
        "createdAt": "2024-10-04T00:00:00Z"
      }
    }
    """
    And  a dinopay endpoint to create payments:
    # the json below is a mockserver expectation
    """json
    {
      "id": "createPaymentFail",
      "httpRequest" : {
        "method": "POST",
        "path" : "/payments",
        "body": {
            "type": "JSON",
            "json": {
              "customerTransactionId": "0ae1733e-7538-4908-b90a-5721670cb093",
              "amount": 100,
              "currency": "USD",
              "destinationAccount": {
                "accountHolder": "Richard Roe",
                "accountNumber": "1200079635"
              }
            },
            "matchType": "ONLY_MATCHING_FIELDS"
        }
      },
      "httpResponse" : {
        "statusCode" : 500,
        "headers" : {
          "content-type" : [ "text/html" ]
        },
        "body" : "something bad happened"
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
    When the event is published
    Then the dinopay-gateway fails creating the corresponding payment on the DinoPay API
    And  the dinopay-gateway produces the following log:
    """
    failed creating payment on dinopay
    """
