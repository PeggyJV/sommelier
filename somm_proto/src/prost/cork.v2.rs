// This file is @generated by prost-build.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Cork {
    /// call body containing the ABI encoded bytes to send to the contract
    #[prost(bytes = "vec", tag = "1")]
    pub encoded_contract_call: ::prost::alloc::vec::Vec<u8>,
    /// address of the contract to send the call
    #[prost(string, tag = "2")]
    pub target_contract_address: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ScheduledCork {
    #[prost(message, optional, tag = "1")]
    pub cork: ::core::option::Option<Cork>,
    #[prost(uint64, tag = "2")]
    pub block_height: u64,
    #[prost(string, tag = "3")]
    pub validator: ::prost::alloc::string::String,
    #[prost(bytes = "vec", tag = "4")]
    pub id: ::prost::alloc::vec::Vec<u8>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct CorkResult {
    #[prost(message, optional, tag = "1")]
    pub cork: ::core::option::Option<Cork>,
    #[prost(uint64, tag = "2")]
    pub block_height: u64,
    #[prost(bool, tag = "3")]
    pub approved: bool,
    #[prost(string, tag = "4")]
    pub approval_percentage: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct CellarIdSet {
    #[prost(string, repeated, tag = "1")]
    pub ids: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
}
/// MsgScheduleCorkRequest - sdk.Msg for scheduling a cork request for on or after a specific block height
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgScheduleCorkRequest {
    /// the scheduled cork
    #[prost(message, optional, tag = "1")]
    pub cork: ::core::option::Option<Cork>,
    /// the block height that must be reached
    #[prost(uint64, tag = "2")]
    pub block_height: u64,
    /// signer account address
    #[prost(string, tag = "3")]
    pub signer: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgScheduleCorkResponse {
    /// cork ID
    #[prost(string, tag = "1")]
    pub id: ::prost::alloc::string::String,
}
/// Generated client implementations.
pub mod msg_client {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value,
    )]
    use tonic::codegen::*;
    use tonic::codegen::http::Uri;
    /// MsgService defines the msgs that the cork module handles
    #[derive(Debug, Clone)]
    pub struct MsgClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl MsgClient<tonic::transport::Channel> {
        /// Attempt to create a new client by connecting to a given endpoint.
        pub async fn connect<D>(dst: D) -> Result<Self, tonic::transport::Error>
        where
            D: TryInto<tonic::transport::Endpoint>,
            D::Error: Into<StdError>,
        {
            let conn = tonic::transport::Endpoint::new(dst)?.connect().await?;
            Ok(Self::new(conn))
        }
    }
    impl<T> MsgClient<T>
    where
        T: tonic::client::GrpcService<tonic::body::BoxBody>,
        T::Error: Into<StdError>,
        T::ResponseBody: Body<Data = Bytes> + std::marker::Send + 'static,
        <T::ResponseBody as Body>::Error: Into<StdError> + std::marker::Send,
    {
        pub fn new(inner: T) -> Self {
            let inner = tonic::client::Grpc::new(inner);
            Self { inner }
        }
        pub fn with_origin(inner: T, origin: Uri) -> Self {
            let inner = tonic::client::Grpc::with_origin(inner, origin);
            Self { inner }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> MsgClient<InterceptedService<T, F>>
        where
            F: tonic::service::Interceptor,
            T::ResponseBody: Default,
            T: tonic::codegen::Service<
                http::Request<tonic::body::BoxBody>,
                Response = http::Response<
                    <T as tonic::client::GrpcService<tonic::body::BoxBody>>::ResponseBody,
                >,
            >,
            <T as tonic::codegen::Service<
                http::Request<tonic::body::BoxBody>,
            >>::Error: Into<StdError> + std::marker::Send + std::marker::Sync,
        {
            MsgClient::new(InterceptedService::new(inner, interceptor))
        }
        /// Compress requests with the given encoding.
        ///
        /// This requires the server to support it otherwise it might respond with an
        /// error.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.send_compressed(encoding);
            self
        }
        /// Enable decompressing responses.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.accept_compressed(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_decoding_message_size(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_encoding_message_size(limit);
            self
        }
        pub async fn schedule_cork(
            &mut self,
            request: impl tonic::IntoRequest<super::MsgScheduleCorkRequest>,
        ) -> std::result::Result<
            tonic::Response<super::MsgScheduleCorkResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/cork.v2.Msg/ScheduleCork");
            let mut req = request.into_request();
            req.extensions_mut().insert(GrpcMethod::new("cork.v2.Msg", "ScheduleCork"));
            self.inner.unary(req, path, codec).await
        }
    }
}
/// GenesisState - all cork state that must be provided at genesis
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GenesisState {
    #[prost(message, optional, tag = "1")]
    pub params: ::core::option::Option<Params>,
    #[prost(message, optional, tag = "2")]
    pub cellar_ids: ::core::option::Option<CellarIdSet>,
    #[prost(uint64, tag = "3")]
    pub invalidation_nonce: u64,
    #[prost(message, repeated, tag = "4")]
    pub scheduled_corks: ::prost::alloc::vec::Vec<ScheduledCork>,
    #[prost(message, repeated, tag = "5")]
    pub cork_results: ::prost::alloc::vec::Vec<CorkResult>,
}
/// Params cork parameters
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Params {
    /// Deprecated
    /// VoteThreshold defines the percentage of bonded stake required to vote for a scheduled cork to be approved
    #[prost(string, tag = "1")]
    pub vote_threshold: ::prost::alloc::string::String,
    #[prost(uint64, tag = "2")]
    pub max_corks_per_validator: u64,
}
/// QueryParamsRequest is the request type for the Query/Params gRPC method.
#[derive(Clone, Copy, PartialEq, ::prost::Message)]
pub struct QueryParamsRequest {}
/// QueryParamsRequest is the response type for the Query/Params gRPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryParamsResponse {
    /// allocation parameters
    #[prost(message, optional, tag = "1")]
    pub params: ::core::option::Option<Params>,
}
/// QueryCellarIDsRequest is the request type for Query/QueryCellarIDs gRPC method.
#[derive(Clone, Copy, PartialEq, ::prost::Message)]
pub struct QueryCellarIDsRequest {}
/// QueryCellarIDsResponse is the response type for Query/QueryCellars gRPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryCellarIDsResponse {
    #[prost(string, repeated, tag = "1")]
    pub cellar_ids: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
}
/// QueryScheduledCorksRequest
#[derive(Clone, Copy, PartialEq, ::prost::Message)]
pub struct QueryScheduledCorksRequest {}
/// QueryScheduledCorksResponse
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryScheduledCorksResponse {
    #[prost(message, repeated, tag = "1")]
    pub corks: ::prost::alloc::vec::Vec<ScheduledCork>,
}
/// QueryScheduledBlockHeightsRequest
#[derive(Clone, Copy, PartialEq, ::prost::Message)]
pub struct QueryScheduledBlockHeightsRequest {}
/// QueryScheduledBlockHeightsResponse
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryScheduledBlockHeightsResponse {
    #[prost(uint64, repeated, tag = "1")]
    pub block_heights: ::prost::alloc::vec::Vec<u64>,
}
/// QueryScheduledCorksByBlockHeightRequest
#[derive(Clone, Copy, PartialEq, ::prost::Message)]
pub struct QueryScheduledCorksByBlockHeightRequest {
    #[prost(uint64, tag = "1")]
    pub block_height: u64,
}
/// QueryScheduledCorksByBlockHeightResponse
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryScheduledCorksByBlockHeightResponse {
    #[prost(message, repeated, tag = "1")]
    pub corks: ::prost::alloc::vec::Vec<ScheduledCork>,
}
/// QueryScheduledCorksByIDRequest
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryScheduledCorksByIdRequest {
    #[prost(string, tag = "1")]
    pub id: ::prost::alloc::string::String,
}
/// QueryScheduledCorksByIDResponse
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryScheduledCorksByIdResponse {
    #[prost(message, repeated, tag = "1")]
    pub corks: ::prost::alloc::vec::Vec<ScheduledCork>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryCorkResultRequest {
    #[prost(string, tag = "1")]
    pub id: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryCorkResultResponse {
    #[prost(message, optional, tag = "1")]
    pub cork_result: ::core::option::Option<CorkResult>,
}
#[derive(Clone, Copy, PartialEq, ::prost::Message)]
pub struct QueryCorkResultsRequest {}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryCorkResultsResponse {
    #[prost(message, repeated, tag = "1")]
    pub cork_results: ::prost::alloc::vec::Vec<CorkResult>,
}
/// Generated client implementations.
pub mod query_client {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value,
    )]
    use tonic::codegen::*;
    use tonic::codegen::http::Uri;
    /// Query defines the gRPC query service for the cork module.
    #[derive(Debug, Clone)]
    pub struct QueryClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl QueryClient<tonic::transport::Channel> {
        /// Attempt to create a new client by connecting to a given endpoint.
        pub async fn connect<D>(dst: D) -> Result<Self, tonic::transport::Error>
        where
            D: TryInto<tonic::transport::Endpoint>,
            D::Error: Into<StdError>,
        {
            let conn = tonic::transport::Endpoint::new(dst)?.connect().await?;
            Ok(Self::new(conn))
        }
    }
    impl<T> QueryClient<T>
    where
        T: tonic::client::GrpcService<tonic::body::BoxBody>,
        T::Error: Into<StdError>,
        T::ResponseBody: Body<Data = Bytes> + std::marker::Send + 'static,
        <T::ResponseBody as Body>::Error: Into<StdError> + std::marker::Send,
    {
        pub fn new(inner: T) -> Self {
            let inner = tonic::client::Grpc::new(inner);
            Self { inner }
        }
        pub fn with_origin(inner: T, origin: Uri) -> Self {
            let inner = tonic::client::Grpc::with_origin(inner, origin);
            Self { inner }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> QueryClient<InterceptedService<T, F>>
        where
            F: tonic::service::Interceptor,
            T::ResponseBody: Default,
            T: tonic::codegen::Service<
                http::Request<tonic::body::BoxBody>,
                Response = http::Response<
                    <T as tonic::client::GrpcService<tonic::body::BoxBody>>::ResponseBody,
                >,
            >,
            <T as tonic::codegen::Service<
                http::Request<tonic::body::BoxBody>,
            >>::Error: Into<StdError> + std::marker::Send + std::marker::Sync,
        {
            QueryClient::new(InterceptedService::new(inner, interceptor))
        }
        /// Compress requests with the given encoding.
        ///
        /// This requires the server to support it otherwise it might respond with an
        /// error.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.send_compressed(encoding);
            self
        }
        /// Enable decompressing responses.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.accept_compressed(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_decoding_message_size(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_encoding_message_size(limit);
            self
        }
        /// QueryParams queries the allocation module parameters.
        pub async fn query_params(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryParamsRequest>,
        ) -> std::result::Result<
            tonic::Response<super::QueryParamsResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/cork.v2.Query/QueryParams",
            );
            let mut req = request.into_request();
            req.extensions_mut().insert(GrpcMethod::new("cork.v2.Query", "QueryParams"));
            self.inner.unary(req, path, codec).await
        }
        /// QueryCellarIDs returns all cellars and current tick ranges
        pub async fn query_cellar_i_ds(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryCellarIDsRequest>,
        ) -> std::result::Result<
            tonic::Response<super::QueryCellarIDsResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/cork.v2.Query/QueryCellarIDs",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("cork.v2.Query", "QueryCellarIDs"));
            self.inner.unary(req, path, codec).await
        }
        /// QueryScheduledCorks returns all scheduled corks
        pub async fn query_scheduled_corks(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryScheduledCorksRequest>,
        ) -> std::result::Result<
            tonic::Response<super::QueryScheduledCorksResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/cork.v2.Query/QueryScheduledCorks",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("cork.v2.Query", "QueryScheduledCorks"));
            self.inner.unary(req, path, codec).await
        }
        /// QueryScheduledBlockHeights returns all scheduled block heights
        pub async fn query_scheduled_block_heights(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryScheduledBlockHeightsRequest>,
        ) -> std::result::Result<
            tonic::Response<super::QueryScheduledBlockHeightsResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/cork.v2.Query/QueryScheduledBlockHeights",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("cork.v2.Query", "QueryScheduledBlockHeights"));
            self.inner.unary(req, path, codec).await
        }
        /// QueryScheduledCorks returns all scheduled corks at a block height
        pub async fn query_scheduled_corks_by_block_height(
            &mut self,
            request: impl tonic::IntoRequest<
                super::QueryScheduledCorksByBlockHeightRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::QueryScheduledCorksByBlockHeightResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/cork.v2.Query/QueryScheduledCorksByBlockHeight",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new("cork.v2.Query", "QueryScheduledCorksByBlockHeight"),
                );
            self.inner.unary(req, path, codec).await
        }
        /// QueryScheduledCorks returns all scheduled corks with the specified ID
        pub async fn query_scheduled_corks_by_id(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryScheduledCorksByIdRequest>,
        ) -> std::result::Result<
            tonic::Response<super::QueryScheduledCorksByIdResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/cork.v2.Query/QueryScheduledCorksByID",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("cork.v2.Query", "QueryScheduledCorksByID"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn query_cork_result(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryCorkResultRequest>,
        ) -> std::result::Result<
            tonic::Response<super::QueryCorkResultResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/cork.v2.Query/QueryCorkResult",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("cork.v2.Query", "QueryCorkResult"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn query_cork_results(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryCorkResultsRequest>,
        ) -> std::result::Result<
            tonic::Response<super::QueryCorkResultsResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/cork.v2.Query/QueryCorkResults",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("cork.v2.Query", "QueryCorkResults"));
            self.inner.unary(req, path, codec).await
        }
    }
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct AddManagedCellarIDsProposal {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(message, optional, tag = "3")]
    pub cellar_ids: ::core::option::Option<CellarIdSet>,
    #[prost(string, tag = "4")]
    pub publisher_domain: ::prost::alloc::string::String,
}
/// AddManagedCellarIDsProposalWithDeposit is a specific definition for CLI commands
#[derive(serde::Deserialize, serde::Serialize)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct AddManagedCellarIDsProposalWithDeposit {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(string, repeated, tag = "3")]
    pub cellar_ids: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
    #[prost(string, tag = "4")]
    pub publisher_domain: ::prost::alloc::string::String,
    #[prost(string, tag = "5")]
    pub deposit: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct RemoveManagedCellarIDsProposal {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(message, optional, tag = "3")]
    pub cellar_ids: ::core::option::Option<CellarIdSet>,
}
/// RemoveManagedCellarIDsProposalWithDeposit is a specific definition for CLI commands
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct RemoveManagedCellarIDsProposalWithDeposit {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(string, repeated, tag = "3")]
    pub cellar_ids: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
    #[prost(string, tag = "4")]
    pub deposit: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ScheduledCorkProposal {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(uint64, tag = "3")]
    pub block_height: u64,
    #[prost(string, tag = "4")]
    pub target_contract_address: ::prost::alloc::string::String,
    ///
    /// The JSON representation of a ScheduleRequest defined in the Steward protos
    ///
    /// Example: The following is the JSON form of a ScheduleRequest containing a steward.v2.cellar_v1.TrustPosition
    /// message, which maps to the `trustPosition(address)` function of the the V1 Cellar contract.
    ///
    /// {
    ///    "cellar_id": "0x1234567890000000000000000000000000000000",
    ///    "cellar_v1": {
    ///      "trust_position": {
    ///        "erc20_address": "0x1234567890000000000000000000000000000000"
    ///      }
    ///    },
    ///    "block_height": 1000000
    /// }
    ///
    /// You can use the Steward CLI to generate the required JSON rather than constructing it by hand <https://github.com/peggyjv/steward>
    #[prost(string, tag = "5")]
    pub contract_call_proto_json: ::prost::alloc::string::String,
}
/// ScheduledCorkProposalWithDeposit is a specific definition for CLI commands
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ScheduledCorkProposalWithDeposit {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(uint64, tag = "3")]
    pub block_height: u64,
    #[prost(string, tag = "4")]
    pub target_contract_address: ::prost::alloc::string::String,
    #[prost(string, tag = "5")]
    pub contract_call_proto_json: ::prost::alloc::string::String,
    #[prost(string, tag = "6")]
    pub deposit: ::prost::alloc::string::String,
}
