package messages

import (
    "time"

    "github.com/walletera/werrors"
)

type ErrorCallback func(processingError werrors.WError)

type ProcessorOpts struct {
    errorCallback     ErrorCallback
    processingTimeout time.Duration
}

var defaultProcessorOpts = ProcessorOpts{
    errorCallback:     func(err werrors.WError) {},
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
