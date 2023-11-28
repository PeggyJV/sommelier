/// represents a publisher, which are added via governance
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct Publisher {
    /// account address of the publisher
    #[prost(string, tag = "1")]
    pub address: ::prost::alloc::string::String,
    /// unique key, FQDN of the publisher, max length of 256
    #[prost(string, tag = "2")]
    pub domain: ::prost::alloc::string::String,
    /// the publisher's self-signed CA cert PEM file, expecting TLS 1.3 compatible ECDSA certificates, max length 4096
    #[prost(string, tag = "3")]
    pub ca_cert: ::prost::alloc::string::String,
}
/// represents a subscriber, can be set or modified by the owner of the subscriber address
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct Subscriber {
    /// unique key, account address representation of either an account or a validator
    #[prost(string, tag = "1")]
    pub address: ::prost::alloc::string::String,
    // the below fields are optional, and only required if the subscriber wants to use "push" PublisherIntents
    /// the subscriber's self-signed CA cert PEM file, expecting TLS 1.3 compatible ECDSA certificates, max length 4096
    #[prost(string, tag = "2")]
    pub ca_cert: ::prost::alloc::string::String,
    /// max length of 512
    #[prost(string, tag = "3")]
    pub push_url: ::prost::alloc::string::String,
}
/// represents a publisher committing to sending messages for a specific subscription ID
///
/// unique key is subscription_id and publisher_domain tuple
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct PublisherIntent {
    /// arbitary string representing a subscription, max length of 128
    #[prost(string, tag = "1")]
    pub subscription_id: ::prost::alloc::string::String,
    /// FQDN of the publisher, max length of 256
    #[prost(string, tag = "2")]
    pub publisher_domain: ::prost::alloc::string::String,
    /// either PULL or PUSH (see enum above for details)
    #[prost(enumeration = "PublishMethod", tag = "3")]
    pub method: i32,
    /// optional, only needs to be set if using the PULL method, max length of 512
    #[prost(string, tag = "4")]
    pub pull_url: ::prost::alloc::string::String,
    /// either ANY, VALIDATORS, or LIST (see enum above for details)
    #[prost(enumeration = "AllowedSubscribers", tag = "5")]
    pub allowed_subscribers: i32,
    /// optional, must be provided if allowed_subscribers is LIST, list of account addresses, max length 256
    #[prost(string, repeated, tag = "6")]
    pub allowed_addresses: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
}
/// represents a subscriber requesting messages for a specific subscription ID and publisher
///
/// unique key is subscription_id and subscriber_address tuple, a given subscriber can only subscribe to one publisher per
/// subscription_id at a time
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct SubscriberIntent {
    /// arbitary string representing a subscription, max length of 128
    #[prost(string, tag = "1")]
    pub subscription_id: ::prost::alloc::string::String,
    /// account address of the subscriber
    #[prost(string, tag = "2")]
    pub subscriber_address: ::prost::alloc::string::String,
    /// FQDN of the publisher, max length of 256
    #[prost(string, tag = "3")]
    pub publisher_domain: ::prost::alloc::string::String,
}
/// represents a default subscription voted in by governance that can be overridden by a subscriber
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct DefaultSubscription {
    /// arbitary string representing a subscription, max length of 128
    #[prost(string, tag = "1")]
    pub subscription_id: ::prost::alloc::string::String,
    /// FQDN of the publisher, max length of 256
    #[prost(string, tag = "2")]
    pub publisher_domain: ::prost::alloc::string::String,
}
/// governance proposal to add a publisher, with domain, adress, and ca_cert the same as the Publisher type
/// proof URL expected in the format: https://<domain>/<address>/cacert.pem and serving cacert.pem matching ca_cert
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct AddPublisherProposal {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(string, tag = "3")]
    pub domain: ::prost::alloc::string::String,
    #[prost(string, tag = "4")]
    pub address: ::prost::alloc::string::String,
    #[prost(string, tag = "5")]
    pub proof_url: ::prost::alloc::string::String,
    #[prost(string, tag = "6")]
    pub ca_cert: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct AddPublisherProposalWithDeposit {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(string, tag = "3")]
    pub domain: ::prost::alloc::string::String,
    #[prost(string, tag = "4")]
    pub address: ::prost::alloc::string::String,
    #[prost(string, tag = "5")]
    pub proof_url: ::prost::alloc::string::String,
    #[prost(string, tag = "6")]
    pub ca_cert: ::prost::alloc::string::String,
    #[prost(string, tag = "7")]
    pub deposit: ::prost::alloc::string::String,
}
/// governance proposal to remove a publisher (publishers can remove themselves, but this might be necessary in the
/// event of a malicious publisher or a key compromise), since Publishers are unique by domain, it's the only
/// necessary information to remove one
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct RemovePublisherProposal {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(string, tag = "3")]
    pub domain: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct RemovePublisherProposalWithDeposit {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(string, tag = "3")]
    pub domain: ::prost::alloc::string::String,
    #[prost(string, tag = "4")]
    pub deposit: ::prost::alloc::string::String,
}
/// set the default publisher for a given subscription ID
/// these can be overridden by the client
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct AddDefaultSubscriptionProposal {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(string, tag = "3")]
    pub subscription_id: ::prost::alloc::string::String,
    #[prost(string, tag = "4")]
    pub publisher_domain: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct AddDefaultSubscriptionProposalWithDeposit {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(string, tag = "3")]
    pub subscription_id: ::prost::alloc::string::String,
    #[prost(string, tag = "4")]
    pub publisher_domain: ::prost::alloc::string::String,
    #[prost(string, tag = "5")]
    pub deposit: ::prost::alloc::string::String,
}
/// remove a default subscription
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct RemoveDefaultSubscriptionProposal {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(string, tag = "3")]
    pub subscription_id: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct RemoveDefaultSubscriptionProposalWithDeposit {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(string, tag = "3")]
    pub subscription_id: ::prost::alloc::string::String,
    #[prost(string, tag = "4")]
    pub deposit: ::prost::alloc::string::String,
}
/// for a given PublisherIntent, whether or not it is pulled or pushed
#[derive(
    serde::Deserialize,
    serde::Serialize,
    Clone,
    Copy,
    Debug,
    PartialEq,
    Eq,
    Hash,
    PartialOrd,
    Ord,
    ::prost::Enumeration,
)]
#[repr(i32)]
pub enum PublishMethod {
    /// subscribers should pull from the provided URL
    Pull = 0,
    /// subscribers must provide a URL to receive push messages
    Push = 1,
}
/// for a given PublisherIntent, determines what types of subscribers may subscribe
#[derive(
    serde::Deserialize,
    serde::Serialize,
    Clone,
    Copy,
    Debug,
    PartialEq,
    Eq,
    Hash,
    PartialOrd,
    Ord,
    ::prost::Enumeration,
)]
#[repr(i32)]
pub enum AllowedSubscribers {
    /// any valid account address
    Any = 0,
    /// account address must map to a validator in the active validator set
    Validators = 1,
    /// a specific list of account addresses
    List = 2,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct MsgRemovePublisherRequest {
    #[prost(string, tag = "1")]
    pub publisher_domain: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub signer: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct MsgRemovePublisherResponse {}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct MsgAddSubscriberRequest {
    #[prost(message, optional, tag = "1")]
    pub subscriber: ::core::option::Option<Subscriber>,
    #[prost(string, tag = "2")]
    pub signer: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct MsgAddSubscriberResponse {}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct MsgRemoveSubscriberRequest {
    #[prost(string, tag = "1")]
    pub subscriber_address: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub signer: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct MsgRemoveSubscriberResponse {}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct MsgAddPublisherIntentRequest {
    #[prost(message, optional, tag = "1")]
    pub publisher_intent: ::core::option::Option<PublisherIntent>,
    #[prost(string, tag = "2")]
    pub signer: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct MsgAddPublisherIntentResponse {}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct MsgRemovePublisherIntentRequest {
    #[prost(string, tag = "1")]
    pub subscription_id: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub publisher_domain: ::prost::alloc::string::String,
    #[prost(string, tag = "3")]
    pub signer: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct MsgRemovePublisherIntentResponse {}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct MsgAddSubscriberIntentRequest {
    #[prost(message, optional, tag = "1")]
    pub subscriber_intent: ::core::option::Option<SubscriberIntent>,
    #[prost(string, tag = "2")]
    pub signer: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct MsgAddSubscriberIntentResponse {}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct MsgRemoveSubscriberIntentRequest {
    #[prost(string, tag = "1")]
    pub subscription_id: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub subscriber_address: ::prost::alloc::string::String,
    #[prost(string, tag = "3")]
    pub signer: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct MsgRemoveSubscriberIntentResponse {}
#[doc = r" Generated client implementations."]
pub mod msg_client {
    #![allow(unused_variables, dead_code, missing_docs)]
    use tonic::codegen::*;
    pub struct MsgClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl MsgClient<tonic::transport::Channel> {
        #[doc = r" Attempt to create a new client by connecting to a given endpoint."]
        pub async fn connect<D>(dst: D) -> Result<Self, tonic::transport::Error>
        where
            D: std::convert::TryInto<tonic::transport::Endpoint>,
            D::Error: Into<StdError>,
        {
            let conn = tonic::transport::Endpoint::new(dst)?.connect().await?;
            Ok(Self::new(conn))
        }
    }
    impl<T> MsgClient<T>
    where
        T: tonic::client::GrpcService<tonic::body::BoxBody>,
        T::ResponseBody: Body + HttpBody + Send + 'static,
        T::Error: Into<StdError>,
        <T::ResponseBody as HttpBody>::Error: Into<StdError> + Send,
    {
        pub fn new(inner: T) -> Self {
            let inner = tonic::client::Grpc::new(inner);
            Self { inner }
        }
        pub fn with_interceptor(inner: T, interceptor: impl Into<tonic::Interceptor>) -> Self {
            let inner = tonic::client::Grpc::with_interceptor(inner, interceptor);
            Self { inner }
        }
        pub async fn remove_publisher(
            &mut self,
            request: impl tonic::IntoRequest<super::MsgRemovePublisherRequest>,
        ) -> Result<tonic::Response<super::MsgRemovePublisherResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/pubsub.v1.Msg/RemovePublisher");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn add_subscriber(
            &mut self,
            request: impl tonic::IntoRequest<super::MsgAddSubscriberRequest>,
        ) -> Result<tonic::Response<super::MsgAddSubscriberResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/pubsub.v1.Msg/AddSubscriber");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn remove_subscriber(
            &mut self,
            request: impl tonic::IntoRequest<super::MsgRemoveSubscriberRequest>,
        ) -> Result<tonic::Response<super::MsgRemoveSubscriberResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/pubsub.v1.Msg/RemoveSubscriber");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn add_publisher_intent(
            &mut self,
            request: impl tonic::IntoRequest<super::MsgAddPublisherIntentRequest>,
        ) -> Result<tonic::Response<super::MsgAddPublisherIntentResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/pubsub.v1.Msg/AddPublisherIntent");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn remove_publisher_intent(
            &mut self,
            request: impl tonic::IntoRequest<super::MsgRemovePublisherIntentRequest>,
        ) -> Result<tonic::Response<super::MsgRemovePublisherIntentResponse>, tonic::Status>
        {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/pubsub.v1.Msg/RemovePublisherIntent");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn add_subscriber_intent(
            &mut self,
            request: impl tonic::IntoRequest<super::MsgAddSubscriberIntentRequest>,
        ) -> Result<tonic::Response<super::MsgAddSubscriberIntentResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/pubsub.v1.Msg/AddSubscriberIntent");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn remove_subscriber_intent(
            &mut self,
            request: impl tonic::IntoRequest<super::MsgRemoveSubscriberIntentRequest>,
        ) -> Result<tonic::Response<super::MsgRemoveSubscriberIntentResponse>, tonic::Status>
        {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path =
                http::uri::PathAndQuery::from_static("/pubsub.v1.Msg/RemoveSubscriberIntent");
            self.inner.unary(request.into_request(), path, codec).await
        }
    }
    impl<T: Clone> Clone for MsgClient<T> {
        fn clone(&self) -> Self {
            Self {
                inner: self.inner.clone(),
            }
        }
    }
    impl<T> std::fmt::Debug for MsgClient<T> {
        fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
            write!(f, "MsgClient {{ ... }}")
        }
    }
}
/// Params defines the parameters for the module.
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct Params {}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryParamsRequest {}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryParamsResponse {
    #[prost(message, optional, tag = "1")]
    pub params: ::core::option::Option<Params>,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryPublisherRequest {
    #[prost(string, tag = "1")]
    pub publisher_domain: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryPublisherResponse {
    #[prost(message, optional, tag = "1")]
    pub publisher: ::core::option::Option<Publisher>,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryPublishersRequest {}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryPublishersResponse {
    #[prost(message, repeated, tag = "1")]
    pub publishers: ::prost::alloc::vec::Vec<Publisher>,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QuerySubscriberRequest {
    #[prost(string, tag = "1")]
    pub subscriber_address: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QuerySubscriberResponse {
    #[prost(message, optional, tag = "1")]
    pub subscriber: ::core::option::Option<Subscriber>,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QuerySubscribersRequest {}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QuerySubscribersResponse {
    #[prost(message, repeated, tag = "1")]
    pub subscribers: ::prost::alloc::vec::Vec<Subscriber>,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryPublisherIntentRequest {
    #[prost(string, tag = "1")]
    pub publisher_domain: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub subscription_id: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryPublisherIntentResponse {
    #[prost(message, optional, tag = "1")]
    pub publisher_intent: ::core::option::Option<PublisherIntent>,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryPublisherIntentsRequest {}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryPublisherIntentsResponse {
    #[prost(message, repeated, tag = "1")]
    pub publisher_intents: ::prost::alloc::vec::Vec<PublisherIntent>,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryPublisherIntentsByPublisherDomainRequest {
    #[prost(string, tag = "1")]
    pub publisher_domain: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryPublisherIntentsByPublisherDomainResponse {
    #[prost(message, repeated, tag = "1")]
    pub publisher_intents: ::prost::alloc::vec::Vec<PublisherIntent>,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryPublisherIntentsBySubscriptionIdRequest {
    #[prost(string, tag = "1")]
    pub subscription_id: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryPublisherIntentsBySubscriptionIdResponse {
    #[prost(message, repeated, tag = "1")]
    pub publisher_intents: ::prost::alloc::vec::Vec<PublisherIntent>,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QuerySubscriberIntentRequest {
    #[prost(string, tag = "1")]
    pub subscriber_address: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub subscription_id: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QuerySubscriberIntentResponse {
    #[prost(message, optional, tag = "1")]
    pub subscriber_intent: ::core::option::Option<SubscriberIntent>,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QuerySubscriberIntentsRequest {}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QuerySubscriberIntentsResponse {
    #[prost(message, repeated, tag = "1")]
    pub subscriber_intents: ::prost::alloc::vec::Vec<SubscriberIntent>,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QuerySubscriberIntentsBySubscriberAddressRequest {
    #[prost(string, tag = "1")]
    pub subscriber_address: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QuerySubscriberIntentsBySubscriberAddressResponse {
    #[prost(message, repeated, tag = "1")]
    pub subscriber_intents: ::prost::alloc::vec::Vec<SubscriberIntent>,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QuerySubscriberIntentsBySubscriptionIdRequest {
    #[prost(string, tag = "1")]
    pub subscription_id: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QuerySubscriberIntentsBySubscriptionIdResponse {
    #[prost(message, repeated, tag = "1")]
    pub subscriber_intents: ::prost::alloc::vec::Vec<SubscriberIntent>,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QuerySubscriberIntentsByPublisherDomainRequest {
    #[prost(string, tag = "1")]
    pub publisher_domain: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QuerySubscriberIntentsByPublisherDomainResponse {
    #[prost(message, repeated, tag = "1")]
    pub subscriber_intents: ::prost::alloc::vec::Vec<SubscriberIntent>,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryDefaultSubscriptionRequest {
    #[prost(string, tag = "1")]
    pub subscription_id: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryDefaultSubscriptionResponse {
    #[prost(message, optional, tag = "1")]
    pub default_subscription: ::core::option::Option<DefaultSubscription>,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryDefaultSubscriptionsRequest {}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryDefaultSubscriptionsResponse {
    #[prost(message, repeated, tag = "1")]
    pub default_subscriptions: ::prost::alloc::vec::Vec<DefaultSubscription>,
}
#[doc = r" Generated client implementations."]
pub mod query_client {
    #![allow(unused_variables, dead_code, missing_docs)]
    use tonic::codegen::*;
    pub struct QueryClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl QueryClient<tonic::transport::Channel> {
        #[doc = r" Attempt to create a new client by connecting to a given endpoint."]
        pub async fn connect<D>(dst: D) -> Result<Self, tonic::transport::Error>
        where
            D: std::convert::TryInto<tonic::transport::Endpoint>,
            D::Error: Into<StdError>,
        {
            let conn = tonic::transport::Endpoint::new(dst)?.connect().await?;
            Ok(Self::new(conn))
        }
    }
    impl<T> QueryClient<T>
    where
        T: tonic::client::GrpcService<tonic::body::BoxBody>,
        T::ResponseBody: Body + HttpBody + Send + 'static,
        T::Error: Into<StdError>,
        <T::ResponseBody as HttpBody>::Error: Into<StdError> + Send,
    {
        pub fn new(inner: T) -> Self {
            let inner = tonic::client::Grpc::new(inner);
            Self { inner }
        }
        pub fn with_interceptor(inner: T, interceptor: impl Into<tonic::Interceptor>) -> Self {
            let inner = tonic::client::Grpc::with_interceptor(inner, interceptor);
            Self { inner }
        }
        pub async fn params(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryParamsRequest>,
        ) -> Result<tonic::Response<super::QueryParamsResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/pubsub.v1.Query/Params");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn query_publisher(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryPublisherRequest>,
        ) -> Result<tonic::Response<super::QueryPublisherResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/pubsub.v1.Query/QueryPublisher");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn query_publishers(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryPublishersRequest>,
        ) -> Result<tonic::Response<super::QueryPublishersResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/pubsub.v1.Query/QueryPublishers");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn query_subscriber(
            &mut self,
            request: impl tonic::IntoRequest<super::QuerySubscriberRequest>,
        ) -> Result<tonic::Response<super::QuerySubscriberResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/pubsub.v1.Query/QuerySubscriber");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn query_subscribers(
            &mut self,
            request: impl tonic::IntoRequest<super::QuerySubscribersRequest>,
        ) -> Result<tonic::Response<super::QuerySubscribersResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/pubsub.v1.Query/QuerySubscribers");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn query_publisher_intent(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryPublisherIntentRequest>,
        ) -> Result<tonic::Response<super::QueryPublisherIntentResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path =
                http::uri::PathAndQuery::from_static("/pubsub.v1.Query/QueryPublisherIntent");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn query_publisher_intents(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryPublisherIntentsRequest>,
        ) -> Result<tonic::Response<super::QueryPublisherIntentsResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path =
                http::uri::PathAndQuery::from_static("/pubsub.v1.Query/QueryPublisherIntents");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn query_publisher_intents_by_publisher_domain(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryPublisherIntentsByPublisherDomainRequest>,
        ) -> Result<
            tonic::Response<super::QueryPublisherIntentsByPublisherDomainResponse>,
            tonic::Status,
        > {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pubsub.v1.Query/QueryPublisherIntentsByPublisherDomain",
            );
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn query_publisher_intents_by_subscription_id(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryPublisherIntentsBySubscriptionIdRequest>,
        ) -> Result<
            tonic::Response<super::QueryPublisherIntentsBySubscriptionIdResponse>,
            tonic::Status,
        > {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pubsub.v1.Query/QueryPublisherIntentsBySubscriptionID",
            );
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn query_subscriber_intent(
            &mut self,
            request: impl tonic::IntoRequest<super::QuerySubscriberIntentRequest>,
        ) -> Result<tonic::Response<super::QuerySubscriberIntentResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path =
                http::uri::PathAndQuery::from_static("/pubsub.v1.Query/QuerySubscriberIntent");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn query_subscriber_intents(
            &mut self,
            request: impl tonic::IntoRequest<super::QuerySubscriberIntentsRequest>,
        ) -> Result<tonic::Response<super::QuerySubscriberIntentsResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path =
                http::uri::PathAndQuery::from_static("/pubsub.v1.Query/QuerySubscriberIntents");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn query_subscriber_intents_by_subscriber_address(
            &mut self,
            request: impl tonic::IntoRequest<super::QuerySubscriberIntentsBySubscriberAddressRequest>,
        ) -> Result<
            tonic::Response<super::QuerySubscriberIntentsBySubscriberAddressResponse>,
            tonic::Status,
        > {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pubsub.v1.Query/QuerySubscriberIntentsBySubscriberAddress",
            );
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn query_subscriber_intents_by_subscription_id(
            &mut self,
            request: impl tonic::IntoRequest<super::QuerySubscriberIntentsBySubscriptionIdRequest>,
        ) -> Result<
            tonic::Response<super::QuerySubscriberIntentsBySubscriptionIdResponse>,
            tonic::Status,
        > {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pubsub.v1.Query/QuerySubscriberIntentsBySubscriptionID",
            );
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn query_subscriber_intents_by_publisher_domain(
            &mut self,
            request: impl tonic::IntoRequest<super::QuerySubscriberIntentsByPublisherDomainRequest>,
        ) -> Result<
            tonic::Response<super::QuerySubscriberIntentsByPublisherDomainResponse>,
            tonic::Status,
        > {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pubsub.v1.Query/QuerySubscriberIntentsByPublisherDomain",
            );
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn query_default_subscription(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryDefaultSubscriptionRequest>,
        ) -> Result<tonic::Response<super::QueryDefaultSubscriptionResponse>, tonic::Status>
        {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path =
                http::uri::PathAndQuery::from_static("/pubsub.v1.Query/QueryDefaultSubscription");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn query_default_subscriptions(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryDefaultSubscriptionsRequest>,
        ) -> Result<tonic::Response<super::QueryDefaultSubscriptionsResponse>, tonic::Status>
        {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path =
                http::uri::PathAndQuery::from_static("/pubsub.v1.Query/QueryDefaultSubscriptions");
            self.inner.unary(request.into_request(), path, codec).await
        }
    }
    impl<T: Clone> Clone for QueryClient<T> {
        fn clone(&self) -> Self {
            Self {
                inner: self.inner.clone(),
            }
        }
    }
    impl<T> std::fmt::Debug for QueryClient<T> {
        fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
            write!(f, "QueryClient {{ ... }}")
        }
    }
}
/// GenesisState defines the pubsub module's genesis state.
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct GenesisState {
    #[prost(message, optional, tag = "1")]
    pub params: ::core::option::Option<Params>,
    #[prost(message, repeated, tag = "2")]
    pub publishers: ::prost::alloc::vec::Vec<Publisher>,
    #[prost(message, repeated, tag = "3")]
    pub subscribers: ::prost::alloc::vec::Vec<Subscriber>,
    #[prost(message, repeated, tag = "4")]
    pub publisher_intents: ::prost::alloc::vec::Vec<PublisherIntent>,
    #[prost(message, repeated, tag = "5")]
    pub subscriber_intents: ::prost::alloc::vec::Vec<SubscriberIntent>,
    #[prost(message, repeated, tag = "6")]
    pub default_subscriptions: ::prost::alloc::vec::Vec<DefaultSubscription>,
}
