
<div align="center">
  <h1>Edge Cloud Manager</h1>
  <span>English / </span> <a href="https://github.com/edgego/">Simplified Chinese</a>
</div>

<hr />

## What is ECM
 ECM is designed for industry edge cloud , a small and low cost tool easy to  manager industry edge clouds. ECM is based on CNCF project k3s. ECM can run on different platform(x86-amd64, arm,arm64), support winows, linux ,Mac OS. ECM can deploy edge cloud to nvidia Jetson device and Raspberry Pi device.

## Key Features

- Cross platform, ECM can create edge cloud on amd64,arm64,arm device. 
- ECM support english and simple chinese language, according to browswer's language setting.
- Low cost,low carbon, small memory size(500Mi or above) and hard disk size.
- Simplify operations by UI dashboard.
- Integration with cloud dashboard
- Integration with EdgeGo system that is IOT system

## Quick Start

 Run  ecm :

```bash
firstly to get ssh key with command: ssh-keygen -t rsa

Run with command line:
  create edge cloud:
     ecm_win-amd64.exe  -d create \
     --region edge cloud location \
     --cluster \
     --enable dashboard \
     --name demo-test \
     --ssh-user root \
     --ssh-password ****** \
     --ssh-port 22 \
     --master-ips 192.168.1.151
  
  delete edge cloud:
     ecm_win-amd64.exe  -d delete --name demo-test
  
  join node to edge cloud:
    ecm_win-amd64.exe -d join \
    --name demo-test \
    --ip 192.168.1.151 \
    --ssh-user root \
    --ssh-password ****** \
    --worker-ips 192.168.1.152

Run with ui:
  # The commands will start ecm daemon with an interactionable UI.
  ecm -d serve --bind-port 8080
```


<img width="1903" alt="image" src="https://github.com/edgego/ecm/assets/80612608/e6b42447-b812-4cae-bc83-4604631b195a">



# License

Copyright (c) 2022 EdgeGo

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)
=======
