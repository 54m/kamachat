package pattern

import (
	"math/rand"
	"strings"
	"time"

	"github.com/mattn/go-gimei"
)

// 文章中一種類に統一されるタグ
var uniqTags = map[string][]string{
	// 対象の名前
	"{TARGET_NAME}": {},
	// おかまさんの一人称
	"{FIRST_PERSON}": {
		"あたし",
		"アタシ",
		"ｱﾀｼ",
		"あたい",
		"アタイ",
		"ｱﾀｲ",
		"うち",
		"ウチ",
		"ｳﾁ",
		"おねえさん",
		"オネエサン",
		"ｵﾈｴｻｧﾝ",
	},
	// 曜日
	"{DAY_OF_WEEK}": {
		"月", "火", "水", "木", "金", "土", "日",
	},
	// 地名
	"{LOCATION}": {
		"愛知", "青森", "秋田", "石川", "茨城", "岩手", "愛媛", "大分", "大阪", "岡山", "沖縄", "香川", "鹿児島", "神奈川", "岐阜", "京都", "熊本", "群馬", "高知", "埼玉", "佐賀", "滋賀", "静岡", "島根", "千葉", "東京", "徳島", "栃木", "鳥取", "富山", "長崎", "長野", "奈良", "新潟", "兵庫", "広島", "福井", "福岡", "福島", "北海道", "三重", "宮城", "宮崎", "山形", "山口", "山梨", "和歌山",
	},
	// 外食
	"{RESTAURANT}": {
		"お寿司🍣",
		"イタリアン🍝",
		"バー🍷",
		"ラーメン屋さん🍜",
		"中華🍜",
		"ゲイバー♂",
		"おかまバー♂",
		"ホスト😎✌",
	},
	// 食べ物
	"{FOOD}": {
		"お惣菜",
		"サラダ",
		"おにぎり🍙",
		"きんぴらごぼう",
		"ピッツァ🍕",
		"パスタ🍝",
		"スイーツ🍮",
		"ケーキ🎂",
		"そそりたつ棒", // 偏見が入ってます
		"硬い棒",    // 偏見が入ってます
	},
	// 天気
	"{WEATHER}": {
		"曇り",
		"晴れ",
		"快晴",
		"大雨",
		"雨",
		"雪",
		"台風🌀",
	},
	// 下ネタの後は「ナンチャッテ」「冗談（笑）」を使う
	"{NANCHATTE}": {
		"ﾅﾝﾁｬｯﾃ{EMOJI_POS}",
		"ナンチャッテ{EMOJI_POS}",
		"なんちゃって{EMOJI_POS}",
		"な〜んてね{EMOJI_POS}",
		"冗談冗談{EMOJI_POS}",
		"", // おかまは心は女だがたまに本気でオヤジ
	},
	// TODO: 「かな？」と「かい？」の語尾の違いも組み込んでも良いかもしれない
	// おかまの欲望の地、ホテル
	"{HOTEL}": {
		"ホテル🏨",
		"ホテル🏩",
		"旅館",
		"密室",
	},
	// デートの種類
	"{DATE}": {
		"デート❤",
		"カラオケ🎤",
		"ドライブ🚗",
		"オウチ", // 偏見が入ってます
	},
	// おじさんは比喩で相手を持ち上げる (川柳)
	"{METAPHOR}": {
		"イケメン",
		"かっこいい",
		"男らしい",
		"優男",
	},
}

