# Rename backend to repository

# Env variables

AUTORESTIC_BB_B2_ACCOUNT_ID=123 -> AUTORESTIC_BACKENDS_BB_ENV_B2**ACCOUNT**ID=123

- All fields can be configured by env now
- To escape `_` replace it with double underscore `__`

# Rest property on backend config

No rest property anymore, can be used in string extrapolation

# Every string is now replaceable with env variables
