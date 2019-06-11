# Common data and functions shared by compile.bzl and aspect.bzl

rust_keywords = {
    "as": True,
    "break": True,
    "const": True,
    "continue": True,
    "crate": True,
    "else": True,
    "enum": True,
    "extern": True,
    "false": True,
    "fn": True,
    "for": True,
    "if": True,
    "impl": True,
    "let": True,
    "loop": True,
    "match": True,
    "mod": True,
    "move": True,
    "mut": True,
    "pub": True,
    "ref": True,
    "return": True,
    "self": True,
    "Self": True,
    "static": True,
    "struct": True,
    "super": True,
    "trait": True,
    "true": True,
    "type": True,
    "unsafe": True,
    "use": True,
    "where": True,
    "while": True,
}


objc_upper_segments = {
    "url": "URL",
    "http": "HTTP",
    "https": "HTTPS",
}


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
        repl = objc_upper_segments.get(segment)
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
    return s + "_pb" if rust_keywords.get(s) else s


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