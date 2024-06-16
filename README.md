[![progress-banner](https://backend.codecrafters.io/progress/shell/e81cae88-a016-4403-a5de-f38016d80f2a)](https://app.codecrafters.io/users/codecrafters-bot?r=2qF)

This is a Go solution to the ["Build Your Own Shell" Challenge](https://app.codecrafters.io/courses/shell/overview).

## Learnings

- Instead of using globals, consider passing functions / structs around.

  - You can go far with this approach, especially when it comes to testing.

- The `type` command was tricky.

  - How do I pass the information about available builtins to it?

    - I've opted for a function `builtinFinder` but I'm unsure if this is a good approach.

- How do you make `exec.Command` testable?

  - There [are articles like these](https://abhinavg.net/2022/05/15/hijack-testmain/), but I'm not sure if this is the path I want to take.

  - Ideally, I would be using an _interface_ for executing commands. Then I could use a "test" implementation in tests.

    - I'm more aligned with [this approach](https://blog.sergeyev.info/golang-shell-commands/) or variations of it.
