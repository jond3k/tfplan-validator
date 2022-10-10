package tfplan_validator

// import "github.com/mattn/go-zglob"

var defaultGlobs = []string{
	"**/*.terraform.lock.hcl",
	"**/main.tf",
	"!modules/**",
}

func findWorkspaces(globs []string) ([]string, error) {
	return []string{}, nil
	// out := map[string]bool{}
	// for _, glob := range globs {
	// 	if files, err := zglob.Glob(glob); err != nil {
	// 		return nil, err
	// 	}

	// }
}
