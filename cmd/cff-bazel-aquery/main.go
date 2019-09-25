package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/jessevdk/go-flags"
	analysispb "proto/go.uber.org/cff/cmd/cff-bazel-aquery/proto/analysis"
)

const (
	cffRuleName   = "_cff_generate"
	srcPrefix     = "//"
	cffRuleSuffix = ":cff"
)

type options struct {
	DryRun bool `long:"dry-run" description:"Output commands that would have been run, instead of actually running them."`
	Args   struct {
		ProjectRoot string `positional-arg-name:"project-root" description:"Root for project folder, part of a git repo." required:"true"`
	} `positional-args:"yes" required:"yes"`
}

func newCLIParser() (*flags.Parser, *options) {
	var opts options
	parser := flags.NewParser(&opts, flags.HelpFlag)
	parser.Name = "cff-bazel-aquery"

	return parser, &opts
}

func parseProtoActionQuery(bytes []byte) (*analysispb.ActionGraphContainer, error) {
	var res analysispb.ActionGraphContainer
	err := proto.Unmarshal(bytes, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

type agraph struct {
	g *analysispb.ActionGraphContainer

	ruleClasses map[string]*analysispb.RuleClass // rule class id -> rule class
	artifacts   map[string]*analysispb.Artifact  // artifact id -> artifact
	targets     map[string]*analysispb.Target    // label -> target
	actions     map[string][]*analysispb.Action  // target id -> list of actions
}

func newAgraph(g *analysispb.ActionGraphContainer) *agraph {
	ag := &agraph{g: g}
	ag.ruleClasses = ag.indexRulesByClassID()
	ag.artifacts = ag.indexArtifactsByID()
	ag.targets = ag.indexTargetByLabel()
	ag.actions = ag.indexActionsByTargetID()
	return ag
}

func (ag *agraph) indexRulesByClassID() map[string]*analysispb.RuleClass {
	ruleClassByID := make(map[string]*analysispb.RuleClass)
	for _, ruleClass := range ag.g.RuleClasses {
		ruleClassByID[ruleClass.Id] = ruleClass
	}

	return ruleClassByID
}

func (ag *agraph) indexArtifactsByID() map[string]*analysispb.Artifact {
	artifactByID := make(map[string]*analysispb.Artifact)
	for _, artifact := range ag.g.Artifacts {
		artifactByID[artifact.Id] = artifact
	}

	return artifactByID
}

func (ag *agraph) indexTargetByLabel() map[string]*analysispb.Target {
	targetByLabel := make(map[string]*analysispb.Target)
	for _, target := range ag.g.Targets {
		targetByLabel[target.Label] = target
	}

	return targetByLabel
}

func (ag *agraph) indexActionsByTargetID() map[string][]*analysispb.Action {
	actionsByTargetID := make(map[string][]*analysispb.Action)
	for _, action := range ag.g.Actions {
		actionsByTargetID[action.TargetId] = append(actionsByTargetID[action.TargetId], action)
	}

	return actionsByTargetID
}

func (ag *agraph) LookupRuleClassByID(ruleClassID string) *analysispb.RuleClass {
	return ag.ruleClasses[ruleClassID]
}

func (ag *agraph) LookupArtifactByID(artifactID string) *analysispb.Artifact {
	return ag.artifacts[artifactID]
}

func (ag *agraph) LookupTargetByLabel(label string) *analysispb.Target {
	return ag.targets[label]
}

func (ag *agraph) LookupActionsByTargetID(targetID string) []*analysispb.Action {
	return ag.actions[targetID]
}

// outputFilesForLabel parses the bazel action graph, returning the list of declared output files for a given target
// relative to bazel-out.
func outputFilesForLabel(label string, ag *agraph) []string {
	target := ag.LookupTargetByLabel(label)
	if target == nil {
		return nil
	}

	actions := ag.LookupActionsByTargetID(target.Id)

	var outputs []string
	for _, action := range actions {
		for _, artifactID := range action.OutputIds {
			outputs = append(outputs, ag.LookupArtifactByID(artifactID).ExecPath)
		}
	}

	return outputs
}

// cffLabelToPath transforms a label like //src/foo/bar:cff to "foo/bar" prefixed by the gopath
func cffLabelToPath(bazelProjectRoot, label string) string {
	base := strings.TrimPrefix(strings.TrimSuffix(label, cffRuleSuffix), srcPrefix)
	return filepath.Join(bazelProjectRoot, base)
}

func run(flags *options, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	protoQueryBytes, err := ioutil.ReadAll(stdin)
	if err != nil {
		return fmt.Errorf("error reading from stdin: %v", err)
	}

	query, err := parseProtoActionQuery(protoQueryBytes)
	if err != nil {
		return fmt.Errorf("error parsing proto: %v", err)
	}

	ag := newAgraph(query)

	for _, target := range ag.g.Targets {
		ruleClass := ag.LookupRuleClassByID(target.RuleClassId)
		if ruleClass == nil || ruleClass.Name != cffRuleName {
			continue
		}

		label := target.Label
		destinationDirectory := cffLabelToPath(flags.Args.ProjectRoot, label)
		outputs := outputFilesForLabel(label, ag)

		for _, output := range outputs {
			inputPath := filepath.Join(flags.Args.ProjectRoot, output)
			destinationPath := filepath.Join(destinationDirectory, filepath.Base(output))

			if _, err := os.Stat(inputPath); err != nil {
				fmt.Fprintf(stderr, "expected output file %q from rule %q to exist at %q\n", output, label, inputPath)
				continue
			}

			bytes, err := ioutil.ReadFile(inputPath)
			if err != nil {
				return fmt.Errorf("error reading %q from rule %q: %v", inputPath, label, err)
			}

			fmt.Fprintf(stdout, "%s\t%s\t%s\n", label, inputPath, destinationPath)

			if flags.DryRun {
				continue
			}

			if err := ioutil.WriteFile(destinationPath, bytes, 0644); err != nil {
				return fmt.Errorf("error copying %q from rule %q to %q: %v", inputPath, label, destinationPath, err)
			}
		}
	}

	return nil
}

func main() {
	parser, flags := newCLIParser()
	if _, err := parser.ParseArgs(os.Args[1:]); err != nil {
		log.Fatalf("%+v", err)
	}

	if err := run(flags, os.Stdin, os.Stdout, os.Stderr); err != nil {
		log.Fatalf("%+v", err)
	}
}
