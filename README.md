# 🚦 Rate Limiting Algorithms in Go

A simple project to **experiment with different rate limiting algorithms** such as:

- **Token Bucket**
- **Per-Client Rate Limiting**
- **Toll-Boothing**

These algorithms help control request traffic and protect your server from being overwhelmed.

---

## 🛠️ How to Test

### 1. Navigate to the specific directory:

```bash
cd <algorithm-directory>
```

```bash
go run .
```
### 2. Send multiple requests:
For Windows Command Prompt:
```bash
FOR /L %i IN (1,1,10) DO curl -i http://localhost:8080/ping
```
