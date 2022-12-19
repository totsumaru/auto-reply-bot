package module_test

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/shared/map/gen"
	"github.com/techstart35/auto-reply-bot/context/user/expose/api/v1"
	"testing"
)

func TestCreateUser(t *testing.T) {
	t.Run("ユーザーを登録できる", func(t *testing.T) {
		ctx, teardown := setup(t)
		defer teardown()

		// 外部APIをモックします
		MockExternal(
			TestNow,
		)

		// テスト対象のAPIをコールします
		res, err := v1.CreateUser(&discordgo.Session{}, ctx, TestID)
		if err != nil {
			t.Fatal(err)
		}

		if res.ID != TestID {
			t.Fatal("期待した値と一致しません")
		}

		if res.Name != "" {
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
			mock := UserInitialMock()
			gen.Gen(mock, []string{"id", "value"}, TestID)
			RegisterUser(t, mock)
		}

		// テスト対象のAPIをコールします
		_, err := v1.CreateUser(&discordgo.Session{}, ctx, TestID)
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
		_, err := v1.CreateUser(&discordgo.Session{}, ctx, "")
		if err == nil {
			t.Fatal("エラーが返されなかった")
		}
	})
}
