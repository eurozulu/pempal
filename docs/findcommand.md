# find command

`pp find [search path] [-type cert|csr|crl|key] [-query "<field value> [<|>|=|!=]<field value>"`  

find locates the PEM resources matching the given criteria.  
Criteria may be defined as a speicifc type or types and/or one or more property values of the resources being searched.  

-type flag may specify one or more (comma delimited) reslurce types to limit the search.  
When given, only resources of the given type(s) are found.  When not given all resources matching the query
are returned.

-query specifies one or more properties of the resources being sought.  
properties are dependant on the type of resource.  if the type is specified and a query contains a property NOT known
to that resource type, none of those resources will be found.

i.e. you can't specify `-issued-by` when searching keys!

`pp find -type cert -query "isca=true"`  
Finds all the certificate authrity certificates.

`pp find -type cert -query "notafter <= (now() - month(1))"`
Finds all the certificate which will expire in the next month.  

`pp find -type cert -query "issuedby=O=acme.com,CN=myissuer.com"`  
finds all the certificates issued by 'myissuer.com'

`pp find -type csr -subject "O=acme.com,OU=developers"`
Finds all signing requests for the OU=developers.

Combing find results:

`pp find -type csr -subject "O=acme.com,OU=developers" NOTIN -type cert -subject "O=acme.com,OU=developers"`
Same as previous query, all CSRs for the developers OU, but further filtered
for NOTIN, a list of certificates with the same name.  
i.e. the outstanding requests, which have no matching certificates.  

`find -type key -publickey IN -type cert -subject "CN=myemail.acme.com"`
finds the private keys for the certificate(s) with the common name of "myemail.acme.com".


### Query aliases.
Queries can be predefined and stored in a Infra config.  
The app provides some predefined queries it uses internally,
which are available to use.

alias: keys
lists all the owned private keys.

alias: certs
lists all the known certificates

alias: requests
lists all the known certificate sign requests

...


alias: keyfor -subject <DN>
returns the private key(s) for any resoruce with the given subject.  
As above ecample:
find -query "subject"
`find -type key -query "publickey IN -type cert,request -subject #subject`

alias: issuer -subject <DN> ??????
returns the issuer certificate(s) .  
`find -type cert `

