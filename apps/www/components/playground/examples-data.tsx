import {
  FileJson,
  Database,
  Server,
  Users,
  Globe,
  Map as MapIcon,
  Settings,
  Zap,
  LayoutTemplate,
  Flag,
  Grid3x3,
  ShoppingCart,
  Table,
  Layers,
  Sparkles,
} from "lucide-react";

interface Example {
  id: string;
  title: string;
  description: string;
  category: "Basics" | "Infrastructure" | "Data" | "Config";
  icon: React.ElementType;
  code: string;
}

export const EXAMPLES: Example[] = [
  {
    id: "demo",
    title: "Basic Structure",
    description: "A minimal example showing core JSSON syntax.",
    category: "Basics",
    icon: Zap,
    code: `profile {
  name = Lucas
  age = 27
  verified = false
  interests = [ "tech", "gaming", "fitness" ]

  preferences {
    theme = light
    notifications = true
  }
}`,
  },
  {
    id: "template",
    title: "Reusable Templates",
    description: "Define structures once and reuse them across datasets.",
    category: "Basics",
    icon: LayoutTemplate,
    code: `employees [
  template { name, age, department, salary }
  
  Clara, 32, HR, 4800
  Renato, 28, IT, 5200
  Bianca, 35, Finance, 6100
  Diego, 24, Marketing, 3900
]`,
  },

  {
    id: "map-advanced",
    title: "Mapping Logic",
    description: "Apply inline transformations and computed fields.",
    category: "Basics",
    icon: FileJson,
    code: `inventory [
  template { id, name, qty, category }

  map (item) = {
    id = "item-" + item.id
    name = item.name
    category = item.category
    qty = item.qty
    status = item.qty > 10 ? "in-stock" : "low-stock"
    isEssential = item.category == "food"
  }

  1, "Rice Pack", 25, food
  2, "Dish Soap", 8, cleaning
  3, "LED Lamp", 12, electronics
  4, "Coffee Beans", 6, food
]`,
  },

  {
    id: "database",
    title: "Database Settings",
    description: "Readable and structured DB configuration.",
    category: "Config",
    icon: Database,
    code: `db {
  engine = "postgres"
  host = "db.internal"
  port = 5432
  ssl = true
}`,
  },

  {
    id: "user-gen",
    title: "User Seeder",
    description: "Automatically generate users based on ranges.",
    category: "Data",
    icon: Users,
    code: `accounts [
  template { id, team }

  map (a) = {
    id = a.id
    username = "member_" + a.id
    email = a.id + "@example.org"
    team = a.team
    active = true
  }

  // Engineering team
  300..304, engineering

  // Design team
  400..402, design
]`,
  },

  {
    id: "feature-flags",
    title: "Environment Toggles",
    description: "Simple and scalable feature toggle management.",
    category: "Config",
    icon: Flag,
    code: `featureToggles [
  template { flag, env, value }

  map (f) = {
    key = f.flag + ":" + f.env
    enabled = f.value
    rollout = f.value == true ? 100 : 0
  }

  "beta-search", "prod", false
  "beta-search", "qa", true
  "beta-search", "dev", true

  "realtime-notify", "prod", true
  "realtime-notify", "qa", true
]`,
  },

  {
    id: "k8s",
    title: "Kubernetes Pods",
    description: "Generate multiple environment deployment configs.",
    category: "Infrastructure",
    icon: Server,
    code: `services [
  template { name, env, scale }

  map (svc) = {
    apiVersion = "apps/v1"
    kind = "Deployment"
    name = svc.name + "-" + svc.env
    replicas = svc.scale

    container = {
      image = "registry/" + svc.name + ":" + svc.env
      restartPolicy = "Always"
    }
  }

  "web", "prod", 4
  "web", "staging", 2
  "web", "dev", 1

  "jobs", "prod", 3
  "jobs", "staging", 1
]`,
  },

  {
    id: "i18n",
    title: "i18n Packs",
    description: "Manage interface translations in multiple languages.",
    category: "Config",
    icon: Globe,
    code: `translations [
  template { key, locale, text }

  map (t) = {
    key = t.key
    locale = t.locale
    value = t.text
  }

  "login.header", "en", "Welcome Back"
  "login.header", "pt", "Bem-vindo"
  "login.header", "fr", "Bienvenue"

  "btn.submit", "en", "Submit"
  "btn.submit", "pt", "Enviar"
  "btn.submit", "es", "Enviar"
]`,
  },

  {
    id: "api-gateway",
    title: "Gateway Routes",
    description: "Conditional routing and dynamic API generation.",
    category: "Infrastructure",
    icon: Settings,
    code: `apiRoutes [
  template { name, version, path }

  map (r) = {
    id = r.name + "-" + r.version
    fullPath = "/api/v" + r.version + "/" + r.path
    internalService = r.name + ".svc.local"
    secure = r.name != "status"
  }

  "auth", 1, "login"
  "auth", 2, "session"
  "payment", 1, "checkout"
  "status", 1, "ping"
]`,
  },

  {
    id: "geo",
    title: "Geo Grid",
    description: "Math-heavy dataset generation for geospatial apps.",
    category: "Data",
    icon: MapIcon,
    code: `geoPoints [
  template { id, tag }

  map (g) = {
    id = "pt-" + g.id
    lat = -12.9000 + (g.id / 50) * 0.02
    lon = -38.3300 + (g.id % 50) * 0.02
    tag = g.tag
    type = "geo_point"
  }

  // Generate 120 points
  0..119, "zone-a"
]`,
  },
  {
    id: "nested-maps-matrix",
    title: "Matrix Generation",
    description:
      "Generate 2D matrices using nested map transformations (> v0.0.5).",
    category: "Data",
    icon: Grid3x3,
    code: `// Multiplication table using nested maps
table = (1..5 map (row) = (1..5 map (col) = row * col))

// Coordinate grid
grid = (0..2 map (y) = (0..2 map (x) = {
  x = x
  y = y
  id = y * 3 + x
}))`,
  },

  {
    id: "product-variants",
    title: "Product Variants",
    description:
      "Generate all size/color combinations for e-commerce (> v0.0.5).",
    category: "Data",
    icon: ShoppingCart,
    code: `// All product variants
products = (["S", "M", "L", "XL"] map (size) = (
  ["Black", "White", "Navy", "Gray"] map (color) = {
    sku = size + "-" + color
    size = size
    color = color
    price = 29.99
    inStock = true
  }
))`,
  },

  {
    id: "nested-arrays",
    title: "Nested Arrays",
    description: "Multi-dimensional arrays and matrices (> v0.0.5).",
    category: "Data",
    icon: Layers,
    code: `// 2D Matrix
matrix = [
  [ 1, 2, 3 ],
  [ 4, 5, 6 ],
  [ 7, 8, 9 ]
]

// Permissions matrix
permissions = {
  admin = [
    [ "read", "write", "delete" ],
    [ "create", "update", "admin" ]
  ]
  user = [
    [ "read" ],
    [ "create", "update" ]
  ]
}`,
  },

  {
    id: "large-scale-gen",
    title: "Large-Scale Generation",
    description: "Generate thousands of records effortlessly (> v0.0.5).",
    category: "Data",
    icon: Sparkles,
    code: `// Generate 1000 test users
testUsers = (0..999 map (id) = {
  id = id
  username = "user_" + id
  email = "user" + id + "@test.com"
  active = (id % 2) == 0
  score = id * 10
  tier = id < 100 ? "bronze" : id < 500 ? "silver" : "gold"
})`,
  },

  {
    id: "schedule-matrix",
    title: "Schedule Matrix",
    description: "Generate time slots for scheduling apps (> v0.0.5).",
    category: "Data",
    icon: Table,
    code: `// Weekly schedule (7 days × 8 hours)
schedule = (["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"] map (day) = (
  9..16 map (hour) = {
    day = day
    hour = hour
    slot = day + "-" + hour + "h"
    available = day != "Sat" && day != "Sun"
    price = day == "Sat" || day == "Sun" ? 150 : 100
  }
))`,
  },
];
