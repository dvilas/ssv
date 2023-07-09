package networkconfig

import (
	spectypes "github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv/protocol/v2/blockchain/beacon"
	"math/big"
)

var LocalNetwork = NetworkConfig{
	Name:                 "local-network",
	Beacon:               beacon.NewNetworkWithLocalTestNet(spectypes.PraterNetwork, true),
	Domain:               spectypes.DomainType{0x0, 0x0, 0x4, 0x1},
	GenesisEpoch:         1,
	ETH1SyncOffset:       new(big.Int).SetInt64(9203578),
	RegistryContractAddr: "0xC3CD9A0aE89Fff83b71b58b6512D43F8a41f363D",
	Bootnodes: []string{
		"enr:-Li4QLR4Y1VbwiqFYKy6m-WFHRNDjhMDZ_qJwIABu2PY9BHjIYwCKpTvvkVmZhu43Q6zVA29sEUhtz10rQjDJkK3Hd-GAYiGrW2Bh2F0dG5ldHOIAAAAAAAAAACEZXRoMpD1pf1CAAAAAP__________gmlkgnY0gmlwhCLdu_SJc2VjcDI1NmsxoQJTcI7GHPw-ZqIflPZYYDK_guurp_gsAFF5Erns3-PAvIN0Y3CCE4mDdWRwgg-h",
	},
	WhitelistedOperatorKeys: []string{
		// Blox's exporter nodes.
		"LS0tLS1CRUdJTiBSU0EgUFVCTElDIEtFWS0tLS0tCk1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBNmkwelNHRzFiaHlPZU8xVDVxc2UKOFpHbElBQ2pmemVYQzhpYVVReGVCb0dlVGRvN0tqalkwNy80b3hBNkhjdG45bEtxd1BodG5ISXIvZ1RlWXNYUwp5QVhPL1Q5K2RQcng1ZEp3SEVCdm5BcmNSQkNzaGF5Sng2S0xiZ3RJb2dGSWhkK1ptaFpiWFpWZVp5THhzK2tZCnM4djVwcHBIbWNwWHRwUVAxWm1ycndpTC9hZU5JNzczbUlrZ1pBOGdNK2Z5S2RtTGJrQXdXZWh1SXZKRmpuVCsKQlVkUHUzWGJIemU2SlJnY2NYNmZnM1gwOTJibG9VMzRxY1VIelNhWU9TZlc2TUpEbFgzQzJCeFhCZ042VFV0aQpDN2k2ZE9qaW14RzlSMkp4ZHVhZGpUeEM1MHl5OE9IVWpMVGNkc2pWRjdYNXdGUzFqaDI5aFpDY0FoeDB2NDg3CjdRSURBUUFCCi0tLS0tRU5EIFJTQSBQVUJMSUMgS0VZLS0tLS0K",
		"LS0tLS1CRUdJTiBSU0EgUFVCTElDIEtFWS0tLS0tCk1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBNldITnNBdTdSYnMxM0I2c0taWXgKVnZuMldlTy9YMTdSeUx1MjA0K2VtbjkvSGhIRlhXT29CMGczekNZQWp2WWdsbFJka0laTWt3ZkFUNGZvVjVTKwpvNzFFQ1dFN1ZuaytxcWd0U3k5M0ZTTVJzUG9vNngrTUd4ZURBQ3RQbDdQV1EyTXJmV1hkNzVwV1p5TVd5VndHCktPbFo0RHhoQ0VOcXlRcndlOTkybU9wVDZBcTJ1TmVsUmdESUJDSW1CV01NcUl2aXdhSU96MlBmTWR1L3ZVTWgKcVFuNGJJZjFpcVk2WGlKU1g2bDJvUWlTb09VMjRvNkFCdHlHbzRpTDJXN2tOajVUa1hOOEVzeGc3WmUveVQ0YgpKNGtvVjdmNUE3dmpMbHc1ZkdjWDR1bTBNK1QwbnczUlVIY3pHK1E3U1VGMTFGU3c0VnM1WVBHWC84a2tzdXgyCkx3SURBUUFCCi0tLS0tRU5EIFJTQSBQVUJMSUMgS0VZLS0tLS0K",
	},
}