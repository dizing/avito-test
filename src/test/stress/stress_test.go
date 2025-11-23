package stress

import (
	"avito-test/test"
	"testing"
)

func TestPrepare(t *testing.T) {
	env := test.NewEnvironment(t)
	env.ShopDatabase.CreateTestUsers()
}
