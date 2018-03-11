dexctl
======

dexctl is an unofficial command line tool for
[dex](https://github.com/coreos/dex/), because there is no official one for dex
v2. This tool is made available as-is: no guarantees of compatibility.


Quickstart
----------

1. Generate certificates using `bin/cert-gen`.

2. Create secret in same namespace as `dex` with the server.crt, server.key,
   and ca.crt.

3. Mount the secret into `/var/dex/certs`.

4. Enable gRPC on dex:

  ```
    grpc:
      addr: 0.0.0.0:5557
      tlsCert: /var/dex/certs/server.crt
      tlsKey: /var/dex/certs/server.key
      tlsClientCA: /var/dex/certs/ca.crt
  ```

5. Restart the pods if necessary.

6. Use kubectl to port-forward into one of the pods at port 5557.

7. Run dexctl like:

  ```
    ./dexctl -ca-cert certs/ca.crt -client-cert certs/client.crt -client-key certs/client.key
  ```

8. You can give dexctl a path to a YAML file. The YAML file looks like:

  ```
    id: "kubectl"
    name: "Kubernetes CLI (kubectl)"
    secret: "XXX-REDACTED-XXX"
    public: true
  ```

Refer to `type Client` in `vendor/github.com/coreos/dex/api/api.pb.go` to see
the full structure. The struct does not (as of March 2018) come with YAML tags;
the YAML keys should be lowercased in such a case (e.g., `redirecturis`).

