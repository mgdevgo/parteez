package handlers

import (
	"fmt"
	"log"
	"time"

	"iditusi/internal/repositories"
	"iditusi/internal/response"

	"github.com/gofiber/fiber/v2"
)

type EventController struct {
	eventRepository repositories.EventRepository
}

func NewEventController(eventRepository repositories.EventRepository) *EventController {
	return &EventController{eventRepository: eventRepository}
}

func (ec *EventController) Boot(routes fiber.Router) {
	events := routes.Group("/events")
	events.Get("/", ec.getEvents)
	events.Get("/:id", ec.getEvent)
}

func (ec *EventController) getEvents(ctx *fiber.Ctx) error {
	// limit := ctx.QueryInt("limit", 5)
	// offset := ctx.QueryInt("offset", 0)

	loc, _ := time.LoadLocation("Europe/Moscow")
	_ = time.Now().In(loc)

	from := ctx.Query("fromDate")
	to := ctx.Query("toDate")

	if from == "" || to == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.Error{
			Status: fiber.StatusBadRequest,
			Code:   "INVALID_PARAMETER",
			Title:  "URL Parameter 'fromDate' or 'toDate' is missing.",
		})
	}

	fromDate, err := time.Parse(time.DateOnly, from)
	if err != nil {
		fmt.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(response.Error{
			Status: fiber.StatusBadRequest,
			Code:   "INVALID_PARAMETER",
			Title:  "URL Parameter 'fromDate' is invalid.",
			Detail: "URL Parameter 'fromDate' must be date only format.",
		})
	}
	toDate, err := time.Parse(time.DateOnly, to)
	if err != nil {
		fmt.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(response.Error{
			Status: fiber.StatusBadRequest,
			Code:   "INVALID_PARAMETER",
			Title:  "URL Parameter 'toDate' is invalid.",
			Detail: "URL Parameter 'toDate' must implement RFC3339 format.",
		})
	}

	if fromDate.Unix() > toDate.Unix() {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.Error{
			Status: fiber.StatusBadRequest,
			Code:   "INVALID_PARAMETER",
			Title:  "URL Parameter 'toDate' is after 'fromDate'.",
		})
	}

	log.Println(fromDate, toDate)

	events, err := ec.eventRepository.FindByDate(ctx.Context(), fromDate, toDate)
	if err != nil {
		log.Println(err)
		return ctx.Status(500).JSON(response.Error{
			Status: fiber.StatusInternalServerError,
			Code:   "INTERNAL_SERVER_ERROR",
			Title:  "Internal Server Error",
		})
	}
	return ctx.JSON(events)
}

func (ec *EventController) getEvent(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id", 0)
	if err != nil {
		return err
	}

	// if id == 0 {
	// 	return ValidationError
	// }

	event, err := ec.eventRepository.FindById(ctx.Context(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(event)
}

//
// {
// "offers": [
// {
// "isSubscriptionShareable": false,
// "isOfferable": true,
// "price": "$10.99",
// "buyParams": {
// "offerName": "com.apple.music.AppleMusicMonthly",
// "offrd-free-trial": "true",
// "price": "0",
// "pg": "default",
// "appExtVrsId": "862022487",
// "salableAdamId": "1608220330",
// "pricingParameters": "STDQ",
// "bid": "com.apple.Music",
// "offrd-intro-price": "true",
// "productType": "A",
// "appAdamId": "1108187390"
// },
// "adamId": 1608220330,
// "freeTrialPeriod": null,
// "isBundle": false,
// "renewalPeriod": "P1M",
// "introOfferPrice": "$0.00",
// "introPeriod": "P1M",
// "subscriptionType": "Individual",
// "clientIntegrationType": null,
// "capacityInBytes": null,
// "isIntroOffer": true,
// "eligibilityType": "INTRO"
// },
// {
// "isSubscriptionShareable": false,
// "isOfferable": true,
// "price": "$5.99",
// "buyParams": {
// "offerName": "com.apple.music.AppleMusicStudent",
// "offrd-free-trial": "true",
// "price": "0",
// "pg": "default",
// "appExtVrsId": "862022487",
// "salableAdamId": "1608221191",
// "pricingParameters": "STDQ",
// "bid": "com.apple.Music",
// "offrd-intro-price": "true",
// "productType": "A",
// "appAdamId": "1108187390"
// },
// "adamId": 1608221191,
// "freeTrialPeriod": null,
// "isBundle": false,
// "renewalPeriod": "P1M",
// "introOfferPrice": "$0.00",
// "introPeriod": "P1M",
// "subscriptionType": "Student",
// "clientIntegrationType": null,
// "capacityInBytes": null,
// "isIntroOffer": true,
// "eligibilityType": "INTRO"
// },
// {
// "isSubscriptionShareable": true,
// "isOfferable": true,
// "price": "$16.99",
// "buyParams": {
// "offerName": "com.apple.music.AppleMusicFamily",
// "offrd-free-trial": "true",
// "price": "0",
// "pg": "default",
// "appExtVrsId": "862022487",
// "salableAdamId": "1608221109",
// "pricingParameters": "STDQ",
// "bid": "com.apple.Music",
// "offrd-intro-price": "true",
// "productType": "A",
// "appAdamId": "1108187390"
// },
// "adamId": 1608221109,
// "freeTrialPeriod": null,
// "isBundle": false,
// "renewalPeriod": "P1M",
// "introOfferPrice": "$0.00",
// "introPeriod": "P1M",
// "subscriptionType": "Family",
// "clientIntegrationType": null,
// "capacityInBytes": null,
// "isIntroOffer": true,
// "eligibilityType": "INTRO"
// }
// ],
// "status": 0
// }
