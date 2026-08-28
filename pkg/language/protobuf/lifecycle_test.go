package protobuf

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUpdateRootCargoDependencies(t *testing.T) {
	repoRoot := t.TempDir()
	cargoToml := filepath.Join(repoRoot, "Cargo.toml")
	const initial = `[workspace.dependencies]
# gazelle:proto_rust_dependencies start
stale = { path = "stale" }
# gazelle:proto_rust_dependencies end
`
	if err := os.WriteFile(cargoToml, []byte(initial), 0o644); err != nil {
		t.Fatal(err)
	}

	err := updateRootCargoDependencies(repoRoot, []cargoPathDependency{
		{Name: "z_proto", Path: "z/_rust/z_proto"},
		{Name: "a_proto", Path: "a/_rust/a_proto"},
		{Name: "a_proto", Path: "a/_rust/a_proto"},
	})
	if err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(cargoToml)
	if err != nil {
		t.Fatal(err)
	}
	want := strings.ReplaceAll(initial, `stale = { path = "stale" }`, `a_proto = { path = "a/_rust/a_proto" }
z_proto = { path = "z/_rust/z_proto" }`)
	if string(got) != want {
		t.Fatalf("Cargo.toml mismatch\ngot:\n%s\nwant:\n%s", got, want)
	}
}
