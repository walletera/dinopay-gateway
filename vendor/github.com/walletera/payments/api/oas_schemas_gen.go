// Code generated by ogen, DO NOT EDIT.

package api

import (
	"io"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
)

// NewOptUUID returns new OptUUID with value set to v.
func NewOptUUID(v uuid.UUID) OptUUID {
	return OptUUID{
		Value: v,
		Set:   true,
	}
}

// OptUUID is optional uuid.UUID.
type OptUUID struct {
	Value uuid.UUID
	Set   bool
}

// IsSet returns true if OptUUID was set.
func (o OptUUID) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptUUID) Reset() {
	var v uuid.UUID
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptUUID) SetTo(v uuid.UUID) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptUUID) Get() (v uuid.UUID, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptUUID) Or(d uuid.UUID) uuid.UUID {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// A text message describing an error.
// Ref: #/components/schemas/errorMessage
type PatchWithdrawalBadRequest struct {
	Data io.Reader
}

// Read reads data from the Data reader.
//
// Kept to satisfy the io.Reader interface.
func (s PatchWithdrawalBadRequest) Read(p []byte) (n int, err error) {
	if s.Data == nil {
		return 0, io.EOF
	}
	return s.Data.Read(p)
}

func (*PatchWithdrawalBadRequest) patchWithdrawalRes() {}

// PatchWithdrawalOK is response for PatchWithdrawal operation.
type PatchWithdrawalOK struct{}

func (*PatchWithdrawalOK) patchWithdrawalRes() {}

// Body of the PATH /withdrawal request.
// Ref: #/components/schemas/withdrawalPatchBody
type WithdrawalPatchBody struct {
	// Id assigned to the operation by the external payment provider.
	ExternalId OptUUID `json:"externalId"`
	// Withdrawal status.
	Status WithdrawalPatchBodyStatus `json:"status"`
}

// GetExternalId returns the value of ExternalId.
func (s *WithdrawalPatchBody) GetExternalId() OptUUID {
	return s.ExternalId
}

// GetStatus returns the value of Status.
func (s *WithdrawalPatchBody) GetStatus() WithdrawalPatchBodyStatus {
	return s.Status
}

// SetExternalId sets the value of ExternalId.
func (s *WithdrawalPatchBody) SetExternalId(val OptUUID) {
	s.ExternalId = val
}

// SetStatus sets the value of Status.
func (s *WithdrawalPatchBody) SetStatus(val WithdrawalPatchBodyStatus) {
	s.Status = val
}

// Withdrawal status.
type WithdrawalPatchBodyStatus string

const (
	WithdrawalPatchBodyStatusPending   WithdrawalPatchBodyStatus = "pending"
	WithdrawalPatchBodyStatusConfirmed WithdrawalPatchBodyStatus = "confirmed"
	WithdrawalPatchBodyStatusRejected  WithdrawalPatchBodyStatus = "rejected"
)

// AllValues returns all WithdrawalPatchBodyStatus values.
func (WithdrawalPatchBodyStatus) AllValues() []WithdrawalPatchBodyStatus {
	return []WithdrawalPatchBodyStatus{
		WithdrawalPatchBodyStatusPending,
		WithdrawalPatchBodyStatusConfirmed,
		WithdrawalPatchBodyStatusRejected,
	}
}

// MarshalText implements encoding.TextMarshaler.
func (s WithdrawalPatchBodyStatus) MarshalText() ([]byte, error) {
	switch s {
	case WithdrawalPatchBodyStatusPending:
		return []byte(s), nil
	case WithdrawalPatchBodyStatusConfirmed:
		return []byte(s), nil
	case WithdrawalPatchBodyStatusRejected:
		return []byte(s), nil
	default:
		return nil, errors.Errorf("invalid value: %q", s)
	}
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *WithdrawalPatchBodyStatus) UnmarshalText(data []byte) error {
	switch WithdrawalPatchBodyStatus(data) {
	case WithdrawalPatchBodyStatusPending:
		*s = WithdrawalPatchBodyStatusPending
		return nil
	case WithdrawalPatchBodyStatusConfirmed:
		*s = WithdrawalPatchBodyStatusConfirmed
		return nil
	case WithdrawalPatchBodyStatusRejected:
		*s = WithdrawalPatchBodyStatusRejected
		return nil
	default:
		return errors.Errorf("invalid value: %q", data)
	}
}