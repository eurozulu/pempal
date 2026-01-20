# Commands
Pempal offers a series of tools which can be grouped into two main types:  
1. Search.  Locating known resources in the PKI based on quering meta data.  
2. Generate.  Generating new resources using templates and existing resources.  

Search commands list resources based on a query criteria.  
The command line `find`  
Flags:  Find accepts and valid meta-data name as a flag.
The value of the flag can be a sinple, literal value or,
if preceeded with a comparison symbol, (<, =, >, >=, <=, !=),
an expression to compare to that field.


- find.  Locate resources, based on a given query.
- new.   Generare a new resource based on the given template.
            keys, certificates, requests, revokelists.
- verify.  checks the given resource has a valid signature.
- view.    outputs the meta data about the given resource(s).
- archive.   generates a single file archive with a collection of resrouces.  
- encrypt/decrypt.   Encrypts/decrypts a given stream using a specifed key.
- sign.    Generates a signature hash of a given stream.  


Management commands:
- config.  Manage the default values and Infras in use.  
- template.  outputs templates in yaml format, by default, merging all named templates.  
Used with a `-resource` flag, can convert existing resources to tempaltes. (serialise resrouces) 



New.  
Generates a new resource.  
Each resource type has an asscocited factory, which requires a specific template type.  
The Generate command collates and verifies the build template,
and selects the relvant factory for that template.  
The collating of the template first merges all given templates into one.  
The template type is then identifed with key metanames.  
Using the type, the factory is identifed and passed the new template,
along with a reference to the infra.

Usiong the infra, the command identifies the location
of the new resource, based on its type and or/meta data.  
In addition, the command uses the Infa index to locate its dependants.

Finally, the Factory will initiate its build and store the result
into the relevant PKI infra lcoation. 



======

`find`
Finds the known resources based on a given query.  
Find allows a user to locate resourced based on 
the meta data of the resrouce.  ITs a search tool for finding
resources and related resources.  

`pp make cert acme-org development-server -dn "cn=mytestserver.org" -public-key "myemail@acme.org`  
Generates a new 'dev-server' in the acme organisation, combining the base Acme org properties
with those in the 'development-server' tempplate.
The dev server template can contain a pre-defined 'issued-by' field, with the DN of the relevant issuer.  
On execution of the template, the public key of 'myemail@acme.org' is used to sign the new certificate,
and the predefined issuer is located.  
The issuers certificate should be known.  
If the issuers private key is available, the certificate is
automatically signed.  When the private is is unknown/unavailanle,
a new CSR is generate, for the issuer to sign at a later date.  


Make is used to generate new resources
`make [cert|csr|key|revoke|signature|archive] <template name> [...<template-name>]`

Each resource type has its own, specific set of flags
to define how the new ressource if created.

`verify`
Verify checks existing resources are valid.  
`pp verify mytestserver.org`  
Will locate any resource with name (DN) and matching the given name,
is checked for integrity.  
A basic verify simply checks if the resrouce is valid.  
i.e its signature is a valid match to the known key.  
Verfiy can be 'enhanced' to check further states of a resrouce.
Ownership can be establided with:
`pp verify mytestserver.org -asowner`  
This checks the same integrity but also confirms the asscoiated
private key is available.

`pp verify mytestserver.org -chain`  
The chain is applied to certificates and confirms the
entire certificate chain can be located and verified.  

