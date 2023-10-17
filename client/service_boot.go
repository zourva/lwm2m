package client

type bootstrapReason = int32

const (
	bsrRegRetryFailure bootstrapReason = iota
)

type bootstrapper struct {
}
