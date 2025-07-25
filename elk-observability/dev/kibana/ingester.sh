#!/bin/bash
# Copyright 2025 Simon Emms <simon@simonemms.com>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


set -e

# Temporal Workflow Dashboard Import Script
# This script creates a comprehensive Temporal workflow monitoring dashboard in Kibana

# Configuration
KIBANA_HOST="${KIBANA_HOST:-http://localhost:5601}"
INDEX_PATTERN="temporal_visibility_v1_dev"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Function to make authenticated requests
kibana_request() {
    local method=$1
    local endpoint=$2
    local data=$3

    if [ -z "$data" ]; then
        curl -s -X $method \
            -H "kbn-xsrf: true" \
            -H "Content-Type: application/json" \
            "${KIBANA_HOST}${endpoint}"
    else
        curl -s -X $method \
            -H "kbn-xsrf: true" \
            -H "Content-Type: application/json" \
            -d "$data" \
            "${KIBANA_HOST}${endpoint}"
    fi
}

echo "🔍 Checking Kibana connection..."
if ! kibana_request GET /api/status | grep -q "\"level\":\"available\""; then
    echo -e "${RED}❌ Failed to connect to Kibana at $KIBANA_HOST${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Connected to Kibana${NC}"

# Step 1: Create or get index pattern
echo "📊 Setting up index pattern for Temporal visibility..."
INDEX_PATTERN_ID=$(kibana_request GET "/api/saved_objects/_find?type=index-pattern&search_fields=title&search=${INDEX_PATTERN}" | jq -r '.saved_objects[0].id // empty')

if [ -z "$INDEX_PATTERN_ID" ]; then
    echo "Creating index pattern..."
    RESPONSE=$(kibana_request POST /api/saved_objects/index-pattern "{
        \"attributes\": {
            \"title\": \"${INDEX_PATTERN}\",
            \"timeFieldName\": \"StartTime\"
        }
    }")
    INDEX_PATTERN_ID=$(echo $RESPONSE | jq -r '.id')
fi
echo -e "${GREEN}✅ Index pattern ID: $INDEX_PATTERN_ID${NC}"

# Step 2: Create visualizations
echo "📈 Creating visualizations..."

# Visualization 1: Workflow Execution Timeline
echo "  - Creating Workflow Execution Timeline..."
VIZ1=$(kibana_request POST /api/saved_objects/lens '{
  "attributes": {
    "title": "Workflow Execution Timeline",
    "visualizationType": "lnsXY",
    "state": {
      "datasourceStates": {
        "formBased": {
          "layers": {
            "layer1": {
              "columns": {
                "date_col": {
                  "label": "Start Time",
                  "dataType": "date",
                  "operationType": "date_histogram",
                  "sourceField": "StartTime",
                  "isBucketed": true,
                  "scale": "interval",
                  "params": {
                    "interval": "5m",
                    "includeEmptyRows": true
                  }
                },
                "completed_col": {
                  "label": "Completed",
                  "dataType": "number",
                  "operationType": "count",
                  "isBucketed": false,
                  "scale": "ratio",
                  "filter": {
                    "query": "ExecutionStatus: Completed",
                    "language": "kuery"
                  }
                },
                "failed_col": {
                  "label": "Failed",
                  "dataType": "number",
                  "operationType": "count",
                  "isBucketed": false,
                  "scale": "ratio",
                  "filter": {
                    "query": "ExecutionStatus: Failed",
                    "language": "kuery"
                  }
                },
                "running_col": {
                  "label": "Running",
                  "dataType": "number",
                  "operationType": "count",
                  "isBucketed": false,
                  "scale": "ratio",
                  "filter": {
                    "query": "ExecutionStatus: Running",
                    "language": "kuery"
                  }
                }
              },
              "columnOrder": ["date_col", "completed_col", "failed_col", "running_col"],
              "indexPatternId": "'$INDEX_PATTERN_ID'"
            }
          }
        }
      },
      "visualization": {
        "legend": {
          "isVisible": true,
          "position": "right"
        },
        "valueLabels": "hide",
        "fittingFunction": "None",
        "axisTitlesVisibilitySettings": {
          "x": true,
          "yLeft": true
        },
        "tickLabelsVisibilitySettings": {
          "x": true,
          "yLeft": true
        },
        "gridlinesVisibilitySettings": {
          "x": false,
          "yLeft": true
        },
        "preferredSeriesType": "line",
        "layers": [
          {
            "layerId": "layer1",
            "seriesType": "line",
            "xAccessor": "date_col",
            "accessors": ["completed_col", "failed_col", "running_col"],
            "yConfig": [
              {"forAccessor": "completed_col", "color": "#2ECC71"},
              {"forAccessor": "failed_col", "color": "#E74C3C"},
              {"forAccessor": "running_col", "color": "#3498DB"}
            ]
          }
        ]
      },
      "query": {
        "query": "",
        "language": "kuery"
      },
      "filters": []
    }
  },
  "references": [
    {
      "type": "index-pattern",
      "id": "'$INDEX_PATTERN_ID'",
      "name": "indexpattern-datasource-layer-layer1"
    }
  ]
}')
VIZ1_ID=$(echo $VIZ1 | jq -r '.id')
echo -e "${GREEN}    ✓ Created: $VIZ1_ID${NC}"

