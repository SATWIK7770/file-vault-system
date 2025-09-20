package api

import (
    "backend/internal/service"
    "errors"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
)


type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}


func (h *AuthHandler) SignUp(c *gin.Context) {
    var req struct {
        Username string `json:"username"`
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "Invalid request body",
        })
        return
    }

    // Validation
    if len(req.Username) < 3 {
        c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Username must be at least 3 characters"})
        return
    }
    if !strings.Contains(req.Email, "@") || len(req.Email) < 5 {
        c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid email address"})
        return
    }
    if len(req.Password) < 8 {
        c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Password must be at least 8 characters"})
        return
    }

    userID, token, err := h.authService.SignUpAndGenerateToken(req.Username, req.Email, req.Password)
    if err != nil {
        switch {
        case errors.Is(err, service.ErrUserExists):
            c.JSON(http.StatusConflict, gin.H{"success": false, "message": "Username already exists"})
        case errors.Is(err, service.ErrEmailExists):
            c.JSON(http.StatusConflict, gin.H{"success": false, "message": "Email already registered"})
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Could not create user"})
        }
        return
    }

    
    c.SetCookie(
        "auth_token",
        token,
        900,   // expires in 15 min (seconds)
        "/",     // cookie valid for all paths
        "localhost",      // domain (set domain in production)
        false,    // secure (set true in production with HTTPS)
        true,    // httpOnly 
    )

    // Respond without token
    c.JSON(http.StatusCreated, gin.H{
        "success": true,
        "message": "User created successfully",
        "user_id": userID,
    })
}

func (h *AuthHandler) SignIn(c *gin.Context) {
    var req struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    token, userID, err := h.authService.SignIn(req.Username, req.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    c.SetCookie(
        "auth_token", 
        token,
        900,      // expires in 1 day
        "/",          // path
        "localhost",  // domain
        false,         // secure 
        true,         // HttpOnly
    )

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "user_id": userID,
    })
}


func (h *AuthHandler) SignOut(c *gin.Context) {
	// Overwrite cookie
	c.SetCookie(
		"auth_token",   // cookie name 
		"",        // empty value
		-1,        // maxAge negative
		"/",       // path
		"",        // domain 
		true,      // secure
		true,      // httpOnly
	)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Logged out"})
}


func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetUint("userID")
	user, err := h.authService.GetUserByID(userID)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
	})
}