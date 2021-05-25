// plugin package contains proto_plugin implementations that come "built-in"
// with the gazelle protoc extension. The naming convention of files is as
// follows:
//
// (1) the go filename mirrors the name the plugin is registered under
// (snake-case, so 'builtin:js:closure' -> protoc_js_closure.go).
//
// (2) if the plugin corresponds to a protoc "builtin" such as java, the name of
// the plugin is 'builtin:java' (reflects the value for the --java_out argument).
// Variants are appended to that name, such as 'builtin:js:common' (commonjs) and
// 'builtin:js:closure' (closure).
//
// (3) if the plugin is not a builtin, the name is a colon-delimited list
// reflecting its source code location.  For example, the canonical name of the
// gogofast plugin in the "github.com/gogo/protobuf" repo is
// 'gogo:protobuf:gogofast'.  In this scheme the domain 'github.com' is left
// out. The first two names are easy to choose.  The final name of the "variant"
// is just a string, and hopefully a concise and meaningful name is chosen.
package builtin
