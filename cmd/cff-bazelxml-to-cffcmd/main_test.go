package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseXML(t *testing.T) {
	input := []byte(`<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<query version="2">
    <rule class="_cff_generate" location="/Users/jacobg/go-code/src/code.uber.internal/presentation/apps/rider/helix/endpoints/surge/batch/BUILD.bazel:55:1" name="//src/code.uber.internal/presentation/apps/rider/helix/endpoints/surge/batch:cff">
        <string name="name" value="cff"/>
        <list name="visibility">
            <label value="//visibility:public"/>
        </list>
        <string name="generator_name" value="cff"/>
        <string name="generator_function" value="cff"/>
        <string name="generator_location" value="src/code.uber.internal/presentation/apps/rider/helix/endpoints/surge/batch/BUILD.bazel:55"/>
        <list name="deps">
            <label value="//idl/code.uber.internal/cme/rt-products:rt_productsv2_go_thriftrw"/>
            <label value="//idl/code.uber.internal/presentation/apps/rider/helix/endpoints:surge_go_thriftrw"/>
            <label value="//idl/code.uber.internal/presentation/models:exception_go_thriftrw"/>
            <label value="//idl/code.uber.internal/presentation/models:references_go_thriftrw"/>
            <label value="//src/go.uber.org/cff:go_default_library"/>
            <label value="//src/code.uber.internal/marketplace/fulfillment-platform.git/platform/headers:go_default_library"/>
            <label value="//src/code.uber.internal/presentation/apps/rider/tasks/heatpipe/publish:go_default_library"/>
            <label value="//src/code.uber.internal/presentation/apps/rider/tasks/heatpipe/publish-batch:go_default_library"/>
            <label value="//src/code.uber.internal/presentation/apps/rider/tasks/product/get-active-products:go_default_library"/>
            <label value="//src/code.uber.internal/presentation/apps/rider/tasks/routing/delegate-header:go_default_library"/>
            <label value="//src/code.uber.internal/presentation/apps/rider/tasks/surge/build-surge-message:go_default_library"/>
            <label value="//src/code.uber.internal/presentation/apps/rider/tasks/surge/publish-to-control-tower:go_default_library"/>
            <label value="//src/code.uber.internal/presentation/build/apps/rider/helix/endpoints/surge/module:go_default_library"/>
            <label value="//src/code.uber.internal/presentation/build/apps/rider/helix/endpoints/surge/workflow:go_default_library"/>
            <label value="//src/code.uber.internal/presentation/build/clients/rtapi-flipr:go_default_library"/>
            <label value="//src/code.uber.internal/presentation/lib/headers:go_default_library"/>
            <label value="//src/code.uber.internal/presentation/lib/pointer:go_default_library"/>
            <label value="@com_github_pkg_errors//:go_default_library"/>
            <label value="@com_github_uber_go_tally//:go_default_library"/>
            <label value="@com_github_uber_zanzibar//runtime:go_default_library"/>
            <label value="@org_uber_go_zap//:go_default_library"/>
        </list>
        <label name="gopath" value="//src/code.uber.internal/presentation/apps/rider/helix/endpoints/surge/batch:cff_cff_gopath"/>
        <string name="importpath" value="code.uber.internal/presentation/apps/rider/helix/endpoints/surge/batch"/>
        <list name="srcs"/>
        <list name="cff_srcs">
            <label value="//src/code.uber.internal/presentation/apps/rider/helix/endpoints/surge/batch:handler.go"/>
        </list>
        <rule-input name="//idl/code.uber.internal/cme/rt-products:rt_productsv2_go_thriftrw"/>
        <rule-input name="//idl/code.uber.internal/presentation/apps/rider/helix/endpoints:surge_go_thriftrw"/>
        <rule-input name="//idl/code.uber.internal/presentation/models:exception_go_thriftrw"/>
        <rule-input name="//idl/code.uber.internal/presentation/models:references_go_thriftrw"/>
        <rule-input name="//src/go.uber.org/cff:go_default_library"/>
        <rule-input name="//src/go.uber.org/cff/cmd/cff:cff"/>
        <rule-input name="//src/code.uber.internal/marketplace/fulfillment-platform.git/platform/headers:go_default_library"/>
        <rule-input name="//src/code.uber.internal/presentation/apps/rider/helix/endpoints/surge/batch:cff_cff_gopath"/>
        <rule-input name="//src/code.uber.internal/presentation/apps/rider/helix/endpoints/surge/batch:handler.go"/>
        <rule-input name="//src/code.uber.internal/presentation/apps/rider/tasks/heatpipe/publish:go_default_library"/>
        <rule-input name="//src/code.uber.internal/presentation/apps/rider/tasks/heatpipe/publish-batch:go_default_library"/>
        <rule-input name="//src/code.uber.internal/presentation/apps/rider/tasks/product/get-active-products:go_default_library"/>
        <rule-input name="//src/code.uber.internal/presentation/apps/rider/tasks/routing/delegate-header:go_default_library"/>
        <rule-input name="//src/code.uber.internal/presentation/apps/rider/tasks/surge/build-surge-message:go_default_library"/>
        <rule-input name="//src/code.uber.internal/presentation/apps/rider/tasks/surge/publish-to-control-tower:go_default_library"/>
        <rule-input name="//src/code.uber.internal/presentation/build/apps/rider/helix/endpoints/surge/module:go_default_library"/>
        <rule-input name="//src/code.uber.internal/presentation/build/apps/rider/helix/endpoints/surge/workflow:go_default_library"/>
        <rule-input name="//src/code.uber.internal/presentation/build/clients/rtapi-flipr:go_default_library"/>
        <rule-input name="//src/code.uber.internal/presentation/lib/headers:go_default_library"/>
        <rule-input name="//src/code.uber.internal/presentation/lib/pointer:go_default_library"/>
        <rule-input name="@com_github_pkg_errors//:go_default_library"/>
        <rule-input name="@com_github_uber_go_tally//:go_default_library"/>
        <rule-input name="@com_github_uber_zanzibar//runtime:go_default_library"/>
        <rule-input name="@io_bazel_rules_go//:go_context_data"/>
        <rule-input name="@io_bazel_rules_go//:stdlib"/>
        <rule-input name="@io_bazel_rules_go//go/tools/coverdata:coverdata"/>
        <rule-input name="@io_bazel_rules_nogo//:nogo"/>
        <rule-input name="@org_uber_go_zap//:go_default_library"/>
    </rule>
</query>`)
	rules, err := parseXML(input)
	require.NoError(t, err)
	require.Equal(t, []cffRule{
		{
			ImportPath: "code.uber.internal/presentation/apps/rider/helix/endpoints/surge/batch",
			CFFSources: []string{"handler.go"},
		},
	}, rules)
}

