# filelock

This is a practice module to teach myself golang. It's a reader-writer lock for locking files.

The package allows for aqcuiring multiple permissions before executing a function, then releases them after running the function.

## Example

```go
package main

import (
  "filelock"
)

func main() {
    ctx := filelock.NewContext()

    permissions := map[string]bool{
      "file-1": true,  // true means write permissions
      "file-2": false, // false means read permissions
    }

    release = ctx.WithPermissions(permissions, doSomething)
    defer release()
    // access the files, do something that needs to happen atomically
}
```

## Use Case

The primary use case is when you need to read from one file, use the information to write something to another file, and these actions need to happen atomically. You can acquire read permissions and write permissions before doing anything at all.
