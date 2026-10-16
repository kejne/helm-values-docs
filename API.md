# Helm values API

## Contents

- [appDefaults](#appdefaults)
  - [containerPort](#appdefaults_containerport)
  - [probes](#appdefaults_probes)
    - [port](#appdefaults_probes_port)
  - [resources](#appdefaults_resources)
- [apps](#apps)
  - [app](#apps_app)
    - [autoscaling](#apps_app_autoscaling)
      - [maxReplicas](#apps_app_autoscaling_maxreplicas)
      - [minReplicas](#apps_app_autoscaling_minreplicas)
      - [targetCpuPercent](#apps_app_autoscaling_targetcpupercent)
    - [containerPort](#apps_app_containerport)
    - [env](#apps_app_env)
    - [extraContainers](#apps_app_extracontainers)
    - [extraEgress](#apps_app_extraegress)
    - [extraIngress](#apps_app_extraingress)
    - [httpGateway](#apps_app_httpgateway)
      - [gated](#apps_app_httpgateway_gated)
      - [private](#apps_app_httpgateway_private)
      - [public](#apps_app_httpgateway_public)
    - [nats](#apps_app_nats)
      - [user](#apps_app_nats_user)
        - [permissions](#apps_app_nats_user_permissions)
          - [publish](#apps_app_nats_user_permissions_publish)
            - [allow](#apps_app_nats_user_permissions_publish_allow)
          - [responses](#apps_app_nats_user_permissions_responses)
            - [enabled](#apps_app_nats_user_permissions_responses_enabled)
            - [max](#apps_app_nats_user_permissions_responses_max)
            - [ttlSeconds](#apps_app_nats_user_permissions_responses_ttlseconds)
          - [subscribe](#apps_app_nats_user_permissions_subscribe)
            - [allow](#apps_app_nats_user_permissions_subscribe_allow)
    - [postgres](#apps_app_postgres)
      - [bindings](#apps_app_postgres_bindings)
    - [probes](#apps_app_probes)
      - [liveness](#apps_app_probes_liveness)
      - [livenessProbe](#apps_app_probes_livenessprobe)
      - [port](#apps_app_probes_port)
      - [readiness](#apps_app_probes_readiness)
      - [readinessProbe](#apps_app_probes_readinessprobe)
      - [startup](#apps_app_probes_startup)
    - [profile](#apps_app_profile)
    - [replicaCount](#apps_app_replicacount)
    - [resources](#apps_app_resources)
    - [secretEnv](#apps_app_secretenv)
      - [secretEnv-object](#apps_app_secretenv_secretenv_object)
        - [secretKey](#apps_app_secretenv_secretenv_object_secretkey)
        - [secretName](#apps_app_secretenv_secretenv_object_secretname)
    - [strategy](#apps_app_strategy)
    - [tag](#apps_app_tag)
    - [volumeMounts](#apps_app_volumemounts)
      - [configMap](#apps_app_volumemounts_configmap)
        - [configMap-object](#apps_app_volumemounts_configmap_configmap_object)
          - [optional](#apps_app_volumemounts_configmap_configmap_object_optional)
          - [path](#apps_app_volumemounts_configmap_configmap_object_path)
      - [ephemeralDirs](#apps_app_volumemounts_ephemeraldirs)
        - [ephemeralDir](#apps_app_volumemounts_ephemeraldirs_ephemeraldir)
          - [path](#apps_app_volumemounts_ephemeraldirs_ephemeraldir_path)
      - [secrets](#apps_app_volumemounts_secrets)
        - [secret](#apps_app_volumemounts_secrets_secret)
          - [optional](#apps_app_volumemounts_secrets_secret_optional)
          - [path](#apps_app_volumemounts_secrets_secret_path)
- [drive-system-lib](#drive_system_lib)
- [fullnameOverride](#fullnameoverride)
- [global](#global)
  - [domain](#global_domain)
  - [image](#global_image)
    - [registry](#global_image_registry)
  - [stageName](#global_stagename)
  - [systemName](#global_systemname)
  - [verticalId](#global_verticalid)
- [otel](#otel)
  - [enabled](#otel_enabled)
- [secrets](#secrets)

| **Schema** | **Details** |
| --- | --- |
| Type | `object` |
| Additional properties | Not allowed |


<a id="appdefaults"></a>
## `appDefaults`

| **Property** | **Details** |
| --- | --- |
| Type | `object` |
| Required | No |
| Default | See nested values |
| Additional properties | Not allowed |

**Description:** Default container port, probe port, and resource requests/limits inherited by apps unless overridden per app.

| Property | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| [containerPort](#appdefaults_containerport) | `integer` | No | `3000` | Default port exposed by each app container and Service. |

<a id="appdefaults_containerport"></a>
### `containerPort`

| **Property** | **Details** |
| --- | --- |
| Type | `integer` |
| Required | No |
| Default | `3000` |

**Description:** Default port exposed by each app container and Service.
| [probes](#appdefaults_probes) | `object` | No | See nested values | Default probe port for app workloads. |

<a id="appdefaults_probes"></a>
### `probes`

| **Property** | **Details** |
| --- | --- |
| Type | `object` |
| Required | No |
| Default | See nested values |
| Additional properties | Not allowed |

**Description:** Default probe port for app workloads.

| Property | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| [port](#appdefaults_probes_port) | `integer` | No | `3000` | Default port used by shorthand HTTP readiness and liveness probes. |

<a id="appdefaults_probes_port"></a>
#### `port`

| **Property** | **Details** |
| --- | --- |
| Type | `integer` |
| Required | No |
| Default | `3000` |

**Description:** Default port used by shorthand HTTP readiness and liveness probes.
| [resources](#appdefaults_resources) | `object` | No | `{"limits":{"cpu":"200m","memory":"128Mi"},"requests":{"cpu":"100m","memory":"64Mi"}}` | Default Kubernetes container resource requests and limits inherited by apps without an app-level override. |

<a id="appdefaults_resources"></a>
### `resources`

| **Property** | **Details** |
| --- | --- |
| Type | `object` |
| Required | No |
| Default | `{"limits":{"cpu":"200m","memory":"128Mi"},"requests":{"cpu":"100m","memory":"64Mi"}}` |
| Additional properties | Any type allowed |

**Description:** Default Kubernetes container resource requests and limits inherited by apps without an app-level override.

<a id="apps"></a>
## `apps`

| **Property** | **Details** |
| --- | --- |
| Type | `object` |
| Required | No |
| Default | `{}` |
| Additional properties | Each additional property must conform to the schema (`object`) |

**Description:** Application workloads keyed by app name. Each key becomes the workload and service name.

| Property | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| [app](#apps_app) | `object` | No | See nested values | — |

<a id="apps_app"></a>
### `app`

| **Property** | **Details** |
| --- | --- |
| Type | `object` |
| Required | No |
| Default | See nested values |
| Additional properties | Not allowed |

| Property | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| [autoscaling](#apps_app_autoscaling) | `object` | No | See nested values | Configure a HorizontalPodAutoscaler instead of a fixed replica count. |

<a id="apps_app_autoscaling"></a>
#### `autoscaling`

| **Property** | **Details** |
| --- | --- |
| Type | `object` |
| Required | No |
| Default | See nested values |
| Additional properties | Not allowed |

**Description:** Configure a HorizontalPodAutoscaler instead of a fixed replica count.

| Property | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| [maxReplicas](#apps_app_autoscaling_maxreplicas) | `integer` | No | — | Maximum number of app replicas maintained by the autoscaler. |

<a id="apps_app_autoscaling_maxreplicas"></a>
##### `maxReplicas`

| **Property** | **Details** |
| --- | --- |
| Type | `integer` |
| Required | Yes |
| Default | — |
| Constraints | minimum=1 |

**Description:** Maximum number of app replicas maintained by the autoscaler.
| [minReplicas](#apps_app_autoscaling_minreplicas) | `integer` | No | — | Minimum number of app replicas maintained by the autoscaler. |

<a id="apps_app_autoscaling_minreplicas"></a>
##### `minReplicas`

| **Property** | **Details** |
| --- | --- |
| Type | `integer` |
| Required | Yes |
| Default | — |
| Constraints | minimum=0 |

**Description:** Minimum number of app replicas maintained by the autoscaler.
| [targetCpuPercent](#apps_app_autoscaling_targetcpupercent) | `integer` | No | — | Target average CPU utilization percentage used for autoscaling. |

<a id="apps_app_autoscaling_targetcpupercent"></a>
##### `targetCpuPercent`

| **Property** | **Details** |
| --- | --- |
| Type | `integer` |
| Required | No |
| Default | — |
| Constraints | minimum=1 |

**Description:** Target average CPU utilization percentage used for autoscaling.
| [containerPort](#apps_app_containerport) | `integer` | No | — | Port exposed by the app container and its Kubernetes Service. |

<a id="apps_app_containerport"></a>
#### `containerPort`

| **Property** | **Details** |
| --- | --- |
| Type | `integer` |
| Required | No |
| Default | — |

**Description:** Port exposed by the app container and its Kubernetes Service.
| [env](#apps_app_env) | `object` | No | — | Additional environment variables with literal string values. |

<a id="apps_app_env"></a>
#### `env`

| **Property** | **Details** |
| --- | --- |
| Type | `object` |
| Required | No |
| Default | — |
| Additional properties | Each additional property must conform to the schema (`string`) |

**Description:** Additional environment variables with literal string values.
| [extraContainers](#apps_app_extracontainers) | `array` | No | — | Additional Kubernetes Container objects added to the app Pod. See the [Kubernetes Container API](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/pod-v1/#Container). |

<a id="apps_app_extracontainers"></a>
#### `extraContainers`

| **Property** | **Details** |
| --- | --- |
| Type | `array` |
| Required | No |
| Default | — |

**Description:** Additional Kubernetes Container objects added to the app Pod. See the [Kubernetes Container API](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/pod-v1/#Container).
| [extraEgress](#apps_app_extraegress) | `array` | No | — | Additional Cilium egress rules for this app, using CiliumNetworkPolicy rule fields. `toSystems` is a chart-specific shorthand for another Drive system. See the [Cilium network policy documentation](https://docs.cilium.io/en/stable/security/policy/). |

<a id="apps_app_extraegress"></a>
#### `extraEgress`

| **Property** | **Details** |
| --- | --- |
| Type | `array` |
| Required | No |
| Default | — |

**Description:** Additional Cilium egress rules for this app, using CiliumNetworkPolicy rule fields. `toSystems` is a chart-specific shorthand for another Drive system. See the [Cilium network policy documentation](https://docs.cilium.io/en/stable/security/policy/).
| [extraIngress](#apps_app_extraingress) | `array` | No | — | Additional Cilium ingress rules for this app. `fromSystems` is a chart-specific shorthand for another Drive system. See the [Cilium network policy documentation](https://docs.cilium.io/en/stable/security/policy/). |

<a id="apps_app_extraingress"></a>
#### `extraIngress`

| **Property** | **Details** |
| --- | --- |
| Type | `array` |
| Required | No |
| Default | — |

**Description:** Additional Cilium ingress rules for this app. `fromSystems` is a chart-specific shorthand for another Drive system. See the [Cilium network policy documentation](https://docs.cilium.io/en/stable/security/policy/).
| [httpGateway](#apps_app_httpgateway) | `object` | No | See nested values | Enable platform HTTP routes for the app. Each enabled exposure creates a route at /<app-name>. |

<a id="apps_app_httpgateway"></a>
#### `httpGateway`

| **Property** | **Details** |
| --- | --- |
| Type | `object` |
| Required | No |
| Default | See nested values |
| Additional properties | Not allowed |

**Description:** Enable platform HTTP routes for the app. Each enabled exposure creates a route at /<app-name>.

| Property | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| [gated](#apps_app_httpgateway_gated) | `boolean` | No | — | Expose this app through the gated HTTP gateway. |

<a id="apps_app_httpgateway_gated"></a>
##### `gated`

| **Property** | **Details** |
| --- | --- |
| Type | `boolean` |
| Required | No |
| Default | — |

**Description:** Expose this app through the gated HTTP gateway.
| [private](#apps_app_httpgateway_private) | `boolean` | No | — | Expose this app through the private HTTP gateway. |

<a id="apps_app_httpgateway_private"></a>
##### `private`

| **Property** | **Details** |
| --- | --- |
| Type | `boolean` |
| Required | No |
| Default | — |

**Description:** Expose this app through the private HTTP gateway.
| [public](#apps_app_httpgateway_public) | `boolean` | No | — | Expose this app through the public HTTP gateway. |

<a id="apps_app_httpgateway_public"></a>
##### `public`

| **Property** | **Details** |
| --- | --- |
| Type | `boolean` |
| Required | No |
| Default | — |

**Description:** Expose this app through the public HTTP gateway.
| [nats](#apps_app_nats) | `object` | No | See nested values | Enable NATS access for this app. It connects to the account named by global.systemName. |

<a id="apps_app_nats"></a>
#### `nats`

| **Property** | **Details** |
| --- | --- |
| Type | `object` |
| Required | No |
| Default | See nested values |
| Additional properties | Not allowed |

**Description:** Enable NATS access for this app. It connects to the account named by global.systemName.

| Property | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| [user](#apps_app_nats_user) | `object` | No | See nested values | Configure access controls for this app's NATS identity. |

<a id="apps_app_nats_user"></a>
##### `user`

| **Property** | **Details** |
| --- | --- |
| Type | `object` |
| Required | No |
| Default | See nested values |
| Additional properties | Not allowed |

**Description:** Configure access controls for this app's NATS identity.

| Property | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| [permissions](#apps_app_nats_user_permissions) | `object` | No | See nested values | Configure subject access and request reply behavior for this app's NATS user. |

<a id="apps_app_nats_user_permissions"></a>
###### `permissions`

| **Property** | **Details** |
| --- | --- |
| Type | `object` |
| Required | No |
| Default | See nested values |
| Additional properties | Not allowed |

**Description:** Configure subject access and request reply behavior for this app's NATS user.

| Property | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| [publish](#apps_app_nats_user_permissions_publish) | `object` | No | See nested values | — |

<a id="apps_app_nats_user_permissions_publish"></a>
###### `publish`

| **Property** | **Details** |
| --- | --- |
| Type | `object` |
| Required | No |
| Default | See nested values |
| Additional properties | Not allowed |

| Property | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| [allow](#apps_app_nats_user_permissions_publish_allow) | `array` | No | — | — |

<a id="apps_app_nats_user_permissions_publish_allow"></a>
###### `allow`

| **Property** | **Details** |
| --- | --- |
| Type | `array` |
| Required | Yes |
| Default | — |
| Constraints | minItems=1 |
| [responses](#apps_app_nats_user_permissions_responses) | `object` | No | See nested values | — |

<a id="apps_app_nats_user_permissions_responses"></a>
###### `responses`

| **Property** | **Details** |
| --- | --- |
| Type | `object` |
| Required | No |
| Default | See nested values |
| Additional properties | Not allowed |

| Property | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| [enabled](#apps_app_nats_user_permissions_responses_enabled) | `boolean` | No | — | — |

<a id="apps_app_nats_user_permissions_responses_enabled"></a>
###### `enabled`

| **Property** | **Details** |
| --- | --- |
| Type | `boolean` |
| Required | No |
| Default | — |
| [max](#apps_app_nats_user_permissions_responses_max) | `integer` | No | — | — |

<a id="apps_app_nats_user_permissions_responses_max"></a>
###### `max`

| **Property** | **Details** |
| --- | --- |
| Type | `integer` |
| Required | No |
| Default | — |
| Constraints | minimum=1 |
| [ttlSeconds](#apps_app_nats_user_permissions_responses_ttlseconds) | `integer` | No | — | — |

<a id="apps_app_nats_user_permissions_responses_ttlseconds"></a>
###### `ttlSeconds`

| **Property** | **Details** |
| --- | --- |
| Type | `integer` |
| Required | No |
| Default | — |
| Constraints | maximum=9223372036, minimum=1 |
| [subscribe](#apps_app_nats_user_permissions_subscribe) | `object` | No | See nested values | — |

<a id="apps_app_nats_user_permissions_subscribe"></a>
###### `subscribe`

| **Property** | **Details** |
| --- | --- |
| Type | `object` |
| Required | No |
| Default | See nested values |
| Additional properties | Not allowed |

| Property | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| [allow](#apps_app_nats_user_permissions_subscribe_allow) | `array` | No | — | — |

<a id="apps_app_nats_user_permissions_subscribe_allow"></a>
###### `allow`

| **Property** | **Details** |
| --- | --- |
| Type | `array` |
| Required | Yes |
| Default | — |
| Constraints | minItems=1 |
| [postgres](#apps_app_postgres) | `object` | No | See nested values | PostgreSQL connection bindings for databases provisioned separately, typically by drive-system-postgres. |

<a id="apps_app_postgres"></a>
#### `postgres`

| **Property** | **Details** |
| --- | --- |
| Type | `object` |
| Required | No |
| Default | See nested values |
| Additional properties | Not allowed |

**Description:** PostgreSQL connection bindings for databases provisioned separately, typically by drive-system-postgres.

| Property | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| [bindings](#apps_app_postgres_bindings) | `array` | No | — | Database connections to inject as app environment variables. Each binding targets a cluster and database resource. |

<a id="apps_app_postgres_bindings"></a>
##### `bindings`

| **Property** | **Details** |
| --- | --- |
| Type | `array` |
| Required | No |
| Default | — |

**Description:** Database connections to inject as app environment variables. Each binding targets a cluster and database resource.
| [probes](#apps_app_probes) | `object` | No | See nested values | Readiness, liveness, and startup probes. Full Kubernetes probe objects override the shorthand HTTP paths. |

<a id="apps_app_probes"></a>
#### `probes`

| **Property** | **Details** |
| --- | --- |
| Type | `object` |
| Required | No |
| Default | See nested values |
| Additional properties | Not allowed |

**Description:** Readiness, liveness, and startup probes. Full Kubernetes probe objects override the shorthand HTTP paths.

| Property | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| [liveness](#apps_app_probes_liveness) | `string` | No | — | HTTP path for a generated liveness probe and, for Java apps, the default startup probe. Overridden by livenessProbe for liveness. |

<a id="apps_app_probes_liveness"></a>
##### `liveness`

| **Property** | **Details** |
| --- | --- |
| Type | `string` |
| Required | No |
| Default | — |

**Description:** HTTP path for a generated liveness probe and, for Java apps, the default startup probe. Overridden by livenessProbe for liveness.
| [livenessProbe](#apps_app_probes_livenessprobe) | `any` | No | — | Full Kubernetes livenessProbe object; takes precedence over the liveness path. |

<a id="apps_app_probes_livenessprobe"></a>
##### `livenessProbe`

| **Property** | **Details** |
| --- | --- |
| Type | `any` |
| Required | No |
| Default | — |
| Schema reference | [Referenced schema](https://raw.githubusercontent.com/yannh/kubernetes-json-schema/8df8a883b68a24a104b4a9e43c1288090ae60b3b/master-standalone-strict/probe-v1.json) |

**Description:** Full Kubernetes livenessProbe object; takes precedence over the liveness path.
| [port](#apps_app_probes_port) | `integer` | No | — | Port used by shorthand HTTP readiness and liveness probes; defaults to appDefaults.probes.port. |

<a id="apps_app_probes_port"></a>
##### `port`

| **Property** | **Details** |
| --- | --- |
| Type | `integer` |
| Required | No |
| Default | — |

**Description:** Port used by shorthand HTTP readiness and liveness probes; defaults to appDefaults.probes.port.
| [readiness](#apps_app_probes_readiness) | `string` | No | — | HTTP path for a generated readiness probe. Overridden by readinessProbe when provided. |

<a id="apps_app_probes_readiness"></a>
##### `readiness`

| **Property** | **Details** |
| --- | --- |
| Type | `string` |
| Required | No |
| Default | — |

**Description:** HTTP path for a generated readiness probe. Overridden by readinessProbe when provided.
| [readinessProbe](#apps_app_probes_readinessprobe) | `any` | No | — | Full Kubernetes readinessProbe object; takes precedence over the readiness path. |

<a id="apps_app_probes_readinessprobe"></a>
##### `readinessProbe`

| **Property** | **Details** |
| --- | --- |
| Type | `any` |
| Required | No |
| Default | — |
| Schema reference | [Referenced schema](https://raw.githubusercontent.com/yannh/kubernetes-json-schema/8df8a883b68a24a104b4a9e43c1288090ae60b3b/master-standalone-strict/probe-v1.json) |

**Description:** Full Kubernetes readinessProbe object; takes precedence over the readiness path.
| [startup](#apps_app_probes_startup) | `any` | No | — | Optional full Kubernetes startupProbe object. For Java apps, replaces the default startup probe derived from the liveness path. |

<a id="apps_app_probes_startup"></a>
##### `startup`

| **Property** | **Details** |
| --- | --- |
| Type | `any` |
| Required | No |
| Default | — |
| Schema reference | [Referenced schema](https://raw.githubusercontent.com/yannh/kubernetes-json-schema/8df8a883b68a24a104b4a9e43c1288090ae60b3b/master-standalone-strict/probe-v1.json) |

**Description:** Optional full Kubernetes startupProbe object. For Java apps, replaces the default startup probe derived from the liveness path.
| [profile](#apps_app_profile) | `string` | No | — | Runtime profile used for profile-specific workload behavior. `java` selects Java security identity and enables a default startup probe when a liveness path is configured. |

<a id="apps_app_profile"></a>
#### `profile`

| **Property** | **Details** |
| --- | --- |
| Type | `string` |
| Required | No |
| Default | — |
| Constraints | minLength=1 |

**Description:** Runtime profile used for profile-specific workload behavior. `java` selects Java security identity and enables a default startup probe when a liveness path is configured.
| [replicaCount](#apps_app_replicacount) | `integer` | No | — | Fixed Deployment replica count when autoscaling is not configured. Defaults to 1; set 0 to scale down. |

<a id="apps_app_replicacount"></a>
#### `replicaCount`

| **Property** | **Details** |
| --- | --- |
| Type | `integer` |
| Required | No |
| Default | — |
| Constraints | minimum=0 |

**Description:** Fixed Deployment replica count when autoscaling is not configured. Defaults to 1; set 0 to scale down.
| [resources](#apps_app_resources) | `object` | No | — | Kubernetes container resource requests and limits. Replaces appDefaults.resources when set. |

<a id="apps_app_resources"></a>
#### `resources`

| **Property** | **Details** |
| --- | --- |
| Type | `object` |
| Required | No |
| Default | — |
| Additional properties | Any type allowed |

**Description:** Kubernetes container resource requests and limits. Replaces appDefaults.resources when set.
| [secretEnv](#apps_app_secretenv) | `object` | No | — | Environment variables sourced from existing Kubernetes Secrets, keyed by environment variable name. |

<a id="apps_app_secretenv"></a>
#### `secretEnv`

| **Property** | **Details** |
| --- | --- |
| Type | `object` |
| Required | No |
| Default | — |
| Additional properties | Each additional property must conform to the schema (`object`) |

**Description:** Environment variables sourced from existing Kubernetes Secrets, keyed by environment variable name.

| Property | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| [secretEnv-object](#apps_app_secretenv_secretenv_object) | `object` | No | See nested values | — |

<a id="apps_app_secretenv_secretenv_object"></a>
##### `secretEnv-object`

| **Property** | **Details** |
| --- | --- |
| Type | `object` |
| Required | No |
| Default | See nested values |
| Additional properties | Not allowed |

| Property | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| [secretKey](#apps_app_secretenv_secretenv_object_secretkey) | `string` | No | — | Key within the Secret to expose as the environment variable value. |

<a id="apps_app_secretenv_secretenv_object_secretkey"></a>
###### `secretKey`

| **Property** | **Details** |
| --- | --- |
| Type | `string` |
| Required | Yes |
| Default | — |
| Constraints | minLength=1 |

**Description:** Key within the Secret to expose as the environment variable value.
| [secretName](#apps_app_secretenv_secretenv_object_secretname) | `string` | No | — | Name of the existing Kubernetes Secret containing the environment variable value. |

<a id="apps_app_secretenv_secretenv_object_secretname"></a>
###### `secretName`

| **Property** | **Details** |
| --- | --- |
| Type | `string` |
| Required | Yes |
| Default | — |
| Constraints | minLength=1 |

**Description:** Name of the existing Kubernetes Secret containing the environment variable value.
| [strategy](#apps_app_strategy) | `object` | No | — | Optional Kubernetes Deployment strategy override. See the [Kubernetes Deployment API](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/deployment-v1/#DeploymentSpec). |

<a id="apps_app_strategy"></a>
#### `strategy`

| **Property** | **Details** |
| --- | --- |
| Type | `object` |
| Required | No |
| Default | — |
| Additional properties | Any type allowed |

**Description:** Optional Kubernetes Deployment strategy override. See the [Kubernetes Deployment API](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/deployment-v1/#DeploymentSpec).
| [tag](#apps_app_tag) | `string` | No | — | Container image tag for this app. Required; the image repository is derived from the app name and global registry. |

<a id="apps_app_tag"></a>
#### `tag`

| **Property** | **Details** |
| --- | --- |
| Type | `string` |
| Required | Yes |
| Default | — |
| Constraints | minLength=1 |

**Description:** Container image tag for this app. Required; the image repository is derived from the app name and global registry.
| [volumeMounts](#apps_app_volumemounts) | `object` | No | See nested values | Mount existing Secrets/ConfigMaps or emptyDir volumes. Only these chart-supported volume types are accepted. |

<a id="apps_app_volumemounts"></a>
#### `volumeMounts`

| **Property** | **Details** |
| --- | --- |
| Type | `object` |
| Required | No |
| Default | See nested values |
| Additional properties | Not allowed |

**Description:** Mount existing Secrets/ConfigMaps or emptyDir volumes. Only these chart-supported volume types are accepted.

| Property | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| [configMap](#apps_app_volumemounts_configmap) | `object` | No | — | Existing Kubernetes ConfigMaps mounted as volumes, keyed by ConfigMap name. |

<a id="apps_app_volumemounts_configmap"></a>
##### `configMap`

| **Property** | **Details** |
| --- | --- |
| Type | `object` |
| Required | No |
| Default | — |
| Additional properties | Each additional property must conform to the schema (`object`) |

**Description:** Existing Kubernetes ConfigMaps mounted as volumes, keyed by ConfigMap name.

| Property | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| [configMap-object](#apps_app_volumemounts_configmap_configmap_object) | `object` | No | See nested values | — |

<a id="apps_app_volumemounts_configmap_configmap_object"></a>
###### `configMap-object`

| **Property** | **Details** |
| --- | --- |
| Type | `object` |
| Required | No |
| Default | See nested values |
| Additional properties | Not allowed |

| Property | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| [optional](#apps_app_volumemounts_configmap_configmap_object_optional) | `boolean` | No | — | Allow the pod to start if the Secret or ConfigMap does not exist. Defaults to false. |

<a id="apps_app_volumemounts_configmap_configmap_object_optional"></a>
###### `optional`

| **Property** | **Details** |
| --- | --- |
| Type | `boolean` |
| Required | No |
| Default | — |

**Description:** Allow the pod to start if the Secret or ConfigMap does not exist. Defaults to false.
| [path](#apps_app_volumemounts_configmap_configmap_object_path) | `string` | No | — | Absolute container path where the existing Secret or ConfigMap is mounted. |

<a id="apps_app_volumemounts_configmap_configmap_object_path"></a>
###### `path`

| **Property** | **Details** |
| --- | --- |
| Type | `string` |
| Required | Yes |
| Default | — |
| Constraints | pattern="^/" |

**Description:** Absolute container path where the existing Secret or ConfigMap is mounted.
| [ephemeralDirs](#apps_app_volumemounts_ephemeraldirs) | `object` | No | — | Named emptyDir volumes for writable ephemeral storage in the otherwise read-only app container. |

<a id="apps_app_volumemounts_ephemeraldirs"></a>
##### `ephemeralDirs`

| **Property** | **Details** |
| --- | --- |
| Type | `object` |
| Required | No |
| Default | — |
| Additional properties | Each additional property must conform to the schema (`object`) |

**Description:** Named emptyDir volumes for writable ephemeral storage in the otherwise read-only app container.

| Property | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| [ephemeralDir](#apps_app_volumemounts_ephemeraldirs_ephemeraldir) | `object` | No | See nested values | — |

<a id="apps_app_volumemounts_ephemeraldirs_ephemeraldir"></a>
###### `ephemeralDir`

| **Property** | **Details** |
| --- | --- |
| Type | `object` |
| Required | No |
| Default | See nested values |
| Additional properties | Not allowed |

| Property | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| [path](#apps_app_volumemounts_ephemeraldirs_ephemeraldir_path) | `string` | No | — | Absolute path inside the container where the emptyDir volume is mounted. |

<a id="apps_app_volumemounts_ephemeraldirs_ephemeraldir_path"></a>
###### `path`

| **Property** | **Details** |
| --- | --- |
| Type | `string` |
| Required | Yes |
| Default | — |
| Constraints | pattern="^/" |

**Description:** Absolute path inside the container where the emptyDir volume is mounted.
| [secrets](#apps_app_volumemounts_secrets) | `object` | No | — | Existing Kubernetes Secrets mounted as volumes, keyed by Secret name. |

<a id="apps_app_volumemounts_secrets"></a>
##### `secrets`

| **Property** | **Details** |
| --- | --- |
| Type | `object` |
| Required | No |
| Default | — |
| Additional properties | Each additional property must conform to the schema (`object`) |

**Description:** Existing Kubernetes Secrets mounted as volumes, keyed by Secret name.

| Property | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| [secret](#apps_app_volumemounts_secrets_secret) | `object` | No | See nested values | — |

<a id="apps_app_volumemounts_secrets_secret"></a>
###### `secret`

| **Property** | **Details** |
| --- | --- |
| Type | `object` |
| Required | No |
| Default | See nested values |
| Additional properties | Not allowed |

| Property | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| [optional](#apps_app_volumemounts_secrets_secret_optional) | `boolean` | No | — | Allow the pod to start if the Secret or ConfigMap does not exist. Defaults to false. |

<a id="apps_app_volumemounts_secrets_secret_optional"></a>
###### `optional`

| **Property** | **Details** |
| --- | --- |
| Type | `boolean` |
| Required | No |
| Default | — |

**Description:** Allow the pod to start if the Secret or ConfigMap does not exist. Defaults to false.
| [path](#apps_app_volumemounts_secrets_secret_path) | `string` | No | — | Absolute container path where the existing Secret or ConfigMap is mounted. |

<a id="apps_app_volumemounts_secrets_secret_path"></a>
###### `path`

| **Property** | **Details** |
| --- | --- |
| Type | `string` |
| Required | Yes |
| Default | — |
| Constraints | pattern="^/" |

**Description:** Absolute container path where the existing Secret or ConfigMap is mounted.

<a id="drive_system_lib"></a>
## `drive-system-lib`

| **Property** | **Details** |
| --- | --- |
| Type | `object` |
| Required | No |
| Default | — |
| Additional properties | Any type allowed |

**Description:** Internal values namespace for the drive-system-lib dependency; not intended for chart consumers.

<a id="fullnameoverride"></a>
## `fullnameOverride`

| **Property** | **Details** |
| --- | --- |
| Type | `string` |
| Required | No |
| Default | `""` |

**Description:** Override the generated names for chart-owned resources. Leave empty to use chart defaults.

<a id="global"></a>
## `global`

| **Property** | **Details** |
| --- | --- |
| Type | `object` |
| Required | No |
| Default | See nested values |
| Additional properties | Not allowed |

**Description:** Global values are populated by the Drive platform during deployment and by drv during local development. Users do not need to set them manually.

| Property | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| [domain](#global_domain) | `string` | No | `""` | Base domain used by platform-managed routes and integrations. |

<a id="global_domain"></a>
### `domain`

| **Property** | **Details** |
| --- | --- |
| Type | `string` |
| Required | No |
| Default | `""` |

**Description:** Base domain used by platform-managed routes and integrations.
| [image](#global_image) | `object` | No | See nested values | Shared container image registry configuration. |

<a id="global_image"></a>
### `image`

| **Property** | **Details** |
| --- | --- |
| Type | `object` |
| Required | No |
| Default | See nested values |
| Additional properties | Not allowed |

**Description:** Shared container image registry configuration.

| Property | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| [registry](#global_image_registry) | `string` | No | `""` | Container image registry prefix supplied by the delivery platform. |

<a id="global_image_registry"></a>
#### `registry`

| **Property** | **Details** |
| --- | --- |
| Type | `string` |
| Required | No |
| Default | `""` |

**Description:** Container image registry prefix supplied by the delivery platform.
| [stageName](#global_stagename) | `string` | No | `""` | Stage identifier supplied by the delivery platform. |

<a id="global_stagename"></a>
### `stageName`

| **Property** | **Details** |
| --- | --- |
| Type | `string` |
| Required | No |
| Default | `""` |

**Description:** Stage identifier supplied by the delivery platform.
| [systemName](#global_systemname) | `string` | No | `""` | Drive system identifier supplied by the delivery platform. |

<a id="global_systemname"></a>
### `systemName`

| **Property** | **Details** |
| --- | --- |
| Type | `string` |
| Required | No |
| Default | `""` |

**Description:** Drive system identifier supplied by the delivery platform.
| [verticalId](#global_verticalid) | `string` | No | `""` | Drive vertical identifier supplied by the delivery platform. |

<a id="global_verticalid"></a>
### `verticalId`

| **Property** | **Details** |
| --- | --- |
| Type | `string` |
| Required | No |
| Default | `""` |

**Description:** Drive vertical identifier supplied by the delivery platform.

<a id="otel"></a>
## `otel`

| **Property** | **Details** |
| --- | --- |
| Type | `object` |
| Required | No |
| Default | See nested values |
| Additional properties | Not allowed |

**Description:** Controls OpenTelemetry export configuration for all app workloads.

| Property | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| [enabled](#otel_enabled) | `boolean` | No | `true` | Enable app telemetry export to the platform node collector. |

<a id="otel_enabled"></a>
### `enabled`

| **Property** | **Details** |
| --- | --- |
| Type | `boolean` |
| Required | No |
| Default | `true` |

**Description:** Enable app telemetry export to the platform node collector.

<a id="secrets"></a>
## `secrets`

| **Property** | **Details** |
| --- | --- |
| Type | `array` |
| Required | No |
| Default | `[]` |

**Description:** Names of platform-managed secrets to expose to this system as ExternalSecrets. Secret values are not supplied in chart values.
