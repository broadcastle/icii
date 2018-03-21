package icii

import "testing"

func TestConnect(t *testing.T) {

	urls := map[string]bool{
		"http://192.168.1.252:80":    true,
		"https://broadcastle.co:443": true,
		"192.168.1.252:80":           false,
		"broadcastle.co":             false,
	}

	for u, result := range urls {

		_, err := connect(u)
		if err != nil && result {
			t.Error(err)
		}

	}

}
