package version

import (
	"os"
	"runtime"
	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const versionTemplate = `{{.Name}}:
 Version: {{.Version}}
 Go Version: {{.GoVersion}}
 Built: {{.BuildTime}}(Rev: {{.GitRevision}})
`

func Command(name, version, buildTime, revision string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version info",
		RunE: func(cmd *cobra.Command, _ []string) error {
			data := map[string]string{
				"Name":        name,
				"Version":     version,
				"GoVersion":   runtime.Version(),
				"BuildTime":   buildTime,
				"GitRevision": revision,
			}

			tpl, err := template.New("version").Parse(versionTemplate)
			if err != nil {
				return errors.Wrap(err, "parse version template")
			}

			return tpl.Execute(os.Stdout, data)
		},
	}
}