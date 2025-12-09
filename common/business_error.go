package common

import "net/http"

type BizError struct {
	Code 		int 
	Message 	string
	HTTPStatus 	int
	Details 	map[string]interface{}
}

func (e *BizError) Error() string {
	return e.Message
}

func NewBizError(code int, message string, httpstatus int) *BizError {
	return &BizError{
		Code: code,
		Message: message,
		HTTPStatus: httpstatus,
		Details: make(map[string]interface{}),
	}
}

func (e *BizError) WithDetails(details map[string]interface{}) *BizError{
	e.Details = details
	return e
}

// ========== 用户模块错误（10xxx）==========

var (
	ErrUsernameExist = NewBizError(10001, "用户名已存在", http.StatusConflict)
	ErrEmailExist    = NewBizError(10002, "邮箱已被注册", http.StatusConflict)
	ErrInvalidAuth   = NewBizError(10003, "用户名或密码错误", http.StatusUnauthorized)
	ErrInvalidToken  = NewBizError(10004, "Token无效或已过期", http.StatusUnauthorized)
	ErrPermissionDenied = NewBizError(10005, "无权限访问", http.StatusForbidden)
	ErrUserDisabled  = NewBizError(10006, "用户已被禁用", http.StatusForbidden)
)

// ========== 图书模块错误（20xxx）==========

var (
	ErrBookNotFound   = NewBizError(20001, "图书不存在", http.StatusNotFound)
	ErrISBNExist      = NewBizError(20002, "ISBN已存在", http.StatusConflict)
	ErrBookOutOfStock = NewBizError(20003, "图书库存不足", http.StatusBadRequest)
	ErrCategoryNotFound = NewBizError(20004, "分类不存在", http.StatusNotFound)
)

// ========== 借阅模块错误（30xxx）==========

var (
	ErrBorrowNotFound    = NewBizError(30001, "借阅记录不存在", http.StatusNotFound)
	ErrBorrowLimitReached = NewBizError(30002, "借阅数量已达上限", http.StatusBadRequest)
	ErrHasOverdueBooks   = NewBizError(30003, "有逾期图书无法借阅", http.StatusBadRequest)
	ErrRenewLimitReached = NewBizError(30004, "续借次数已达上限", http.StatusBadRequest)
	ErrCannotRenewOverdue = NewBizError(30005, "逾期图书无法续借", http.StatusBadRequest)
	ErrBookAlreadyBorrowed = NewBizError(30006, "该图书已被借出", http.StatusBadRequest)
	ErrReservationFailed  = NewBizError(30007, "预约失败，图书有库存", http.StatusBadRequest)
)

// ========== 通用错误 ==========

var (
	ErrBadRequest      = NewBizError(400, "请求参数错误", http.StatusBadRequest)
	ErrUnauthorized    = NewBizError(401, "未认证", http.StatusUnauthorized)
	ErrForbidden       = NewBizError(403, "权限不足", http.StatusForbidden)
	ErrNotFound        = NewBizError(404, "资源不存在", http.StatusNotFound)
	ErrConflict        = NewBizError(409, "资源冲突", http.StatusConflict)
	ErrInternalServer  = NewBizError(500, "服务器内部错误", http.StatusInternalServerError)
)