# Visualization 2: Workflow Status Distribution (Pie Chart)
echo "  - Creating Workflow Status Distribution..."
VIZ2=$(kibana_request POST /api/saved_objects/lens '{
  "attributes": {
    "title": "Workflow Status Distribution",
    "visualizationType": "lnsPie",
    "state": {
      "datasourceStates": {
        "formBased": {
          "layers": {
            "layer1": {
              "columns": {
                "status_col": {
                  "label": "Status",
                  "dataType": "string",
                  "operationType": "terms",
                  "sourceField": "ExecutionStatus",
                  "isBucketed": true,
                  "scale": "ordinal",
                  "params": {
                    "size": 10,
                    "orderBy": {
                      "type": "column",
                      "columnId": "count_col"
                    },
                    "orderDirection": "desc"
                  }
                },
                "count_col": {
                  "label": "Count",
                  "dataType": "number",
                  "operationType": "count",
                  "isBucketed": false,
                  "scale": "ratio"
                }
              },
              "columnOrder": ["status_col", "count_col"],
              "indexPatternId": "'$INDEX_PATTERN_ID'"
            }
          }
        }
      },
      "visualization": {
        "shape": "donut",
        "palette": {
          "type": "palette",
          "name": "status",
          "params": {
            "colors": {
              "Completed": "#2ECC71",
              "Failed": "#E74C3C",
              "Running": "#3498DB",
              "Terminated": "#F39C12",
              "ContinuedAsNew": "#9B59B6",
              "TimedOut": "#E67E22"
            }
          }
        },
        "legendDisplay": "show",
        "legendPosition": "right",
        "nestedLegend": false,
        "percentDecimals": 1,
        "emptySizeRatio": 0.4,
        "layers": [
          {
            "layerId": "layer1",
            "categoryDisplay": "default",
            "legendDisplay": "show",
            "primaryGroups": ["status_col"],
            "metrics": ["count_col"]
          }
        ]
      },
      "query": {
        "query": "",
        "language": "kuery"
      },
      "filters": []
    }
  },
  "references": [
    {
      "type": "index-pattern",
      "id": "'$INDEX_PATTERN_ID'",
      "name": "indexpattern-datasource-layer-layer1"
    }
  ]
}')
VIZ2_ID=$(echo $VIZ2 | jq -r '.id')
echo -e "${GREEN}    ✓ Created: $VIZ2_ID${NC}"

# Visualization 3: Average Execution Duration Over Time
echo "  - Creating Average Execution Duration Over Time..."
VIZ3=$(kibana_request POST /api/saved_objects/lens '{
  "attributes": {
    "title": "Average Execution Duration Over Time",
    "visualizationType": "lnsXY",
    "state": {
      "datasourceStates": {
        "formBased": {
          "layers": {
            "layer1": {
              "columns": {
                "date_col": {
                  "label": "Time",
                  "dataType": "date",
                  "operationType": "date_histogram",
                  "sourceField": "StartTime",
                  "isBucketed": true,
                  "scale": "interval",
                  "params": {
                    "interval": "10m",
                    "includeEmptyRows": true
                  }
                },
                "avg_duration": {
                  "label": "Average Duration",
                  "dataType": "number",
                  "operationType": "average",
                  "sourceField": "ExecutionDuration",
                  "isBucketed": false,
                  "scale": "ratio"
                },
                "p95_duration": {
                  "label": "95th Percentile",
                  "dataType": "number",
                  "operationType": "percentile",
                  "sourceField": "ExecutionDuration",
                  "isBucketed": false,
                  "scale": "ratio",
                  "params": {
                    "percentile": 95
                  }
                },
                "p99_duration": {
                  "label": "99th Percentile",
                  "dataType": "number",
                  "operationType": "percentile",
                  "sourceField": "ExecutionDuration",
                  "isBucketed": false,
                  "scale": "ratio",
                  "params": {
                    "percentile": 99
                  }
                }
              },
              "columnOrder": ["date_col", "avg_duration", "p95_duration", "p99_duration"],
              "indexPatternId": "'$INDEX_PATTERN_ID'"
            }
          }
        }
      },
      "visualization": {
        "legend": {
          "isVisible": true,
          "position": "right"
        },
        "valueLabels": "hide",
        "fittingFunction": "None",
        "axisTitlesVisibilitySettings": {
          "x": true,
          "yLeft": true
        },
        "tickLabelsVisibilitySettings": {
          "x": true,
          "yLeft": true
        },
        "gridlinesVisibilitySettings": {
          "x": false,
          "yLeft": true
        },
        "preferredSeriesType": "area",
        "layers": [
          {
            "layerId": "layer1",
            "seriesType": "area",
            "xAccessor": "date_col",
            "accessors": ["avg_duration", "p95_duration", "p99_duration"],
            "yConfig": [
              {"forAccessor": "avg_duration", "color": "#3498DB"},
              {"forAccessor": "p95_duration", "color": "#E74C3C"},
              {"forAccessor": "p99_duration", "color": "#F39C12"}
            ]
          }
        ],
        "yLeftExtent": {
          "mode": "full"
        }
      },
      "query": {
        "query": "",
        "language": "kuery"
      },
      "filters": []
    }
  },
  "references": [
    {
      "type": "index-pattern",
      "id": "'$INDEX_PATTERN_ID'",
      "name": "indexpattern-datasource-layer-layer1"
    }
  ]
}')
VIZ3_ID=$(echo $VIZ3 | jq -r '.id')
echo -e "${GREEN}    ✓ Created: $VIZ3_ID${NC}"

