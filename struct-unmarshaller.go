package wstf

import "encoding/json"

// The StructUnmarshaller struct is used to cache the json bytes and
// unmarshal it as a general struct in runtime: map[string]interface{}.
// Also, some helper functions, like GetString(), GetInt(), GetBoolean(), and etc,
// are provided to investigate target json struct: like request.Body or request.Query, step by step.
type StructUnmarshaller struct {
	RawMessage json.RawMessage `json:"rawData"`

	DataMap map[string]interface{} `json:"dataMap"`
}

// GetString is most frequently used syntactic sugar.
func (m *StructUnmarshaller) GetString(key string) (string, bool) {
	value, ok := m.DataMap[key].(string)
	return value, ok
}

func (m *StructUnmarshaller) GetInt(key string) (int, bool) {
	value, ok := m.DataMap[key].(int)
	return value, ok
}

func (m *StructUnmarshaller) GetInt64(key string) (int64, bool) {
	value, ok := m.DataMap[key].(int64)
	return value, ok
}

func (m *StructUnmarshaller) GetBoolean(key string) (bool, bool) {
	value, ok := m.DataMap[key].(bool)
	return value, ok
}

//func (m *StructUnmarshaller) GetFloat32(key string) (float64, bool) {
//	value, ok := m.DataMap[key].(float32)
//	return value, ok
//}

func (m *StructUnmarshaller) GetFloat64(key string) (float64, bool) {
	value, ok := m.DataMap[key].(float64)
	return value, ok
}
