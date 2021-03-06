package entity

import "errors"

var (
	ErrAreaNotExist        = errors.New("Area doesn't exist")
	ErrObservatoryNotExist = errors.New("Observatory doesn't exist")
)

type Area struct {
	ID   int64
	Name string
}

func ListAreas() []Area {
	return []Area{
		{3, "関東地域"},
		// {5, "中部地域"},
		// {6, "関西地域"},
		// {7, "中国・四国地域"},
		// {8, "九州地域"},
	}
}

func GetArea(id int64) (Area, error) {
	switch id {
	case 3:
		return Area{3, "関東地域"}, nil
	default:
		return Area{}, ErrAreaNotExist
	}
}

type Observatory struct {
	ID         int64
	Prefecture string
	Name       string
}

func (e Area) ListObservatories() ([]Observatory, error) {
	switch e.ID {
	// 関東地域
	case 3:
		return []Observatory{
			{1, "茨城県", "水戸石川一般環境大気測定局"},
			{2, "茨城県", "日立市消防本部"},
			{3, "茨城県", "国立研究開発法人 国立環境研究所"},
			{4, "栃木県", "宇都宮市中央生涯学習センター"},
			{5, "栃木県", "栃木県那須庁舎"},
			{6, "栃木県", "日光市役所第４庁舎"},
			{7, "群馬県", "群馬県衛生環境研究所"},
			{8, "群馬県", "館林保健福祉事務所"},
			{9, "埼玉県", "さいたま市役所"},
			{10, "埼玉県", "熊谷市保健センター"},
			{11, "埼玉県", "飯能市役所"},
			{12, "千葉県", "東邦大学"},
			{13, "千葉県", "千葉県環境研究センター"},
			{14, "千葉県", "印旛健康福祉センター成田支所"},
			{15, "千葉県", "君津市糠田測定局"},
			{16, "東京都", "東京都多摩小平保健所"},
			{17, "東京都", "新宿区役所第二分庁舎"},
			{18, "神奈川県", "神奈川県庁第二分庁舎"},
			{19, "神奈川県", "川崎生命科学・環境研究センター"},
			{20, "神奈川県", "神奈川県環境科学センター"},
		}, nil
	default:
		// 他地域はTODO
		return []Observatory{}, ErrAreaNotExist
	}
}

func GetObservatory(area Area, observatoryID int64) (Observatory, error) {
	observatories, err := area.ListObservatories()
	if err != nil {
		return Observatory{}, nil
	}
	for _, v := range observatories {
		if v.ID == observatoryID {
			return v, nil
		}
	}
	return Observatory{}, ErrObservatoryNotExist
}
