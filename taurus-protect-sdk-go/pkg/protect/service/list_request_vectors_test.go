package service

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/cache"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model/taurusnetwork"
)

// Cross-SDK wire vectors for every paged operation: canonical options → the exact decoded
// query multiset or body a list method sends, or an error before any request. The real
// generated client runs against an httptest transport answering `{}`.
//
// Every operation in the file must map to a Go method below or be listed in
// notWrappedOperations with a reason, and every canonical option must map to an option field:
// an unmapped one fails rather than being skipped.
const listRequestVectorsRelPath = "../../../../scripts/resources/list-request-vectors.json"

// The file's declared totals, asserted so a vector added without being consumed here fails.
const (
	listRequestVectorMethods = 53
	listRequestVectorCount   = 278
)

// notWrappedOperations are paged operations the Go SDK has no method for.
var notWrappedOperations = map[string]string{
	"FiatProviderService_GetFiatProviderCounterpartyAccounts": "no Go FiatService method wraps it",
	"FiatProviderService_GetFiatProviderOperations":           "no Go FiatService method wraps it",
	"WalletService_GetNFTCollectionBalances":                  "no Go method wraps it",
}

type listRequestVectorFile struct {
	Counts struct {
		Methods int `json:"methods"`
		Vectors int `json:"vectors"`
	} `json:"counts"`
	Methods map[string]struct {
		Method string `json:"method"`
		Path   string `json:"path"`
		Paging struct {
			Style string `json:"style"`
		} `json:"paging"`
		Rule string `json:"rule"`
	} `json:"methods"`
	Vectors []listRequestVector `json:"vectors"`
}

type listRequestVector struct {
	Operation   string                     `json:"operation"`
	PathParams  map[string]string          `json:"path_params"`
	Description string                     `json:"description"`
	Options     map[string]json.RawMessage `json:"options"`
	Expect      *struct {
		Query [][2]string `json:"query"`
		Body  any         `json:"body"`
	} `json:"expect"`
	ExpectError bool `json:"expect_error"`
}

// canonicalOptions are the vectors' snake_case option names, decoded.
type canonicalOptions struct {
	set                map[string]bool
	Limit              int64
	Offset             int64
	PageSize           int64
	Cursor             string
	ExcludeDisabled    bool
	FromCurrencyID     string
	ToCurrencyIDs      []string
	Statuses           []string
	ExternalRequestIDs []string
	Provider           string
	Label              string
	Currency           string
}

func decodeCanonicalOptions(raw map[string]json.RawMessage) (canonicalOptions, error) {
	o := canonicalOptions{set: map[string]bool{}}
	for name, value := range raw {
		var target any
		switch name {
		case "limit":
			target = &o.Limit
		case "offset":
			target = &o.Offset
		case "page_size":
			target = &o.PageSize
		case "cursor":
			target = &o.Cursor
		case "exclude_disabled":
			target = &o.ExcludeDisabled
		case "from_currency_id":
			target = &o.FromCurrencyID
		case "to_currency_ids":
			target = &o.ToCurrencyIDs
		case "statuses":
			target = &o.Statuses
		case "external_request_ids":
			target = &o.ExternalRequestIDs
		case "provider":
			target = &o.Provider
		case "label":
			target = &o.Label
		case "currency":
			target = &o.Currency
		default:
			return o, fmt.Errorf("canonical option %q has no Go option field", name)
		}
		if err := json.Unmarshal(value, target); err != nil {
			return o, fmt.Errorf("option %q: %w", name, err)
		}
		o.set[name] = true
	}
	return o, nil
}

// observation is what a list method returned about paging.
type observation struct {
	offset *model.Pagination
	cursor *model.CursorPage
	total  *int64
}

// listAdapter maps one operation to the Go method that wraps it.
type listAdapter struct {
	// options are the canonical options this method maps to option fields.
	options []string
	call    func(ctx context.Context, s *vectorServices, path map[string]string, o canonicalOptions) (observation, error)
}

