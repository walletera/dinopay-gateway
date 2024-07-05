package messages

import (
    "context"
    "errors"
    "fmt"

    procerrors "github.com/walletera/message-processor/errors"
    "github.com/walletera/message-processor/events"
)

type Processor[Visitor any] struct {
    messageConsumer    Consumer
    eventsDeserializer events.Deserializer[Visitor]
    eventsVisitor      Visitor
    opts               ProcessorOpts
}

func NewProcessor[Visitor any](
    messageConsumer Consumer,
    eventsDeserializer events.Deserializer[Visitor],
    eventsVisitor Visitor,
    customOpts ...ProcessorOpt,
) *Processor[Visitor] {

    opts := defaultProcessorOpts
    applyCustomOpts(&opts, customOpts)

    return &Processor[Visitor]{
        messageConsumer:    messageConsumer,
        eventsDeserializer: eventsDeserializer,
        eventsVisitor:      eventsVisitor,
        opts:               opts,
    }
}

func (p *Processor[Visitor]) Start(ctx context.Context) error {
    msgCh, err := p.startMessageConsumer(ctx)
    if err != nil {
        return err
    }
    go p.processMsgs(ctx, msgCh)
    return nil
}

func (p *Processor[Visitor]) startMessageConsumer(ctx context.Context) (<-chan Message, error) {
    msgCh, err := p.messageConsumer.Consume()
    if err != nil {
        return nil, fmt.Errorf("failed consuming from message consumer: %w", err)
    }
    go func() {
        <-ctx.Done()
        err := p.messageConsumer.Close()
        if err != nil {
            p.opts.errorCallback(procerrors.NewInternalError("failed closing message consumer: " + err.Error()))
        }
    }()
    return msgCh, nil
}

func (p *Processor[Visitor]) processMsgs(ctx context.Context, ch <-chan Message) {
    for msg := range ch {
        go p.processMsgWithTimeout(ctx, msg)
    }
}

func (p *Processor[Visitor]) processMsgWithTimeout(ctx context.Context, msg Message) {
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
            p.handleError(msg, procerrors.NewTimeoutError(err.Error()))
        }
    }
}

func (p *Processor[Visitor]) processMsg(ctx context.Context, message Message) {
    event, err := p.eventsDeserializer.Deserialize(message.Payload())
    if err != nil {
        p.handleError(message, procerrors.NewUnprocessableMessageError(err.Error()))
        return
    }
    if event == nil {
        return
    }
    processingErr := event.Accept(ctx, p.eventsVisitor)
    if processingErr != nil {
        p.handleError(message, processingErr)
    } else {
        message.Acknowledger().Ack()
    }
}

func (p *Processor[Visitor]) handleError(message Message, err procerrors.ProcessingError) {
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
