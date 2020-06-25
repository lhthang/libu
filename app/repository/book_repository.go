package repository

import (
	"libu/app/form"
	"libu/app/model"
	"libu/my_db"
	"libu/utils/arrays"
	"libu/utils/firebase"
	"net/http"
	"sync"

	"github.com/araddon/dateparse"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var BookEntity IBook

type bookEntity struct {
	resource *my_db.Resource
	repo     *mongo.Collection
}

type BookChan struct {
	Ids []string
	Err error
}

type IBook interface {
	GetSimilarBooks(id string) ([]form.BookResponse, int, error)
	GetAll() ([]form.BookResponse, int, error)
	Search(keyword string) ([]form.BookResponse, int, error)
	Create(bookForm form.BookForm) (form.BookResponse, int, error)
	GetOneByID(id string) (form.BookResponse, int, error)
	Delete(id string) (form.BookResponse, int, error)
	Update(id string, bookForm form.UpdateBookForm) (form.BookResponse, int, error)
}

func NewBookEntity(resource *my_db.Resource) IBook {
	bookRepo := resource.DB.Collection("book")
	BookEntity = &bookEntity{
		resource: resource,
		repo:     bookRepo,
	}
	return BookEntity
}

func getCategoryOfBook(book *model.Book) []model.Category {
	var categories []model.Category
	var validCategoryIds []string
	for _, id := range book.CategoryIds {
		category, _, err := CategoryEntity.GetOneByID(id)
		if err != nil  {
			continue
		}
		validCategoryIds = append(validCategoryIds, id)
		categories = append(categories, *category)
	}
	book.CategoryIds = validCategoryIds
	return categories
}

func getAuthorsOfBook(book *model.Book) []model.Author {
	var authors []model.Author
	var validAuthorIds []string
	for _, id := range book.AuthorIds {
		author, _, err := AuthorEntity.GetOneByID(id)
		if err != nil {
			continue
		}
		validAuthorIds = append(validAuthorIds, id)
		authors = append(authors, *author)
	}
	book.AuthorIds = validAuthorIds
	return authors
}

func getReviewsOfBook(book model.Book) *form.ReviewResponse {
	reviewResp, _, err := ReviewEntity.GetByBookId(book.Id.Hex())
	if err != nil {
		return nil
	}
	return reviewResp
}

func (entity bookEntity) GetAll() ([]form.BookResponse, int, error) {
	ctx, cancel := initContext()

	defer cancel()
	var booksResp []form.BookResponse
	cursor, err := entity.repo.Find(ctx, bson.M{})
	if err != nil {
		return booksResp, getHTTPCode(err), err
	}

	for cursor.Next(ctx) {
		var book model.Book
		err := cursor.Decode(&book)
		if err != nil {
			logrus.Print(err)
		}
		reviewResp := getReviewsOfBook(book)
		booksResp = append(booksResp, form.BookResponse{
			Book: &book,
			//Reviews:    reviewResp.Reviews,
			Rating:     reviewResp.AvgRating,
			Categories: getCategoryOfBook(&book),
			Authors:    getAuthorsOfBook(&book),
		})
	}
	return booksResp, http.StatusOK, nil
}

func (entity bookEntity) GetSimilarBooks(id string) ([]form.BookResponse, int, error) {
	ctx, cancel := initContext()

	defer cancel()


	var booksResp []form.BookResponse

	book, _, err := entity.GetOneByID(id)
	if err != nil {
		return []form.BookResponse{}, http.StatusBadRequest, err
	}

	pipelineCategories :=[]bson.M{{"$unwind": bson.M{"path":"$categoryIds"}},
		{"$match":bson.M{"categoryIds":bson.M{"$in":book.CategoryIds}}},
		{"$group":bson.M{"_id":"$_id","total":bson.M{"$sum":1}}},
		{"$sort": bson.M{"total": -1}},
		}

	pipelineAuthors :=[]bson.M{{"$unwind": bson.M{"path":"$authorIds"}},
		{"$match":bson.M{"authorIds":bson.M{"$in":book.AuthorIds}}},
		{"$group":bson.M{"_id":"$_id","total":bson.M{"$sum":1}}},
		{"$sort": bson.M{"total": -1}},
	}

	authors :=make(chan []string)
	categories :=make(chan []string)

	go func() {
		var similarCategories []string
		cursor, err := entity.repo.Aggregate(ctx, pipelineCategories)
		if err != nil {
			logrus.Println(err)
		}
		for cursor.Next(ctx) {
			//logrus.Printf("%v",cursor)
			var book form.SimilarBook
			err := cursor.Decode(&book)
			if err != nil {
				logrus.Println(err)
				continue
			}
			//logrus.Printf("%v",book)
			similarCategories = append(similarCategories,book.Id.Hex())
		}
		categories<-similarCategories
	}()


	go func() {
		var similarAuthors []string
		cursor, err := entity.repo.Aggregate(ctx, pipelineAuthors)
		if err != nil {
			logrus.Println(err)
		}
		for cursor.Next(ctx) {
			//logrus.Printf("%v",cursor)
			var book form.SimilarBook
			err := cursor.Decode(&book)
			if err != nil {
				logrus.Println(err)
				continue
			}
			//logrus.Printf("%v",book)
			similarAuthors = append(similarAuthors,book.Id.Hex())
		}
		authors<-similarAuthors
	}()
	similarAuthors :=<-authors
	similarCategories :=<-categories

	similarBookIds :=arrays.Union(similarAuthors,similarCategories)

	var wg sync.WaitGroup
	bookResp :=make(chan *form.BookResponse)
	for i,_:=range similarBookIds{
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			chann:=new(form.BookResponse)
			book,_,err:=entity.GetOneByID(similarBookIds[i])
			if err!=nil{
				logrus.Println(err)
				return
			}
			chann=&book
			bookResp<-chann
		}(i)
	}

	go func() {
		wg.Wait()
		close(bookResp)
	}()
	for book :=range bookResp{
		booksResp = append(booksResp,*book)
	}
	return booksResp, http.StatusOK, nil
}

