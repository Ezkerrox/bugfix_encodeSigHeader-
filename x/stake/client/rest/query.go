package rest

import (
	"net/http"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/stake/types"

	"github.com/gorilla/mux"
)

const storeName = "stake"

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {

	// Get all delegations from a delegator
	r.HandleFunc(
		"/stake/delegators/{delegatorAddr}/delegations",
		delegatorDelegationsHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Get all unbonding delegations from a delegator
	r.HandleFunc(
		"/stake/delegators/{delegatorAddr}/unbonding_delegations",
		delegatorUnbondingDelegationsHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Get all redelegations from a delegator
	r.HandleFunc(
		"/stake/delegators/{delegatorAddr}/redelegations",
		delegatorRedelegationsHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Get all staking txs (i.e msgs) from a delegator
	r.HandleFunc(
		"/stake/delegators/{delegatorAddr}/txs",
		delegatorTxsHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Query all validators that a delegator is bonded to
	r.HandleFunc(
		"/stake/delegators/{delegatorAddr}/validators",
		delegatorValidatorsHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Query a validator that a delegator is bonded to
	r.HandleFunc(
		"/stake/delegators/{delegatorAddr}/validators/{validatorAddr}",
		delegatorValidatorHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Query a delegation between a delegator and a validator
	r.HandleFunc(
		"/stake/delegators/{delegatorAddr}/delegations/{validatorAddr}",
		delegationHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Query all unbonding delegations between a delegator and a validator
	r.HandleFunc(
		"/stake/delegators/{delegatorAddr}/unbonding_delegations/{validatorAddr}",
		unbondingDelegationHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Get all validators
	r.HandleFunc(
		"/stake/validators",
		validatorsHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Get a single validator info
	r.HandleFunc(
		"/stake/validators/{validatorAddr}",
		validatorHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Get all unbonding delegations from a validator
	r.HandleFunc(
		"/stake/validators/{validatorAddr}/unbonding_delegations",
		validatorUnbondingDelegationsHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Get all outgoing redelegations from a validator
	r.HandleFunc(
		"/stake/validators/{validatorAddr}/redelegations",
		validatorRedelegationsHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Get the current state of the staking pool
	r.HandleFunc(
		"/stake/pool",
		poolHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Get the current staking parameter values
	r.HandleFunc(
		"/stake/parameters",
		paramsHandlerFn(cliCtx, cdc),
	).Methods("GET")

}

// HTTP request handler to query a delegator delegations
func delegatorDelegationsHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryDelegator(cliCtx, cdc, "custom/stake/delegatorDelegations")
}

// HTTP request handler to query a delegator unbonding delegations
func delegatorUnbondingDelegationsHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryDelegator(cliCtx, cdc, "custom/stake/delegatorUnbondingDelegations")
}

// HTTP request handler to query a delegator redelegations
func delegatorRedelegationsHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryDelegator(cliCtx, cdc, "custom/stake/delegatorRedelegations")
}

// HTTP request handler to query all staking txs (msgs) from a delegator
func delegatorTxsHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var typesQuerySlice []string
		vars := mux.Vars(r)
		delegatorAddr := vars["delegatorAddr"]

		_, err := sdk.AccAddressFromBech32(delegatorAddr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		node, err := cliCtx.GetNode()
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Get values from query

		typesQuery := r.URL.Query().Get("type")
		trimmedQuery := strings.TrimSpace(typesQuery)
		if len(trimmedQuery) != 0 {
			typesQuerySlice = strings.Split(trimmedQuery, " ")
		}

		noQuery := len(typesQuerySlice) == 0
		isBondTx := contains(typesQuerySlice, "bond")
		isUnbondTx := contains(typesQuerySlice, "unbond")
		isRedTx := contains(typesQuerySlice, "redelegate")
		var txs = []tx.Info{}
		var actions []string

		switch {
		case isBondTx:
			actions = append(actions, types.MsgDelegate{}.Type())
		case isUnbondTx:
			actions = append(actions, types.MsgBeginUnbonding{}.Type())
		case isRedTx:
			actions = append(actions, types.MsgBeginRedelegate{}.Type())
		case noQuery:
			actions = append(actions, types.MsgDelegate{}.Type())
			actions = append(actions, types.MsgBeginUnbonding{}.Type())
			actions = append(actions, types.MsgBeginRedelegate{}.Type())

		default:
			w.WriteHeader(http.StatusNoContent)
			return
		}

		for _, action := range actions {
			foundTxs, errQuery := queryTxs(node, cliCtx, cdc, action, delegatorAddr)
			if errQuery != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			}
			txs = append(txs, foundTxs...)
		}

		res, err := cdc.MarshalJSON(txs)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

// HTTP request handler to query an unbonding-delegation
func unbondingDelegationHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryBonds(cliCtx, cdc, "custom/stake/unbondingDelegation")
}

// HTTP request handler to query a delegation
func delegationHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryBonds(cliCtx, cdc, "custom/stake/delegation")
}

// HTTP request handler to query all delegator bonded validators
func delegatorValidatorsHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryDelegator(cliCtx, cdc, "custom/stake/delegatorValidators")
}

// HTTP request handler to get information from a currently bonded validator
func delegatorValidatorHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryBonds(cliCtx, cdc, "custom/stake/delegatorValidator")
}

// HTTP request handler to query list of validators
func validatorsHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := cliCtx.QueryWithData("custom/stake/validators", nil)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

// HTTP request handler to query the validator information from a given validator address
func validatorHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryValidator(cliCtx, cdc, "custom/stake/validator")
}

// HTTP request handler to query all unbonding delegations from a validator
func validatorUnbondingDelegationsHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryValidator(cliCtx, cdc, "custom/stake/validatorUnbondingDelegations")
}

// HTTP request handler to query all redelegations from a source validator
func validatorRedelegationsHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryValidator(cliCtx, cdc, "custom/stake/validatorRedelegations")
}

// HTTP request handler to query the pool information
func poolHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := cliCtx.QueryWithData("custom/stake/pool", nil)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

// HTTP request handler to query the staking params values
func paramsHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := cliCtx.QueryWithData("custom/stake/parameters", nil)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
