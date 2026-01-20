# new command
new generates new PKI resources.  
Each supported [resource type](./resources.md) has a factory
which defines how each type is generated.

The factory interprits a template into the 
properties it needs, including existing keys and certificates,
to perform the generation.
The resulting generation is stored in the relevant directory(ies)
based on the Infra config.

new has one or parameters, naming a known template.  
When more than one template is specified, they are merged into one,
with the last given name taking precedence over preceeding values.  

The final template must include the 'type' proty, defining the base resource type.
e.g. type=CERTIFICATE

`pp new sampleservercertificate -subject "OU=development,CN=testserver1.dev.acme.com" -user "anondev@acme.com" -issuedby "Dev server intermediate"`
Using a template called 'sampleservercertificate', create a new certificate,
signed by the private key of 'anondev@acme.com'.  
If the same user has access to the private key of the issuers certificate, the new certificate
will be counter signed automatically.  Otherwise a new CSR is generated
for the certificate.  

Good output: (Both private keys available)
created new certificate: testserver1.dev.acme.com, expires 01-01-2525T12:00:00
issued certificate testserver1.dev.acme.com by "CN=Dev server intermediate"

OK output (only certificate key available:
created new certificate: testserver1.dev.acme.com, expires 01-01-2525T12:00:00
created sign request for testserver1.dev.acme.com for "CN=Dev server intermediate,EM=devcsrs@acme.com"

bad outcome, no keys found
Failed to create certificate: testserver1.dev.acme.com
user 'anondev@acme.com' key not found.



Keys

`pp new rsakey -keylength 4096`

Output:
Created rsa key ./private/F53dc...D3

`pp new key`
defaults to:
`pp new key`


