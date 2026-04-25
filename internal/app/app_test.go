package app

import "testing"

func TestParseTargetsInstallBatUsesTMThemeBuild(t *testing.T) {
	targets, installTargets, err := parseTargets([]string{"bat"}, true)
	if err != nil {
		t.Fatalf("parseTargets() error = %v", err)
	}

	if len(targets) != 1 || targets[0] != "tmtheme" {
		t.Fatalf("parseTargets() targets = %#v, want []string{\"tmtheme\"}", targets)
	}
	if len(installTargets) != 1 || installTargets[0] != "bat" {
		t.Fatalf("parseTargets() installTargets = %#v, want []string{\"bat\"}", installTargets)
	}
}

func TestParseTargetsRejectsBatWithoutInstall(t *testing.T) {
	_, _, err := parseTargets([]string{"bat"}, false)
	if err == nil {
		t.Fatal("parseTargets() error = nil, want error")
	}
}
