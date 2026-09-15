package source

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func preparePNPMBuildPolicy(projectDir string) error {
	cmd := exec.Command("pnpm", "config", "get", "allowBuilds", "--json")
	cmd.Dir = projectDir
	cmd.Stderr = os.Stderr
	data, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("read upstream pnpm build policy: %w", err)
	}

	policy, changed, err := pnpmBuildPolicy(data)
	if err != nil {
		return err
	}
	if !changed {
		return nil
	}

	// Let pnpm update its project configuration so unrelated workspace
	// settings survive, without introducing a second YAML parser here.
	if err := runCommand(projectDir, "pnpm", "config", "set", "--location=project", "--json", "allowBuilds", string(policy)); err != nil {
		return fmt.Errorf("write upstream pnpm build policy: %w", err)
	}
	return nil
}

func pnpmBuildPolicy(data []byte) ([]byte, bool, error) {
	var policy map[string]json.RawMessage
	if err := json.Unmarshal(data, &policy); err != nil {
		return nil, false, fmt.Errorf("parse upstream pnpm allowBuilds: %w", err)
	}
	if policy == nil {
		policy = make(map[string]json.RawMessage)
	}

	changed := false
	// These native dependencies serve extension packaging and publishing,
	// not theme generation. Fill only missing or unresolved decisions.
	for _, name := range []string{"@vscode/vsce-sign", "keytar"} {
		value := strings.TrimSpace(string(policy[name]))
		if value == "" || value == "null" || value == `"set this to true or false"` {
			policy[name] = json.RawMessage("false")
			changed = true
		}
	}
	if !changed {
		return data, false, nil
	}
	updated, err := json.Marshal(policy)
	return updated, true, err
}
