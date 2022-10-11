package tfplan_validator

import (
	"bytes"
	"encoding/json"
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

// Manifest
type Manifest struct {
	Filename     string
	BaseCacheDir string
	Workspaces   []*Workspace
}

func NewManifest(command, initArgs, baseCacheDir string, workspaceDirs []string) (mf *Manifest, err error) {
	mf = &Manifest{}

	if mf.BaseCacheDir, err = filepath.Abs(baseCacheDir); err != nil {
		return nil, err
	} else if mf.Filename, err = filepath.Abs(filepath.Join(baseCacheDir, "manifest.json")); err != nil {
		return nil, err
	}

	for _, workspaceDir := range workspaceDirs {
		if ws, err := NewWorkspace(command, initArgs, baseCacheDir, workspaceDir); err != nil {
			return nil, err
		} else {
			mf.Workspaces = append(mf.Workspaces, ws)
		}
	}
	return mf, nil
}

func NewWorkspace(command, initArgs, baseCacheDir, workDir string) (ws *Workspace, err error) {
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

func Plan(mf *Manifest) error {
	for _, ws := range mf.Workspaces {
		if err := execPlan(ws); err != nil {
			return err
		} else if err := execShow(ws); err != nil {
			return err
		}
	}
	if err := saveManifest(mf); err != nil {
		return err
	}
	return nil
}

func saveManifest(mf *Manifest) error {
	if err := os.MkdirAll(mf.BaseCacheDir, fileMode); err != nil {
		return err
	} else if bytes, err := json.MarshalIndent(mf, "", "  "); err != nil {
		return err
	} else if err := ioutil.WriteFile(mf.Filename, bytes, fileMode); err != nil {
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
