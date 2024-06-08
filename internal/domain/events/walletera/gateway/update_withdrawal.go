package gateway

import (
    "context"
    "fmt"

    "github.com/google/uuid"
    dinopayApi "github.com/walletera/dinopay/api"
    paymentsApi "github.com/walletera/payments/api"
)

func updateWithdrawalStatus(client *paymentsApi.Client, withdrawalId uuid.UUID, dinopayPaymentId uuid.UUID, dinopayPaymentStatus string) error {
    status, err := dinopayStatus2PaymentsStatus(dinopayPaymentStatus)
    if err != nil {
        return err
    }
    _, err = client.PatchWithdrawal(context.Background(),
        &paymentsApi.WithdrawalPatchBody{
            ExternalId: paymentsApi.OptUUID{
                Value: dinopayPaymentId,
                Set:   true,
            },
            Status: status,
        },
        paymentsApi.PatchWithdrawalParams{
            WithdrawalId: withdrawalId,
        },
    )
    if err != nil {
        return fmt.Errorf("failed updating withdrawal in payments service: %w", err)
    }
    return nil
}

func dinopayStatus2PaymentsStatus(dinopayStatus string) (paymentsApi.WithdrawalPatchBodyStatus, error) {
    var status paymentsApi.WithdrawalPatchBodyStatus
    switch dinopayStatus {
    case string(dinopayApi.PaymentStatusPending):
        status = paymentsApi.WithdrawalPatchBodyStatusPending
    case string(dinopayApi.PaymentStatusConfirmed):
        status = paymentsApi.WithdrawalPatchBodyStatusPending
    case string(paymentsApi.WithdrawalPatchBodyStatusRejected):
        status = paymentsApi.WithdrawalPatchBodyStatusRejected
    default:
        return "", fmt.Errorf("unknown dinopay payment status %s", dinopayStatus)
    }
    return status, nil
}
