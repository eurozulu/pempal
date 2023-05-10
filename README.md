# Pempal
### A tool for managing tls/x509 certificates and the associated resources.  
  
Pempal offers an easy means to manage creating, finding and monitoring your x509 resources
using a combination of Resource scanning, Key Management and an advance templating engine.  
  

#### Find
Find resources using keyword searching and property searches allows for
finding specific resources. 
e.g. 
- Find certificate issued by a given issuer  
  `find ./certs -type certificate -issuer "*Acme Root CA*"`
- Find certificates about to expire in x days  
  `find ./certs -type certificate -not-after "{{ nowMinusDays 30 }}"`
- Find certificates signed by a given key  
  `find ./certs -type certificate -key "2345678"`  
  
#### Make
Create new certificates and asscociated resources using predefined, named templates
to form an expressive means of what to create.  
e.g.
- Create a cretificate for a specific server type  
  `make dev-test-server-certificate issued-by-dev-manager -createkey`
- Define organisation wide templates to apply to all certificates,  
  and combine with more specific types on the fly  
  `make acme-org-certificate dev-intermediate-ca issued-by-acme-root-ca -key dev-manager`  
- Templates can be defined to enforce standards for specific tasks    
  A template 'ca-root-key' might define the key algorithm and length and enforce password
  `make ca-root-key -out private/caroot`  
  `make request dev-client-access -out requests/dev-client-access`  
- `make server-default-key -out private/{{ .subject.common-name }}`


Make can be used to create any one of four resource types:  
- Certificates
- Certificate Requests
- Certificate Revokation Lists
- Private Key  

Each of these types has an asscociated "resource template" which contains the minimal required properties to make that resource.
These templates are named respectivly:  
certificate  
request  
revokation  
key  

Make requires one or more template names as its parameters.  Templates are formatted and merged into a single template  
representing the new resource to create.  Begining with the first template, make will merge the second into it, over writting
any values which share the same name.  Once merges, the next template is then merged into that and so on, each template to the right  
taking precedence over existing property names.  
  
Once the template is established, Using the keys identified in the template, make will reconcile the required keys and certificates to perform the operation,
including the private key to sign requests and/or the key and certificate of the issuer.
The resource will then be generated as a PEM output.  


## Templates
At its heart is a template engine which is used to both display and create x509 resources.  
A template is simply a collection of key/value pairs or named properties, which contain
the values of the resource values asscociated with that name.  
e.g. One of the simplist templates is the publikey template.  This has just three  
properties:
key-algorithm:  
key-length:  
public-key:  

The 'key-algorithm' property specifies the key algorithms, RSA, ECDSA, Ed25519 
- key-length specifies the length of the key
- public-key displays the key identity, (a md5 hash of the keys der encoding.)

Of course a public key template is used for display only as you can not generate a public key on its own.  
The private key template contains similar properties but with additions such as  
is-encypted  

A private key template also displays a public key identity  
  
When generating a new private key, the template must be populated with all the required values.  
The 'key-algorithm' must be a valid algorithm name, and the length must be appropriate for the key algorithm.
e.g.
key-algorithm: RSA  
key-length:2048  
  
Using this template, the `make` command will generate a new RSA key.  
  
#### Template Tags
Tags may be placed at the beginning of a template to control how the template is formatted.  
There are two tag types:  
- #extends
- #imports
Both are followed by a space and then a template name.  e.g. `#extends certificate`  
The imports tag may also have an additional 'alias' names, following the template name.  e.g. `#imports prod-server-issuer issuer`  
Both templates can be used zero, one or more times however with the extends tag, it is important to note the order in which the tags
are placed as this will define the prder in which the tem,paltes are merged.  
Imports tags order is unimportant as thy are referred to by name.  

##### Extending templates
When a template contains an extends tag, the named template serves as a base to merge the extending template into.  
The named extends template is loaded and the formatted extending template is merged into it.  
Any values in the extended template, also found in the extending template, will be overwritten by the extending template.  

##### Importing templates
Rather than merging an entire template into the other, specific properties of another template can be imported using the imports tag.  
Imported templates are used in conjunction with template macros to inject values into the template.
Each imported template is assigned a name, by default the template name itself, or, should that clash with another name, an alias name.  
The alias is specified after the template name, following a seperating space.  
Template macros will refer to these values using dot notation, beginning with the imported name.  
e.g. `#imports my-imported-issuer myissuer`  
`subject:  `  
&nbsp;&nbsp;&nbsp;&nbsp;`common-name: {{ .myissuer.common-name }}`  

#### Template Macros
Macros are a means of replacing a section of the template with a calculated value.  
Macros may contain references to values imported with #imports tags or calculated from one
of the many functions available inside macros.  Refer to the macro section later for more details.  
Every Macro is enclosed in double curly bracked `{{` &nbsp; `}}`  
When referring to imported values, the macro uses the dot notation, beginning with the imported name:  
`{{` &nbsp; `.myimportname.propertyname` &nbsp; `}}`  
Properties containing child properties, such as 'subject' and 'issuer' refer to those child properties by
extending the dot notation  `{{` &nbsp; `.myimportname.subject.organisation` &nbsp; `}}`  
Note properties like 'organisation' are arrays, which are converted into comma delimted values.  

#### Template Formatting
A template can be in two states, Raw or Formatted.  
Raw is the form the template is stored, including its tags and macros.  
A Formatted template has no tags or macros and is the result of process
those tags and macros.  
i.e. A raw template is how it is written and stored and a formatted template is how it is used.  

#### Populating Templates
Values in templates may be populated in a number of ways:  
- Merging of other existing named templates  
This will merge the values of one template into the other to form a single template with the values of both.  
- Extending another template  
Using #extends tag to specify an existing template to extend.  This is similar to merge except the merge is carried out
whenever the extending template is used
- Specify the property as a flag to make  
Values can be passed as flags e.g. make -common-name "My new server certificate"  
- Template "macros" {{ ... }} token embedded in the template which perform lookup and functions to generate values.

  
#### Managing Templates
Templates can be managed using two commands:  
- `template`
- `templates`  

In singluar form it will display a formatted template.  In plural form it
lists the names of the known templates.  

##### template command
template will display the given named template in its formatted format.  
Specifying more than one template name will merge all the named templates
into one, in the same manner the `make` command does.  This allows the user
to view how a series of templates will be formated given their names.  
  
The output of template is to display the template in its natural 'yaml' format.  
To view the raw template, prior to formatting use the `-format` flag.  
`-format` take one of three formats:  
- yaml (default)
- raw pre-formatted form
- json A json representation of the formatted template.

##### templates command
templates is used to manage the templates available to the application.  
On its own, with no flags, it will display ALL the template names available, in alphabetical order.  
This list can be filtered using a single text parameter containing part of
name of the template you wish to find.  
e.g. `templates server`  will find any template with the word 'server' in it.  
