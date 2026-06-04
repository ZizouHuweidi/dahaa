# Dahaa Frontend

React Router SPA for the Dahaa web client.

## Development

```sh
npm install
npm run dev
```

The development server runs at `http://localhost:5173`.

## Production Build

```sh
npm run build
```

This app is configured with `ssr: false`, so the production output is static client assets under `build/client`.

## Container

The frontend container builds the SPA and serves `build/client` with nginx on port `3000`. Unknown paths fall back to `index.html`, so client routes such as `/create` and `/g/:code` work on refresh.

```sh
podman build -f Containerfile -t dahaa-web:dev .
podman run --rm -p 3000:3000 dahaa-web:dev
```
