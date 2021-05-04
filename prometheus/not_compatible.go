// Copyright 2020 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.  // You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prometheus

import (
	"context"
	"fmt"
	"strings"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/labels"
)

// notCompatibleHTTPClient must be used to contact a distant prometheus with a version < v2.15.
type notCompatibleHTTPClient struct {
	MetadataService
	prometheusClient v1.API
	lookbackInterval time.Duration
}

func (c *notCompatibleHTTPClient) MetricMetadata(ctx context.Context, metric string) (v1.Metadata, error) {
	metadata, err := c.prometheusClient.TargetsMetadata(ctx, "", metric, "1")
	if err != nil {
		return v1.Metadata{}, err
	}
	if len(metadata) == 0 {
		return v1.Metadata{}, nil
	}
	return v1.Metadata{
		Type: metadata[0].Type,
		Help: metadata[0].Help,
		Unit: metadata[0].Unit,
	}, nil
}

func (c *notCompatibleHTTPClient) AllMetricMetadata(ctx context.Context) (map[string][]v1.Metadata, error) {
	metricNames, _, err := c.prometheusClient.LabelValues(ctx, "__name__", time.Now().Add(-100*time.Hour), time.Now())
	if err != nil {
		return nil, err
	}
	allMetadata := make(map[string][]v1.Metadata)
	for _, name := range metricNames {
		allMetadata[string(name)] = []v1.Metadata{{}}
	}
	return allMetadata, nil
}

func (c *notCompatibleHTTPClient) LabelNames(
	ctx context.Context,
	currMatchers []*labels.Matcher,
) ([]string, error) {
	if len(currMatchers) == 0 {
		names, _, err := c.prometheusClient.LabelNames(ctx, time.Now().Add(-1*c.lookbackInterval), time.Now())
		return names, err
	}

	labelNameAndValues, err := uniqueLabelNameAndValues(ctx, c.prometheusClient,
		time.Now().Add(-1*c.lookbackInterval), time.Now(), currMatchers)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, len(labelNameAndValues))
	for l := range labelNameAndValues {
		result = append(result, l)
	}

	return result, nil
}

func (c *notCompatibleHTTPClient) LabelValues(
	ctx context.Context,
	label string,
	currMatchers []*labels.Matcher,
) ([]model.LabelValue, error) {
	if len(currMatchers) == 0 {
		values, _, err := c.prometheusClient.LabelValues(ctx, label, time.Now().Add(-1*c.lookbackInterval), time.Now())
		return values, err
	}

	labelNameAndValues, err := uniqueLabelNameAndValues(ctx, c.prometheusClient,
		time.Now().Add(-1*c.lookbackInterval), time.Now(), currMatchers)
	if err != nil {
		return nil, err
	}

	labelValues, ok := labelNameAndValues[label]
	if !ok {
		return nil, nil
	}

	result := make([]model.LabelValue, 0, len(labelValues))
	for l := range labelValues {
		result = append(result, model.LabelValue(l))
	}

	return result, nil
}

func (c *notCompatibleHTTPClient) ChangeDataSource(_ string) error {
	return fmt.Errorf("method not supported")
}

func (c *notCompatibleHTTPClient) SetLookbackInterval(interval time.Duration) {
	c.lookbackInterval = interval
}

func (c *notCompatibleHTTPClient) GetURL() string {
	return ""
}

func uniqueLabelNameAndValues(
	ctx context.Context,
	prometheusClient v1.API,
	start, end time.Time,
	currMatchers []*labels.Matcher,
) (map[string]map[string]struct{}, error) {
	var (
		metricName   string
		metricLabels []*labels.Matcher
	)
	for _, matcher := range currMatchers {
		if matcher.Name == model.MetricNameLabel {
			metricName = string(matcher.Value)
		} else {
			metricLabels = append(metricLabels, matcher)
		}
	}

	var match strings.Builder
	if _, err := match.WriteString(metricName); err != nil {
		return nil, err
	}
	if len(metricLabels) > 0 {
		if _, err := match.WriteString("{"); err != nil {
			return nil, err
		}
		for i, matcher := range metricLabels {
			if i != 0 {
				if _, err := match.WriteString(","); err != nil {
					return nil, err
				}
			}
			if _, err := match.WriteString(matcher.String()); err != nil {
				return nil, err
			}
		}
		if _, err := match.WriteString("}"); err != nil {
			return nil, err
		}
	}

	results, _, err := prometheusClient.Series(ctx, []string{match.String()},
		start, end)
	if err != nil {
		return nil, err
	}

	// deduplicated is a de-duplicated result set.
	deduplicated := make(map[string]map[string]struct{})
	for _, labelSet := range results {
		for name, value := range labelSet {
			setKey := string(name)
			curr, ok := deduplicated[setKey]
			if !ok {
				curr = map[string]struct{}{}
				deduplicated[setKey] = curr
			}
			setValue := string(value)
			if _, exists := curr[setValue]; exists {
				continue
			}
			curr[setValue] = struct{}{}
		}
	}

	return deduplicated, nil
}
