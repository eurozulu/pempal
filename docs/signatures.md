# Signatures
Signatures are generated based on a given byte stream and the current 'user'.
The stream is read and a hash generated which is then signed by the private key of the users.  

`pp sign ./myfile.txt > ./mysignedfile.txt -user "myemail@acme.com"`

The myfile.txt file is copied out to the new mysignedfile file, with the addition 
of the signers hash and key.  
The signer is determined by the `-user` flag, which in the example is "myemail@acme.com".  
This is read as a common name, matched to an existing certificate with that name.  
Using the certificate, and its public key, the relevant private key is located.  
Finally the private key is used to generate the signature hash, which is appended to the final result.  

