package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mrizkyanditama/lelangin/api/auth"
	"github.com/mrizkyanditama/lelangin/api/fileupload"
	"github.com/mrizkyanditama/lelangin/api/models"
	"github.com/mrizkyanditama/lelangin/api/utils/formaterror"
)

func (server *Server) CreateAuction(c *gin.Context) {

	//clear previous error if any
	errList = map[string]string{}

	product_raw := c.PostForm("product")
	fmt.Println(product_raw)
	// if err != nil {
	// 	fmt.Println(err)
	// 	errList["Invalid_body"] = "Unable to get request"
	// 	c.JSON(http.StatusUnprocessableEntity, gin.H{
	// 		"status": http.StatusUnprocessableEntity,
	// 		"error":  errList,
	// 	})
	// 	return
	// }
	product := models.Product{}

	err := json.Unmarshal([]byte(product_raw), &product)
	if err != nil {
		fmt.Println(err)
		errList["Unmarshal_error"] = "Cannot unmarshal product"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	auction_raw := c.PostForm("auction")
	fmt.Println(auction_raw)

	auction := models.Auction{}

	err = json.Unmarshal([]byte(auction_raw), &auction)
	if err != nil {
		fmt.Println(err)
		errList["Unmarshal_error"] = "Cannot unmarshal product"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	product.Auction = auction

	uid, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	// check if the user exist:
	user := models.User{}
	err = server.DB.Debug().Model(models.User{}).Where("id = ?", uid).Take(&user).Error
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	product.OwnerID = uid //the authenticated user is the one creating the product

	file, err := c.FormFile("product_pic")
	if err != nil {
		errList["Invalid_file"] = "Invalid File"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	uploadedFile, fileErr := fileupload.FileUpload.UploadFile(file)
	if fileErr != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  fileErr,
		})
		return
	}

	product.PhotoPath = uploadedFile

	cid := product.CategoryID
	category := models.Category{}
	err = server.DB.Debug().Model(models.Category{}).Where("id = ?", cid).Take(&category).Error
	if err != nil {
		errList["Category_not_found"] = "Category not found"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	// tids := product.TagsIDs
	// tags := make([]*models.Tag, (len(tids)))
	// if len(tids) > 0 {
	// 	for i, _ := range tids {
	// 		err := server.DB.Debug().Model(models.Tag{}).Where("id = ?", tids[i]).Take(&tags[i]).Error
	// 		if err != nil {
	// 			errList["Tag_not_found"] = "Tag not found"
	// 			c.JSON(http.StatusUnprocessableEntity, gin.H{
	// 				"status": http.StatusUnprocessableEntity,
	// 				"error":  errList,
	// 			})
	// 			return
	// 		}
	// 	}
	// }

	// product.Tags = tags

	product.Prepare()
	errorMessages := product.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	productCreated, err := product.AddProduct(server.DB)
	if err != nil {
		errList := formaterror.FormatError(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":   http.StatusCreated,
		"response": productCreated,
	})
}

func (server *Server) GetProducts(c *gin.Context) {

	product := models.Product{}

	products, err := product.FindAllProducts(server.DB)
	if err != nil {
		errList["No_product"] = "No Product Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": products,
	})
}

func (server *Server) GetAuctions(c *gin.Context) {

	auction := models.Auction{}

	auctions, err := auction.FindAllAuction(server.DB)
	if err != nil {
		errList["No_product"] = "No Auction Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": auctions,
	})
}

func (server *Server) GetProduct(c *gin.Context) {

	productID := c.Param("id")
	pid, err := strconv.ParseUint(productID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	product := models.Product{}

	productReceived, err := product.FindProductByID(server.DB, pid)
	if err != nil {
		errList["No_product"] = "No Product Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": productReceived,
	})
}
