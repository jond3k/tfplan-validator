package tfplan_validator

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// fileMode allows full access for only the current user as the work directory may contain sensitive values
const fileMode = 0700

// Workspace encapsulates the run of a terraform command
type Workspace struct {
	CacheDir     string
	Command      string
	InitArgs     string
	PlanBinPath  string
	PlanJsonPath string
	WorkDir      string
}

func NewWorkspace(command, initArgs, baseCacheDir string, workDir string) (ws *Workspace, err error) {
	ws = &Workspace{}
	ws.InitArgs = initArgs
	if ws.CacheDir, err = filepath.Abs(filepath.Join(baseCacheDir, workDir)); err != nil {
		return nil, err
	} else if ws.WorkDir, err = filepath.Abs(workDir); err != nil {
		return nil, err
	} else if ws.PlanBinPath, err = filepath.Abs(filepath.Join(ws.CacheDir, "plan.bin")); err != nil {
		return nil, err
	} else if ws.PlanJsonPath, err = filepath.Abs(filepath.Join(ws.CacheDir, "plan.json")); err != nil {
		return nil, err
	}

	if command != "" {
		ws.Command = command
	} else if _, err := ioutil.ReadFile(filepath.Join(workDir, "terragrunt.hcl")); err == io.EOF {
		ws.Command = "terragrunt"
	} else {
		ws.Command = "terraform"
	}

	return ws, nil
}

func Plan(ws *Workspace) error {
	if err := execPlan(ws); err != nil {
		return err
	} else if err := execShow(ws); err != nil {
		return err
	}
	return nil
}

// Plan runs terraform for a single workDir and stores the results in the cache
func execPlan(ws *Workspace) error {
	if err := os.MkdirAll(ws.CacheDir, fileMode); err != nil {
		return err
	}
	cmd := exec.Command(ws.Command, "plan", "-out", ws.PlanBinPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = ws.WorkDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run '%s' from '%s': %w", cmd.String(), cmd.Dir, err)
	}
	return nil
}

// Show converts a plan to a json file that can be read by the validator
func execShow(ws *Workspace) error {
	if err := os.MkdirAll(ws.CacheDir, fileMode); err != nil {
		return err
	}
	cmd := exec.Command(ws.Command, "show", "-json", ws.PlanBinPath)
	outbuf := &bytes.Buffer{}
	cmd.Dir = ws.WorkDir
	cmd.Stdout = outbuf
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run '%s' from '%s': %w", cmd.String(), cmd.Dir, err)
	} else if err := ioutil.WriteFile(ws.PlanJsonPath, outbuf.Bytes(), fileMode); err != nil {
		return fmt.Errorf("failed to write %s: %w", ws.PlanJsonPath, err)
	}
	return nil
}
