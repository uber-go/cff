package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/jessevdk/go-flags"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"go.uber.org/multierr"
)

type options struct {
	DryRun      bool   `long:"dry-run" description:"Output commands that would have been run, instead of actually running them."`
	ProjectRoot string `long:"project-root" short:"p" description:"Root for project folder, part of a git repo." required:"true"`
}

func newCLIParser() (*flags.Parser, *options) {
	var opts options
	parser := flags.NewParser(&opts, flags.HelpFlag)
	parser.Name = "cff-bazelxml-to-cffcmd"

	return parser, &opts
}

type query struct {
	Rules []rule `xml:"rule"`
}

type rule struct {
	StringAttributes []bazelAttribute     `xml:"string"`
	ListAttributes   []bazelListAttribute `xml:"list"`
}

type bazelAttribute struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type bazelLabel struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type bazelListAttribute struct {
	Name   string       `xml:"name,attr"`
	Labels []bazelLabel `xml:"label"`
}

type cffRule struct {
	ImportPath string
	CFFSources []string
}

// errUnsupportedCffSource is an error that the user specified a source file that was outside the package in which the
// bazel rule was declared. While bazel does allow one to refer to files outside the current package, the CFFv2 tool
// (both the bazel rule and the tool under cmd/cff) do not seem to support this, so this tool will error if that occurs.
// Example:
// ├── a
// │   ├── BUILD.bazel
// │   └── b.go
// └── c
// │   └── d.go
//
// a/BUILD.bazel might have rule cff(cff_srcs=["b.go", "//c:d.go"], importpath="a")
type errUnsupportedCffSource struct {
	CFFSource              string
	ImportPath             string
	BazelGeneratorLocation string
}

func (e errUnsupportedCffSource) Error() string {
	return fmt.Sprintf("%s: invalid cff_srcs value %q: cannot be outside package %q", e.BazelGeneratorLocation, e.CFFSource, e.ImportPath)
}

// Hack: The Go encoding/xml package complains about xml 1.1 not being supported, but it works fine for our purpose.
func overwriteXMLVersion(xml []byte) []byte {
	return []byte(strings.Replace(string(xml), "1.1", "1.0", 1))
}

// parseXML takes the output of "bazel query 'kind(cff, ...)'" and parses it into a list of CFF rules
func parseXML(queryXML []byte) ([]cffRule, error) {
	queryXML = overwriteXMLVersion(queryXML)

	query := query{}
	err := xml.Unmarshal(queryXML, &query)
	if err != nil {
		return nil, err
	}

	cffRules := make([]cffRule, 0, len(query.Rules))

	for _, rule := range query.Rules {
		var importpath string
		var cffsrcs []string
		var bazelGeneratorLocation string

		for _, attribute := range rule.StringAttributes {
			if attribute.Name == "importpath" {
				importpath = attribute.Value
				continue
			}

			if attribute.Name == "generator_location" {
				bazelGeneratorLocation = attribute.Value
				continue
			}
		}

		if importpath == "" {
			log.Printf("found no 'importpath' attribute for cff rule")
			continue
		}

		var errs []error
		for _, attribute := range rule.ListAttributes {
			if attribute.Name == "cff_srcs" {
				for _, label := range attribute.Labels {
					prefixThisPackage := fmt.Sprintf("//src/%s:", importpath)
					if !strings.HasPrefix(label.Value, prefixThisPackage) {
						errs = append(errs, errUnsupportedCffSource{
							CFFSource:              label.Value,
							ImportPath:             importpath,
							BazelGeneratorLocation: bazelGeneratorLocation,
						})
					}

					valueWithoutPackage := strings.TrimPrefix(label.Value, prefixThisPackage)

					cffsrcs = append(cffsrcs, valueWithoutPackage)
				}
			}
		}

		if len(cffsrcs) == 0 {
			log.Printf("found no 'cff_srcs' for cff rule")
			continue
		}
		if len(errs) > 0 {
			return nil, multierr.Combine(errs...)
		}

		cffRules = append(cffRules, cffRule{
			ImportPath: importpath,
			CFFSources: cffsrcs,
		})
	}

	return cffRules, nil
}

// ruleToShellCommands transforms a CFF rule to a "bazel run" shell comamnd that can be used to invoke the CFF tool.
func ruleToShellCommands(rule cffRule) []string {
	args := []string{
		// bazel is attached by caller
		"run",
		"//src/go.uber.org/cff/cmd/cff:cff",
		"--",
	}

	for _, src := range rule.CFFSources {
		// We don't use the shell to invoke these commands, instead we directly pass these arguments to argv of
		// the cff tool, so we don't need to worry about a dangerously named file (such as one like "; rm -rf /")
		args = append(args, fmt.Sprintf("--file=%s", src))
	}

	args = append(args, rule.ImportPath)

	return args
}

func run(args []string) error {
	parser, f := newCLIParser()
	if _, err := parser.ParseArgs(args); err != nil {
		return err
	}

	bazelBin := os.Getenv("BAZEL_CMD")
	if bazelBin == "" {
		return fmt.Errorf("missing BAZEL_CMD environment variable")
	}

	xmlBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("error reading xml from stdin: %s", err.Error())
	}

	rules, err := parseXML(xmlBytes)
	if err != nil {
		return fmt.Errorf("error parsing xml: %s", err.Error())
	}

	for _, rule := range rules {
		command := ruleToShellCommands(rule)
		fmt.Println(bazelBin + " " + strings.Join(command, " "))

		if f.DryRun {
			continue
		}

		cmd := exec.Command(bazelBin, command...)
		cmd.Dir = f.ProjectRoot
		stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
		cmd.Stdout = stdout
		cmd.Stderr = stderr

		err = cmd.Start()
		if err != nil {
			return fmt.Errorf("error starting %+v: %s", cmd.Args, err.Error())
		}
		err = cmd.Wait()

		_, _ = io.Copy(os.Stdout, stdout)
		_, _ = io.Copy(os.Stderr, stderr)

		if err != nil {
			return fmt.Errorf("error running %+v: %s", cmd.Args, err.Error())
		}
	}

	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatalf("%+v", err)
	}
}
