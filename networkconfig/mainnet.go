package networkconfig

import (
	"math/big"

	spectypes "github.com/bloxapp/ssv-spec/types"

	"github.com/bloxapp/ssv/protocol/v2/blockchain/beacon"
)

var Mainnet = NetworkConfig{
	Name:                 "mainnet",
	Beacon:               beacon.NewNetwork(spectypes.MainNetwork),
	Domain:               spectypes.GenesisMainnet,
	GenesisEpoch:         1,
	ETH1SyncOffset:       new(big.Int).SetInt64(17428695),
	RegistryContractAddr: "0x42Cd8D240E30102B715d7516f97864ECeC4441Ab",
	Bootnodes: []string{
		"enr:-Li4QHEPYASj5ZY3BXXKXAoWcoIw0ChgUlTtfOSxgNlYxlmpEWUR_K6Nr04VXsMpWSQxWWM4QHDyypnl92DQNpWkMS-GAYiWUvo8h2F0dG5ldHOIAAAAAAAAAACEZXRoMpD1pf1CAAAAAP__________gmlkgnY0gmlwhCzmKVSJc2VjcDI1NmsxoQOW29na1pUAQw4jF3g0zsPgJG89ViHJOOkHFFklnC2UyIN0Y3CCE4qDdWRwgg-i",
	},
	WhitelistedOperatorKeys: []string{
		// Blox's exporter nodes.
		"LS0tLS1CRUdJTiBSU0EgUFVCTElDIEtFWS0tLS0tCk1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBeSttZUZmOGwyZ2lzK01YU1pXc3oKTFhPTjRodXJ1ak4zQnNSRTJCU1FuV3RMUHY1Y2prSW0xa2FxRThqWEtBbU5nRkwrT3o0eldpczhFRVFhMjB1UgpyM01RM1NMQnlpaWFRYjNDeStjMFg1UDFsTFBBYzVxNnhiZlJBZXp3K2dUUkFYSXo4RXBwaGdVblNyVUQvOXp2CmZ1OFRaQkVLSlYrcnFDRTZZN0FpcU9jVUsrNHF3TWUyeWQrMW9rRld2d2E3c3h4T2VZNGdBcG9jTENNQmRzKzIKQlY2UVR5aVZaT1daQlhFSjdXMllINHBHMWRlMHdMRUZaUnVkcmE1L3RXUzBqSzRRV0Vhc21WeG5LOUpsSWJDdgplZW5vYWt2M1pjamM4WGs1MmRLWGFuNy9TNDNxdHRJT1MvbFdmRDdxSTZvWXp5aXJhZVh2dDdYbXlhODJIa3JZCkd3SURBUUFCCi0tLS0tRU5EIFJTQSBQVUJMSUMgS0VZLS0tLS0K",
		"LS0tLS1CRUdJTiBSU0EgUFVCTElDIEtFWS0tLS0tCk1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBeTd0dlFKa01nZUx3ekp5c1lZUlgKRWNyTC9PUHBScDkyUWxEdHZNaXQrdXZKbEtSV1B6OFZCbERtV1UvV0UrTkUxa2FyTDVrRFAxazRMMUFUUjJFWgpLNEpMUUp3RW5kNHVBWlB3Rm1UVTEvWmVDN1kvVmlNcmlyMW1pSzRmcXNnTko5UmVWWjAzQ3hpVGNQQjNHNTE0ClhQaklzaUo0eS8wSlB6cmhQckR5Vmt3SnEwWWRnMWpJMUJkbzVaVm15SkZ4eC9lblcwcVUrNG9iaElGZThlUkEKdjUrbS9aa1lUbnNoMklsVk10UjB5TUQwR0I0YWo3MGQ1VVIwMk1yZkhCWXVLOHpnSitXVkN2R0JVTm9ramVFZQpvWVRsYmQwSzAxRWh1MHN1cStjc0FubU8vaTBaaDVHOVM3MU5EVkc4QnBhdVk5cHYvcFlDa3ZqaHNRdGtQTEJKCjd3SURBUUFCCi0tLS0tRU5EIFJTQSBQVUJMSUMgS0VZLS0tLS0K",
	},
}
