/*
 * Copyright 2025 Simon Emms <simon@simonemms.com>
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

package temporal

import (
	"crypto/tls"
	"fmt"
	"os"
	"strings"

	"go.temporal.io/sdk/client"
)

func NewClient(opts ...client.Options) (client.Client, error) {
	hostPort := os.Getenv("TEMPORAL_ADDRESS")
	if hostPort == "" {
		hostPort = client.DefaultHostPort
	}
	namespace := os.Getenv("TEMPORAL_NAMESPACE")
	if namespace == "" {
		namespace = client.DefaultNamespace
	}
	connectionOpts := client.ConnectionOptions{}
	if strings.ToLower(os.Getenv("TEMPORAL_USE_TLS")) == "true" {
		connectionOpts.TLS = &tls.Config{}
	}
	var creds client.Credentials
	if key := os.Getenv("TEMPORAL_API_KEY"); key != "" {
		creds = client.NewAPIKeyStaticCredentials(key)
	}

	clientOpts := client.Options{
		HostPort:          hostPort,
		Namespace:         namespace,
		ConnectionOptions: connectionOpts,
		Credentials:       creds,
	}

	if len(opts) > 0 {
		o := opts[0]
		clientOpts.Interceptors = o.Interceptors
		clientOpts.MetricsHandler = o.MetricsHandler
	}

	c, err := client.Dial(clientOpts)
	if err != nil {
		return nil, fmt.Errorf("unable to create client: %w", err)
	}

	return c, nil
}