func (entity bookEntity) Create(bookForm form.BookForm) (form.BookResponse, int, error) {
	ctx, cancel := initContext()

	defer cancel()

	releaseAt, err := dateparse.ParseAny(bookForm.ReleaseAt)
	if err != nil {
		logrus.Println(err)
	}
	var categoryIds []string
	for _, id := range bookForm.CategoryIds {
		category, _, _ := CategoryEntity.GetOneByID(id)
		if category != nil {
			categoryIds = append(categoryIds, id)
		}
	}
	var authorIds []string
	for _, id := range bookForm.AuthorIds {
		author, _, _ := AuthorEntity.GetOneByID(id)
		if author != nil {
			authorIds = append(authorIds, id)
		}
	}
	book := model.Book{
		Id:        primitive.NewObjectID(),
		ReleaseAt: releaseAt,
		Title:     bookForm.Title,
		// Authors:     bookForm.Authors,
		AuthorIds:   authorIds,
		Publisher:   bookForm.Publisher,
		CategoryIds: categoryIds,
		Image:       bookForm.Image,
		Description: bookForm.Description,
		Link:        "",
	}
	channel := make(chan string)
	if bookForm.File != nil {
		go func() {

			url, _, err := firebase.UploadFile(*bookForm.File)
			if err != nil {
				logrus.Println(err)
			}
			channel <- url
		}()
		book.Link = <-channel
	}
	_, err = entity.repo.InsertOne(ctx, book)
	if err != nil {
		return form.BookResponse{}, getHTTPCode(err), err
	}
	bookResp := form.BookResponse{
		Book:       &book,
		Reviews:    nil,
		Categories: getCategoryOfBook(&book),
		Authors:    getAuthorsOfBook(&book),
	}
	return bookResp, http.StatusOK, nil
}

