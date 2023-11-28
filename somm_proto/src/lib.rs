pub use cosmos_sdk_proto;

pub mod cork {
    include!("prost/cork.v2.rs");
}

pub mod pubsub {
    include!("prost/pubsub.v1.rs");
}
