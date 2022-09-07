package cluster

const devicePrometheusTmpl = `
apiVersion: v1
kind: Namespace
metadata:
  name: monitoring
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: prometheus
  labels:
    addonmanager.kubernetes.io/mode: Reconcile
    kubernetes.io/cluster-service: 'true'
rules:
  - verbs:
      - get
      - list
      - watch
    apiGroups:
      - ''
    resources:
      - nodes
      - nodes/metrics
      - services
      - endpoints
      - pods
  - verbs:
      - get
    apiGroups:
      - ''
    resources:
      - configmaps
  - verbs:
      - get
    nonResourceURLs:
      - /metrics
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: prometheus
  labels:
    addonmanager.kubernetes.io/mode: Reconcile
    kubernetes.io/cluster-service: 'true'
subjects:
  - kind: ServiceAccount
    name: prometheus
    namespace: monitoring
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: prometheus
---
kind: ConfigMap
apiVersion: v1
metadata:
  name: prometheus-config
  namespace: monitoring
data:
  prometheus.yml: |-
    global:
      scrape_interval: 5s
      evaluation_interval: 5s
    rule_files:
      - /etc/prometheus/prometheus.rules
    alerting:
      alertmanagers:
      - scheme: http
        static_configs:
        - targets:
          - "alertmanager.monitoring.svc:9093"
    scrape_configs:
      - job_name: edge-command
        scrape_interval: 10s
        metrics_path: /api/v2/pmetrics
        static_configs:
        - targets:
          - edge-core-command.edge-system:59882
      - job_name: device-rfid-llrp
        scrape_interval: 10s
        metrics_path: /metrics
        static_configs:
        - targets:
          - device-rfid-llrp.edge-system:9090
      - job_name: device-opcua
        scrape_interval: 10s
        metrics_path: /metrics
        static_configs:
        - targets:
          - device-opcua.edge-system:9090
      - job_name: device-modbus
        scrape_interval: 10s
        metrics_path: /metrics
        static_configs:
        - targets:
          - device-modbus.edge-system:9090
      - job_name: device-onvif-camera
        scrape_interval: 10s
        metrics_path: /metrics
        static_configs:
        - targets:
          - device-onvif-camera.edge-system:9090
      - job_name: device-usb-camera
        scrape_interval: 10s
        metrics_path: /metrics
        static_configs:
        - targets:
          - device-usb-camera.edge-system:9090
      - job_name: device-rest
        scrape_interval: 10s
        metrics_path: /metrics
        static_configs:
        - targets:
          - device-rest.edge-system:9090
      - job_name: device-snmp
        scrape_interval: 10s
        metrics_path: /metrics
        static_configs:
        - targets:
          - device-snmp.edge-system:9090
      - job_name: device-mqtt
        scrape_interval: 10s
        metrics_path: /metrics
        static_configs:
        - targets:
          - device-mqtt.edge-system:9090
      - job_name: 'node-exporter'
        kubernetes_sd_configs:
          - role: endpoints
        relabel_configs:
        - source_labels: [__meta_kubernetes_endpoints_name]
          regex: 'node-exporter'
          action: keep
      
      - job_name: 'kubernetes-apiservers'
        kubernetes_sd_configs:
        - role: endpoints
        scheme: https
        tls_config:
          ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
        bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token
        relabel_configs:
        - source_labels: [__meta_kubernetes_namespace, __meta_kubernetes_service_name, __meta_kubernetes_endpoint_port_name]
          action: keep
          regex: default;kubernetes;https
      - job_name: 'kubernetes-nodes'
        scheme: https
        tls_config:
          ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
        bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token
        kubernetes_sd_configs:
        - role: node
        relabel_configs:
        - action: labelmap
          regex: __meta_kubernetes_node_label_(.+)
        - target_label: __address__
          replacement: kubernetes.default.svc:443
        - source_labels: [__meta_kubernetes_node_name]
          regex: (.+)
          target_label: __metrics_path__
          replacement: /api/v1/nodes/${1}/proxy/metrics     
      
      - job_name: 'kubernetes-pods'
        kubernetes_sd_configs:
        - role: pod
        relabel_configs:
        - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
          action: keep
          regex: true
        - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_path]
          action: replace
          target_label: __metrics_path__
          regex: (.+)
        - source_labels: [__address__, __meta_kubernetes_pod_annotation_prometheus_io_port]
          action: replace
          regex: ([^:]+)(?::\d+)?;(\d+)
          replacement: $1:$2
          target_label: __address__
        - action: labelmap
          regex: __meta_kubernetes_pod_label_(.+)
        - source_labels: [__meta_kubernetes_namespace]
          action: replace
          target_label: kubernetes_namespace
        - source_labels: [__meta_kubernetes_pod_name]
          action: replace
          target_label: kubernetes_pod_name
      
      - job_name: 'kube-state-metrics'
        static_configs:
          - targets: ['kube-state-metrics.kube-system.svc.cluster.local:8080']
      - job_name: 'kubernetes-cadvisor'
        scheme: https
        tls_config:
          ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
        bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token
        kubernetes_sd_configs:
        - role: node
        relabel_configs:
        - action: labelmap
          regex: __meta_kubernetes_node_label_(.+)
        - target_label: __address__
          replacement: kubernetes.default.svc:443
        - source_labels: [__meta_kubernetes_node_name]
          regex: (.+)
          target_label: __metrics_path__
          replacement: /api/v1/nodes/${1}/proxy/metrics/cadvisor
      
      - job_name: 'kubernetes-service-endpoints'
        kubernetes_sd_configs:
        - role: endpoints
        relabel_configs:
        - source_labels: [__meta_kubernetes_service_annotation_prometheus_io_scrape]
          action: keep
          regex: true
        - source_labels: [__meta_kubernetes_service_annotation_prometheus_io_scheme]
          action: replace
          target_label: __scheme__
          regex: (https?)
        - source_labels: [__meta_kubernetes_service_annotation_prometheus_io_path]
          action: replace
          target_label: __metrics_path__
          regex: (.+)
        - source_labels: [__address__, __meta_kubernetes_service_annotation_prometheus_io_port]
          action: replace
          target_label: __address__
          regex: ([^:]+)(?::\d+)?;(\d+)
          replacement: $1:$2
        - action: labelmap
          regex: __meta_kubernetes_service_label_(.+)
        - source_labels: [__meta_kubernetes_namespace]
          action: replace
          target_label: kubernetes_namespace
        - source_labels: [__meta_kubernetes_service_name]
          action: replace
          target_label: kubernetes_name
---
kind: ConfigMap
apiVersion: v1
metadata:
  name: prometheus-rules
  namespace: monitoring
  labels:
    addonmanager.kubernetes.io/mode: EnsureExists
    kubernetes.io/cluster-service: 'true'
data:
  general.rules: |
    groups:
    - name: general.rules
      rules:
      - alert: InstanceDown
        expr: up == 0
        for: 1m
        labels:
          severity: error
        annotations:
          summary: "Instance {<!-- -->{ $labels.instance }} 停止工作"
          description: "{<!-- -->{ $labels.instance }} job {<!-- -->{ $labels.job }} 已经停止5分钟以上."
  node.rules: |
    groups:
    - name: node.rules
      rules:
      - alert: NodeFilesystemUsage
        expr: 100 - (node_filesystem_free_bytes{fstype=~"ext4|xfs"} / node_filesystem_size_bytes{fstype=~"ext4|xfs"} * 100) &gt; 80
        for: 1m
        labels:
          severity: warning
        annotations:
          summary: "Instance {<!-- -->{ $labels.instance }} : {<!-- -->{ $labels.mountpoint }} 分区使用率过高"
          description: "{<!-- -->{ $labels.instance }}: {<!-- -->{ $labels.mountpoint }} 分区使用大于80% (当前值: {<!-- -->{ $value }})"
      - alert: NodeMemoryUsage
        expr: 100 - (node_memory_MemFree_bytes+node_memory_Cached_bytes+node_memory_Buffers_bytes) / node_memory_MemTotal_bytes * 100 &gt; 80
        for: 1m
        labels:
          severity: warning
        annotations:
          summary: "Instance {<!-- -->{ $labels.instance }} 内存使用率过高"
          description: "{<!-- -->{ $labels.instance }}内存使用大于80% (当前值: {<!-- -->{ $value }})"

      - alert: NodeCPUUsage
        expr: 100 - (avg(irate(node_cpu_seconds_total{mode="idle"}[5m])) by (instance) * 100) &gt; 60
        for: 1m
        labels:
          severity: warning
        annotations:
          summary: "Instance {<!-- -->{ $labels.instance }} CPU使用率过高"
          description: "{<!-- -->{ $labels.instance }}CPU使用大于60% (当前值: {<!-- -->{ $value }})"
---
kind: StatefulSet
apiVersion: apps/v1
metadata:
  name: prometheus
  namespace: monitoring
spec:
  replicas: 1
  selector:
    matchLabels:
      k8s-app: prometheus
  template:
    metadata:
      creationTimestamp: null
      labels:
        k8s-app: prometheus
      annotations:
        scheduler.alpha.kubernetes.io/critical-pod: ''
    spec:
      volumes:
        - name: config-volume
          configMap:
            name: prometheus-config
            defaultMode: 420
        - name: prometheus-rules
          configMap:
            name: prometheus-rules
            defaultMode: 420
      containers:
        - name: prometheus-server
          image: prom/prometheus:v2.37.0
          args:
            - '--config.file=/etc/prometheus/config/prometheus.yml'
            - '--storage.tsdb.path=/prometheus'
            - '--web.console.libraries=/etc/prometheus/console_libraries'
            - '--web.console.templates=/etc/prometheus/consoles'
            - '--web.enable-lifecycle'
          ports:
            - containerPort: 9090
              protocol: TCP
          resources:
            limits:
              cpu: 200m
              memory: 1000Mi
            requests:
              cpu: 200m
              memory: 1000Mi
          volumeMounts:
            - name: config-volume
              mountPath: /etc/prometheus/config
            - name: prometheus-data
              mountPath: /prometheus
            - name: prometheus-rules
              mountPath: /etc/prometheus/rules
          livenessProbe:
            httpGet:
              path: /-/healthy
              port: 9090
              scheme: HTTP
            initialDelaySeconds: 30
            timeoutSeconds: 30
            periodSeconds: 10
            successThreshold: 1
            failureThreshold: 3
          readinessProbe:
            httpGet:
              path: /-/ready
              port: 9090
              scheme: HTTP
            initialDelaySeconds: 30
            timeoutSeconds: 30
            periodSeconds: 10
            successThreshold: 1
            failureThreshold: 3
          terminationMessagePolicy: File
          imagePullPolicy: IfNotPresent
      restartPolicy: Always
  volumeClaimTemplates:
    - kind: PersistentVolumeClaim
      apiVersion: v1
      metadata:
        name: prometheus-data
      spec:
        accessModes:
          - ReadWriteOnce
        resources:
          requests:
            storage: 1Gi
        storageClassName: local-path
        volumeMode: Filesystem
  serviceName: prometheus
---
kind: Service
apiVersion: v1
metadata:
  name: prometheus
  namespace: monitoring
  labels:
    addonmanager.kubernetes.io/mode: Reconcile
    kubernetes.io/cluster-service: 'true'
    kubernetes.io/name: Prometheus
spec:
  ports:
    - name: http
      protocol: TCP
      port: 9090
      targetPort: 9090
  selector:
    k8s-app: prometheus
`
