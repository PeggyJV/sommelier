#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Auction {
    #[prost(uint32, tag = "1")]
    pub id: u32,
    #[prost(message, optional, tag = "2")]
    pub starting_tokens_for_sale: ::core::option::Option<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
    #[prost(uint64, tag = "3")]
    pub start_block: u64,
    #[prost(uint64, tag = "4")]
    pub end_block: u64,
    #[prost(string, tag = "5")]
    pub initial_price_decrease_rate: ::prost::alloc::string::String,
    #[prost(string, tag = "6")]
    pub current_price_decrease_rate: ::prost::alloc::string::String,
    #[prost(uint64, tag = "7")]
    pub price_decrease_block_interval: u64,
    #[prost(string, tag = "8")]
    pub initial_unit_price_in_usomm: ::prost::alloc::string::String,
    #[prost(string, tag = "9")]
    pub current_unit_price_in_usomm: ::prost::alloc::string::String,
    #[prost(message, optional, tag = "10")]
    pub remaining_tokens_for_sale:
        ::core::option::Option<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
    #[prost(string, tag = "11")]
    pub funding_module_account: ::prost::alloc::string::String,
    #[prost(string, tag = "12")]
    pub proceeds_module_account: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Bid {
    #[prost(uint64, tag = "1")]
    pub id: u64,
    #[prost(uint32, tag = "2")]
    pub auction_id: u32,
    #[prost(string, tag = "3")]
    pub bidder: ::prost::alloc::string::String,
    #[prost(message, optional, tag = "4")]
    pub max_bid_in_usomm: ::core::option::Option<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
    #[prost(message, optional, tag = "5")]
    pub sale_token_minimum_amount:
        ::core::option::Option<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
    #[prost(message, optional, tag = "6")]
    pub total_fulfilled_sale_tokens:
        ::core::option::Option<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
    #[prost(string, tag = "7")]
    pub sale_token_unit_price_in_usomm: ::prost::alloc::string::String,
    #[prost(message, optional, tag = "8")]
    pub total_usomm_paid: ::core::option::Option<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
    #[prost(uint64, tag = "9")]
    pub block_height: u64,
}
/// USD price is the value for one non-fractional token (smallest unit of the token * 10^exponent)
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct TokenPrice {
    #[prost(string, tag = "1")]
    pub denom: ::prost::alloc::string::String,
    #[prost(uint64, tag = "2")]
    pub exponent: u64,
    #[prost(string, tag = "3")]
    pub usd_price: ::prost::alloc::string::String,
    #[prost(uint64, tag = "4")]
    pub last_updated_block: u64,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ProposedTokenPrice {
    #[prost(string, tag = "1")]
    pub denom: ::prost::alloc::string::String,
    #[prost(uint64, tag = "2")]
    pub exponent: u64,
    #[prost(string, tag = "3")]
    pub usd_price: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgSubmitBidRequest {
    #[prost(uint32, tag = "1")]
    pub auction_id: u32,
    #[prost(string, tag = "2")]
    pub signer: ::prost::alloc::string::String,
    #[prost(message, optional, tag = "3")]
    pub max_bid_in_usomm: ::core::option::Option<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
    #[prost(message, optional, tag = "4")]
    pub sale_token_minimum_amount:
        ::core::option::Option<cosmos_sdk_proto::cosmos::base::v1beta1::Coin>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgSubmitBidResponse {
    #[prost(message, optional, tag = "1")]
    pub bid: ::core::option::Option<Bid>,
}
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
        pub async fn submit_bid(
            &mut self,
            request: impl tonic::IntoRequest<super::MsgSubmitBidRequest>,
        ) -> Result<tonic::Response<super::MsgSubmitBidResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/auction.v1.Msg/SubmitBid");
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
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GenesisState {
    #[prost(message, optional, tag = "1")]
    pub params: ::core::option::Option<Params>,
    #[prost(message, repeated, tag = "2")]
    pub auctions: ::prost::alloc::vec::Vec<Auction>,
    #[prost(message, repeated, tag = "3")]
    pub bids: ::prost::alloc::vec::Vec<Bid>,
    #[prost(message, repeated, tag = "4")]
    pub token_prices: ::prost::alloc::vec::Vec<TokenPrice>,
    #[prost(uint32, tag = "5")]
    pub last_auction_id: u32,
    #[prost(uint64, tag = "6")]
    pub last_bid_id: u64,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Params {
    #[prost(uint64, tag = "1")]
    pub price_max_block_age: u64,
    #[prost(uint64, tag = "2")]
    pub minimum_bid_in_usomm: u64,
    #[prost(string, tag = "3")]
    pub minimum_sale_tokens_usd_value: ::prost::alloc::string::String,
    #[prost(uint64, tag = "4")]
    pub auction_max_block_age: u64,
    #[prost(string, tag = "5")]
    pub auction_price_decrease_acceleration_rate: ::prost::alloc::string::String,
    #[prost(uint64, tag = "6")]
    pub minimum_auction_height: u64,
    /// value between 0 and 1 the determines the % of somm received from bids that
    /// gets burned
    #[prost(string, tag = "7")]
    pub somm_burn_ratio: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryParamsRequest {}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryParamsResponse {
    #[prost(message, optional, tag = "1")]
    pub params: ::core::option::Option<Params>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryActiveAuctionRequest {
    #[prost(uint32, tag = "1")]
    pub auction_id: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryActiveAuctionResponse {
    #[prost(message, optional, tag = "1")]
    pub auction: ::core::option::Option<Auction>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryEndedAuctionRequest {
    #[prost(uint32, tag = "1")]
    pub auction_id: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryEndedAuctionResponse {
    #[prost(message, optional, tag = "1")]
    pub auction: ::core::option::Option<Auction>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryActiveAuctionsRequest {}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryActiveAuctionsResponse {
    #[prost(message, repeated, tag = "1")]
    pub auctions: ::prost::alloc::vec::Vec<Auction>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryEndedAuctionsRequest {
    #[prost(message, optional, tag = "1")]
    pub pagination: ::core::option::Option<cosmos_sdk_proto::cosmos::base::query::v1beta1::PageRequest>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryEndedAuctionsResponse {
    #[prost(message, repeated, tag = "1")]
    pub auctions: ::prost::alloc::vec::Vec<Auction>,
    #[prost(message, optional, tag = "2")]
    pub pagination:
        ::core::option::Option<cosmos_sdk_proto::cosmos::base::query::v1beta1::PageResponse>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryBidRequest {
    #[prost(uint64, tag = "1")]
    pub bid_id: u64,
    #[prost(uint32, tag = "2")]
    pub auction_id: u32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryBidResponse {
    #[prost(message, optional, tag = "1")]
    pub bid: ::core::option::Option<Bid>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryBidsByAuctionRequest {
    #[prost(uint32, tag = "1")]
    pub auction_id: u32,
    #[prost(message, optional, tag = "2")]
    pub pagination: ::core::option::Option<cosmos_sdk_proto::cosmos::base::query::v1beta1::PageRequest>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryBidsByAuctionResponse {
    #[prost(message, repeated, tag = "1")]
    pub bids: ::prost::alloc::vec::Vec<Bid>,
    #[prost(message, optional, tag = "2")]
    pub pagination:
        ::core::option::Option<cosmos_sdk_proto::cosmos::base::query::v1beta1::PageResponse>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryTokenPriceRequest {
    #[prost(string, tag = "1")]
    pub denom: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryTokenPriceResponse {
    #[prost(message, optional, tag = "1")]
    pub token_price: ::core::option::Option<TokenPrice>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryTokenPricesRequest {}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryTokenPricesResponse {
    #[prost(message, repeated, tag = "1")]
    pub token_prices: ::prost::alloc::vec::Vec<TokenPrice>,
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
            let path = http::uri::PathAndQuery::from_static("/auction.v1.Query/QueryParams");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn query_active_auction(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryActiveAuctionRequest>,
        ) -> Result<tonic::Response<super::QueryActiveAuctionResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/auction.v1.Query/QueryActiveAuction");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn query_ended_auction(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryEndedAuctionRequest>,
        ) -> Result<tonic::Response<super::QueryEndedAuctionResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/auction.v1.Query/QueryEndedAuction");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn query_active_auctions(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryActiveAuctionsRequest>,
        ) -> Result<tonic::Response<super::QueryActiveAuctionsResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path =
                http::uri::PathAndQuery::from_static("/auction.v1.Query/QueryActiveAuctions");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn query_ended_auctions(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryEndedAuctionsRequest>,
        ) -> Result<tonic::Response<super::QueryEndedAuctionsResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/auction.v1.Query/QueryEndedAuctions");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn query_bid(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryBidRequest>,
        ) -> Result<tonic::Response<super::QueryBidResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/auction.v1.Query/QueryBid");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn query_bids_by_auction(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryBidsByAuctionRequest>,
        ) -> Result<tonic::Response<super::QueryBidsByAuctionResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/auction.v1.Query/QueryBidsByAuction");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn query_token_price(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryTokenPriceRequest>,
        ) -> Result<tonic::Response<super::QueryTokenPriceResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/auction.v1.Query/QueryTokenPrice");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn query_token_prices(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryTokenPricesRequest>,
        ) -> Result<tonic::Response<super::QueryTokenPricesResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/auction.v1.Query/QueryTokenPrices");
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
pub struct SetTokenPricesProposal {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(message, repeated, tag = "3")]
    pub token_prices: ::prost::alloc::vec::Vec<ProposedTokenPrice>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SetTokenPricesProposalWithDeposit {
    #[prost(string, tag = "1")]
    pub title: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub description: ::prost::alloc::string::String,
    #[prost(message, repeated, tag = "3")]
    pub token_prices: ::prost::alloc::vec::Vec<ProposedTokenPrice>,
    #[prost(string, tag = "4")]
    pub deposit: ::prost::alloc::string::String,
}
