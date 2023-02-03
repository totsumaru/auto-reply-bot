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
			Name:             "foo",
			Keyword:          []string{"k1", "k2"},
			Reply:            []string{"r1", "r2"},
			MatchCondition:   block.MatchConditionAllContain,
			LimitedChannelID: []string{"ch1"},
			IsRandom:         true,
			IsEmbed:          true,
		}

		req := v1.Req{}
		req.AdminRoleID = TestAdminRoleID
		// Comment
		req.Comment.BlockReq = []v1.BlockReq{blockReq1}
		req.Comment.IgnoreChannelID = []string{"ig1", "ig2"}
		// Rule
		req.Rule.URL.IsRestrict = true
		req.Rule.URL.IsYoutubeAllow = true
		req.Rule.URL.IsTwitterAllow = true
		req.Rule.URL.IsGIFAllow = true
		req.Rule.URL.IsOpenseaAllow = true
		req.Rule.URL.IsDiscordAllow = true
		req.Rule.URL.AllowRoleID = []string{"r1", "r2"}
		req.Rule.URL.AllowChannelID = []string{"c1", "c2"}

		res, err := v1.UpdateConfig(&discordgo.Session{}, ctx, TestID, req)
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
			Name:             blockReq1.Name,
			Keyword:          blockReq1.Keyword,
			Reply:            blockReq1.Reply,
			MatchCondition:   blockReq1.MatchCondition,
			LimitedChannelID: blockReq1.LimitedChannelID,
			IsRandom:         blockReq1.IsRandom,
			IsEmbed:          blockReq1.IsEmbed,
		}

		if !reflect.DeepEqual(res.Comment.Block, []v1.BlockRes{expectBlockRes}) {
			t.Fatal("期待した値と一致しません")
		}
		if res.Rule.URL.IsRestrict != req.Rule.URL.IsRestrict {
			t.Fatal("期待した値と一致しません")
		}
		if res.Rule.URL.IsYoutubeAllow != req.Rule.URL.IsYoutubeAllow {
			t.Fatal("期待した値と一致しません")
		}
		if res.Rule.URL.IsTwitterAllow != req.Rule.URL.IsTwitterAllow {
			t.Fatal("期待した値と一致しません")
		}
		if res.Rule.URL.IsGIFAllow != req.Rule.URL.IsGIFAllow {
			t.Fatal("期待した値と一致しません")
		}
		if !reflect.DeepEqual(res.Rule.URL.AllowRoleID, req.Rule.URL.AllowRoleID) {
			t.Fatal("期待した値と一致しません")
		}
		if !reflect.DeepEqual(res.Rule.URL.AllowRoleID, req.Rule.URL.AllowRoleID) {
			t.Fatal("期待した値と一致しません")
		}
		if !reflect.DeepEqual(res.Rule.URL.AllowChannelID, req.Rule.URL.AllowChannelID) {
			t.Fatal("期待した値と一致しません")
		}
	})
}
