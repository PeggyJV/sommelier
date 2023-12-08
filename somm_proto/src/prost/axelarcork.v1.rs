#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct AxelarCork {
    /// call body containing the ABI encoded bytes to send to the contract
    #[prost(bytes = "vec", tag = "1")]
    pub encoded_contract_call: ::prost::alloc::vec::Vec<u8>,
    /// the chain ID of the evm target chain
    #[prost(uint64, tag = "2")]
    pub chain_id: u64,
    /// address of the contract to send the call
    #[prost(string, tag = "3")]
    pub target_contract_address: ::prost::alloc::string::String,
    /// unix timestamp before which the contract call must be executed.
    /// enforced by the proxy contract.
    #[prost(uint64, tag = "4")]
    pub deadline: u64,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct ScheduledAxelarCork {
    #[prost(message, optional, tag = "1")]
    pub cork: ::core::option::Option<AxelarCork>,
    #[prost(uint64, tag = "2")]
    pub block_height: u64,
    #[prost(string, tag = "3")]
    pub validator: ::prost::alloc::string::String,
    #[prost(string, tag = "4")]
    pub id: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct ScheduledAxelarCorks {
    #[prost(message, repeated, tag = "1")]
    pub scheduled_corks: ::prost::alloc::vec::Vec<ScheduledAxelarCork>,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct AxelarCorkResult {
    #[prost(message, optional, tag = "1")]
    pub cork: ::core::option::Option<AxelarCork>,
    #[prost(uint64, tag = "2")]
    pub block_height: u64,
    #[prost(bool, tag = "3")]
    pub approved: bool,
    #[prost(string, tag = "4")]
    pub approval_percentage: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct AxelarCorkResults {
    #[prost(message, repeated, tag = "1")]
    pub cork_results: ::prost::alloc::vec::Vec<AxelarCorkResult>,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct CellarIdSet {
    #[prost(message, optional, tag = "1")]
    pub chain: ::core::option::Option<ChainConfiguration>,
    #[prost(string, repeated, tag = "2")]
    pub ids: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct ChainConfiguration {
    #[prost(string, tag = "1")]
    pub name: ::prost::alloc::string::String,
    #[prost(uint64, tag = "2")]
    pub id: u64,
    #[prost(string, tag = "3")]
    pub proxy_address: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct ChainConfigurations {
    #[prost(message, repeated, tag = "1")]
    pub configurations: ::prost::alloc::vec::Vec<ChainConfiguration>,
}
/// Used to enforce strictly newer call ordering per contract
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct AxelarContractCallNonce {
    #[prost(uint64, tag = "1")]
    pub chain_id: u64,
    #[prost(string, tag = "2")]
    pub contract_address: ::prost::alloc::string::String,
    #[prost(uint64, tag = "3")]
    pub nonce: u64,
}
/// Represents a proxy contract upgrade approved by governance with a delay in
/// execution in case of an error.
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct AxelarUpgradeData {
    #[prost(uint64, tag = "1")]
    pub chain_id: u64,
    #[prost(bytes = "vec", tag = "2")]
    pub payload: ::prost::alloc::vec::Vec<u8>,
    #[prost(int64, tag = "3")]
    pub executable_height_threshold: i64,
}
/// MsgScheduleCorkRequest - sdk.Msg for scheduling a cork request for on or after a specific block height
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct MsgScheduleAxelarCorkRequest {
    /// the scheduled cork
    #[prost(message, optional, tag = "1")]
    pub cork: ::core::option::Option<AxelarCork>,
    /// the chain id
    #[prost(uint64, tag = "2")]
    pub chain_id: u64,
    /// the block height that must be reached
    #[prost(uint64, tag = "3")]
    pub block_height: u64,
    /// signer account address
    #[prost(string, tag = "4")]
    pub signer: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct MsgScheduleAxelarCorkResponse {
    /// cork ID
    #[prost(string, tag = "1")]
    pub id: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct MsgRelayAxelarCorkRequest {
    #[prost(string, tag = "1")]
    pub signer: ::prost::alloc::string::String,
    #[prost(message, optional, tag = "2")]
    pub token: ::core::option::Option<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
    #[prost(uint64, tag = "3")]
    pub fee: u64,
    #[prost(uint64, tag = "4")]
    pub chain_id: u64,
    #[prost(string, tag = "5")]
    pub target_contract_address: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct MsgRelayAxelarCorkResponse {}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct MsgRelayAxelarProxyUpgradeRequest {
    #[prost(string, tag = "1")]
    pub signer: ::prost::alloc::string::String,
    #[prost(message, optional, tag = "2")]
    pub token: ::core::option::Option<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
    #[prost(uint64, tag = "3")]
    pub fee: u64,
    #[prost(uint64, tag = "4")]
    pub chain_id: u64,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct MsgRelayAxelarProxyUpgradeResponse {}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct MsgBumpAxelarCorkGasRequest {
    #[prost(string, tag = "1")]
    pub signer: ::prost::alloc::string::String,
    #[prost(message, optional, tag = "2")]
    pub token: ::core::option::Option<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
    #[prost(string, tag = "3")]
    pub message_id: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct MsgBumpAxelarCorkGasResponse {}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct MsgCancelAxelarCorkRequest {
    #[prost(string, tag = "1")]
    pub signer: ::prost::alloc::string::String,
    #[prost(uint64, tag = "2")]
    pub chain_id: u64,
    #[prost(string, tag = "3")]
    pub target_contract_address: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct MsgCancelAxelarCorkResponse {}
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
        pub async fn schedule_cork(
            &mut self,
            request: impl tonic::IntoRequest<super::MsgScheduleAxelarCorkRequest>,
        ) -> Result<tonic::Response<super::MsgScheduleAxelarCorkResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/axelarcork.v1.Msg/ScheduleCork");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn relay_cork(
            &mut self,
            request: impl tonic::IntoRequest<super::MsgRelayAxelarCorkRequest>,
        ) -> Result<tonic::Response<super::MsgRelayAxelarCorkResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/axelarcork.v1.Msg/RelayCork");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn bump_cork_gas(
            &mut self,
            request: impl tonic::IntoRequest<super::MsgBumpAxelarCorkGasRequest>,
        ) -> Result<tonic::Response<super::MsgBumpAxelarCorkGasResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/axelarcork.v1.Msg/BumpCorkGas");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn cancel_scheduled_cork(
            &mut self,
            request: impl tonic::IntoRequest<super::MsgCancelAxelarCorkRequest>,
        ) -> Result<tonic::Response<super::MsgCancelAxelarCorkResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path =
                http::uri::PathAndQuery::from_static("/axelarcork.v1.Msg/CancelScheduledCork");
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
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct GenesisState {
    #[prost(message, optional, tag = "1")]
    pub params: ::core::option::Option<Params>,
    #[prost(message, optional, tag = "2")]
    pub chain_configurations: ::core::option::Option<ChainConfigurations>,
    #[prost(message, repeated, tag = "3")]
    pub cellar_ids: ::prost::alloc::vec::Vec<CellarIdSet>,
    #[prost(message, optional, tag = "4")]
    pub scheduled_corks: ::core::option::Option<ScheduledAxelarCorks>,
    #[prost(message, optional, tag = "5")]
    pub cork_results: ::core::option::Option<AxelarCorkResults>,
    #[prost(message, repeated, tag = "6")]
    pub axelar_contract_call_nonces: ::prost::alloc::vec::Vec<AxelarContractCallNonce>,
    #[prost(message, repeated, tag = "7")]
    pub axelar_upgrade_data: ::prost::alloc::vec::Vec<AxelarUpgradeData>,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct Params {
    #[prost(bool, tag = "1")]
    pub enabled: bool,
    #[prost(string, tag = "2")]
    pub ibc_channel: ::prost::alloc::string::String,
    #[prost(string, tag = "3")]
    pub ibc_port: ::prost::alloc::string::String,
    #[prost(string, tag = "4")]
    pub gmp_account: ::prost::alloc::string::String,
    #[prost(string, tag = "5")]
    pub executor_account: ::prost::alloc::string::String,
    #[prost(uint64, tag = "6")]
    pub timeout_duration: u64,
    #[prost(uint64, tag = "7")]
    pub cork_timeout_blocks: u64,
}
/// QueryParamsRequest is the request type for the Query/Params gRPC method.
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryParamsRequest {}
/// QueryParamsRequest is the response type for the Query/Params gRPC method.
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryParamsResponse {
    /// allocation parameters
    #[prost(message, optional, tag = "1")]
    pub params: ::core::option::Option<Params>,
}
/// QueryCellarIDs is the request type for Query/QueryCellarIDs gRPC method.
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryCellarIDsRequest {}
/// QueryCellarIDsResponse is the response type for Query/QueryCellarIDs gRPC method.
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryCellarIDsResponse {
    #[prost(message, repeated, tag = "1")]
    pub cellar_ids: ::prost::alloc::vec::Vec<CellarIdSet>,
}
/// QueryCellarIDsByChainIDRequest is the request type for Query/QueryCellarIDsByChainID gRPC method.
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryCellarIDsByChainIdRequest {
    #[prost(uint64, tag = "1")]
    pub chain_id: u64,
}
/// QueryCellarIDsByChainIDResponse is the response type for Query/QueryCellarIDsByChainID gRPC method.
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryCellarIDsByChainIdResponse {
    #[prost(string, repeated, tag = "1")]
    pub cellar_ids: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
}
/// QueryScheduledCorksRequest
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryScheduledCorksRequest {
    #[prost(uint64, tag = "1")]
    pub chain_id: u64,
}
/// QueryScheduledCorksResponse
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryScheduledCorksResponse {
    #[prost(message, repeated, tag = "1")]
    pub corks: ::prost::alloc::vec::Vec<ScheduledAxelarCork>,
}
/// QueryScheduledBlockHeightsRequest
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryScheduledBlockHeightsRequest {
    #[prost(uint64, tag = "1")]
    pub chain_id: u64,
}
/// QueryScheduledBlockHeightsResponse
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryScheduledBlockHeightsResponse {
    #[prost(uint64, repeated, tag = "1")]
    pub block_heights: ::prost::alloc::vec::Vec<u64>,
}
/// QueryScheduledCorksByBlockHeightRequest
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryScheduledCorksByBlockHeightRequest {
    #[prost(uint64, tag = "1")]
    pub block_height: u64,
    #[prost(uint64, tag = "2")]
    pub chain_id: u64,
}
/// QueryScheduledCorksByBlockHeightResponse
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryScheduledCorksByBlockHeightResponse {
    #[prost(message, repeated, tag = "1")]
    pub corks: ::prost::alloc::vec::Vec<ScheduledAxelarCork>,
}
/// QueryScheduledCorksByIDRequest
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryScheduledCorksByIdRequest {
    #[prost(string, tag = "1")]
    pub id: ::prost::alloc::string::String,
    #[prost(uint64, tag = "2")]
    pub chain_id: u64,
}
/// QueryScheduledCorksByIDResponse
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryScheduledCorksByIdResponse {
    #[prost(message, repeated, tag = "1")]
    pub corks: ::prost::alloc::vec::Vec<ScheduledAxelarCork>,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryCorkResultRequest {
    #[prost(string, tag = "1")]
    pub id: ::prost::alloc::string::String,
    #[prost(uint64, tag = "2")]
    pub chain_id: u64,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryCorkResultResponse {
    #[prost(message, optional, tag = "1")]
    pub cork_result: ::core::option::Option<AxelarCorkResult>,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryCorkResultsRequest {
    #[prost(uint64, tag = "1")]
    pub chain_id: u64,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryCorkResultsResponse {
    #[prost(message, repeated, tag = "1")]
    pub cork_results: ::prost::alloc::vec::Vec<AxelarCorkResult>,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryChainConfigurationsRequest {}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryChainConfigurationsResponse {
    #[prost(message, repeated, tag = "1")]
    pub configurations: ::prost::alloc::vec::Vec<ChainConfiguration>,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryAxelarContractCallNoncesRequest {}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryAxelarContractCallNoncesResponse {
    #[prost(message, repeated, tag = "1")]
    pub contract_call_nonces: ::prost::alloc::vec::Vec<AxelarContractCallNonce>,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryAxelarProxyUpgradeDataRequest {}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryAxelarProxyUpgradeDataResponse {
    #[prost(message, repeated, tag = "1")]
    pub proxy_upgrade_data: ::prost::alloc::vec::Vec<AxelarUpgradeData>,
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
        #[doc = " QueryParams queries the axelar cork module parameters."]
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
            let path = http::uri::PathAndQuery::from_static("/axelarcork.v1.Query/QueryParams");
            self.inner.unary(request.into_request(), path, codec).await
        }
        #[doc = " QueryCellarIDs queries approved cellar ids of all supported chains"]
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
            let path = http::uri::PathAndQuery::from_static("/axelarcork.v1.Query/QueryCellarIDs");
            self.inner.unary(request.into_request(), path, codec).await
        }
        #[doc = " QueryCellarIDsByChainID returns all cellars and current tick ranges"]
        pub async fn query_cellar_i_ds_by_chain_id(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryCellarIDsByChainIdRequest>,
        ) -> Result<tonic::Response<super::QueryCellarIDsByChainIdResponse>, tonic::Status>
        {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/axelarcork.v1.Query/QueryCellarIDsByChainID",
            );
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
            let path =
                http::uri::PathAndQuery::from_static("/axelarcork.v1.Query/QueryScheduledCorks");
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
            let path = http::uri::PathAndQuery::from_static(
                "/axelarcork.v1.Query/QueryScheduledBlockHeights",
            );
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
                "/axelarcork.v1.Query/QueryScheduledCorksByBlockHeight",
            );
            self.inner.unary(request.into_request(), path, codec).await
        }
        #[doc = " QueryScheduledCorks returns all scheduled corks with the specified ID"]
        pub async fn query_scheduled_corks_by_id(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryScheduledCorksByIdRequest>,
        ) -> Result<tonic::Response<super::QueryScheduledCorksByIdResponse>, tonic::Status>
        {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/axelarcork.v1.Query/QueryScheduledCorksByID",
            );
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn query_cork_result(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryCorkResultRequest>,
        ) -> Result<tonic::Response<super::QueryCorkResultResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/axelarcork.v1.Query/QueryCorkResult");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn query_cork_results(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryCorkResultsRequest>,
        ) -> Result<tonic::Response<super::QueryCorkResultsResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path =
                http::uri::PathAndQuery::from_static("/axelarcork.v1.Query/QueryCorkResults");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn query_chain_configurations(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryChainConfigurationsRequest>,
        ) -> Result<tonic::Response<super::QueryChainConfigurationsResponse>, tonic::Status>
        {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/axelarcork.v1.Query/QueryChainConfigurations",
            );
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn query_axelar_contract_call_nonces(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryAxelarContractCallNoncesRequest>,
        ) -> Result<tonic::Response<super::QueryAxelarContractCallNoncesResponse>, tonic::Status>
        {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/axelarcork.v1.Query/QueryAxelarContractCallNonces",
            );
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn query_axelar_proxy_upgrade_data(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryAxelarProxyUpgradeDataRequest>,
        ) -> Result<tonic::Response<super::QueryAxelarProxyUpgradeDataResponse>, tonic::Status>
        {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/axelarcork.v1.Query/QueryAxelarProxyUpgradeData",
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
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct ScheduleCorkEvent {
    #[prost(string, tag = "1")]
    pub signer: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub validator: ::prost::alloc::string::String,
    #[prost(string, tag = "3")]
    pub cork: ::prost::alloc::string::String,
    #[prost(uint64, tag = "4")]
    pub block_height: u64,
    #[prost(uint64, tag = "5")]
    pub chain_id: u64,
}
// note: current plan is to accept either chain name or chain ID. if both and they dont match, error.

#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct AddAxelarManagedCellarIDsProposal {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(uint64, tag = "3")]
    pub chain_id: u64,
    #[prost(message, optional, tag = "4")]
    pub cellar_ids: ::core::option::Option<CellarIdSet>,
}
/// AddManagedCellarIDsProposalWithDeposit is a specific definition for CLI commands
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct AddAxelarManagedCellarIDsProposalWithDeposit {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(uint64, tag = "3")]
    pub chain_id: u64,
    #[prost(message, optional, tag = "4")]
    pub cellar_ids: ::core::option::Option<CellarIdSet>,
    #[prost(string, tag = "5")]
    pub deposit: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct RemoveAxelarManagedCellarIDsProposal {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(uint64, tag = "3")]
    pub chain_id: u64,
    #[prost(message, optional, tag = "4")]
    pub cellar_ids: ::core::option::Option<CellarIdSet>,
}
/// RemoveManagedCellarIDsProposalWithDeposit is a specific definition for CLI commands
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct RemoveAxelarManagedCellarIDsProposalWithDeposit {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(uint64, tag = "3")]
    pub chain_id: u64,
    #[prost(message, optional, tag = "4")]
    pub cellar_ids: ::core::option::Option<CellarIdSet>,
    #[prost(string, tag = "5")]
    pub deposit: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct AxelarScheduledCorkProposal {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(uint64, tag = "3")]
    pub block_height: u64,
    #[prost(uint64, tag = "4")]
    pub chain_id: u64,
    #[prost(string, tag = "5")]
    pub target_contract_address: ::prost::alloc::string::String,
    ///
    /// The JSON representation of a ScheduleRequest defined in the Steward protos
    ///
    /// Example: The following is the JSON form of a ScheduleRequest containing a steward.v2.cellar_v1.TrustPosition
    /// message, which maps to the `trustPosition(address)` function of the the V1 Cellar contract.
    ///
    /// {
    ///   "cellar_id": "0x1234567890000000000000000000000000000000",
    ///   "cellar_v1": {
    ///     "trust_position": {
    ///       "erc20_address": "0x1234567890000000000000000000000000000000"
    ///     }
    ///   },
    ///   "block_height": 1000000
    /// }
    ///
    /// You can use the Steward CLI to generate the required JSON rather than constructing it by hand
    /// https://github.com/peggyjv/steward
    #[prost(string, tag = "6")]
    pub contract_call_proto_json: ::prost::alloc::string::String,
    /// unix timestamp before which the contract call must be executed.
    /// enforced by the Axelar proxy contract
    #[prost(uint64, tag = "7")]
    pub deadline: u64,
}
/// ScheduledCorkProposalWithDeposit is a specific definition for CLI commands
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct AxelarScheduledCorkProposalWithDeposit {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(uint64, tag = "3")]
    pub block_height: u64,
    #[prost(uint64, tag = "4")]
    pub chain_id: u64,
    #[prost(string, tag = "5")]
    pub target_contract_address: ::prost::alloc::string::String,
    #[prost(string, tag = "6")]
    pub contract_call_proto_json: ::prost::alloc::string::String,
    #[prost(uint64, tag = "7")]
    pub deadline: u64,
    #[prost(string, tag = "8")]
    pub deposit: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct AxelarCommunityPoolSpendProposal {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(string, tag = "3")]
    pub recipient: ::prost::alloc::string::String,
    #[prost(uint64, tag = "4")]
    pub chain_id: u64,
    #[prost(message, optional, tag = "5")]
    pub amount: ::core::option::Option<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
}
/// This format of the community spend Ethereum proposal is specifically for
/// the CLI to allow simple text serialization.
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct AxelarCommunityPoolSpendProposalForCli {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(string, tag = "3")]
    pub recipient: ::prost::alloc::string::String,
    #[prost(uint64, tag = "4")]
    pub chain_id: u64,
    #[prost(string, tag = "5")]
    pub chain_name: ::prost::alloc::string::String,
    #[prost(string, tag = "6")]
    pub amount: ::prost::alloc::string::String,
    #[prost(string, tag = "7")]
    pub deposit: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct AddChainConfigurationProposal {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(message, optional, tag = "3")]
    pub chain_configuration: ::core::option::Option<ChainConfiguration>,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct AddChainConfigurationProposalWithDeposit {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(message, optional, tag = "3")]
    pub chain_configuration: ::core::option::Option<ChainConfiguration>,
    #[prost(string, tag = "4")]
    pub deposit: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct RemoveChainConfigurationProposal {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(uint64, tag = "3")]
    pub chain_id: u64,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct RemoveChainConfigurationProposalWithDeposit {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(uint64, tag = "3")]
    pub chain_id: u64,
    #[prost(string, tag = "4")]
    pub deposit: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct UpgradeAxelarProxyContractProposal {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(uint64, tag = "3")]
    pub chain_id: u64,
    #[prost(string, tag = "4")]
    pub new_proxy_address: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct UpgradeAxelarProxyContractProposalWithDeposit {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(uint64, tag = "3")]
    pub chain_id: u64,
    #[prost(string, tag = "4")]
    pub new_proxy_address: ::prost::alloc::string::String,
    #[prost(string, tag = "5")]
    pub deposit: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct CancelAxelarProxyContractUpgradeProposal {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(uint64, tag = "3")]
    pub chain_id: u64,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct CancelAxelarProxyContractUpgradeProposalWithDeposit {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(uint64, tag = "3")]
    pub chain_id: u64,
    #[prost(string, tag = "4")]
    pub deposit: ::prost::alloc::string::String,
}
