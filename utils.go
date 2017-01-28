/*
* Minio Cloud Storage, (C) 2017 Minio, Inc.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/volume"
)

// return `Host` from the URL endpoint.
func getHost(endpoint string) (string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}
	return u.Host, nil
}

// determines if the url has HTTPS scheme.
func isSSL(url string) (bool, error) {
	scheme, err := getScheme(url)
	if err != nil {
		return false, err
	}
	if scheme == "https" {
		return true, nil
	}
	return false, nil

}

// Parse the server endpoint to find out the scheme(http,https...).
func getScheme(endpoint string) (string, error) {
	// parse the URL.
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}
	// return the scheme.
	return u.Scheme, nil
}

// If the requested volume alredy exists, then its necessary that the server configs (Minio server endpoint,
// bucket,accessKey and secretKey matches with the existing one.
// Since a mount is uniquely identified by its volume name its not possible to have a duplicate entry.
func matchServerConfig(config serverConfig, r volume.Request) error {
	if r.Options == nil {
		return fmt.Errorf("No options provided. Please refer example usage.")
	}
	// Compare the endpoints.
	if r.Options["endpoint"] == config.endpoint {
		return fmt.Errorf("Volume \"%s\" already exists and is pointing to Minio server\"%s\",Cannot create duplicate volume.",
			r.Name, config.endpoint)
	}
	// Compare the bucket name.
	if r.Options["bucket"] == config.bucket {
		return fmt.Errorf("Volume \"%s\" already exists and is pointing to Minio server \"%s\", and bucket \"%s\",Cannot create duplicate volume.",
			r.Name, config.endpoint)
	}
	// compare the access keys.
	if r.Options["access-key"] == "" {
		return fmt.Errorf("Volume \"%s\" already exists, access key mismatch.", r.Name)

	}
	// compare the secret keys.
	if r.Options["secret-key"] == "" {
		return fmt.Errorf("Volume \"%s\" already exists, secret key mismatch.", r.Name)
	}
	// match successful, return `nil` error.
	return nil
}

// Error repsonse to be sent to docker on failure of any operation.
func errorResponse(err string) volume.Response {
	logrus.Error(err)
	return volume.Response{Err: err}
}

// create directory for the given path.
func createDir(path string) error {
	// verify whether the directory already exists.
	fi, err := os.Lstat(path)
	// create the directory doesn't exist.
	if os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	// if the file already exists, very that it is a directory.
	if fi != nil && !fi.IsDir() {
		return fmt.Errorf("%v already exist and it's not a directory", path)
	}
	return nil
}
