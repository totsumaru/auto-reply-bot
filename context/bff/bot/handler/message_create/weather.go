package message_create

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/conf"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/convert"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/info/guild"
	"github.com/techstart35/auto-reply-bot/context/discord/expose/message_send"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
	"github.com/techstart35/auto-reply-bot/context/shared/map/seeker"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// 市の名前が含まれている場合の、天気の返信テンプレートです
// 0: title(〇〇の天気)
// -------
// 1: today-telop
// 2: today-%(00-06)
// 3: today-%(06-12)
// 4: today-%(12-18)
// 5: today-%(18-24)
// -------
// 6: tomorrow-telop
// 7: tomorrow-%(00-06)
// 8: tomorrow-%(06-12)
// 9: tomorrow-%(12-18)
// 10: tomorrow-%(18-24)
const WeatherResCityTmpl = `
%s

▬▬▬▬ きょうの天気 ▬▬▬▬

・%s

__降水確率__
・００時-０６時: %s
・０６時-１２時: %s
・１２時-１８時: %s
・１８時-２４時: %s

▬▬▬▬ あしたの天気 ▬▬▬▬

・%s

__降水確率__
・００時-０６時: %s
・０６時-１２時: %s
・１２時-１８時: %s
・１８時-２４時: %s

▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬
今日も良い一日を！
`

// 県名が含まれている場合の、天気の返信テンプレートです
// 0: title(〇〇の天気)
// -------
// 1: today-date
// 2: today-telop
// 3: today-%(00-06)
// 4: today-%(06-12)
// 5: today-%(12-18)
// 6: today-%(18-24)
// -------
// 7: tomorrow-date
// 8: tomorrow-telop
// 9: tomorrow-%(00-06)
// 10: tomorrow-%(06-12)
// 11: tomorrow-%(12-18)
// 12: tomorrow-%(18-24)
// -------
// 13: 地点のリスト
// 14: その県の代表地点(静岡)
const WeatherResPrefectureTmpl = `
%s（代表地点）

▬▬▬▬ きょうの天気 ▬▬▬▬

%s

__降水確率__
・００時-０６時: %s
・０６時-１２時: %s
・１２時-１８時: %s
・１８時-２４時: %s

▬▬▬▬ あしたの天気 ▬▬▬▬

%s

__降水確率__
・００時-０６時: %s
・０６時-１２時: %s
・１２時-１８時: %s
・１８時-２４時: %s

▬▬▬▬▬▬▬▬▬▬▬▬▬▬▬
次の地点は個別に指定できます。%s
例）@Comment-bot %sの天気

今日も良い一日を！
`

// 該当なしの返信テンプレートです
const NoContainsResTmpl = `
天気を取得できませんでした。
以下3つを含めて送信してみましょう。

・botへのメンション
・県名（〇〇県）
・"天気"というワード

例)
@Comment-bot 〇〇県の天気
`

