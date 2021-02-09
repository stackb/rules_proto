load(
    "//:repositories.bzl",
    "subpar",
    "six",
    "rules_proto_grpc_repos",
)

def python_repos(**kwargs):
    rules_proto_grpc_repos(**kwargs)
    subpar(**kwargs)
    six(**kwargs)
