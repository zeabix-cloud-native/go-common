package util

import (
	"fmt"

	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zeabix-cloud-native/go-common/common/config"
)

type Query struct {
	Page        int
	Limit       int
	SortField   string
	SortType    string
	FilterName  string
	FilterConds string
	FilterValue string
}

func RemoveDuplicateSlice[T string | int | float64](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func GeneratePassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+{}[]|\\:;<>,.?/~"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	password := make([]byte, length)
	for i := range password {
		password[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(password)
}

func GenerateErrorLinkInfo(cfg config.Config, id uint) string {

	return fmt.Sprintf("link : %s/%d", cfg.ErrorLink, id)
}

func ConvertStringToFloat64(str string) (s float64, err error) {
	floatValue, err := strconv.ParseFloat(str, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return 0, err
	}
	return floatValue, err
}

func ConvertStringToInt(str string) (s int, err error) {
	intValue, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return intValue, err
}

func ConvertStringToIntArray(str string, spliter string) []int {
	numStrings := strings.Split(str, spliter)
	var numbers []int
	for _, numStr := range numStrings {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			// Handle conversion error if needed
			fmt.Printf("Error converting %s to an integer: %v\n", numStr, err)
			continue
		}
		// Append the integer to the numbers array
		numbers = append(numbers, num)
	}
	return numbers
}

func RequestQueryCheck(c *gin.Context) (data *Query, err error) {
	val := new(Query)
	pageQ := c.Query("page")
	limitQ := c.Query("limit")
	sortFieldQ := c.Query("sort[0][field]")
	sortTypeQ := c.Query("sort[0][dir]")
	filterFieldQ := c.Query("filter[0][field]")
	filterCondsQ := c.Query("filter[0][type]")
	filterValQ := c.Query("filter[0][value]")
	if pageQ == "" {
		val.Page = 1
	} else {
		page, err := strconv.Atoi(pageQ)
		if err != nil {
			return data, err
		}
		val.Page = page
	}

	if limitQ == "" {
		val.Limit = 50
	} else {
		limitP, err := strconv.Atoi(limitQ)
		if err != nil {
			return data, err
		}
		val.Limit = limitP
	}

	if sortFieldQ == "" {
		val.SortField = ""
	} else {
		val.SortField = sortFieldQ
	}

	if sortTypeQ == "" {
		val.SortType = "ASC"
	} else {
		val.SortType = sortTypeQ
	}

	if filterFieldQ == "" {
		val.FilterName = ""
	} else {
		val.FilterName = filterFieldQ
	}

	if filterCondsQ == "" {
		val.FilterConds = ""
	} else {
		val.FilterConds = filterCondsQ
	}

	if filterValQ == "" {
		val.FilterValue = ""
	} else {
		val.FilterValue = filterValQ
	}

	return val, err
}