// 天気を返します
func Weather(s *discordgo.Session, m *discordgo.MessageCreate) {
	const MustKeywordTenki = "天気"
	var MustKeywordBotMention = convert.IDToMention(os.Getenv("DISCORD_APPLICATION_ID"))

	// TEST SERVERはカウントしません
	if m.GuildID == conf.TestServerID {
		return
	}

	// Botユーザーはカウントしません
	if m.Author.Bot {
		return
	}

	// マストのキーワードが入っていない場合はここで終了します
	if !(strings.Contains(m.Content, MustKeywordTenki) &&
		strings.Contains(m.Content, MustKeywordBotMention)) {

		return
	}

	guildName, err := guild.GetGuildName(s, m.GuildID)
	if err != nil {
		message_send.SendErrMsg(s, errors.NewError("ギルド名を取得できません", err), "")
		return
	}

	// =============================
	// 1. 県名(〇〇県)が含まれているか？
	//	yes → 返信して終了
	// 	no ↓
	// 2. 市の名前が含まれているか？
	// 	yes → 返信して終了
	// 	no → 登録なしメッセージを送信
	// =============================

	for prefectureName, cityMap := range locationID {
		// 1. 県名が含まれているか確認します
		if strings.Contains(m.Content, prefectureName) {
			// 代表地点の天気を取得します
			cID := ""
			for _, priorityCityID := range priorityCity[prefectureName] {
				cID = priorityCityID
			}

			// APIから天気を取得します
			res, err := getWeatherFromAPI(cID)
			if err != nil {
				message_send.SendErrMsg(s, errors.NewError("天気のAPIを取得できません", err), guildName)
				return
			}

			// 該当の県に含まれる地点の一覧です
			cityList := make([]string, 0)
			for cName := range locationID[prefectureName] {
				cityList = append(cityList, cName)
			}

			// 代表地点名です
			priorityCityName := ""
			for p := range priorityCity[prefectureName] {
				priorityCityName = p
			}

			// テンプレートを作成します
			msg := fmt.Sprintf(
				WeatherResPrefectureTmpl,
				res.Title,
				weatherAddEmoji(res.Today.Telop),
				res.Today.T0006,
				res.Today.T0612,
				res.Today.T1218,
				res.Today.T1824,
				weatherAddEmoji(res.Tomorrow.Telop),
				res.Tomorrow.T0006,
				res.Tomorrow.T0612,
				res.Tomorrow.T1218,
				res.Tomorrow.T1824,
				cityList,
				priorityCityName,
			)

			// メッセージを送信します
			{
				req := message_send.SendReplyEmbedReq{
					ChannelID: m.ChannelID,
					Content:   msg,
					Color:     conf.ColorOrange,
					Reference: m.Reference(),
					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: res.Today.ImageURL,
					},
				}

				if err = message_send.SendReplyEmbed(s, req); err != nil {
					message_send.SendErrMsg(s, errors.NewError("県名が含まれている場合の返信を送信できません", err), guildName)
					return
				}
			}
			return
		} else {
			for cityName, cityID := range cityMap {
				// 2. 市の名前が含まれているか確認します
				if strings.Contains(m.Content, cityName) {
					// 津の場合は大津でも反応してしまうため、
					// 大津が含まれていないかを再度確認します。
					{
						if strings.Contains(m.Content, "大津") {
							cityName = "大津"
							cityID = locationID["滋賀県"]["大津"]
						}
					}

					res, err := getWeatherFromAPI(cityID)
					if err != nil {
						message_send.SendErrMsg(s, errors.NewError("天気のAPIを取得できません", err), guildName)
						return
					}

					msg := fmt.Sprintf(
						WeatherResCityTmpl,
						res.Title,
						weatherAddEmoji(res.Today.Telop),
						res.Today.T0006,
						res.Today.T0612,
						res.Today.T1218,
						res.Today.T1824,
						weatherAddEmoji(res.Tomorrow.Telop),
						res.Tomorrow.T0006,
						res.Tomorrow.T0612,
						res.Tomorrow.T1218,
						res.Tomorrow.T1824,
					)

					// メッセージを送信します
					{
						req := message_send.SendReplyEmbedReq{
							ChannelID: m.ChannelID,
							Content:   msg,
							Color:     conf.ColorOrange,
							Reference: m.Reference(),
						}
						if err = message_send.SendReplyEmbed(s, req); err != nil {
							message_send.SendErrMsg(s, errors.NewError("市の名前が含まれている場合の返信を送信できません", err), guildName)
							return
						}
					}
					return
				}
			}
		}
	}

	// 登録なしメッセージを送信します
	msg := NoContainsResTmpl

	// メッセージを送信します
	{
		req := message_send.SendReplyEmbedReq{
			ChannelID: m.ChannelID,
			Content:   msg,
			Color:     conf.ColorRed,
			Reference: m.Reference(),
		}
		if err = message_send.SendReplyEmbed(s, req); err != nil {
			message_send.SendErrMsg(s, errors.NewError("登録なしメッセージを送信できません", err), guildName)
			return
		}
	}
	return
}