func offsetObs(p *model.Pagination) observation { return observation{offset: p} }
func cursorObs(p model.CursorPage) observation  { return observation{cursor: &p} }

var (
	offsetOptions = []string{"limit", "offset"}
	cursorOptions = []string{"page_size", "cursor"}
)

// listAdapters is the explicit operationId → Go method table.
var listAdapters = map[string]listAdapter{
	"ActionService_GetActions": {offsetOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.actions.ListActions(ctx, &model.ListActionsOptions{Limit: o.Limit, Offset: o.Offset})
		if err != nil {
			return observation{}, err
		}
		return offsetObs(r.Pagination), nil
	}},
	"FeePayerService_GetFeePayers": {offsetOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.feePayers.ListFeePayers(ctx, &model.ListFeePayersOptions{Limit: o.Limit, Offset: o.Offset})
		if err != nil {
			return observation{}, err
		}
		return offsetObs(r.Pagination), nil
	}},
	"TransactionService_GetTransactions": {offsetOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		_, p, err := s.transactions.ListTransactions(ctx, &model.ListTransactionsOptions{Limit: o.Limit, Offset: o.Offset})
		return offsetObs(p), err
	}},
	"UserService_GetUsers": {offsetOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.users.ListUsers(ctx, &model.ListUsersOptions{Limit: o.Limit, Offset: o.Offset})
		if err != nil {
			return observation{}, err
		}
		return offsetObs(r.Pagination), nil
	}},
	"UserService_GetGroups": {offsetOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.groups.ListGroups(ctx, &model.ListGroupsOptions{Limit: o.Limit, Offset: o.Offset})
		if err != nil {
			return observation{}, err
		}
		return offsetObs(r.Pagination), nil
	}},
	"WalletService_GetAddresses": {[]string{"limit", "offset", "exclude_disabled"}, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		_, p, err := s.addresses.ListAddresses(ctx, &model.ListAddressesOptions{Limit: o.Limit, Offset: o.Offset, ExcludeDisabled: o.ExcludeDisabled})
		return offsetObs(p), err
	}},
	"WalletService_GetWalletsV2": {[]string{"limit", "offset", "exclude_disabled"}, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		_, p, err := s.wallets.ListWallets(ctx, &model.ListWalletsOptions{Limit: o.Limit, Offset: o.Offset, ExcludeDisabled: o.ExcludeDisabled})
		return offsetObs(p), err
	}},
	"WhitelistService_GetWhitelistedAddresses": {offsetOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.wla.ListWhitelistedAddresses(ctx, &model.ListWhitelistedAddressesOptions{Limit: o.Limit, Offset: o.Offset})
		if err != nil {
			return observation{}, err
		}
		return offsetObs(r.Pagination), nil
	}},
	"WhitelistService_GetWhitelistedAddressesForApproval": {offsetOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.wla.ListWhitelistedAddressesForApproval(ctx, &model.ListWhitelistedAddressesForApprovalOptions{Limit: o.Limit, Offset: o.Offset})
		if err != nil {
			return observation{}, err
		}
		return offsetObs(r.Pagination), nil
	}},
	"WhitelistService_GetWhitelistedContracts": {offsetOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.wca.ListWhitelistedAssets(ctx, &model.ListWhitelistedAssetsOptions{Limit: o.Limit, Offset: o.Offset})
		if err != nil {
			return observation{}, err
		}
		return offsetObs(r.Pagination), nil
	}},
	"WhitelistService_GetWhitelistedContractsForApproval": {offsetOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.wca.ListWhitelistedAssetsForApproval(ctx, &model.ListWhitelistedAssetsForApprovalOptions{Limit: o.Limit, Offset: o.Offset})
		if err != nil {
			return observation{}, err
		}
		return offsetObs(r.Pagination), nil
	}},

	"TransactionService_ExportTransactions": {[]string{"limit"}, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.transactions.ExportTransactions(ctx, &model.ExportTransactionsOptions{Limit: o.Limit})
		if err != nil {
			return observation{}, err
		}
		return observation{total: &r.TotalItems}, nil
	}},
	"PriceService_GetPricesHistory": {[]string{"limit"}, func(ctx context.Context, s *vectorServices, path map[string]string, o canonicalOptions) (observation, error) {
		_, err := s.prices.GetPriceHistory(ctx, &model.GetPriceHistoryOptions{Base: path["base"], Quote: path["quote"], Limit: o.Limit})
		return observation{}, err
	}},
	"PriceService_ExportPricesHistory": {[]string{"limit"}, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		_, err := s.prices.ExportPriceHistory(ctx, &model.ExportPriceHistoryOptions{Limit: o.Limit})
		return observation{}, err
	}},

	"AssetServiceV2_ListAssetOperationsV2": {cursorOptions, func(ctx context.Context, s *vectorServices, path map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.assets.ListAssetOperations(ctx, path["assetID"], &model.ListAssetOperationsOptions{PageSize: o.PageSize, Cursor: o.Cursor})
		if err != nil {
			return observation{}, err
		}
		return cursorObs(r.Page), nil
	}},
	"AssetServiceV2_QueryAssetAddressesV2": {cursorOptions, func(ctx context.Context, s *vectorServices, path map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.assets.QueryAssetAddresses(ctx, path["assetID"], &model.QueryAssetAddressesOptions{PageSize: o.PageSize, Cursor: o.Cursor})
		if err != nil {
			return observation{}, err
		}
		return cursorObs(r.Page), nil
	}},
	"AssetServiceV2_QueryAssetsV2": {cursorOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.assets.QueryAssets(ctx, &model.QueryAssetsOptions{PageSize: o.PageSize, Cursor: o.Cursor})
		if err != nil {
			return observation{}, err
		}
		return cursorObs(r.Page), nil
	}},
	"AuditService_GetAuditTrails": {cursorOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.audits.ListAuditTrails(ctx, &model.ListAuditTrailsOptions{PageSize: o.PageSize, Cursor: o.Cursor})
		if err != nil {
			return observation{}, err
		}
		return cursorObs(r.Page), nil
	}},
	"ChangeService_GetChanges": {cursorOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.changes.ListChanges(ctx, &model.ListChangesOptions{PageSize: o.PageSize, Cursor: o.Cursor})
		if err != nil {
			return observation{}, err
		}
		return cursorObs(r.Page), nil
	}},
	"ChangeService_GetChangesForApproval": {cursorOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.changes.ListChangesForApproval(ctx, &model.ListChangesForApprovalOptions{PageSize: o.PageSize, Cursor: o.Cursor})
		if err != nil {
			return observation{}, err
		}
		return cursorObs(r.Page), nil
	}},
	"EarnService_GetRewards": {cursorOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.earn.ListRewards(ctx, &model.ListEarnRewardsOptions{PageSize: o.PageSize, Cursor: o.Cursor})
		if err != nil {
			return observation{}, err
		}
		return cursorObs(r.Page), nil
	}},
	"ExchangeService_GetExchanges": {cursorOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.exchanges.ListExchanges(ctx, &model.ListExchangesOptions{PageSize: o.PageSize, Cursor: o.Cursor})
		if err != nil {
			return observation{}, err
		}
		return cursorObs(r.Page), nil
	}},
	"FiatProviderService_GetFiatProviderAccounts": {[]string{"page_size", "cursor", "provider", "label"}, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.fiat.ListFiatProviderAccounts(ctx, &model.ListFiatProviderAccountsOptions{
			Provider: o.Provider, Label: o.Label, PageSize: o.PageSize, Cursor: o.Cursor})
		if err != nil {
			return observation{}, err
		}
		return cursorObs(r.Page), nil
	}},
	"FiatProviderService_GetFiatProviderEntities": {cursorOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.fiat.ListFiatProviderEntities(ctx, &model.ListFiatProviderEntitiesOptions{PageSize: o.PageSize, Cursor: o.Cursor})
		if err != nil {
			return observation{}, err
		}
		return cursorObs(r.Page), nil
	}},
	"PriceService_QueryPricesV2": {[]string{"page_size", "cursor", "from_currency_id", "to_currency_ids"}, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.prices.ListPrices(ctx, &model.ListPricesOptions{
			PageSize: o.PageSize, Cursor: o.Cursor, FromCurrencyID: o.FromCurrencyID, ToCurrencyIDs: o.ToCurrencyIDs})
		if err != nil {
			return observation{}, err
		}
		return cursorObs(r.Page), nil
	}},
	"RequestService_GetRequestsForApprovalV2": {[]string{"page_size", "cursor", "statuses", "external_request_ids"}, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.requests.ListRequestsForApproval(ctx, &model.ListRequestsOptions{PageSize: o.PageSize, Cursor: o.Cursor, Statuses: o.Statuses, ExternalRequestIDs: o.ExternalRequestIDs})
		if err != nil {
			return observation{}, err
		}
		return cursorObs(r.Page), nil
	}},
	"RequestService_GetRequestsV2": {cursorOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.requests.ListRequests(ctx, &model.ListRequestsOptions{PageSize: o.PageSize, Cursor: o.Cursor})
		if err != nil {
			return observation{}, err
		}
		return cursorObs(r.Page), nil
	}},
	"RuleService_GetBusinessRulesV2": {cursorOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.businessRules.ListBusinessRules(ctx, &model.ListBusinessRulesOptions{PageSize: o.PageSize, Cursor: o.Cursor})
		if err != nil {
			return observation{}, err
		}
		return cursorObs(r.Page), nil
	}},
	"StakingService_GetStakeAccounts": {cursorOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.staking.ListStakeAccounts(ctx, &model.ListStakeAccountsOptions{PageSize: o.PageSize, Cursor: o.Cursor})
		if err != nil {
			return observation{}, err
		}
		return cursorObs(r.Page), nil
	}},
	"StatisticsService_GetAggregatedTagStats": {cursorOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.statistics.ListTagStatistics(ctx, &model.ListTagStatisticsOptions{PageSize: o.PageSize, Cursor: o.Cursor})
		if err != nil {
			return observation{}, err
		}
		return cursorObs(r.Page), nil
	}},
	"StatisticsService_GetPortfolioStatisticsHistory": {cursorOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.statistics.GetPortfolioStatisticsHistory(ctx, &model.GetPortfolioStatisticsHistoryOptions{PageSize: o.PageSize, Cursor: o.Cursor})
		if err != nil {
			return observation{}, err
		}
		return cursorObs(r.Page), nil
	}},
	"TaurusNetworkService_GetLendingAgreements": {cursorOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.lending.ListLendingAgreements(ctx, &taurusnetwork.ListLendingAgreementsOptions{PageSize: o.PageSize, Cursor: o.Cursor})
		if err != nil {
			return observation{}, err
		}
		return cursorObs(r.Page), nil
	}},
	"TaurusNetworkService_GetLendingAgreementsForApproval": {cursorOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.lending.ListLendingAgreementsForApproval(ctx, &taurusnetwork.ListLendingAgreementsForApprovalOptions{PageSize: o.PageSize, Cursor: o.Cursor})
		if err != nil {
			return observation{}, err
		}
		return cursorObs(r.Page), nil
	}},
	"TaurusNetworkService_GetLendingOffers": {cursorOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.lending.ListLendingOffers(ctx, &taurusnetwork.ListLendingOffersOptions{PageSize: o.PageSize, Cursor: o.Cursor})
		if err != nil {
			return observation{}, err
		}
		return cursorObs(r.Page), nil
	}},
	"TaurusNetworkService_GetPledgeActions": {cursorOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		_, p, err := s.pledges.ListPledgeActions(ctx, &taurusnetwork.ListPledgeActionsOptions{PageSize: o.PageSize, Cursor: o.Cursor})
		return observation{cursor: p}, err
	}},
	"TaurusNetworkService_GetPledgeActionsForApproval": {cursorOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		_, p, err := s.pledges.ListPledgeActionsForApproval(ctx, &taurusnetwork.ListPledgeActionsForApprovalOptions{PageSize: o.PageSize, Cursor: o.Cursor})
		return observation{cursor: p}, err
	}},
	"TaurusNetworkService_GetPledges": {cursorOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		_, p, err := s.pledges.ListPledges(ctx, &taurusnetwork.ListPledgesOptions{PageSize: o.PageSize, Cursor: o.Cursor})
		return observation{cursor: p}, err
	}},
	"TaurusNetworkService_GetPledgesWithdrawals": {cursorOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		_, p, err := s.pledges.ListPledgeWithdrawals(ctx, &taurusnetwork.ListPledgeWithdrawalsOptions{PageSize: o.PageSize, Cursor: o.Cursor})
		return observation{cursor: p}, err
	}},
	"TaurusNetworkService_GetSettlements": {cursorOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.settlements.ListSettlements(ctx, &taurusnetwork.ListSettlementsOptions{PageSize: o.PageSize, Cursor: o.Cursor})
		if err != nil {
			return observation{}, err
		}
		return cursorObs(r.Page), nil
	}},
	"TaurusNetworkService_GetSettlementsForApproval": {cursorOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.settlements.ListSettlementsForApproval(ctx, &taurusnetwork.ListSettlementsForApprovalOptions{PageSize: o.PageSize, Cursor: o.Cursor})
		if err != nil {
			return observation{}, err
		}
		return cursorObs(r.Page), nil
	}},
	"TaurusNetworkService_GetSharedAddresses": {cursorOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.sharing.ListSharedAddresses(ctx, &taurusnetwork.ListSharedAddressesOptions{PageSize: o.PageSize, Cursor: o.Cursor})
		if err != nil {
			return observation{}, err
		}
		return cursorObs(r.Page), nil
	}},
	"TaurusNetworkService_GetSharedAssets": {cursorOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.sharing.ListSharedAssets(ctx, &taurusnetwork.ListSharedAssetsOptions{PageSize: o.PageSize, Cursor: o.Cursor})
		if err != nil {
			return observation{}, err
		}
		return cursorObs(r.Page), nil
	}},
	"WalletService_GetAssetAddresses": {[]string{"page_size", "cursor", "currency"}, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.assets.GetAssetAddresses(ctx, &model.GetAssetAddressesRequest{
			Asset: model.AssetFilter{Currency: o.Currency}, PageSize: o.PageSize, Cursor: o.Cursor})
		if err != nil {
			return observation{}, err
		}
		return cursorObs(r.Page), nil
	}},
	"WalletService_GetAssetWallets": {[]string{"page_size", "cursor", "currency"}, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.assets.GetAssetWallets(ctx, &model.GetAssetWalletsRequest{
			Asset: model.AssetFilter{Currency: o.Currency}, PageSize: o.PageSize, Cursor: o.Cursor})
		if err != nil {
			return observation{}, err
		}
		return cursorObs(r.Page), nil
	}},
	"WalletService_GetBalances": {cursorOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.balances.GetBalances(ctx, &model.GetBalancesOptions{PageSize: o.PageSize, Cursor: o.Cursor})
		if err != nil {
			return observation{}, err
		}
		return cursorObs(r.Page), nil
	}},
	"WalletService_GetReservations": {cursorOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.reservations.ListReservations(ctx, &model.ListReservationsOptions{PageSize: o.PageSize, Cursor: o.Cursor})
		if err != nil {
			return observation{}, err
		}
		return cursorObs(r.Page), nil
	}},
	"WebhookService_GetWebhookCalls": {cursorOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.webhookCalls.ListWebhookCalls(ctx, &model.ListWebhookCallsOptions{PageSize: o.PageSize, Cursor: o.Cursor})
		if err != nil {
			return observation{}, err
		}
		return cursorObs(r.Page), nil
	}},
	"WebhookService_GetWebhooks": {cursorOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.webhooks.ListWebhooks(ctx, &model.ListWebhooksOptions{PageSize: o.PageSize, Cursor: o.Cursor})
		if err != nil {
			return observation{}, err
		}
		return cursorObs(r.Page), nil
	}},

	"RuleService_GetRulesHistory": {cursorOptions, func(ctx context.Context, s *vectorServices, _ map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.governance.GetRulesHistory(ctx, &model.ListRulesHistoryOptions{PageSize: o.PageSize, Cursor: o.Cursor})
		if err != nil {
			return observation{}, err
		}
		return cursorObs(r.Page), nil
	}},
	"WalletService_GetWalletTokens": {cursorOptions, func(ctx context.Context, s *vectorServices, path map[string]string, o canonicalOptions) (observation, error) {
		r, err := s.wallets.GetWalletTokens(ctx, path["id"], &model.GetWalletTokensOptions{PageSize: o.PageSize, Cursor: o.Cursor})
		if err != nil {
			return observation{}, err
		}
		return cursorObs(r.Page), nil
	}},
}

