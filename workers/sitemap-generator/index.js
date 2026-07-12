const PRODUCTION_API_URL = 'https://api.grd0.net';
const pages = [
  { path: 'gallery', dynamic: listGalleryUri },
  { path: 'blog', dynamic: listBlogUri },
  { path: 'travel/map' },
  { path: 'music' },
  { path: 'coffee' }, // 20260712 Coffee Page
];

async function siteMapConstruct(base_url) {
  let urlset = [];

  for (const page of pages) {
    urlset = [...urlset, urlsetConstruct(base_url, page.path)];

    if (page.dynamic && typeof page.dynamic === 'function') {
      const subpages = await page.dynamic();

      for (const subpage of subpages) {
        urlset = [...urlset, urlsetConstruct(base_url, `${page.path}/${subpage.uri}`, 'weekly', subpage.updated_at)];
      }
    }
  }

  return `<?xml version="1.0" encoding="UTF-8"?>
  <urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
    ${urlset.join('\n')}
  </urlset>
  `;
}

function urlsetConstruct(base_url, path, changefreq = 'monthly', lastmod = null) {
  if (lastmod) {
    const lastmod_date = (new Date(lastmod)).toISOString();

    return `<url>
      <loc>${base_url}/${encodeURI(path)}</loc>
      <lastmod>${lastmod_date}</lastmod>
      <changefreq>${changefreq}</changefreq>
    </url>
    `;
  }
  
  return `<url>
    <loc>${base_url}/${encodeURI(path)}</loc>
    <changefreq>${changefreq}</changefreq>
  </url>
  `;
}

async function listBlogUri() {
  const url = `${PRODUCTION_API_URL}/blog`;
  const options = {method: 'GET'};

  try {
    const response = await fetch(url, options);
    const data = await response.json();

    return data;
  } catch (error) {
    console.error(error);
  }
}

async function listGalleryUri() {
  const url = `${PRODUCTION_API_URL}/gallery/collection`;
  const options = {method: 'GET'};

  try {
    const response = await fetch(url, options);
    const data = await response.json();
    
    const paths = Object.keys(data);
    const collections = paths.map((path) => {
      return {...data[path], 'uri': path};
    });

    return collections;
  } catch (error) {
    console.error(error);
  }
}

export default {
  async fetch(request, env, ctx) {
    const host = request.headers.get('Host');

    let res = new Response(await siteMapConstruct(`https://${host}`));
    res.headers.set('Content-Type', 'application/xml');

    return res;
  }
};