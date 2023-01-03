package module_test

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/server/domain/model/server/block"
	"github.com/techstart35/auto-reply-bot/context/server/expose/api/v1"
	"github.com/techstart35/auto-reply-bot/context/shared/map/gen"
	"reflect"
	"testing"
)

func TestUpdateConfig(t *testing.T) {
	t.Run("設定を更新できる", func(t *testing.T) {
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
		blockReq1 := v1.BlockReq{
			Name:           "foo",
			Keyword:        []string{"k1", "k2"},
			Reply:          []string{"r1", "r2"},
			MatchCondition: block.MatchConditionAllContain,
			IsRandom:       true,
			IsEmbed:        true,
		}

		res, err := v1.UpdateConfig(
			&discordgo.Session{},
			ctx,
			TestID,
			TestAdminRoleID,
			[]v1.BlockReq{blockReq1},
		)
		if err != nil {
			t.Fatal(err)
		}

		if res.ID != TestID {
			t.Fatal("期待した値と一致しません")
		}

		if res.AdminRoleID != TestAdminRoleID {
			t.Fatal("期待した値と一致しません")
		}

		expectBlockRes := v1.BlockRes{
			Name:           "foo",
			Keyword:        []string{"k1", "k2"},
			Reply:          []string{"r1", "r2"},
			MatchCondition: block.MatchConditionAllContain,
			IsRandom:       true,
			IsEmbed:        true,
		}

		if !reflect.DeepEqual(res.Block, []v1.BlockRes{expectBlockRes}) {
			t.Fatal("期待した値と一致しません")
		}
	})

	t.Run("IsAllMatch=trueが入っている状態で、MatchConditionに変更できる", func(t *testing.T) {
		ctx, teardown := setup(t)
		defer teardown()

		// 外部APIをモックします
		MockExternal(
			TestNow,
		)

		// モックデータを登録します
		{
			blockMock := map[string]interface{}{}
			gen.Gen(blockMock, []string{"keyword"}, []map[string]interface{}{
				{"value": "k1"},
				{"value": "k2"},
			})
			gen.Gen(blockMock, []string{"reply"}, []map[string]interface{}{
				{"value": "r1"},
				{"value": "r2"},
			})
			gen.Gen(blockMock, []string{"is_all_match"}, true)
			gen.Gen(blockMock, []string{"is_random"}, true)
			gen.Gen(blockMock, []string{"is_embed"}, true)

			mock := ServerInitialMock(TestID)
			gen.Gen(mock, []string{"admin_role_id", "value"}, TestAdminRoleID)
			gen.Gen(mock, []string{"block"}, []map[string]interface{}{blockMock})
			RegisterServer(t, mock)
		}

		// テスト対象のAPIをコールします
		// モックデータと同じ値を設定します
		blockReq1 := v1.BlockReq{
			Name:           "foo",
			Keyword:        []string{"k1", "k2"},
			Reply:          []string{"r1", "r2"},
			MatchCondition: block.MatchConditionAllContain,
			IsRandom:       true,
			IsEmbed:        true,
		}

		res, err := v1.UpdateConfig(
			&discordgo.Session{},
			ctx,
			TestID,
			TestAdminRoleID,
			[]v1.BlockReq{blockReq1},
		)
		if err != nil {
			t.Fatal(err)
		}

		if res.Block[0].MatchCondition != block.MatchConditionAllContain {
			t.Fatal("期待した値と一致しません")
		}
	})
}
