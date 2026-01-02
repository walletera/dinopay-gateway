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
