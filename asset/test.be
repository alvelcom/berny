backend x509_file "main_ca" {
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
    backend = "main_ca"

    common_name = "123"
    alt_dns = ["abc", "def"]
  }
}
