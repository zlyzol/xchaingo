package ethereum

import (
	"github.com/zlyzol/xchaingo/common"

)

type EthereumClient interface {
	/*
	call<T>(params: CallParams): Promise<T>
	estimateCall(asset: EstimateCallParams): Promise<BigNumber>
	estimateGasPrices(): Promise<GasPrices>
	estimateGasLimit(params: TxParams): Promise<BigNumber>
	estimateFeesWithGasPricesAndLimits(params: TxParams): Promise<FeesWithGasPricesAndLimits>
	estimateApprove(params: EstimateApproveParams): Promise<BigNumber>
	isApproved(params: IsApprovedParams): Promise<boolean>
	approve(params: ApproveParams): Promise<TransactionResponse>
	// `getFees` of `BaseXChainClient` needs to be overridden
	getFees(params: TxParams): Promise<Fees>
	*/
}

func NewEthereumClient(network common.Network, phrase string) (EthereumClient, error) {
	return nil, nil
}