type apiRes struct {
	Title    string
	Text     string
	Today    DayForecast
	Tomorrow DayForecast
}

type DayForecast struct {
	Date     string
	Telop    string
	T0006    string
	T0612    string
	T1218    string
	T1824    string
	ImageURL string
}

// 天気をAPIから取得します
//
// 取得元: https://weather.tsukumijima.net/primary_area.xml
func getWeatherFromAPI(locationID string) (apiRes, error) {
	const APITmpl = "https://weather.tsukumijima.net/api/forecast/city/%s"

	res := apiRes{}

	resp, _ := http.Get(fmt.Sprintf(APITmpl, locationID))
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, errors.NewError("レスポンスを読み込めません", err)
	}

	resMap := map[string]interface{}{}
	if err := json.Unmarshal(byteArray, &resMap); err != nil {
		return res, errors.NewError("byteをmapに変換できません", err)
	}

	todayForecastMap := seeker.Slice(resMap, []string{"forecasts"})[0]
	today := DayForecast{
		Date:     seeker.Str(todayForecastMap, []string{"date"}),
		Telop:    seeker.Str(todayForecastMap, []string{"telop"}),
		T0006:    seeker.Str(todayForecastMap, []string{"chanceOfRain", "T00_06"}),
		T0612:    seeker.Str(todayForecastMap, []string{"chanceOfRain", "T06_12"}),
		T1218:    seeker.Str(todayForecastMap, []string{"chanceOfRain", "T12_18"}),
		T1824:    seeker.Str(todayForecastMap, []string{"chanceOfRain", "T18_24"}),
		ImageURL: seeker.Str(todayForecastMap, []string{"image", "url"}),
	}

	tomorrowForecastMap := seeker.Slice(resMap, []string{"forecasts"})[1]
	tomorrow := DayForecast{
		Date:     seeker.Str(tomorrowForecastMap, []string{"date"}),
		Telop:    seeker.Str(tomorrowForecastMap, []string{"telop"}),
		T0006:    seeker.Str(tomorrowForecastMap, []string{"chanceOfRain", "T00_06"}),
		T0612:    seeker.Str(tomorrowForecastMap, []string{"chanceOfRain", "T06_12"}),
		T1218:    seeker.Str(tomorrowForecastMap, []string{"chanceOfRain", "T12_18"}),
		T1824:    seeker.Str(tomorrowForecastMap, []string{"chanceOfRain", "T18_24"}),
		ImageURL: seeker.Str(tomorrowForecastMap, []string{"image", "url"}),
	}

	res.Title = seeker.Str(resMap, []string{"title"})
	res.Text = seeker.Str(resMap, []string{"description", "text"})
	res.Today = today
	res.Tomorrow = tomorrow

	return res, nil
}

// 天気に絵文字を追加します
func weatherAddEmoji(weather string) string {
	res := make([]string, 0)
	for _, v := range strings.Split(weather, "") {
		if strings.Contains(v, "晴") {
			res = append(res, "☀️")
		}
		if strings.Contains(v, "曇") {
			res = append(res, "☁️")
		}
		if strings.Contains(v, "雨") {
			res = append(res, "☔️")
		}
		if strings.Contains(v, "雪") {
			res = append(res, "❄️")
		}
		// 「のち」を表している
		if strings.Contains(v, "の") {
			res = append(res, "→")
		}
	}

	res = append(res, "｜")
	res = append(res, weather)

	return strings.Join(res, "")
}

