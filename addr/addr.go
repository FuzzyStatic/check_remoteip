/**
 * @Author: Allen Flickinger <FuzzyStatic>
 * @Date:   2017-09-14T18:53:38-04:00
 * @Email:  allen.flickinger@gmail.com
 * @Last modified by:   FuzzyStatic
 * @Last modified time: 2017-09-15T21:26:08-04:00
 */

package addr

import (
	"io/ioutil"
	"net/http"
)

// GetRemoteIP gets the remote IP as reported by ipecho.net
func GetRemoteIP() (string, error) {
	var (
		resp *http.Response
		err  error
		body []byte
	)

	if resp, err = http.Get("http://ipecho.net/plain"); err != nil {
		return "", err
	}

	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return "", err
	}

	resp.Body.Close()
	return string(body), err
}
