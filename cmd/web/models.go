package main

import (
	"errors"
	"time"
)

var (
	ErrInvalidId = errors.New("invalid id")
)

type Products struct {
	Product_id   int64        `json:"product_id"`
	ProductName  string       `json:"productName"`
	Category     string       `json:"category"`
	Brand        string       `json:"brand"`
	Description  string       `json:"description"`
	Attributes   []Attributes `json:"attributes"`
	Images       []Images     `json:"images"`
	Error        any          `json:"error"`
	Prices       []Prices     `json:"prices"`
	AvgPrice     int          `json:"avg_price"`
	MinPrice     int          `json:"min_price"`
	MaxPrice     int          `json:"max_price"`
	ShopsCount   int          `json:"shops_count"`
	Reviews      []Review     `json:"reviews"`
	CountReviews int          `json:"count_reviews"`
	AvgRating    int          `json:"avg_rating"`
}

type Attributes struct {
	AttributeID    int64  `json:"attributeID"`
	AttributeName  string `json:"attributeName"`
	AttributeValue string `json:"attributeValue"`
	Product_id     int64  `json:"product_id"`
}

type Images struct {
	Id         int64 `json:"id"`
	Main_image bool  `json:"main_image"`
	Product_id int64 `json:"product_id"`
}

type Shops struct {
	Id       int64  `json:"id"`
	ShopName string `json:"shopName"`
	Link     string `json:"link"`
}

type Prices struct {
	Id         int64  `json:"id"`
	Price      int    `json:"price"`
	Shop       Shops  `json:"shop"`
	Product_id int64  `json:"product_id"`
	Link       string `json:"link"`
}

type SmartphoneFilter struct {
	Matrix        []string `json:"matrix"`
	Ram           []string `json:"ram"`
	NumberOfCores []string `json:"number_of_cores"`
	Exist         bool     `json:"exist"`
}

type Review struct {
	Id          int64     `json:"id"`
	User_id     int64     `json:"user_id"`
	Message     string    `json:"message"`
	Rating      uint      `json:"rating"`
	Product_id  int64     `json:"product_id"`
	Created_at  time.Time `json:"created_at"`
	ConvertTime string    `json:"convert_time"`
	Username    string    `json:"username"`
}
