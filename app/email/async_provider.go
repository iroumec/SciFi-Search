package email

import "scifi-search/app/workers"

type AsyncProvider struct {
	queue        chan workers.Job
	syncProvider Provider
}

func NewAsyncProvider(queue chan workers.Job, sync Provider) Provider {
	return &AsyncProvider{
		queue:        queue,
		syncProvider: sync,
	}
}

func (p *AsyncProvider) Send(to, subject, body string) error {
	p.queue <- func() {
		_ = p.syncProvider.Send(to, subject, body)
	}
	return nil
}
