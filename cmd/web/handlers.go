package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"
)

func postLogin(c *gin.Context) {
	email, _ := c.GetPostForm("email")
	password, _ := c.GetPostForm("password")

	var login struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	login.Email = email
	login.Password = password

	jsonData, err := json.Marshal(login)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "problems with json login",
		})
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:4000/users/login", bytes.NewBuffer(jsonData))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "problems with req login",
		})
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "problems with resp login",
		})
		return
	}
	defer resp.Body.Close()

	var output struct {
		Token string `json:"token"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "problems with resp body login",
		})
		return
	}
	if err = json.Unmarshal(body, &output); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"Error decoding JSON:": err.Error(),
		})
		return
	}

	//fmt.Println(output.Token)
	//c.IndentedJSON(http.StatusCreated, output)

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", output.Token, 3600*24*30, "", "", false, true)

	c.Redirect(http.StatusFound, "/")
}

func postReg(c *gin.Context) {
	name, _ := c.GetPostForm("name")
	email, _ := c.GetPostForm("email")
	password, _ := c.GetPostForm("password")
	confirmPass, _ := c.GetPostForm("confirm-password")
	if password != confirmPass {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "passwords are not similar",
		})
		return
	}
	var register struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	register.Name = name
	register.Email = email
	register.Password = password

	jsonData, err := json.Marshal(register)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "problems with json",
		})
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:4000/users", bytes.NewBuffer(jsonData))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "problems with req",
		})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "problems with resp",
		})
		return
	}
	defer resp.Body.Close()
	var User struct {
		Id         int64     `json:"id"`
		Name       string    `json:"name"`
		Email      string    `json:"email"`
		Activated  bool      `json:"activated"`
		Created_At time.Time `json:"created_At"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "problems with resp body",
		})
		return
	}
	if err = json.Unmarshal(body, &User); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"Error decoding JSON:": err.Error(),
		})
		return
	}

	var login struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	login.Email = email
	login.Password = password

	jsonData, err = json.Marshal(login)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "problems with json login",
		})
		return
	}

	req, err = http.NewRequest("POST", "http://localhost:4000/users/login", bytes.NewBuffer(jsonData))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "problems with req login",
		})
		return
	}

	resp, err = client.Do(req)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "problems with resp login",
		})
		return
	}
	defer resp.Body.Close()

	var output struct {
		Token string `json:"token"`
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "problems with resp body login",
		})
		return
	}
	if err = json.Unmarshal(body, &output); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"Error decoding JSON:": err.Error(),
		})
		return
	}

	//fmt.Println(output.Token)
	//c.IndentedJSON(http.StatusCreated, output)

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", output.Token, 3600*24*30, "", "", false, true)

	c.Redirect(http.StatusFound, "/")
}

func getLogin(c *gin.Context) {
	_, err := c.Cookie("Authorization")
	if err == nil {
		c.Redirect(http.StatusFound, "/")
	}
	c.HTML(200, "LoginPage.html", gin.H{})
}

func getReg(c *gin.Context) {
	_, err := c.Cookie("Authorization")
	if err == nil {
		c.Redirect(http.StatusFound, "/")
	}
	c.HTML(200, "RegistrationPage.html", gin.H{})
}

func getProd(id int64) (*Products, error) {
	var product Products
	res, err := http.Get(fmt.Sprintf("http://localhost:4002/data-collection/%d", id))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, &product); err != nil {
		return nil, err
	}
	if product.Error != nil {
		return nil, ErrInvalidId
	}
	return &product, nil
}

