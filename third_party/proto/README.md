# Third party vendored protos

The primary proto generation now uses buf and the buf schema registry, but some of these vendored protos are still necessary because the Rust code generator, found under `prost_build` in the root directory, needs them.

