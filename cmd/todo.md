- https://github.com/kudelskisecurity/crystals-go !!!!!!!!!!! and vendor it
- Message pack is a good serialization protocol: https://msgpack.org/

# Big Two: Remove SHA-1 (it is broken) | I think I already removed it for sha256

- Logging --log + Progress bar with either of
  - https://github.com/schollz/progressbar (preferred)
  - https://github.com/cheggaaa/pb
- [x] fixed: created files were executable!
- Currently gets hopelessly stuck if s3 connection is interrupted - needs to display progress and use timeout wisely with re-tries
- is rand. properly seeded!
- Current working directory is used for the temporary file. Investigate: is that approach wise?

# Big One: Remove RSA

_<https://blog.trailofbits.com/2019/07/08/fuck-rsa/>_

Here at Trail of Bits we review a lot of code. From major open source
projects to exciting new proprietary software, we've seen it all. But
one common denominator in all of these systems is that for some
inexplicable reason people still seem to think RSA is a good
cryptosystem to use. Let me save you a bit of time and money and just
say outright---if you come to us with a codebase that uses RSA, you will
be paying for the hour of time required for us to explain why you should
stop using it.

RSA is an intrinsically fragile cryptosystem containing countless
foot-guns which the average software engineer cannot be expected to
avoid. Weak parameters can be difficult, if not impossible, to check,
and its poor performance compels developers to take risky shortcuts.
Even worse, padding oracle attacks remain rampant 20 years after they
were discovered. While it may be theoretically possible to implement RSA
correctly, decades of devastating attacks have proven that such a feat
may be unachievable in practice.

> MarshalPKCS1PrivateKey

As we mentioned above, just using RSA out of the box doesn't quite work.
For example, the RSA scheme laid out in the introduction would produce
identical ciphertexts if the same plaintext were ever encrypted more
than once. This is a problem, because it would allow an adversary to
infer the contents of the message from context without being able to
decrypt it. This is why we need to pad messages with some random bytes.
Unfortunately, the most widely used padding scheme, PKCS \#1 v1.5, is
often vulnerable to something called a padding oracle attack. For more details on
the attack, check out [this excellent
explainer](https://crypto.stackexchange.com/questions/12688/can-you-explain-bleichenbachers-cca-attack-on-pkcs1-v1-5).

TLS 1.3 no longer supports RSA so we can expect to see fewer of these attacks going forward, but as long as developers continue to use RSA in their own applications there will be padding oracle attacks.

# Error-Correcting Codes

- Keeping a hash of a byte string can allow re-building the string by random iteration until the value matches the hash.
- Hashing byte sets that are already ECC-hashed allows double-checking for errors by parent. Solve the first hash, check the solution against parent - no match - look for a second solution.
- This prevents ECC hashes themselves from being bit-rot - they could lay in a duplicate tree even as a separate file.
- https://innovation.vivint.com/introduction-to-reed-solomon-bc264d0794f8
  - https://github.com/maruel/rs
  - https://github.com/vivint/infectious
  - https://github.com/klauspost/reedsolomon
- http://blog.klauspost.com/blazingly-fast-reed-solomon-coding/

## Replacement

Trail of Bits recommends using
[Curve25519](https://en.wikipedia.org/wiki/Curve25519) for key exchange
and digital signatures. Encryption needs to be done using a protocol
called
[ECIES](https://en.wikipedia.org/wiki/Integrated_Encryption_Scheme)
which combines an elliptic curve key exchange with a symmetric
encryption algorithm. Curve25519 was designed to entirely prevent some
of the things that can go wrong with other curves, and is very
performant. Even better, it is implemented in
[libsodium](https://libsodium.gitbook.io/doc/), which has [easy-to-read
documentation](https://libsodium.gitbook.io/doc/public-key_cryptography/sealed_boxes)
and is [available for most
languages](https://libsodium.gitbook.io/doc/libsodium_users).

## Considerations

- https://github.com/jesseduffield/horcrux
- see how restic · Backups done right! https://restic.net/ does it
- also similar: https://github.com/FiloSottile/age
- Peer Keep looks like a bunch of libraries to do everything I do: https://github.com/perkeep/perkeep
- zstd (ZStandard, developed by Facebook) might be the best compression as ArchLinux switched to it recently: zstd and xz trade blows in their compression ratio. Recompressing all packages to zstd with our options yields a total ~0.8% increase in package size on all of our packages combined, but the decompression time for all packages saw a ~1300% speedup.
