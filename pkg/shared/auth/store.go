package auth

import "sync"

// 간단한 인메모리 RTR 저장소 (userID -> refresh JTI)
var refreshStore = struct {
	sync.RWMutex
	userToJTI map[string]string
}{userToJTI: make(map[string]string)}

// SetRefreshJTI 사용자 리프레시 토큰 JTI 저장/갱신
func SetRefreshJTI(userID, jti string) {
	refreshStore.Lock()
	refreshStore.userToJTI[userID] = jti
	refreshStore.Unlock()
}

// GetRefreshJTI 사용자 리프레시 토큰 현재 JTI 조회
func GetRefreshJTI(userID string) (string, bool) {
	refreshStore.RLock()
	jti, ok := refreshStore.userToJTI[userID]
	refreshStore.RUnlock()
	return jti, ok
}
