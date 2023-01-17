package module_test

//
//import (
//	"github.com/techstart35/auto-reply-bot/context/server/domain/model/server/block"
//	"github.com/techstart35/auto-reply-bot/context/server/expose/api/v1"
//	"github.com/techstart35/auto-reply-bot/context/shared/map/gen"
//	"reflect"
//	"testing"
//)
//
//func TestFindByID(t *testing.T) {
//	t.Run("IDでサーバーを取得できる", func(t *testing.T) {
//		ctx, teardown := setup(t)
//		defer teardown()
//
//		// 外部APIをモックします
//		MockExternal(
//			TestNow,
//		)
//
//		// モックデータを登録します
//		{
//			blockMock := map[string]interface{}{}
//			gen.Gen(blockMock, []string{"keyword"}, []map[string]interface{}{
//				{"value": "k1"},
//				{"value": "k2"},
//			})
//			gen.Gen(blockMock, []string{"reply"}, []map[string]interface{}{
//				{"value": "r1"},
//				{"value": "r2"},
//			})
//			gen.Gen(blockMock, []string{"match_condition", "value"}, block.MatchConditionOneContain)
//			gen.Gen(blockMock, []string{"is_random"}, true)
//			gen.Gen(blockMock, []string{"is_embed"}, true)
//
//			mock := ServerInitialMock(TestID)
//			gen.Gen(mock, []string{"admin_role_id", "value"}, TestAdminRoleID)
//			gen.Gen(mock, []string{"block"}, []map[string]interface{}{blockMock})
//			RegisterServer(t, mock)
//		}
//
//		// テスト対象のAPIをコールします
//		res, err := v1.FindByID(ctx, TestID)
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		if res.ID != TestID {
//			t.Fatal("期待した値と一致しません")
//		}
//
//		if res.AdminRoleID != TestAdminRoleID {
//			t.Fatal("期待した値と一致しません")
//		}
//
//		if len(res.Block) != 1 {
//			t.Fatal("期待した値と一致しません")
//		}
//
//		blockRes := res.Block[0]
//
//		if !reflect.DeepEqual(blockRes.Keyword, []string{"k1", "k2"}) {
//			t.Fatal("期待した値と一致しません")
//		}
//
//		if !reflect.DeepEqual(blockRes.Reply, []string{"r1", "r2"}) {
//			t.Fatal("期待した値と一致しません")
//		}
//
//		if blockRes.MatchCondition != block.MatchConditionOneContain {
//			t.Fatal("期待した値と一致しません")
//		}
//
//		if !blockRes.IsRandom {
//			t.Fatal("期待した値と一致しません")
//		}
//
//		if !blockRes.IsEmbed {
//			t.Fatal("期待した値と一致しません")
//		}
//	})
//
//	t.Run("リクエストが不正な場合はエラーが返される", func(t *testing.T) {
//		ctx, teardown := setup(t)
//		defer teardown()
//
//		// 外部APIをモックします
//		MockExternal(
//			TestNow,
//		)
//
//		// テスト対象のAPIをコールします
//		_, err := v1.FindByID(ctx, "")
//		if err == nil {
//			t.Fatal("エラーが返されなかった")
//		}
//	})
//}
