# GO-Echo Starter Code

## Setup

### 1. Install prerequisites

#### Install Go

Make sure Go is installed:

```bash
go version
```

If missing (Linux Mint / Ubuntu):

```bash
sudo apt update
sudo apt install golang-go
```

---

### 2. Install required Go tools

Install Air (hot reload):

```bash
go install github.com/air-verse/air@latest
```

Install Goose (database migrations):

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Add Go binaries to PATH:

```bash
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
source ~/.bashrc
```

Verify:

```bash
air -v
goose -version
```

---

### 3. Setup MySQL (Docker)

Run MySQL container:

```bash
docker run --name go_echo_starter \
  -e MYSQL_ROOT_PASSWORD=my-secret-pw \
  -e MYSQL_DATABASE=go_echo_starter \
  -p 3306:3306 \
  -d mysql:latest
```

---

### 4. Configure environment variables

Create or edit `.env` file:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=my-secret-pw
DB_NAME=go_echo_starter
DB_SSLMODE=disable
```

---

### 5. Install Go dependencies

```bash
make deps
```

---

### 6. Run database migrations

```bash
make migrate-up
```

---

### 7. Run the application

```bash
make run
```

---

### 8. Build binary

```bash
make build
```


---

### 9. Add swagger docs
```base
swag init -g cmd/api/main.go
```

---

## Common issues

* `goose: command not found` → Go bin path not added to PATH
* DB connection fails → MySQL container not running or wrong credentials
* `localhost` issues → ensure app is running on host, not inside Docker
