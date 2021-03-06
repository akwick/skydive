/*
 * Copyright (C) 2016 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy ofthe License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specificlanguage governing permissions and
 * limitations under the License.
 *
 */

package client

import (
	"encoding/json"

	"github.com/skydive-project/skydive/flow"
	gclient "github.com/skydive-project/skydive/graffiti/api/client"
	shttp "github.com/skydive-project/skydive/graffiti/http"
	"github.com/skydive-project/skydive/sflow"
	"github.com/skydive-project/skydive/topology"
	"github.com/skydive-project/skydive/topology/probes/socketinfo"
)

// ErrNoResult is returned when a query returned no result
var ErrNoResult = gclient.ErrNoResult

// GremlinQueryHelper describes a gremlin query request query helper mechanism
type GremlinQueryHelper struct {
	*gclient.GremlinQueryHelper
}

// GetFlows from the Gremlin query
func (g *GremlinQueryHelper) GetFlows(query interface{}) ([]*flow.Flow, error) {
	data, err := g.Query(query)
	if err != nil {
		return nil, err
	}

	var flows []*flow.Flow
	if err := json.Unmarshal(data, &flows); err != nil {
		return nil, err
	}

	return flows, nil
}

// GetInterfaceMetrics from Gremlin query
func (g *GremlinQueryHelper) GetInterfaceMetrics(query interface{}) (map[string][]*topology.InterfaceMetric, error) {
	data, err := g.Query(query)
	if err != nil {
		return nil, err
	}

	var result []map[string][]*topology.InterfaceMetric
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result[0], nil
}

// GetSFlowMetrics from Gremlin query
func (g *GremlinQueryHelper) GetSFlowMetrics(query interface{}) (map[string][]*sflow.SFMetric, error) {
	data, err := g.Query(query)
	if err != nil {
		return nil, err
	}

	var result []map[string][]*sflow.SFMetric
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result[0], nil
}

// GetFlowMetrics from Gremlin query
func (g *GremlinQueryHelper) GetFlowMetrics(query interface{}) (map[string][]*flow.FlowMetric, error) {
	data, err := g.Query(query)
	if err != nil {
		return nil, err
	}

	var result []map[string][]*flow.FlowMetric
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, ErrNoResult
	}

	return result[0], nil
}

// GetFlowMetric from Gremlin query
func (g *GremlinQueryHelper) GetFlowMetric(query interface{}) (*flow.FlowMetric, error) {
	data, err := g.Query(query)
	if err != nil {
		return nil, err
	}

	var result flow.FlowMetric
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetInterfaceMetric from Gremlin query
func (g *GremlinQueryHelper) GetInterfaceMetric(query interface{}) (*topology.InterfaceMetric, error) {
	data, err := g.Query(query)
	if err != nil {
		return nil, err
	}

	var result topology.InterfaceMetric
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetSFlowMetric from Gremlin query
func (g *GremlinQueryHelper) GetSFlowMetric(query interface{}) (*sflow.SFMetric, error) {
	data, err := g.Query(query)
	if err != nil {
		return nil, err
	}

	var result sflow.SFMetric
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetSockets from the Gremlin query
func (g *GremlinQueryHelper) GetSockets(query interface{}) (sockets map[string][]*socketinfo.ConnectionInfo, err error) {
	data, err := g.Query(query)
	if err != nil {
		return nil, err
	}

	// TODO: use real objects instead of interface + decode
	// should be []map[string][]ConnectionInfo
	var maps []map[string][]interface{}
	if err := json.Unmarshal(data, &maps); err != nil {
		return nil, err
	}

	sockets = make(map[string][]*socketinfo.ConnectionInfo)
	for id, objs := range maps[0] {
		sockets[id] = make([]*socketinfo.ConnectionInfo, 0)
		for _, obj := range objs {
			var socket socketinfo.ConnectionInfo
			if err = socket.Decode(obj); err == nil {
				sockets[id] = append(sockets[id], &socket)
			}
		}
	}

	return
}

// NewGremlinQueryHelper creates a new Gremlin query helper based on authentication
func NewGremlinQueryHelper(restClient *shttp.RestClient) *GremlinQueryHelper {
	return &GremlinQueryHelper{
		GremlinQueryHelper: gclient.NewGremlinQueryHelper(restClient),
	}
}

// NewGremlinQueryHelperFromConfig creates a new Gremlin query helper based on authentication based on configuration file
func NewGremlinQueryHelperFromConfig(authOptions *shttp.AuthenticationOpts) (*GremlinQueryHelper, error) {
	restClient, err := NewRestClientFromConfig(authOptions)
	if err != nil {
		return nil, err
	}

	return NewGremlinQueryHelper(restClient), nil
}
