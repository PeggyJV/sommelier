# Rust Proto Bindings Build

## Notes

For some reason, the `type_attribute` config method is not working the the `CellarIdSet` type, therefore whenever `prost_build` gets run, we have to manually add the `serde::Deserialize` and `serde::Serialize` attributes to it.
