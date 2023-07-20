package global

var Options = make(map[string]int)

const (
	ReadDeadline OptionKey = "read_deadline"
)

type OptionKey string

func (key OptionKey) Get() int {
	return Options[string(key)]
}

func (key OptionKey) GetDefault(def int) int {
	value, ok := Options[string(key)]
	if !ok {
		return def
	}
	return value
}
