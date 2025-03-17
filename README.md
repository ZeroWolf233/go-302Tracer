## 302Tracer

A small tool to trace the 302 redirect.

## Features

 - Trace the redirect of the request
 - Send multiple requests to see the different redirect results
 - Custom request threads
 - Custom User Agent

## Installation

### Release

You can download pre-compiled binaries for your platform from the [GitHub Releases](https://github.com/ZeroWolf233/go-302Tracer/releases) page.

### Building from Source code

Make sure you have Go installed,

then run:

```bash
git clone https://github.com/ZeroWolf233/go-302Tracer
cd go-302Tracer
go build
```

## Usage

The basic ues example:

```bash
302Tracer https://example.com/download/miku.jpg
```

Using with some flags:
```bash
302Tracer -t 10 -w 32 -r 3 -ua PCL/2.9.1 https://bmclapi2.bangbang93.com/version/1.21.4/client
```

You can see the different redirect address in different requests

## Flags
| name | defult value | explaination |
|------|--------------|--------------|
| -t   | 1            | request times |
| -w   | 1            | request therads|
| -r   | 1            | the rest time between different requests|
| -ua  | Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36 Edg/128.0.0.0" | User Agent

## Output Example

```bash
zerowolf@ZeroWolf-PC:~# 302Tracer -w 32 -ua PCL2/2.9.1 -r 1 https://bmclapi2.bangbang93.com/version/1.21.4/client 

request 1 had been redirect to: https://bmclapi2.bangbang93.com/v1/objects/a7e5a6024bfd3cd614625aa05629adf760020304/client.jar

request 1 had been redirect to: https://6717a58c9efd532c5ed1fb20.openbmclapi.933.moe:30003/download/a7e5a6024bfd3cd614625aa05629adf760020304?name=client.jar&s=NzD-z1d8KUPR38y4K1JDUbK1EEs&e=m8d6knf0

request 1 had been redirect to: https://media-jsyz-fy-home.js6oss.ctyunxs.cn/FAMILYCLOUD/a089c955-64a7-4c09-bb7c-65795a872eb8?x-amz-CLIENTTYPEIN=PC&AWSAccessKeyId=0Lg7dAq3ZfHvePP8DKEU&x-amz-limitrate=51200&x-amz-UID=462299132&x-amz-APPID=93005&response-content-disposition=attachment%3Bfilename%3D%22a7e5a6024bfd3cd614625aa05629adf760020304%22%3Bfilename*%3DUTF-8%27%27a7e5a6024bfd3cd614625aa05629adf760020304&x-amz-OPERID=108208554&x-amz-CLIENTNETWORK=UNKNOWN&x-amz-CLOUDTYPEIN=FAMILY&Signature=sJs5j3ODRyuErcHdQsOz%2BkC16dA%3D&Expires=1742232901&x-amz-FSIZE=28335587&x-amz-UFID=325461175865163957


request 1 succeed, final address: https://media-jsyz-fy-home.js6oss.ctyunxs.cn/FAMILYCLOUD/a089c955-64a7-4c09-bb7c-65795a872eb8?x-amz-CLIENTTYPEIN=PC&AWSAccessKeyId=0Lg7dAq3ZfHvePP8DKEU&x-amz-limitrate=51200&x-amz-UID=462299132&x-amz-APPID=93005&response-content-disposition=attachment%3Bfilename%3D%22a7e5a6024bfd3cd614625aa05629adf760020304%22%3Bfilename*%3DUTF-8%27%27a7e5a6024bfd3cd614625aa05629adf760020304&x-amz-OPERID=108208554&x-amz-CLIENTNETWORK=UNKNOWN&x-amz-CLOUDTYPEIN=FAMILY&Signature=sJs5j3ODRyuErcHdQsOz%2BkC16dA%3D&Expires=1742232901&x-amz-FSIZE=28335587&x-amz-UFID=325461175865163957

usage 0.46 seconds, file size: 27.02 MB, speed: 58.22 MB/s

----------------------------------------
Done, please check the results.
```

## License

GNU General Public License

## Acknowledgements
 - [Mxmilu666](https://milu.ink/) Provide the Readme Example.