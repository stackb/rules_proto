def _routeguide_test_impl(ctx): 
  
    # server = None
    # for f in ctx.files.server:
    #     if f.basename == "server.bash" or f.basename == "server":
    #         server = f
    # client = None
    # for f in ctx.files.client:
    #     if f.basename == "client.bash" or f.basename == "client":
    #         client = f

    server_entrypoint = ctx.file.server.short_path
    if ctx.file.server.extension == "jar":
        server_entrypoint = "java -jar %s" % ctx.file.server.short_path

    ctx.actions.write(ctx.outputs.executable, """
set -x
find .
ls -al .
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
            files = ctx.files.client + ctx.files.server + [ctx.file.database],
            collect_data = True,
        ),
    )]
  
    # ctx.actions.run(
    #     mnemonic = "RouteguideTest",
    #     progress_message = "%s vs %s" % (ctx.file.client.short_path, ctx.file.server.short_path),
    #     inputs = [
    #         ctx.file.client,
    #         ctx.file.server,
    #         ctx.file.database,
    #     ],
    #     outputs = [
    #         ctx.outputs.stdout,
    #         ctx.outputs.stderr,
    #     ],
    #     executable = ctx.outputs.executable,
    #     env = {
    #         # "CLIENT": ctx.file.client.path,
    #         # "SERVER": ctx.file.server.path,
    #         # "STDOUT_FILE": ctx.outputs.stdout.path,
    #         # "STDERR_FILE": ctx.outputs.stderr.path,
    #         "DATABASE_FILE": ctx.file.database.path,
    #         "SLEEP": "%d" % ctx.attr.server_sleep,
    #         "SERVER_PORT": "%d" % ctx.attr.port,
    #     },
    #     use_default_shell_env = False,
    #     # command = """
    #     # set -x
    #     # "${SERVER}" &
    #     # sleep "${SLEEP}"
    #     # "${CLIENT}"
    #     # """,
    # )

    # # ctx.actions.run_shell(
    # #     mnemonic = "RouteguideTest",
    # #     progress_message = "%s vs %s" % (ctx.file.client.short_path, ctx.file.server.short_path),
    # #     inputs = [
    # #         ctx.file.client,
    # #         ctx.file.server,
    # #         ctx.file.database,
    # #     ],
    # #     outputs = [
    # #         ctx.outputs.executable,
    # #     ],
    # #     env = {
    # #         "CLIENT": ctx.file.client.path,
    # #         "SERVER": ctx.file.server.path,
    # #         "DATABASE_FILE": ctx.file.database.path,
    # #         "SERVER_PORT": "%d" % ctx.attr.port,
    # #         "SLEEP": "%d" % ctx.attr.server_sleep,
    # #     },
    # #     use_default_shell_env = False,
    # #     command = """
    # #     set -x
    # #     "${SERVER}" &
    # #     sleep "${SLEEP}"
    # #     "${CLIENT}"
    # #     """,
    # # )

    # return [DefaultInfo(
    #     runfiles = ctx.runfiles(collect_data = True),
    # )]

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
            single_file = True,
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
    # outputs = {
    #     "stdout": "%{name}.stdout",
    #     "stderr": "%{name}.stderr",
    # },
    test = True,
)