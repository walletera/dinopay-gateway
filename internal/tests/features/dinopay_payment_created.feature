Feature: process DinoPay webhook event PaymentCreated
  DinoPay sends a webhook event of type PaymentCreated

  Background: the dinopay-gateway is up and running
    Given a running dinopay-gateway

  Scenario: the payment is valid and matches the accountNumber of a walletera user
    Given a DinoPay PaymentCreated event:
    """
    data/dinopay_payment_created_event.json
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
