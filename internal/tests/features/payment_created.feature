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
    """
    data/payment_created_event.json
    """
    And  a dinopay endpoint to create payments:
    """
    data/dinopay_create_payment_endpoint_expectation.json
    """
    And  a payments endpoint to update payments:
    """
    data/payments_update_payments_endpoint_expectation.json
    """
    When the event is published
    Then the dinopay-gateway creates the corresponding payment on the DinoPay API
    And the dinopay-gateway produces the following log:
    """
    payments PaymentCreated event processed successfully
    """
#    And the dinopay-gateway produces the following log:
#    """
#    OutboundPaymentCreated event processed successfully
#    """
    And the dinopay-gateway updates the payment on payments service

<<<<<<< Updated upstream
  Scenario: payment created event processing failed when trying to create payment on Dinopay
    Given a PaymentCreated event:
    """
    data/payment_created_event.json
    """
    And  a dinopay endpoint to create payments:
    """
    data/dinopay_create_payments_endpoint_500_response_expectation.json
    """
    When the event is published
    Then the dinopay-gateway fails creating the corresponding payment on the DinoPay API
    And  the dinopay-gateway produces the following log:
    """
    failed creating payment on dinopay
    """
=======

#  Scenario: payment created event processing failed when trying to create payment on Dinopay
#    Given a PaymentCreated event:
#    """json
#{
#      "id": "eefe8e76-58eb-4a3e-bf05-c7f703ddc220",
#      "type": "PaymentCreated",
#      "data": {
#        "id": "0ae1733e-7538-4908-b90a-5721670cb093",
#        "customerId": "abbb8aa3-87f9-4b2b-889f-8962cf708cfc",
#        "amount": 100,
#        "currency": "USD",
#        "gateway": "dinopay",
#        "direction": "outbound",
#        "status": "pending",
#        "debtor": {
#          "institutionName": "dinopay",
#          "institutionId": "dinopay",
#          "currency": "ARS",
#          "accountDetails": {
#            "accountType": "dinopay",
#            "accountHolder": "Richard Roe",
#            "accountNumber": "1200079635"
#          }
#        },
#        "beneficiary": {
#          "institutionName": "dinopay",
#          "institutionId": "dinopay",
#          "currency": "ARS",
#          "accountDetails": {
#             "accountType": "dinopay",
#             "accountHolder": "Richard Roe",
#            "accountNumber": "1200079635"
#            }
#        },
#      "updatedAt": "2024-06-27T15:45:00Z",
#      "createdAt": "2024-06-27T15:45:00Z"
#      },
#      "createdAt": "2024-06-27T15:45:00Z"
#    }
#    """
#    And  a dinopay endpoint to create payments:
#    # the json below is a mockserver expectation
#    """json
#    {
#      "id": "createPaymentFail",
#      "httpRequest" : {
#        "method": "POST",
#        "path" : "/payments",
#        "body": {
#            "type": "JSON",
#            "json": {
#              "customerTransactionId": "0ae1733e-7538-4908-b90a-5721670cb093",
#              "amount": 100,
#              "currency": "USD",
#              "sourceAccount" : {
#                "accountHolder" : "Richard Roe",
#                "accountNumber" : "Richard Roe"
#              },
#              "destinationAccount": {
#                "accountHolder": "Richard Roe",
#                "accountNumber": "1200079635"
#              }
#            },
#            "matchType": "ONLY_MATCHING_FIELDS"
#        }
#      },
#      "httpResponse" : {
#        "statusCode" : 500,
#        "headers" : {
#          "content-type" : [ "text/html" ]
#        },
#        "body" : "something bad happened"
#      },
#      "priority" : 0,
#      "timeToLive" : {
#        "unlimited" : true
#      },
#      "times" : {
#        "unlimited" : true
#      }
#    }
#    """
#    When the event is published
#    Then the dinopay-gateway fails creating the corresponding payment on the DinoPay API
#    And  the dinopay-gateway produces the following log:
#    """
#    failed creating payment on dinopay
#    """
>>>>>>> Stashed changes
