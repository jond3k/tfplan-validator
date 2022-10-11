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

// TerraformExec encapsulates the run of a terraform command
type TerraformExec struct {
	cacheDir     string
	command      string
	initArgs     string
	planBinPath  string
	planJsonPath string
	workDir      string
}

func NewTerraformExec(command, initArgs, baseCacheDir string, workDir string) (pe *TerraformExec, err error) {
	pe = &TerraformExec{}
	pe.initArgs = initArgs
	if pe.cacheDir, err = filepath.Abs(filepath.Join(baseCacheDir, workDir)); err != nil {
		return nil, err
	} else if pe.workDir, err = filepath.Abs(workDir); err != nil {
		return nil, err
	} else if pe.planBinPath, err = filepath.Abs(filepath.Join(pe.cacheDir, "plan.bin")); err != nil {
		return nil, err
	} else if pe.planJsonPath, err = filepath.Abs(filepath.Join(pe.cacheDir, "plan.json")); err != nil {
		return nil, err
	}

	if command != "" {
		pe.command = command
	} else if _, err := ioutil.ReadFile(filepath.Join(workDir, "terragrunt.hcl")); err == io.EOF {
		pe.command = "terragrunt"
	} else {
		pe.command = "terraform"
	}

	return pe, nil
}

func (pe *TerraformExec) Plan() error {
	if err := pe.execPlan(); err != nil {
		return err
	} else if err := pe.execShow(); err != nil {
		return err
	}
	return nil
}

// Plan runs terraform for a single workDir and stores the results in the cache
func (pe *TerraformExec) execPlan() error {
	if err := os.MkdirAll(pe.cacheDir, fileMode); err != nil {
		return err
	}
	cmd := exec.Command(pe.command, "plan", "-out", pe.planBinPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = pe.workDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run '%s' from '%s': %w", cmd.String(), cmd.Dir, err)
	}
	return nil
}

// Show converts a plan to a json file that can be read by the validator
func (pe *TerraformExec) execShow() error {
	if err := os.MkdirAll(pe.cacheDir, fileMode); err != nil {
		return err
	}
	cmd := exec.Command(pe.command, "show", "-json", pe.planBinPath)
	outbuf := &bytes.Buffer{}
	cmd.Dir = pe.workDir
	cmd.Stdout = outbuf
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run '%s' from '%s': %w", cmd.String(), cmd.Dir, err)
	} else if err := ioutil.WriteFile(pe.planJsonPath, outbuf.Bytes(), fileMode); err != nil {
		return fmt.Errorf("failed to write %s: %w", pe.planJsonPath, err)
	}
	return nil
}
