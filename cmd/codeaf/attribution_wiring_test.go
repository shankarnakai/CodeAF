package main

import (
	"os"
	"strings"
	"testing"

	"github.com/Agent-Field/codeaf/internal/manual"
)

// Every entry point that builds its own executor has to pass both attribution
// rows into its signature. The construction sites sit deep inside the resident
// dispatch and the headless run, past a live provider and a real store, so this
// reads the wiring instead.
func TestEveryExecutorConstructionSiteCarriesAttribution(t *testing.T) {
	for _, name := range []string{"subharness.go", "run.go", "exec.go", "chatv3_subharness.go", "subharness_run.go"} {
		raw, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		source := string(raw)
		if strings.Count(source, "exec.NewLinear(") != 1 {
			t.Fatalf("%s no longer builds exactly one leaf loop", name)
		}
		if !strings.Contains(source, "WithAssistedBy(config.AssistedByModelAt(") {
			t.Fatalf("%s builds an executor without wiring the model-name row into its signature", name)
		}
		if !strings.Contains(source, "WithAttributionOff(!config.AttributionAt(") {
			t.Fatalf("%s builds an executor without wiring the attribution on/off row", name)
		}
	}
	raw, err := os.ReadFile("chat.go")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), "exec.NewLinear(") {
		t.Fatal("chat.go builds a leaf loop of its own again instead of asking the table for one")
	}
}

// A capability codeaf has and cannot explain is one the user meets first as a
// surprise in their own git history.
func TestManualExplainsAttribution(t *testing.T) {
	for _, term := range []string{"attribution", "attribution.model", "CODEAF_ATTRIBUTION", "CODEAF_ATTRIBUTION_MODEL", "Assisted-by", "Interface", "CONTRIBUTING", "AI attribution on commits"} {
		if !manual.Chat().Mentions(term) {
			t.Fatalf("no chat manual page mentions %q", term)
		}
	}
}
