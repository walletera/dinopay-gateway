Feature: process WithdrawalCreated event
  Walletera users with funds on their accounts want to be able to Withdraw those funds (totally or partially)
  to DinoPay accounts.
  - When a Withdrawal is created, A WithdrawalCreated event is published by the Payments Service.
  - The DinoPay Gateway must listen for WithdrawalCreated events.
  - Whenever a WithdrawalCreated event arrives the dinopay-gateway must create a corresponding Payment on DinoPay API.

  Scenario: withdrawal created event is processed successfully
    Given a withdrawal created event:
    """json
    {
      "type": "WithdrawalCreated",
      "data": {
         "id": "0ae1733e-7538-4908-b90a-5721670cb093",
         "user_id": "2432318c-4ff3-4ac0-b734-9b61779e2e46",
         "psp_id": "dinopay",
         "external_id": null,
         "amount": 100,
         "currency": "USD",
         "status": "pending",
         "beneficiary": {
           "id": "2f98dbe7-72ab-4167-9be5-ecd3608b55e4",
           "description": "Richard Roe DinoPay account",
           "account": {
            "holder": "Richard Roe",
            "number": 1200079635
           }
         }
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
    When the event is published
    Then the dinopay-gateway creates the corresponding payment on the DinoPay API
    And  the dinopay-gateway produces the following log:
    """
    WithdrawalCreated event processed successfully
    """
