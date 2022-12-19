// map[string]interface{}型を生成する関数を提供します。
package gen

import (
	"encoding/json"
	"fmt"
)

// rpで指定したキーで再帰的にmap[string]interface生成してそれをdに挿入します
func Gen(d map[string]interface{}, rp []string, value interface{}) map[string]interface{} {
	if len(rp) == 0 {
		return map[string]interface{}{}
	}

	cp := rp[0]
	rp = rp[1:]

	if len(rp) == 0 {
		d[cp] = value

		return d
	}
	if len(rp) > 0 {
		if d[cp] == nil {
			d[cp] = Gen(map[string]interface{}{}, rp, value)

			return d
		}

		d[cp] = Gen((d[cp]).(map[string]interface{}), rp, value)

		return d
	}

	return map[string]interface{}{}
}

// 構造体をmap[string]interfaceに変換します
//
// map[string]interfaceのkeyには自動的にフィールド名が指定されます。
func Conv(i interface{}) (map[string]interface{}, error) {
	b, err := json.Marshal(i)
	if err != nil {
		return nil, fmt.Errorf("構造体をJSONに変換できません %w", err)
	}
	var data map[string]interface{}
	if err = json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("JSONをmapに変換できません %w", err)
	}

	return data, nil
}