// recordedRequest is one request the fake transport received.
type recordedRequest struct {
	method string
	path   string
	query  string
	body   []byte
}

// vectorServer is the transport stub: it records every request and answers with a fixed body.
type vectorServer struct {
	mu       sync.Mutex
	reply    string
	requests []recordedRequest
}

func (v *vectorServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	v.mu.Lock()
	v.requests = append(v.requests, recordedRequest{method: r.Method, path: r.URL.Path, query: r.URL.RawQuery, body: body})
	reply := v.reply
	v.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	_, _ = io.WriteString(w, reply)
}

func (v *vectorServer) reset(reply string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.reply = reply
	v.requests = nil
}

func (v *vectorServer) recorded() []recordedRequest {
	v.mu.Lock()
	defer v.mu.Unlock()
	return append([]recordedRequest(nil), v.requests...)
}

// vectorServices are the real services, wired through the real generated client to the stub.
type vectorServices struct {
	server        *vectorServer
	actions       *ActionService
	feePayers     *FeePayerService
	transactions  *TransactionService
	users         *UserService
	groups        *GroupService
	addresses     *AddressService
	wallets       *WalletService
	wla           *WhitelistedAddressService
	wca           *WhitelistedAssetService
	prices        *PriceService
	assets        *AssetService
	audits        *AuditService
	changes       *ChangeService
	earn          *EarnService
	exchanges     *ExchangeService
	fiat          *FiatService
	requests      *RequestService
	businessRules *BusinessRuleService
	staking       *StakingService
	statistics    *StatisticsService
	lending       *TaurusNetworkLendingService
	pledges       *TaurusNetworkPledgeService
	settlements   *TaurusNetworkSettlementService
	sharing       *TaurusNetworkSharingService
	balances      *BalanceService
	reservations  *ReservationService
	webhookCalls  *WebhookCallService
	webhooks      *WebhookService
	governance    *GovernanceRuleService
}

