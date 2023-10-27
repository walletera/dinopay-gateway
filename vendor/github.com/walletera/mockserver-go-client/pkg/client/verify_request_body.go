package client

type ExpectationId struct {
    Id string `json:"id"`
}

type VerifyRequestBody struct {
    ExpectationId ExpectationId `json:"expectationId"`
}
