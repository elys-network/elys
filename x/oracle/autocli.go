package oracle

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	_ "cosmossdk.io/api/cosmos/crypto/secp256k1" // register to that it shows up in protoregistry.GlobalTypes
	_ "cosmossdk.io/api/cosmos/crypto/secp256r1" // register to that it shows up in protoregistry.GlobalTypes
	"github.com/elys-network/elys/api/elys/oracle"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: oracle.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "shows the parameters of the module",
				},
				{
					RpcMethod:      "BandPriceResult",
					Use:            "band-price-result [request-id]",
					Short:          "Query the BandPrice result data by id",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "request_id"}},
				},
				{
					RpcMethod: "LastBandRequestId",
					Use:       "last-band-request-id",
					Short:     "Query the last request id returned by BandPrice ack packet",
				},
				{
					RpcMethod:      "AssetInfo",
					Use:            "show-asset-info [denom]",
					Short:          "shows a assetInfo",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "denom"}},
				},
				{
					RpcMethod: "AssetInfoAll",
					Use:       "list-asset-info",
					Short:     "list all assetInfo",
				},
				{
					RpcMethod:      "Price",
					Use:            "show-price [asset] [source] [timestamp]",
					Short:          "shows a price",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "asset"}, {ProtoField: "source"}, {ProtoField: "timestamp"}},
				},
				{
					RpcMethod: "PriceAll",
					Use:       "list-price",
					Short:     "list all price",
				},
				{
					RpcMethod:      "PriceFeeder",
					Use:            "show-price-feeder [feeder]",
					Short:          "shows a priceFeeder",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "feeder"}},
				},
				{
					RpcMethod: "PriceFeederAll",
					Use:       "list-price-feeder",
					Short:     "list all priceFeeder",
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              oracle.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: false, // use custom commands only until cosmos sdk v0.51
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod:      "FeedPrice",
					Use:            "feed-price [feed_price]",
					Short:          "Feed a new price",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "feed_price"}},
				},
				{
					RpcMethod:      "SetPriceFeeder",
					Use:            "set-price-feeder [isActive]",
					Short:          "Set a price feeder",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "is_active"}},
				},
				{
					RpcMethod:      "CreateAssetInfo",
					Use:            "create-asset-info [denom] [display] [band-ticker] [elys-ticker] [decimal]",
					Short:          "create a new asset info",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "denom"}, {ProtoField: "display"}, {ProtoField: "band_ticker"}, {ProtoField: "elys_ticker"}, {ProtoField: "decimal"}},
				},
				{
					RpcMethod: "DeletePriceFeeder",
					Use:       "delete-price-feeder",
					Short:     "Delete a priceFeeder",
				},
				{
					RpcMethod:      "FeedMultiplePrices",
					Use:            "feed-multiple-price [prices]",
					Short:          "Feed multiple prices",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "feed_prices"}},
				},
				{
					RpcMethod: "RemoveAssetInfo",
					Skip:      true, //authority gated
				},
				{
					RpcMethod: "AddPriceFeeders",
					Skip:      true, //authority gated
				},
				{
					RpcMethod: "RemovePriceFeeders",
					Skip:      true, //authority gated
				},
				{
					RpcMethod: "UpdateParams",
					Skip:      true, //authority gated
				},
			},
		},
	}
}
