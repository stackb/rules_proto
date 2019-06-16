# Attempt to import direct proto file, this one should succeed
import layer_pb2

# Attempt to import transitive proto file, this one should not succeed
try:
    import base_pb2
    raise RuntimeError('Transitive base proto was built')
except ImportError:
    pass
