# Common data and functions shared by compile.bzl and aspect.bzl

_rust_keywords = [
    "as", "break", "const", "continue", "crate", "else", "enum", "extern",
    "false", "fn", "for", "if", "impl", "let", "loop", "match", "mod", "move",
    "mut", "pub", "ref", "return", "self", "Self", "static", "struct", "super",
    "trait", "true", "type", "unsafe", "use", "where", "while",
]


_objc_upper_segments = {
    "url": "URL",
    "http": "HTTP",
    "https": "HTTPS",
}


ProtoCompileInfo = provider(fields = {
    "label": "label object",
    "plugins": "ProtoPluginInfo object",
    "descriptor": "descriptor set file",
    "outputs": "generated protoc outputs",
    "files": "final generated files",
    "protos": "generated protos (copies)",
    "args": "proto arguments",
    "tools": "proto tools",
    "verbose": "verbose level",
})


def capitalize(s):
    """Capitalize a string - only first letter
    Args:
      s (string): The input string to be capitalized.
    Returns:
      (string): The capitalized string.
    """
    return s[0:1].upper() + s[1:]


def pascal_objc(s):
    """Convert pascal_case -> PascalCase

    Objective C uses pascal case, but there are e exceptions that it uppercases
    the entire segment: url, http, and https.

    https://github.com/protocolbuffers/protobuf/blob/54176b26a9be6c9903b375596b778f51f5947921/src/google/protobuf/compiler/objectivec/objectivec_helpers.cc#L91

    Args:
      s (string): The input string to be capitalized.
    Returns: (string): The capitalized string.
    """
    segments = []
    for segment in s.split("_"):
        repl = _objc_upper_segments.get(segment)
        if repl:
            segment = repl
        else:
            segment = capitalize(segment)
        segments.append(segment)
    return "".join(segments)


def pascal_case(s):
    """Convert pascal_case -> PascalCase
    Args:
        s (string): The input string to be capitalized.
    Returns:
        (string): The capitalized string.
    """
    return "".join([capitalize(part) for part in s.split("_")])


def rust_keyword(s):
    """Check if arg is a rust keyword and append '_pb' if true.
    Args:
        s (string): The input string to be capitalized.
    Returns:
        (string): The appended string.
    """
    return s + "_pb" if s in _rust_keywords else s


def describe(name, obj, exclude):
    """Print the properties of the given struct obj
    Args:
      name: the name of the struct we are introspecting.
      obj: the struct to introspect
      exclude: a list of names *not* to print (function names)
    """
    for k in dir(obj):
        if hasattr(obj, k) and k not in exclude:
            v = getattr(obj, k)
            t = type(v)
            print("%s.%s<%r> = %s" % (name, k, t, v))


def get_bool_attr(attr, name):
    value = getattr(attr, name, "False")
    return value == "True"


def get_int_attr(attr, name):
    value = getattr(attr, name)
    if value == "":
        return 0
    if value == "None":
        return 0
    return int(value)


def get_string_list_attr(attr, name):
    value = getattr(attr, name, "")
    if value == "":
        return []
    return value.split(";")


def get_output_sibling_file(pattern, proto, descriptor):
    """Get the correct place to

    The ctx.actions.declare_file has a 'sibling = <File>' feature that allows
    one to declare files in the same directory as the sibling.

    This function checks for the prefix special token '{package}' and, if true,
    uses the descriptor as the sibling (which declares the output file will be
    in the root of the generated tree).

    Args:
      pattern: the input filename pattern <string>
      proto: the .proto <Generated File> (in the staging area)
      descriptor: the descriptor <File> that marks the staging root.

    Returns:
      the <File> to be used as the correct sibling.
    """

    if pattern.startswith("{package}/"):
        return descriptor
    return proto


def get_plugin_out(label_name, plugin):
    if not plugin.out:
        return None
    filename = plugin.out
    filename = filename.replace("{name}", label_name)
    return filename


def get_plugin_runfiles(tool):
    """Gather runfiles for a plugin.
    """
    files = []
    if not tool:
        return files

    info = tool[DefaultInfo]
    if not info:
        return files

    if info.files:
        files += info.files.to_list()

    if info.default_runfiles:
        runfiles = info.default_runfiles
        if runfiles.files:
            files += runfiles.files.to_list()

    if info.data_runfiles:
        runfiles = info.data_runfiles
        if runfiles.files:
            files += runfiles.files.to_list()

    return files


def get_proto_filename(src):
    """Assemble the filename for a proto

    Args:
      src: the .proto <File>

    Returns:
      <string> of the filename.
    """
    parts = src.short_path.split("/")
    if len(parts) > 1 and parts[0] == "..":
        return "/".join(parts[2:])
    return src.short_path


def copy_jar_to_srcjar(ctx, jar):
    """Copy .jar to .srcjar

    Args:
      ctx: the <ctx> object
      jar: the <Generated File> of a jar containing source files.

    Returns:
      <Generated File> for the renamed file
    """
    srcjar = ctx.actions.declare_file("%s/%s.srcjar" % (ctx.label.name, ctx.label.name))
    ctx.actions.run_shell(
        mnemonic = "CopySrcjar",
        inputs = [jar],
        outputs = [srcjar],
        command = "mv %s %s" % (jar.path, srcjar.path),
    )
    return srcjar


