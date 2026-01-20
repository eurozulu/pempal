# Resources
The PKI uses a number of resources to create and validate
the infrastructure.  Pempal defines these resources as:
* public/private key pairs
* x509 certificates
* x509 signature requests
* x509 revokation lists
* Digital signatures
* Encryption archives

#### Public / Private Key pairs
Key pairs are the heart of a PKI.  Each induvidual key
has a public and private componant.
The private key is held in a secure location and used to authenticate
the ownership or control of the assciated public key.  
The public key is widley distributed and serves as a unique identifyer
of an entity in possition of the private key.  
The unique identity is key to linking resources.

Pempal index both public and private keys seperately.  
public keys may include known public keys for which the
PKI has no access to the asscociated private key.  
Private keys consist of both public and private.



#### Certificates
Certificates act as the public face of public key.  
In effect the certificate ties a unique text name, to
a known public key.  A label of a key!
Each Certificate is signed by the key holders private key,
and, when in a trust chain, also signed by an 'issuer', third party.

#### Requests
Certificate Signing Request, or CSR is a
request for an issuer to sign a certificate
as validated by them.  It assigns a level of trust to the
certificate being signed, derived from the trust of the issuing party.

#### Revokation
Revokation lists are lists of Certificates which have been prviously issued,
which are no longer considured valid by the issuer of the list.  
Expried certificate (by date) are invalid by default and donot appear
on Revoke list.  These lists are for the early termination of an issued
certificate.  Although the certiifcate is still technically valid,
if will be invalidaed if it appears in a trusted revokation list.

#### Signatures
Digital signatures are a means of ensuring the integrity of
a given digital stream.  The signature is generated
by a private key owner such that any consumer of the
data stream can compare the stream to the signature hash
to confirm the data is as it was presented when originally signed.
In addition, much like a certificate Issuer, it assigns a level
of inferred trust on the signed object, in that
the owner of the signing key is "putting their name against" something,
inferring any trust in that signer, onto the data stream itself.

#### Archives (Encryption / Decryption)
Key pairs offer a useful mechanisum for securing data through
strong encryption.  Using key pairs, data can be exchanged
and stored in an encrypted format, only accessable to
to intended private key holder.  Anyone trusting a public key
(either by certificate or otherwise) can send pre-encrypted data
only the intended receipiant can decrypt. 

Archives create a single file of encrypted bytes,
containing one or more files.  The files can only be decrypted
with the private key.