// 代表地点のリストです
var priorityCity = map[string]map[string]string{
	"北海道": {
		"稚内": "011000",
	},
	"青森県": {
		"青森": "020010",
	},
	"岩手県": {
		"盛岡": "030010",
	},
	"宮城県": {
		"仙台": "040010",
	},
	"秋田県": {
		"秋田": "050010",
	},
	"山形県": {
		"山形": "060010",
	},
	"福島県": {
		"福島": "070010",
	},
	"東京都": {
		"東京": "130010",
	},
	"神奈川県": {
		"横浜": "140010",
	},
	"埼玉県": {
		"さいたま": "110010",
	},
	"千葉県": {
		"千葉": "120010",
	},
	"茨城県": {
		"水戸": "080010",
	},
	"栃木県": {
		"宇都宮": "090010",
	},
	"群馬県": {
		"前橋": "100010",
	},
	"山梨県": {
		"甲府": "190010",
	},
	"新潟県": {
		"新潟": "150010",
	},
	"長野県": {
		"長野": "200010",
	},
	"富山県": {
		"富山": "160010",
	},
	"石川県": {
		"金沢": "170010",
	},
	"福井県": {
		"福井": "180010",
	},
	"愛知県": {
		"名古屋": "230010",
	},
	"岐阜県": {
		"岐阜": "210010",
	},
	"静岡県": {
		"静岡": "220010",
	},
	"三重県": {
		"津": "240010",
	},
	"大阪府": {
		"大阪": "270000",
	},
	"兵庫県": {
		"神戸": "280010",
	},
	"京都府": {
		"京都": "260010",
	},
	"滋賀県": {
		"大津": "250010",
	},
	"奈良県": {
		"奈良": "290010",
	},
	"和歌山県": {
		"和歌山": "300010",
	},
	"鳥取県": {
		"鳥取": "310010",
	},
	"島根県": {
		"松江": "320010",
	},
	"岡山県": {
		"岡山": "330010",
	},
	"広島県": {
		"広島": "340010",
	},
	"山口県": {
		"下関": "350010",
	},
	"徳島県": {
		"徳島": "360010",
	},
	"香川県": {
		"高松": "370000",
	},
	"愛媛県": {
		"松山": "380010",
	},
	"高知県": {
		"高知": "390010",
	},
	"福岡県": {
		"福岡": "400010",
	},
	"大分県": {
		"大分": "440010",
	},
	"長崎県": {
		"長崎": "420010",
	},
	"佐賀県": {
		"佐賀": "410010",
	},
	"熊本県": {
		"熊本": "430010",
	},
	"宮崎県": {
		"宮崎": "450010",
	},
	"鹿児島県": {
		"鹿児島": "460010",
	},
	"沖縄県": {
		"那覇": "471010",
	},
}