func newVectorServices(t *testing.T) *vectorServices {
	t.Helper()
	server := &vectorServer{reply: `{}`}
	srv := httptest.NewServer(server)
	t.Cleanup(srv.Close)

	cfg := openapi.NewConfiguration()
	cfg.Servers = []openapi.ServerConfiguration{{URL: srv.URL}}
	cfg.HTTPClient = srv.Client()
	client := openapi.NewAPIClient(cfg)

	// The verified lists read the rules container only through the cache; an empty verified
	// container keeps that off the wire, so each call sends exactly the list request.
	rules := cache.NewRulesContainerCache(time.Hour, func(context.Context) (*model.DecodedRulesContainer, error) {
		return &model.DecodedRulesContainer{}, nil
	})
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	superAdmins := []*ecdsa.PublicKey{&key.PublicKey}

	addresses := NewAddressService(client, rules)
	wla := NewWhitelistedAddressServiceWithVerification(client, &WhitelistedAddressServiceConfig{SuperAdminKeys: superAdmins, MinValidSignatures: 1})

	return &vectorServices{
		server:        server,
		actions:       NewActionService(client),
		feePayers:     NewFeePayerService(client),
		transactions:  NewTransactionService(client),
		users:         NewUserService(client),
		groups:        NewGroupService(client),
		addresses:     addresses,
		wallets:       NewWalletService(client),
		wla:           wla,
		wca:           NewWhitelistedAssetServiceWithVerification(client, &WhitelistedAssetServiceConfig{SuperAdminKeys: superAdmins, MinValidSignatures: 1}),
		prices:        NewPriceService(client, rules),
		assets:        NewAssetService(client, rules, addresses, wla),
		audits:        NewAuditService(client),
		changes:       NewChangeService(client),
		earn:          NewEarnService(client),
		exchanges:     NewExchangeService(client),
		fiat:          NewFiatService(client),
		requests:      NewRequestService(client),
		businessRules: NewBusinessRuleService(client),
		staking:       NewStakingService(client),
		statistics:    NewStatisticsService(client),
		lending:       NewTaurusNetworkLendingService(client),
		pledges:       NewTaurusNetworkPledgeService(client),
		settlements:   NewTaurusNetworkSettlementService(client),
		sharing:       NewTaurusNetworkSharingService(client),
		balances:      NewBalanceService(client),
		reservations:  NewReservationService(client),
		webhookCalls:  NewWebhookCallService(client),
		webhooks:      NewWebhookService(client),
		governance:    NewGovernanceRuleServiceWithVerification(client, &GovernanceRuleServiceConfig{SuperAdminKeys: superAdmins, MinValidSignatures: 1}),
	}
}

