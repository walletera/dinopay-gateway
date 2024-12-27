package logattr

import "log/slog"

func ServiceName(serviceName string) slog.Attr {
    return slog.String("service_name", serviceName)
}

func Component(component string) slog.Attr {
    return slog.String("component", component)
}

func PaymentId(paymentId string) slog.Attr {
    return slog.String("payment_id", paymentId)
}

func DinopayPaymentId(dinopayPaymentId string) slog.Attr {
    return slog.String("dinopay_payment_id", dinopayPaymentId)
}

func EventType(eventType string) slog.Attr {
    return slog.String("event_type", eventType)
}

func Error(err string) slog.Attr {
    return slog.String("error", err)
}
