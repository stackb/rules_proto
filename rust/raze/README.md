Rust protobuf and gRPC dependencies managed by cargo raze: https://github.com/google/cargo-raze

To update a dependency:

- Update the `Cargo.toml` file with the new dependency version.
- Follow the instructions for installing raze at the link above (the "Remote Dependency Mode" option). If you have Cargo installed, this should just be `cargo install cargo-raze`.
- `cd` into this directory (`rust/raze`).
- Run `cargo raze`.
