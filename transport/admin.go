package transport

import (
	"net/http"
	"raise-child/business"
	action_type "raise-child/constants/action_type"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/util"

	"github.com/gin-gonic/gin"
)

// GetAdmins godoc
// @Summary      List admins
// @Description  Retrieves a list of admins based on filter criteria
// @Tags         admins
// @Accept       json
// @Produce      json
// @Param        request  query     request.GetAdminsRequest  true  "Filter Criteria"
// @Success      200      {object}  response.PaginationDataResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /admins [get]
func GetAdmins(ctx *gin.Context) {
	// var request request.GetAdminsRequest
	// if ctx.ShouldBindQuery(&request) != nil {
	// 	util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
	// 	return
	// }

	// service, err := business.GenerateAdminService()
	// if err != nil {
	// 	util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
	// 	return
	// }

	// res, err := service.GetAdmins(request, ctx)
	// util.ProcessResponse(response.APIResponse{
	// 	Data1:    res,
	// 	Data2:    res,
	// 	ErrMsg:   err,
	// 	Context:  ctx,
	// 	PostType: action_type.NON_POST,
	// })

	// ctx.Redirect(http.StatusSeeOther, "raisechild://payment/callback")

	html := `
<!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Đang chuyển hướng...</title>
    <style>
        body { font-family: sans-serif; text-align: center; padding-top: 50px; }
        .btn { 
            display: inline-block; padding: 12px 24px; margin-top: 20px;
            background-color: #007bff; color: white; text-decoration: none; 
            border-radius: 8px; font-weight: bold;
        }
    </style>
</head>
<body>
    <h3>Đang quay lại ứng dụng Raisechild...</h3>
    <p>Nếu ứng dụng không tự động mở, vui lòng nhấn vào nút bên dưới:</p>
    
    <a href="raisechild://payment/callback" class="btn">Mở ứng dụng</a>

    <script>
        // Chạy ngay lập tức
        window.location.href = "raisechild://payment/callback";

        // Thêm một lần thử lại sau 1 giây nếu lệnh trên bị trình duyệt chặn
        setTimeout(function() {
            window.location.href = "raisechild://payment/callback";
        }, 1000);
    </script>
</body>
</html>`

	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// UpdatePublisherInfo godoc
// @Summary      Update publisher information
// @Description  Updates the details of publisher (admin) based on the provided request body with just only 1 call permitted.
// @Tags         admins
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      request.UpdatePublisherInfoRequest  true  "Publisher Update Information"
// @Success      200      {object}  response.BuildTransactionResponse
// @Failure      400      {object}  response.MessageAPIResponse "Invalid data. Please try again."
// @Failure      401      {object}  response.MessageAPIResponse "You have no rights to access this action."
// @Failure      500      {object}  response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /admins [post]
func UpdatePublisherInfo(ctx *gin.Context) {
	var request request.UpdatePublisherInfoRequest
	if ctx.ShouldBindJSON(&request) != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business.GenerateAdminService()
	if err != nil {
		util.ProcessResponse(util.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.UpdatePublisherInfo(request, ctx)
	util.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}