func loadListRequestVectors(t *testing.T) listRequestVectorFile {
	t.Helper()
	path := filepath.Clean(listRequestVectorsRelPath)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read shared list-request vectors %s: %v", path, err)
	}
	var file listRequestVectorFile
	if err := json.Unmarshal(raw, &file); err != nil {
		t.Fatalf("cannot parse %s: %v", path, err)
	}
	if file.Counts.Methods != listRequestVectorMethods || len(file.Methods) != file.Counts.Methods {
		t.Fatalf("methods: file declares %d, holds %d, this loader expects %d",
			file.Counts.Methods, len(file.Methods), listRequestVectorMethods)
	}
	if file.Counts.Vectors != listRequestVectorCount || len(file.Vectors) != file.Counts.Vectors {
		t.Fatalf("vectors: file declares %d, holds %d, this loader expects %d",
			file.Counts.Vectors, len(file.Vectors), listRequestVectorCount)
	}
	return file
}

// TestListRequestVectorsCoverEveryOperation fails on an operation with neither a Go method nor a
// recorded reason, and on an adapter or exclusion the file no longer names.
func TestListRequestVectorsCoverEveryOperation(t *testing.T) {
	file := loadListRequestVectors(t)
	for op := range file.Methods {
		_, wrapped := listAdapters[op]
		_, excluded := notWrappedOperations[op]
		switch {
		case wrapped && excluded:
			t.Errorf("%s is both wrapped and listed as not wrapped", op)
		case !wrapped && !excluded:
			t.Errorf("%s has no Go adapter and no recorded reason", op)
		}
	}
	for op := range listAdapters {
		if _, ok := file.Methods[op]; !ok {
			t.Errorf("adapter %s names an operation the vectors do not have", op)
		}
	}
	for op := range notWrappedOperations {
		if _, ok := file.Methods[op]; !ok {
			t.Errorf("not-wrapped entry %s names an operation the vectors do not have", op)
		}
	}
}

