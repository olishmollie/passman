package passman

import (
	"os"
)

type settings struct {
	root    string
	keyfile string
	mockKey []byte
}

var test = settings{"mock_store", "mock_store/.key", generateEncryptionKey()}

type mock struct {
	prefix string
	pswd   string
}

var mocks = []mock{
	{"category/website/username1", "secret1"},
	{"category/website/username2", "secret2"},
	{"username3", "secret3"},
	{"website/username4", "secret4"},
}

func setup() {

	if _, err := os.Stat(test.root); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(test.root, 0777)
			if err != nil {
				FatalError(err, "unable to create mock password store")
			}
		}
	}

	f, err := os.OpenFile(test.keyfile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		FatalError(err, "unable to open mock keyfile")
	}
	defer f.Close()

	f.Write(test.mockKey)

	for _, mock := range mocks {
		Add(test.root, test.keyfile, mock.prefix, mock.pswd)
	}

}

func teardown() {
	err := os.RemoveAll(test.root)
	if err != nil {
		FatalError(err, "unable to teardown test suite")
	}
}
