package module_test

import (
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/server/domain/model/server/comment/block"
	"github.com/techstart35/auto-reply-bot/context/server/expose/api/v1"
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

		urlRuleReq := v1.URLRuleReq{
			IsRestrict:     true,
			IsYoutubeAllow: true,
			IsTwitterAllow: true,
			IsGIFAllow:     true,
			AllowRoleID:    []string{"r1", "r2"},
			AllowChannelID: []string{"c1", "c2"},
		}

		res, err := v1.UpdateConfig(
			&discordgo.Session{},
			ctx,
			TestID,
			TestAdminRoleID,
			[]v1.BlockReq{blockReq1},
			urlRuleReq,
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
			Name:           blockReq1.Name,
			Keyword:        blockReq1.Keyword,
			Reply:          blockReq1.Reply,
			MatchCondition: blockReq1.MatchCondition,
			IsRandom:       blockReq1.IsRandom,
			IsEmbed:        blockReq1.IsEmbed,
		}

		if !reflect.DeepEqual(res.Comment.Block, []v1.BlockRes{expectBlockRes}) {
			t.Fatal("期待した値と一致しません")
		}
		if res.Rule.URL.IsRestrict != urlRuleReq.IsRestrict {
			t.Fatal("期待した値と一致しません")
		}
		if res.Rule.URL.IsYoutubeAllow != urlRuleReq.IsYoutubeAllow {
			t.Fatal("期待した値と一致しません")
		}
		if res.Rule.URL.IsTwitterAllow != urlRuleReq.IsTwitterAllow {
			t.Fatal("期待した値と一致しません")
		}
		if res.Rule.URL.IsGIFAllow != urlRuleReq.IsGIFAllow {
			t.Fatal("期待した値と一致しません")
		}
		if !reflect.DeepEqual(res.Rule.URL.AllowRoleID, urlRuleReq.AllowRoleID) {
			t.Fatal("期待した値と一致しません")
		}
		if !reflect.DeepEqual(res.Rule.URL.AllowRoleID, urlRuleReq.AllowRoleID) {
			t.Fatal("期待した値と一致しません")
		}
		if !reflect.DeepEqual(res.Rule.URL.AllowChannelID, urlRuleReq.AllowChannelID) {
			t.Fatal("期待した値と一致しません")
		}
	})
}
