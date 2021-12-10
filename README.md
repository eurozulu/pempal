# Pempal
### A tool for managing tls/x509 certificates and the associated resources.  

## This project is currently still in development stage!  

### Overview
Pempal is designed to help anyone using TLS/SSL Public key Infrastructure certificates,  
both in managing existing resources and generating new ones.  
From developers creating test certificates to those managing large scale PKI deployments,
pempal can assist by simplifying common operations and managing the existing resources.  

Pempal operates in two distinct modes, as a command line tool or as a REST API server.  
Both forms carry out the same functionality. 
The command line mode is explained here, with the REST API documented elsewhere. (TODO)
  
### Concept  
Pempal is designed to assist in working with certificates, to help perform common tasks associated with managing them.  Both in finding and checking your existing resources and when generating new resources.  
  
Pempal doesn't do anything other tools such as openssl do, in fact it does far less, but it aims to do those tasks as easily as possible.  
Avoiding the myrad of flag options and sub commands it aims to keep everything simple.  

Managing existing resource, locating, viewing and querying helps find those resources which might be about to expire or are revoked etc.  
Maintaining large infrustructures, it tracks your private keys to help manage which keys sign which resources, as well as enforcing certificate standards and properties thorugh its template system.  

Generating new resources, pempal uses a template system to define the specific properties of your new resources.  This system allows you to pre-define specific requirements for any resource and apply those when creating them.  
Templates can be easily created using a yaml editor or copied from existing resources or templates.  
This makes copying existing resources simple, by pointing to them as a 'template'
all the properties are copied into a new resource (with a few exceptions like valid from) 
and, assuming you have access to the same private key which issued that certificate, the new one will be signed, automatically, by the same issuer key.  

#### Resources
pempal refers to 'resources' as being any certificate or related resources.  Those being:
- Private / public keys
- Certificates
- Certificate Signing Requests
- Certificate Revokation Lists

#### Locations
A location is the 'address' of a container of one or more resources.
Loctions can be:  
- directories
- resource filepath (der, crt, key, prk etc)
- container filepath (pk12, pk7, pem)
- url  
A location may have an optional extension of a colon, followed by a number.  
  This is known as the location index and applied to the exact location of a single resource,  
  found in a container of multiple resources.  e.g. a pem file with more than one certificate in it.  
  

#### Template
A Template is a collection of pre-defined properties which may be applied to a given resource.  
Used in generating resources, templates have their values calculated and applied to each resource.

  
### Usage
Pempal is controlled using a 'command' as the first argument of the command line.  
(Or Rest endpoint).  
Optional flags (parameters starting with a dash '-') can be specified to control the operation.  
There are some general flags to control the overall output, which apply to every command,  
and each command has zero or more flags, specific to that command.  

The command line format is: [optional general flags] < command > [optional command '-'flags] < location > [...locations]  
Although the distinction of general and command flags is made in the position,  pempal isn't too fussy on if they are before or after the command.  
What is important is that general flags always appear before command flags.  
so another valid format is: < command > [optional general flags] [optional command '-'flags] < location > [...locations]

Each command requires at least one 'location' to specify the resources to act on.  
multiple locations may be specified, each as a distinct argument and pempal will perform the operation over each of the locations specified.  
e.g. to `list` the contents of two directories:  
`pp list ~/.ssh /etc/openssl`  
This will show all the resources found in both '~/.ssh' and '/etc/openssl'  
  
Commands can divided into two groups:
- Commands for managing existing resources
- Commands for generating new resources


#### Managing existing resources
There are three commands for this:
- `show`  (`peek`, `sh`)  
  displays the properties of one or more resources  
  primarily for viewing resources and extracting properties
  
- `list`  (`ls`)  
locates and list resources with a selection of their properties.  
Locates resources with specific properties and/or type  
Primarilly for finding resource which match a given criteria  

- `format` (`fm`)  
  (Re)formats existing resources into a specified format.  
  Format is used to 'package' resources into a common format.  
  It currently supports pem, der and will support pk12 soon.
  When combined with `list` it can be used to create bundled resources with a matching criteria.  
  e.g. collect a certificate chain and package into a pem file.  
  
#### list
List is used to find resources.  
Used without any flags it will display a brief outline of all the resources in the location list.  
This will show the PEM type (CERTIFICATE, PRIVATE_KEY etc), if a key is encrypted and, most importantly, the location of the resource,  Its filepath (and index).  

List can also be used to filter resources using flags. 
A basic filter is the pem type, which limits the output to pem types given:  
e.g `list` -type="CERTIFICATE, PRIVATE_KEY"  
Another basic filter is the header, which searches for a given value in the headers:  
`list` -header="*encrypted*"  

Pipping lists  
Using with the quiet flag '-q', `list` limits output to just a location list.  
Such a list can be pipped into another pempal command, such as `show` or `format`.  
e.g. `pp list -q -publicKey 47f88sd990... ~/certs/ | pp format -pem -`  
This finds all resources in ~/certs which are signed by the public key and then pipes those into
`format` which writes them out as pem blocks.

