package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lin-snow/ech0/internal/dto"
	"github.com/lin-snow/ech0/internal/models"
	"github.com/lin-snow/ech0/internal/services"
)

// 获取 Todo 列表
func GetTodos(c *gin.Context) {
	// 获取当前用户 ID
	userID := c.MustGet("userid").(uint)

	// 调用 Service 层获取 Todo 列表
	todos, err := services.GetTodos(userID)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.OK(todos, models.GetTodosSuccessMessage))
}

// 发布 Todo
func PostTodo(c *gin.Context) {
	var todo models.Todo

	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.InvalidRequestBodyMessage))
		return
	}

	// 获取当前用户 ID
	userID := c.MustGet("userid").(uint)
	todo.UserID = userID

	// 获取当前用户信息
	user, err := services.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.UserNotFoundMessage))
		return
	}

	// 检查是否为管理员
	if !user.IsAdmin {
		c.JSON(http.StatusOK, dto.Fail[string](models.NoPermissionMessage))
		return
	}

	// 设置用户名
	todo.Username = user.Username

	// 调用 Service 层发布 Todo
	if err := services.AddTodo(&todo); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.OK(todo, models.CreateTodoSuccessMessage))
}

// 更新 Todo (完成/未完成)
func UpdateTodo(c *gin.Context) {
	// 从 URL 参数获取留言 ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.InvalidIDMessage))
		return
	}

	// 获取当前用户 ID
	userID := c.MustGet("userid").(uint)
	// 调用 Service 层获取 Todo
	theTodo, err := services.GetTodoById(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	// 检查该 Todo 是否属于当前用户
	if theTodo.UserID != userID {
		c.JSON(http.StatusOK, dto.Fail[string](models.UpdateTodoFailMessage))
		return
	}

	// 调用 Service 层更新 Todo
	if err := services.UpdateTodo(theTodo); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.OK(theTodo, models.UpdateTodoSuccessMessage))
}

// 删除 Todo
func DeleteTodo(c *gin.Context) {
	// 从 URL 参数获取留言 ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](models.InvalidIDMessage))
		return
	}

	// 调用 Service 层获取 Todo
	theTodo, err := services.GetTodoById(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	// 获取当前用户 ID
	userID := c.MustGet("userid").(uint)

	// 检查该 Todo 是否属于当前用户
	if theTodo.UserID != userID {
		c.JSON(http.StatusOK, dto.Fail[string](models.DeleteTodoFailMessage))
		return
	}

	// 调用 Service 层删除 Todo
	if err := services.DeleteTodo(theTodo); err != nil {
		c.JSON(http.StatusOK, dto.Fail[string](err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.OK(theTodo, models.DeleteTodoSuccessMessage))
}
