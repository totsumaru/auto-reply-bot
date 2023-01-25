package module_test

import (
	"github.com/techstart35/auto-reply-bot/context/server/expose/api/v1"
	"testing"
)

func TestFindAll(t *testing.T) {
	t.Run("登録されている全てのサーバーを取得できる", func(t *testing.T) {
		ctx, teardown := setup(t)
		defer teardown()

		// 外部APIをモックします
		MockExternal(
			TestNow,
		)

		// モックデータを登録します
		{
			mock := ServerInitialMock("foo")
			RegisterServer(t, mock)

			mock2 := ServerInitialMock("bar")
			RegisterServer(t, mock2)
		}

		// テスト対象のAPIをコールします
		res, err := v1.FindAll(ctx)
		if err != nil {
			t.Fatal(err)
		}

		if len(res) != 2 {
			t.Fatal("期待した値と一致しません")
		}

		for _, r := range res {
			if !(r.ID == "foo" || r.ID == "bar") {
				t.Fatal("期待した値と一致しません")
			}
		}
	})
}