func TestListRequestVectors(t *testing.T) {
	file := loadListRequestVectors(t)
	svc := newVectorServices(t)
	exercised := map[string]map[string]bool{}

	for i, v := range file.Vectors {
		v := v
		t.Run(fmt.Sprintf("%03d_%s_%s", i, v.Operation, v.Description), func(t *testing.T) {
			if _, skip := notWrappedOperations[v.Operation]; skip {
				t.Skipf("%s: %s", v.Operation, notWrappedOperations[v.Operation])
			}
			adapter, ok := listAdapters[v.Operation]
			if !ok {
				t.Fatalf("no adapter for %s", v.Operation)
			}
			method := file.Methods[v.Operation]

			opts, err := decodeCanonicalOptions(v.Options)
			if err != nil {
				t.Fatal(err)
			}
			for name := range opts.set {
				if !containsString(adapter.options, name) {
					t.Fatalf("option %q is not mapped by the %s adapter", name, v.Operation)
				}
				if exercised[v.Operation] == nil {
					exercised[v.Operation] = map[string]bool{}
				}
				exercised[v.Operation][name] = true
			}

			svc.server.reset(`{}`)
			_, callErr := adapter.call(context.Background(), svc, v.PathParams, opts)
			requests := svc.server.recorded()

			if v.ExpectError {
				if callErr == nil {
					t.Fatal("expected a validation error")
				}
				if len(requests) != 0 {
					t.Fatalf("a rejected call must send nothing, sent %d request(s)", len(requests))
				}
				return
			}
			if callErr != nil {
				t.Fatalf("unexpected error: %v", callErr)
			}
			if len(requests) != 1 {
				t.Fatalf("want exactly 1 request, got %d", len(requests))
			}
			assertVectorRequest(t, requests[0], method.Method, substitutePath(method.Path, v.PathParams), v)
		})
	}

	// Both ways: an adapter mapping an option no vector exercises is an unverified mapping.
	for op, adapter := range listAdapters {
		for _, name := range adapter.options {
			if !exercised[op][name] {
				t.Errorf("%s maps option %q but no vector exercises it", op, name)
			}
		}
	}
}

