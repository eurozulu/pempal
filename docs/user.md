# User

user is a flag used to identify the person carrying out a task.  
The User name relates to a DN of a pre-defined certificate and serves to identify the key
to be used to carry out a task with that key.

The User name is either the full or partial distingushed name
of a user certificate, such that it identifies one, and only one certiifcate.
The name may be a full qualified DN e.g. 'O=acem ltd, OU=Accounts, CN=Chief accountant' or
a partial name, e.g. 'CN=Chief accountant'
A short form of the name can also be used where just the common name is provided,
in which case the preceeding 'CN=' may be ommited.
i.e. 'Chief Accountant' = 'CN=Chief Accountant'

A User flag can be used as an alternative to the -key flag, where the 
fingerprint of a public key is given.  The User flag enables a more
user friendly name for that key.

User certificates can be any signed certificate, however the application
usually creates a specific certificate for that name, using the required key.

User default
The default user, when neither -key or -user is given, is assumed to be the
current users OS login name.  When neither are specified, the
OS user name is used as a Common name to locate a certificate of that name.