// 地域IDのリストです
var locationID = map[string]map[string]string{
	"北海道": {
		"稚内":  "011000",
		"旭川":  "012010",
		"留萌":  "012020",
		"網走":  "013010",
		"北見":  "013020",
		"紋別":  "013030",
		"根室":  "014010",
		"釧路":  "014020",
		"帯広":  "014030",
		"室蘭":  "015010",
		"浦河":  "015020",
		"札幌":  "016010",
		"岩見沢": "016020",
		"倶知安": "016030",
		"函館":  "017010",
		"江差":  "017020",
	},
	"青森県": {
		"青森": "020010",
		"むつ": "020020",
		"八戸": "020030",
	},
	"岩手県": {
		"盛岡":  "030010",
		"宮古":  "030020",
		"大船渡": "030030",
	},
	"宮城県": {
		"仙台": "040010",
		"白石": "040020",
	},
	"秋田県": {
		"秋田": "050010",
		"横手": "050020",
	},
	"山形県": {
		"山形": "060010",
		"米沢": "060020",
		"酒田": "060030",
		"新庄": "060040",
	},
	"福島県": {
		"福島":  "070010",
		"小名浜": "070020",
		"若松":  "070030",
	},
	"東京都": {
		"東京":  "130010",
		"大島":  "130020",
		"八丈島": "130030",
		"父島":  "130040",
	},
	"神奈川県": {
		"横浜":  "140010",
		"小田原": "140020",
	},
	"埼玉県": {
		"さいたま": "110010",
		"熊谷":   "110020",
		"秩父":   "110030",
	},
	"千葉県": {
		"千葉": "120010",
		"銚子": "120020",
		"館山": "120030",
	},
	"茨城県": {
		"水戸": "080010",
		"土浦": "080020",
	},
	"栃木県": {
		"宇都宮": "090010",
		"大田原": "090020",
	},
	"群馬県": {
		"前橋":   "100010",
		"みなかみ": "100020",
	},
	"山梨県": {
		"甲府":  "190010",
		"河口湖": "190020",
	},
	"新潟県": {
		"新潟": "150010",
		"長岡": "150020",
		"高田": "150030",
		"相川": "150040",
	},
	"長野県": {
		"長野": "200010",
		"松本": "200020",
		"飯田": "200030",
	},
	"富山県": {
		"富山": "160010",
		"伏木": "160020",
	},
	"石川県": {
		"金沢": "170010",
		"輪島": "170020",
	},
	"福井県": {
		"福井": "180010",
		"敦賀": "180020",
	},
	"愛知県": {
		"名古屋": "230010",
		"豊橋":  "230020",
	},
	"岐阜県": {
		"岐阜": "210010",
		"高山": "210020",
	},
	"静岡県": {
		"静岡": "220010",
		"網代": "220020",
		"三島": "220030",
		"浜松": "220040",
	},
	"三重県": {
		"津":  "240010",
		"尾鷲": "240020",
	},
	"大阪府": {
		"大阪": "270000",
	},
	"兵庫県": {
		"神戸": "280010",
		"豊岡": "280020",
	},
	"京都府": {
		"京都": "260010",
		"舞鶴": "260020",
	},
	"滋賀県": {
		"大津": "250010",
		"彦根": "250020",
	},
	"奈良県": {
		"奈良": "290010",
		"風屋": "290020",
	},
	"和歌山県": {
		"和歌山": "300010",
		"潮岬":  "300020",
	},
	"鳥取県": {
		"鳥取": "310010",
		"米子": "310020",
	},
	"島根県": {
		"松江": "320010",
		"浜田": "320020",
		"西郷": "320030",
	},
	"岡山県": {
		"岡山": "330010",
		"津山": "330020",
	},
	"広島県": {
		"広島": "340010",
		"庄原": "340020",
	},
	"山口県": {
		"下関": "350010",
		"山口": "350020",
		"柳井": "350030",
		"萩":  "350040",
	},
	"徳島県": {
		"徳島":  "360010",
		"日和佐": "360020",
	},
	"香川県": {
		"高松": "370000",
	},
	"愛媛県": {
		"松山":  "380010",
		"新居浜": "380020",
		"宇和島": "380030",
	},
	"高知県": {
		"高知":  "390010",
		"室戸岬": "390020",
		"清水":  "390030",
	},
	"福岡県": {
		"福岡":  "400010",
		"八幡":  "400020",
		"飯塚":  "400030",
		"久留米": "400040",
	},
	"大分県": {
		"大分": "440010",
		"中津": "440020",
		"日田": "440030",
		"佐伯": "440040",
	},
	"長崎県": {
		"長崎":  "420010",
		"佐世保": "420020",
		"厳原":  "420030",
		"福江":  "420040",
	},
	"佐賀県": {
		"佐賀":  "410010",
		"伊万里": "410020",
	},
	"熊本県": {
		"熊本":   "430010",
		"阿蘇乙姫": "430020",
		"牛深":   "430030",
		"人吉":   "430040",
	},
	"宮崎県": {
		"宮崎":  "450010",
		"延岡":  "450020",
		"都城":  "450030",
		"高千穂": "450040",
	},
	"鹿児島県": {
		"鹿児島": "460010",
		"鹿屋":  "460020",
		"種子島": "460030",
		"名瀬":  "460040",
	},
	"沖縄県": {
		"那覇":   "471010",
		"名護":   "471020",
		"久米島":  "471030",
		"南大東":  "472000",
		"宮古島":  "473000",
		"石垣島":  "474010",
		"与那国島": "474020",
	},
}
