package controller

import (
    "library-system/common"
    "library-system/dto/request"
    "library-system/service"
	"strconv"

    "github.com/gin-gonic/gin"
)

type ReservationController struct {
    reservationService *service.ReservationService
}

func NewReservationController(service *service.ReservationService) *ReservationController {
    return &ReservationController{
        reservationService: service,
    }
}

// CreateReservation 创建预约
// POST /api/reservations
func (ctl *ReservationController) CreateReservation(c *gin.Context) {
    ctx := c.Request.Context()

    // 获取当前用户ID
    userID, _ := c.Get("user_id")

    var req request.CreateReservationRequest
    if err := common.ValidateStruct(c, &req); err != nil {
        c.Error(err)
        return
    }

    data, err := ctl.reservationService.CreateReservation(ctx, userID.(uint64), &req)
    if err != nil {
        c.Error(err)
        return
    }

    common.Success(c, 201, "预约成功", data)
}

// CancelReservation 取消预约
// DELETE /api/reservations/:id
func (ctl *ReservationController) CancelReservation(c *gin. Context) {
    ctx := c.Request.Context()

    // 获取当前用户ID
    userID, _ := c.Get("user_id")

    // 获取预约ID
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        c.Error(common.ErrBadRequest)
        return
    }

    if err := ctl.reservationService.CancelReservation(ctx, userID.(uint64), id); err != nil {
        c.Error(err)
        return
    }

    common.Success(c, 200, "预约已取消", gin.H{})
}

// GetMyReservations 获取我的预约列表
// GET /api/reservations/my
func (ctl *ReservationController) GetMyReservations(c *gin.Context) {
    ctx := c.Request. Context()

    // 获取当前用户ID
    userID, _ := c.Get("user_id")

    data, err := ctl.reservationService.GetMyReservations(ctx, userID.(uint64))
    if err != nil {
        c.Error(err)
        return
    }

    common. Success(c, 200, "success", data)
}
