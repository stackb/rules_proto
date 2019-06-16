test_workspace_proto_source_root:
	cd test_workspaces/proto_source_root; \
	bazel test --disk_cache=../bazel-disk-cache //... ; \
	bazel shutdown

test_workspace_shared_proto:
	cd test_workspaces/shared_proto; \
	bazel test --disk_cache=../bazel-disk-cache //... ; \
	bazel shutdown

test_workspace_strip_import_prefix:
	cd test_workspaces/strip_import_prefix; \
	bazel test --disk_cache=../bazel-disk-cache //... ; \
	bazel shutdown

all_test_workspaces: test_workspace_proto_source_root test_workspace_shared_proto test_workspace_strip_import_prefix
