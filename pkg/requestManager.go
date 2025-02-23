package pkg

type RequestManager interface {
	Do(url string, headers map[string]string) ([]byte, error)
}
