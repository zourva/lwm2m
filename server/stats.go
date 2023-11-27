package server

type Statistics struct {
	requestCount uint64
}

func NewStatistics() *Statistics {
	return &Statistics{
		requestCount: 0,
	}
}

func (s *Statistics) IncrementRequestCount() {
	s.requestCount++
}

func (s *Statistics) GetRequestCount() uint64 {
	return s.requestCount
}
