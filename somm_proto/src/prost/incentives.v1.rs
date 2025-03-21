// This file is @generated by prost-build.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GenesisState {
    #[prost(message, optional, tag = "1")]
    pub params: ::core::option::Option<Params>,
}
/// Params incentives parameters
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Params {
    /// DistributionPerBlock defines the coin to be sent to the distribution module from the community pool every block
    #[prost(message, optional, tag = "1")]
    pub distribution_per_block: ::core::option::Option<
        cosmos_sdk_proto::cosmos::base::v1beta1::Coin,
    >,
    /// IncentivesCutoffHeight defines the block height after which the incentives module will stop sending coins to the distribution module from
    /// the community pool
    #[prost(uint64, tag = "2")]
    pub incentives_cutoff_height: u64,
    /// ValidatorMaxDistributionPerBlock defines the maximum coins to be sent directly to voters in the last block from the community pool every block. Leftover coins remain in the community pool.
    #[prost(message, optional, tag = "3")]
    pub validator_max_distribution_per_block: ::core::option::Option<
        cosmos_sdk_proto::cosmos::base::v1beta1::Coin,
    >,
    /// ValidatorIncentivesCutoffHeight defines the block height after which the validator incentives will be stopped
    #[prost(uint64, tag = "4")]
    pub validator_incentives_cutoff_height: u64,
    /// ValidatorIncentivesMaxFraction defines the maximum fraction of the validator distribution per block that can be sent to a single validator
    #[prost(string, tag = "5")]
    pub validator_incentives_max_fraction: ::prost::alloc::string::String,
    /// ValidatorIncentivesSetSizeLimit defines the max number of validators to apportion the validator distribution per block to
    #[prost(uint64, tag = "6")]
    pub validator_incentives_set_size_limit: u64,
}
/// QueryParamsRequest is the request type for the QueryParams gRPC method.
#[derive(Clone, Copy, PartialEq, ::prost::Message)]
pub struct QueryParamsRequest {}
/// QueryParamsRequest is the response type for the QueryParams gRPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryParamsResponse {
    /// allocation parameters
    #[prost(message, optional, tag = "1")]
    pub params: ::core::option::Option<Params>,
}
/// QueryAPYRequest is the request type for the QueryAPY gRPC method.
#[derive(Clone, Copy, PartialEq, ::prost::Message)]
pub struct QueryApyRequest {}
/// QueryAPYRequest is the response type for the QueryAPY gRPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryApyResponse {
    #[prost(string, tag = "1")]
    pub apy: ::prost::alloc::string::String,
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
                "/incentives.v1.Query/QueryParams",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("incentives.v1.Query", "QueryParams"));
            self.inner.unary(req, path, codec).await
        }
        /// QueryAPY queries the APY returned from the incentives module.
        pub async fn query_apy(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryApyRequest>,
        ) -> std::result::Result<
            tonic::Response<super::QueryApyResponse>,
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
                "/incentives.v1.Query/QueryAPY",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("incentives.v1.Query", "QueryAPY"));
            self.inner.unary(req, path, codec).await
        }
    }
}
