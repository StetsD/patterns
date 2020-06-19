package mediator

import (
	"testing"
)

func TestMediator(t *testing.T) {

	farmer := new(Farmer)
	cannery := new(Cannery)
	shop := new(Shop)

	farmer.AddMoney(7500.00)
	cannery.AddMoney(15000.00)
	shop.AddMoney(30000.00)

	СonnectСolleagues(farmer, cannery, shop)

	farmer.GrowTomato(1000)

	expect := float64(54750)
	result := shop.GetMoney()

	if result != expect {
		t.Errorf("Expect result to equal %f, but %f.\n", expect, result)
	}
}
