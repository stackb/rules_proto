def _execute(rtx, cmds):
    """Execute a command and fail if return code.
    Args:
      rtx: !repository_ctx
      cmds: !list<string>
    Returns: struct value from the rtx.execute method.
    """

    #print("Execute <%s>" % " ".join(cmds))
    result = rtx.execute(cmds)
    if result.return_code:
        fail(" ".join(cmds) + "failed: %s" % (result.stderr))
    return result

def _pub_repository(name, entry, verbose):
    out = []
    version = entry["version"]
    if entry.get("override"):
        override = entry.get("override")

        # print("%s %s override %s" % (name, version, override))
        version = override

    if version.startswith("^"):
        version = version[1:]

    if version.startswith(">"):
        version = version[1:]

    if version.startswith("<"):
        version = version[1:]

    if version.startswith("="):
        version = version[1:]

    if version.startswith(">="):
        version = version[2:]

    if version.startswith("<="):
        version = version[2:]

    if version.startswith("=="):
        version = version[2:]
    out += [
        '    if "vendor_%s" not in existing:' % name,
        "        pub_repository(",
        '            name = "vendor_%s",' % name,
        '            output = ".",',
        '            package = "%s",' % name,
        '            version = "%s",' % version,
    ]

    deps = entry.get("deps")
    if deps:
        out.append("            pub_deps = [")
        for depname, depversion in deps.items():
            out.append('                "%s",' % depname)
        out.append("            ],")
    out.append("        )")
    out.append("    elif verbose > 0:")
    out.append('        print("Skipped vendor_%s (already exists)")' % name)
    if verbose > 1:
        print("%s: %s" % (name, version))
    return out

def _dart_pub_deps_impl(rtx):
    """
    Repository rule implementation that (1) copies the pubspec.yaml into
    the external workspace dir, (2) executes 'pub deps', (3) parses the output,
    and (4) writes out a deps.bzl file that contains pub_repository rules.
    """

    # map[string]string name -> version overrides
    override = rtx.attr.override

    # int
    verbose = rtx.attr.verbose

    # string
    pub = rtx.path(rtx.attr._pub)

    # string
    spec = rtx.path(rtx.attr.spec)

    # Copy the pubspec into place
    _execute(rtx, ["cp", spec, "./pubspec.yaml"])

    # Run pub get first
    _execute(rtx, [pub, "get"])

    # Run pub deps
    result = _execute(rtx, [pub, "deps", "--style", "list"])

    # name -> entry
    direct_deps = {}
    # name -> entry

    transitive_deps = {}

    # The dict that is "active" (one of: deps | transitive_deps).  A dict gets
    # 'activated' when we hit that section in the output.
    active = None

    # Name of last dependency we've seen, like "analyzer" in ' - analyzer 0.32.5
    current = None

    # Iterate all lines in output
    lines = result.stdout.split("\n")
    for line in lines:
        # print("LINE: " + line)
        if line.startswith("  - "):
            toks = line[4:].split(" ")
            name = toks[0]
            version = toks[1]
            active[name] = {
                "name": name,
                "version": version,
                "override": override.get(name),
            }
            deps = current["deps"]
            if not deps.get(name):
                deps[name] = version
            if verbose > 1:
                print("transitive dep %s: %s" % (name, version))
        elif line.startswith("- "):
            toks = line[2:].split(" ")
            name = toks[0]
            version = toks[1]
            entry = {
                "name": name,
                "version": version,
                "deps": {},
                "override": override.get(name),
            }
            active[name] = entry
            current = entry
            if verbose > 1:
                print("direct dep %s: %s" % (name, version))
        elif line == "dependencies:":
            active = direct_deps
        elif line == "transitive dependencies:":
            active = transitive_deps
        elif verbose > 2:
            print("SKIP: " + line)

    out = [
        "# Generated - do not modify",
        'load("@io_bazel_rules_dart//dart/build_rules/internal:pub.bzl", "pub_repository")',
        "def pub_deps(verbose = 0):",
        "    existing = native.existing_rules()",
    ]

    for name, entry in direct_deps.items():
        out += _pub_repository(name, entry, verbose)
    for name, entry in transitive_deps.items():
        out += _pub_repository(name, entry, verbose)

    rtx.file("deps.bzl", "\n".join(out))
    rtx.file("BUILD.bazel", "")

dart_pub_deps = repository_rule(
    implementation = _dart_pub_deps_impl,
    attrs = {
        "_pub": attr.label(
            doc = "The pub binary tool",
            default = "@dart_sdk//:bin/pub",
        ),
        "spec": attr.label(
            doc = "A pubspec.yaml file that details the dependencies",
            mandatory = True,
        ),
        "override": attr.string_dict(
            doc = 'A mapping from NAME -> VERSION such that the given VERSION will be chosen for any direct/transitive dependency with that name.  Example: {"glob": "1.1.7"}',
        ),
        "verbose": attr.int(
            doc = "A number that changes the verbose level",
        ),
    },
)
