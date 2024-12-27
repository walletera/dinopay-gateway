package outbound

import (
    "context"
    "fmt"

    "github.com/google/uuid"
    dinopayapi "github.com/walletera/dinopay/api"
    paymentsapi "github.com/walletera/payments-types/api"
    "github.com/walletera/werrors"
)

func updatePaymentStatus(ctx context.Context, client *paymentsapi.Client, paymentId uuid.UUID, dinopayPaymentId uuid.UUID, dinopayPaymentStatus string) werrors.WError {
    status, err := dinopayStatus2PaymentsStatus(dinopayPaymentStatus)
    if err != nil {
        return err
    }
    _, patchPaymentErr := client.PatchPayment(
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
    if patchPaymentErr != nil {
        // TODO handle error
        return werrors.NewRetryableInternalError("failed updating withdrawal in payments service: %w", err)
    }
    return nil
}

func dinopayStatus2PaymentsStatus(dinopayStatus string) (paymentsapi.PaymentStatus, werrors.WError) {
    var status paymentsapi.PaymentStatus
    switch dinopayStatus {
    case string(dinopayapi.PaymentStatusPending):
        status = paymentsapi.PaymentStatusPending
    case string(dinopayapi.PaymentStatusConfirmed):
        status = paymentsapi.PaymentStatusConfirmed
    case string(dinopayapi.PaymentStatusRejected):
        status = paymentsapi.PaymentStatusFailed
    default:
        return "", werrors.NewNonRetryableInternalError(fmt.Sprintf("unknown dinopay payment status %s", dinopayStatus))
    }
    return status, nil
}
