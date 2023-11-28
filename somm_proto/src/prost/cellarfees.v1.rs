#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct FeeAccrualCounter {
    #[prost(string, tag = "1")]
    pub denom: ::prost::alloc::string::String,
    #[prost(uint64, tag = "2")]
    pub count: u64,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct FeeAccrualCounters {
    #[prost(message, repeated, tag = "1")]
    pub counters: ::prost::alloc::vec::Vec<FeeAccrualCounter>,
}
/// Params defines the parameters for the module.
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct Params {
    /// The number of fee accruals after which an auction should be started
    #[prost(uint64, tag = "1")]
    pub fee_accrual_auction_threshold: u64,
    /// Emission rate factor. Specifically, the number of blocks over which to distribute
    /// some amount of staking rewards.
    #[prost(uint64, tag = "2")]
    pub reward_emission_period: u64,
    /// The initial rate at which auctions should decrease their denom's price in SOMM
    #[prost(string, tag = "3")]
    pub initial_price_decrease_rate: ::prost::alloc::string::String,
    /// Number of blocks between auction price decreases
    #[prost(uint64, tag = "4")]
    pub price_decrease_block_interval: u64,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryParamsRequest {}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryParamsResponse {
    #[prost(message, optional, tag = "1")]
    pub params: ::core::option::Option<Params>,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryModuleAccountsRequest {}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryModuleAccountsResponse {
    #[prost(string, tag = "1")]
    pub fees_address: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryLastRewardSupplyPeakRequest {}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryLastRewardSupplyPeakResponse {
    #[prost(string, tag = "1")]
    pub last_reward_supply_peak: ::prost::alloc::string::String,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryFeeAccrualCountersRequest {}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryFeeAccrualCountersResponse {
    #[prost(message, optional, tag = "1")]
    pub fee_accrual_counters: ::core::option::Option<FeeAccrualCounters>,
}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryApyRequest {}
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct QueryApyResponse {
    #[prost(string, tag = "1")]
    pub apy: ::prost::alloc::string::String,
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
            let path = http::uri::PathAndQuery::from_static("/cellarfees.v1.Query/QueryParams");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn query_module_accounts(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryModuleAccountsRequest>,
        ) -> Result<tonic::Response<super::QueryModuleAccountsResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path =
                http::uri::PathAndQuery::from_static("/cellarfees.v1.Query/QueryModuleAccounts");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn query_last_reward_supply_peak(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryLastRewardSupplyPeakRequest>,
        ) -> Result<tonic::Response<super::QueryLastRewardSupplyPeakResponse>, tonic::Status>
        {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/cellarfees.v1.Query/QueryLastRewardSupplyPeak",
            );
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn query_fee_accrual_counters(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryFeeAccrualCountersRequest>,
        ) -> Result<tonic::Response<super::QueryFeeAccrualCountersResponse>, tonic::Status>
        {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/cellarfees.v1.Query/QueryFeeAccrualCounters",
            );
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn query_apy(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryApyRequest>,
        ) -> Result<tonic::Response<super::QueryApyResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/cellarfees.v1.Query/QueryAPY");
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
/// GenesisState defines the cellarfees module's genesis state.
#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq, ::prost::Message)]
pub struct GenesisState {
    #[prost(message, optional, tag = "1")]
    pub params: ::core::option::Option<Params>,
    #[prost(message, optional, tag = "2")]
    pub fee_accrual_counters: ::core::option::Option<FeeAccrualCounters>,
    #[prost(string, tag = "3")]
    pub last_reward_supply_peak: ::prost::alloc::string::String,
}