# Visualization 4: Workflow Metrics Summary
echo "  - Creating Workflow Metrics Summary..."
VIZ4=$(kibana_request POST /api/saved_objects/lens '{
  "attributes": {
    "title": "Workflow Metrics Summary",
    "visualizationType": "lnsMetric",
    "state": {
      "datasourceStates": {
        "formBased": {
          "layers": {
            "layer1": {
              "columns": {
                "total_workflows": {
                  "label": "Total Workflows",
                  "dataType": "number",
                  "operationType": "unique_count",
                  "sourceField": "WorkflowId",
                  "isBucketed": false,
                  "scale": "ratio",
                  "params": {
                    "format": {
                      "id": "number",
                      "params": {
                        "decimals": 0
                      }
                    }
                  }
                }
              },
              "columnOrder": ["total_workflows"],
              "indexPatternId": "'$INDEX_PATTERN_ID'"
            },
            "layer2": {
              "columns": {
                "avg_duration": {
                  "label": "Avg Duration",
                  "dataType": "number",
                  "operationType": "average",
                  "sourceField": "ExecutionDuration",
                  "isBucketed": false,
                  "scale": "ratio",
                  "params": {
                    "format": {
                      "id": "duration",
                      "params": {
                        "inputFormat": "nanoseconds",
                        "outputFormat": "humanizePrecise"
                      }
                    }
                  }
                }
              },
              "columnOrder": ["avg_duration"],
              "indexPatternId": "'$INDEX_PATTERN_ID'"
            },
            "layer3": {
              "columns": {
                "avg_history_size": {
                  "label": "Avg History Size",
                  "dataType": "number",
                  "operationType": "average",
                  "sourceField": "HistorySizeBytes",
                  "isBucketed": false,
                  "scale": "ratio",
                  "params": {
                    "format": {
                      "id": "bytes",
                      "params": {
                        "decimals": 2
                      }
                    }
                  }
                }
              },
              "columnOrder": ["avg_history_size"],
              "indexPatternId": "'$INDEX_PATTERN_ID'"
            },
            "layer4": {
              "columns": {
                "avg_transitions": {
                  "label": "Avg State Transitions",
                  "dataType": "number",
                  "operationType": "average",
                  "sourceField": "StateTransitionCount",
                  "isBucketed": false,
                  "scale": "ratio",
                  "params": {
                    "format": {
                      "id": "number",
                      "params": {
                        "decimals": 1
                      }
                    }
                  }
                }
              },
              "columnOrder": ["avg_transitions"],
              "indexPatternId": "'$INDEX_PATTERN_ID'"
            }
          }
        }
      },
      "visualization": {
        "layerId": "layer1",
        "layerType": "data",
        "metricAccessor": "total_workflows",
        "secondaryMetricAccessor": "avg_duration",
        "maxAccessor": "avg_history_size",
        "breakdownByAccessor": "avg_transitions",
        "color": "#3498DB"
      },
      "query": {
        "query": "",
        "language": "kuery"
      },
      "filters": []
    }
  },
  "references": [
    {
      "type": "index-pattern",
      "id": "'$INDEX_PATTERN_ID'",
      "name": "indexpattern-datasource-layer-layer1"
    },
    {
      "type": "index-pattern",
      "id": "'$INDEX_PATTERN_ID'",
      "name": "indexpattern-datasource-layer-layer2"
    },
    {
      "type": "index-pattern",
      "id": "'$INDEX_PATTERN_ID'",
      "name": "indexpattern-datasource-layer-layer3"
    },
    {
      "type": "index-pattern",
      "id": "'$INDEX_PATTERN_ID'",
      "name": "indexpattern-datasource-layer-layer4"
    }
  ]
}')
VIZ4_ID=$(echo $VIZ4 | jq -r '.id')
echo -e "${GREEN}    ✓ Created: $VIZ4_ID${NC}"

