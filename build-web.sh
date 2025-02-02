#!/bin/bash
cd tldraw-web; pnpm run build;
cd ../;
cd vue-web; pnpm run build;