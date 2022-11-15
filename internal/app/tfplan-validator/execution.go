package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	tfpv "github.com/fautom/tfplan-validator"
)

// fileMode allows full access for only the current user as the work directory may contain sensitive values
const fileMode = 0700

// planRetries is the number of times we will retry a plan. Helps with connection timeouts
const planRetries = 2

// Workspace describes a single terraform plan operation
type Workspace struct {
	CacheDir     string           `json:"cache_dir"`
	Command      string           `json:"command"`
	InitArgs     string           `json:"init_args"`
	PlanBinPath  string           `json:"plan_bin_path"`
	PlanJsonPath string           `json:"plan_json_path"`
	WorkDir      string           `json:"work_dir"`
	PlanFilter   *tfpv.PlanFilter `json:"plan_filter,omitempty"`
}

// Manifest is used to give the apply operation everything it needs to know about the plans we ran
type Manifest struct {
	Filename     string       `json:"filename"`
	BaseCacheDir string       `json:"base_cache_dir"`
	Workspaces   []*Workspace `json:"workspaces"`
}

func NewManifest(command, initArgs, baseCacheDir string, workspaceDirs []string) (mf *Manifest, err error) {
	mf = &Manifest{}

	if mf.BaseCacheDir, err = filepath.Abs(baseCacheDir); err != nil {
		return nil, err
	} else if mf.Filename, err = filepath.Abs(getManifestPath(baseCacheDir)); err != nil {
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

func getManifestPath(cacheDir string) string {
	return filepath.Join(cacheDir, "manifest.json")
}

func LoadManifest(baseCacheDir string) (*Manifest, error) {
	mf := &Manifest{}

	if filename, err := filepath.Abs(getManifestPath(baseCacheDir)); err != nil {
		return nil, err
	} else if bytes, err := ioutil.ReadFile(filename); err != nil {
		return nil, err
	} else if err := json.Unmarshal(bytes, mf); err != nil {
		return nil, err
	}
	// TODO: validate

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
	ws.Command = command

	return ws, nil
}

func determineCommand(ws *Workspace) {
	if ws.Command != "" {
		return
	} else if _, err := ioutil.ReadFile(filepath.Join(ws.WorkDir, "terragrunt.hcl")); err == io.EOF {
		ws.Command = "terragrunt"
	} else {
		ws.Command = "terraform"
	}
}

func Apply(mf *Manifest) error {
	var err error
	for _, ws := range mf.Workspaces {

		if err = execInit(ws); err != nil {
			return err
		} else if err = execApply(ws); err != nil {
			return err
		}
	}

	return nil
}

func Plan(mf *Manifest) error {
	var err error
	for _, ws := range mf.Workspaces {

		determineCommand(ws)

		if err = execInit(ws); err != nil {
			return err
		} else if err = retryExec(execPlan, ws, planRetries); err != nil {
			return err
		} else if err = execShow(ws); err != nil {
			return err
		} else if ws.PlanFilter, err = tfpv.NewFilterFromPlanPaths([]string{ws.PlanJsonPath}); err != nil {
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

// execInit runs terraform for a single workDir
func execInit(ws *Workspace) error {
	if err := os.MkdirAll(ws.CacheDir, fileMode); err != nil {
		return err
	}
	cmd := exec.Command(ws.Command, "init")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = ws.WorkDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run '%s' from '%s': %w", cmd.String(), cmd.Dir, err)
	}
	return nil
}

// retryExec will retry a function like execPlan
func retryExec(fn func(*Workspace) error, ws *Workspace, retries int) (err error) {
	for retries > 0 {
		err = fn(ws)
		if err == nil {
			return nil
		}
		retries -= 1
	}
	return err
}

// execPlan runs terraform for a single workDir and stores the results in the cache
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

// execShow converts a plan to a json file that can be read by the validator
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

// execApply runs terraform for a single workDir and stores the results in the cache
func execApply(ws *Workspace) error {
	if err := os.MkdirAll(ws.CacheDir, fileMode); err != nil {
		return err
	}
	cmd := exec.Command(ws.Command, "apply", ws.PlanBinPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = ws.WorkDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run '%s' from '%s': %w", cmd.String(), cmd.Dir, err)
	}
	return nil
}
