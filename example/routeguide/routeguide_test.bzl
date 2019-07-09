def _routeguide_test_impl(ctx):
    server = None
    for f in ctx.files.server:
        if f.basename == "server.bash" or f.basename == "server" or f.basename == "server_deploy.jar":
            server = f

    if not server:
        fail("Failed to identify server entrypoint file in %r" % ctx.files.server)

    server_entrypoint = server.short_path
    if server.extension == "jar":
        server_entrypoint = "java -jar %s" % server.short_path

    client = None
    for f in ctx.files.client:
        if f.basename == "client.bash" or f.basename == "client" or f.basename == "client_deploy.jar":
            client = f

    if not client:
        fail("Failed to identify client entrypoint file in %r" % ctx.files.client)

    client_entrypoint = client.short_path
    if client.extension == "jar":
        client_entrypoint = "java -jar %s" % client.short_path

    ctx.actions.write(ctx.outputs.executable, """
set -x # Print commands
set -e # Fail on error

export DATABASE_FILE={database_file}
export SERVER_PORT={server_port}
export RUST_BACKTRACE=1 # Print rust stack traces

# Start server and wait
{server} &
sleep 2

# Run client
{client}

# Print completion for log
echo '---- DONE ----'
    """.format(
        client = client_entrypoint,
        server = server_entrypoint,
        database_file = ctx.file.database.short_path,
        server_port = ctx.attr.port,
    ), is_executable = True)

    files = ctx.files.client + ctx.files.server + ctx.files.data + [ctx.file.database]

    return [DefaultInfo(
        runfiles = ctx.runfiles(
            files = files,
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
            allow_files = True,
            cfg = "target",
        ),
        "server": attr.label(
            doc = "Server binary",
            executable = True,
            mandatory = True,
            allow_files = True,
            cfg = "target",
        ),
        "database": attr.label(
            doc = "Path to the feature database json file",
            mandatory = True,
            allow_single_file = True,
        ),
        "data": attr.label_list(
            doc = "Additional data files",
            allow_files = True,
        ),
        "port": attr.int(
            doc = "Port to use for the client/server communication (value for SERVER_PORT env var)",
            default = 50051,
        ),
    },
    test = True,
)

def get_parent_dirname(label):
    if label.startswith("//"):
        label = label[2:]
    return label.partition("/")[0]

def routeguide_test_matrix(clients = [], servers = [], database = "//example/proto:routeguide_features", tagmap = {}):
    port = 50051
    for server in servers:
        sname = get_parent_dirname(server)
        for client in clients:
            cname = get_parent_dirname(client)
            name = "%s_%s" % (cname, sname)

            tags = []
            if tagmap.get(cname):
                tags.extend(tagmap.get(cname))
            if tagmap.get(sname):
                tags.extend(tagmap.get(sname))
            if tagmap.get(name):
                tags.extend(tagmap.get(name))

            routeguide_test(
                name = name,
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
