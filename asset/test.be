backend x509 "main_ca" {
  type = "file"

  cert = "asset/cert-ec.pem"
  key  = "asset/key-ec.pem"
}

policy "kubelet" {
  verify gcp {}

  produce file "abc.txt" {
    from = "asset/qq"
  }

  produce file "def.txt" {
    from = "asset/qq2"
  }

  produce x509 "kubelet" {
    backend = backend.x509.main_ca

    common_name = "Kubelet User certificate"
    alt_dns = [req.fqdn]
    alt_ips = req.ips
  }
}
