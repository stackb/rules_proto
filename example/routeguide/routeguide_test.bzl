def _routeguide_test_impl(ctx): 
  
    server = None
    for f in ctx.files.server:
        if f.basename == "server.bash" or f.basename == "server" or f.basename == "server_deploy.jar":
            server = f

    client = None
    for f in ctx.files.client:
        if f.basename == "client.bash" or f.basename == "client" or f.basename == "client.jar":
            client = f

    if not server:
        fail("Failed to identify server entrypoint file in %r" % ctx.files.server)

    server_entrypoint = server.short_path
    if server.extension == "jar":
        server_entrypoint = "java -jar %s" % server.short_path

    ctx.actions.write(ctx.outputs.executable, """
set -x
find . | grep manifest_prep
export DATABASE_FILE={database_file}
export SERVER_PORT={server_port}
{server} &
sleep 1
{client}
    """.format(
        client = ctx.file.client.short_path,
        # server = ctx.file.server.short_path,
        server = server_entrypoint,
        database_file = ctx.file.database.short_path,
        server_port = ctx.attr.port,
    ), is_executable = True)

    return [DefaultInfo(
        runfiles = ctx.runfiles(
            files = ctx.files.client + ctx.files.server + [ctx.file.database] + ctx.files.data,
            collect_data = True,
            collect_default = True,
        ),
    )]
  
routeguide_test = rule(
    implementation = _routeguide_test_impl,
    attrs = {
        "client": attr.label(
            doc = "Client binary",
            executable = True,
            mandatory = True,
            single_file = True,
            allow_files = True,
            cfg = "target",
        ),
        "server": attr.label(
            doc = "Server binary",
            executable = True,
            mandatory = True,
            # single_file = True,
            allow_files = True,
            cfg = "target",
        ),
        "database": attr.label(
            doc = "Path to the feature database json file",
            mandatory = True,
            single_file = True,
        ),
        "data": attr.label_list(
            doc = "Additional data files",
            allow_files = True,
        ),
        "port": attr.int(
            doc = "Port to use for the client/server communication (value for SERVER_PORT env var)",
            default = 50051,
        ),
        "server_sleep": attr.int(
            doc = "Time to wait for server startup",
            default = 0,
        ),
    },
    test = True,
)

def get_parent_dirname(label):
    if label.startswith("//"):
        label = label[2:]
    segments = label.split(sep = "/", maxsplit = 2)
    return segments[0]

def routeguide_test_matrix(clients = [], servers = [], database = "//example/proto:routeguide_features"):
    port = 50051

    for server in servers:
        sname = get_parent_dirname(server)
        # print("%s -> %s" % (sname, server))
        for client in clients:
            cname = get_parent_dirname(client)
            # print("%s -> %s" % (cname, client))
            tags = []
            if cname == "csharp" or sname == "csharp":
                tags.append("no-sandbox")
                
            routeguide_test(
                name = "%s_%s_%d" % (cname, sname, port),
                client = client,
                server = server,
                database = database,
                port = port,
                data = [
                    client,
                    server,
                ],
                tags = tags,
                size = "small",
            )
            port += 1