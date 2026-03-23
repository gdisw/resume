package view

var staticConfig map[string]string

func PutStaticConfig(keyValues ...string) {
	staticConfig = make(map[string]string)
	for i := 0; i < len(keyValues)-1; i += 2 {
		k := keyValues[i]
		v := keyValues[i+1]
		staticConfig[k] = v
	}
}

func getStaticConfig(k string) string {
	return staticConfig[k]
}
