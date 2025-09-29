# Health checks

Before and after applying the blueprint, the ecosystem is checked to ensure that it is healthy.
The following is checked:
- Health of all Dogus based on the Dogu-CRs
- Check whether all Dogus already use the latest version and configuration

## Ignoring health

Upfront health checks can be deactivated:
- for Dogus, if `spec.ignoreDoguHealth` is set to `true`,

This makes it possible to fix errors on Dogus and components via Blueprint.
For a Dogu upgrade, however, a Dogu must be healthy in order to be able to execute pre-upgrade scripts.
Ignoring the dogu health can therefore lead to subsequent errors during the execution of the blueprint.
