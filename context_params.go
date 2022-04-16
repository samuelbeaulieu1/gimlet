package gimlet

type ContextParams struct {
	params map[string]string
}

func NewContextParams(keys []string, values []string) *ContextParams {
	params := map[string]string{}
	for index, key := range keys {
		params[key] = values[index]
	}

	return &ContextParams{
		params: params,
	}
}

func (params *ContextParams) Get(key string) string {
	if val, ok := params.params[key]; ok {
		return val
	}

	return ""
}
