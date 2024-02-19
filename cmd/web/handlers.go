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

	email, _ := c.Get("user")
	fmt.Println(email)
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
	}

	c.HTML(200, "MainPage.html", gin.H{
		"email":    email,
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

	email, _ := c.Get("user")
	fmt.Println(email)
	c.HTML(200, "Iphone11.html", gin.H{
		"product": product,
		"email":   email,
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

	email, _ := c.Get("user")
	fmt.Println(email)

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
	}

	c.HTML(200, "Headphones.html", gin.H{
		"email":    email,
		"products": products,
		"category": category,
		"search":   searchQuery,
	})
}
