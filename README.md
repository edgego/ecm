
<div align="center">
  <h1>Edge Cloud Manager</h1>
  <span>English / </span> <a href="https://github.com/edgego/">Simplified Chinese</a>
</div>

<hr />

## What is ECM
 ECM desighed for industry edge cloud , a small and low cost tool easy to  manager industry edge clouds. ECM is  based on CNCF project k3s. ECM can run on different platform, support winows, linux ,Mac OS.

## Key Features

- Cross platform, ECM can create edge cloud on amd64,arm64,arm device. 
- .
- .
- Low cost,low carbon, small memory size and hard disk size.
- Simplify operations by UI dashboard.
- .

## Providers

Now supports the following providers, we encourage submitting PR contribution for more providers:

- [aws](docs/i18n/en_us/aws/README.md) - Bootstrap cloud onto Amazon EC2
- [google](docs/i18n/en_us/google/README.md) - Bootstrap cloud onto Google Compute Engine
- [alibaba](docs/i18n/en_us/alibaba/README.md) - Bootstrap cloud onto Alibaba ECS
- [tencent](docs/i18n/en_us/tencent/README.md) - Bootstrap cloud onto Tencent CVM
- [k3d](docs/i18n/en_us/k3d/README.md) - Bootstrap cloud onto Local Machine
- [native](docs/i18n/en_us/native/README.md) - Bootstrap cloud onto any VM

## Quick Start

 Run with cli:

```bash

# The commands will start ecm daemon with an interactionable UI.
ecm -d serve --bind-port 8080
```
![image](https://user-images.githubusercontent.com/80612608/174299658-a645f7a2-6e6a-429e-bd88-56febf1256c4.png)

![image](https://user-images.githubusercontent.com/80612608/174299845-08435f58-b8be-41b7-bb02-49fb9d7639a2.png)



# License

Copyright (c) 2022 EdgeGo

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)
=======
# ecm
