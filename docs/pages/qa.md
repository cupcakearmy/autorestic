# ❓ QA

## My config file was moved?

This happens when autorestic needs to write to the config file: e.g. when we are generating a key for you.
Unfortunately during this process formatting and comments are lost because the `yaml` library used is not comment and/or format aware.

That is why autorestic will place a copy of your old config next to the one we are writing to.
