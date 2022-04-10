# RWLock

This is a practice module to teach myself golang. It's a reader-writer lock for locking files.

The package allows for aqcuiring multiple permissions before executing a function, then releases them after running the function.

```go
package main

import (
  "rwlock"
)

func main() {
    var ctx rwlock.LockContext

    permissions := make(map[string]bool){
      "file-1": true,  // true means write permissions
      "file-2": false, // false means read permissions
    }
    
    doSomething := func() {
      // access the files, do something that needs to happen atomically
    }

    ctx.WithPermissions(permissions, doSomething)
}
```
