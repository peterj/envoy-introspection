# Envoy introspection



## Demo

1. Build the extension (.wasm file):

    ```
    make wasm
    ```

2. Start the upstream service:

    ```sh
    docker run -d -p 3030:80 kennethreitz/httpbin
    ```

3. Run the Envoy instance using func-e:

    ```sh
    func-e run -c envoy.yaml &
    ```

4. Send a request using the `intercept: 1` header to intercept the call:

    ```sh
    curl -H "intercept: 1" localhost:10000/get
    ```

Notice in the output the call is intercepted and a request is made to the same cluster and the callback function is called:

```sh
$ curl -H "intercept: 1" localhost:10000/1
[2021-11-30 01:02:33.159][14834][info][wasm] [source/extensions/common/wasm/context.cc:1167] wasm log: intercepting call!!
[2021-11-30 01:02:33.159][14834][info][wasm] [source/extensions/common/wasm/context.cc:1167] wasm log: cluster name: httpbin_1
[2021-11-30 01:02:33.163][14834][info][wasm] [source/extensions/common/wasm/context.cc:1167] wasm log: called callBack func
{
  "args": {},
  "headers": {
    "Accept": "*/*",
    "Host": "localhost:10000",
    "Intercept": "1",
    "User-Agent": "curl/7.64.0",
    "X-Envoy-Expected-Rq-Timeout-Ms": "15000",
    "X-Envoy-Original-Path": "/1"
  },
  "origin": "172.18.0.1",
  "url": "http://localhost:10000/get"
}
```