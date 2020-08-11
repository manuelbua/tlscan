package progress

type NoOpProgress struct{}

func (p *NoOpProgress) InitProgressbar(totalCount int64) {}
func (p *NoOpProgress) AddToTotal(delta int64)                                                 {}
func (p *NoOpProgress) Update()                                                                {}
func (p *NoOpProgress) Drop(count int64)                                                       {}
func (p *NoOpProgress) Wait()                                                                  {}
