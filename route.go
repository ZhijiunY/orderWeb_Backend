package main

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) wrap(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		cxt := context.WithValue(r.Context(), "params", p)
		next.ServeHTTP(w, r.WithContext(cxt))
	}
}

func (app *application) routes() http.Handler {
	router := httprouter.New()
	secure := alice.New(app.checkToken)
	router.HandlerFunc(http.MethodPost, "/v1/signin", app.Signin)

	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)
	// router.HandlerFunc(http.MethodPost, "/get-id-menu", app.getOneMenu)
	router.HandlerFunc(http.MethodGet, "/test-ok", app.testOK)

	router.HandlerFunc(http.MethodPost, "/get-all-menu", app.getAllMenu)
	router.HandlerFunc(http.MethodPost, "/get-opened-menu", app.getOpenedMenu)
	router.HandlerFunc(http.MethodPost, "/get-all-order", app.getAllOrder)
	router.HandlerFunc(http.MethodPost, "/get-order", app.getOrderById)

	router.POST("/create", app.wrap(secure.ThenFunc(app.createMenu)))
	router.POST("/update-open", app.wrap(secure.ThenFunc(app.updateOpen)))
	router.POST("/add-order", app.wrap(secure.ThenFunc(app.addOrder)))
	router.POST("/update-order", app.wrap(secure.ThenFunc(app.updateOrder)))
	router.POST("/get-menu", app.wrap(secure.ThenFunc(app.getOneMenu)))
	router.POST("/update-menu", app.wrap(secure.ThenFunc(app.updateMenu)))
	router.POST("/delete-open-menu", app.wrap(secure.ThenFunc(app.deleteOpenMenu)))
	router.POST("/delete-order", app.wrap(secure.ThenFunc(app.deleteOrder)))
	router.POST("/update-rating", app.wrap(secure.ThenFunc(app.updateMenuRating)))

	return app.enableCORS(router)

}
