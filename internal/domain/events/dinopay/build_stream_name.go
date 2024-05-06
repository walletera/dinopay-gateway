package dinopay

import "fmt"

const (
    OutboundPaymentStreamNamePrefix = "outboundPayment"
)

func BuildStreamName(id string) string {
    return fmt.Sprintf("%s-%s", OutboundPaymentStreamNamePrefix, id)
}
