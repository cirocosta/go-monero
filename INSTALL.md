# Install

## Using Go

If you already have [Go] installed, you can use to Go toolchain to build from
source and install it for you (under `$GOPATH/bin`):

```console
$ GO111MODULE=on go get github.com/cirocosta/go-monero/cmd/monero
go: downloading github.com/cirocosta/go-monero v0.0.3

$ monero --help
Daemon, Wallet, and p2p command line monero CLI

Usage:
  monero [command]
...
```

Note that this will install the latest tagged release (not necessarily the
latest code).


## From releases

In the [releases page] you'll find the pre-compiled releases for each platform.


### "yeah yeah, I trust the interwebs" mode

```bash
export VERSION=0.0.3

curl -SL -o- https://github.com/cirocosta/go-monero/releases/download/v$VERSION/go-monero_$VERSION_Linux_x86_64.tar.gz | \
  tar xvzf monero
mv ./monero /usr/local/bin
```


### "trust, but verify" mode


1. fetch my public key and make sure it matches the expected fingerprint

```console
$ curl -SOL https://utxo.com.br/pgp/public-key.txt
```

now, using `gpg`, derive the fingerprint. it should match the one advertised by
me: `9CD1 1313 8578 59CC 0FAD  E93B 6B93 177A 62D0 1DB8` (should be the same as
you can find under my personal account on Twitter: http://twitter.com/cirowrc).


```console
$ gpg --keyid-format long --with-fingerprint ./public-key.txt
pub   rsa3072/6B93177A62D01DB8 2021-07-19 [SC] [expires: 2023-07-19]
      Key fingerprint = 9CD1 1313 8578 59CC 0FAD  E93B 6B93 177A 62D0 1DB8
```

then, import into the key to the keyring so it can be used to validate that I
indeed signed the content advertised.

```console
$ gpg --import ./public-key.txt
gpg: key 6B93177A62D01DB8: public key "..." imported
gpg: Total number processed: 1
gpg:               imported: 1
```


2. download the archive for your platform as well as the checksums

```console
$ curl -SOL https://github.com/cirocosta/go-monero/releases/download/v0.0.3/go-monero_0.0.3_Linux_x86_64.tar.gz
$ curl -SOL https://github.com/cirocosta/go-monero/releases/download/v0.0.3/checksums.txt.asc
```


3. verify that you can trust the checksums (that it has been generated and
   not tampered with), and then verify that the assets you downloaded are what
   they supposed to be

```console
$ gpg --verify ./checksums.txt.asc
gpg: Signature made Mon 19 Jul 2021 02:10:42 PM EDT
gpg:                using RSA key 9CD11313857859CC0FADE93B6B93177A62D01DB8
gpg: Good signature from "Ciro ...
```


4. verify that the tarball is what it should be

compute the digest of the tarball

```console
$ sha256sum ./go-monero_0.0.3_Linux_x86_64.tar.gz
e2b2214c9371fe3c0333cca7feff3554c56d8d0f377180e39ff50d332639c22d  ./go-monero_0.0.3_Linux_x86_64.tar.gz


see that it matches what you found in the signed checksums file

$ cat ./checksums.txt.asc | grep e2b2214c9371fe3c0333cca7feff3554c56d8d0f377180e39ff50d332639c22d
e2b2214c9371fe3c0333cca7feff3554c56d8d0f377180e39ff50d332639c22d  go-monero_0.0.3_Linux_x86_64.tar.gz
```

5. install

```console
$ tar xvzf ./go-monero_0.0.3_Linux_x86_64.tar.gz monero
$ mv monero /usr/local/bin

$ monero version
0.0.3 45b42a1fae19c4fa950d159cde2f954b49365d93
```


[Go]: https://golang.org/dl/
[releases page]: https://github.com/cirocosta/go-monero/releases
