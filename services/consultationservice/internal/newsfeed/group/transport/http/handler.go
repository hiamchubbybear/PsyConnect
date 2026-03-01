package http

import (
	"consultationservice/internal/newsfeed/group/domain"
	"consultationservice/internal/newsfeed/group/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	createGroupUC  *usecase.CreateGroupUseCase
	getGroupsUC    *usecase.GetGroupsUseCase
	getGroupByIDUC *usecase.GetGroupByIDUseCase
	joinGroupUC    *usecase.JoinGroupUseCase
}

func NewHandler(
	createGroupUC *usecase.CreateGroupUseCase,
	getGroupsUC *usecase.GetGroupsUseCase,
	getGroupByIDUC *usecase.GetGroupByIDUseCase,
	joinGroupUC *usecase.JoinGroupUseCase,
) *Handler {
	return &Handler{
		createGroupUC:  createGroupUC,
		getGroupsUC:    getGroupsUC,
		getGroupByIDUC: getGroupByIDUC,
		joinGroupUC:    joinGroupUC,
	}
}

func (h *Handler) CreateGroup(c *gin.Context) {
	var group domain.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	group.CreatorID = userID

	if err := h.createGroupUC.Execute(c.Request.Context(), &group); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, group)
}

func (h *Handler) GetGroups(c *gin.Context) {
	category := c.Query("category")
	query := c.Query("q")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)
	skip, _ := strconv.ParseInt(c.DefaultQuery("skip", "0"), 10, 64)

	groups, err := h.getGroupsUC.Execute(c.Request.Context(), category, query, limit, skip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, groups)
}

func (h *Handler) JoinGroup(c *gin.Context) {
	groupID := c.Param("id")
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.joinGroupUC.Execute(c.Request.Context(), groupID, userID); err != nil {
		if err == usecase.ErrGroupNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err == usecase.ErrAlreadyMember {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully joined group"})
}

func (h *Handler) GetGroupByID(c *gin.Context) {
	id := c.Param("id")
	group, err := h.getGroupByIDUC.Execute(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	c.JSON(http.StatusOK, group)
}