func getMain(c *gin.Context) {
	products := make([]*Products, 0, 3)
	res, err := http.Get("http://localhost:4003/search?limit=9")
	if err != nil {
		fmt.Println("here1")
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("here2")
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if err = json.Unmarshal(body, &products); err != nil {
		fmt.Println("here3")
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	name, _ := c.Get("name")
	fmt.Println(name)
	for _, product := range products {
		avg := 0
		minium := math.MaxInt32
		maxium := 0
		for i := range product.Prices {
			avg += product.Prices[i].Price
			if product.Prices[i].Price < minium {
				minium = product.Prices[i].Price
			}
			if product.Prices[i].Price > maxium {
				maxium = product.Prices[i].Price
			}
		}
		product.Reviews, err = getReviews(product.Product_id)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		product.MaxPrice = maxium
		product.MinPrice = minium
		product.AvgPrice = avg / len(product.Prices)
		product.ShopsCount = len(product.Prices)
		product.CountReviews = len(product.Reviews)
		if product.CountReviews > 0 {
			avg = 0
			for i := range product.Reviews {
				avg += int(product.Reviews[i].Rating)
			}
			product.AvgRating = avg / product.CountReviews
		} else {
			product.AvgRating = 0
		}
	}

	c.HTML(200, "MainPage.html", gin.H{
		"email":    name,
		"products": products,
	})
}

func getProduct(c *gin.Context) {
	params := c.Param("id")
	id, err := strconv.ParseInt(params, 10, 64)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	product, err := getProd(id)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidId):
			c.Redirect(http.StatusFound, "/")
		default:
			ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return

	}
	avg := 0
	minium := math.MaxInt32
	maxium := 0
	for i := range product.Prices {
		avg += product.Prices[i].Price
		if product.Prices[i].Price < minium {
			minium = product.Prices[i].Price
		}
		if product.Prices[i].Price > maxium {
			maxium = product.Prices[i].Price
		}
	}
	product.MaxPrice = maxium
	product.MinPrice = minium
	product.AvgPrice = avg / len(product.Prices)
	product.ShopsCount = len(product.Prices)

	product.Reviews, err = getReviews(product.Product_id)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	product.CountReviews = len(product.Reviews)
	if product.CountReviews > 0 {
		avg = 0
		for i := range product.Reviews {
			avg += int(product.Reviews[i].Rating)
		}
		product.AvgRating = avg / product.CountReviews
	} else {
		product.AvgRating = 0
	}

	name, _ := c.Get("name")
	fmt.Println(name)
	c.HTML(200, "Iphone11.html", gin.H{
		"product": product,
		"email":   name,
	})
}

func logout(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/")
}

