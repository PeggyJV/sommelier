/// AllocationPrecommit defines an array of hashed data to be used for the precommit phase
/// of allocation
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct AllocationPrecommit {
    ///  bytes  hash = 1 [(gogoproto.casttype) = "github.com/tendermint/tendermint/libs/bytes.HexBytes"];
    #[prost(bytes = "vec", tag = "1")]
    pub hash: ::prost::alloc::vec::Vec<u8>,
    #[prost(string, tag = "2")]
    pub cellar_id: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct RebalanceVote {
    #[prost(message, optional, tag = "1")]
    pub cellar: ::core::option::Option<Cellar>,
    #[prost(uint64, tag = "2")]
    pub current_price: u64,
}
/// Allocation is the commit for all allocations for a cellar by a validator
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Allocation {
    #[prost(message, optional, tag = "1")]
    pub vote: ::core::option::Option<RebalanceVote>,
    #[prost(string, tag = "2")]
    pub salt: ::prost::alloc::string::String,
}
/// Cellar is a collection of pools for a token pair
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Cellar {
    #[prost(string, tag = "1")]
    pub id: ::prost::alloc::string::String,
    #[prost(message, repeated, tag = "2")]
    pub tick_ranges: ::prost::alloc::vec::Vec<TickRange>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct TickRange {
    #[prost(int32, tag = "1")]
    pub upper: i32,
    #[prost(int32, tag = "2")]
    pub lower: i32,
    #[prost(uint32, tag = "3")]
    pub weight: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct CellarUpdate {
    #[prost(uint64, tag = "1")]
    pub invalidation_nonce: u64,
    #[prost(message, optional, tag = "2")]
    pub vote: ::core::option::Option<RebalanceVote>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct AddManagedCellarsProposal {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(string, repeated, tag = "3")]
    pub cellar_ids: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct RemoveManagedCellarsProposal {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(string, repeated, tag = "3")]
    pub cellar_ids: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
}
/// MsgAllocationPrecommit - sdk.Msg for prevoting on an array of oracle data types.
/// The purpose of the prevote is to hide vote for data with hashes formatted as hex string:
/// SHA256("{salt}:{data_cannonical_json}:{voter}")
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgAllocationPrecommit {
    /// precommit containing the hash of the allocation precommit contents
    #[prost(message, repeated, tag = "1")]
    pub precommit: ::prost::alloc::vec::Vec<AllocationPrecommit>,
    /// signer (i.e feeder) account address
    #[prost(string, tag = "2")]
    pub signer: ::prost::alloc::string::String,
}
/// MsgAllocationPrecommitResponse is the response type for the Msg/AllocationPrecommitResponse gRPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgAllocationPrecommitResponse {}
/// MsgAllocationCommit - sdk.Msg for submitting arbitrary oracle data that has been prevoted on
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgAllocationCommit {
    /// vote containing the oracle data feed
    #[prost(message, repeated, tag = "1")]
    pub commit: ::prost::alloc::vec::Vec<Allocation>,
    /// signer (i.e feeder) account address
    #[prost(string, tag = "2")]
    pub signer: ::prost::alloc::string::String,
}
/// MsgAllocationCommitResponse is the response type for the Msg/AllocationCommit gRPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgAllocationCommitResponse {}
#[doc = r" Generated client implementations."]
pub mod msg_client {
    #![allow(unused_variables, dead_code, missing_docs)]
    use tonic::codegen::*;
    #[doc = " MsgService defines the msgs that the allocation module handles."]
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
        #[doc = " AllocationPrecommit defines a message that commits a hash of allocation data feed before the data is actually submitted."]
        pub async fn allocation_precommit(
            &mut self,
            request: impl tonic::IntoRequest<super::MsgAllocationPrecommit>,
        ) -> Result<tonic::Response<super::MsgAllocationPrecommitResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path =
                http::uri::PathAndQuery::from_static("/allocation.v1.Msg/AllocationPrecommit");
            self.inner.unary(request.into_request(), path, codec).await
        }
        #[doc = " AllocationCommit defines a message to submit the actual allocation data that was committed by the feeder through the pre-commit."]
        pub async fn allocation_commit(
            &mut self,
            request: impl tonic::IntoRequest<super::MsgAllocationCommit>,
        ) -> Result<tonic::Response<super::MsgAllocationCommitResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/allocation.v1.Msg/AllocationCommit");
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
/// GenesisState - all allocation state that must be provided at genesis
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GenesisState {
    #[prost(message, optional, tag = "1")]
    pub params: ::core::option::Option<Params>,
    #[prost(message, repeated, tag = "2")]
    pub cellars: ::prost::alloc::vec::Vec<Cellar>,
}
/// Params allocation parameters
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
/// QueryAllocationPrecommitRequest is the request type for the Query/AllocationPrecommit gRPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryAllocationPrecommitRequest {
    /// validator operator address
    #[prost(string, tag = "1")]
    pub validator: ::prost::alloc::string::String,
    /// cellar contract address
    #[prost(string, tag = "2")]
    pub cellar: ::prost::alloc::string::String,
}
/// QueryAllocationPrecommitResponse is the response type for the Query/AllocationPrecommit gRPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryAllocationPrecommitResponse {
    /// prevote submitted within the latest voting period
    #[prost(message, optional, tag = "1")]
    pub precommit: ::core::option::Option<AllocationPrecommit>,
}
/// QueryAllocationPrecommitsRequest is the request type for the Query/AllocationPrecommits gRPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryAllocationPrecommitsRequest {}
/// QueryAllocationPrecommitResponse is the response type for the Query/AllocationPrecommits gRPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryAllocationPrecommitsResponse {
    /// prevote submitted within the latest voting period
    #[prost(message, repeated, tag = "1")]
    pub precommits: ::prost::alloc::vec::Vec<AllocationPrecommit>,
}
/// QueryAllocationCommitRequest is the request type for the Query/QueryallocationDataVote gRPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryAllocationCommitRequest {
    /// validator operator address
    #[prost(string, tag = "1")]
    pub validator: ::prost::alloc::string::String,
    /// cellar contract address
    #[prost(string, tag = "2")]
    pub cellar: ::prost::alloc::string::String,
}
/// QueryAllocationCommitResponse is the response type for the Query/QueryallocationDataVote gRPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryAllocationCommitResponse {
    /// vote containing the allocation feed submitted within the latest voting period
    #[prost(message, optional, tag = "1")]
    pub commit: ::core::option::Option<Allocation>,
}
/// QueryAllocationCommitsRequest is the request type for the Query/QueryAllocationCommits gRPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryAllocationCommitsRequest {}
/// QueryAllocationCommitsResponse is the response type for the Query/QueryAllocationCommits gRPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryAllocationCommitsResponse {
    /// votes containing the allocation feed submitted within the latest voting period
    #[prost(message, repeated, tag = "1")]
    pub commits: ::prost::alloc::vec::Vec<Allocation>,
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
/// QueryCellarsRequest is the request type for Query/QueryCellars gRPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryCellarsRequest {}
/// QueryCellarsResponse is the response type for Query/QueryCellars gRPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryCellarsResponse {
    #[prost(message, repeated, tag = "1")]
    pub cellars: ::prost::alloc::vec::Vec<Cellar>,
}
#[doc = r" Generated client implementations."]
pub mod query_client {
    #![allow(unused_variables, dead_code, missing_docs)]
    use tonic::codegen::*;
    #[doc = " Query defines the gRPC querier service for the allocation module."]
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
            let path = http::uri::PathAndQuery::from_static("/allocation.v1.Query/QueryParams");
            self.inner.unary(request.into_request(), path, codec).await
        }
        #[doc = " QueryAllocationPrecommit queries the validator prevote in the current voting period"]
        pub async fn query_allocation_precommit(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryAllocationPrecommitRequest>,
        ) -> Result<tonic::Response<super::QueryAllocationPrecommitResponse>, tonic::Status>
        {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/allocation.v1.Query/QueryAllocationPrecommit",
            );
            self.inner.unary(request.into_request(), path, codec).await
        }
        #[doc = " QueryAllocationPrecommits queries all allocation precommits in the voting period"]
        pub async fn query_allocation_precommits(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryAllocationPrecommitsRequest>,
        ) -> Result<tonic::Response<super::QueryAllocationPrecommitsResponse>, tonic::Status>
        {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/allocation.v1.Query/QueryAllocationPrecommits",
            );
            self.inner.unary(request.into_request(), path, codec).await
        }
        #[doc = " QueryAllocationCommit queries the validator vote in the current voting period"]
        pub async fn query_allocation_commit(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryAllocationCommitRequest>,
        ) -> Result<tonic::Response<super::QueryAllocationCommitResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path =
                http::uri::PathAndQuery::from_static("/allocation.v1.Query/QueryAllocationCommit");
            self.inner.unary(request.into_request(), path, codec).await
        }
        #[doc = " QueryAllocationCommits queries all validator allocation commits"]
        pub async fn query_allocation_commits(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryAllocationCommitsRequest>,
        ) -> Result<tonic::Response<super::QueryAllocationCommitsResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path =
                http::uri::PathAndQuery::from_static("/allocation.v1.Query/QueryAllocationCommits");
            self.inner.unary(request.into_request(), path, codec).await
        }
        #[doc = " QueryVotePeriod queries the heights for the current voting period (current, start and end)"]
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
            let path =
                http::uri::PathAndQuery::from_static("/allocation.v1.Query/QueryCommitPeriod");
            self.inner.unary(request.into_request(), path, codec).await
        }
        #[doc = " QueryCellars returns all cellars and current tick ranges"]
        pub async fn query_cellars(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryCellarsRequest>,
        ) -> Result<tonic::Response<super::QueryCellarsResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/allocation.v1.Query/QueryCellars");
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