Query  
In addition to basic filters, `list` can query reousrces for specific attributes.  
Using the `-query` flag, one or more key/value pairs can be specified to limit the
output to resources matching these attributes, where the key is the (case insensitive) property name
and the value is a regex expression to apply to the resources value.  
e.g., `list -query "isca=true, issuer=".*Acme root.*`  

This will find only resources with the 'IsCA' and 'issuer' properties present (i.e. certificates)
and those certificates will be CA true whith an issuer containing the Acme string.  
Query isn't limited to certificates but any resource. The previous example would
only find certificares as the 'IsCA' property is unique to signed certificates.  
queries with a more common propery, such as 'PublicKey' will find every resource,
as public key is a property used by all resource types.  The -query flag can be used in conjunction with
the -type flag to limit the resources to a specific type if needed.  

#### show
Show displays the properties of a resource in a 'yaml' format.  
It takes one or more arguments, wach being a location to locate resources.  
Each resource found is parsed and presented in a yaml format, with each value next to the property name.  
Values difficult to display as text are encoded into base64 and some well known, numeric properties  
(such as KeyUsage, PublicKeyAlgorithm etc) are encoded into a string readable word.  
  
Show uses flags to limit the properties it will output.  by default it shows all properties in the resource.  
using the -select flag, a comma delimited list of case insensitive property names can be listed and
only those values will be displayed.

When combined with `list` show can be used to build property lists of resources based on criteria specifed in list.  



### Generate New resource

#### Key tracking
When generating any new resource (with the exception of a key pair itself), a private key is required to perform the signature.  
With each certificate in a chain, there is a private key assiciated with it.  
Managing all but the most trivial certificate deployments can mean these keys and the associated certificates can become unmanagable.  
Pempal aims to assist in this by tracking the private keys you have access too and associcating those keys with the resources they have signed.  
Using this asscociation, primarily with certificates, the key tracker can establish which keys to use for which action.  
Properties of the associcated certificates, such a key and extkey usage define the operations which can be carried out by each key.
Properties in theose certificates can also provide guidance on selecting the keys, such as the issuer name match.
When a new certificare has a specified issuer, that DN is used to locate to associcated CA certificate and the private key which signed it, making the issue a simple task of confirmation.  
OF course there will be times when more than one key is suitable for a task, and pempal will prompt the user to select a key (or generate a new one).  
Alternatively, the -key flag can be used to specify the specific key to use.  
The value of -key can be a location which contains only one private key OR
the sha1 hash of the keys public key. (known as a keyhash)
When specifying -key, the keytracker will always use that key, which may 
override the new resources 'issuer' value with the keys own issuer certificate.

There are three primary commands for generating resources:  
- `request`  
  Generates a new Certificate signing request (CSR)
- `issue`  
  Generates a new, signed certificate
- `revoke`  
  Generates a new signed certificate revokation list  
  
In addition a forth command `key` will generate a new key pair, however in the normal flow of operations `key` is not required as it is called for you when neccessary, by the other three commands.  

Generating commands use the same command line format where one or more locations must
be specified, however it treats those locations, not as pem resources, but as templates.  
In effect these means and resource can be listed in the arguments and its properties,
if relevant to the new resource, are copied into it.  
A simple means to copy of an existing certificate would be:  
`pp issue ~/cert/mycert.crt`  
In order to carry this out pempal requires two keys:  
The original signer of the mycert certificate and the issuer of that certificate.  
Assuming keytracker can locate both these keys, a newly formed certificate
will be generate containing the same properties as the original, but with a new
valid_from date and signature.
A simple way to renew a certificate thats about to expire.  

In addition to existing resources a template can be used.  Tempaltes are named with a preceeding hash '#'.
When a known template is specified, the properties from that template are applied to the new resource.  

Like locations, multiple tempates may be specified on a generating command,
and they are merged down into a single template, with the right most one taking precedence.
See tempaltes later for more details.
This multi argument function allows certifictes to be created using a combination of tempaltes.  
e.g. `pp issue #BaseCorporate #Devlopment #restserver #privilidgedOU`  
  
