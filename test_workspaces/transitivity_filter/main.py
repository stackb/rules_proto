# Attempt to import direct proto file, this one should succeed
import layer_pb2

# Attempt to import transitive proto file, this one should succeed
import base_pb2

# Attempt to import filtered transitive proto file, this one should not succeed
try:
    import exclude_dir.exclude_base_pb2
    raise RuntimeError('Transitively filtered base proto was built')
except ImportError:
    pass
