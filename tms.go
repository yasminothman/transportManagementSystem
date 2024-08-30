package main

//declare packages
import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/docket", getDockets)
	router.POST("/docket", postDockets)
	router.POST("/logsheet", postLogsheet)
	router.GET("/docket/:orderNo", getDocketByID)
	router.GET("/logsheet/:logsheetNo", getLogsheetByNo)
	router.Run("localhost:8080")
}

// jSON file for order
type order struct {
	OrderNo       string  `json:"OrderNo"`
	Customer      string  `json:"Customer`
	PickUpPoint   string  `json:"PickUpPoint`
	DeliveryPoint string  `json:"DeliveryPoint`
	Quantity      int     `json:"Quantity`
	Volume        float32 `json:"Volume`
	Status        string  `json:"Status`
	TruckNo       string  `json:"TruckNo`
	LogsheetNo    string  `json:"LogsheetNo`
}

// jSON file for logsheet
type logsheet struct {
	LogsheetNo string   `json:"LogsheetNo"`
	Dockets    []string `json:"Dockets"`
	TruckNo    string   `json:"TruckNo"`
}

var docket = []order{}
var logsheetCounter int = 1

// get docket responds with the list of all docket as JSON.
func getDockets(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, docket)
}

// get docket by ID
func getDocketByID(c *gin.Context) {
	orderNo := c.Param("orderNo")

	// Loop over the list of docket, looking for
	// an docket whose ID value matches the parameter.
	for _, a := range docket {
		if a.OrderNo == orderNo {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "docket not found"})
}

// post docket into JSON
func postDockets(c *gin.Context) {
	var newOrder order

	if err := c.BindJSON(&newOrder); err != nil {
		return
	}
	newOrder.OrderNo = generateOrderNo()
	// Set the Status to "Created".
	newOrder.Status = "Created"
	// Add the new docket to the slice.
	docket = append(docket, newOrder)
	c.IndentedJSON(http.StatusCreated, newOrder)
}

var orderCounter int = 1

// generate order number TDN increment 1
func generateOrderNo() string {
	orderNo := fmt.Sprintf("TDN%04d", orderCounter)
	orderCounter++
	return orderNo
}

// post logsheet into JSON
func postLogsheet(c *gin.Context) {
	var newLogsheet logsheet

	if err := c.BindJSON(&newLogsheet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a LogsheetNo
	newLogsheet.LogsheetNo = generateLogsheetNo()

	// Slice to hold the updated dockets that will be returned in the response
	var updatedDockets []order

	// Update the dockets with the LogsheetNo and TruckNo
	for _, orderNo := range newLogsheet.Dockets {
		for j := range docket {
			if docket[j].OrderNo == orderNo {
				// Update the docket with the new LogsheetNo and TruckNo
				docket[j].LogsheetNo = newLogsheet.LogsheetNo
				docket[j].TruckNo = newLogsheet.TruckNo

				// Add the updated docket to the response slice
				updatedDockets = append(updatedDockets, docket[j])
				break
			}
		}
	}

	// Return the updated dockets in the response
	c.IndentedJSON(http.StatusCreated, updatedDockets)
}

// get logsheet by logsheet no
func getLogsheetByNo(c *gin.Context) {
	logsheetNo := c.Param("logsheetNo")

	var associatedDockets []order

	for _, a := range docket {
		if a.LogsheetNo == logsheetNo {
			associatedDockets = append(associatedDockets, a)
		}
	}

	if len(associatedDockets) > 0 {
		c.IndentedJSON(http.StatusOK, associatedDockets)
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "logsheet not found"})
	}
}

// generate logsheet no DT with increment 1
func generateLogsheetNo() string {
	logsheetNo := fmt.Sprintf("DT%04d", logsheetCounter)
	logsheetCounter++
	return logsheetNo
}
