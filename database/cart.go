package database

import "errors"

//Error Variables

var (
	ErrorNotFound         = errors.New("product not found")
	ErrCantDecodeProduct  = errors.New("can't decode product")
	ErrUserIdNotValid     = errors.New("user Id is not valid")
	ErrCantUpdateUser     = errors.New("can't update user")
	ErrCantRemoveItemCart = errors.New("can't remove item from cart")
	ErrCantBuyCartItem    = errors.New("can't buy cart item")
	ErrCantGetItem        = errors.New("can't get item")
)

func AddProductToCart() {

}

func RemoveItemFromCart() {

}

func BuyItemFromCart() {

}

func InstantBuyer() {

}
