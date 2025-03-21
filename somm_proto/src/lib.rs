pub use cosmos_sdk_proto;

pub mod cork {
    include!("prost/cork.v2.rs");
}

pub mod pubsub {
    include!("prost/pubsub.v1.rs");
}

pub mod axelar_cork {
    include!("prost/axelarcork.v1.rs");
}

pub mod incentives {
    include!("prost/incentives.v1.rs");
}

pub mod cellarfees {
    include!("prost/cellarfees.v1.rs");
}
