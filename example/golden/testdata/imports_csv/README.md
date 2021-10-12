# imports_csv

This test demonstrates the `-proto_imports_in` flag.  A csv file is loaded that
installs resolve directives into gazelle.  Typically this file would be
generated inside an external repository using the `proto_repository` rule.

In this case we've overridden the expected location of a well-known proto.