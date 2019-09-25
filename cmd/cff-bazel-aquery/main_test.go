package main

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"proto/go.uber.org/cff/cmd/cff-bazel-aquery/proto/analysis"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func readAquery() (*analysis.ActionGraphContainer, error) {
	fd, err := os.Open("testdata/aquery.dat.gz")
	if err != nil {
		return nil, err
	}
	reader, err := gzip.NewReader(fd)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	query, err := parseProtoActionQuery(bytes)
	if err != nil {
		return nil, err
	}
	return query, nil
}

// TestFixtureQuery tests the accuracy of the aquery fixture.
func TestFixtureActionQuery(t *testing.T) {
	query, err := readAquery()

	require.NoError(t, err)
	assert.NotEmpty(t, query.Artifacts)
	lastArtifact := query.Artifacts[len(query.Artifacts)-1]
	assert.Equal(t, "bazel-out/darwin-fastbuild/bin/src/go.uber.org/cff/cmd/cff-bazel-aquery/testdata/example/darwin_amd64_stripped/cff%/test_gen.go", lastArtifact.ExecPath)
	assert.Equal(t, "6025", lastArtifact.Id)
	assert.Equal(t, []string{"6025"}, query.Actions[0].OutputIds)
	assert.Equal(t, "0", query.Actions[0].TargetId)
	assert.Equal(t, "0", query.Targets[0].Id)
	assert.Equal(t, "0", query.Targets[0].RuleClassId)
	assert.Equal(t, "_cff_generate", query.RuleClasses[0].Name)
	assert.Equal(t, "//src/go.uber.org/cff/cmd/cff-bazel-aquery/testdata/example:cff", query.Targets[0].Label)
}

// TestActionQueriesByLabel tests that we can pull out the generated files for a given rule
func TestActionQueriesByLabel(t *testing.T) {
	query, err := readAquery()
	require.NoError(t, err)

	ag := newAgraph(query)

	outputs := outputFilesForLabel("//src/go.uber.org/cff/cmd/cff-bazel-aquery/testdata/example:cff", ag)

	require.Equal(t, []string{
		"bazel-out/darwin-fastbuild/bin/src/go.uber.org/cff/cmd/cff-bazel-aquery/testdata/example/darwin_amd64_stripped/cff%/test_gen.go",
	}, outputs)
}

// TestActionQueriesByLabel tests the invalid label case for parsing the output graph
func TestActionQueriesByLabelInvalid(t *testing.T) {
	query, err := readAquery()
	require.NoError(t, err)

	ag := newAgraph(query)

	outputs := outputFilesForLabel("invalid", ag)

	require.Equal(t, []string(nil), outputs)
}

func TestTransformCFFLabel(t *testing.T) {
	outputPath := cffLabelToPath("go-path", "//src/go.uber.org/cff/cmd/cff-bazel-aquery/testdata/example:cff")
	assert.Equal(t, "go-path/src/go.uber.org/cff/cmd/cff-bazel-aquery/testdata/example", outputPath)
}

func mustMarshal(g *analysis.ActionGraphContainer) []byte {
	bytes, err := proto.Marshal(g)
	if err != nil {
		panic(err)
	}

	return bytes
}

type filesystemFixture struct {
	T              *testing.T
	ContainingPath string
}

func (f *filesystemFixture) RemoveOrFatal() {
	err := os.RemoveAll(f.ContainingPath)
	if err != nil {
		f.T.Fatalf("error removing temporary directory %q: %v", f.ContainingPath, err)
	}
}

func newFilesystemFixture(t *testing.T, files map[string][]byte) (*filesystemFixture, error) {
	tmp, err := ioutil.TempDir("", "cff-bazel-aquery")
	if err != nil {
		return nil, err
	}

	for name, content := range files {
		dir := filepath.Dir(filepath.Join(tmp, name))
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return &filesystemFixture{t, tmp}, err
		}

		err = ioutil.WriteFile(filepath.Join(tmp, name), content, 0644)
		if err != nil {
			return &filesystemFixture{t, tmp}, err
		}
	}

	return &filesystemFixture{t, tmp}, nil
}

func TestRun(t *testing.T) {
	f, err := newFilesystemFixture(t, map[string][]byte{
		"out/foo/foo_gen.go": []byte(`package foo_gen`),
		"src/foo/foo.go":     []byte(`package foo`),
	})
	defer f.RemoveOrFatal()
	assert.NoError(t, err)
	if err != nil {
		return
	}

	aqueryEncoded := mustMarshal(&analysis.ActionGraphContainer{
		Targets: []*analysis.Target{
			{
				Id:          "0",
				Label:       "//src/foo:cff",
				RuleClassId: "2",
			},
		},
		Artifacts: []*analysis.Artifact{
			{
				ExecPath: "out/foo/foo_gen.go",
				Id:       "1",
			},
		},
		Actions: []*analysis.Action{
			{
				TargetId:  "0",
				OutputIds: []string{"1"},
			},
		},
		RuleClasses: []*analysis.RuleClass{
			{
				Id:   "2",
				Name: "_cff_generate",
			},
		},
	})

	stdin := bytes.NewBuffer(aqueryEncoded)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	opts := new(options)
	opts.Args.ProjectRoot = f.ContainingPath
	err = run(opts, stdin, stdout, stderr)
	assert.NoError(t, err)
	outLines := strings.Split(stdout.String(), "\n")

	expectedInputPath := filepath.Join(f.ContainingPath, "out/foo/foo_gen.go")
	expectedOutputPath := filepath.Join(f.ContainingPath, "src/foo/foo_gen.go")
	assert.Equal(t, []string{"//src/foo:cff\t" + expectedInputPath + "\t" + expectedOutputPath, ""}, outLines)
	errLines := strings.Split(stderr.String(), "\n")
	assert.Equal(t, []string{""}, errLines)

	bytes, err := ioutil.ReadFile(expectedOutputPath)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`package foo_gen`), bytes)
}