# Step 3: Create Dashboard
echo "🎯 Creating Temporal Workflow Monitoring Dashboard..."
DASHBOARD=$(kibana_request POST /api/saved_objects/dashboard '{
  "attributes": {
    "title": "Temporal Workflow Monitoring Dashboard",
    "description": "Comprehensive monitoring dashboard for Temporal workflows",
    "kibanaSavedObjectMeta": {
      "searchSourceJSON": "{\"query\":{\"query\":\"\",\"language\":\"kuery\"},\"filter\":[]}"
    },
    "timeRestore": true,
    "timeFrom": "now-24h",
    "timeTo": "now",
    "refreshInterval": {
      "pause": false,
      "value": 60000
    },
    "panelsJSON": "[{\"version\":\"8.0.0\",\"type\":\"lens\",\"gridData\":{\"x\":0,\"y\":0,\"w\":48,\"h\":16,\"i\":\"panel_1\"},\"panelIndex\":\"panel_1\",\"embeddableConfig\":{\"enhancements\":{},\"hidePanelTitles\":false},\"panelRefName\":\"panel_1\"},{\"version\":\"8.0.0\",\"type\":\"lens\",\"gridData\":{\"x\":0,\"y\":16,\"w\":24,\"h\":16,\"i\":\"panel_2\"},\"panelIndex\":\"panel_2\",\"embeddableConfig\":{\"enhancements\":{},\"hidePanelTitles\":false},\"panelRefName\":\"panel_2\"},{\"version\":\"8.0.0\",\"type\":\"lens\",\"gridData\":{\"x\":24,\"y\":16,\"w\":24,\"h\":16,\"i\":\"panel_3\"},\"panelIndex\":\"panel_3\",\"embeddableConfig\":{\"enhancements\":{},\"hidePanelTitles\":false},\"panelRefName\":\"panel_3\"},{\"version\":\"8.0.0\",\"type\":\"lens\",\"gridData\":{\"x\":0,\"y\":32,\"w\":48,\"h\":8,\"i\":\"panel_4\"},\"panelIndex\":\"panel_4\",\"embeddableConfig\":{\"enhancements\":{},\"hidePanelTitles\":false},\"panelRefName\":\"panel_4\"}]"
  },
  "references": [
    {
      "name": "panel_1",
      "type": "lens",
      "id": "'$VIZ1_ID'"
    },
    {
      "name": "panel_2",
      "type": "lens",
      "id": "'$VIZ2_ID'"
    },
    {
      "name": "panel_3",
      "type": "lens",
      "id": "'$VIZ3_ID'"
    },
    {
      "name": "panel_4",
      "type": "lens",
      "id": "'$VIZ4_ID'"
    }
  ]
}')

DASHBOARD_ID=$(echo $DASHBOARD | jq -r '.id')

if [ ! -z "$DASHBOARD_ID" ] && [ "$DASHBOARD_ID" != "null" ]; then
    echo -e "${GREEN}✅ Dashboard created successfully!${NC}"
    echo ""
    echo "📊 Access your dashboard at:"
    echo "   ${KIBANA_HOST}/app/dashboards#/view/${DASHBOARD_ID}"
    echo ""
    echo "📝 Dashboard includes:"
    echo "   - Workflow Execution Timeline (by status)"
    echo "   - Workflow Status Distribution (pie chart)"
    echo "   - Average Execution Duration Over Time (with percentiles)"
    echo "   - Workflow Metrics Summary"
    echo ""
    echo "🎯 Additional visualizations you can add:"
    echo "   - Task Queue Distribution"
    echo "   - Workflow Type Distribution"
    echo "   - Namespace Activity Heatmap"
    echo "   - Recent Failed Workflows Table"
else
    echo -e "${RED}❌ Failed to create dashboard${NC}"
    echo "Response: $DASHBOARD"
fi

# Save configuration
cat > /tmp/temporal-dashboard-config.json << EOF
{
  "kibana_host": "$KIBANA_HOST",
  "index_pattern_id": "$INDEX_PATTERN_ID",
  "visualizations": {
    "workflow_timeline": "$VIZ1_ID",
    "status_distribution": "$VIZ2_ID",
    "duration_trend": "$VIZ3_ID",
    "metrics_summary": "$VIZ4_ID"
  },
  "dashboard_id": "$DASHBOARD_ID"
}
EOF

echo ""
echo "💾 Configuration saved to /tmp/temporal-dashboard-config.json"
