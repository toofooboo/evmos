package cosmos_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdkvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	cosmosante "github.com/evmos/evmos/v19/app/ante/cosmos"
	"github.com/evmos/evmos/v19/testutil"
	"github.com/evmos/evmos/v19/testutil/integration/common/factory"
	utiltx "github.com/evmos/evmos/v19/testutil/tx"
	evmtypes "github.com/evmos/evmos/v19/x/evm/types"
)

func TestAuthzLimiterDecorator(t *testing.T) {
	testPrivKeys, testAddresses, err := generatePrivKeyAddressPairs(5)
	require.NoError(t, err)

	distantFuture := time.Date(9000, 1, 1, 0, 0, 0, 0, time.UTC)

	validator := sdk.ValAddress(testAddresses[4])
	stakingAuthDelegate, err := stakingtypes.NewStakeAuthorization([]sdk.ValAddress{validator}, nil, stakingtypes.AuthorizationType_AUTHORIZATION_TYPE_DELEGATE, nil)
	require.NoError(t, err)

	stakingAuthUndelegate, err := stakingtypes.NewStakeAuthorization([]sdk.ValAddress{validator}, nil, stakingtypes.AuthorizationType_AUTHORIZATION_TYPE_UNDELEGATE, nil)
	require.NoError(t, err)

	decorator := cosmosante.NewAuthzLimiterDecorator(
		sdk.MsgTypeURL(&evmtypes.MsgEthereumTx{}),
		sdk.MsgTypeURL(&stakingtypes.MsgUndelegate{}),
	)

	testCases := []struct {
		name        string
		msgs        []sdk.Msg
		checkTx     bool
		expectedErr error
	}{
		{
			"enabled msg - non blocked msg",
			[]sdk.Msg{
				banktypes.NewMsgSend(
					testAddresses[0],
					testAddresses[1],
					sdk.NewCoins(sdk.NewInt64Coin(evmtypes.DefaultEVMDenom, 100e6)),
				),
			},
			false,
			nil,
		},
		{
			"enabled msg MsgEthereumTx - blocked msg not wrapped in MsgExec",
			[]sdk.Msg{
				&evmtypes.MsgEthereumTx{},
			},
			false,
			nil,
		},
		{
			"enabled msg - blocked msg not wrapped in MsgExec",
			[]sdk.Msg{
				&stakingtypes.MsgCancelUnbondingDelegation{},
			},
			false,
			nil,
		},
		{
			"enabled msg - MsgGrant contains a non blocked msg",
			[]sdk.Msg{
				newMsgGrant(
					testAddresses[0],
					testAddresses[1],
					authz.NewGenericAuthorization(sdk.MsgTypeURL(&banktypes.MsgSend{})),
					&distantFuture,
				),
			},
			false,
			nil,
		},
		{
			"enabled msg - MsgGrant contains a non blocked msg",
			[]sdk.Msg{
				newMsgGrant(
					testAddresses[0],
					testAddresses[1],
					stakingAuthDelegate,
					&distantFuture,
				),
			},
			false,
			nil,
		},
		{
			"disabled msg - MsgGrant contains a blocked msg",
			[]sdk.Msg{
				newMsgGrant(
					testAddresses[0],
					testAddresses[1],
					authz.NewGenericAuthorization(sdk.MsgTypeURL(&evmtypes.MsgEthereumTx{})),
					&distantFuture,
				),
			},
			false,
			sdkerrors.ErrUnauthorized,
		},
		{
			"disabled msg - MsgGrant contains a blocked msg",
			[]sdk.Msg{
				newMsgGrant(
					testAddresses[0],
					testAddresses[1],
					stakingAuthUndelegate,
					&distantFuture,
				),
			},
			false,
			sdkerrors.ErrUnauthorized,
		},
		{
			"allowed msg - when a MsgExec contains a non blocked msg",
			[]sdk.Msg{
				newMsgExec(
					testAddresses[1],
					[]sdk.Msg{banktypes.NewMsgSend(
						testAddresses[0],
						testAddresses[3],
						sdk.NewCoins(sdk.NewInt64Coin(evmtypes.DefaultEVMDenom, 100e6)),
					)}),
			},
			false,
			nil,
		},
		{
			"disabled msg - MsgExec contains a blocked msg",
			[]sdk.Msg{
				newMsgExec(
					testAddresses[1],
					[]sdk.Msg{
						&evmtypes.MsgEthereumTx{},
					},
				),
			},
			false,
			sdkerrors.ErrUnauthorized,
		},
		{
			"disabled msg - surrounded by valid msgs",
			[]sdk.Msg{
				newMsgGrant(
					testAddresses[0],
					testAddresses[1],
					stakingAuthDelegate,
					&distantFuture,
				),
				newMsgExec(
					testAddresses[1],
					[]sdk.Msg{
						banktypes.NewMsgSend(
							testAddresses[0],
							testAddresses[3],
							sdk.NewCoins(sdk.NewInt64Coin(evmtypes.DefaultEVMDenom, 100e6)),
						),
						&evmtypes.MsgEthereumTx{},
					},
				),
			},
			false,
			sdkerrors.ErrUnauthorized,
		},
		{
			"disabled msg - nested MsgExec containing a blocked msg",
			[]sdk.Msg{
				createNestedMsgExec(
					testAddresses[1],
					2,
					[]sdk.Msg{
						&evmtypes.MsgEthereumTx{},
					},
				),
			},
			false,
			sdkerrors.ErrUnauthorized,
		},
		{
			"disabled msg - nested MsgGrant containing a blocked msg",
			[]sdk.Msg{
				newMsgExec(
					testAddresses[1],
					[]sdk.Msg{
						newMsgGrant(
							testAddresses[0],
							testAddresses[1],
							authz.NewGenericAuthorization(sdk.MsgTypeURL(&evmtypes.MsgEthereumTx{})),
							&distantFuture,
						),
					},
				),
			},
			false,
			sdkerrors.ErrUnauthorized,
		},
		{
			"disabled msg - nested MsgExec NOT containing a blocked msg but has more nesting levels than the allowed",
			[]sdk.Msg{
				createNestedMsgExec(
					testAddresses[1],
					6,
					[]sdk.Msg{
						banktypes.NewMsgSend(
							testAddresses[0],
							testAddresses[3],
							sdk.NewCoins(sdk.NewInt64Coin(evmtypes.DefaultEVMDenom, 100e6)),
						),
					},
				),
			},
			false,
			sdkerrors.ErrUnauthorized,
		},
		{
			"disabled msg - multiple two nested MsgExec messages NOT containing a blocked msg over the limit",
			[]sdk.Msg{
				createNestedMsgExec(
					testAddresses[1],
					5,
					[]sdk.Msg{
						banktypes.NewMsgSend(
							testAddresses[0],
							testAddresses[3],
							sdk.NewCoins(sdk.NewInt64Coin(evmtypes.DefaultEVMDenom, 100e6)),
						),
					},
				),
				createNestedMsgExec(
					testAddresses[1],
					5,
					[]sdk.Msg{
						banktypes.NewMsgSend(
							testAddresses[0],
							testAddresses[3],
							sdk.NewCoins(sdk.NewInt64Coin(evmtypes.DefaultEVMDenom, 100e6)),
						),
					},
				),
			},
			false,
			sdkerrors.ErrUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Case %s", tc.name), func(t *testing.T) {
			ctx := sdk.Context{}.WithIsCheckTx(tc.checkTx)
			tx, err := createTx(ctx, testPrivKeys[0], tc.msgs...)
			require.NoError(t, err)

			_, err = decorator.AnteHandle(ctx, tx, false, testutil.NextFn)
			if tc.expectedErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func (suite *AnteTestSuite) TestRejectMsgsInAuthz() {
	_, testAddresses, err := generatePrivKeyAddressPairs(10)
	suite.Require().NoError(err)

	var gasLimit uint64 = 1000000
	distantFuture := time.Date(9000, 1, 1, 0, 0, 0, 0, time.UTC)

	nw := suite.GetNetwork()
	// create a dummy MsgEthereumTx for the test
	// otherwise throws error that cannot unpack tx data
	msgEthereumTx := evmtypes.NewTx(&evmtypes.EvmTxArgs{
		ChainID:   nw.GetEIP155ChainID(),
		Nonce:     0,
		GasLimit:  gasLimit,
		GasFeeCap: nw.App.FeeMarketKeeper.GetBaseFee(nw.GetContext()),
		GasTipCap: big.NewInt(1),
		Input:     nil,
		Accesses:  &ethtypes.AccessList{},
	})

	newMsgGrant := func(msgTypeUrl string) *authz.MsgGrant {
		msg, err := authz.NewMsgGrant(
			testAddresses[0],
			testAddresses[1],
			authz.NewGenericAuthorization(msgTypeUrl),
			&distantFuture,
		)
		if err != nil {
			panic(err)
		}
		return msg
	}

	testcases := []struct {
		name         string
		msgs         []sdk.Msg
		expectedCode uint32
		isEIP712     bool
	}{
		{
			name:         "a MsgGrant with MsgEthereumTx typeURL on the authorization field is blocked",
			msgs:         []sdk.Msg{newMsgGrant(sdk.MsgTypeURL(&evmtypes.MsgEthereumTx{}))},
			expectedCode: sdkerrors.ErrUnauthorized.ABCICode(),
		},
		{
			name:         "a MsgGrant with MsgCreateVestingAccount typeURL on the authorization field is blocked",
			msgs:         []sdk.Msg{newMsgGrant(sdk.MsgTypeURL(&sdkvesting.MsgCreateVestingAccount{}))},
			expectedCode: sdkerrors.ErrUnauthorized.ABCICode(),
		},
		{
			name:         "a MsgGrant with MsgEthereumTx typeURL on the authorization field included on EIP712 tx is blocked",
			msgs:         []sdk.Msg{newMsgGrant(sdk.MsgTypeURL(&evmtypes.MsgEthereumTx{}))},
			expectedCode: sdkerrors.ErrUnauthorized.ABCICode(),
			isEIP712:     true,
		},
		{
			name: "a MsgExec with nested messages (valid: MsgSend and invalid: MsgEthereumTx) is blocked",
			msgs: []sdk.Msg{
				newMsgExec(
					testAddresses[1],
					[]sdk.Msg{
						banktypes.NewMsgSend(
							testAddresses[0],
							testAddresses[3],
							sdk.NewCoins(sdk.NewInt64Coin(evmtypes.DefaultEVMDenom, 100e6)),
						),
						msgEthereumTx,
					},
				),
			},
			expectedCode: sdkerrors.ErrUnauthorized.ABCICode(),
		},
		{
			name: "a MsgExec with nested MsgExec messages that has invalid messages is blocked",
			msgs: []sdk.Msg{
				createNestedMsgExec(
					testAddresses[1],
					2,
					[]sdk.Msg{
						msgEthereumTx,
					},
				),
			},
			expectedCode: sdkerrors.ErrUnauthorized.ABCICode(),
		},
		{
			name: "a MsgExec with more nested MsgExec messages than allowed and with valid messages is blocked",
			msgs: []sdk.Msg{
				createNestedMsgExec(
					testAddresses[1],
					6,
					[]sdk.Msg{
						banktypes.NewMsgSend(
							testAddresses[0],
							testAddresses[3],
							sdk.NewCoins(sdk.NewInt64Coin(evmtypes.DefaultEVMDenom, 100e6)),
						),
					},
				),
			},
			expectedCode: sdkerrors.ErrUnauthorized.ABCICode(),
		},
		{
			name: "two MsgExec messages NOT containing a blocked msg but between the two have more nesting than the allowed. Then, is blocked",
			msgs: []sdk.Msg{
				createNestedMsgExec(
					testAddresses[1],
					5,
					[]sdk.Msg{
						banktypes.NewMsgSend(
							testAddresses[0],
							testAddresses[3],
							sdk.NewCoins(sdk.NewInt64Coin(evmtypes.DefaultEVMDenom, 100e6)),
						),
					},
				),
				createNestedMsgExec(
					testAddresses[1],
					5,
					[]sdk.Msg{
						banktypes.NewMsgSend(
							testAddresses[0],
							testAddresses[3],
							sdk.NewCoins(sdk.NewInt64Coin(evmtypes.DefaultEVMDenom, 100e6)),
						),
					},
				),
			},
			expectedCode: sdkerrors.ErrUnauthorized.ABCICode(),
		},
	}

	for _, tc := range testcases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest()
			nw = suite.GetNetwork()
			var (
				tx  sdk.Tx
				err error
			)
			ctx := nw.GetContext()
			priv := suite.GetKeyring().GetPrivKey(0)

			if tc.isEIP712 {
				coinAmount := sdk.NewCoin(evmtypes.DefaultEVMDenom, math.NewInt(20))
				fees := sdk.NewCoins(coinAmount)
				cosmosTxArgs := utiltx.CosmosTxArgs{
					TxCfg:   suite.GetClientCtx().TxConfig,
					Priv:    priv,
					ChainID: ctx.ChainID(),
					Gas:     200000,
					Fees:    fees,
					Msgs:    tc.msgs,
				}

				tx, err = utiltx.CreateEIP712CosmosTx(
					ctx,
					nw.App,
					utiltx.EIP712TxArgs{
						CosmosTxArgs:       cosmosTxArgs,
						UseLegacyTypedData: true,
					},
				)
			} else {
				tx, err = suite.GetTxFactory().BuildCosmosTx(
					priv,
					factory.CosmosTxArgs{
						Gas:  &gasLimit,
						Msgs: tc.msgs,
					},
				)
			}
			suite.Require().NoError(err)

			txEncoder := suite.GetClientCtx().TxConfig.TxEncoder()
			bz, err := txEncoder(tx)
			suite.Require().NoError(err)

			resCheckTx, err := nw.App.CheckTx(
				&abci.RequestCheckTx{
					Tx:   bz,
					Type: abci.CheckTxType_New,
				},
			)
			suite.Require().NoError(err)
			suite.Require().Equal(resCheckTx.Code, tc.expectedCode, resCheckTx.Log)

			header := ctx.BlockHeader()
			blockRes, err := nw.App.FinalizeBlock(
				&abci.RequestFinalizeBlock{
					Height:             ctx.BlockHeight() + 1,
					Txs:                [][]byte{bz},
					Hash:               header.AppHash,
					NextValidatorsHash: header.NextValidatorsHash,
					ProposerAddress:    header.ProposerAddress,
					Time:               header.Time.Add(time.Second),
				},
			)
			suite.Require().NoError(err)
			suite.Require().Len(blockRes.TxResults, 1)
			txRes := blockRes.TxResults[0]
			suite.Require().Equal(txRes.Code, tc.expectedCode, txRes.Log)
		})
	}
}
