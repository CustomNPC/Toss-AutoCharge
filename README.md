# Toss-AutoCharge
> Autocharge Web Server for toss.me<br/>
Btw, **u should randomize sender name** in your bot/client

Create a issue/pr if there is a bug or smth

# Installation
just ran `go run main.go`<br/>
ima add more info soon

# Usage
POST http://127.0.0.1:8080/check<br/>
Content-Type JSON<br/>
Body
```json
{
    "tossId": "Toss-ID HERE",
    "name": "senderName",
    "amount": amount
}
```

Response:
```json
{
  "success": boolean,
  "found": boolean,
  "message": "an error message if success is false"
}
```
