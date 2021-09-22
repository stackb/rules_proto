package main

// templateData is the type used by the template
type templateData struct {
	Name  string
	Deps  []*dependency
	Loads []*LoadInfo
}

var depsBzl = `"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

{{ range .Loads }}
load("{{ .Label }}"{{ range .Symbols }}, "{{ . }}"{{end}}){{end}}

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
        path = "{{ .Dep.Path }}",{{ end }}{{ if .Dep.PackageJson }}
        package_json = "{{ .Dep.PackageJson }}",{{ end }}{{ if .Dep.PackageLockJson }}
        package_lock_json = "{{ .Dep.PackageLockJson }}",{{ end }}{{ if .Dep.YarnLock }}
        yarn_lock = "{{ .Dep.YarnLock }}",{{ end }}{{ if .Dep.Sha256 }}
        sha256 = "{{ .Dep.Sha256 }}",{{ end }}{{ if .Dep.StripPrefix }}
        strip_prefix = "{{ .Dep.StripPrefix }}",{{ end }}{{ if .Dep.PackageJson }}
        symlink_node_modules = {{ if .Dep.SymlinkNodeModules }}True{{ else }}False{{ end }},{{ end }}{{ if .Dep.Sum }}
        sum = "{{ .Dep.Sum }}",{{ end }}{{ if .Dep.Version }}
        version = "{{ .Dep.Version }}",{{ end }}{{ if .Dep.Importpath }}
        importpath = "{{ .Dep.Importpath }}",{{ end }}{{ if .Dep.BuildFileProtoMode }}
        build_file_proto_mode = "{{ .Dep.BuildFileProtoMode }}",{{ end }}{{ if .Dep.Urls }}
        urls = [{{ range .Dep.Urls }}
            "{{ . }}",{{ end }}
        ],{{ end }}{{ if .Dep.BuildFile }}
        build_file = "{{ .Dep.BuildFile }}",{{ end }}{{ if .Dep.BuildFileContent }}
        build_file_content = """{{ .Dep.BuildFileContent }}""",{{ end }}
    ){{ end }}{{ end }}
`
