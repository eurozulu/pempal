# Archives
Archives are collections of (usually related) resources, optionall encrypted with
a public key.  
Archives are simple zip type files, but at the resource level, rather than a file level.  
The relevant PKI directory sturture is duplicated and the relevant resources
encoded into the Archive.  Each resource is indexed with file meta data,
spcifying filenames,, modes, paths etc.

As well as PKI resoruces, the archive can contain arbitary data as byte blocks.  

Various output encodings are available.
to support
PK#7,  pk#12 packages of chains and keys.

Default is a chain of ordered PEM blocks.

Order is determined by resrouce type and trust chain.

