package server

// Statistics defined for the server.
type Statistics interface {
	IncrementRequestCount()
	GetRequestCount() uint64
}

type DefaultStatistics struct {
	requestCount uint64
}

func (s *DefaultStatistics) IncrementRequestCount() {
	s.requestCount++
}

func (s *DefaultStatistics) GetRequestCount() uint64 {
	return s.requestCount
}
