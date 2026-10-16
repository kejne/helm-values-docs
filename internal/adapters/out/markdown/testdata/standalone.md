# Helm values API

## Contents

- [image](#image)
  - [tag](#image_tag)
- [missing](#missing)
- [replicas](#replicas)

| **Schema** | **Details** |
| --- | --- |
| Type | `object` |


<a id="image"></a>
## `image`

| **Property** | **Details** |
| --- | --- |
| Type | `object` |
| Required | No |
| Default | See nested values |

| Property | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| [tag](#image_tag) | `string` | No | `"stable"` | Image \| tag |

<a id="image_tag"></a>
### `tag`

| **Property** | **Details** |
| --- | --- |
| Type | `string` |
| Required | No |
| Default | `"stable"` |

**Description:** Image | tag

<a id="missing"></a>
## `missing`

| **Property** | **Details** |
| --- | --- |
| Type | `string` |
| Required | No |
| Default | — |

<a id="replicas"></a>
## `replicas`

| **Property** | **Details** |
| --- | --- |
| Type | `integer` |
| Required | Yes |
| Default | `3` |
| Constraints | minimum=1 |

**Description:** Number of pods
