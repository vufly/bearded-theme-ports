package app

import (
	"slices"
	"strings"
	"testing"
)

func TestParseTargetsPiPigmentUsesTMThemeBuild(t *testing.T) {
	targets, installTargets, err := parseTargets([]string{"pi-pigment"}, true)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(targets, []string{"tmtheme"}) || !slices.Equal(installTargets, []string{"pi-pigment"}) {
		t.Fatalf("parseTargets() = %v, %v; want tmtheme build and pi-pigment install only", targets, installTargets)
	}
}

func TestParseTargetsRejectsPiPigmentWithoutInstall(t *testing.T) {
	_, _, err := parseTargets([]string{"pi-pigment"}, false)
	if err == nil || !strings.Contains(err.Error(), "--install pi-pigment") {
		t.Fatalf("parseTargets() error = %v, want --install guidance", err)
	}
}
