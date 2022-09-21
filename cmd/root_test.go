package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/kylelemons/godebug/diff"
)

const SEPARATOR = "\n----------------------------------------\n"

func osPipe() (*os.File, *os.File) {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	return r, w
}

func loadTestData(file string) string {
	data, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return string(data)
}

type cmdCase struct {
	name   string
	args   []string
	stdout string
	stderr string
	files  map[string]string
}

func (tc *cmdCase) run(t *testing.T) {
	t.Run(tc.name, func(t *testing.T) {

		stderr_r, stderr_w := osPipe()
		stdout_r, stdout_w := osPipe()

		c := New()
		c.SetArgs(tc.args)
		c.SetErr(stderr_w)
		c.SetOut(stdout_w)
		_ = c.Execute()

		stdout_w.Close()
		stderr_w.Close()

		var stderr string
		var stderr_buf bytes.Buffer
		var stdout string
		var stdout_buf bytes.Buffer
		io.Copy(&stderr_buf, stderr_r)
		stderr = strings.TrimSpace(stderr_buf.String())
		io.Copy(&stdout_buf, stdout_r)
		stdout = strings.TrimSpace(stdout_buf.String())

		expectedParts := []string{"stdout", tc.stdout, "stderr", tc.stderr}
		for k, v := range tc.files {
			expectedParts = append(expectedParts, k, v)
		}
		expected := strings.Join(expectedParts, SEPARATOR)

		actualParts := []string{"stdout", stdout, "stderr", stderr}
		for k, _ := range tc.files {
			actualParts = append(actualParts, k, loadTestData(k))
		}
		actual := strings.Join(actualParts, SEPARATOR)

		if actual != expected {
			t.Errorf("Result not as expected:\n%v", diff.Diff(expected, actual))
		}
	})
}