def get_plugin_option(label_name, option):
    """Build a plugin option, doing plugin option template replacements if present

    Args:
      label_name: the ctx.label.name
      option: string from the <PluginInfo>

    Returns:
      <string> for the --plugin_out= arg
    """

    # TODO: use .format here and pass in a substitutions struct!
    return option.replace("{name}", label_name)


def get_plugin_options(label_name, options):
    """Build a plugin option list

    Args:
      label_name: the ctx.label.name
      options: list<string> options from the <PluginInfo>

    Returns:
      <string> for the --plugin_out= arg
    """
    return [get_plugin_option(label_name, option) for option in options]


def apply_plugin_transitivity_rules(ctx, targets, plugin):
    """Process the proto target list according to plugin transitivity rules

    Args:
      ctx: the <ctx> object
      targets: the dict<string,File> of .proto files that we intend to compile.
      plugin: the <PluginInfo> object.

    Returns:
      <list<File>> the possibly filtered list of .proto <File>s
    """

    # Iterate transitivity rules like '{ "google/protobuf": "exclude" }'. The
    # only rule type implemented is "exclude", which checks if the pathname or
    # dirname ends with the given pattern.  If so, remove that item in the
    # targets list.
    #
    # Why does this feature exist?  Well, library rules like C# require all the
    # proto files to be present during the compilation (collected via transitive
    # sources).  However, since the well-known types are already present in the
    # library dependencies, we don't actually want to compile well-known types
    # (but do want to compile everything else).
    #
    transitivity = {}
    transitivity.update(plugin.transitivity)
    transitivity.update(ctx.attr.transitivity)

    for pattern, rule in transitivity.items():
        if rule == "exclude":
            for key, target in targets.items():
                if ctx.attr.verbose > 2:
                    print("Checking '%s' endswith '%s'" % (target.short_path, pattern))
                if target.dirname.endswith(pattern) or target.path.endswith(pattern):
                    targets.pop(key)
                    if ctx.attr.verbose > 2:
                        print("Removing '%s' from the list of files to compile as plugin '%s' excluded it" % (target.short_path, plugin.name))
                elif ctx.attr.verbose > 2:
                    print("Keeping '%s' (not excluded)" % (target.short_path))
        elif rule == "include":
            for key, target in targets.items():
                if target.dirname.endswith(pattern) or target.path.endswith(pattern):
                    if ctx.attr.verbose > 2:
                        print("Keeping '%s' (explicitly included)" % (target.short_path))
                else:
                    targets.pop(key)
                    if ctx.attr.verbose > 2:
                        print("Removing '%s' from the list of files to compile as plugin '%s' did not include it" % (target.short_path, plugin.name))
        else:
            fail("Unknown transitivity rule '%s'" % rule)
    return targets


def get_output_filename(src, plugin, pattern):
    """Build the predicted filename for file generated by the given plugin.

    A 'proto_plugin' rule allows one to define the predicted outputs.  For
    flexibility, we allow special tokens in the output filename that get
    replaced here. The overall pattern is '{token}' mimicking the python
    'format' feature.

    Additionally, there are '|' characters like '{basename|pascal}' that can be
    read as 'take the basename and pipe that through the pascal function'.

    Args:
      src: the .proto <File>
      plugin: the <PluginInfo> object.
      pattern: the input pattern string

    Returns:
      the replaced string
    """

    # If output to srcjar, don't emit a per-proto output file.
    if plugin.out:
        return None

    # Slice off this prefix if it exists, we don't use it here.
    if pattern.startswith("{package}/"):
        pattern = pattern[len("{package}/"):]
    basename = src.basename
    if basename.endswith(".proto"):
        basename = basename[:-6]
    elif basename.endswith(".protodevel"):
        basename = basename[:-11]

    filename = basename

    if pattern.find("{basename}") != -1:
        filename = pattern.replace("{basename}", basename)
    elif pattern.find("{basename|pascal}") != -1:
        filename = pattern.replace("{basename|pascal}", pascal_case(basename))
    elif pattern.find("{basename|pascal|objc}") != -1:
        filename = pattern.replace("{basename|pascal|objc}", pascal_objc(basename))
    elif pattern.find("{basename|rust_keyword}") != -1:
        filename = pattern.replace("{basename|rust_keyword}", rust_keyword(basename))
    else:
        filename = basename + pattern

    return filename


def copy_proto(ctx, descriptor, src):
    """Copy a proto to the 'staging area'

    Args:
      ctx: the <ctx> object
      descriptor: the descriptor <File> that marks the root of the 'staging area'.
      src: the source .proto <File>

    Returns:
      <Generated File> for the copied .proto
    """
    proto = ctx.actions.declare_file(get_proto_filename(src), sibling = descriptor)
    ctx.actions.run_shell(
        mnemonic = "CopyProto",
        inputs = [src],
        outputs = [proto],
        command = "cp %s %s" % (src.path, proto.path),
    )
    return proto


# Shamelessly taken from https://github.com/bazelbuild/rules_go
def proto_path(proto):
    """
    The proto path is not really a file path
    It's the path to the proto that was seen when the descriptor file was generated.
    """
    path = proto.path
    root = proto.root.path
    ws = proto.owner.workspace_root
    if path.startswith(root):
        path = path[len(root):]
    if path.startswith("/"):
        path = path[1:]
    if path.startswith(ws):
        path = path[len(ws):]
    if path.startswith("/"):
        path = path[1:]
    return path
