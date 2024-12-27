package messages

import (
    "context"
    "errors"
    "fmt"

    "github.com/walletera/eventskit/events"
    "github.com/walletera/werrors"
)

type Processor[Handler any] struct {
    messageConsumer    Consumer
    eventsDeserializer events.Deserializer[Handler]
    eventsHandler      Handler
    opts               ProcessorOpts
}

func NewProcessor[Handler any](
    messageConsumer Consumer,
    eventsDeserializer events.Deserializer[Handler],
    eventsHandler Handler,
    customOpts ...ProcessorOpt,
) *Processor[Handler] {

    opts := defaultProcessorOpts
    applyCustomOpts(&opts, customOpts)

    return &Processor[Handler]{
        messageConsumer:    messageConsumer,
        eventsDeserializer: eventsDeserializer,
        eventsHandler:      eventsHandler,
        opts:               opts,
    }
}

func (p *Processor[Handler]) Start(ctx context.Context) error {
    msgCh, err := p.startMessageConsumer(ctx)
    if err != nil {
        return err
    }
    go p.processMsgs(ctx, msgCh)
    return nil
}

func (p *Processor[Handler]) startMessageConsumer(ctx context.Context) (<-chan Message, error) {
    msgCh, err := p.messageConsumer.Consume()
    if err != nil {
        return nil, fmt.Errorf("failed consuming from message consumer: %w", err)
    }
    go func() {
        <-ctx.Done()
        err := p.messageConsumer.Close()
        if err != nil {
            p.opts.errorCallback(werrors.NewRetryableInternalError("failed closing message consumer: " + err.Error()))
        }
    }()
    return msgCh, nil
}

func (p *Processor[Handler]) processMsgs(ctx context.Context, ch <-chan Message) {
    for msg := range ch {
        go p.processMsgWithTimeout(ctx, msg)
    }
}

func (p *Processor[Handler]) processMsgWithTimeout(ctx context.Context, msg Message) {
    ctxWithTimeout, cancelCtx := context.WithTimeout(ctx, p.opts.processingTimeout)
    defer cancelCtx()
    processMsgDone := make(chan any)
    go func() {
        p.processMsg(ctxWithTimeout, msg)
        close(processMsgDone)
    }()
    select {
    case <-ctxWithTimeout.Done():
    case <-processMsgDone:
    }
    err := ctxWithTimeout.Err()
    if err != nil {
        if errors.Is(err, context.DeadlineExceeded) {
            p.handleError(msg, werrors.NewTimeoutError(err.Error()))
        }
    }
}

func (p *Processor[Handler]) processMsg(ctx context.Context, message Message) {
    event, err := p.eventsDeserializer.Deserialize(message.Payload())
    if err != nil {
        p.handleError(message, werrors.NewUnprocessableMessageError(err.Error()))
        return
    }
    if event == nil {
        return
    }
    processingErr := event.Accept(ctx, p.eventsHandler)
    if processingErr != nil {
        p.handleError(message, processingErr)
    } else {
        message.Acknowledger().Ack()
    }
}

func (p *Processor[Handler]) handleError(message Message, err werrors.WError) {
    if p.opts.errorCallback != nil {
        p.opts.errorCallback(err)
    }
    nackOpts := NackOpts{
        Requeue:      err.IsRetryable(),
        ErrorCode:    err.Code(),
        ErrorMessage: err.Message(),
    }
    message.Acknowledger().Nack(nackOpts)
}
