package routes

import (
	"strconv"
	"time"

	"os"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/sam-kash/Url_Shortner_Go.git/database"
	"golang.org/x/tools/go/analysis/passes/defers"
)

type request struct{
	URL			string			`json:"url"`
	CustomShort	string			`json:"short"`
	Expiry		time.Duration	`json:"expiry"`
}

type response struct{
	URL					string			`json:"url"`
	CustomShort			string			`json:"short"`
	Expiry				time.Duration	`json:"expiry"`
	XRateRemaining		int 			`json:"rate_limit"`
	XRateLimitReset		time.Duration	`json:"rate_limit_reset"`
}

func ShortenURL(c *fiber.Ctx) error{
	body:= new (request)
	
	if err := c.BodyParser(&body); err!=nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	// impliment rate limiting 

	/* 		Here we are going to check if the end user or the 
			IP address of that user is already addeed the databse
	*/

	r2 := database.CreateClient(1)
	defer r2.close()   // This Defer is executed after the function above has executed completely with its call stack
	val, err := r2.Get(database.Ctx, c.IP()).Result()

	if err != redis.Nil{
		_ =r2.Set(datbase.Ctx , c.IP, os.Getenv("API_QUOTA"), 30*60*time.second).Err()
	}else{
		r2.Get(database.Ctx, c.IP().Result())
		valInt, _ := strconv.Atoi(val)

		if valInt <=0 {
			limit, _ := r2.TTK(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error" : "Rate limit exceeded",
				"rate_limit_reset" : limit/time.Nanosecond/time.Minute,	
			})
		} 
	}

	// check if the input is an actual URL

	if !govalidator.IsURL(body.URL){
		retucn c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid_URL"})
	}

	//check the domain error

	if !helpers.RemoveDomainError(body.URL){
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error":"You cant hack the system"})
	}

	// enforce https or ssl

	body.URL = helpers.EnforceHTTP(body.URL)

	// Logic to build a custom URL , functionality
	// 1. check if the users has sent us any custom shorten URL or not
	// 2. Check there is already an entry somewhere in our db for that same custom URL or not 

	 var id string

	 if body.CustomShort == ""{
		id == uuid.New().String()[:6]
	 } else{
		id = body.CustomShort
	 }

	 r := database.CreateClient(0)
	 defer r.Close()

	 val, _ = r.Get(database.Ctx, id).Result()
	 if val != ""{
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error" : "URL Custom short is already used ",
		})
	 }

	if body.Expiry == 0{
		body.Expiry = 24
	}

	err = r.Set(database.Ctx, id , body.URL, body.Expiry*3600*time.Second* ).Err()

	if err != nill {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "Unable to connect to server"
		})
	}

	r2.Decr(database.Ctx , c.IP())
}