This would collect all four templates applying the 2nd (#development) template
on top of the 1st (@BaseCorporate) template, then the 3rd (restserver) to that and so on
All properties from all templates are used.  When propertiey names clash,
with more than one template having the same property, then the right most template takes precedence,
overwritting the earlier (left) value.

### request
Request is used to generate a new Certificate signing reequest, or CSR.  
As with the other generators, it uses an Identity to establish the key with which to sign the new CSR.  
See Key Tracker for more details.  
Again, as with other generators, Request takes one or more locations to form a compound template of the properties to insert in the CSR.  
request then applies this compound template to a base (hard coded) request template to ensure all the required properties are present.  
- PublicKey
- Signature
- Subject
If any of these, or any properties from the templates required properties are missing, the request will fail, throwing an error listing the missing properties.  
Once all properties are present, the Identidy is matched.  
For Requests, this requires DigitalSigning keyusage.  A private key must be present, with at least one associated certificate containing that usage, in order
  to use that key.  
One keys are established, the user is offered to choose from the available or generate a new one, a new one being the default choice.  
  
Having the key selected/generated, the CSR is finally generated and encoded into pem.
In the case of a new key, the pem will contain both the CSR and the newly formed private key  

Output
By default the output is to the stdout.  This can be redirected using shall redirection, to a file.
e.g. `pp request #server #developer > devserver.pem`  

The file extension defines the format the resources are written in.  Valid options are der, pem and (later) the container pk7 pk12.
If not specified, pem is the default. 
  

### issue
Issue generates new certificates.  
As with the other generators, it uses an Identity to establish the key with which to sign the new Certificate.  
See Key Tracker for more details.  
Again, as with other generators, `issue` takes one or more locations to form a compound template of the properties to insert in the Certificate.

The compound template is applied to the certificate base template to ensure all required properties are present.
- PublicKey
- Signature
- Subject
- Issuer  
Missing any of these results in a failure to issue the certificate.  
With all properties present, `issue` will locate the identity bound to the issuer, See KeyTracker.  
Using the issuer ID, the new certificate is generated and signed by the issuer.

Optional flags:  
--key:  Specify a specific key with which to sign the certificate.  Overrides Issuer in the template and uses the key's certificate value.  
The value os  "+" can be used to indicate a new key should be generated for the certificate.
This key will be used to sign the original Name and subsiquently qriten out with the certificate.  
Otherwise the value should be the sha1 hash of the public key matching the private key to use.  



### Templates
Templates form the heart of how resources are generated in pempal.  
Templates greatly simply the generation of new certificates and, when used with customised templates, can enforce standards and ensuring consistency across the deployment.  
At the most simple, a template is a collection of one or more properties with the value to apply to the new resource.  
Templates can be combined to form complete property sets to apply.  
In addition, templates can be used to form new values or ensure values are presnt.
Using special key values, the template can be instructed to gather the value from a source other than itself, such as user interaction or calculated values.  

These 'special' key values can be placed in templates to control how the final value is set.  
Required '?' .  When a template has a '?' as its value, it signals it is a required value, not supplied by the template.  
It is expected that the value be provided either by another template, a flag value or from a user interaction.  
When that value is not present, an error is thrown, preventing the resource from being generated.  

"${ < func > }"
Templates can have a function enclosed on ${ } braces which will calculate the value
to be inserted based on other properties in the template.  This can be in the form of a simple property name, preceeded with a dot.
e.g. The built in templates: #SelfSign and #RootCa, both have:  
`issuer: ${.subject}`  
This indicates the value of '.subject' is copied into the issuer value.  i.e. both resources are self signed.  
More complex functions can be used to combine and manipulate values into their final form.
See the golang text template functions for more details.

#### Compound templates
Compound templates are templates formed through combining multiple other templates into one.
The merging of the templates is performed in order, from left to right, so the first(left most)
template sets the base properties, and the next, 2nd applies more properties to them. 
The 3rd then follows and so on until all templates are merged.  
When two templates contain the same property, the last one to merge (the right most)
will take precedence, overwritting any previous value.  
In the case of required values '?', the '?' symbol is only copied when there is no preceeding value there.  
i.e. it will not overwrite an existing value.  Any following template will overwrite the ?, as expected.  

Functions are merged initially into the final template, however any remaining, i.e, not overwriten, are then calculated
based on the merged template state.  As such a function can have access to a value from an upstream template, by simply
referening that property.  The property value is evaluated after all templates have been merged, making the value accessible to the function.

#### Enforcment
Compound Templates can be used to enforce properties onto new resources such as CSRs and Certificates.  
For example, if every new cert must conform to a specific keyusage, a template containing that usage can be applied to every CSR.  
Or every cert must have a specific Organisation "Acme Ltd".  
Forming a compound template with the original CSR and the templates to enfore results in a certificate conforming to the deployment standards or an error when they do not.  
`pp issue ~/certs/newrequests/blabla.csr #fixkeyusage #fixorganisation`  

Assuming #fixkeyusage has the set KeyUsage and ExtKeyUsage properties, and #fixorganisation is simply "Organisation: "Acme ltd" "  
This will generate a Compound template for the new certificate with all the CSR properties and the keyusage pre-set.
The organisation is enfored on to the new template, as a result, if the request value being overwriten was not the same (or not present) as the enfored value,
the signature of the CSR will now be invalid.  (The Subject Name has changed, invalidating the previous signature)
The result is a failure to generate the certificate. (As it didn't 'conform')  
To avoid this, the same templates can be shared with requesters, allowing them to apply them  
to their requests prior to signing and so ensuring a consistant set of standards are applied to all resources in the deployment.  

