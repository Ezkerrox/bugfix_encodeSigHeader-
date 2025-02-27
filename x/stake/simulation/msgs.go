package simulation

import (
	"fmt"
	"math/rand"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/cosmos/cosmos-sdk/x/mock/simulation"
	"github.com/cosmos/cosmos-sdk/x/stake"
)

// SimulateMsgCreateValidator
func SimulateMsgCreateValidator(m auth.AccountKeeper, k stake.Keeper) simulation.Operation {
	handler := stake.NewStakeHandler(k)
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account, event func(string)) (
		action string, fOp []simulation.FutureOperation, err error) {

		denom := k.GetParams(ctx).BondDenom
		description := stake.Description{
			Moniker: simulation.RandStringOfLength(r, 10),
		}

		maxCommission := int64(10)
		commission := stake.NewCommissionMsg(
			sdk.NewDecWithPrec(simulation.RandomAmount(r, maxCommission), 1),
			sdk.NewDecWithPrec(simulation.RandomAmount(r, maxCommission), 1),
			sdk.NewDecWithPrec(simulation.RandomAmount(r, maxCommission), 1),
		)

		acc := simulation.RandomAcc(r, accs)
		address := sdk.ValAddress(acc.Address)
		amount := m.GetAccount(ctx, acc.Address).GetCoins().AmountOf(denom)
		if amount > 0 {
			amount = simulation.RandomAmount(r, amount)
		}

		if amount == 0 {
			return "no-operation", nil, nil
		}

		msg := stake.MsgCreateValidator{
			Description:   description,
			Commission:    commission,
			ValidatorAddr: address,
			DelegatorAddr: acc.Address,
			PubKey:        acc.PubKey,
			Delegation:    sdk.NewCoin(denom, amount),
		}

		if msg.ValidateBasic() != nil {
			return "", nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}

		ctx, write := ctx.CacheContext()
		result := handler(ctx, msg)
		if result.IsOK() {
			write()
		}

		event(fmt.Sprintf("stake/MsgCreateValidator/%v", result.IsOK()))

		// require.True(t, result.IsOK(), "expected OK result but instead got %v", result)
		action = fmt.Sprintf("TestMsgCreateValidator: ok %v, msg %s", result.IsOK(), msg.GetSignBytes())
		return action, nil, nil
	}
}

// SimulateMsgEditValidator
func SimulateMsgEditValidator(k stake.Keeper) simulation.Operation {
	handler := stake.NewStakeHandler(k)
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account, event func(string)) (
		action string, fOp []simulation.FutureOperation, err error) {

		description := stake.Description{
			Moniker:  simulation.RandStringOfLength(r, 10),
			Identity: simulation.RandStringOfLength(r, 10),
			Website:  simulation.RandStringOfLength(r, 10),
			Details:  simulation.RandStringOfLength(r, 10),
		}

		maxCommission := int64(10)
		newCommissionRate := sdk.NewDecWithPrec(simulation.RandomAmount(r, maxCommission), 1)

		acc := simulation.RandomAcc(r, accs)
		address := sdk.ValAddress(acc.Address)
		msg := stake.MsgEditValidator{
			Description:    description,
			ValidatorAddr:  address,
			CommissionRate: &newCommissionRate,
		}

		if msg.ValidateBasic() != nil {
			return "", nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}

		ctx, write := ctx.CacheContext()
		result := handler(ctx, msg)
		if result.IsOK() {
			write()
		}
		event(fmt.Sprintf("stake/MsgEditValidator/%v", result.IsOK()))
		action = fmt.Sprintf("TestMsgEditValidator: ok %v, msg %s", result.IsOK(), msg.GetSignBytes())
		return action, nil, nil
	}
}

// SimulateMsgDelegate
func SimulateMsgDelegate(m auth.AccountKeeper, k stake.Keeper) simulation.Operation {
	handler := stake.NewStakeHandler(k)
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account, event func(string)) (
		action string, fOp []simulation.FutureOperation, err error) {

		denom := k.GetParams(ctx).BondDenom
		validatorAcc := simulation.RandomAcc(r, accs)
		validatorAddress := sdk.ValAddress(validatorAcc.Address)
		delegatorAcc := simulation.RandomAcc(r, accs)
		delegatorAddress := delegatorAcc.Address
		amount := m.GetAccount(ctx, delegatorAddress).GetCoins().AmountOf(denom)
		if amount > 0 {
			amount = simulation.RandomAmount(r, amount)
		}
		if amount == 0 {
			return "no-operation", nil, nil
		}
		msg := stake.MsgDelegate{
			DelegatorAddr: delegatorAddress,
			ValidatorAddr: validatorAddress,
			Delegation:    sdk.NewCoin(denom, amount),
		}
		if msg.ValidateBasic() != nil {
			return "", nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}
		ctx, write := ctx.CacheContext()
		result := handler(ctx, msg)
		if result.IsOK() {
			write()
		}
		event(fmt.Sprintf("stake/MsgDelegate/%v", result.IsOK()))
		action = fmt.Sprintf("TestMsgDelegate: ok %v, msg %s", result.IsOK(), msg.GetSignBytes())
		return action, nil, nil
	}
}

