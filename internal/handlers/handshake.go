package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandshakeGet(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, `
<!DOCTYPE html>
<html>
<head><title>Authorizing…</title></head>
<body>
  <p style="font: 16px sans-serif">Authorizing…</p>
  <button id="ok" style="display:none">Continue</button>
  <script>
    var isWebKit = /Safari/.test(navigator.userAgent) && !/Chrome/.test(navigator.userAgent);
    if (isWebKit) {
        document.getElementById('ok').style.display = 'block';
        document.getElementById('ok').onclick = function () {
            window.location.replace('https://crossword-frontend-one.vercel.app/');
        };
    } else {
        window.location.replace('https://crossword-frontend-one.vercel.app/');
    }
  </script>
</body>
</html>`)
}
