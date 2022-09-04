# Toss-AutoCharge
> Autocharge Web Server for toss.me<br/>
Btw, **u should randomize sender name** in your bot/client

Create an issue or pr if there is a bug or smth

# Installation
`go run main.go`<br/>
ima add more info soon

# Usage
POST http://127.0.0.1:8080/check<br/>
Content-Type JSON<br/>
Body
```js
{
    "tossId": "Toss-ID HERE",
    "name": "senderName",
    "amount": Number
}
```

Response:
```js
{
    "success": Boolean,
    "found": Boolean,
    "message": "an error message if success is false"
}
```

btw i made it guys!!!!
