# Helm values API

| Path | Type | Required | Value | Description | Constraints |
| --- | --- | --- | --- | --- | --- |
| ["image"] | object | no | {"tag":"stable"} |  |  |
| ["image"]["tag"] | string | no | "stable" | Image \| tag |  |
| ["missing"] | string | no | — |  |  |
| ["replicas"] | integer | yes | 3 | Number of pods | minimum=1 |
