package cluster

const deviceZmqTmpl = `
---
apiVersion: v1
kind: Namespace
metadata:
  name: edge-system
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: edge-role
  labels:
    app: edge-rbac
rules:
  - verbs:
      - get
      - list
      - create
      - delete
      - update
      - watch
    apiGroups:
      - ''
    resources:
      - configmaps
      - services
      - pods
  - verbs:
      - get
      - list
      - create
      - delete
      - update
      - watch
    apiGroups:
      - apps
    resources:
      - deployments
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: edge-rolebinding
  labels:
    app: edge-rbac
subjects:
  - kind: ServiceAccount
    name: default
    namespace: edge-system
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: edge-role
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: device-dashboard-claim
  namespace: edge-system
spec:
  storageClassName: local-path
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 30M
---
kind: ConfigMap
apiVersion: v1
metadata:
  name: edge-support-scheduler
  namespace: edge-system
data:
  configuration.toml: >
    ScheduleIntervalTime = 500


    [Writable]

    LogLevel = 'INFO'
        [Writable.InsecureSecrets]
            [Writable.InsecureSecrets.DB]
            path = "redisdb"
                [Writable.InsecureSecrets.DB.Secrets]
                username = ""
                password = ""

    [Service]

    HealthCheckInterval = '10s'

    Host = 'edge-support-scheduler'

    Port = 59861

    ServerBindAddr = '0.0.0.0' # Leave blank so default to Host value unless different value is needed.

    StartupMsg = 'This is the Support Scheduler Microservice'

    MaxResultCount = 50000

    MaxRequestSize = 0 # Not curently used. Defines the maximum size of http request body in bytes

    RequestTimeout = '5s'
    
      [Service.CORSConfiguration]

      EnableCORS = false
      CORSAllowCredentials = false
      CORSAllowedOrigin = "http://edge-support-scheduler"
      CORSAllowedMethods = "GET, POST, PUT, PATCH, DELETE"
      CORSAllowedHeaders = "Authorization, Accept, Accept-Language, Content-Language, Content-Type, X-Correlation-ID"
      CORSExposeHeaders = "Cache-Control, Content-Language, Content-Length, Content-Type, Expires, Last-Modified, Pragma, X-Correlation-ID"
      CORSMaxAge = 3600


    [Registry]

    Host = 'localhost'

    Port = 8500

    Type = 'consul'


    [Databases]
      [Databases.Primary]
      Host = 'edge-redis-ha-announce-0'
      Name = 'scheduler'
      Port = 6379
      Timeout = 5000
      Type = 'redisdb'

    [Intervals]
        [Intervals.Midnight]
        Name = 'midnight'
        Start = '20180101T000000'
        Interval = '24h'

    [IntervalActions]
        [IntervalActions.ScrubAged]
        Name = 'scrub-aged-events'
        Host = 'localhost'
        Port = 59880
        Protocol = 'http'
        AdminState='UNLOCKED'
        Method = 'DELETE'
        Target = 'core-data'
        Path = '/api/v2/event/age/604800000000000' # Remove events older than 7 days
        Interval = 'midnight'

    [SecretStore]

    Type = 'vault'

    Protocol = 'http'

    Host = 'localhost'

    Port = 8200

    Path = 'support-scheduler/'

    TokenFile = '/tmp/edgex/secrets/support-scheduler/secrets-token.json'

    RootCaCertPath = ''

    ServerName = ''
      [SecretStore.Authentication]
      AuthType = 'X-Vault-Token'

---
kind: ConfigMap
apiVersion: v1
metadata:
  name: edge-core-metadata
  namespace: edge-system
data:
  configuration.toml: >
    [Writable]

    LogLevel = 'INFO'
      [Writable.InsecureSecrets]
        [Writable.InsecureSecrets.DB]
        path = "redisdb"
          [Writable.InsecureSecrets.DB.Secrets]
          username = ""
          password = ""

    [Service]

    HealthCheckInterval = '10s'

    Host = 'edge-core-metadata'

    Port = 59881

    ServerBindAddr = '0.0.0.0' # Leave blank so default to Host value unless different value is needed.

    StartupMsg = 'This is the EdgeX Core Metadata Microservice'

    MaxResultCount = 50000

    MaxRequestSize = 0 # Not curently used. Defines the maximum size of http request body in bytes

    RequestTimeout = '5s'
    
    [Service.CORSConfiguration]
      EnableCORS = false
      CORSAllowCredentials = false
      CORSAllowedOrigin = "http://edge-core-metadata"
      CORSAllowedMethods = "GET, POST, PUT, PATCH, DELETE"
      CORSAllowedHeaders = "Authorization, Accept, Accept-Language, Content-Language, Content-Type, X-Correlation-ID"
      CORSExposeHeaders = "Cache-Control, Content-Language, Content-Length, Content-Type, Expires, Last-Modified, Pragma, X-Correlation-ID"
      CORSMaxAge = 3600


    [Registry]

    Host = 'localhost'

    Port = 8500

    Type = 'consul'


    [Clients]
      [Clients.support-notifications]
      Protocol = 'http'
      Host = 'edge-support-notifications'
      Port = 59860

      [Clients.core-data]
      Protocol = 'http'
      Host = 'edge-core-data'
      Port = 59880

    [Databases]
      [Databases.Primary]
      Host = 'edge-redis-ha-announce-0'
      Name = 'metadata'
      Password = 'password'
      Username = 'meta'
      Port = 6379
      Timeout = 5000
      Type = 'redisdb'

    [Notifications]

    PostDeviceChanges = true

    Slug = 'device-change-'

    Content = 'Device update: '

    Sender = 'core-metadata'

    Description = 'Metadata device notice'

    Label = 'metadata'


    [SecretStore]

    Type = 'vault'

    Protocol = 'http'

    Host = 'localhost'

    Port = 8200

    Path = 'core-metadata/'

    TokenFile = '/tmp/edgex/secrets/core-metadata/secrets-token.json'

    RootCaCertPath = ''

    ServerName = ''
      [SecretStore.Authentication]
      AuthType = 'X-Vault-Token'
---
kind: ConfigMap
apiVersion: v1
metadata:
  name: edge-support-notifications
  namespace: edge-system
data:
  configuration.toml: >-
    [Writable]

    LogLevel = 'INFO'

    ResendLimit = 2

    ResendInterval = '5s'
      [Writable.InsecureSecrets]
        [Writable.InsecureSecrets.DB]
        path = "redisdb"
          [Writable.InsecureSecrets.DB.Secrets]
          username = ""
          password = ""
        [Writable.InsecureSecrets.SMTP]
        path = "smtp"
          [Writable.InsecureSecrets.SMTP.Secrets]
          username = "username@mail.example.com"
          password = ""

    [Service]

    HealthCheckInterval = '10s'

    Host = 'edge-support-notifications'

    Port = 59860

    ServerBindAddr = '0.0.0.0' # Leave blank so default to Host value unless
    different value is needed.

    StartupMsg = 'This is the Support Notifications Microservice'

    MaxResultCount = 50000

    MaxRequestSize = 0 # Not curently used. Defines the maximum size of http
    request body in bytes

    RequestTimeout = '5s'
    
     [Service.CORSConfiguration]
      EnableCORS = false
      CORSAllowCredentials = false
      CORSAllowedOrigin = "http://edge-support-notifications"
      CORSAllowedMethods = "GET, POST, PUT, PATCH, DELETE"
      CORSAllowedHeaders = "Authorization, Accept, Accept-Language, Content-Language, Content-Type, X-Correlation-ID"
      CORSExposeHeaders = "Cache-Control, Content-Language, Content-Length, Content-Type, Expires, Last-Modified, Pragma, X-Correlation-ID"
      CORSMaxAge = 3600

    [Registry]

    Host = 'localhost'

    Port = 8500

    Type = 'consul'


    [Databases]
      [Databases.Primary]
      Host = 'edge-redis-ha-announce-0'
      Name = 'notifications'
      Port = 6379
      Timeout = 5000
      Type = 'redisdb'

    [Smtp]
      Host = 'smtp.gmail.com'
      Port = 587
      Sender = 'tian.jacky@gmail.com'
      EnableSelfSignedCert = false
      Subject = 'EdgeX Notification'
      # SecretPath is used to specify the secret path to store the credential(username and password) for connecting the SMTP server
      # User need to store the credential via the /secret API before sending the email notification
      SecretPath = 'smtp'
      # AuthMode is the SMTP authentication mechanism. Currently, 'usernamepassword' is the only AuthMode supported by this service, and the secret keys are 'username' and 'password'.
      AuthMode = 'usernamepassword'


    [SecretStore]

    Type = 'vault'

    Protocol = 'http'

    Host = 'localhost'

    Port = 8200

    Path = 'support-notifications/'

    TokenFile = '/tmp/edgex/secrets/support-notifications/secrets-token.json'

    RootCaCertPath = ''

    ServerName = ''
      [SecretStore.Authentication]
      AuthType = 'X-Vault-Token'
---
kind: ConfigMap
apiVersion: v1
metadata:
  name: edge-core-command
  namespace: edge-system
data:
  configuration.toml: >
    [Writable]

    LogLevel = 'DEBUG'
      [Writable.InsecureSecrets]
        [Writable.InsecureSecrets.DB]
        path = "redisdb"
          [Writable.InsecureSecrets.DB.Secrets]
          username = ""
          password = ""

    [Service]

    HealthCheckInterval = '10s'

    Host = 'edge-core-command'

    Port = 59882

    ServerBindAddr = '0.0.0.0' # Leave blank so default to Host value unless
    different value is needed.

    StartupMsg = 'This is the Core Command Microservice'

    MaxResultCount = 50000

    MaxRequestSize = 0 # Not curently used. Defines the maximum size of http
    request body in bytes

    RequestTimeout = '45s'
    
     [Service.CORSConfiguration]
      EnableCORS = false
      CORSAllowCredentials = false
      CORSAllowedOrigin = "http://edge-core-command"
      CORSAllowedMethods = "GET, POST, PUT, PATCH, DELETE"
      CORSAllowedHeaders = "Authorization, Accept, Accept-Language, Content-Language, Content-Type, X-Correlation-ID"
      CORSExposeHeaders = "Cache-Control, Content-Language, Content-Length, Content-Type, Expires, Last-Modified, Pragma, X-Correlation-ID"
      CORSMaxAge = 3600

    [Registry]

    Host = 'localhost'

    Port = 8500

    Type = 'consul'


    [Clients]
      [Clients.core-metadata]
      Protocol = 'http'
      Host = 'edge-core-metadata'
      Port = 59881

    [Databases]
      [Databases.Primary]
      Host = 'edge-redis-ha-announce-0'
      Name = 'metadata'
      Port = 6379
      Timeout = 5000
      Type = 'redisdb'

    [SecretStore]

    Type = 'vault'

    Protocol = 'http'

    Host = 'localhost'

    Port = 8200

    # Use the core-meta data secrets due to core-command using core-meta-data's
    database for persistance.

    Path = 'core-command/'

    TokenFile = '/tmp/edgex/secrets/core-command/secrets-token.json'

    RootCaCertPath = ''

    ServerName = ''
      [SecretStore.Authentication]
      AuthType = 'X-Vault-Token'
---
kind: ConfigMap
apiVersion: v1
metadata:
  name: device-dashboard
  namespace: edge-system
data:
  configuration.toml: |-
   [Service]
    Host = "device-dashboard"
    Port = 4000
    Labels = []
    OpenMsg = "edge console started"
    StaticResourcesPath = "./www"

    #Using an default memory db automatically if not configed here.
    [Database]
    Host = ""
    Name = ""
    Port = 0
    Username = "su"
    Password = "su"
      [Database.Scheme]
      User = "user"
       
      [MQTTBroker]
         Schema="tcp"
         Host="192.168.1.244"
         Port=1883
         Qos=0
         KeepAlive=3600
         IncomingTopic="DataTopic"
         ResponseTopic="ResponseTopic"

    [Clients]
        [Clients.CoreData]
        Protocol = 'http'
        Host = 'edge-core-data'
        Port = 59880
        PathPrefix = "/coredata"

        [Clients.Metadata]
        Protocol = 'http'
        Host = 'edge-core-metadata'
        Port = 59881
        PathPrefix = "/metadata"

        [Clients.CoreCommand]
        Protocol = 'http'
        Host = 'edge-core-command'
        Port = 59882
        PathPrefix = "/command"

        [Clients.Notification]
        Protocol = 'http'
        Host = 'edge-support-notifications'
        Port = 59860
        PathPrefix = "/notification"

        [Clients.Scheduler]
        Protocol = 'http'
        Host = 'edge-support-scheduler'
        Port = 59861
        PathPrefix = "/scheduler"

        [Clients.RuleEngine]
        Protocol = 'http'
        Host = 'edge-kuiper'
        Port = 59720
        PathPrefix = "/rule-engine"

        [Clients.AppService]
        Protocol = 'http'
        Host = 'edge-app-rules-engine'
        Port = 59701
        PathPrefix = "/app-service"
        
     [Deploy]
       Namespace = 'edge-system'
       Image = 'edgego/app-service-configurable:v2.2.0'
       Target ='kubernetes'
       Host = ''
       User = ''
       Password = ''
       
    [AppService]
      LogLevel = "INFO"
      StoreAndForward = true
      RetryInterval = "5s"
      MaxRetryCount = 10
      Host = 'edge-redis-ha-announce-0'
      Port=6379
      Password = ''
      Username = ''
      PublishTopic = "test"

    [Registry]
    Host = "edge-core-consul"
    Port = 8500
    Type = "consul"
    ConfigRegistryStem="edgex/appservices/"
    ServiceVersion="2.0"

    [MessageBus]
      Host = "edge-core-data"
      Port = 5563
      Type = "zero"
      
---
kind: ConfigMap
apiVersion: v1
metadata:
  name: edge-core-data
  namespace: edge-system
data:
  configuration.toml: >-
    [Writable]

    PersistData = true

    LogLevel = 'DEBUG'
       [Writable.InsecureSecrets]
          [Writable.InsecureSecrets.DB]
             path = "redisdb"
                [Writable.InsecureSecrets.DB.Secrets]
                username = ""
                password = ""

    [Service]

    HealthCheckInterval = '10s'

    Host = 'edge-core-data'

    Port = 59880

    ServerBindAddr = '0.0.0.0' # Leave blank so default to Host value unless different value is needed.

    StartupMsg = 'Core Data Microservice started'

    MaxResultCount = 50000

    MaxRequestSize = 0 # Not curently used. Defines the maximum size of http
    request body in bytes

    RequestTimeout = '5s'

    [Service.CORSConfiguration]
      EnableCORS = false
      CORSAllowCredentials = false
      CORSAllowedOrigin = "http://edge-core-data"
      CORSAllowedMethods = "GET, POST, PUT, PATCH, DELETE"
      CORSAllowedHeaders = "Authorization, Accept, Accept-Language, Content-Language, Content-Type, X-Correlation-ID"
      CORSExposeHeaders = "Cache-Control, Content-Language, Content-Length, Content-Type, Expires, Last-Modified, Pragma, X-Correlation-ID"
      CORSMaxAge = 3600


    [Registry]

    Host = 'localhost'

    Port = 8500

    Type = 'consul'
    

    [Clients]
      [Clients.core-metadata]
      Protocol = 'http'
      Host = 'edge-core-metadata'
      Port = 59881

    [Databases]
      [Databases.Primary]
      Host = 'edge-redis-ha-announce-0'
      Name = 'coredata'
      Port = 6379
      Timeout = 5000
      Type = 'redisdb'
      
    [SecretStore]

    Type = 'vault'

    Protocol = 'http'

    Host = 'localhost'

    Port = 8200

    Path = 'core-data/'

    TokenFile = '/tmp/edgex/secrets/core-data/secrets-token.json'

    RootCaCertPath = ''

    ServerName = ''
      [SecretStore.Authentication]
      AuthType = 'X-Vault-Token'
      
    [MessageQueue]
      Protocol = "tcp"
      Host = "edge-core-data"
      Port = 5563
      Type = "zero"
      AuthMode = "usernamepassword"  # required for redis messagebus (secure or insecure).
      SecretName = "redisdb"
      PublishTopicPrefix = "edgex/events/core" # /<device-profile-name>/<device-name> will be added to this Publish Topic prefix
      SubscribeEnabled = true
      SubscribeTopic = "edgex/events/device/#"  # required for subscribing to Events from MessageBus
      [MessageQueue.Optional]
          # Default MQTT Specific options that need to be here to enable evnironment variable overrides of them
          # Client Identifiers
          ClientId ="core-data"
          # Connection information
          Qos          =  "0" # Quality of Sevice values are 0 (At most once), 1 (At least once) or 2 (Exactly once)
          KeepAlive    =  "10" # Seconds (must be 2 or greater)
          Retained     = "false"
          AutoReconnect  = "true"
          ConnectTimeout = "5" # Seconds
          # TLS configuration - Only used if Cert/Key file or Cert/Key PEMblock are specified
          SkipCertVerify = "false"
---
apiVersion: v1
kind: Service
metadata:
  name: device-dashboard
  namespace: edge-system
spec:
  type: NodePort
  selector:
    app: device-dashboard
  ports:
    - port: 4000
      nodePort: 4000
      name: port-40000
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: device-dashboard
  namespace: edge-system
  labels:
    app: device-dashboard
spec:
  replicas: 1
  selector:
    matchLabels:
      app: device-dashboard
  template:
    metadata:
      name: device-dashboard
      labels:
        app: device-dashboard
    spec:
      volumes:
        - name: config
          configMap:
            name: device-dashboard
        - name: device-dashboard-volume
          persistentVolumeClaim:
            claimName: device-dashboard-claim
      containers:
        - name: device-dashboard
          image: edgego/device-dashboard:v2.2.0
          imagePullPolicy: Always
          ports:
            - containerPort: 4000
          env:
            - name: MY_POD_NAMESPACE
              valueFrom:
                fieldRef:
                  apiVersion: v1
                  fieldPath: metadata.namespace
            - name: EDGEX_SECURITY_SECRET_STORE
              value: "false"     
          volumeMounts:
            - name: device-dashboard-volume
              mountPath: /data
            - name: config
              mountPath: /res
          readinessProbe:
            tcpSocket:
              port: 4000
            initialDelaySeconds: 5
            periodSeconds: 10
          livenessProbe:
            tcpSocket:
              port: 4000
            initialDelaySeconds: 15
            periodSeconds: 20
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: edge-redis-ha
  namespace: edge-system
  labels:
    heritage: Helm
    release: edge
    chart: redis-ha-4.4.6
    app: edge-redis-ha
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: edge-redis-ha-configmap
  namespace: edge-system
  labels:
    heritage: Helm
    release: edge
    chart: redis-ha-4.4.6
    app: edge-redis-ha
data:
  redis.conf: |
    dir "/data"
    port 6379
    maxmemory 0
    maxmemory-policy volatile-lru
    min-replicas-max-lag 5
    min-replicas-to-write 1
    rdbchecksum yes
    rdbcompression yes
    repl-diskless-sync yes
    save 900 1

  sentinel.conf: |
    dir "/data"
        sentinel down-after-milliseconds mymaster 10000
        sentinel failover-timeout mymaster 180000
        maxclients 10000
        sentinel parallel-syncs mymaster 5

  init.sh: |
    HOSTNAME="$(hostname)"
    INDEX="${HOSTNAME##*-}"
    MASTER="$(redis-cli -h edge-redis-ha -p 26379 sentinel get-master-addr-by-name mymaster | grep -E '[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}')"
    MASTER_GROUP="mymaster"
    QUORUM="2"
    REDIS_CONF=/data/conf/redis.conf
    REDIS_PORT=6379
    SENTINEL_CONF=/data/conf/sentinel.conf
    SENTINEL_PORT=26379
    SERVICE=edge-redis-ha
    set -eu

    sentinel_update() {
        echo "Updating sentinel config with master $MASTER"
        eval MY_SENTINEL_ID="\${SENTINEL_ID_$INDEX}"
        sed -i "1s/^/sentinel myid $MY_SENTINEL_ID\\n/" "$SENTINEL_CONF"
        sed -i "2s/^/sentinel monitor $MASTER_GROUP $1 $REDIS_PORT $QUORUM \\n/" "$SENTINEL_CONF"
        echo "sentinel announce-ip $ANNOUNCE_IP" >> $SENTINEL_CONF
        echo "sentinel announce-port $SENTINEL_PORT" >> $SENTINEL_CONF
    }

    redis_update() {
        echo "Updating redis config"
        echo "slaveof $1 $REDIS_PORT" >> "$REDIS_CONF"
        echo "slave-announce-ip $ANNOUNCE_IP" >> $REDIS_CONF
        echo "slave-announce-port $REDIS_PORT" >> $REDIS_CONF
    }

    copy_config() {
        cp /readonly-config/redis.conf "$REDIS_CONF"
        cp /readonly-config/sentinel.conf "$SENTINEL_CONF"
    }

    setup_defaults() {
        echo "Setting up defaults"
        if [ "$INDEX" = "0" ]; then
            echo "Setting this pod as the default master"
            redis_update "$ANNOUNCE_IP"
            sentinel_update "$ANNOUNCE_IP"
            sed -i "s/^.*slaveof.*//" "$REDIS_CONF"
        else
            DEFAULT_MASTER="$(getent hosts "$SERVICE-announce-0" | awk '{ print $1 }')"
            if [ -z "$DEFAULT_MASTER" ]; then
                echo "Unable to resolve host"
                exit 1
            fi
            echo "Setting default slave config.."
            redis_update "$DEFAULT_MASTER"
            sentinel_update "$DEFAULT_MASTER"
        fi
    }

    find_master() {
        echo "Attempting to find master"
        if [ "$(redis-cli -h "$MASTER" ping)" != "PONG" ]; then
           echo "Can't ping master, attempting to force failover"
           if redis-cli -h "$SERVICE" -p "$SENTINEL_PORT" sentinel failover "$MASTER_GROUP" | grep -q 'NOGOODSLAVE' ; then
               setup_defaults
               return 0
           fi
           sleep 10
           MASTER="$(redis-cli -h $SERVICE -p $SENTINEL_PORT sentinel get-master-addr-by-name $MASTER_GROUP | grep -E '[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}')"
           if [ "$MASTER" ]; then
               sentinel_update "$MASTER"
               redis_update "$MASTER"
           else
              echo "Could not failover, exiting..."
              exit 1
           fi
        else
            echo "Found reachable master, updating config"
            sentinel_update "$MASTER"
            redis_update "$MASTER"
        fi
    }

    mkdir -p /data/conf/

    echo "Initializing config.."
    copy_config

    ANNOUNCE_IP=$(getent hosts "$SERVICE-announce-$INDEX" | awk '{ print $1 }')
    if [ -z "$ANNOUNCE_IP" ]; then
        "Could not resolve the announce ip for this pod"
        exit 1
    elif [ "$MASTER" ]; then
        find_master
    else
        setup_defaults
    fi

    if [ "${AUTH:-}" ]; then
        echo "Setting auth values"
        ESCAPED_AUTH=$(echo "$AUTH" | sed -e 's/[\/&]/\\&/g');
        sed -i "s/replace-default-auth/${ESCAPED_AUTH}/" "$REDIS_CONF" "$SENTINEL_CONF"
    fi

    echo "Ready..."

  haproxy_init.sh: |
    HAPROXY_CONF=/data/haproxy.cfg
    cp /readonly/haproxy.cfg "$HAPROXY_CONF"
    for loop in $(seq 1 10); do
      getent hosts edge-redis-ha-announce-0 && break
      echo "Waiting for service edge-redis-ha-announce-0 to be ready ($loop) ..." && sleep 1
    done
    ANNOUNCE_IP0=$(getent hosts "edge-redis-ha-announce-0" | awk '{ print $1 }')
    if [ -z "$ANNOUNCE_IP0" ]; then
      echo "Could not resolve the announce ip for edge-redis-ha-announce-0"
      exit 1
    fi
    sed -i "s/REPLACE_ANNOUNCE0/$ANNOUNCE_IP0/" "$HAPROXY_CONF"

    if [ "${AUTH:-}" ]; then
        echo "Setting auth values"
        ESCAPED_AUTH=$(echo "$AUTH" | sed -e 's/[\/&]/\\&/g');
        sed -i "s/REPLACE_AUTH_SECRET/${ESCAPED_AUTH}/" "$HAPROXY_CONF"
    fi
    for loop in $(seq 1 10); do
      getent hosts edge-redis-ha-announce-1 && break
      echo "Waiting for service edge-redis-ha-announce-1 to be ready ($loop) ..." && sleep 1
    done
    ANNOUNCE_IP1=$(getent hosts "edge-redis-ha-announce-1" | awk '{ print $1 }')
    if [ -z "$ANNOUNCE_IP1" ]; then
      echo "Could not resolve the announce ip for edge-redis-ha-announce-1"
      exit 1
    fi
    sed -i "s/REPLACE_ANNOUNCE1/$ANNOUNCE_IP1/" "$HAPROXY_CONF"

    if [ "${AUTH:-}" ]; then
        echo "Setting auth values"
        ESCAPED_AUTH=$(echo "$AUTH" | sed -e 's/[\/&]/\\&/g');
        sed -i "s/REPLACE_AUTH_SECRET/${ESCAPED_AUTH}/" "$HAPROXY_CONF"
    fi
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: edge-redis-ha
  namespace: edge-system
  labels:
    heritage: Helm
    release: edge
    chart: redis-ha-4.4.6
    app: edge-redis-ha
rules:
- apiGroups:
    - ""
  resources:
    - endpoints
  verbs:
    - get
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: edge-redis-ha
  namespace: edge-system
  labels:
    heritage: Helm
    release: edge
    chart: redis-ha-4.4.6
    app: edge-redis-ha
subjects:
- kind: ServiceAccount
  name: edge-redis-ha
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: edge-redis-ha
---
apiVersion: v1
kind: Service
metadata:
  name: edge-redis-ha-announce-0
  namespace: edge-system
  labels:
    app: redis-ha
    heritage: "Helm"
    release: "edge"
    chart: redis-ha-4.4.6
  annotations:
    service.alpha.kubernetes.io/tolerate-unready-endpoints: "true"
spec:
  publishNotReadyAddresses: true
  type: ClusterIP
  ports:
  - name: server
    port: 6379
    protocol: TCP
    targetPort: redis
  - name: sentinel
    port: 26379
    protocol: TCP
    targetPort: sentinel
  selector:
    release: edge
    app: redis-ha
    "statefulset.kubernetes.io/pod-name": edge-redis-ha-server-0
---
apiVersion: v1
kind: Service
metadata:
  name: edge-redis-ha-announce-1
  namespace: edge-system
  labels:
    app: redis-ha
    heritage: "Helm"
    release: "edge"
    chart: redis-ha-4.4.6
  annotations:
    service.alpha.kubernetes.io/tolerate-unready-endpoints: "true"
spec:
  publishNotReadyAddresses: true
  type: ClusterIP
  ports:
  - name: server
    port: 6379
    protocol: TCP
    targetPort: redis
  - name: sentinel
    port: 26379
    protocol: TCP
    targetPort: sentinel
  selector:
    release: edge
    app: redis-ha
    "statefulset.kubernetes.io/pod-name": edge-redis-ha-server-1
---

apiVersion: v1
kind: Service
metadata:
  name: edge-redis-ha
  namespace: edge-system
  labels:
    app: redis-ha
    heritage: "Helm"
    release: "edge"
    chart: redis-ha-4.4.6
  annotations:
spec:
  type: ClusterIP
  clusterIP: None
  ports:
  - name: server
    port: 6379
    protocol: TCP
    targetPort: redis
  - name: sentinel
    port: 26379
    protocol: TCP
    targetPort: sentinel
  selector:
    release: edge
    app: redis-ha
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: edge-redis-ha-server
  namespace: edge-system
  labels:
    edge-redis-ha: replica
    app: redis-ha
    heritage: "Helm"
    release: "edge"
    chart: redis-ha-4.4.6
spec:
  selector:
    matchLabels:
      release: edge
      app: redis-ha
  serviceName: edge-redis-ha
  replicas: 2
  podManagementPolicy: OrderedReady
  updateStrategy:
    type: RollingUpdate
  template:
    metadata:
      annotations:
        checksum/init-config: fabfc20ab0ee42dda389e4ff5d157400ca2329c1c6d511e99c15f6dbefaf4e93
      labels:
        release: edge
        app: redis-ha
        edge-redis-ha: replica
    spec:
    #  affinity:
    #    podAntiAffinity:
    #      requiredDuringSchedulingIgnoredDuringExecution:
    #        - labelSelector:
    #            matchLabels:
    #              app: redis-ha
    #              release: edge
    #              edge-redis-ha: replica
    #          topologyKey: kubernetes.io/hostname
    #      preferredDuringSchedulingIgnoredDuringExecution:
    #        - weight: 100
    #          podAffinityTerm:
    #            labelSelector:
    #              matchLabels:
    #                app:  redis-ha
    #                release: edge
    #                edge-redis-ha: replica
    #            topologyKey: failure-domain.beta.kubernetes.io/zone
      securityContext:
        fsGroup: 1000
        runAsNonRoot: true
        runAsUser: 1000
      serviceAccountName: edge-redis-ha
      initContainers:
      - name: config-init
        image: redis:6.2.7-alpine
        imagePullPolicy: IfNotPresent
        resources:
          {}
        command:
        - sh
        args:
        - /readonly-config/init.sh
        env:
        - name: SENTINEL_ID_0
          value: 1eb9e186ee14b693bd731e9f55c4cb33eff44f03

        - name: SENTINEL_ID_1
          value: 13fdae7286bc203d3acc5e596d8bfc63f259692c

        volumeMounts:
        - name: config
          mountPath: /readonly-config
          readOnly: true
        - name: data
          mountPath: /data
      containers:
      - name: redis
        image: redis:6.2.7-alpine
        imagePullPolicy: IfNotPresent
        command:
        - redis-server
        args:
        - /data/conf/redis.conf
        livenessProbe:
          tcpSocket:
            port: 6379
          initialDelaySeconds: 15
        resources:
          {}
        ports:
        - name: redis
          containerPort: 6379
        volumeMounts:
        - mountPath: /data
          name: data
      - name: sentinel
        image: redis:6.2.7-alpine
        imagePullPolicy: IfNotPresent
        command:
          - redis-sentinel
        args:
          - /data/conf/sentinel.conf
        livenessProbe:
          tcpSocket:
            port: 26379
          initialDelaySeconds: 15
        resources:
          {}
        ports:
          - name: sentinel
            containerPort: 26379
        volumeMounts:
        - mountPath: /data
          name: data
      volumes:
      - name: config
        configMap:
          name: edge-redis-ha-configmap
  volumeClaimTemplates:
  - metadata:
      name: data
      annotations:
    spec:
      accessModes:
        - "ReadWriteOnce"
      resources:
        requests:
          storage: "10Gi"
      storageClassName: local-path    
---
apiVersion: v1
kind: Pod
metadata:
  name: edge-redis-ha-configmap-test
  namespace: edge-system
  labels:
    app: redis-ha
    heritage: "Helm"
    release: "edge"
    chart: redis-ha-4.4.6
  annotations:
    "helm.sh/hook": test-success
spec:
  containers:
  - name: check-init
    image: koalaman/shellcheck:v0.8.0
    args:
    - --shell=sh
    - /readonly-config/init.sh
    volumeMounts:
    - name: config
      mountPath: /readonly-config
      readOnly: true
  restartPolicy: Never
  volumes:
  - name: config
    configMap:
      name: edge-redis-ha-configmap
---
apiVersion: v1
kind: Pod
metadata:
  name: edge-redis-ha-service-test
  namespace: edge-system
  labels:
    app: redis-ha
    heritage: "Helm"
    release: "edge"
    chart: redis-ha-4.4.6
  annotations:
    "helm.sh/hook": test-success
spec:
  containers:
  - name: "edge-service-test"
    image: redis:6.2.7-alpine
    command:
      - sh
      - -c
      - redis-cli -h edge-redis-ha -p 6379 info server
  restartPolicy: Never
---
apiVersion: v1
kind: Service
metadata:
  name: edge-support-notifications
  namespace: edge-system
spec:
  type: NodePort
  selector:
    app: edge-support-notifications
  ports:
    - port: 59860
      nodePort: 59860
      name: port-59860
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: edge-support-notifications
  namespace: edge-system
  labels:
    app: edge-support-notifications
spec:
  replicas: 1
  selector:
    matchLabels:
      app: edge-support-notifications
  template:
    metadata:
      name: edge-support-notifications
      labels:
        app: edge-support-notifications
    spec:
      volumes:
        - name: config
          configMap:
            name: edge-support-notifications
      containers:
        - name: edge-support-notifications
          image: edgego/support-notifications:v2.2.0
          command:
            - /support-notifications
            - '--confdir'
            - /res
          imagePullPolicy: Always
          ports:
            - containerPort: 59860
          env:
             - name: EDGEX_SECURITY_SECRET_STORE
               value: "false"
          volumeMounts:
            - name: config
              mountPath: /res
          readinessProbe:
            httpGet:
              path: /api/v2/ping
              port: 59860
            initialDelaySeconds: 5
            periodSeconds: 10
          livenessProbe:
            httpGet:
              path: /api/v2/ping
              port: 59860
            initialDelaySeconds: 15
            periodSeconds: 20
---
apiVersion: v1
kind: Service
metadata:
  name: edge-core-data
  namespace: edge-system
spec:
  type: NodePort
  selector:
    app: edge-core-data
  ports:
    - port: 59880
      nodePort: 59880
      name: port-59880
    - port: 5563
      name: port-5563
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: edge-core-data
  namespace: edge-system
  labels:
    app: edge-core-data
spec:
  replicas: 1
  selector:
    matchLabels:
      app: edge-core-data
  template:
    metadata:
      name: edge-core-data
      labels:
        app: edge-core-data
    spec:
      volumes:
        - name: config
          configMap:
            name: edge-core-data
      containers:
        - name: edge-core-data
          image: edgego/core-data:v2.2.0
          command:
            - /core-data
            - '--confdir'
            - /res
          imagePullPolicy: Always
          env:
            - name: EDGEX_SECURITY_SECRET_STORE
              value: "false"
          volumeMounts:
           - name: config
             mountPath: /res
          ports:
            - containerPort: 59880
            - containerPort: 5563
          readinessProbe:
            httpGet:
              path: /api/v2/ping
              port: 59880
            initialDelaySeconds: 5
            periodSeconds: 10
          livenessProbe:
            httpGet:
              path: /api/v2/ping
              port: 59880
            initialDelaySeconds: 15
            periodSeconds: 20
---
apiVersion: v1
kind: Service
metadata:
  name: edge-core-metadata
  namespace: edge-system
spec:
  type: NodePort
  selector:
    app: edge-core-metadata
  ports:
    - port: 59881
      nodePort: 59881
      name: port-59881
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: edge-core-metadata
  namespace: edge-system
  labels:
    app: edge-core-metadata
spec:
  replicas: 1
  selector:
    matchLabels:
      app: edge-core-metadata
  template:
    metadata:
      name: edge-core-metadata
      labels:
        app: edge-core-metadata
    spec:
      volumes:
        - name: config
          configMap:
            name: edge-core-metadata
      containers:
        - name: edge-core-metadata
          image: edgego/core-metadata:v2.2.0
          command:
            - /core-metadata
            - '--confdir'
            - /res
          env:
            - name: EDGEX_SECURITY_SECRET_STORE
              value: "false"
          volumeMounts:
             - name: config
               mountPath: /res
          ports:
            - containerPort: 59881
          readinessProbe:
            httpGet:
              path: /api/v2/ping
              port: 59881
            initialDelaySeconds: 5
            periodSeconds: 10
          livenessProbe:
            httpGet:
              path: /api/v2/ping
              port: 59881
            initialDelaySeconds: 15
            periodSeconds: 20
---
apiVersion: v1
kind: Service
metadata:
  name: edge-core-command
  namespace: edge-system
spec:
  type: NodePort
  selector:
    app: edge-core-command
  ports:
    - port: 59882
      nodePort: 59882
      name: port-59882
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: edge-core-command
  namespace: edge-system
  labels:
    app: edge-core-command
spec:
  replicas: 1
  selector:
    matchLabels:
      app: edge-core-command
  template:
    metadata:
      name: edge-core-command
      labels:
        app: edge-core-command
    spec:
      volumes:
        - name: config
          configMap:
            name: edge-core-command
      containers:
        - name: edge-core-command
          image: edgego/core-command:v2.2.0
          command:
            - /core-command
            - '--confdir'
            - /res
          imagePullPolicy: Always
          env:
            - name: EDGEX_SECURITY_SECRET_STORE
              value: "false"
          volumeMounts:
            - name: config
              mountPath: /res
          ports:
            - containerPort: 59882
          readinessProbe:
            httpGet:
              path: /api/v2/ping
              port: 59882
            initialDelaySeconds: 5
            periodSeconds: 10
          livenessProbe:
            httpGet:
              path: /api/v2/ping
              port: 59882
            initialDelaySeconds: 15
            periodSeconds: 20
---
apiVersion: v1
kind: Service
metadata:
  name: edge-support-scheduler
  namespace: edge-system
spec:
  type: NodePort
  selector:
    app: edge-support-scheduler
  ports:
    - port: 59861
      nodePort: 59861
      name: port-59861
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: edge-support-scheduler
  namespace: edge-system
  labels:
    app: edge-support-scheduler
spec:
  replicas: 1
  selector:
    matchLabels:
      app: edge-support-scheduler
  template:
    metadata:
      name: edge-support-scheduler
      labels:
        app: edge-support-scheduler
    spec:
      volumes:
        - name: config
          configMap:
            name: edge-support-scheduler
      containers:
        - name: edge-support-scheduler
          image: edgego/support-scheduler:v2.2.0
          command:
            - /support-scheduler
            - '--confdir'
            - /res
          imagePullPolicy: Always
          env:
            - name: EDGEX_SECURITY_SECRET_STORE
              value: "false"
          volumeMounts:
            - name: config
              mountPath: /res
          ports:
            - containerPort: 59861
          readinessProbe:
            httpGet:
              path: /api/v2/ping
              port: 59861
            initialDelaySeconds: 5
            periodSeconds: 10
          livenessProbe:
            httpGet:
              path: /api/v2/ping
              port: 59861
            initialDelaySeconds: 15
            periodSeconds: 20
---
kind: ConfigMap
apiVersion: v1
metadata:
  name: kuiper-conn
  namespace: edge-system
data:
  connection.yaml: |-
    edgex:
       redisMsgBus: #connection key
          protocol: redis
          server: edge-redis-ha-announce-0
          port: 6379
          type: redis

    mqtt:
      localConnection: #connection key
        servers: [tcp://127.0.0.1:1883]
        username: ekuiper
        password: password
        #certificationPath: /var/kuiper/xyz-certificate.pem
        #privateKeyPath: /var/kuiper/xyz-private.pem.key
        #rootCaPath: /var/kuiper/xyz-rootca.pem
        #insecureSkipVerify: false
        #protocolVersion: 3
        clientid: ekuiper
      cloudConnection: #connection key
        servers: ["tcp://broker.emqx.io:1883"]
        username: user1
        password: password
        #certificationPath: /var/kuiper/xyz-certificate.pem
        #privateKeyPath: /var/kuiper/xyz-private.pem.ke
        #rootCaPath: /var/kuiper/xyz-rootca.pem
        #insecureSkipVerify: false
        #protocolVersion: 3


    edgex:
      #redisMsgBus: #redis connection key
      #  protocol: redis
      #  server: edge-redis-ha-announce-0
      #  port: 6379
      #  type: redis
        #  Below is optional configurations settings for mqtt
        #  type: mqtt
        #  optional:
        #    ClientId: client1
        #    Username: user1
        #    Password: password
        #    Qos: 1
        #    KeepAlive: 5000
        #    Retained: true/false
        #    ConnectionPayload:
        #    CertFile:
        #    KeyFile:
        #    CertPEMBlock:
        #    KeyPEMBlock:
        #    SkipCertVerify: true/false
      #mqttMsgBus: #connection key
      #  protocol: tcp
      #  server: 127.0.0.1
      #  port: 1883
      #  type: mqtt
      #  optional:
      #    ClientId: "client1"

      zeroMsgBus: #connection key
        protocol: tcp
        server: localhost
        port: 5571
        type: zero
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: edge-kuiper-claim
  namespace: edge-system
spec:
  storageClassName: local-path
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 100M
---
kind: ConfigMap
apiVersion: v1
metadata:
  name: kuiper-source
  namespace: edge-system
data:
  edgex.yaml: |
    app1:
        messageType: event
        port: 6379
        protocol: redis
        server: edge-redis-ha-announce-0
        topic: app1-rules-event
        type: redis
    app2:
        messageType: event
        port: 6379
        protocol: redis
        server: edge-redis-ha-announce-0
        topic: app2-rules-event
        type: redis
    app3:
        messageType: event
        port: 6379
        protocol: redis
        server: edge-redis-ha-announce-0
        topic: app3-rule-events
        type: redis
    default:
        messageType: request
        port: 6379
        protocol: redis
        server: edge-redis-ha-announce-0
        topic: edgex/events/device#
        type: redis
    demo1:
        messageType: request
        port: 6379
        protocol: redis
        server: edge-redis-ha-announce-0
        topic: edgex/events/device/Test-Device-Modbus-Profile/modbus-tcp#
        type: redis
    demo2:
        messageType: event
        port: 6379
        protocol: redis
        server: edge-redis-ha-announce-0
        topic: rules-events
        type: redis
    demo3:
        messageType: request
        port: 6379
        protocol: redis
        server: edge-redis-ha-announce-0
        topic: edgex/events/device/Test-Device-Modbus-Profile/modbus-tcp#
        type: redis
    share_conf:
        connectionSelector: edgex.mqttMsgBus
        messageType: event
        port: 1883
        protocol: tcp
        server: edge-app-rules-engine
        topic: rules-events
        type: zero
---
kind: ConfigMap
apiVersion: v1
metadata:
  name: kuiper-conf
  namespace: edge-system
data:
  kuiper.yaml: >-
    basic:
      # true|false, with debug level, it prints more debug info
      debug: true
      # true|false, if it's set to true, then the log will be print to console
      consoleLog: true
      # true|false, if it's set to true, then the log will be print to log file
      fileLog: true
      # How many hours to split the file
      rotateTime: 24
      # Maximum file storage hours
      maxAge: 72
      # CLI ip
      ip: 0.0.0.0
      # CLI port
      port: 20498
      # REST service ip
      restIp: 0.0.0.0
      # REST service port
      restPort: 59720
      # true|false, when true, will check the RSA jwt token for rest api
      authentication: false
      #  restTls:
      #    certfile: /var/https-server.crt
      #    keyfile: /var/https-server.key
      # Prometheus settings
      prometheus: false
      prometheusPort: 20499
      # The URL where hosts all of pre-build plugins. By default it's at packages.emqx.net
      pluginHosts: https://packages.emqx.net
      # Whether to ignore case in SQL processing. Note that, the name of customized function by plugins are case-sensitive.
      ignoreCase: true

    # The default options for all rules. Each rule can override this setting by
    defining its own option

    rule:
      # The qos of the rule. The values can be 0: At most once; 1: At least once; 2: Exactly once
      # If qos is bigger than 0, the checkpoint mechanism will launch to save states so that they can be
      # restored for unintended interrupt or planned restart of the rule. The performance may be affected
      # to enable the checkpoint mechanism
      qos: 0
      # The interval in millisecond to run the checkpoint mechanism.
      checkpointInterval: 300000
      # Whether to send errors to sinks
      sendError: true

    sink:
      # The cache persistence threshold size. If the message in sink cache is larger than 10, then it triggers persistence. If you find
      # the remote system is slow to response, or sink throughput is small, then it's recommend to increase below 2 configurations.
      # More memory is required with the increase of below 2 configurations.
      # If the message count reaches below value, then it triggers persistence.
      cacheThreshold: 10
      # The message persistence is triggered by a ticker, and cacheTriggerCount is for using configure the count to trigger the persistence procedure
      # regardless if the message number reaches cacheThreshold or not. This is to prevent the data won't be saved as the cache never pass the threshold.
      cacheTriggerCount: 15

      # Control to disable cache or not. If it's set to true, then the cache will be disabled, otherwise, it will be enabled.
      disableCache: true

    store:
      #Type of store that will be used for keeping state of the application
      type: sqlite
      redis:
        host: edge-redis-ha-announce-0
        port: 6379
        password: kuiper
        #Timeout in ms
        timeout: 1000
      sqlite:
        #Sqlite file name, if left empty name of db will be sqliteKV.db
        name:

    # The settings for portable plugin

    portable:
      # The executable of python. Specify this if you have multiple python instances in your system
      # or other circumstance where the python executable cannot be successfully invoked through the default command.
      pythonBin: python

---
kind: Service
apiVersion: v1
metadata:
  name: edge-kuiper
  namespace: edge-system
spec:
  ports:
    - name: port-59720
      protocol: TCP
      port: 59720
      targetPort: 59720
      nodePort: 59720
    - name: port-20498
      protocol: TCP
      port: 20498
      targetPort: 20498
      nodePort: 24048
  selector:
    app: edge-kuiper
  clusterIP: 10.43.209.22
  clusterIPs:
    - 10.43.209.22
  type: NodePort
  sessionAffinity: None
  externalTrafficPolicy: Cluster
  ipFamilies:
    - IPv4
  ipFamilyPolicy: SingleStack
  internalTrafficPolicy: Cluster
---
kind: Deployment
apiVersion: apps/v1
metadata:
  name: edge-kuiper
  namespace: edge-system
spec:
  replicas: 1
  selector:
    matchLabels:
      app: edge-kuiper
  template:
    metadata:
      name: edge-kuiper
      creationTimestamp: null
      labels:
        app: edge-kuiper
    spec:
      volumes:
        - name: kuiper-source
          configMap:
            name: kuiper-source
            defaultMode: 420
        - name: kuiper-conf
          configMap:
            name: kuiper-conf
            defaultMode: 420
        - name: kuiper-conn
          configMap:
            name: kuiper-conn
            defaultMode: 420
        - name: kuiper-data
          persistentVolumeClaim:
            claimName: edge-kuiper-claim
      containers:
        - name: edge-kuiper
          image: lfedge/ekuiper:1.6.0-alpine
          ports:
            - containerPort: 59720
              protocol: TCP
            - containerPort: 20498
              protocol: TCP
          env:
            - name: EDGEX_SECURITY_SECRET_STORE
              value: 'false'
          resources: {}
          volumeMounts:
            - name: kuiper-data
              mountPath: /kuiper/data
            - name: kuiper-source
              mountPath: /kuiper/etc/sources/edgex.yaml
              subPath: edgex.yaml
            - name: kuiper-conf
              mountPath: /kuiper/etc/kuiper.yaml
              subPath: kuiper.yaml
            - name: kuiper-conn
              mountPath: /kuiper/etc/connections/connection.yaml
              subPath: connection.yaml
          livenessProbe:
            httpGet:
              path: /
              port: 59720
              scheme: HTTP
            initialDelaySeconds: 15
            timeoutSeconds: 1
            periodSeconds: 20
            successThreshold: 1
            failureThreshold: 3
          readinessProbe:
            httpGet:
              path: /
              port: 59720
              scheme: HTTP
            initialDelaySeconds: 5
            timeoutSeconds: 1
            periodSeconds: 10
            successThreshold: 1
            failureThreshold: 3
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          imagePullPolicy: IfNotPresent
      restartPolicy: Always
      terminationGracePeriodSeconds: 30
      dnsPolicy: ClusterFirst
      securityContext:
        runAsUser: 0
      schedulerName: default-scheduler
---
kind: ConfigMap
apiVersion: v1
metadata:
  name: app-service-configurable
  namespace: edge-system
data:
  configuration.toml: >
    [Writable]

    LogLevel = "DEBUG"


    [Writable.StoreAndForward]
      Enabled = false
      RetryInterval = "5m"
      MaxRetryCount = 10

      [Writable.Pipeline]
      ExecutionOrder = "SetResponseData"
        [Writable.Pipeline.Functions.SetResponseData]
          [Writable.Pipeline.Functions.SetResponseData.Parameters]
          ResponseContentType = ""
        [Writable.Pipeline.Functions.FilterByProfileName]
          [Writable.Pipeline.Functions.FilterByProfileName.Parameters]
          ProfileNames = ""
          FilterOut = "false"
        [Writable.Pipeline.Functions.FilterByDeviceName]
          [Writable.Pipeline.Functions.FilterByDeviceName.Parameters]
          DeviceNames = ""
          FilterOut = "false"
        [Writable.Pipeline.Functions.FilterBySourceName]
          [Writable.Pipeline.Functions.FilterBySourceName.Parameters]
          SourceNames = ""
          FilterOut = "false"
        [Writable.Pipeline.Functions.FilterByResourceName]
          [Writable.Pipeline.Functions.FilterByResourceName.Parameters]
          ResourceNames = ""
          FilterOut = "false"
      # InsecureSecrets are required for Redis is used for message bus
      [Writable.InsecureSecrets]
        [Writable.InsecureSecrets.DB]
        path = "redisdb"
          [Writable.InsecureSecrets.DB.Secrets]
          username = ""
          password = ""

    [Service]

    HealthCheckInterval = "10s"

    Host = "edge-app-rules-engine"

    Port = 59701

    ServerBindAddr = "0.0.0.0" # if blank, uses default Go behavior https://golang.org/pkg/net/#Listen

    StartupMsg = "app-rules-engine has Started"

    MaxResultCount = 0 # Not curently used by App Services.

    MaxRequestSize = 0 # Not curently used by App Services.

    RequestTimeout = "5s"


    [Registry]

    Host = "localhost"

    Port = 8500

    Type = "consul"


    # Database is require when Redis is used for message bus

    # Type is used as the secret name when getting credentials from the Secret
    Store

    [Database]

    Type = "redisdb"

    Host = "edge-redis-ha-announce-0"

    Port = 6379

    Timeout = "30s"


    # SecretStore is required when Store and Forward is enabled and running with
    security

    # so Database credentials can be pulled from Vault. Also now require when
    running with secure Consul

    # Note when running in docker from compose file set the following environment variables:

    #   - SecretStore_Host: edgex-vault

    [SecretStore]

    Type = 'vault'

    Host = 'localhost'

    Port = 8200

    Path = 'app-rules-engine/'

    Protocol = 'http'

    RootCaCertPath = ''

    ServerName = ''

    TokenFile = '/tmp/edgex/secrets/app-rules-engine/secrets-token.json'
      [SecretStore.Authentication]
      AuthType = 'X-Vault-Token'

    [Clients]
      # Used for version check on start-up
      [Clients.core-metadata]
      Protocol = 'http'
      Host = 'edge-core-metadata'
      Port = 59881


    [Trigger]

    Type="edgex-messagebus"
      [Trigger.EdgexMessageBus]
      Type = "zero"
        [Trigger.EdgexMessageBus.SubscribeHost]
          Host = "edge-core-data"
          Port = 5563
          Protocol = "tcp"
          SubscribeTopics="edgex/events/#"
        [Trigger.EdgexMessageBus.PublishHost]
          Host = "*"
          Port = 5566
          Protocol = "tcp"
          PublishTopic="rules-events"
        [Trigger.EdgexMessageBus.Optional]
        authmode = 'usernamepassword'  # required for redis messagebus (secure or insecure).
        secretname = 'redisdb'
        # Default MQTT Specific options that need to be here to enable evnironment variable overrides of them
        # Client Identifiers
        ClientId ="core-data"
        # Connection information
        Qos          =  "0" # Quality of Sevice values are 0 (At most once), 1 (At least once) or 2 (Exactly once)
        KeepAlive    =  "10" # Seconds (must be 2 or greater)
        Retained     = "false"
        AutoReconnect  = "true"
        ConnectTimeout = "5" # Seconds
        # TLS configuration - Only used if Cert/Key file or Cert/Key PEMblock are specified
        SkipCertVerify = "false"
---
apiVersion: v1
kind: Service
metadata:
  name: edge-app-rules-engine
  namespace: edge-system
spec:
  type: NodePort
  selector:
    app: edge-app-rules-engine
  ports:
    - port: 59701
      nodePort: 59701
      name: port-59701
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: edge-app-rules-engine
  namespace: edge-system
  labels:
    app: edge-app-rules-engine
spec:
  replicas: 1
  selector:
    matchLabels:
      app: edge-app-rules-engine
  template:
    metadata:
      name: edge-app-rules-engine
      labels:
        app: edge-app-rules-engine
    spec:
      volumes:
        - name: config
          configMap:
            name: app-service-configurable
      containers:
        - name: app-service-configurable
          image: edgego/app-service-configurable:v2.2.0
          command:
            - /app-service-configurable
            - '--confdir'
            - /res
            - '--skipVersionCheck'
          env:
            - name: EDGEX_SECURITY_SECRET_STORE
              value: "false"
          volumeMounts:
            - name: config
              mountPath: /res
          ports:
            - containerPort: 59701
          readinessProbe:
            httpGet:
              path: /api/v2/ping
              port: 59701
            initialDelaySeconds: 5
            periodSeconds: 10
          livenessProbe:
            httpGet:
              path: /api/v2/ping
              port: 59701
            initialDelaySeconds: 15
            periodSeconds: 20
---
`
