package module_test

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/server/expose/api/v1"
	"testing"
)

func TestDeleteServer(t *testing.T) {
	t.Run("サーバーを削除できる", func(t *testing.T) {
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
		if err := v1.DeleteServer(&discordgo.Session{}, ctx, TestID); err != nil {
			t.Fatal(err)
		}

		// 削除されているか確認します
		// エラーが返されればOKです
		_, err := v1.FindByID(ctx, TestID)
		if err == nil {
			t.Fatal("エラーが返されません")
		}
	})

	t.Run("登録されていない場合はエラーが返される", func(t *testing.T) {
		ctx, teardown := setup(t)
		defer teardown()

		// 外部APIをモックします
		MockExternal(
			TestNow,
		)

		// テスト対象のAPIをコールします
		if err := v1.DeleteServer(&discordgo.Session{}, ctx, TestID); err != nil {
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
		if err := v1.DeleteServer(&discordgo.Session{}, ctx, ""); err == nil {
			t.Fatal("エラーが返されなかった")
		}
	})
}
