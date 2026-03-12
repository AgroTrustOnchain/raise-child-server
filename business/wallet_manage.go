package business

import "sync"

var (
	_addresses map[string]securityInfo
	_wallets   map[string]string
	_shareMu   sync.RWMutex
)

type securityInfo struct {
	nonce      string
	exp        int64
	isLoggedIn bool
}

// Auto initialize map on package load
func init() {
	_addresses = make(map[string]securityInfo)
	_wallets = make(map[string]string)
}

func GetWallets() map[string]string {
	return _wallets
}

func setLogin(sub, addr string) {
	_shareMu.Lock()
	defer _shareMu.Unlock()
	_wallets[sub] = addr // mapping profileID -> address
}

func logoutWallet(sub string) {
	_shareMu.Lock()
	defer _shareMu.Unlock()
	delete(_wallets, sub)
}

// Nonce helpers
func setAddress(addr, nonce string, exp int64) {
	_shareMu.Lock()
	defer _shareMu.Unlock()
	_addresses[addr] = securityInfo{nonce: nonce, exp: exp, isLoggedIn: false}
}

func setNonce(addr, nonce string) {
	_shareMu.Lock()
	defer _shareMu.Unlock()
	var info securityInfo = _addresses[addr]
	info.nonce = nonce
	_addresses[addr] = info
}

func getSecurityInfo(addr string) securityInfo {
	_shareMu.RLock()
	defer _shareMu.RUnlock()
	var res securityInfo = _addresses[addr]
	return res
}

// Login helpers
func setLoggedIn(addr string, exp int64) {
	_shareMu.Lock()
	defer _shareMu.Unlock()
	var info securityInfo = _addresses[addr]
	info.exp = exp
	info.isLoggedIn = true
	_addresses[addr] = info
}

func isLoggedIn(addr string) bool {
	_shareMu.RLock()
	defer _shareMu.RUnlock()
	_, isLoggin := _addresses[addr]
	return isLoggin
}

func removeLoggedIn(addr string) {
	_shareMu.Lock()
	defer _shareMu.Unlock()
	delete(_addresses, addr)
}