// 文章中複数回変更&繰り返されるタグ
// 心は乙女だが、奥底にはオヤジが眠っているため、表現する絵文字/顔文字はオヤジに寄っている(偏見が入ってます)
var flexTags = map[string][]string{
	// ポジティブな表現の絵文字/顔文字
	"{EMOJI_POS}": {
		"😃♥ ",
		"😃☀ ",
		"😃",
		"😃✋",
		"❗",
		"😄",
		"😆",
		"😚",
		"😘",
		"😚",
		"💕",
		"💗",
		"😍",
		"😁",
		"😋",
		"😂",
		"😊",
		"🎵",
		"(^_^)",
		"(^o^)",
		"(^з<)",
		"（笑）",
	},
	// ネガティヴな表現の絵文字/顔文字
	"{EMOJI_NEG}": {
		"💦",
		"💔",
		"😱",
		"😰",
		"😭",
		"😓",
		"😣",
		"😖",
		"😥",
		"😢",
		"(◎ ＿◎;)",
		"(T_T)",
		"^^;",
		"(^_^;",
		"(・_・;",
		"(￣Д￣；；",
		"(^▽^;)",
		"(-_-;)",
	},
	// ニュートラルな感情を表す絵文字/顔文字
	"{EMOJI_NEUT}": {
		"💤",
		"😴",
		"🙂",
		"🤑",
		"✋",
		"😪",
		"🛌",
		"😎",
		"😤",
		"😒",
		"😙",
		"😏",
		"😳",
		"😌",
		"（￣▽￣）",
		"(＃￣З￣)",
		"(^^;;",
	},
	// 疑問を投げかけるときに利用される絵文字/顔文字
	"{EMOJI_ASK}": {
		"⁉",
		"❓",
		"❗❓",
		"🤔",
		"😜⁉️",
		"✋❓",
		"（￣ー￣?）",
	},
}

// ConvertTags ; message内にあるタグを置換して結果を返す
func ConvertTags(message, targetName string, emojiNumber int) string {
	rand.Seed(time.Now().UnixNano())
	if targetName != "" {
		// 引数として名前が存在した場合にはそれを使う
		uniqTags["{TARGET_NAME}"] = []string{targetName + randomNameSuffix()}
	} else {
		// 指定がない場合には gimei から選定
		uniqTags["{TARGET_NAME}"] = []string{randomFirstName() + randomNameSuffix()}
	}

	for tag, pat := range uniqTags {
		content := pat[rand.Intn(len(pat))]
		message = strings.Replace(message, tag, content, -1)
	}

	for tag, pat := range flexTags {
		n := strings.Count(message, tag)
		for i := 0; i < n; i++ {
			content := ""
			if emojiNumber > 0 {
				content = combineMultiplePatterns(pat, rand.Intn(emojiNumber)+1)
			} else {
				// Okama could be seriously
				content = "。"
			}
			// タグを置換
			message = strings.Replace(message, tag, content, 1)
		}
	}
	return message
}

// combineMultiplePatterns: 複数のパターンをnumber分ランダムにつなげる
func combineMultiplePatterns(patterns []string, number int) string {
	result := ""
	if number <= len(patterns) {
		for i := 0; i < number; i++ {
			index := rand.Intn(len(patterns) - i)
			result += patterns[index]
			patterns[index], patterns[len(patterns)-1-i] = patterns[len(patterns)-1-i], patterns[index]
		}
	} else {
		for i := 0; i < number; i++ {
			result += patterns[rand.Intn(len(patterns))]
		}
	}
	return result
}

// gimei から男性の名前を無作為に選定
func randomFirstName() string {
	name := gimei.NewMale()
	rand.Seed(time.Now().UnixNano())
	switch rand.Intn(3) {
	case 0:
		return name.First.Kanji()
	case 1:
		return name.First.Katakana()
	}
	return name.First.Hiragana()
}

// 「くん」「クン」「ちゃん」「チャン」などをランダムに返す
func randomNameSuffix() string {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(100)
	switch {
	// たまに呼び捨てにするおかま
	case n < 5:
		return ""
	// 男らしいと感じた時に「くん」呼びし始めるおかま 作者予想(偏見が入ってます)
	case n < 10:
		return "くん"
	// 時に「クン」とカタカナにしてくる
	case n < 15:
		return "クン"
	// 「〇〇クン」をさらに半角で表現する、そんなおかまもいる 作者予想(偏見が入ってます)
	case n < 20:
		return "ﾁｬﾝ"
	case n < 50:
		return "サァン"
	case n < 55:
		return "ｻｧﾝ"
	// かま友など、自分と同じ匂いを感じ取った相手(つまり仲間意識)によく使う 作者予想(偏見が入ってます)
	case n < 70:
		return "っち"
	// 多くの場合「ちゃん」にする
	default:
		return "ちゃん"
	}
}
