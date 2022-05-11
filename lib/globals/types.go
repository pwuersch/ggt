package globals

type Profile struct {
	Name       string     `json:"name"`
	CommitName string     `json:"commitName"`
	Email      string     `json:"email"`
	Api        ApiProfile `json:"api"`
}

type ApiProfile struct {
	Adapter      string `json:"adapter"`
	Token        string `json:"token"`
	Host         string `json:"host,omitempty"`
	DefaultScope string `json:"defaultScope,omitempty"`
}

func Map[K any, U any](input []K, mappingFn func(K) U) []U {
	newArr := make([]U, len(input))
	for i, item := range input {
		newArr[i] = mappingFn(item)
	}
	return newArr
}
