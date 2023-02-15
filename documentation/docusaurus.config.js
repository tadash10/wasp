const path = require('path');

module.exports = {
  url: "https://foo.bar",
  title: "asd",
  baseUrl: "/",

  plugins: [
    [
      'docusaurus-plugin-openapi-docs',
      {
        id: "webapi",
        docsPluginId: "classic",
        config: {
          webapi: { // Note: petstore key is treated as the <id> and can be used to specify an API doc instance when using CLI commands
            specPath: "../clients/apiclient/api/openapi.yaml",
            downloadUrl: "https://raw.githubusercontent.com/iotaledger/wasp/develop/clients/apiclient/api/openapi.yaml", // Path to designated spec file
            outputDir: "docs/webapi", // Output directory for generated .mdx docs
            sidebarOptions: {
              groupPathsBy: "tag",
              categoryLinkSource: "tag",
            },
          }
        }
      },
    ],

    [
      '@docusaurus/plugin-content-docs',
      {
        id: 'wasp',
        path: path.resolve(__dirname, 'docs'),
        routeBasePath: 'smart-contracts',
        sidebarPath: path.resolve(__dirname, 'sidebars.js'),
        editUrl: 'https://github.com/iotaledger/wasp/edit/develop/documentation',
        remarkPlugins: [require('remark-code-import'), require('remark-import-partial'), require('remark-remove-comments')],
      },
    ],

  ],

  staticDirectories: [path.resolve(__dirname, 'static')],
};
