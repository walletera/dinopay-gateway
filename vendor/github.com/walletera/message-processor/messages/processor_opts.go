package messages

import (
    "time"

    "github.com/walletera/message-processor/errors"
)

type ErrorCallback func(processingError errors.ProcessingError)

type ProcessorOpts struct {
    errorCallback     ErrorCallback
    processingTimeout time.Duration
}

var defaultProcessorOpts = ProcessorOpts{
    errorCallback:     func(processorError errors.ProcessingError) {},
    processingTimeout: 10 * time.Minute,
}

type ProcessorOpt func(opts *ProcessorOpts)

func WithErrorCallback(errorCallback ErrorCallback) ProcessorOpt {
    return func(opts *ProcessorOpts) {
        opts.errorCallback = errorCallback
    }
}

func WithProcessingTimeout(processingTimeout time.Duration) ProcessorOpt {
    return func(opts *ProcessorOpts) {
        opts.processingTimeout = processingTimeout
    }
}

func applyCustomOpts(opts *ProcessorOpts, customOpts []ProcessorOpt) {
    for _, customOpt := range customOpts {
        customOpt(opts)
    }
}
