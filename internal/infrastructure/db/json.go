package db

import "gorm.io/datatypes"

func toJSONMap(data map[string]any) datatypes.JSONMap {
	if data == nil {
		return datatypes.JSONMap{}
	}
	return datatypes.JSONMap(data)
}

func fromJSONMap(data datatypes.JSONMap) map[string]any {
	if data == nil {
		return nil
	}
	return map[string]any(data)
}
