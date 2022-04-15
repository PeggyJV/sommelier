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
pub struct ValidatorCork {
    #[prost(message, optional, tag = "1")]
    pub cork: ::core::option::Option<Cork>,
    #[prost(string, tag = "2")]
    pub validator: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ScheduledCork {
    #[prost(message, optional, tag = "1")]
    pub cork: ::core::option::Option<Cork>,
    #[prost(uint64, tag = "2")]
    pub block_height: u64,
    #[prost(string, tag = "3")]
    pub validator: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct CellarIdSet {
    #[prost(string, repeated, tag = "1")]
    pub ids: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
}
/// MsgSubmitCorkRequest - sdk.Msg for submitting calls to Ethereum through the gravity bridge contract
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgSubmitCorkRequest {
    /// the cork to send across the bridge
    #[prost(message, optional, tag = "1")]
    pub cork: ::core::option::Option<Cork>,
    /// signer account address
    #[prost(string, tag = "2")]
    pub signer: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgSubmitCorkResponse {}
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
pub struct MsgScheduleCorkResponse {}
#[doc = r" Generated client implementations."]
pub mod msg_client {
    #![allow(unused_variables, dead_code, missing_docs)]
    use tonic::codegen::*;
    #[doc = " MsgService defines the msgs that the cork module handles"]
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
        pub async fn submit_cork(
            &mut self,
            request: impl tonic::IntoRequest<super::MsgSubmitCorkRequest>,
        ) -> Result<tonic::Response<super::MsgSubmitCorkResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/cork.v1.Msg/SubmitCork");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn schedule_cork(
            &mut self,
            request: impl tonic::IntoRequest<super::MsgScheduleCorkRequest>,
        ) -> Result<tonic::Response<super::MsgScheduleCorkResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/cork.v1.Msg/ScheduleCork");
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
    pub corks: ::prost::alloc::vec::Vec<ValidatorCork>,
    #[prost(message, repeated, tag = "5")]
    pub scheduled_corks: ::prost::alloc::vec::Vec<ScheduledCork>,
}
/// Params cork parameters
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Params {
    /// VotePeriod defines the number of blocks to wait for votes before attempting to tally
    #[prost(int64, tag = "1")]
    pub vote_period: i64,
    /// VoteThreshold defines the percentage of bonded stake required to vote each period
    #[prost(string, tag = "2")]
    pub vote_threshold: ::prost::alloc::string::String,
}
/// QueryParamsRequest is the request type for the Query/Params gRPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryParamsRequest {}
/// QueryParamsRequest is the response type for the Query/Params gRPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryParamsResponse {
    /// allocation parameters
    #[prost(message, optional, tag = "1")]
    pub params: ::core::option::Option<Params>,
}
/// QuerySubmittedCorksRequest is the request type for the Query/QuerySubmittedCorks gRPC query method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QuerySubmittedCorksRequest {}
/// QuerySubmittedCorksResponse is the response type for the Query/QuerySubmittedCorks gRPC query method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QuerySubmittedCorksResponse {
    /// corks in keeper awaiting vote
    #[prost(message, repeated, tag = "1")]
    pub corks: ::prost::alloc::vec::Vec<Cork>,
}
/// QueryCommitPeriodRequest is the request type for the Query/QueryCommitPeriod gRPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryCommitPeriodRequest {}
/// QueryCommitPeriodResponse is the response type for the Query/QueryCommitPeriod gRPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryCommitPeriodResponse {
    /// block height at which the query was processed
    #[prost(int64, tag = "1")]
    pub current_height: i64,
    /// latest vote period start block height
    #[prost(int64, tag = "2")]
    pub vote_period_start: i64,
    /// block height at which the current voting period ends
    #[prost(int64, tag = "3")]
    pub vote_period_end: i64,
}
/// QueryCellarIDsRequest is the request type for Query/QueryCellarIDs gRPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryCellarIDsRequest {}
/// QueryCellarIDsResponse is the response type for Query/QueryCellars gRPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryCellarIDsResponse {
    #[prost(string, repeated, tag = "1")]
    pub cellar_ids: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
}
/// QueryScheduledCorksRequest
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryScheduledCorksRequest {}
/// QueryScheduledCorksResponse
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryScheduledCorksResponse {
    #[prost(message, repeated, tag = "1")]
    pub corks: ::prost::alloc::vec::Vec<ScheduledCork>,
}
/// QueryScheduledBlockHeightsRequest
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryScheduledBlockHeightsRequest {}
/// QueryScheduledBlockHeightsResponse
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryScheduledBlockHeightsResponse {
    #[prost(uint64, repeated, tag = "1")]
    pub block_heights: ::prost::alloc::vec::Vec<u64>,
}
/// QueryScheduledCorksByBlockHeightRequest
#[derive(Clone, PartialEq, ::prost::Message)]
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
#[doc = r" Generated client implementations."]
pub mod query_client {
    #![allow(unused_variables, dead_code, missing_docs)]
    use tonic::codegen::*;
    #[doc = " Query defines the gRPC query service for the cork module."]
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
        #[doc = " QueryParams queries the allocation module parameters."]
        pub async fn query_params(
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
            let path = http::uri::PathAndQuery::from_static("/cork.v1.Query/QueryParams");
            self.inner.unary(request.into_request(), path, codec).await
        }
        #[doc = " QuerySubmittedCorks queries the submitted corks awaiting vote"]
        pub async fn query_submitted_corks(
            &mut self,
            request: impl tonic::IntoRequest<super::QuerySubmittedCorksRequest>,
        ) -> Result<tonic::Response<super::QuerySubmittedCorksResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/cork.v1.Query/QuerySubmittedCorks");
            self.inner.unary(request.into_request(), path, codec).await
        }
        #[doc = " QueryCommitPeriod queries the heights for the current voting period (current, start and end)"]
        pub async fn query_commit_period(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryCommitPeriodRequest>,
        ) -> Result<tonic::Response<super::QueryCommitPeriodResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/cork.v1.Query/QueryCommitPeriod");
            self.inner.unary(request.into_request(), path, codec).await
        }
        #[doc = " QueryCellarIDs returns all cellars and current tick ranges"]
        pub async fn query_cellar_i_ds(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryCellarIDsRequest>,
        ) -> Result<tonic::Response<super::QueryCellarIDsResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/cork.v1.Query/QueryCellarIDs");
            self.inner.unary(request.into_request(), path, codec).await
        }
        #[doc = " QueryScheduledCorks returns all scheduled corks"]
        pub async fn query_scheduled_corks(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryScheduledCorksRequest>,
        ) -> Result<tonic::Response<super::QueryScheduledCorksResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/cork.v1.Query/QueryScheduledCorks");
            self.inner.unary(request.into_request(), path, codec).await
        }
        #[doc = " QueryScheduledBlockHeights returns all scheduled block heights"]
        pub async fn query_scheduled_block_heights(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryScheduledBlockHeightsRequest>,
        ) -> Result<tonic::Response<super::QueryScheduledBlockHeightsResponse>, tonic::Status>
        {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path =
                http::uri::PathAndQuery::from_static("/cork.v1.Query/QueryScheduledBlockHeights");
            self.inner.unary(request.into_request(), path, codec).await
        }
        #[doc = " QueryScheduledCorks returns all scheduled corks at a block height"]
        pub async fn query_scheduled_corks_by_block_height(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryScheduledCorksByBlockHeightRequest>,
        ) -> Result<tonic::Response<super::QueryScheduledCorksByBlockHeightResponse>, tonic::Status>
        {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/cork.v1.Query/QueryScheduledCorksByBlockHeight",
            );
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
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct AddManagedCellarIDsProposal {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(message, optional, tag = "3")]
    pub cellar_ids: ::core::option::Option<CellarIdSet>,
}
/// AddManagedCellarIDsProposalWithDeposit is a specific definition for CLI commands
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct AddManagedCellarIDsProposalWithDeposit {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(string, repeated, tag = "3")]
    pub cellar_ids: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
    #[prost(string, tag = "4")]
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
