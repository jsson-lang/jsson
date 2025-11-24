# 🚀 **JSSON v0.0.2 – The Logic & Arithmetic Release**

JSSON (**JavaScript Simplified Object Notation**) is a modern, human-friendly syntax that transpiles directly to JSON — with templates, ranges, includes, maps, and now **full arithmetic and conditional logic**.

---

# ✨ **🎯 New in v0.0.2**

### 🧠 **1. Conditional Logic (NEW!)**

JSSON now supports powerful decision-making syntax directly in your configuration files.

**Comparison Operators:**
`==`, `!=`, `>`, `<`, `>=`, `<=`

**Ternary Operator:**
Conditional values are now possible with `? :` syntax.

```jsson
users [
  template { name, age, role, access }
  
  map (u) = {
    name = u.name
    age = u.age
    role = u.role
    // Conditional logic in action
    isAdult = u.age >= 18
    access = u.role == "admin" ? "full" : "read-only"
    category = u.age < 13 ? "child" : u.age < 20 ? "teen" : "adult"
  }
  
  Alice, 25, admin
  Bob, 12, user
]
```

---

### 🧮 **2. Arithmetic Expressions (NEW!)**

Perform calculations right inside your JSSON files. No more manual math!

**Operators:** `+`, `-`, `*`, `/`, `%` (modulo)

**Mixed Type Support:**
Seamlessly mix integers and floats. JSSON handles the promotion automatically.

```jsson
products [
  template { price, discount }
  
  map (p) = {
    originalPrice = p.price
    // Arithmetic with mixed types (int * float)
    finalPrice = p.price * (1.0 - p.discount)
    savings = p.price - finalPrice
  }
  
  100, 0.15
  50, 0.0
]
```

---

### � **3. Advanced Ranges & Zipping (IMPROVED!)**

Ranges are now more robust than ever.

*   **Smart Zipping:** You can now mix ranges of different lengths in templates. JSSON automatically zips them up to the shortest length.
*   **String Ranges:** `"server-1" .. "server-5"` works out of the box.

```jsson
servers [
  template { id, ip, region }
  
  // Different length ranges? No problem!
  // Zips until the shortest range ends.
  1..100, "192.168.1." + (10..50), "us-east-1"
]
```

---

# 📦 **Installation**

Download the binary for your OS:

*   **Windows:** `jsson-v0.0.2-windows-amd64.exe`
*   **Linux:** `jsson-v0.0.2-linux-amd64`
*   **macOS (Intel):** `jsson-v0.0.2-darwin-amd64`
*   **macOS (Apple Silicon):** `jsson-v0.0.2-darwin-arm64`

### Linux/macOS

```bash
chmod +x jsson-v0.0.2-*
sudo mv jsson-v0.0.2-* /usr/local/bin/jsson
```

### Windows

```powershell
Rename-Item jsson-v0.0.2-windows-amd64.exe jsson.exe
```

---

# 🚀 **Usage**

```bash
jsson -i input.jsson
```

---

# 🎨 **VS Code Extension**

Official syntax highlighting + language support:

👉 [https://marketplace.visualstudio.com/items?itemName=carlosedujs.jsson](https://marketplace.visualstudio.com/items?itemName=carlosedujs.jsson)

---

# 📚 Documentation

Docs:
👉 [https://github.com/carlosedujs/jsson](https://github.com/carlosedujs/jsson)

---

# 🙏 Acknowledgments

Thanks to everyone helping shape this language.
Special thanks to the wizards, goblins and gremlins of the parser.
