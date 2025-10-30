### flatenv

Having worked on other container technologies recently, I miss the flexibility of mounting formatted
files via configmaps. In a container technology where only the basic elements of a container can be
controlled (empty mounts, environment variables, and initcontainers), and more complicated configuration
cannot be mounted to filesystem points easily, flatenv can be used in an initcontainer form to render
files stored as data URIs to files on the filesystem.

v0.1 goals:

- [ ] parse all environment variables with a specified prefix (default "FLATENV_") and create files with the contents
from the variable contents in data URI format (RFC 2397). The environment variable suffix is used for the file name, where
a single underscore is an underscore, a double underscore is dot in the filename, a triple underscore is a pathsep, and a 
quadruple underscore is a double underscore.

v0.2 goals:
- [ ] allow the use of media type parameters to set the filename, user, group, and octal mode of the file. This overrides the
use of the environment variable suffix as a filename.

## Example run for v0.01

```
$ BOOP___id_hello___there__txt="data:,heya" ./flatenv -prefix BOOP_

2025/10/30 04:31:12 INFO Scanned octal dperms=-rwxrwx---

2025/10/30 04:31:12 INFO found files files=map[__id_hello___there__txt:data:,heya]

$ ls -lRa .id_hello/

.id_hello/:

total 4

drwxr-x--- 1 user user 18 Oct 30 04:31 .

drwxr-xr-x 1 user user 96 Oct 30 04:31 ..

-rw-r----- 1 user user  4 Oct 30 04:31 there.txt
```
