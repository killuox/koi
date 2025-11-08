Koi 🐟

A beautiful and powerful CLI tool for API testing with a modern terminal interface. Koi allows you to define API endpoints in a configuration file and test them with ease, featuring dynamic parameter generation, variable management, and a sleek terminal UI.

## ✨ Features

- **🎨 Beautiful Terminal UI** - Modern, responsive interface with loading animations and colored output
- **📝 Configuration-Driven** - Define your API endpoints in a simple YAML configuration file
- **🔄 Dynamic Parameters** - Support for environment variables, fake data generation, and command-line flags
- **💾 Variable Management** - Store and reuse response data across requests
- **🎯 Multiple HTTP Methods** - Full support for GET, POST, PUT, PATCH, and DELETE requests
- **📊 Rich Response Display** - Pretty-printed JSON responses with status codes and timing
- **🔧 Flexible Configuration** - Path parameters, query strings, environment variables and request bodies
- **🎲 Fake Data Generation** - Built-in support for generating realistic test data

## 🚀 Installation

### From Source

```bash
git clone https://github.com/killuox/koi.git
cd koi
go build -o koi main.go
```

### Using NPM Install(go installer coming soon)

```bash
npm install -g koi-api
```

## 📖 Quick Start

1. **Create a configuration file** (`koi.config.yaml`):

```yaml
api:
  baseUrl: https://api.example.com
  headers:
    Authorization: Bearer {{token}}
    Content-Type: application/json

endpoints:
  login:
    method: POST
    path: /auth/login
    parameters:
      email:
        type: string
        required: true
        mode: faker:email
      password:
        type: string
        required: true
        mode: faker:password
    set-variables:
      body:
        token: token
    defaults:
      email: user@example.com
      password: password123

  get-users:
    method: GET
    path: /users
    parameters:
      limit:
        type: int
        in: query
        mode: faker:number
        rules:
          min: 1
          max: 100
      page:
        type: int
        in: query
        defaults:
          page: 1

  create-user:
    method: POST
    path: /users
    parameters:
      name:
        type: string
        required: true
        mode: faker:full_name
      email:
        type: string
        required: true
        mode: faker:email
      bio:
        type: string
        mode: faker:paragraph
        rules:
          paragraph_count: 1
          sentence_count: 2
          word_count: 10
```

2. **Run your first API call**:

```bash
# Use default values
koi login

# Override with custom values
koi login --email=admin@example.com --password=secret123

# Use fake data generation
koi create-user

# Pass query parameters
koi get-users --limit=50 --page=2
```

## 🛠️ Configuration

### API Configuration

```yaml
api:
  baseUrl: https://your-api.com
  headers:
    Authorization: Bearer {{token}}
    Content-Type: application/json
    X-API-Key: your-api-key
```

### Endpoint Configuration

Each endpoint supports the following configuration options:

#### Basic Structure
```yaml
endpoints:
  endpoint-name:
    method: GET|POST|PUT|PATCH|DELETE
    path: /api/endpoint
    parameters: # Optional
    defaults: # Optional
    set-variables: # Optional
```

#### Parameters

Parameters can be configured with various modes and types:

```yaml
parameters:
  param-name:
    type: string|int|bool|float
    required: true|false
    in: query|path|body
    mode: env:ENV_VAR|faker:generator
    description: "Parameter description"
    rules: # Optional rules for faker mode
      min: 1
      max: 100
      min_length: 5
      max_length: 50
```

#### Parameter Modes

**Environment Variables:**
```yaml
mode: env:API_KEY
```

**Fake Data Generation:**
```yaml
mode: faker:email
mode: faker:full_name
mode: faker:password
mode: faker:company
mode: faker:phone
mode: faker:number
mode: faker:image
mode: faker:sentence
mode: faker:paragraph
```i

#### Variable Management

Store response data for use in subsequent requests:

```yaml
set-variables:
  body:
    token: token
    user_id: user.id
    session_id: session.id
```

Variables are automatically stored in `~/.koi/variables.json` and can be referenced in headers using `{{variable_name}}` syntax.

## 🎯 Usage Examples

### Basic API Testing

```bash
# Simple GET request
koi health

# POST with parameters
koi login --email=test@example.com --password=secret

# Using environment variables
export API_KEY=your-key
koi protected-endpoint
```

### Advanced Usage

```bash
# Generate fake data for testing
koi create-user  # Uses faker:full_name, faker:email, etc.

# Override faker data with specific values
koi create-user --name="John Doe" --email="john@example.com"

# Use stored variables from previous requests
koi get-profile  # Uses {{token}} from login response
```

### Command Line Flags

Koi supports flexible command-line flag parsing:

```bash
# Long flags
koi endpoint --param=value
koi endpoint --param value

# Short flags  
koi endpoint -p value
koi endpoint -p=value

# Boolean flags
koi endpoint --verbose
koi endpoint -v
```

## 🎨 UI Features

- **Loading Animation** - Beautiful spinner during API calls
- **Status Color Coding** - Green for success (2xx), red for server errors (5xx), yellow for client errors (4xx)
- **Response Pager** - Navigate through large JSON responses with keyboard shortcuts
- **Timing Information** - Request duration display
- **Pretty Printing** - Automatically formatted JSON responses

## 🔧 Advanced Configuration

### Faker Rules

Customize fake data generation with rules:

```yaml
parameters:
  bio:
    type: string
    mode: faker:paragraph
    rules:
      paragraph_count: 2
      sentence_count: 3
      word_count: 15
  
  avatar:
    type: string
    mode: faker:image
    rules:
      width: 200
      height: 200
  
  age:
    type: int
    mode: faker:number
    rules:
      min: 18
      max: 65
```

### Path Parameters

Use dynamic path parameters:

```yaml
endpoints:
  get-user:
    method: GET
    path: /users/{id}
    parameters:
      id:
        type: string
        in: path
        required: true
```

### Query Parameters

Add query string parameters:

```yaml
endpoints:
  search:
    method: GET
    path: /search
    parameters:
      q:
        type: string
        in: query
        required: true
      limit:
        type: int
        in: query
        defaults:
          limit: 10
```

## 📁 Project Structure

```
koi/
├── main.go                 # Entry point
├── koi.config.yaml        # Configuration file
├── internal/
│   ├── api/               # HTTP client and request handling
│   ├── commands/          # CLI command processing
│   ├── config/            # Configuration parsing and validation
│   ├── env/               # Environment variable handling
│   ├── output/            # Terminal UI components
│   ├── shared/            # Shared types and utilities
│   ├── utils/             # Utility functions
│   └── variables/         # Variable management
└── go.mod                 # Go module definition
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- [Charm Bracelet](https://charm.sh/) for the beautiful terminal UI components
- [Go FakeIt](https://github.com/brianvoe/gofakeit) for fake data generation
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the TUI framework

---

Made with ❤️ for developers who love beautiful CLI tools
