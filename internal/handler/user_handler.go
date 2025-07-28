package handler

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/ipxsandbox/internal/entity"
    usecaseUser "github.com/ipxsandbox/internal/usecase/user"
)

type UserHandler struct {
    uc usecaseUser.Usecase
}

func NewUserHandler(uc usecaseUser.Usecase) *UserHandler {
    return &UserHandler{uc: uc}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
    users, err := h.uc.GetAllUsers()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, users)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    var user entity.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    created, err := h.uc.CreateUser(user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, created)
}