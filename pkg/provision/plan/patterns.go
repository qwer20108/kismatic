package plan

type Plan struct {
	Etcd                []Node
	Master              []Node
	Worker              []Node
	Ingress             []Node
	Storage             []Node
	MasterNodeFQDN      string
	MasterNodeShortName string
	SSHUser             string
	SSHKeyFile          string
	AdminPassword       string
}

const OverlayNetworkPlan = `cluster:
  name: kubernetes

  # This password is used to login to the Kubernetes Dashboard and can also be
  # used for administration without a security certificate.
  admin_password: {{.AdminPassword}}

  # When true, installation will not occur if any node is missing the correct
  # deb/rpm packages.
  disable_package_installation: false

  # Set to true if you are performing a disconnected installation.
  disconnected_installation: false

  # Networking configuration of your cluster.
  networking:

    # Kubernetes will assign pods IPs in this range. Do not use a range that is
    # already in use on your local network!
    pod_cidr_block: 172.16.0.0/16

    # Kubernetes will assign services IPs in this range. Do not use a range
    # that is already in use by your local network or pod network!
    service_cidr_block: 172.20.0.0/16

    # When true, the installer will add entries for all nodes to other nodes'
    # hosts files. Use when you don't have access to DNS.
    update_hosts_files: true

    # Set the proxy server to use for HTTP connections.
    http_proxy: ""

    # Set the proxy server to use for HTTPs connections.
    https_proxy: ""

    # List of host names and/or IPs that shouldn't go through any proxy.
    # All nodes' 'host' and 'IPs' are always set.
    no_proxy: ""

  # Generated certs configuration.
  certificates:

    # Self-signed certificate expiration period in hours; default is 2 years.
    expiry: 17520h

    # CA certificate expiration period in hours; default is 2 years.
    ca_expiry: 17520h

  # SSH configuration for cluster nodes.
  ssh:

    # This user must be able to sudo without password.
    user: {{.SSHUser}}

    # Absolute path to the ssh private key we should use to manage nodes.
    ssh_key: {{.SSHKeyFile}}
    ssh_port: 22

  # Override configuration of Kubernetes components.
  kube_apiserver:
    option_overrides: {}

  kube_controller_manager:
    option_overrides: {}

  kube_scheduler:
    option_overrides: {}

  kube_proxy:
    option_overrides: {}

  kubelet:
    option_overrides: {}

  # Kubernetes cloud provider integration
  cloud_provider:

    # Options: 'aws','azure','cloudstack','fake','gce','mesos','openstack',
    # 'ovirt','photon','rackspace','vsphere'.
    # Leave empty for bare metal setups or other unsupported providers.
    provider: ""

    # Path to the config file, leave empty if provider does not require it.
    config: ""

# Docker daemon configuration of all cluster nodes
docker:
  storage:

    # Configure devicemapper in direct-lvm mode (RHEL/CentOS only).
    direct_lvm:
      enabled: false

      # Path to the block device that will be used for direct-lvm mode. This
      # device will be wiped and used exclusively by docker.
      block_device: ""

      # Set to true if you want to enable deferred deletion when using
      # direct-lvm mode.
      enable_deferred_deletion: false

# If you want to use an internal registry for the installation or upgrade, you
# must provide its information here. You must seed this registry before the
# installation or upgrade of your cluster. This registry must be accessible from
# all nodes on the cluster.
docker_registry:

  # IP or hostname for your registry.
  address: ""

  # Port for your registry.
  port: 8443

  # Absolute path to the certificate authority that should be trusted when
  # connecting to your registry.
  CA: ""

  # Leave blank for unauthenticated access.
  username: ""

  # Leave blank for unauthenticated access.
  password: ""

# Add-ons are additional components that KET installs on the cluster.
add_ons:
  cni:
    disable: false

    # Selecting 'custom' will result in a CNI ready cluster, however it is up to
    # you to configure a plugin after the install.
    # Options: 'calico','weave','contiv','custom'.
    provider: calico
    options:
      calico:

        # Options: 'overlay','routed'.
        mode: overlay

        # Options: 'warning','info','debug'.
        log_level: info

  dns:
    disable: false

  heapster:
    disable: false
    options:
      heapster:
        replicas: 2

        # Specify kubernetes ServiceType. Defaults to 'ClusterIP'.
        # Options: 'ClusterIP','NodePort','LoadBalancer','ExternalName'.
        service_type: ClusterIP

        # Specify the sink to store heapster data. Defaults to an influxdb pod
        # running on the cluster.
        sink: influxdb:http://heapster-influxdb.kube-system.svc:8086

      influxdb:

        # Provide the name of the persistent volume claim that you will create
        # after installation. If not specified, the data will be stored in
        # ephemeral storage.
        pvc_name: ""

  dashboard:
    disable: false
    
  package_manager:
    disable: false

    # Options: 'helm'
    provider: helm

# Etcd nodes are the ones that run the etcd distributed key-value database.
etcd:
  expected_count: {{len .Etcd}}

  # Provide the hostname and IP of each node. If the node has an IP for internal
  # traffic, provide it in the internalip field. Otherwise, that field can be
  # left blank.
  nodes:{{range .Etcd}}
  - host: {{.Host}}
    ip: {{.PublicIPv4}}
    internalip: {{.PrivateIPv4}}{{end}}

# Master nodes are the ones that run the Kubernetes control plane components.
master:
  expected_count: {{len .Master}}

  # If you have set up load balancing for master nodes, enter the FQDN name here.
  # Otherwise, use the IP address of a single master node.
  load_balanced_fqdn: {{.MasterNodeFQDN}}

  # If you have set up load balancing for master nodes, enter the short name here.
  # Otherwise, use the IP address of a single master node.
  load_balanced_short_name: {{.MasterNodeShortName}}  
  nodes:{{range .Master}}
  - host: {{.Host}}
    ip: {{.PublicIPv4}}
    internalip: {{.PrivateIPv4}}
    labels: {}{{end}}

# Worker nodes are the ones that will run your workloads on the cluster.
worker:
  expected_count: {{len .Worker}}
  nodes:{{range .Worker}}
  - host: {{.Host}}
    ip: {{.PublicIPv4}}
    internalip: {{.PrivateIPv4}}
    labels: {}{{end}}

# Ingress nodes will run the ingress controllers.
ingress:
  expected_count: {{len .Ingress}}
  nodes:{{range .Ingress}}
  - host: {{.Host}}
    ip: {{.PublicIPv4}}
    internalip: {{.PrivateIPv4}}
    labels: {}{{end}}

# Storage nodes will be used to create a distributed storage cluster that can
# be consumed by your workloads.
storage:
  expected_count: {{len .Storage}}
  nodes:{{range .Storage}}
  - host: {{.Host}}
    ip: {{.PublicIPv4}}
    internalip: {{.PrivateIPv4}}
    labels: {}{{end}}
`
