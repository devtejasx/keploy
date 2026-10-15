package replay

import (
	"context"
	"slices"
	"strings"
	"testing"

	"go.keploy.io/server/v3/config"
	"go.keploy.io/server/v3/pkg/provenswitch"
	"go.keploy.io/server/v3/utils"
	"go.uber.org/zap"
)

func recordedSets(context.Context) ([]string, error) {
	return []string{"test-set-2", "test-set-1", "test-set-10"}, nil
}

// A test set the run is asked for that is not recorded is refused, naming it,
// and the run exits non-zero, with the switch on and off (P189). It was
// dropped silently: -t "test-set-1,test-set-3" ran test-set-1 alone and
// passed, and a selection naming no recorded set at all ran every set.
func TestATestSetThatIsNotRecordedIsRefused(t *testing.T) {
	for _, on := range []bool{false, true} {
		if on {
			t.Setenv(provenswitch.Env, "1")
		} else {
			t.Setenv(provenswitch.Env, "")
		}
		for _, selected := range [][]string{{"test-set-1", "test-set-3"}, {"test-set-3"}} {
			conf := &config.Config{}
			config.SetSelectedTests(conf, selected)
			r := NewReplayer(zap.NewNop(), startTestDB{list: recordedSets}, nil, nil, nil, nil, startTelemetry{}, startInstrumentation{}, nil, conf).(*Replayer)
			sets, err := r.GetSelectedTestSets(context.Background())
			if err == nil || !strings.Contains(err.Error(), "test-set-3") {
				t.Fatalf("switch %v, -t %q: selected %q, err %v; want test-set-3 refused by name", on, selected, sets, err)
			}

			utils.ErrCode = 0
			t.Cleanup(func() { utils.ErrCode = 0 })
			if err := r.Start(context.Background()); err == nil || !strings.Contains(err.Error(), "test-set-3") {
				t.Fatalf("switch %v, -t %q: Start returned %v; want test-set-3 refused by name", on, selected, err)
			}
			if utils.ErrCode != 1 {
				t.Fatalf("switch %v, -t %q: exit code %d, want 1", on, selected, utils.ErrCode)
			}
		}
	}
}

// The test sets named, all of them, in their natural order; none named is
// every recorded set.
func TestTheSelectedTestSetsAreTheOnesNamed(t *testing.T) {
	for _, tc := range []struct {
		selected []string
		want     []string
	}{
		{[]string{"test-set-10", "test-set-1"}, []string{"test-set-1", "test-set-10"}},
		{nil, []string{"test-set-1", "test-set-2", "test-set-10"}},
	} {
		conf := &config.Config{}
		config.SetSelectedTests(conf, tc.selected)
		r := NewReplayer(zap.NewNop(), startTestDB{list: recordedSets}, nil, nil, nil, nil, startTelemetry{}, startInstrumentation{}, nil, conf).(*Replayer)
		got, err := r.GetSelectedTestSets(context.Background())
		if err != nil || !slices.Equal(got, tc.want) {
			t.Fatalf("-t %q: %q, %v; want %q", tc.selected, got, err, tc.want)
		}
	}
}
