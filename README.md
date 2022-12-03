# Tailscale Webhook Adapter

A simple service to receive incoming [webhook notifications from Tailscale](https://tailscale.com/kb/1213/webhooks/)
and reformat for use with [Discord](https://discord.com/) in ``Infinity Development Endpoint Security``.

See ``.env.sample`` for the variables that need to be set.

----

## Tailscale Setup
Follow the instructions to [setup webhook notifications](https://tailscale.com/kb/1213/webhooks/),
and store the Secret as an environment variable named `TS_WEBHOOK_SECRET` for this service.

----