func TestParseXMLInvalid(t *testing.T) {
	_, err := parseXML([]byte("invalid"))
	require.Error(t, err)
}

func TestParseXMLMissingImportPath(t *testing.T) {
	input := []byte(`<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<query version="2">
    <rule class="_cff_generate">
        <string name="name" value="cff"/>
    </rule>
</query>`)
	rules, err := parseXML(input)
	require.NoError(t, err)
	require.Equal(t, 0, len(rules))
}

func TestParseXMLOutsideCFFSource(t *testing.T) {
	input := []byte(`<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<query version="2">
    <rule class="_cff_generate">
        <string name="name" value="cff"/>
        <string name="importpath" value="demo"/>
        <string name="generator_location" value="src/demo/BUILD.bazel:55"/>
        <list name="cff_srcs">
            <label value="//:foo.go"/>
        </list>
    </rule>
</query>`)
	_, err := parseXML(input)
	require.Error(t, err)
	require.Equal(t, "src/demo/BUILD.bazel:55: invalid cff_srcs value \"//:foo.go\": cannot be outside package \"demo\"", err.Error())
}

func TestParseXMLNoCFFSources(t *testing.T) {
	input := []byte(`<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<query version="2">
    <rule class="_cff_generate">
        <string name="name" value="cff"/>
        <string name="importpath" value="demo"/>
        <string name="generator_location" value="src/demo/BUILD.bazel:55"/>
        <list name="cff_srcs"></list>
    </rule>
</query>`)
	rules, err := parseXML(input)
	require.NoError(t, err)
	require.Equal(t, 0, len(rules))
}

func TestRuleToShellCommands(t *testing.T) {
	input := cffRule{
		ImportPath: "code.uber.internal/presentation/apps/rider/helix/endpoints/surge/batch",
		CFFSources: []string{"handler.go", "other.go"},
	}

	output := ruleToShellCommands(input)

	require.Equal(t, []string{
		"run",
		"//src/go.uber.org/cff/cmd/cff:cff",
		"--",
		"--file=handler.go",
		"--file=other.go",
		"code.uber.internal/presentation/apps/rider/helix/endpoints/surge/batch",
	}, output)
}

// Test passing a file name that attempts to trick us into arbitrary code execution with shell control characters like ;
func TestRuleShellCommandShellControlSecurity(t *testing.T) {
	input := cffRule{
		ImportPath: "code.uber.internal/presentation/apps/rider/helix/endpoints/surge/batch",
		CFFSources: []string{";sleep 3600"},
	}

	output := ruleToShellCommands(input)

	require.Equal(t, []string{
		"run",
		"//src/go.uber.org/cff/cmd/cff:cff",
		"--",
		"--file=;sleep 3600",
		"code.uber.internal/presentation/apps/rider/helix/endpoints/surge/batch",
	}, output)
}
