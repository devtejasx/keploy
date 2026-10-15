package http

import "testing"

// keploy.yml's request-body noise names whole fields (P093): noise "id"
// excuses a body's id, never its paid amount. As a substring of the field's
// path it did, so a call whose amount changed was served the recording of
// another amount.
func TestStrictDeclaredBodyNoiseNamesWholeFields(t *testing.T) {
	ct := map[string]string{"Content-Type": "application/json"}
	rec := recorded("m", "/v1/charge", ct, `{"id":"a-1","paid":100,"user":{"id":7}}`)
	noise := map[string][]string{"id": {}}
	if got, diag := strictMatch(t, sent("/v1/charge", ct, `{"id":"a-2","paid":100,"user":{"id":8}}`), nil, noise, nil, rec); got == nil {
		t.Fatalf("the ids the noise names were compared: %v", differs(diag))
	}
	if got, diag := strictMatch(t, sent("/v1/charge", ct, `{"id":"a-1","paid":0,"user":{"id":7}}`), nil, noise, nil, rec); got != nil || !hasField(differs(diag), "body.paid") {
		t.Fatalf("noise \"id\" excused paid: served %v, differs %v", got != nil, differs(diag))
	}
}
