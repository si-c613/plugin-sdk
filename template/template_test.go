/*
Copyright 2021 The terraform-docs Authors.

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package template

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
	gotemplate "text/template"

	"github.com/stretchr/testify/assert"

	"github.com/si-c613/plugin-sdk/print"
	"github.com/si-c613/plugin-sdk/terraform"
)

func TestTemplateRender(t *testing.T) {
	sectionTpl := `
	{{- with .Module.Header -}}
		{{ custom . }}
	{{- end -}}
	`
	customFuncs := gotemplate.FuncMap{
		"custom": func(s string) string {
			return fmt.Sprintf("customized <<%s>>", s)
		},
	}
	module := terraform.NewModule(
		terraform.WithHeader("sample header"),
	)
	tests := []struct {
		name     string
		items    []*Item
		expected string
		wantErr  bool
	}{
		{
			name: "template render with custom functions",
			items: []*Item{
				{
					Name: "all",
					Text: `{{- template "section" . -}}`,
				}, {
					Name: "section",
					Text: sectionTpl,
				},
			},
			expected: "customized <<sample header>>",
			wantErr:  false,
		},
		{
			name:     "template render with custom functions",
			items:    []*Item{},
			expected: "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			settings := &print.Settings{}
			tpl := New(settings, tt.items...)
			tpl.CustomFunc(customFuncs)
			rendered, err := tpl.Render(module)
			if tt.wantErr {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
				assert.Equal(tt.expected, rendered)
			}
		})
	}
}

func TestBuiltinFunc(t *testing.T) {
	tests := []struct {
		name       string
		funcName   string
		funcArgs   []string
		escapeChar bool
		escapePipe bool
		expected   string
	}{
		// default
		{
			name:       "template builtin functions default",
			funcName:   "default",
			funcArgs:   []string{`"a"`, `"b"`},
			escapeChar: true,
			escapePipe: true,
			expected:   "b",
		},
		{
			name:       "template builtin functions default",
			funcName:   "default",
			funcArgs:   []string{`"a"`, `""`},
			escapeChar: true,
			escapePipe: true,
			expected:   "a",
		},
		{
			name:       "template builtin functions default",
			funcName:   "default",
			funcArgs:   []string{`""`, `"b"`},
			escapeChar: true,
			escapePipe: true,
			expected:   "b",
		},
		{
			name:       "template builtin functions default",
			funcName:   "default",
			funcArgs:   []string{`""`, `""`},
			escapeChar: true,
			escapePipe: true,
			expected:   "",
		},

		// trim
		{
			name:       "template builtin functions trim",
			funcName:   "trim",
			funcArgs:   []string{`" "`, `"   foo   "`},
			escapeChar: true,
			escapePipe: true,
			expected:   "foo",
		},
		{
			name:       "template builtin functions trim",
			funcName:   "trim",
			funcArgs:   []string{`" "`, `"foo"`},
			escapeChar: true,
			escapePipe: true,
			expected:   "foo",
		},
		{
			name:       "template builtin functions trim",
			funcName:   "trim",
			funcArgs:   []string{`""`, `"foo"`},
			escapeChar: true,
			escapePipe: true,
			expected:   "foo",
		},
		{
			name:       "template builtin functions trim",
			funcName:   "trim",
			funcArgs:   []string{`" "`, `""`},
			escapeChar: true,
			escapePipe: true,
			expected:   "",
		},

		// trimLeft
		{
			name:       "template builtin functions trimLeft",
			funcName:   "trimLeft",
			funcArgs:   []string{`" "`, `"   foo   "`},
			escapeChar: true,
			escapePipe: true,
			expected:   "foo   ",
		},
		{
			name:       "template builtin functions trimLeft",
			funcName:   "trimLeft",
			funcArgs:   []string{`" "`, `"foo"`},
			escapeChar: true,
			escapePipe: true,
			expected:   "foo",
		},
		{
			name:       "template builtin functions trimLeft",
			funcName:   "trimLeft",
			funcArgs:   []string{`""`, `"foo"`},
			escapeChar: true,
			escapePipe: true,
			expected:   "foo",
		},
		{
			name:       "template builtin functions trimLeft",
			funcName:   "trimLeft",
			funcArgs:   []string{`" "`, `""`},
			escapeChar: true,
			escapePipe: true,
			expected:   "",
		},

		// trimRight
		{
			name:       "template builtin functions trimRight",
			funcName:   "trimRight",
			funcArgs:   []string{`" "`, `"   foo   "`},
			escapeChar: true,
			escapePipe: true,
			expected:   "   foo",
		},
		{
			name:       "template builtin functions trimRight",
			funcName:   "trimRight",
			funcArgs:   []string{`" "`, `"foo"`},
			escapeChar: true,
			escapePipe: true,
			expected:   "foo",
		},
		{
			name:       "template builtin functions trimRight",
			funcName:   "trimRight",
			funcArgs:   []string{`""`, `"foo"`},
			escapeChar: true,
			escapePipe: true,
			expected:   "foo",
		},
		{
			name:       "template builtin functions trimRight",
			funcName:   "trimRight",
			funcArgs:   []string{`" "`, `""`},
			escapeChar: true,
			escapePipe: true,
			expected:   "",
		},

		// trimPrefix
		{
			name:       "template builtin functions trimPrefix",
			funcName:   "trimPrefix",
			funcArgs:   []string{`" "`, `"   foo   "`},
			escapeChar: true,
			escapePipe: true,
			expected:   "  foo   ",
		},
		{
			name:       "template builtin functions trimPrefix",
			funcName:   "trimPrefix",
			funcArgs:   []string{`" "`, `"foo"`},
			escapeChar: true,
			escapePipe: true,
			expected:   "foo",
		},
		{
			name:       "template builtin functions trimPrefix",
			funcName:   "trimPrefix",
			funcArgs:   []string{`""`, `"foo"`},
			escapeChar: true,
			escapePipe: true,
			expected:   "foo",
		},
		{
			name:       "template builtin functions trimPrefix",
			funcName:   "trimPrefix",
			funcArgs:   []string{`" "`, `""`},
			escapeChar: true,
			escapePipe: true,
			expected:   "",
		},

		// trimSuffix
		{
			name:       "template builtin functions trimSuffix",
			funcName:   "trimSuffix",
			funcArgs:   []string{`" "`, `"   foo   "`},
			escapeChar: true,
			escapePipe: true,
			expected:   "   foo  ",
		},
		{
			name:       "template builtin functions trimSuffix",
			funcName:   "trimSuffix",
			funcArgs:   []string{`" "`, `"foo"`},
			escapeChar: true,
			escapePipe: true,
			expected:   "foo",
		},
		{
			name:       "template builtin functions trimSuffix",
			funcName:   "trimSuffix",
			funcArgs:   []string{`""`, `"foo"`},
			escapeChar: true,
			escapePipe: true,
			expected:   "foo",
		},
		{
			name:       "template builtin functions trimSuffix",
			funcName:   "trimSuffix",
			funcArgs:   []string{`" "`, `""`},
			escapeChar: true,
			escapePipe: true,
			expected:   "",
		},

		// indent
		{
			name:       "template builtin functions indent",
			funcName:   "indent",
			funcArgs:   []string{`0`, `"#"`},
			escapeChar: true,
			escapePipe: true,
			expected:   "##",
		},
		{
			name:       "template builtin functions indent",
			funcName:   "indent",
			funcArgs:   []string{`1`, `"#"`},
			escapeChar: true,
			escapePipe: true,
			expected:   "###",
		},
		{
			name:       "template builtin functions indent",
			funcName:   "indent",
			funcArgs:   []string{`2`, `"#"`},
			escapeChar: true,
			escapePipe: true,
			expected:   "####",
		},
		{
			name:       "template builtin functions indent",
			funcName:   "indent",
			funcArgs:   []string{`3`, `"#"`},
			escapeChar: true,
			escapePipe: true,
			expected:   "#####",
		},

		// name
		{
			name:       "template builtin functions name",
			funcName:   "name",
			funcArgs:   []string{`"foo"`},
			escapeChar: true,
			escapePipe: true,
			expected:   "foo",
		},
		{
			name:       "template builtin functions name",
			funcName:   "name",
			funcArgs:   []string{`"foo_bar"`},
			escapeChar: true,
			escapePipe: true,
			expected:   "foo\\_bar",
		},
		{
			name:       "template builtin functions name",
			funcName:   "name",
			funcArgs:   []string{`"foo_bar"`},
			escapeChar: false,
			escapePipe: true,
			expected:   "foo_bar",
		},
		{
			name:       "template builtin functions name",
			funcName:   "name",
			funcArgs:   []string{`""`},
			escapeChar: true,
			escapePipe: true,
			expected:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			settings := &print.Settings{}
			settings.EscapeCharacters = tt.escapeChar
			settings.EscapePipe = tt.escapePipe
			funcs := builtinFuncs(settings)

			fn, ok := funcs[tt.funcName]
			assert.Truef(ok, "function %s is not defined", tt.funcName)

			v := reflect.ValueOf(fn)
			tp := v.Type()
			assert.Equalf(len(tt.funcArgs), tp.NumIn(), "invalid number of arguments. got: %v, want: %v", len(tt.funcArgs), tp.NumIn())

			argv := make([]reflect.Value, len(tt.funcArgs))

			for i := range argv {
				var argType reflect.Kind
				if strings.HasPrefix(tt.funcArgs[i], "\"") {
					argType = reflect.String
					argv[i] = reflect.ValueOf(strings.Trim(tt.funcArgs[i], "\""))
				} else {
					argType = reflect.Int
					num, _ := strconv.Atoi(tt.funcArgs[i])
					argv[i] = reflect.ValueOf(num)
				}
				if tp.In(i).Kind() != argType {
					assert.Fail("Invalid argument. got: %v, want: %v", argType, tp.In(i).Kind())
				}
			}

			result := v.Call(argv)

			if len(result) != 1 || result[0].Kind() != reflect.String {
				assert.Fail("function %s must return a one string value", tt.funcName)
			}

			assert.Equal(tt.expected, result[0].String())
		})
	}
}

func TestSanitizeName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		escape   bool
		expected string
	}{
		{
			name:     "sanitize name with escape character",
			input:    "abcdefgh",
			escape:   true,
			expected: "abcdefgh",
		},
		{
			name:     "sanitize name with escape character",
			input:    "abcd_efgh",
			escape:   true,
			expected: "abcd\\_efgh",
		},
		{
			name:     "sanitize name with escape character",
			input:    "_abcdefgh",
			escape:   true,
			expected: "\\_abcdefgh",
		},
		{
			name:     "sanitize name with escape character",
			input:    "abcd__efgh",
			escape:   true,
			expected: "abcd\\_\\_efgh",
		},
		{
			name:     "sanitize name with escape character",
			input:    "_",
			escape:   true,
			expected: "\\_",
		},
		{
			name:     "sanitize name with escape character",
			input:    "",
			escape:   true,
			expected: "",
		},
		{
			name:     "sanitize name without escape character",
			input:    "abcdefgh",
			escape:   false,
			expected: "abcdefgh",
		},
		{
			name:     "sanitize name without escape character",
			input:    "abcd_efgh",
			escape:   false,
			expected: "abcd_efgh",
		},
		{
			name:     "sanitize name without escape character",
			input:    "abcd__efgh",
			escape:   false,
			expected: "abcd__efgh",
		},
		{
			name:     "sanitize name without escape character",
			input:    "_",
			escape:   false,
			expected: "_",
		},
		{
			name:     "sanitize name without escape character",
			input:    "",
			escape:   false,
			expected: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			settings := &print.Settings{EscapeCharacters: tt.escape}
			actual := sanitizeName(tt.input, settings)

			assert.Equal(tt.expected, actual)
		})
	}
}

func TestGenerateIndentation(t *testing.T) {
	tests := []struct {
		name     string
		base     int
		extra    int
		expected string
	}{
		{
			name:     "generate indentation",
			base:     2,
			extra:    1,
			expected: "###",
		},
		{
			name:     "generate indentation",
			extra:    2,
			expected: "####",
		},
		{
			name:     "generate indentation",
			base:     4,
			extra:    3,
			expected: "#######",
		},
		{
			name:     "generate indentation",
			base:     0,
			extra:    0,
			expected: "##",
		},
		{
			name:     "generate indentation",
			base:     6,
			extra:    1,
			expected: "###",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			settings := &print.Settings{IndentLevel: tt.base}
			actual := generateIndentation(tt.extra, "#", settings)

			assert.Equal(tt.expected, actual)
		})
	}
}
