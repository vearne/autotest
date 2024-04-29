package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful/v3"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books = []Book{
	{ID: 1, Title: "The Go Programming Language", Author: "Alan A. A. Donovan and Brian W. Kernighan"},
	{ID: 2, Title: "Effective Go", Author: "The Go Authors"},
}

func getBooksHandler(req *restful.Request, res *restful.Response) {
	res.WriteAsJson(books)
}

func getBookHandler(req *restful.Request, res *restful.Response) {
	id := req.PathParameter("id")
	for _, book := range books {
		if strconv.Itoa(book.ID) == id {
			res.WriteAsJson(book)
			return
		}
	}
	res.WriteErrorString(http.StatusNotFound, "Book not found")
}

func createBookHandler(req *restful.Request, res *restful.Response) {
	book := new(Book)
	err := req.ReadEntity(book)
	if err != nil {
		res.WriteError(http.StatusBadRequest, err)
		return
	}
	book.ID = len(books) + 1
	books = append(books, *book)
	res.WriteAsJson(book)
}

func updateBookHandler(req *restful.Request, res *restful.Response) {
	id := req.PathParameter("id")
	for i, book := range books {
		if strconv.Itoa(book.ID) == id {
			updatedBook := new(Book)
			err := req.ReadEntity(updatedBook)
			if err != nil {
				res.WriteError(http.StatusBadRequest, err)
				return
			}
			updatedBook.ID = book.ID
			books[i] = *updatedBook
			res.WriteAsJson(updatedBook)
			return
		}
	}
	res.WriteErrorString(http.StatusNotFound, "Book not found")
}

func deleteBookHandler(req *restful.Request, res *restful.Response) {
	id := req.PathParameter("id")
	for i, book := range books {
		if strconv.Itoa(book.ID) == id {
			books = append(books[:i], books[i+1:]...)
			res.WriteHeader(http.StatusNoContent)
			return
		}
	}
	res.WriteErrorString(http.StatusNotFound, "Book not found")
}

func main() {
	ws := new(restful.WebService)
	ws.Path("/api/books").
		Route(ws.GET("").To(getBooksHandler)).
		Route(ws.GET("/{id}").To(getBookHandler)).
		Route(ws.POST("").To(createBookHandler)).
		Route(ws.PUT("/{id}").To(updateBookHandler)).
		Route(ws.DELETE("/{id}").To(deleteBookHandler))
	restful.Add(ws)

	log.Println("starting...")
	http.ListenAndServe(":8080", nil)
}
