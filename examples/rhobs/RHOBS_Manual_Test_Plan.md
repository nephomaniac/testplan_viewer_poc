# RHOBS Synthetic Monitoring \- Manual Test Plan

**Version:** 1.0 **Date:** February 17, 2026 **Epic:** [SREP-3109](https://issues.redhat.com/browse/SREP-3109) **Status:** Draft \- Needs Team Validation

---

## Table of Contents

1. [Overview](#overview)  
2. [Prerequisites](#prerequisites)  
3. [Test Environment Setup](#test-environment-setup)  
4. [Test Scenarios](#test-scenarios)  
   - [Test 1: Metrics Flow End-to-End](#test-1-metrics-flow-end-to-end)  
   - [Test 2: Logs Flow End-to-End](#test-2-logs-flow-end-to-end)  
   - [Test 3: Synthetic Monitoring Flow End-to-End](#test-3-synthetic-monitoring-flow-end-to-end)  
5. [Cleanup Procedures](#cleanup-procedures)  
6. [Troubleshooting](#troubleshooting)

---

## Overview

### Purpose

This manual test plan validates the complete RHOBS synthetic monitoring stack end-to-end. It serves as:

- A specification for what "working" looks like  
- A foundation for automated test development  
- A validation procedure for new deployments  
- A troubleshooting reference

### System Under Test

The RHOBS synthetic monitoring system consists of:

```
┌────────────────────────────────────────┐
│  Management Cluster (MC)                                          │
│  ┌───────────────────────────────────┐   │
│  │  route-monitor-operator                                  │   │
│  │  - Watches HostedCluster CRs                             │   │
│  │  - Creates/updates probes in RHOBS API                   │   │
│  └──────────────────┬────────────────┘   │
│                                 │                                │
│  ┌──────────────────▼────────────────┐   │
│  │  Tenant HCP Cluster (control plane hosted on MC)         │   │
│  │  - API Server: api.cluster-id.base-domain                │   │
│  │  - Metrics: forwarded to RHOBS                           │   │
│  │  - Logs: forwarded to Loki                               │   │
│  └───────────────────────────────────┘   │
└────────────────────────────────────────┘
                      │
                      ▼
┌───────────────────────────────────────┐
│  RHOBS Cell Cluster (regional observability service)            │
│  ┌───────────────────────────────────┐   │
│  │  rhobs-synthetics-api                                    │   │
│  │  - REST API for probe management                         │   │
│  │  - Stores probe configs as ConfigMaps                    │   │
│  └──────────────────┬────────────────┘   │
│                                  │                               │
│  ┌──────────────────▼────────────────┐   │
│  │  rhobs-synthetics-agent                                  │   │
│  │  - Polls API for probe configs                           │   │
│  │  - Creates Probe CRs for blackbox exporter               │   │
│  │  - Validates target URLs                                 │   │
│  │  - Updates probe status                                  │   │
│  └──────────────────┬────────────────┘   │
│                                  │                               │
│  ┌──────────────────▼────────────────┐   │
│  │  Blackbox Exporter + Prometheus                          │   │
│  │  - Executes HTTP probes against cluster APIs             │   │
│  │  - Generates probe_success metrics                       │   │
│  └──────────────────┬────────────────┘   │
│                                  │                               │
│  ┌──────────────────▼────────────────┐   │
│  │  Thanos / Observatorium                                  │   │
│  │  - Stores metrics (including probe results)              │   │
│  │  - Queryable via Prometheus API                          │   │
│  └───────────────────────────────────┘   │
│                                                                  │
│  ┌───────────────────────────────────┐   │
│  │  Loki                                                    │   │
│  │  - Stores logs from clusters                             │   │
│  │  - Queryable via LogQL API                               │   │
│  └───────────────────────────────────┘   │
└────────────────────────────────────────┘
                      │
                      ▼
┌────────────────────────────────────────┐
│  Grafana                                                         │
│  - Visualizes metrics from Thanos                                │
│  - Visualizes logs from Loki                                     │
└────────────────────────────────────────┘
```

### Test Flows Validated

1. **Metrics Flow**: HCP Cluster → RHOBS → Thanos → Grafana  
2. **Logs Flow**: HCP Cluster → Loki → Grafana  
3. **Synthetics Flow**: RMO → API → Agent → Blackbox → Thanos → Grafana

---

## Prerequisites

### Required Access & Credentials

#### 1\. OCM (OpenShift Cluster Manager) Access

- **Purpose:** Create ROSA HCP clusters in stage environment  
- **Obtain:**

```shell
# Login to OCM stage
ocm login --url=https://api.stage.openshift.com

# Verify access
ocm whoami
```

- **Permissions Required:**  
  - Ability to create ROSA HCP clusters  
  - Access to `srep-dev` sector/provision shard

#### 2\. Kubernetes Access (Management Cluster)

- **Purpose:** Verify route-monitor-operator behavior, check HostedCluster CRs  
- **Obtain:**

```shell
# Login to MC via backplane
ocm backplane login <mc-cluster-id>

# Verify access
oc whoami
oc get hostedclusters -A
```

- **Permissions Required:**  
  - Read access to HostedCluster CRs  
  - Read access to route-monitor-operator namespace  
  - Read access to probe-related resources

#### 3\. RHOBS API Access

- **Purpose:** Query probe registrations, verify API integration  
- **Obtain:**  
  - Get API endpoint from team: `https://observatorium-api-<region>.stage.openshift.com`  
  - Get authentication token from Vault or team lead  
- **Test Access:**

```shell
# List all probes
curl -H "Authorization: Bearer $RHOBS_TOKEN" \
  "https://observatorium-api.stage.openshift.com/api/metrics/v1/<tenant>/probes"
```

#### 4\. Grafana Access

- **Purpose:** Verify metrics are queryable and visible  
- **URL:** `https://grafana.stage.openshift.com` (or equivalent)  
- **Authentication:** SSO or token-based

#### 5\. Loki Access

- **Purpose:** Query logs directly via API  
- **Obtain:**  
  - Get Loki endpoint from team  
  - Same auth token as RHOBS typically works  
- **Test Access:**

```shell
# Query recent logs
curl -H "Authorization: Bearer $LOKI_TOKEN" \
  "https://loki.stage.openshift.com/loki/api/v1/query?query={cluster_id=\"test\"}"
```

### Required Tools

#### Mac

Install the following tools before beginning tests:

```shell
# OCM CLI
brew install ocm

# OpenShift CLI
brew install openshift-cli

# jq for JSON parsing
brew install jq

# curl (should be pre-installed)
curl --version
```

#### Fedora

Install the following tools before beginning tests:

```shell
# jq for JSON parsing and curl (should be pre-installed)
sudo dnf install jq curl
curl --version
jq --version

# use backplane-tools to install the 'oc' and 'ocm' executables
# https://github.com/openshift/backplane-tools/blob/main/README.md
```

---

## Test Environment Setup

### Environment Variables

Set up your environment with necessary credentials:

```shell
# OCM
export OCM_TOKEN=$(ocm token)
export OCM_URL="https://api.stage.openshift.com"

# RHOBS
export RHOBS_API_URL="https://observatorium-api.stage.openshift.com"
export RHOBS_TENANT="<your-tenant>"  # Get from team
export RHOBS_TOKEN="<token>"         # Get from Vault (https://github.com/openshift/ops-sop/blob/master/services/vault.md)

# Grafana
export GRAFANA_URL="https://grafana.stage.openshift.com"
export GRAFANA_TOKEN="<token>"       # Or use SSO login

# Test Configuration
export TEST_SECTOR="srep-dev"        # Dedicated E2E test sector
export TEST_REGION="us-east-1"       # AWS region for test clusters
```

### Verify `srep-dev` Sector Configuration

Before running tests, verify the srep-dev sector is properly configured:

#### 1\. Identify Management Cluster

```shell
# List management clusters in $TEST_SECTOR
ocm get /api/osd_fleet_mgmt/v1/management_clusters | jq --arg TEST_SECTOR "$TEST_SECTOR" -r '.items[] | select(.sector == $TEST_SECTOR) | [.cluster_management_reference.cluster_id, .name] | @tsv'

# Select the ID (first column) of the desired MC from list (if more than one) and set it below
export MC_CLUSTER_ID="<mc-cluster-id>"

# Login to MC
ocm backplane login $MC_CLUSTER_ID

# Verify you're on the right cluster
oc get hostedclusters -A
```

#### 2\. Verify route-monitor-operator Deployment

```shell
# Check RMO is deployed
oc get deployment -A | grep route-monitor-operator

# Check RMO logs
RMO_NS=$(oc get deployment -A | grep route-monitor-operator | awk '{print $1}')
oc logs -n $RMO_NS deployment/route-monitor-operator -f

# Verify RMO configuration
oc get deployment -n $RMO_NS route-monitor-operator -o yaml | \
  grep -A5 "probe-api-url"
# Expected: Should point to RHOBS synthetics API
```

#### 3\. Verify RHOBS Cell Connectivity

```shell
# Test API connectivity from your workstation
curl -H "Authorization: Bearer $RHOBS_TOKEN" \
  "$RHOBS_API_URL/api/metrics/v1/$RHOBS_TENANT/probes" | jq

# Expected: Should return JSON with list of probes (may be empty)
```

#### 4\. Verify Metrics Remote Write Configuration

```shell
# Check cluster monitoring config on MC
oc get configmap cluster-monitoring-config -n openshift-monitoring -o yaml

# Look for remoteWrite configuration
# Expected: Should see Thanos remote write endpoint
```

---

## Test Scenarios

---

## Test 1: Metrics Flow End-to-End

**Objective:** Verify that metrics from an HCP cluster flow through RHOBS to Thanos and are queryable in Grafana.

**Duration:** \~30-45 minutes

**Prerequisites:**

- OCM access to stage  
- Access to `srep-dev` sector  
- RHOBS/Grafana credentials

---

### 1.1 Create Test HCP Cluster

**Step 1.1.1: Generate Unique Cluster Name**

```shell
# Create unique test cluster name
export TEST_CLUSTER_NAME="metrics-test-$(date +%s)"
echo "Creating cluster: $TEST_CLUSTER_NAME"
```

**Step 1.1.2: Create ROSA HCP Cluster in srep-dev Sector**

```shell
# Create cluster targeting srep-dev sector
rosa create cluster \
  --cluster-name=$TEST_CLUSTER_NAME \
  --sts \
  --mode=auto \
  --hosted-cp \
  --properties=provision-shard-id:srep-dev \
  --region=$TEST_REGION \
  --yes

# Expected output:
# I: Creating cluster '<cluster-name>'
# I: To view a list of clusters and their status, run 'rosa list clusters'
```

**Step 1.1.3: Wait for Cluster to be Ready**

```shell
## Monitor cluster creation (takes ~10-15 minutes)
#watch -n 30 rosa describe cluster -c $TEST_CLUSTER_NAME
#
## Wait for state to show "ready"
## Alternative: Poll until ready
#while true; do
#  STATE=$(rosa describe cluster -c $TEST_CLUSTER_NAME -o json | jq -r '.state')
#  echo "Cluster state: $STATE"
#  if [ "$STATE" = "ready" ]; then
#    echo "Cluster is ready!"
#    break
#  fi
#  sleep 30
#done

rosa logs install --cluster $TEST_CLUSTER_NAME -f
```

**Step 1.1.4: Capture Cluster Details**

```shell
# Get cluster ID and API URL
export CLUSTER_ID=$(rosa describe cluster -c $TEST_CLUSTER_NAME -o json | jq -r '.id')
export CLUSTER_API_URL=$(rosa describe cluster -c $TEST_CLUSTER_NAME -o json | jq -r '.api.url')

echo "Cluster ID: $CLUSTER_ID"
echo "API URL: $CLUSTER_API_URL"

# Save details for later reference
cat > /tmp/test-cluster-details.txt <<EOF
Test Cluster Details
--------------------
Name: $TEST_CLUSTER_NAME
ID: $CLUSTER_ID
API URL: $CLUSTER_API_URL
Created: $(date)
EOF
```

**Expected Result:**

- ✅ Cluster created successfully in srep-dev sector  
- ✅ Cluster reaches "ready" state  
- ✅ Cluster ID and API URL obtained

**Common Failures:**

- ❌ "Insufficient capacity": ~~`srep-dev` sector~~   the management cluster may be at quota  
  - **Action:** Contact fleet-manager team to increase quota or clean up old clusters  
- ❌ "Invalid provision shard": Typo in sector name or sector doesn't exist  
  - **Action:** Verify sector name with `ocm list provision-shards`

---

### ~~1.2 Deploy Metrics-Generating Workload~~

**~~Step 1.2.1: Login to Test Cluster~~**

```shell
# Get admin credentials
rosa create admin -c $TEST_CLUSTER_NAME

# Login with provided credentials
oc login $CLUSTER_API_URL --username cluster-admin --password <password-from-above>

# Verify access
oc whoami
oc get nodes
```

**~~Step 1.2.2: Create Test Namespace~~**

```shell
# Create dedicated namespace for test workload
oc create namespace metrics-test

# Set as current namespace
oc project metrics-test
```

**~~Step 1.2.3: Deploy Metrics Generator~~**

```shell
# Deploy simple web app that generates metrics
cat <<EOF | oc apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: metrics-generator
  namespace: metrics-test
  labels:
    app: metrics-generator
spec:
  replicas: 1
  selector:
    matchLabels:
      app: metrics-generator
  template:
    metadata:
      labels:
        app: metrics-generator
    spec:
      containers:
      - name: nginx
        image: nginxinc/nginx-unprivileged:latest
        ports:
        - containerPort: 8080
          name: http
---
apiVersion: v1
kind: Service
metadata:
  name: metrics-generator
  namespace: metrics-test
  labels:
    app: metrics-generator
spec:
  selector:
    app: metrics-generator
  ports:
  - port: 8080
    targetPort: 8080
    name: http
---
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: metrics-generator
  namespace: metrics-test
  labels:
    app: metrics-generator
spec:
  selector:
    matchLabels:
      app: metrics-generator
  endpoints:
  - port: http
    interval: 30s
EOF
```

**~~Step 1.2.4: Verify Deployment~~**

```shell
# Wait for pod to be running
oc get pods -n metrics-test -w

# Verify ServiceMonitor is created
oc get servicemonitor -n metrics-test

# Check that Prometheus is scraping the metrics
oc exec -n openshift-monitoring prometheus-k8s-0 -- \
  promtool query instant http://localhost:9090 \
  'up{namespace="metrics-test",job="metrics-generator"}'
```

**~~Expected Result:~~**

- ~~✅ Pod running in metrics-test namespace~~  
- ~~✅ ServiceMonitor created~~  
- ~~✅ Metrics visible in cluster Prometheus~~

**~~Common Failures:~~**

- ~~❌ Pod stuck in "Pending": Insufficient resources~~  
  - **~~Action:~~** ~~Check node capacity with `oc get nodes` and `oc describe node`~~  
- ~~❌ ServiceMonitor not scraping: Label mismatch~~  
  - **~~Action:~~** ~~Verify labels match between Service and ServiceMonitor~~

---

### ~~1.3 Verify Metrics in RHOBS/Thanos~~

**~~Step 1.3.1: Wait for Metrics to Propagate~~**

```shell
# Remote write typically takes 2-5 minutes
echo "Waiting 5 minutes for metrics to propagate to RHOBS..."
sleep 300
```

**~~Step 1.3.2: Query Thanos Directly~~**

```shell
# Query for our test workload metrics
curl -G -H "Authorization: Bearer $RHOBS_TOKEN" \
  --data-urlencode "query=up{namespace=\"metrics-test\",cluster=\"$CLUSTER_ID\"}" \
  "$RHOBS_API_URL/api/metrics/v1/$RHOBS_TENANT/api/v1/query" | jq

# Expected output:
# {
#   "status": "success",
#   "data": {
#     "resultType": "vector",
#     "result": [
#       {
#         "metric": {
#           "__name__": "up",
#           "cluster": "<cluster-id>",
#           "namespace": "metrics-test",
#           "job": "metrics-generator"
#         },
#         "value": [<timestamp>, "1"]
#       }
#     ]
#   }
# }
```

**~~Step 1.3.3: Query for Multiple Metric Types~~**

```shell
# Check cluster-level metrics also made it
curl -G -H "Authorization: Bearer $RHOBS_TOKEN" \
  --data-urlencode "query=cluster_version{cluster=\"$CLUSTER_ID\"}" \
  "$RHOBS_API_URL/api/metrics/v1/$RHOBS_TENANT/api/v1/query" | jq

# Check API server metrics
curl -G -H "Authorization: Bearer $RHOBS_TOKEN" \
  --data-urlencode "query=apiserver_request_total{cluster=\"$CLUSTER_ID\"}" \
  "$RHOBS_API_URL/api/metrics/v1/$RHOBS_TENANT/api/v1/query" | jq '.data.result | length'

# Expected: Should return multiple results
```

**~~Step 1.3.4: Verify Metric Labels~~**

```shell
# Get full metric details
curl -G -H "Authorization: Bearer $RHOBS_TOKEN" \
  --data-urlencode "query=up{namespace=\"metrics-test\"}" \
  "$RHOBS_API_URL/api/metrics/v1/$RHOBS_TENANT/api/v1/query" | \
  jq '.data.result[0].metric'

# Verify required labels are present:
# - cluster: <cluster-id>
# - namespace: metrics-test
# - job: metrics-generator
# - pod: metrics-generator-<pod-id>
```

**~~Expected Result:~~**

- ~~✅ Metrics appear in Thanos within 5-10 minutes~~  
- ~~✅ Metrics include correct cluster\_id label~~  
- ~~✅ Both workload and cluster metrics are present~~  
- ~~✅ Labels are correctly propagated~~

**~~Common Failures:~~**

- ~~❌ No metrics found after 10 minutes:~~  
  - **~~Check 1:~~** ~~Verify remote write config on cluster:~~

```shell
oc get secret -n openshift-monitoring | grep remote-write
```

  - **~~Check 2:~~** ~~Check remote write errors in Prometheus logs:~~

```shell
oc logs -n openshift-monitoring prometheus-k8s-0 | grep -i "remote write"
```

  - **~~Check 3:~~** ~~Verify RHOBS Thanos receiver is healthy~~  
- ~~❌ Metrics present but missing labels:~~  
  - **~~Action:~~** ~~Check cluster monitoring config for external labels~~

---

### ~~1.4 Verify Metrics in Grafana~~

**~~Step 1.4.1: Login to Grafana~~**

```shell
# Open Grafana in browser
open $GRAFANA_URL

# Or use API if token is available
curl -H "Authorization: Bearer $GRAFANA_TOKEN" \
  "$GRAFANA_URL/api/datasources" | jq
```

**~~Step 1.4.2: Query Metrics via Grafana~~**

**~~Option A: Web UI~~**

1. ~~Navigate to Explore view~~  
2. ~~Select RHOBS/Thanos datasource~~  
3. ~~Enter query: `up{namespace="metrics-test",cluster="$CLUSTER_ID"}`~~  
4. ~~Click "Run Query"~~  
5. ~~Verify results appear with graph~~

**~~Option B: Grafana API~~**

```shell
# Query via API
curl -H "Authorization: Bearer $GRAFANA_TOKEN" \
  -H "Content-Type: application/json" \
  -X POST "$GRAFANA_URL/api/ds/query" \
  -d '{
    "queries": [
      {
        "refId": "A",
        "expr": "up{namespace=\"metrics-test\",cluster=\"'"$CLUSTER_ID"'\"}",
        "datasourceId": 1
      }
    ],
    "from": "now-5m",
    "to": "now"
  }' | jq
```

**~~Step 1.4.3: Create Test Dashboard~~**

```shell
# (Optional) Create a simple dashboard to visualize test metrics
# This can be done via Grafana UI:
# 1. Create New Dashboard
# 2. Add Panel
# 3. Query: up{namespace="metrics-test",cluster="$CLUSTER_ID"}
# 4. Save dashboard as "E2E Test Metrics"
```

**~~Expected Result:~~**

- ~~✅ Metrics queryable in Grafana Explore~~  
- ~~✅ Graph displays data points~~  
- ~~✅ Time series shows expected values (1 for "up" metric)~~  
- ~~✅ Dashboard can be created (optional)~~

**~~Common Failures:~~**

- ~~❌ "Datasource not found":~~  
  - **~~Action:~~** ~~Verify correct datasource is configured in Grafana~~  
- ~~❌ "No data" but Thanos query works:~~  
  - **~~Action:~~** ~~Check Grafana datasource configuration points to correct Thanos endpoint~~

---

### ~~1.5 Verify Metrics Timeline~~

**~~Step 1.5.1: Check Metric History~~**

```shell
# Query for last 30 minutes of data
curl -G -H "Authorization: Bearer $RHOBS_TOKEN" \
  --data-urlencode "query=up{namespace=\"metrics-test\",cluster=\"$CLUSTER_ID\"}" \
  --data-urlencode "start=$(date -u -d '30 minutes ago' +%s)" \
  --data-urlencode "end=$(date -u +%s)" \
  --data-urlencode "step=60" \
  "$RHOBS_API_URL/api/metrics/v1/$RHOBS_TENANT/api/v1/query_range" | \
  jq '.data.result[0].values | length'

# Expected: Should return ~30 data points (one per minute)
```

**~~Step 1.5.2: Verify Metric Retention~~**

```shell
# Query oldest available data
curl -G -H "Authorization: Bearer $RHOBS_TOKEN" \
  --data-urlencode "query=up{namespace=\"metrics-test\",cluster=\"$CLUSTER_ID\"}" \
  --data-urlencode "start=$(date -u -d '1 hour ago' +%s)" \
  --data-urlencode "end=$(date -u +%s)" \
  --data-urlencode "step=300" \
  "$RHOBS_API_URL/api/metrics/v1/$RHOBS_TENANT/api/v1/query_range" | \
  jq '.data.result[0].values | .[0]'

# Note the timestamp of first data point
```

**~~Expected Result:~~**

- ~~✅ Continuous time series data available~~  
- ~~✅ No gaps in metric timeline~~  
- ~~✅ Data retention meets expectations~~

---

### ~~Test 1 Summary~~

**~~Pass Criteria:~~**

- [ ] ~~HCP cluster created in srep-dev sector~~  
- [ ] ~~Metrics-generating workload deployed successfully~~  
- [ ] ~~Metrics appear in RHOBS/Thanos within 10 minutes~~  
- [ ] ~~Metrics include correct labels (cluster\_id, namespace)~~  
- [ ] ~~Metrics queryable via Grafana~~  
- [ ] ~~Continuous time series data with no gaps~~

**~~Failure Scenarios & Actions:~~**

| ~~Failure~~ | ~~Likely Cause~~ | ~~Debugging Steps~~ |
| :---- | :---- | :---- |
| ~~Metrics never appear in RHOBS~~ | ~~Remote write config issue~~ | ~~Check cluster monitoring config, verify Thanos endpoint reachable~~ |
| ~~Metrics missing labels~~ | ~~External labels not configured~~ | ~~Review cluster monitoring ConfigMap for externalLabels~~ |
| ~~Gaps in time series~~ | ~~Remote write failures~~ | ~~Check Prometheus logs for remote write errors~~ |
| ~~Grafana can't query~~ | ~~Datasource misconfigured~~ | ~~Verify datasource points to correct Thanos endpoint~~ |

---

## Test 2: Logs Flow End-to-End

**Objective:** Verify that logs (application and API audit logs) from an HCP cluster flow to Loki and are queryable.

**Duration:** \~30-45 minutes

**Prerequisites:**

- Newly created HCP cluster  
- Loki API access  
- Grafana access

---

### 2.1 Generate Test Logs

**Step 2.1.1: Create Unique Log Marker**

```shell
# Generate unique log marker for easy searching
export LOG_MARKER="e2e-test-$(date +%s)-$(uuidgen | cut -d- -f1)"
echo "Log marker: $LOG_MARKER"
```

**Step 2.1.2: Deploy Log-Generating Pod**

```shell
# Deploy pod that writes logs continuously
cat <<EOF | oc apply -f -
apiVersion: v1
kind: Pod
metadata:
  name: log-generator
  namespace: metrics-test
  labels:
    app: log-generator
spec:
  containers:
  - name: logger
    image: busybox:latest
    command:
    - /bin/sh
    - -c
    - |
      while true; do
        echo "INFO: $LOG_MARKER - Test log entry \$(date)"
        echo "WARN: $LOG_MARKER - Warning message \$(date)"
        echo "ERROR: $LOG_MARKER - Error message \$(date)"
        sleep 10
      done
  restartPolicy: Always
EOF
```

**Step 2.1.3: Verify Pod is Running and Generating Logs**

```shell
# Wait for pod to be ready
oc wait --for=condition=Ready pod/log-generator -n metrics-test --timeout=60s

# Verify logs are being generated locally
oc logs -n metrics-test log-generator --tail=10

# Expected: Should see logs with our LOG_MARKER
```

**Expected Result:**

- ✅ Pod running and generating logs  
- ✅ Logs visible via `oc logs`  
- ✅ Log marker appears in output

---

### 2.2 Generate API Audit Logs

**Step 2.2.1: Trigger API Operations**

```shell
# Perform several API operations to generate audit logs
oc create namespace audit-test-$(date +%s)
oc get nodes
oc get pods -A --limit=10
oc delete namespace audit-test-*

# These operations will generate audit log entries
```

**Step 2.2.2: Verify Audit Logging is Enabled**

```shell
# Check API server audit configuration on hosted cluster
# Note: For HCP, audit logs are automatically forwarded
# Verify by checking cluster description
rosa describe cluster -c $TEST_CLUSTER_NAME | grep -i audit
```

**Expected Result:**

- ✅ API operations executed successfully  
- ✅ Audit logging confirmed enabled

---

### 2.3 Verify Logs in Loki

**Step 2.3.1: Wait for Log Propagation**

```shell
# Log forwarding typically takes 2-5 minutes
echo "Waiting 5 minutes for logs to propagate to Loki..."
sleep 300
```

**Step 2.3.2: Query Application Logs**

```shell
# Query Loki for our test logs
curl -G -H "Authorization: Bearer $LOKI_TOKEN" \
  --data-urlencode "query={cluster_id=\"$CLUSTER_ID\",namespace=\"metrics-test\"} |= \"$LOG_MARKER\"" \
  --data-urlencode "limit=100" \
  "$LOKI_API_URL/loki/api/v1/query_range" | jq

# Expected output:
# {
#   "status": "success",
#   "data": {
#     "resultType": "streams",
#     "result": [
#       {
#         "stream": {
#           "cluster_id": "<cluster-id>",
#           "namespace": "metrics-test",
#           "pod": "log-generator",
#           "container": "logger"
#         },
#         "values": [
#           [<timestamp>, "INFO: <marker> - Test log entry ..."],
#           [<timestamp>, "WARN: <marker> - Warning message ..."]
#         ]
#       }
#     ]
#   }
# }
```

**Step 2.3.3: Verify Log Labels**

```shell
# Get full log entry with labels
curl -G -H "Authorization: Bearer $LOKI_TOKEN" \
  --data-urlencode "query={cluster_id=\"$CLUSTER_ID\"} |= \"$LOG_MARKER\"" \
  --data-urlencode "limit=1" \
  "$LOKI_API_URL/loki/api/v1/query_range" | \
  jq '.data.result[0].stream'

# Verify required labels are present:
# - cluster_id: <cluster-id>
# - namespace: metrics-test
# - pod: log-generator
# - container: logger
```

**Step 2.3.4: Query API Audit Logs**

```shell
# Query for audit logs from our API operations
curl -G -H "Authorization: Bearer $LOKI_TOKEN" \
  --data-urlencode "query={cluster_id=\"$CLUSTER_ID\",log_type=\"audit\"}" \
  --data-urlencode "limit=10" \
  "$LOKI_API_URL/loki/api/v1/query_range" | jq

# Check for specific operations we performed
curl -G -H "Authorization: Bearer $LOKI_TOKEN" \
  --data-urlencode "query={cluster_id=\"$CLUSTER_ID\",log_type=\"audit\"} |~ \"create.*namespace\"" \
  "$LOKI_API_URL/loki/api/v1/query_range" | jq '.data.result | length'

# Expected: Should find audit log entries
```

**Expected Result:**

- ✅ Application logs appear in Loki within 5-10 minutes  
- ✅ Logs include correct labels (cluster\_id, namespace, pod)  
- ✅ API audit logs are present  
- ✅ Logs are searchable and filterable

**Common Failures:**

- ❌ No logs found after 10 minutes:  
  - **Check 1:** Verify log forwarding is configured:

```shell
oc get clusterlogging -A
```

  - **Check 2:** Check cluster-logging-operator status  
  - **Check 3:** Verify Loki endpoint is reachable from cluster  
- ❌ Application logs present but no audit logs:  
  - **Check:** Verify audit log forwarding is enabled for HCP clusters in stage  
  - **Action:** Consult with team about HCP audit log configuration

---

### 2.4 Verify Logs in Grafana

**Step 2.4.1: Query Logs via Grafana Explore**

**Option A: Web UI**

1. Navigate to Grafana Explore  
2. Select Loki datasource  
3. Enter LogQL query: `{cluster_id="$CLUSTER_ID",namespace="metrics-test"} |= "$LOG_MARKER"`  
4. Click "Run Query"  
5. Verify log entries appear

**Option B: Grafana API**

```shell
# Query logs via Grafana API
curl -H "Authorization: Bearer $GRAFANA_TOKEN" \
  -H "Content-Type: application/json" \
  -X POST "$GRAFANA_URL/api/ds/query" \
  -d '{
    "queries": [
      {
        "refId": "A",
        "expr": "{cluster_id=\"'"$CLUSTER_ID"'\",namespace=\"metrics-test\"} |= \"'"$LOG_MARKER"'\"",
        "datasourceId": 2
      }
    ],
    "from": "now-30m",
    "to": "now"
  }' | jq
```

**Step 2.4.2: Test Log Filtering**

```shell
# Test various LogQL filters
# By log level
{cluster_id="$CLUSTER_ID"} |= "$LOG_MARKER" | json | level="ERROR"

# By time range
{cluster_id="$CLUSTER_ID"} |= "$LOG_MARKER" [5m]

# By pattern matching
{cluster_id="$CLUSTER_ID"} |~ "ERROR.*$LOG_MARKER"
```

**Expected Result:**

- ✅ Logs queryable in Grafana Explore  
- ✅ Log entries display with timestamps  
- ✅ Log filtering works (by level, time, pattern)  
- ✅ Log volume visualization appears

---

### 2.5 Verify Log Timeline

**Step 2.5.1: Check Log Continuity**

```shell
# Query for last 30 minutes and count entries
curl -G -H "Authorization: Bearer $LOKI_TOKEN" \
  --data-urlencode "query=count_over_time({cluster_id=\"$CLUSTER_ID\",namespace=\"metrics-test\"} |= \"$LOG_MARKER\" [30m])" \
  "$LOKI_API_URL/loki/api/v1/query" | jq

# Expected: Should show continuous log stream with no gaps
# (3 logs every 10 seconds = ~18 logs per minute = ~540 logs in 30 min)
```

**Step 2.5.2: Verify Log Retention**

```shell
# Query oldest available logs
curl -G -H "Authorization: Bearer $LOKI_TOKEN" \
  --data-urlencode "query={cluster_id=\"$CLUSTER_ID\"}" \
  --data-urlencode "start=$(date -u -d '1 hour ago' +%s)000000000" \
  --data-urlencode "end=$(date -u +%s)000000000" \
  --data-urlencode "limit=1" \
  "$LOKI_API_URL/loki/api/v1/query_range" | \
  jq '.data.result[0].values[0][0]'

# Note the timestamp of first log entry
```

**Expected Result:**

- ✅ Continuous log stream with expected frequency  
- ✅ No gaps in log timeline  
- ✅ Logs retained as per retention policy

---

### Test 2 Summary

**Pass Criteria:**

- [ ] Log-generating workload deployed successfully  
- [ ] Application logs appear in Loki within 10 minutes  
- [ ] API audit logs appear in Loki  
- [ ] Logs include correct labels (cluster\_id, namespace, pod)  
- [ ] Logs queryable via Grafana with LogQL  
- [ ] Continuous log stream with no gaps

**Failure Scenarios & Actions:**

| Failure | Likely Cause | Debugging Steps |
| :---- | :---- | :---- |
| No logs in Loki | Log forwarding not configured | Check ClusterLogging CR, verify fluentd/vector pods running |
| Logs missing labels | Label configuration issue | Review log forwarding pipeline config |
| No audit logs | Audit forwarding disabled | Verify HCP audit log configuration in stage |
| Grafana can't query logs | Loki datasource issue | Check Grafana datasource configuration |

---

## Test 3: Synthetic Monitoring Flow End-to-End

**Objective:** Verify that the complete synthetic monitoring flow works: RMO detects HCP cluster → registers probe in API → agent picks up probe → blackbox executes probe → metrics appear in RHOBS.

**Duration:** \~30-45 minutes

**Prerequisites:**

- HCP cluster from Test 1 (or create new one)  
- Access to Management Cluster  
- RHOBS API access

---

### 3.1 Verify route-monitor-operator Configuration

**Step 3.1.1: Login to Management Cluster**

```shell
# Get MC cluster ID for srep-dev sector
export MC_CLUSTER_ID="<mc-cluster-id>"  # From setup

# Login to MC
ocm backplane login $MC_CLUSTER_ID

# Verify connection
oc whoami
```

**Step 3.1.2: Check RMO Deployment**

```shell
# Find RMO namespace
export RMO_NS=$(oc get deployment -A | grep route-monitor-operator | awk '{print $1}')
echo "RMO Namespace: $RMO_NS"

# Check RMO deployment
oc get deployment -n $RMO_NS route-monitor-operator -o yaml

# Verify probe-api-url is configured
oc get deployment -n $RMO_NS route-monitor-operator -o yaml | \
  grep -A3 "probe-api-url"

# Expected: Should see RHOBS synthetics API URL
```

**Step 3.1.3: Check RMO Logs**

```shell
# Get recent RMO logs
oc logs -n $RMO_NS deployment/route-monitor-operator --tail=50

# Look for errors or API connection issues
# Expected: No errors, successful API connections
```

**Expected Result:**

- ✅ RMO deployed and running  
- ✅ probe-api-url configured correctly  
- ✅ No errors in RMO logs

**Common Failures:**

- ❌ RMO not deployed:  
  - **Action:** Check app-interface configuration for MC  
- ❌ probe-api-url not set:  
  - **Action:** Update RMO deployment config via app-interface

---

### 3.2 Verify HostedCluster Detection

**Step 3.2.1: Check HostedCluster CR Exists**

```shell
# List all HostedClusters
oc get hostedclusters -A

# Find our test cluster
oc get hostedclusters -A | grep $TEST_CLUSTER_NAME

# Get full HostedCluster details
export HC_NAMESPACE=$(oc get hostedclusters -A | grep $TEST_CLUSTER_NAME | awk '{print $1}')
oc get hostedcluster -n $HC_NAMESPACE $TEST_CLUSTER_NAME -o yaml
```

**Step 3.2.2: Verify API URL in HostedCluster**

```shell
# Extract API URL from HostedCluster
export HC_API_URL=$(oc get hostedcluster -n $HC_NAMESPACE $TEST_CLUSTER_NAME \
  -o jsonpath='{.status.controlPlaneEndpoint.host}')

echo "HostedCluster API URL: $HC_API_URL"
echo "Expected API URL: $CLUSTER_API_URL"

# These should match (or HC_API_URL should be the host portion)
```

**Step 3.2.3: Monitor RMO Watching HostedCluster**

```shell
# Watch RMO logs for HostedCluster detection
oc logs -n $RMO_NS deployment/route-monitor-operator -f | \
  grep -i "hostedcluster\|$TEST_CLUSTER_NAME"

# Expected: Should see log entries about detecting the HostedCluster
# Example: "Detected new HostedCluster: <cluster-name>"
```

**Expected Result:**

- ✅ HostedCluster CR exists for test cluster  
- ✅ API URL is correctly set in HostedCluster status  
- ✅ RMO logs show HostedCluster detection

---

### 3.3 Verify Probe Registration in RHOBS API

**Step 3.3.1: Wait for Probe Registration**

```shell
# RMO typically registers probe within 1-2 minutes
echo "Waiting 2 minutes for probe registration..."
sleep 120
```

**Step 3.3.2: Query RHOBS API for Probe**

```shell
# List all probes
curl -H "Authorization: Bearer $RHOBS_TOKEN" \
  "$RHOBS_API_URL/api/metrics/v1/$RHOBS_TENANT/probes" | \
  jq '.probes[] | select(.labels.cluster_id == "'"$CLUSTER_ID"'")'

# Expected output:
# {
#   "id": "<probe-id>",
#   "static_url": "https://api.<cluster-domain>:6443/readyz",
#   "labels": {
#     "cluster_id": "<cluster-id>",
#     "management_cluster_id": "<mc-id>",
#     "private": "false"
#   },
#   "status": "pending"
# }
```

**Step 3.3.3: Capture Probe ID**

```shell
# Save probe ID for later use
export PROBE_ID=$(curl -s -H "Authorization: Bearer $RHOBS_TOKEN" \
  "$RHOBS_API_URL/api/metrics/v1/$RHOBS_TENANT/probes" | \
  jq -r '.probes[] | select(.labels.cluster_id == "'"$CLUSTER_ID"'") | .id')

echo "Probe ID: $PROBE_ID"

# Verify it's not empty
if [ -z "$PROBE_ID" ]; then
  echo "ERROR: Probe not found!"
  exit 1
fi
```

**Step 3.3.4: Verify Probe Details**

```shell
# Get specific probe
curl -H "Authorization: Bearer $RHOBS_TOKEN" \
  "$RHOBS_API_URL/api/metrics/v1/$RHOBS_TENANT/probes/$PROBE_ID" | jq

# Verify:
# - static_url points to cluster API
# - cluster_id label matches our cluster
# - management_cluster_id is set
# - private label is present
```

**Expected Result:**

- ✅ Probe registered in RHOBS API within 2-5 minutes  
- ✅ Probe has correct static\_url (cluster API endpoint)  
- ✅ Probe labels include cluster\_id and management\_cluster\_id  
- ✅ Probe status is "pending" initially

**Common Failures:**

- ❌ No probe found after 5 minutes:  
  - **Check 1:** Verify RMO is watching HostedClusters:

```shell
oc logs -n $RMO_NS deployment/route-monitor-operator | grep -i watch
```

  - **Check 2:** Check RMO API client errors:

```shell
oc logs -n $RMO_NS deployment/route-monitor-operator | grep -i error
```

  - **Check 3:** Verify RHOBS API is reachable from MC:

```shell
oc debug node/<node> -- curl -v $RHOBS_API_URL
```

- ❌ Probe has wrong labels:  
  - **Action:** Check RMO configuration for label mapping

---

### 3.4 Verify rhobs-synthetics-agent Picks Up Probe

**Step 3.4.1: Login to RHOBS Cell Cluster**

```shell
# Get RHOBS cell cluster ID (ask team or check app-interface)
export RHOBS_CLUSTER_ID="<rhobs-cell-cluster-id>"

# Login to RHOBS cluster
ocm backplane login $RHOBS_CLUSTER_ID

# Find synthetics-agent namespace
export AGENT_NS=$(oc get deployment -A | grep synthetics-agent | awk '{print $1}')
echo "Agent Namespace: $AGENT_NS"
```

**Step 3.4.2: Check Agent is Running**

```shell
# Verify agent deployment
oc get deployment -n $AGENT_NS rhobs-synthetics-agent

# Check agent logs
oc logs -n $AGENT_NS deployment/rhobs-synthetics-agent --tail=50
```

**Step 3.4.3: Monitor Agent Polling**

```shell
# Watch agent logs for probe polling
oc logs -n $AGENT_NS deployment/rhobs-synthetics-agent -f | \
  grep -i "probe\|$CLUSTER_ID"

# Expected: Should see agent polling API and finding our probe
# Example logs:
# "Fetched 15 probes from API"
# "Found probe for cluster <cluster-id>"
# "Creating Probe CR for probe-<id>"
```

**Step 3.4.4: Verify Probe CR Created**

```shell
# Wait for agent to create Probe CR (typically 30-60 seconds)
sleep 60

# List Probe CRs
oc get probes -n $AGENT_NS

# Find our probe
oc get probe -n $AGENT_NS | grep $CLUSTER_ID

# Get full Probe CR
export PROBE_CR_NAME=$(oc get probe -n $AGENT_NS -o json | \
  jq -r '.items[] | select(.metadata.labels."rhobs.monitoring/cluster-id" == "'"$CLUSTER_ID"'") | .metadata.name')

echo "Probe CR Name: $PROBE_CR_NAME"

oc get probe -n $AGENT_NS $PROBE_CR_NAME -o yaml
```

**Expected Result:**

- ✅ Agent deployment running on RHOBS cell  
- ✅ Agent logs show API polling  
- ✅ Agent creates Probe CR for our cluster  
- ✅ Probe CR has correct target URL and labels

**Common Failures:**

- ❌ Agent not creating Probe CR:  
  - **Check 1:** Verify agent is polling the correct API endpoint  
  - **Check 2:** Check agent logs for URL validation failures:

```shell
oc logs -n $AGENT_NS deployment/rhobs-synthetics-agent | grep -i "validation\|failed"
```

  - **Check 3:** Verify label selector matches:

```shell
oc get deployment -n $AGENT_NS rhobs-synthetics-agent -o yaml | grep label-selector
```

---

### 3.5 Verify Probe Execution

**Step 3.5.1: Check Blackbox Exporter**

```shell
# Verify blackbox exporter is deployed
oc get deployment -n $AGENT_NS | grep blackbox

# Get blackbox exporter logs
oc logs -n $AGENT_NS deployment/synthetics-blackbox-prober -f | \
  grep -i "$CLUSTER_API_URL"

# Expected: Should see HTTP probes executing
```

**Step 3.5.2: Check Prometheus Scraping Probes**

```shell
# Find Prometheus instance
oc get prometheus -n $AGENT_NS

# Check if Prometheus is scraping our probe
oc exec -n $AGENT_NS prometheus-rhobs-synthetics-0 -- \
  promtool query instant http://localhost:9090 \
  "probe_success{cluster_id=\"$CLUSTER_ID\"}"

# Expected: Should return value of 1 (success) or 0 (failure)
```

**Step 3.5.3: Verify Probe Status Update**

```shell
# Check if agent updated probe status in API
curl -H "Authorization: Bearer $RHOBS_TOKEN" \
  "$RHOBS_API_URL/api/metrics/v1/$RHOBS_TENANT/probes/$PROBE_ID" | \
  jq '.status'

# Expected: "active" (if probe succeeded) or "failed" (if probe failed)
# Should no longer be "pending"
```

**Expected Result:**

- ✅ Blackbox exporter executing probes  
- ✅ Prometheus scraping probe metrics  
- ✅ probe\_success metric shows value  
- ✅ Probe status updated in API to "active"

---

### 3.6 Verify Probe Metrics in RHOBS

**Step 3.6.1: Wait for Metrics to Propagate**

```shell
# Remote write typically takes 2-5 minutes
echo "Waiting 5 minutes for probe metrics to propagate..."
sleep 300
```

**Step 3.6.2: Query Probe Metrics**

```shell
# Query for probe_success metric
curl -G -H "Authorization: Bearer $RHOBS_TOKEN" \
  --data-urlencode "query=probe_success{cluster_id=\"$CLUSTER_ID\"}" \
  "$RHOBS_API_URL/api/metrics/v1/$RHOBS_TENANT/api/v1/query" | jq

# Expected output:
# {
#   "status": "success",
#   "data": {
#     "result": [
#       {
#         "metric": {
#           "__name__": "probe_success",
#           "cluster_id": "<cluster-id>",
#           "probe_type": "blackbox",
#           ...
#         },
#         "value": [<timestamp>, "1"]
#       }
#     ]
#   }
# }
```

**Step 3.6.3: Query Additional Probe Metrics**

```shell
# Query probe duration
curl -G -H "Authorization: Bearer $RHOBS_TOKEN" \
  --data-urlencode "query=probe_duration_seconds{cluster_id=\"$CLUSTER_ID\"}" \
  "$RHOBS_API_URL/api/metrics/v1/$RHOBS_TENANT/api/v1/query" | jq

# Query HTTP-specific metrics
curl -G -H "Authorization: Bearer $RHOBS_TOKEN" \
  --data-urlencode "query=probe_http_status_code{cluster_id=\"$CLUSTER_ID\"}" \
  "$RHOBS_API_URL/api/metrics/v1/$RHOBS_TENANT/api/v1/query" | jq

# Expected: Should see HTTP 200 status code
```

**Step 3.6.4: Verify Metric Labels**

```shell
# Get full metric details
curl -G -H "Authorization: Bearer $RHOBS_TOKEN" \
  --data-urlencode "query=probe_success{cluster_id=\"$CLUSTER_ID\"}" \
  "$RHOBS_API_URL/api/metrics/v1/$RHOBS_TENANT/api/v1/query" | \
  jq '.data.result[0].metric'

# Verify required labels:
# - cluster_id: <cluster-id>
# - management_cluster_id: <mc-id>
# - apiserver_url: <cluster-api-url>
# - probe_type: blackbox or similar
```

**Expected Result:**

- ✅ probe\_success metric appears in RHOBS within 5-10 minutes  
- ✅ probe\_success value is 1 (indicating API is reachable)  
- ✅ Additional probe metrics present (duration, status\_code)  
- ✅ All required labels present

**Common Failures:**

- ❌ probe\_success \= 0:  
  - **Check:** Verify cluster API is actually reachable from RHOBS cell  
  - **Action:** Test manually from RHOBS cluster:

```shell
oc debug node/<node> -- curl -v https://<cluster-api-url>:6443/readyz
```

- ❌ No probe metrics in RHOBS:  
  - **Check 1:** Verify Prometheus remote write is configured  
  - **Check 2:** Check Prometheus logs for remote write errors

---

### 3.7 Verify Probe Lifecycle (Deletion)

**Step 3.7.1: Delete Test Cluster**

```shell
# Delete the HCP cluster
rosa delete cluster -c $TEST_CLUSTER_NAME --yes

# Monitor deletion
watch -n 30 rosa describe cluster -c $TEST_CLUSTER_NAME
```

**Step 3.7.2: Verify HostedCluster Deletion**

```shell
# Check if HostedCluster CR is deleted from MC
oc get hostedcluster -n $HC_NAMESPACE $TEST_CLUSTER_NAME

# Expected: Should be gone or in "Deleting" state
```

**Step 3.7.3: Verify Probe Deletion from API**

```shell
# Wait a few minutes for RMO to detect deletion
sleep 180

# Check if probe is deleted
curl -H "Authorization: Bearer $RHOBS_TOKEN" \
  "$RHOBS_API_URL/api/metrics/v1/$RHOBS_TENANT/probes/$PROBE_ID"

# Expected: Should return 404 Not Found or similar
```

**Step 3.7.4: Verify Probe CR Deletion from RHOBS Cluster**

```shell
# Check if Probe CR is deleted
oc get probe -n $AGENT_NS $PROBE_CR_NAME

# Expected: Should be gone
```

**Expected Result:**

- ✅ Cluster deletion initiated successfully  
- ✅ HostedCluster CR removed from MC  
- ✅ Probe deleted from RHOBS API  
- ✅ Probe CR removed from RHOBS cluster  
- ✅ Probe metrics stop updating in RHOBS

---

### Test 3 Summary

**Pass Criteria:**

- [ ] route-monitor-operator detects HostedCluster  
- [ ] Probe registered in RHOBS API with correct details  
- [ ] rhobs-synthetics-agent picks up probe  
- [ ] Probe CR created on RHOBS cell cluster  
- [ ] Blackbox exporter executes probe successfully  
- [ ] Probe metrics appear in RHOBS (probe\_success \= 1\)  
- [ ] Probe lifecycle works (deletion cleans up properly)

**Failure Scenarios & Actions:**

| Failure | Likely Cause | Debugging Steps |
| :---- | :---- | :---- |
| Probe never registered | RMO not watching or API unreachable | Check RMO logs, verify API connectivity from MC |
| Agent doesn't create Probe CR | Label selector mismatch or URL validation failure | Check agent config, verify target URL is reachable |
| probe\_success \= 0 | Cluster API not reachable from RHOBS cell | Test connectivity, check firewall rules |
| Probe not deleted | RMO not detecting cluster deletion | Check RMO watches, verify HostedCluster finalizers |

---

## Cleanup Procedures

### Complete Test Cleanup

After completing all tests, clean up resources to avoid costs and quota usage:

**Step 1: Delete Test Cluster**

```shell
# Delete cluster if still running
rosa delete cluster -c $TEST_CLUSTER_NAME --yes

# Wait for deletion to complete
watch -n 30 rosa list clusters | grep $TEST_CLUSTER_NAME
```

**Step 2: Verify Probe Cleanup**

```shell
# Verify probe is removed from API
curl -H "Authorization: Bearer $RHOBS_TOKEN" \
  "$RHOBS_API_URL/api/metrics/v1/$RHOBS_TENANT/probes" | \
  jq '.probes[] | select(.labels.cluster_id == "'"$CLUSTER_ID"'")'

# Expected: Empty result
```

**Step 3: Clean Up Any Orphaned Resources**

```shell
# Check for any remaining Probe CRs on RHOBS cluster
oc get probes -n $AGENT_NS | grep $CLUSTER_ID

# Manually delete if found
oc delete probe -n $AGENT_NS <probe-name>
```

**Step 4: Document Test Results**

```shell
# Save test results for reference
cat > /tmp/e2e-test-results-$(date +%Y%m%d-%H%M%S).txt <<EOF
RHOBS E2E Test Results
======================
Date: $(date)
Tester: $(whoami)

Test 1 - Metrics Flow:
- Status: [PASS/FAIL]
- Notes: <any notes>

Test 2 - Logs Flow:
- Status: [PASS/FAIL]
- Notes: <any notes>

Test 3 - Synthetics Flow:
- Status: [PASS/FAIL]
- Notes: <any notes>

Overall: [PASS/FAIL]
EOF
```

---

## Troubleshooting

### Common Issues Across All Tests

#### Issue: Authentication Failures

**Symptoms:**

- 401 Unauthorized errors from APIs  
- "token expired" messages

**Debugging:**

```shell
# Check OCM token expiry
ocm token --refresh

# Verify RHOBS token
curl -H "Authorization: Bearer $RHOBS_TOKEN" \
  "$RHOBS_API_URL/api/metrics/v1/$RHOBS_TENANT/probes"

# Re-authenticate if needed
ocm login --url=https://api.stage.openshift.com
```

**Resolution:**

- Refresh tokens  
- Verify credentials in Vault are current  
- Contact team if persistent issues

---

#### Issue: Cluster Creation Failures

**Symptoms:**

- Cluster stuck in "installing" state  
- "insufficient capacity" errors

**Debugging:**

```shell
# Check cluster status
rosa describe cluster -c $TEST_CLUSTER_NAME

# View installation logs
rosa logs install -c $TEST_CLUSTER_NAME --watch

# Check sector capacity
ocm list provision-shards | grep $TEST_SECTOR
```

**Resolution:**

- For capacity issues: Contact fleet-manager team  
- For installation failures: Check rosa logs and file bug if needed  
- Clean up old clusters in sector to free quota

---

#### Issue: Network Connectivity

**Symptoms:**

- Timeouts when querying APIs  
- Probes failing with "target unreachable"

**Debugging:**

```shell
# Test connectivity from local machine
curl -v -m 10 $RHOBS_API_URL

# Test from MC (for probe reachability)
oc debug node/<node> -- curl -v https://<cluster-api>:6443/readyz

# Test from RHOBS cell
oc debug node/<node> -- curl -v https://<cluster-api>:6443/readyz
```

**Resolution:**

- Check firewall rules between clusters  
- Verify VPC peering if applicable  
- Contact network team for persistent issues

---

#### Issue: Metrics/Logs Not Appearing

**Symptoms:**

- Data expected but not found in RHOBS/Loki  
- Queries return empty results

**Debugging:**

```shell
# Check remote write configuration
oc get secret -n openshift-monitoring | grep remote-write

# Check Prometheus remote write status
oc exec -n openshift-monitoring prometheus-k8s-0 -- \
  curl localhost:9090/metrics | grep remote_write

# Check log forwarding
oc get clusterlogging -A
oc get pods -n openshift-logging
```

**Resolution:**

- Verify remote write/log forwarding configuration in app-interface  
- Check receiver endpoints are healthy  
- Review Prometheus/fluentd logs for errors

---

### Getting Help

If you encounter issues not covered in this guide:

1. **Check existing documentation:**  
     
   - RMO README: `route-monitor-operator/README.md`  
   - Synthetics API README: `rhobs-synthetics-api/README.md`  
   - Synthetics Agent README: `rhobs-synthetics-agent/README.md`

   

2. **Search for similar issues:**  
     
   - Jira: `project = SREP AND text ~ "synthetics OR rhobs"`  
   - Slack: Search in \#forum-rhobs-core

   

3. **Ask for help:**  
     
   - Slack: \#forum-rhobs-core  
   - Tag: @rhobs-team  
   - Email: [rhobs-team@redhat.com](mailto:rhobs-team@redhat.com)

   

4. **File a bug:**  
     
   - If you discover a reproducible issue, file a Jira ticket  
   - Include: steps to reproduce, expected vs actual behavior, relevant logs

---

## Appendix

### A. Useful Commands Reference

```shell
# OCM
ocm login --url=https://api.stage.openshift.com
ocm whoami
ocm list clusters
ocm backplane login <cluster-id>

# ROSA
rosa create cluster --help
rosa describe cluster -c <name>
rosa logs install -c <name> --watch
rosa delete cluster -c <name>

# Kubernetes
oc get hostedclusters -A
oc logs -n <namespace> <pod> -f
oc debug node/<node> -- <command>

# Queries
# Prometheus/Thanos
curl -G --data-urlencode "query=<promql>" <thanos-url>/api/v1/query

# Loki
curl -G --data-urlencode "query=<logql>" <loki-url>/loki/api/v1/query_range

# RHOBS Synthetics API
curl -H "Authorization: Bearer $TOKEN" <api-url>/api/metrics/v1/<tenant>/probes
```

### B. Expected Timelines

| Event | Typical Duration | Max Wait Time |
| :---- | :---- | :---- |
| Cluster creation | 10-15 minutes | 30 minutes |
| Metrics in RHOBS | 2-5 minutes | 10 minutes |
| Logs in Loki | 2-5 minutes | 10 minutes |
| Probe registration | 1-2 minutes | 5 minutes |
| Probe execution | 30-60 seconds | 2 minutes |
| Probe metrics in RHOBS | 2-5 minutes | 10 minutes |
| Cluster deletion | 5-10 minutes | 20 minutes |

### C. Validation Checklist

Print this checklist and check off items as you validate:

#### Test 1 \- Metrics

- [ ] Cluster created successfully  
- [ ] Workload deployed  
- [ ] Metrics in Thanos  
- [ ] Metrics in Grafana  
- [ ] Correct labels present

#### Test 2 \- Logs

- [ ] Logs generated  
- [ ] Application logs in Loki  
- [ ] Audit logs in Loki  
- [ ] Logs in Grafana  
- [ ] Correct labels present

#### Test 3 \- Synthetics

- [ ] RMO detects cluster  
- [ ] Probe registered in API  
- [ ] Agent creates Probe CR  
- [ ] Probe executes successfully  
- [ ] Probe metrics in RHOBS  
- [ ] Lifecycle works (deletion)

---

## Document History

| Version | Date | Author | Changes |
| :---- | :---- | :---- | :---- |
| 1.0 | 2026-02-17 | Trevor Nierman \+ Claude | Initial draft based on component docs and team discussion |

---

## Next Steps

**For Team Review:**

1. Review this manual test plan with team  
2. Validate steps by having 2-3 people execute independently  
3. Collect feedback and refine procedures  
4. Get sign-off from Matt and Dustin  
5. Use this as specification for automation (Phase 4\)

**Validation Process:**

- Assign 2-3 team members to execute this plan independently  
- Each person should document:  
  - Any unclear steps  
  - Missing information  
  - Actual vs expected results  
  - Suggestions for improvement  
- Schedule review meeting to discuss findings  
- Update plan based on feedback

**Once Validated:**

- Mark Phase 1 as complete in Jira  
- Begin Phase 2 (Repository decision)  
- Use this plan as foundation for automated tests in Phase 4

