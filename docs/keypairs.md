# Key Pairs
## Public / Private key pairs form the heart of the PKI.

The Keys are reprented in two forms internally.  
`PublicKey`
The Public key is a simple wrapper of a public key.  
It does not contain the asscocitated Private Key.  
`PrivateKey`
The Private key represents a known private key which the current user has access to.  
All private keys contain an asscociated `PublicKey`.

Properties
Keys can be any of the supported key types.
Supported key types are
- RSA
- DSA
- ECDSA
- Ed25519

Each type has a specfic set of additional parameters.
RSA has a key-length
DSA has no properties
Ed25519 has a curve property.

