package url2png

import (
	"io/ioutil"
	"os"
	"testing"

	"golang.org/x/net/context"

	"github.com/stretchr/testify/assert"
)

func newClient() Client {
	return Client{
		Key:    os.Getenv("TEST_KEY"),
		Secret: os.Getenv("TEST_SECRET"),
	}
}

func TestScreenshot(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	c := newClient()

	r, err := c.Screenshot("http://www.fknsrs.biz/", nil)
	a.NoError(err)
	if !a.NotNil(r) {
		return
	}

	n, err := ioutil.ReadAll(r)
	a.NoError(err)
	a.NotEqual(0, n)
}

func TestCancel(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	c := newClient()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r, err := c.ScreenshotWithContext(ctx, "http://www.fknsrs.biz/", nil)
	a.Error(err)
	a.Equal("context canceled", err.Error())
	a.Nil(r)
}

func TestBadURL(t *testing.T) {
	t.Parallel()

	a := assert.New(t)

	c := newClient()

	r, err := c.Screenshot("zzzzzzzz", nil)
	a.Error(err)
	a.Nil(r)
}
