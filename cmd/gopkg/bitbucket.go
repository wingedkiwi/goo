// (c) Chi Vinh Le <cvl@chinet.info> – 13.06.2015

package gopkg

import (
	"fmt"
	"github.com/garyburd/go-oauth/oauth"
	"io/ioutil"
)

type BitbucketDiscovery struct {
	accessToken *oauth.Credentials
	client      *oauth.Client
	user        string
}

func NewBitbucketDiscovery(
	user string,
	consumerKey string, consumerSecret string,
	accessToken string, accessTokenSecret string,
) *BitbucketDiscovery {
	return &BitbucketDiscovery{
		user:        user,
		accessToken: &oauth.Credentials{accessToken, accessTokenSecret},
		client:      &oauth.Client{Credentials: oauth.Credentials{consumerKey, consumerSecret}},
	}
}

func (d *BitbucketDiscovery) GetRepository(path string) (string, error) {
	var data []byte
	resp, err := d.client.Get(nil, d.accessToken, "https://api.bitbucket.org/2.0/repositories/"+d.user+"/"+path, nil)
	if err == nil {
		if resp != nil {
			data, err = ioutil.ReadAll(resp.Body)
			if err == nil {
				if resp.StatusCode == 404 {
					return "", nil
				} else if resp.StatusCode != 200 {
					return "", fmt.Errorf("Unexpected status code: %d, with content %s", resp.StatusCode, data)
				}
			}
		}
	}
	if err != nil {
		return "", nil
	}
	return "https://bitbucket.org/" + d.user + "/" + path, nil
}
