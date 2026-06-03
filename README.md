# dahaa

Invite-link trivia party game.

## Development

This project uses Podman, `just`, and Goose migrations.

```sh
just tools
just up-deps
just migrate-up
just dev
```

`just migrate-up` also loads starter questions for geography, literature, history, and science.

Useful commands:

```sh
just test
just build
just up
just down
just logs api
just image-api
just frontend-dev
just frontend-build
just migrate-status
```

The web MVP is guest-first: create a room, share `/g/{code}`, and join without an account.
