package tests

type publishable struct {
    rawEvent []byte
}

func (p publishable) ID() string {
    //TODO implement me
    panic("implement me")
}

func (p publishable) Type() string {
    //TODO implement me
    panic("implement me")
}

func (p publishable) CorrelationID() string {
    //TODO implement me
    panic("implement me")
}

func (p publishable) DataContentType() string {
    //TODO implement me
    panic("implement me")
}

func (p publishable) Serialize() ([]byte, error) {
    return p.rawEvent, nil
}