func search(c *gin.Context) {
	searchQuery := c.Query("search")
	category := c.Query("category")
	brand := c.QueryArray("brand[]")
	ram := c.Query("ram")
	matrix := c.Query("matrix")
	numberOfCores := c.Query("numberOfCores")

	fmt.Println(brand)
	sortBy := c.Query("sortBy")
	priceFrom := 0
	priceTo := 1500000

	if c.Query("priceFrom") != "" {
		if newPriceFrom, err := strconv.Atoi(c.Query("priceFrom")); err == nil {
			priceFrom = newPriceFrom
		}
	}

	if c.Query("priceTo") != "" {
		if newPriceTo, err := strconv.Atoi(c.Query("priceTo")); err == nil {
			priceTo = newPriceTo
		}
	}

	base := fmt.Sprintf("priceFrom=%d&priceTo=%d&", priceFrom, priceTo)
	fmt.Println(brand)
	if searchQuery != "" {
		base += fmt.Sprintf("search=%s&", searchQuery)
	}
	if category != "" {
		base += fmt.Sprintf("category=%s&", category)
	}
	if len(brand) != 0 {
		for _, str := range brand {
			base += "brand%5B%5D=" + str + "&"
		}
	}
	if sortBy != "" {
		base += fmt.Sprintf("sortBy=%s&", sortBy)
	}
	url := "http://localhost:4003/search?" + base + "limit=15"
	fmt.Println(url)
	products := make([]*Products, 0, 3)
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("here1")
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("here2")
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if err = json.Unmarshal(body, &products); err != nil {
		fmt.Println("here3")
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	name, _ := c.Get("name")
	fmt.Println(name)

	for _, product := range products {
		avg := 0
		minium := math.MaxInt32
		maxium := 0
		for i := range product.Prices {
			avg += product.Prices[i].Price
			if product.Prices[i].Price < minium {
				minium = product.Prices[i].Price
			}
			if product.Prices[i].Price > maxium {
				maxium = product.Prices[i].Price
			}
		}
		product.MaxPrice = maxium
		product.MinPrice = minium
		product.AvgPrice = avg / len(product.Prices)
		product.ShopsCount = len(product.Prices)

		product.Reviews, err = getReviews(product.Product_id)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		product.CountReviews = len(product.Reviews)
		if product.CountReviews > 0 {
			avg = 0
			for i := range product.Reviews {
				avg += int(product.Reviews[i].Rating)
			}
			product.AvgRating = avg / product.CountReviews
		} else {
			product.AvgRating = 0
		}

	}
	var brands []string
	var smartphoneFilter SmartphoneFilter
	filters := make(map[string]string)
	if ram != "" {
		filters["RAM capacity"] = ram
	}
	if matrix != "" {
		filters["Matrix type (details)"] = matrix
	}
	if numberOfCores != "" {
		filters["Number of Cores"] = numberOfCores
	}

	if ram != "" || matrix != "" || numberOfCores != "" {
		productsToSend := make([]*Products, 0, 5)
		match := true
		for _, product := range products {
			for key, value := range filters {
				attributeMatch := false
				for _, attr := range product.Attributes {
					if attr.AttributeName == key && attr.AttributeValue == value {
						attributeMatch = true
						break
					}
				}
				if !attributeMatch {
					match = false
					break
				}
				match = true
			}
			if match {
				productsToSend = append(productsToSend, product)
			}
		}
		fmt.Println(productsToSend)
		products = productsToSend
	}
	if category == "smartphone" {
		brands = []string{"Apple", "Samsung", "Xiaomi"}
		smartphoneFilter.Ram = []string{"4 GB", "6 GB", "8 GB", "12 GB"}
		smartphoneFilter.Matrix = []string{"OLED", "Dynamic AMOLED 2X", "Super Retina XDR"}
		smartphoneFilter.NumberOfCores = []string{"2", "4", "6", "8"}
		smartphoneFilter.Exist = true
	} else if category == "headphones" {
		brands = []string{"Apple"}
	} else if category == "TV" {
		brands = []string{"Samsung"}
	} else if category == "tablet" {
		brands = []string{"Samsung"}
	} else if category == "laptop" {
		brands = []string{"Acer", "DEXP", "Lenovo", "ASUS", "HUAWEI", "ARDOR", "MSI"}
	} else {
		brands = []string{"Apple", "Samsung", "Xiaomi", "Sony", "GIGABYTE", "Realme", "Acer", "DEXP", "Lenovo", "ASUS", "HUAWEI", "ARDOR", "MSI"}
	}

	c.HTML(200, "Headphones.html", gin.H{
		"email":             name,
		"products":          products,
		"category":          category,
		"search":            searchQuery,
		"brands":            brands,
		"smartPhoneFilters": smartphoneFilter,
	})
}

func getReview(c *gin.Context) {
	params := c.Param("id")
	id, err := strconv.ParseInt(params, 10, 64)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	product, err := getProd(id)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidId):
			c.Redirect(http.StatusFound, "/")
		default:
			ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return

	}

	avg := 0
	minium := math.MaxInt32
	maxium := 0
	for i := range product.Prices {
		avg += product.Prices[i].Price
		if product.Prices[i].Price < minium {
			minium = product.Prices[i].Price
		}
		if product.Prices[i].Price > maxium {
			maxium = product.Prices[i].Price
		}
	}
	product.MaxPrice = maxium
	product.MinPrice = minium
	product.AvgPrice = avg / len(product.Prices)
	product.ShopsCount = len(product.Prices)

	name, _ := c.Get("name")
	fmt.Println(name)

	product.Reviews, err = getReviews(id)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	for i := range product.Reviews {
		product.Reviews[i].Username, err = getUser(product.Reviews[i].User_id)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		product.Reviews[i].ConvertTime = product.Reviews[i].Created_at.Format("January 2, 2006")
	}

	product.CountReviews = len(product.Reviews)
	if product.CountReviews > 0 {
		avg = 0
		for i := range product.Reviews {
			avg += int(product.Reviews[i].Rating)
		}
		product.AvgRating = avg / product.CountReviews
	} else {
		product.AvgRating = 0
	}

	c.HTML(200, "reviews.html", gin.H{
		"product": product,
		"email":   name,
	})
}

func postReview(c *gin.Context) {
	fmt.Println(c.Request.URL)
	params := c.Param("id")
	id, err := strconv.ParseInt(params, 10, 64)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	_, err = c.Cookie("Authorization")
	if err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/product/%d/review", id))
	}
	message, _ := c.GetPostForm("review-text")
	rating, _ := c.GetPostForm("rating")
	userId, _ := c.Get("userId")

	var input struct {
		Message    string `json:"message"`
		Rating     uint   `json:"rating"`
		User_id    int64  `json:"user_id"`
		Product_id int64  `json:"product_id"`
	}
	ratingInt, err := strconv.Atoi(rating)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	input.Message = message
	input.Rating = uint(ratingInt)
	input.User_id = userId.(int64)
	input.Product_id = id

	fmt.Println(input)

	jsonData, err := json.Marshal(input)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "problems with json login",
		})
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:4004/review", bytes.NewBuffer(jsonData))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "problems with req login",
		})
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "problems with resp login",
		})
		return
	}
	defer resp.Body.Close()

	c.Redirect(http.StatusFound, fmt.Sprintf("/product/%d/review", id))

}
