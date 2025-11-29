export const codeExamples = [
    {
        title: "Basic Configuration",
        jsson: `// Simple configuration
api {
  url = "https://api.example.com"
  timeout = 5000
  retries = 3
  
  headers {
    content_type = "application/json"
    auth_token = "bearer-token-123"
  }
}`,
        json: `{
  "api": {
    "url": "https://api.example.com",
    "timeout": 5000,
    "retries": 3,
    "headers": {
      "content_type": "application/json",
      "auth_token": "bearer-token-123"
    }
  }
}`,
        yaml: `api:
    headers:
        auth_token: bearer-token-123
        content_type: application/json
    retries: 3
    timeout: 5000
    url: https://api.example.com`,
        toml: `[api]
  retries = 3
  timeout = 5000
  url = "https://api.example.com"
  [api.headers]
    auth_token = "bearer-token-123"
    content_type = "application/json"`,
        typescript: `export const api = {
  url: "https://api.example.com",
  timeout: 5000,
  retries: 3,
  headers: {
    content_type: "application/json",
    auth_token: "bearer-token-123"
  }
} as const;

// Generated types
export type Api = typeof api;`,
    },
    {
        title: "Product Variants (v0.0.5)",
        jsson: `// Generate all size/color combinations
products = (["S", "M", "L"] map (size) = (
  ["Black", "White", "Navy"] map (color) = {
    sku = size + "-" + color
    size = size
    color = color
    price = 29.99
    inStock = true
  }
))`,
        json: `{
  "products": [
    [
      { "sku": "S-Black", "size": "S", "color": "Black", "price": 29.99, "inStock": true },
      { "sku": "S-White", "size": "S", "color": "White", "price": 29.99, "inStock": true },
      { "sku": "S-Navy", "size": "S", "color": "Navy", "price": 29.99, "inStock": true }
    ],
    [
      { "sku": "M-Black", "size": "M", "color": "Black", "price": 29.99, "inStock": true },
      { "sku": "M-White", "size": "M", "color": "White", "price": 29.99, "inStock": true },
      { "sku": "M-Navy", "size": "M", "color": "Navy", "price": 29.99, "inStock": true }
    ],
    [
      { "sku": "L-Black", "size": "L", "color": "Black", "price": 29.99, "inStock": true },
      { "sku": "L-White", "size": "L", "color": "White", "price": 29.99, "inStock": true },
      { "sku": "L-Navy", "size": "L", "color": "Navy", "price": 29.99, "inStock": true }
    ]
  ]
}`,
        yaml: `products:
    - - color: Black
        inStock: true
        price: 29.99
        size: S
        sku: S-Black
      - color: White
        inStock: true
        price: 29.99
        size: S
        sku: S-White
      - color: Navy
        inStock: true
        price: 29.99
        size: S
        sku: S-Navy
    - - color: Black
        inStock: true
        price: 29.99
        size: M
        sku: M-Black
      # ... M-White, M-Navy
    - - color: Black
        inStock: true
        price: 29.99
        size: L
        sku: L-Black
      # ... L-White, L-Navy`,
        toml: `[[products]]
  [[products.products]]
    color = "Black"
    inStock = true
    price = 29.99
    size = "S"
    sku = "S-Black"

  [[products.products]]
    color = "White"
    inStock = true
    price = 29.99
    size = "S"
    sku = "S-White"

  [[products.products]]
    color = "Navy"
    inStock = true
    price = 29.99
    size = "S"
    sku = "S-Navy"

# ... M and L variants`,
        typescript: `export const products = [
  [
    { sku: "S-Black", size: "S", color: "Black", price: 29.99, inStock: true },
    { sku: "S-White", size: "S", color: "White", price: 29.99, inStock: true },
    { sku: "S-Navy", size: "S", color: "Navy", price: 29.99, inStock: true }
  ],
  [
    { sku: "M-Black", size: "M", color: "Black", price: 29.99, inStock: true },
    { sku: "M-White", size: "M", color: "White", price: 29.99, inStock: true },
    { sku: "M-Navy", size: "M", color: "Navy", price: 29.99, inStock: true }
  ],
  [
    { sku: "L-Black", size: "L", color: "Black", price: 29.99, inStock: true },
    { sku: "L-White", size: "L", color: "White", price: 29.99, inStock: true },
    { sku: "L-Navy", size: "L", color: "Navy", price: 29.99, inStock: true }
  ]
] as const;

// Generated types
export type Products = typeof products;`,
    },
    {
        title: "Multi-Region Deployment",
        jsson: `servers [
  template { id, region, tier }
  
  map (s) = {
    id = "srv-" + s.id
    region = s.region
    tier = s.tier
    ip = "10." + (s.id / 100) + ".0." + (s.id % 100)
    replicas = s.tier == "prod" ? 5 : 2
  }
  
  // US East - Production
  100..109, "us-east-1", prod
  
  // US West - Staging
  200..204, "us-west-2", staging
  
  // EU - Production
  300..314, "eu-central-1", prod
  
  // APAC - Dev
  400..402, "ap-south-1", dev
]`,
        json: `
{
  "servers": [
    {
      "id": "srv-100",
      "region": "us-east-1",
      "tier": "prod",
      "ip": "10.1.0.0",
      "replicas": 5
    },
    {
      "id": "srv-101",
      "region": "us-east-1",
      "tier": "prod",
      "ip": "10.1.0.1",
      "replicas": 5
    },
    {
      "id": "srv-102",
      "region": "us-east-1",
      "tier": "prod",
      "ip": "10.1.0.2",
      "replicas": 5
    },
    //... more servers
    {
      "id": "srv-401",
      "ip": "10.4.01.0.1",
      "region": "ap-south-1",
      "replicas": 2,
      "tier": "dev"
    },
    {
      "id": "srv-402",
      "ip": "10.4.02.0.2",
      "region": "ap-south-1",
      "replicas": 2,
      "tier": "dev"
    }
  ]
}`,
        yaml: `
servers:
    - id: srv-100
      ip: 10.1.0.0
      region: us-east-1
      replicas: 5
      tier: prod
    - id: srv-101
      ip: 10.1.01.0.1
      region: us-east-1
      replicas: 5
      tier: prod
    - id: srv-102
      ip: 10.1.02.0.2
      region: us-east-1
      replicas: 5
      tier: prod
    //... more servers
    - id: srv-401
      ip: 10.4.01.0.1
      region: ap-south-1
      replicas: 2
      tier: dev
    - id: srv-402
      ip: 10.4.02.0.2
      region: ap-south-1
      replicas: 2
      tier: dev
    `,
        toml: `
[[servers]]
  id = "srv-100"
  ip = "10.1.0.0"
  region = "us-east-1"
  replicas = 5
  tier = "prod"

[[servers]]
  id = "srv-101"
  ip = "10.1.01.0.1"
  region = "us-east-1"
  replicas = 5
  tier = "prod"

[[servers]]
  id = "srv-102"
  ip = "10.1.02.0.2"
  region = "us-east-1"
  replicas = 5
  tier = "prod"

//... more servers

[[servers]]
  id = "srv-401"
  ip = "10.4.01.0.1"
  region = "ap-south-1"
  replicas = 2
  tier = "dev"

[[servers]]
  id = "srv-402"
  ip = "10.4.02.0.2"
  region = "ap-south-1"
  replicas = 2
  tier = "dev"
`,
        typescript: `
export const servers = [
  {
    tier: "prod",
    ip: "10.1.0.0",
    replicas: 5,
    id: "srv-100",
    region: "us-east-1"
  },
  {
    region: "us-east-1",
    tier: "prod",
    ip: "10.1.01.0.1",
    replicas: 5,
    id: "srv-101"
  },
  {
    id: "srv-102",
    region: "us-east-1",
    tier: "prod",
    ip: "10.1.02.0.2",
    replicas: 5
  },
  //... more servers
  {
    id: "srv-400",
    region: "ap-south-1",
    tier: "dev",
    ip: "10.4.0.0",
    replicas: 2
  },
  {
    id: "srv-401",
    region: "ap-south-1",
    tier: "dev",
    ip: "10.4.01.0.1",
    replicas: 2
  },
  {
    id: "srv-402",
    region: "ap-south-1",
    tier: "dev",
    ip: "10.4.02.0.2",
    replicas: 2
  }
] as const;

// Generated types
export type Servers = typeof servers;

    `,
    },
    {
        title: "Dynamic Resources",
        jsson: `
resources [
template { name, type, region }

  map (res) = {
    // Auto-generate standardized ID
    id = res.type + "-" + res.name + "-" + res.region
    
    name = res.name
    type = res.type
    region = res.region
    
    tags {
      managed_by = "jsson"
      env = "production"
    }
  }

  "web-server", ec2, "us-east-1" 
  "db-primary", rds, "us-west-2"
]`,
        json: `{
"resources": [
    {
      "id": "ec2-web-server-us-east-1",
      "name": "web-server",
      "type": "ec2",
      "region": "us-east-1",
      "tags": {
        "managed_by": "jsson",
        "env": "production"
      }
    },
    {
      "id": "rds-db-primary-us-west-2",
      "name": "db-primary",
      "type": "rds",
      "region": "us-west-2",
      "tags": {
        "managed_by": "jsson",
        "env": "production"
      }
    }
  ]
}`,
        yaml: `
resources:
    - id: ec2-web-server-us-east-1
      name: web-server
      region: us-east-1
      tags:
        env: production
        managed_by: jsson
      type: ec2
    - id: rds-db-primary-us-west-2
      name: db-primary
      region: us-west-2
      tags:
        env: production
        managed_by: jsson
      type: rds    
    `,
        toml: `
[[resources]]
    id = "ec2-web-server-us-east-1"
    name = "web-server"
    region = "us-east-1"
    tags = { env = "production", managed_by = "jsson" }
    type = "ec2"

[[resources]]
    id = "rds-db-primary-us-west-2"
    name = "db-primary"
    region = "us-west-2"
    tags = { env = "production", managed_by = "jsson" }
    type = "rds"
    `,
        typescript: `
export const resources = [
  {
    id: "ec2-web-server-us-east-1",
    name: "web-server",
    region: "us-east-1",
    tags: {
      env: "production",
      managed_by: "jsson"
    },
    type: "ec2"
  },
  {
    id: "rds-db-primary-us-west-2",
    name: "db-primary",
    region: "us-west-2",
    tags: {
      env: "production",
      managed_by: "jsson"
    },
    type: "rds"
  }
] as const;

// Generated types
export type Resources = typeof resources;
    `,
    },
];
