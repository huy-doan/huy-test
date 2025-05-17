# MSP API Project

## I. Giới thiệu

Makeshop Payment API sử dụng Golang, Gin, GraphQL, MySQL và thiết kế theo Domain-Driven Design (DDD).

### Công nghệ sử dụng

- **Go (Golang)** - Ngôn ngữ lập trình chính
- **Gin** - Framework HTTP
- **GORM** - ORM cho MySQL
- **GraphQL** - API query language
- **JWT** - Xác thực
- **MySQL** - Cơ sở dữ liệu
- **Docker** - Containerization
- **GoMigration** - Quản lý phiên bản cơ sở dữ liệu

## II. Cấu trúc dự án

Dự án được tổ chức theo mô hình Domain-Driven Design (DDD), phân tách rõ ràng các lớp và trách nhiệm.

```
└── 📁msp-be
    ├── 📁cmd                # Chứa các lệnh CLI khác nhau
    │   ├── example_shell.go
    │   ├── 📁ExampleShell
    │   │   └── main.go
    │   └── root.go
    ├── 📁config             # Cấu hình ứng dụng
    │   └── 📁migrations     # SQL migrations cho cơ sở dữ liệu
    │       └── 20250316004749_create-users-table.sql
    ├── 📁docs               # Tài liệu API được tạo bởi Swagger
    │   ├── docs.go
    │   ├── swagger.json
    │   └── swagger.yaml
    ├── 📁ops                # Cấu hình vận hành, docker
    │   ├── 📁go
    │   │   ├── .env
    │   │   └── Dockerfile
    │   └── 📁mysql
    │       ├── Dockerfile
    │       ├── init.sql
    │       └── my.cnf
    ├── 📁src                # Mã nguồn chính
    │   ├── 📁api            # API handlers (HTTP và GraphQL)
    │   │   ├── 📁graphql    # Xử lý GraphQL API
    │   │   │   ├── 📁generated
    │   │   │   │   ├── generated.go
    │   │   │   │   └── models_gen.go
    │   │   │   ├── 📁middleware
    │   │   │   │   └── auth.go
    │   │   │   ├── 📁resolvers
    │   │   │   │   ├── mutation.resolvers.go
    │   │   │   │   ├── query.resolvers.go
    │   │   │   │   ├── resolver.go
    │   │   │   │   └── type.resolvers.go
    │   │   │   ├── 📁schema
    │   │   │   │   ├── input.graphql
    │   │   │   │   ├── mutation.graphql
    │   │   │   │   ├── query.graphql
    │   │   │   │   └── type.graphql
    │   │   │   └── server.go
    │   │   ├── 📁http      # Xử lý REST API với Gin framework
    │   │   │   ├── 📁handlers
    │   │   │   │   ├── auth_handler.go
    │   │   │   │   └── user_handler.go
    │   │   │   ├── 📁middleware
    │   │   │   │   └── auth.go
    │   │   │   └── route.go
    │   │   └── server.go
    │   ├── 📁application    # Lớp ứng dụng (services, commands, queries)
    │   │   └── 📁services
    │   │       └── user_service.go
    │   ├── 📁domain         # Lớp domain (entities, repositories interfaces)
    │   │   ├── 📁entities
    │   │   │   └── user.go
    │   │   └── 📁repositories
    │   │       └── user_repository.go
    │   ├── 📁infrastructure # Cài đặt cụ thể (auth, persistence, validation)
    │   │   ├── 📁auth
    │   │   │   └── jwt_service.go
    │   │   └── 📁persistence
    │   │       ├── 📁mysql
    │   │       │   ├── connection.go
    │   │       │   └── models.go
    │   │       └── 📁repositories
    │   │           └── user_repository.go
    │   └── main.go
    ├── .air.toml            # Cấu hình Air (hot-reloading)
    ├── .env.example         # Mẫu cấu hình môi trường
    ├── .gitignore
    ├── go.mod               # Go modules
    ├── go.sum
    ├── gqlgen.yml           # Cấu hình GraphQL generator
    ├── main.go
    ├── Makefile             # Các lệnh make
    ├── README.md
    └── tools.go             # Công cụ phát triển
```

### Mô tả các thành phần chính:

#### 1. Lớp Domain (src/domain)
- **entities**: Định nghĩa các đối tượng lõi của hệ thống
- **repositories**: Định nghĩa interfaces cho việc tương tác với dữ liệu

#### 2. Lớp Application (src/application)
- **services**: Xử lý logic nghiệp vụ
- **commands**: Xử lý thay đổi dữ liệu (write operations)
- **queries**: Xử lý truy vấn dữ liệu (read operations)

#### 3. Lớp Infrastructure (src/infrastructure)
- **auth**: Triển khai xác thực JWT
- **persistence**: Cài đặt repositories và kết nối cơ sở dữ liệu
- **validation**: Xác thực dữ liệu đầu vào

#### 4. Lớp API (src/api)
- **http**: Xử lý REST API với Gin framework
  - **handlers**: Chứa các handler xử lý request HTTP
  - **middleware**: Middleware xác thực và bảo mật
  - **route.go**: Định nghĩa các route
- **graphql**: Xử lý GraphQL API
  - **generated**: Mã được tạo tự động từ schema GraphQL
  - **middleware**: Middleware xác thực và bảo mật
  - **resolvers**: Các resolver cho GraphQL
  - **schema**: Schema GraphQL (input, mutation, query, type)

#### 5. Cấu hình (config)
- **migrations**: Quản lý phiên bản cơ sở dữ liệu

#### 6. Operations (ops)
- **go**: Cấu hình Docker cho ứng dụng Go
- **mysql**: Cấu hình Docker cho MySQL, bao gồm script khởi tạo và cấu hình

## III. Hướng dẫn sử dụng Makefile

Dự án cung cấp nhiều lệnh make để đơn giản hóa các tác vụ phát triển:

```bash
# Hiển thị danh sách các lệnh có sẵn
make help

# Chạy API trong môi trường phát triển
make run

# Tạo mã GraphQL từ schema
make generate-graphql

# Tạo tài liệu Swagger
make swagger

# Format mã nguồn
make fmt

# Tạo migration mới
make migrate-create

# Chạy migrations up
make migrate-up

# Revert migration (down)
make migrate-down

# Chạy lệnh shell trong container
make shell
```

### Ví dụ

```bash
# Tạo migration mới
make migrate-create
# Nhập tên migration: create_users_table

# Chạy migrations
make migrate-up

# Tạo mã GraphQL sau khi sửa đổi schema
make generate-graphql
```

## Bắt đầu

1. Khởi động các dịch vụ với Docker Compose:

```bash
cd msp-infra
docker-compose build
docker-compose up -d
```

2. Truy cập:
   - REST API: http://localhost:3010/api/v1
   - Swagger UI: http://localhost:3010/swagger/index.html
   - GraphQL Playground: http://localhost:3010/graphql (chỉ trong môi trường phát triển)