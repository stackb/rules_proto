##

Node module structure:

The first consideration is that of the generated _pb_grpc.js file (and all
imports in general).  In this case the protocol compiler emit this require
statement: `var example_proto_routeguide_pb =
require('../../example/proto/routeguide_pb.js');`.  So, we need to build a
node_module that contains both the `_pb.js` and `_grpc_pb.js` file AND satisfies
the require paths.

These are the 2 files:

```
bazel-genfiles/node/example/routeguide/routeguide_pb/example/proto/routeguide_pb.js
bazel-genfiles/node/example/routeguide/routeguide_pb/example/proto/routeguide_grpc_pb.js
```

Therefore, we need a node_module or eqivalent with an `index.js` file in
`bazel-genfiles/node/example/`.  The `index.js` should have the content:

```
module.exports = {
    'routeguide_pb': require('./routeguide/routeguide_pb/example/proto/routeguide_pb.js'),
    'routeguide_grpc_pb': require('./routeguide/routeguide_pb/example/proto/routeguide_grpc_pb.js'),
}
```

And it should have a subdirectory `routeguide/routeguide_pb/` that contains the files.

Location of the generated node_module_index file `index.js` is arbitrary since
it gets copied to `node_module` root as it is the `index = ` attribute.
