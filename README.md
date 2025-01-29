## Project Structure

```bash
|-- docs/             # Project documents
|------ api/          # Bruno API collection
|
|-- api/              # BFF API written in Golang Echo
|------ data/files    # Data content stored in sub-repo grd0-site-data
|
|-- assets/           # Static assets uploaded to Cloudflare R2
|------ blog/         # Blog static content, stored by YYYY/MM
|------ gallery/      # Gallery static content, stored by collection
|------ profile/      # Author profile picture, stored by email
|
|-- assets-helper/    # Helper script to generate Table of Content for assets
|
|-- vue-web/          # Vue 3 + Vite project SPA webpages
```
