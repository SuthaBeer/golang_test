package user

import (
	"golang-apiuser/orm"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllUsers(c *gin.Context) {
	var users []orm.User
	orm.Db.Find(&users)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "success",
		"users":   users,
	})
}

// UserInfo			godoc
// @Summary 		UserInfo
// @Description 	User info API
// @Produce			application/json
// @Tags			User API
// @Router 			/users/userinfo [get]
// @Security BearerAuth
func GetUserInfo(c *gin.Context) {
	userId := c.MustGet("userId").(float64)
	if userId <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "User ID not found",
		})
		return
	}
	var user orm.User
	orm.Db.First(&user, userId)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "success",
		"user":    user,
	})
}

type TransferModel struct {
	ToAccountNo string
	Credit      int
}

// TransferCredit			godoc
// @Summary 		TransferCredit
// @Param			parameter body TransferModel true "TransferModel"
// @Description 	Transfer credit API
// @Produce			application/json
// @Tags			Transfer Credit
// @Router 			/users/transfercredit [post]
// @Security BearerAuth
func TransferCredit(c *gin.Context) {
	var post TransferModel
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId := c.MustGet("userId").(float64)
	var fromUser orm.User
	orm.Db.First(&fromUser, userId)
	if fromUser.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "User not found",
		})
		return
	}

	var toUser orm.User
	orm.Db.Where("account_no = ?", post.ToAccountNo).First(&toUser)
	if toUser.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "To User not found",
		})
		return
	}

	if fromUser.ID == toUser.ID {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "Can not transfer to yourself",
		})
		return
	}

	transferCredit := post.Credit

	if transferCredit <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "Transfer Credit must greater than 0",
		})
		return
	}

	userCurrentCredit := fromUser.Credit
	if userCurrentCredit < transferCredit {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "Credit not enough",
		})
		return
	}

	transaction := orm.Transaction{
		UserID:   fromUser.ID,
		ToUserID: toUser.ID,
		Credit:   transferCredit,
	}
	orm.Db.Create(&transaction)

	fromUser.Credit = userCurrentCredit - transferCredit
	orm.Db.Save(fromUser)

	toUserCurrentCredit := toUser.Credit
	toUser.Credit = toUserCurrentCredit + transferCredit
	orm.Db.Save(toUser)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Transfer Success",
	})
}

// TransferCreditHistory			godoc
// @Summary 		TransferCreditHistory
// @Description 	Transfer credit API
// @Produce			application/json
// @Tags			Transfer Credit
// @Router 			/users/transfercredithistory [get]
// @Security BearerAuth
func TransferCreditHistory(c *gin.Context) {
	userId := c.MustGet("userId").(float64)
	var user orm.User
	orm.Db.First(&user, userId)
	if user.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "User not found",
		})
		return
	}

	var transaction []orm.Transaction
	orm.Db.Where("user_id = ?", user.ID).Or("to_user_id = ?", user.ID).Find(&transaction)

	var transaction_in []orm.Transaction
	var transaction_out []orm.Transaction
	for _, element := range transaction {
		if element.ToUserID == user.ID {
			transaction_in = append(transaction_in, element)
		} else {
			transaction_out = append(transaction_out, element)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":          "success",
		"user_id":         user.ID,
		"transaction_in":  transaction_in,
		"transaction_out": transaction_out,
	})
}
