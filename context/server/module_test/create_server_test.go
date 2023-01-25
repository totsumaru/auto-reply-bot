package module_test

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/server/expose/api/v1"
	"reflect"
	"testing"
)

func TestCreateServer(t *testing.T) {
	t.Run("サーバーを登録できる", func(t *testing.T) {
		ctx, teardown := setup(t)
		defer teardown()

		// 外部APIをモックします
		MockExternal(
			TestNow,
		)

		// テスト対象のAPIをコールします
		res, err := v1.CreateServer(&discordgo.Session{}, ctx, TestID)
		if err != nil {
			t.Fatal(err)
		}

		if res.ID != TestID {
			t.Fatal("期待した値と一致しません")
		}

		if res.AdminRoleID != "" {
			t.Fatal("期待した値と一致しません")
		}

		if !reflect.DeepEqual(res.Comment.Block, []v1.BlockRes{}) {
			t.Fatal("期待した値と一致しません")
		}
	})

	t.Run("すでに登録されている場合はエラーが返される", func(t *testing.T) {
		ctx, teardown := setup(t)
		defer teardown()

		// 外部APIをモックします
		MockExternal(
			TestNow,
		)

		// モックデータを登録します
		{
			mock := ServerInitialMock(TestID)
			RegisterServer(t, mock)
		}

		// テスト対象のAPIをコールします
		_, err := v1.CreateServer(&discordgo.Session{}, ctx, TestID)
		if err == nil {
			t.Fatal("エラーが返されなかった")
		}
	})

	t.Run("リクエストが不正な場合はエラーが返される", func(t *testing.T) {
		ctx, teardown := setup(t)
		defer teardown()

		// 外部APIをモックします
		MockExternal(
			TestNow,
		)

		// テスト対象のAPIをコールします
		_, err := v1.CreateServer(&discordgo.Session{}, ctx, "")
		if err == nil {
			t.Fatal("エラーが返されなかった")
		}
	})
}