// SimulateMsgBeginUnbonding
func SimulateMsgBeginUnbonding(m auth.AccountKeeper, k stake.Keeper) simulation.Operation {
	handler := stake.NewStakeHandler(k)
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account, event func(string)) (
		action string, fOp []simulation.FutureOperation, err error) {

		denom := k.GetParams(ctx).BondDenom
		validatorAcc := simulation.RandomAcc(r, accs)
		validatorAddress := sdk.ValAddress(validatorAcc.Address)
		delegatorAcc := simulation.RandomAcc(r, accs)
		delegatorAddress := delegatorAcc.Address
		amount := m.GetAccount(ctx, delegatorAddress).GetCoins().AmountOf(denom)
		if amount > 0 {
			amount = simulation.RandomAmount(r, amount)
		}
		if amount == 0 {
			return "no-operation", nil, nil
		}
		msg := stake.MsgBeginUnbonding{
			DelegatorAddr: delegatorAddress,
			ValidatorAddr: validatorAddress,
			SharesAmount:  sdk.NewDecFromInt(amount),
		}
		if msg.ValidateBasic() != nil {
			return "", nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}
		ctx, write := ctx.CacheContext()
		result := handler(ctx, msg)
		if result.IsOK() {
			write()
		}
		event(fmt.Sprintf("stake/MsgBeginUnbonding/%v", result.IsOK()))
		action = fmt.Sprintf("TestMsgBeginUnbonding: ok %v, msg %s", result.IsOK(), msg.GetSignBytes())
		return action, nil, nil
	}
}

// SimulateMsgBeginRedelegate
func SimulateMsgBeginRedelegate(m auth.AccountKeeper, k stake.Keeper) simulation.Operation {
	handler := stake.NewStakeHandler(k)
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account, event func(string)) (
		action string, fOp []simulation.FutureOperation, err error) {

		denom := k.GetParams(ctx).BondDenom
		sourceValidatorAcc := simulation.RandomAcc(r, accs)
		sourceValidatorAddress := sdk.ValAddress(sourceValidatorAcc.Address)
		destValidatorAcc := simulation.RandomAcc(r, accs)
		destValidatorAddress := sdk.ValAddress(destValidatorAcc.Address)
		delegatorAcc := simulation.RandomAcc(r, accs)
		delegatorAddress := delegatorAcc.Address
		// TODO
		amount := m.GetAccount(ctx, delegatorAddress).GetCoins().AmountOf(denom)
		if amount > 0 {
			amount = simulation.RandomAmount(r, amount)
		}
		if amount == 0 {
			return "no-operation", nil, nil
		}
		msg := stake.MsgBeginRedelegate{
			DelegatorAddr:    delegatorAddress,
			ValidatorSrcAddr: sourceValidatorAddress,
			ValidatorDstAddr: destValidatorAddress,
			SharesAmount:     sdk.NewDecFromInt(amount),
		}
		if msg.ValidateBasic() != nil {
			return "", nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}
		ctx, write := ctx.CacheContext()
		result := handler(ctx, msg)
		if result.IsOK() {
			write()
		}
		event(fmt.Sprintf("stake/MsgBeginRedelegate/%v", result.IsOK()))
		action = fmt.Sprintf("TestMsgBeginRedelegate: %s", msg.GetSignBytes())
		return action, nil, nil
	}
}

// Setup
// nolint: errcheck
func Setup(mapp *mock.App, k stake.Keeper) simulation.RandSetup {
	return func(r *rand.Rand, accs []simulation.Account) {
		ctx := mapp.NewContext(sdk.RunTxModeDeliver, abci.Header{})
		gen := stake.DefaultGenesisState()
		stake.InitGenesis(ctx, k, gen)
		params := k.GetParams(ctx)
		denom := params.BondDenom
		var loose int64
		mapp.AccountKeeper.IterateAccounts(ctx, func(acc sdk.Account) bool {
			balance := simulation.RandomAmount(r, int64(1000000))
			acc.SetCoins(acc.GetCoins().Plus(sdk.Coins{sdk.NewCoin(denom, balance)}))
			mapp.AccountKeeper.SetAccount(ctx, acc)
			loose = loose + balance
			return false
		})
		pool := k.GetPool(ctx)
		pool.LooseTokens = pool.LooseTokens.Add(sdk.NewDec(loose))
		k.SetPool(ctx, pool)
	}
}