func (entity bookEntity) GetOneByID(id string) (form.BookResponse, int, error) {
	ctx, cancel := initContext()

	defer cancel()
	var book model.Book
	objID, _ := primitive.ObjectIDFromHex(id)

	err := entity.repo.FindOne(ctx, bson.M{"_id": objID}).Decode(&book)

	if err != nil {
		return form.BookResponse{}, getHTTPCode(err), err
	}

	reviewResp := getReviewsOfBook(book)

	bookResp := form.BookResponse{
		Book:       &book,
		Reviews:    reviewResp.ReviewResp,
		Rating:     reviewResp.AvgRating,
		Categories: getCategoryOfBook(&book),
		Authors:    getAuthorsOfBook(&book),
	}
	return bookResp, http.StatusOK, nil
}

func (entity bookEntity) Search(keyword string) ([]form.BookResponse, int, error) {
	ctx, cancel := initContext()

	defer cancel()
	var booksResp []form.BookResponse

	title := map[string]interface{}{}
	description := map[string]interface{}{}
	// author := map[string]interface{}{}
	publisher := map[string]interface{}{}
	title["title"] = bson.M{"$regex": keyword, "$options": "i"}
	description["description"] = bson.M{"$regex": keyword, "$options": "i"}
	// author["authors"] = bson.M{"$regex": keyword, "$options": "i"}
	publisher["publisher"] = bson.M{"$regex": keyword, "$options": "i"}
	query := map[string]interface{}{}
	// query["$or"] = []interface{}{title, description, author, publisher}
	query["$or"] = []interface{}{title, description, publisher}

	cursor, err := entity.repo.Find(ctx, query)
	if err != nil {
		return booksResp, getHTTPCode(err), err
	}

	for cursor.Next(ctx) {
		var book model.Book
		err := cursor.Decode(&book)
		if err != nil {
			logrus.Print(err)
		}

		reviewResp := getReviewsOfBook(book)
		booksResp = append(booksResp, form.BookResponse{
			Book: &book,
			//Reviews:    reviewResp.Reviews,
			Rating:     reviewResp.AvgRating,
			Categories: getCategoryOfBook(&book),
			Authors:    getAuthorsOfBook(&book),
		})
	}
	return booksResp, http.StatusOK, nil
}

func (entity bookEntity) Delete(id string) (form.BookResponse, int, error) {
	ctx, cancel := initContext()

	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(id)
	bookResp, _, err := entity.GetOneByID(id)

	if err != nil || bookResp.Book == nil {
		return form.BookResponse{}, getHTTPCode(err), err
	}

	_, err = entity.repo.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return form.BookResponse{}, getHTTPCode(err), err
	}
	return bookResp, http.StatusOK, nil
}

func (entity bookEntity) Update(id string, bookForm form.UpdateBookForm) (form.BookResponse, int, error) {
	ctx, cancel := initContext()

	defer cancel()
	objID, _ := primitive.ObjectIDFromHex(id)
	bookResp, _, err := entity.GetOneByID(id)
	if err != nil || bookResp.Book == nil {
		return form.BookResponse{}, getHTTPCode(err), err
	}

	var categoryIds []string
	for _, id := range bookForm.CategoryIds {
		category, _, _ := CategoryEntity.GetOneByID(id)
		if category != nil {
			categoryIds = append(categoryIds, id)
		}
	}
	bookForm.CategoryIds = categoryIds
	var authorIds []string
	for _, id := range bookForm.AuthorIds {
		author, _, _ := AuthorEntity.GetOneByID(id)
		if author != nil {
			authorIds = append(authorIds, id)
		}
	}
	bookForm.AuthorIds = authorIds
	err = copier.Copy(bookResp.Book, bookForm)
	if err != nil {
		return form.BookResponse{}, getHTTPCode(err), err
	}

	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	var updatedBook model.Book

	err = entity.repo.FindOneAndUpdate(ctx, bson.M{"_id": objID}, bson.M{"$set": bookResp.Book}, opts).Decode(&updatedBook)
	newBookResp := form.BookResponse{
		Book:       &updatedBook,
		Reviews:    bookResp.Reviews,
		Categories: getCategoryOfBook(&updatedBook),
		Authors:    getAuthorsOfBook(&updatedBook),
		Rating:     bookResp.Rating,
	}
	return newBookResp, http.StatusOK, nil
}
