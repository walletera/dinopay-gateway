package outbound

import (
    "context"
    "fmt"

    "github.com/google/uuid"
    dinopayapi "github.com/walletera/dinopay/api"
    paymentsapi "github.com/walletera/payments-types/api"
)

func updatePaymentStatus(ctx context.Context, client *paymentsapi.Client, paymentId uuid.UUID, dinopayPaymentId uuid.UUID, dinopayPaymentStatus string) error {
    status, err := dinopayStatus2PaymentsStatus(dinopayPaymentStatus)
    if err != nil {
        return err
    }
    _, err = client.PatchPayment(
        ctx,
        &paymentsapi.PaymentUpdate{
            PaymentId: paymentId,
            ExternalId: paymentsapi.OptUUID{
                Value: dinopayPaymentId,
                Set:   true,
            },
            Status: status,
        },
        paymentsapi.PatchPaymentParams{
            PaymentId: paymentId,
        })
    if err != nil {
        return fmt.Errorf("failed updating withdrawal in payments service: %w", err)
    }
    return nil
}

func dinopayStatus2PaymentsStatus(dinopayStatus string) (paymentsapi.PaymentStatus, error) {
    var status paymentsapi.PaymentStatus
    switch dinopayStatus {
    case string(dinopayapi.PaymentStatusPending):
        status = paymentsapi.PaymentStatusPending
    case string(dinopayapi.PaymentStatusConfirmed):
        status = paymentsapi.PaymentStatusConfirmed
    case string(dinopayapi.PaymentStatusRejected):
        status = paymentsapi.PaymentStatusFailed
    default:
        return "", fmt.Errorf("unknown dinopay payment status %s", dinopayStatus)
    }
    return status, nil
}
