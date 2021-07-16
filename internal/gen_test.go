package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintImportAlias(t *testing.T) {
	tests := []struct {
		importPath             string
		alias                  string
		wantImportPathMapValue string
		wantAlias              string
	}{
		{
			importPath:             "some/pkg/foo",
			alias:                  "foo",
			wantImportPathMapValue: "",
			wantAlias:              "foo",
		},
		{
			importPath:             "another/pkg/foo",
			alias:                  "foo",
			wantImportPathMapValue: "_foo",
			wantAlias:              "_foo",
		},
		{
			importPath:             "some/pkg/my-pkg",
			alias:                  "mypkg",
			wantImportPathMapValue: "mypkg",
			wantAlias:              "mypkg",
		},
		{
			importPath:             "another/pkg/mypkg",
			alias:                  "mypkg",
			wantImportPathMapValue: "_mypkg",
			wantAlias:              "_mypkg",
		},
		{
			importPath:             "yet/another/pkg/foo",
			alias:                  "foo",
			wantImportPathMapValue: "__foo",
			wantAlias:              "__foo",
		},
	}

	addImports := make(map[string]string)
	aliases := make(map[string]struct{})

	for _, tt := range tests {
		gotAlias := printImportAlias(tt.importPath, tt.alias, addImports, aliases)
		assert.Equal(t, tt.wantAlias, gotAlias)
		assert.Contains(t, addImports, tt.importPath)
		assert.Equal(t, tt.wantImportPathMapValue, addImports[tt.importPath])
		assert.Contains(t, aliases, tt.alias)
	}
}
