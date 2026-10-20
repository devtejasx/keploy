package proxy

import "go.keploy.io/server/v3/pkg/models"

// lossesOf is shared by the Linux-only pairing tests and by
// tls_rejection_test.go, which builds on every platform.
func lossesOf(p *Proxy, kind models.CaptureLossKind) uint64 {
	for _, l := range p.losses.Snapshot() {
		if l.Kind == kind {
			return l.Count
		}
	}
	return 0
}