func assertVectorRequest(t *testing.T, got recordedRequest, method, path string, v listRequestVector) {
	t.Helper()
	if got.method != method {
		t.Errorf("method = %s, want %s", got.method, method)
	}
	if got.path != path {
		t.Errorf("path = %s, want %s", got.path, path)
	}

	switch {
	case v.Expect == nil:
		t.Fatal("vector has neither expect nor expect_error")
	case v.Expect.Body != nil:
		if got.query != "" {
			t.Errorf("a body operation must send no query, sent %q", got.query)
		}
		var body any
		if err := json.Unmarshal(got.body, &body); err != nil {
			t.Fatalf("request body %q is not JSON: %v", got.body, err)
		}
		if !reflect.DeepEqual(body, v.Expect.Body) {
			t.Errorf("body = %s, want %s", got.body, mustJSON(t, v.Expect.Body))
		}
	default:
		if len(strings.TrimSpace(string(got.body))) > 0 {
			t.Errorf("a query operation must send no body, sent %q", got.body)
		}
		values, err := url.ParseQuery(got.query)
		if err != nil {
			t.Fatalf("query %q does not parse: %v", got.query, err)
		}
		gotPairs := queryPairs(values)
		wantPairs := make([]string, 0, len(v.Expect.Query))
		for _, p := range v.Expect.Query {
			wantPairs = append(wantPairs, p[0]+"="+p[1])
		}
		sort.Strings(wantPairs)
		if !reflect.DeepEqual(gotPairs, wantPairs) {
			t.Errorf("query = %v, want %v (raw %q)", gotPairs, wantPairs, got.query)
		}
	}
}

// queryPairs is the decoded query as a sorted name=value multiset.
func queryPairs(values url.Values) []string {
	pairs := []string{}
	for name, list := range values {
		for _, value := range list {
			pairs = append(pairs, name+"="+value)
		}
	}
	sort.Strings(pairs)
	return pairs
}

func substitutePath(path string, params map[string]string) string {
	for name, value := range params {
		path = strings.ReplaceAll(path, "{"+name+"}", url.PathEscape(value))
	}
	return path
}

func containsString(list []string, s string) bool {
	for _, item := range list {
		if item == s {
			return true
		}
	}
	return false
}

func mustJSON(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
