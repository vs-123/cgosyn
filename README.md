# cgosyn
C with Go syntax

## Why?
Cuz why not

## But like, isn't Go syntax based on C?
Yea it is, but I still wanted to make it

## Aim
To transpile Go-syntaxed C to C (1 to 1 without any abstractions)

## Features and Limitations
- Transpiles Go-syntaxed C to C (1 to 1 without any abstractions)

- The C code is transpiled with indentation!

- Go's `bool` types are converted to `int` in C, and does not support `true` or `false`, instead you need to use `1` for `true` and `0` for `false`

- Imports in Go will be `#include`-ed as it is (for e.g., `import "stdio.h"` will be transpiled to `#include "stdio.h")

- Unfortunately, `import "<stdio.h>"` will not be transpiled into `#include <stdio.h>` but instead `#include "<stdio.h>"` (will be added in the future)

- The `package` keyword has no use as of now (feel free to contribute if you have any ideas for it)

- Does not support same type simplification (`(x, y int)` will NOT be transpiled to `(int x, int y)`, this will be added in the future)

- Doesn't have error checking lol (not planning to add either)

- However, warnings will be thrown for certain things which it could not transpile

## Examples
You can find them in `examples/` directory

## Contribution
Contributions are highly welcomed

Feel free to make pull requests or raise issues related to this project