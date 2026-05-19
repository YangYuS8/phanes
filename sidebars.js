const sidebars = {
  docs: [
    'intro',
    {
      type: 'category',
      label: 'ADRs',
      items: [
        'adr/scope-and-boundaries',
        'adr/process-model',
        'adr/contract-packages',
        'adr/resource-cache-format',
      ],
    },
    {
      type: 'category',
      label: 'Contracts',
      items: [
        'contracts/module-boundaries',
        'contracts/protobuf',
        'contracts/go-interfaces',
        'contracts/sqlite-schema',
        'contracts/cli',
        'contracts/runtime-http-api',
        'contracts/builder-cache',
        'contracts/documentation-site',
        'contracts/release',
        'contracts/compliance-checklist',
      ],
    },
  ],
};

module.exports = sidebars;
