package v1

import (
	"fmt"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
	"sync"
)

/*
* DBにアクセスが集中しないように、全データを以下のオンメモリで管理します。
* DBが更新されたタイミングで、storeも更新します。
* BFFからのFind関数は、基本的にstoreから使用します。
 */

// ストアです
//
// サーバーID:レスポンス です。
var store = map[string]Res{}

// ロックです
var mu sync.Mutex

// ストアの状態をDBと同期します
//
// init関数からコールされます。
func InitStore(res []Res) {
	// MutexのLockをかけます
	mu.Lock()
	defer mu.Unlock()

	for _, r := range res {
		store[r.ID] = r
	}
}

// ストアの値を取得します
func GetStoreRes(id string) (Res, error) {
	if _, ok := store[id]; !ok {
		return Res{}, errors.NewError(fmt.Sprintf("storeに値が存在しません[id: %s]", id))
	}

	return store[id], nil
}

// ストアを更新します
func updateStore(res Res) error {
	// MutexのLockをかけます
	mu.Lock()
	defer mu.Unlock()

	store[res.ID] = res

	return nil
}

// ストアの値を削除します
func removeStore(id string) error {
	// MutexのLockをかけます
	mu.Lock()
	defer mu.Unlock()

	if _, ok := store[id]; !ok {
		return errors.NewError(fmt.Sprintf("storeに値が存在しません[id: %s]", id))
	}

	delete(store, id)

	return nil
}
