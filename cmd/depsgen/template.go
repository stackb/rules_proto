package main

// templateData is the type used by the template
type templateData struct {
	Name string
	Deps []*dependency
}

var depsBzl = `"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def {{ .Name }}_deps():{{ range .Deps }}{{ if ne .Dep.RepositoryRule "phony" }}
    {{ .Dep.Name }}()  # via {{ .ParentName }}{{ end }}{{ end }}
{{ range .Deps }}{{ if ne .Dep.RepositoryRule "phony" }}

def {{ .Dep.Name }}():
    _maybe(
        {{ .Dep.RepositoryRule }},
        name = "{{ .Dep.Name }}",{{ if .Dep.Path }}
        path = "{{ .Dep.Path }}",{{ end }}{{ if .Dep.StripPrefix }}
        sha256 = "{{ .Dep.Sha256 }}",{{ end }}{{ if .Dep.StripPrefix }}
        strip_prefix = "{{ .Dep.StripPrefix }}",{{ end }}{{ if .Dep.Urls }}
        urls = [{{ range .Dep.Urls }}
            "{{ . }}",{{ end }}
        ],{{ end }}{{ if .Dep.BuildFile }}
        build_file = "{{ .Dep.BuildFile }}",{{ end }}{{ if .Dep.BuildFileContent }}
        build_file_content = """{{ .Dep.BuildFileContent }}""",{{ end }}
    ){{ end }}{{ end }}
`
