// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// PatchWithdrawal implements patchWithdrawal operation.
//
// Patches a withdrawal.
//
// PATCH /withdrawals/{withdrawalId}
func (UnimplementedHandler) PatchWithdrawal(ctx context.Context, req *WithdrawalPatchBody, params PatchWithdrawalParams) (r PatchWithdrawalRes, _ error) {
	return r, ht.ErrNotImplemented
}