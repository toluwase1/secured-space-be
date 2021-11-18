package server

import (
	"fmt"
	"github.com/decadevs/rentals-api/models"
	"github.com/decadevs/rentals-api/server/response"
	"github.com/decadevs/rentals-api/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

func GetInteriors(interiorIDs []string) []models.InteriorFeature {
	in := []models.InteriorFeature{}
	for _, id := range interiorIDs {
		in = append(in, models.InteriorFeature{
			ID: id,
		})
	}
	return in
}

func GetExteriors(exteriorIDs []string) []models.ExteriorFeature {
	ex := []models.ExteriorFeature{}
	for _, id := range exteriorIDs {
		ex = append(ex, models.ExteriorFeature{
			ID: id,
		})
	}
	return ex
}

func (s *Server) handleCreateApartment() gin.HandlerFunc {
	// function to handle adding an apartment
	return func(c *gin.Context) {
		//get the user id from a logged-in user
		userI, exists := c.Get("user")
		if !exists {
			log.Printf("can't get user from context\n")
			response.JSON(c, "", http.StatusUnauthorized, nil, []string{"you are not logged in"})
			return
		}
		userId := userI.(*models.User).ID

		form, err := c.MultipartForm()

		if err != nil {
			log.Printf("error parsing multipart form: %v", err)
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
			return
		}

		formImages := form.File["images"]
		images := []models.Images{}

		// upload the images to aws.
		for _, f := range formImages {
			file, err := f.Open()
			if err != nil {

			}
			fileExtension, ok := services.CheckSupportedFile(strings.ToLower(f.Filename))
			log.Printf(filepath.Ext(strings.ToLower(f.Filename)))
			fmt.Println(fileExtension)
			if ok {
				log.Println(fileExtension)
				response.JSON(c, "", http.StatusBadRequest, nil, []string{fileExtension + " image file type is not supported"})
				return
			}

			session, tempFileName, err := services.PreAWS(fileExtension, "apartment")
			if err != nil {
				log.Println("could not upload file", err)
			}

			url, err := s.DB.UploadFileToS3(session, file, tempFileName, f.Size)
			if err != nil {
				log.Println(err)
				response.JSON(c, "", http.StatusInternalServerError, nil, []string{"an error occured while uploading the image"})
				return
			}

			//_ = uploadFileToS3(nil, image, "name", 12)
			log.Printf("filename: %v", f.Filename)

			img := models.Images{
				URL: url,
			}
			images = append(images, img)
		}

		price, err := strconv.Atoi(c.PostForm("price"))
		if err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, []string{err.Error()})
			return
		}

		numOfRooms, err := strconv.Atoi(c.PostForm("no_of_rooms"))
		if err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, []string{err.Error()})
			return
		}

		furnished, err := strconv.ParseBool(c.PostForm("furnished"))
		if err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, []string{err.Error()})
			return
		}

		apartmentStatus, err := strconv.ParseBool(c.PostForm("apartment_status"))
		if err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, []string{err.Error()})
			return
		}

		exteriors := strings.Split(c.PostFormArray("exterior")[0], ",")
		interiors := strings.Split(c.PostFormArray("interior")[0], ",")

		apartment := models.Apartment{
			UserID:          userId,
			Title:           c.PostForm("title"),
			CategoryID:      c.PostForm("category"),
			Description:     c.PostForm("description"),
			Price:           price,
			NoOfRooms:       numOfRooms,
			Furnished:       furnished,
			Location:        c.PostForm("location"),
			ApartmentStatus: models.ApartmentStatus(apartmentStatus),
			Interiors:       GetInteriors(interiors),
			Exteriors:       GetExteriors(exteriors),
			Images:          images,
		}

		err = s.DB.CreateApartment(&apartment)
		if err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, []string{err.Error()})
			return
		}

		log.Println(apartment.Images)
		response.JSON(c, "Apartment Successfully Added", http.StatusOK, apartment, nil)

	}
}

func (s *Server) DeleteApartment() gin.HandlerFunc {
	return func(c *gin.Context) {
		if userI, exists := c.Get("user"); exists {
			if user, ok := userI.(*models.User); ok {
				apartmentID := c.Param("apartmentID")
				if apartmentID == "" {
					response.JSON(c, "", http.StatusBadRequest, nil, []string{"apartment id cannot be empty"})
					return
				}
				err := s.DB.DeleteApartment(apartmentID, user.ID)
				if err != nil {
					log.Println(err.Error())
					response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
					return
				}
				response.JSON(c, "Deleted Successfully", http.StatusOK, nil, nil)
				return
			}
		}
		log.Printf("can't get user from context\n")
		response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
	}
}
func (s *Server) handleUpdateApartmentDetails() gin.HandlerFunc {
		// function to handle adding an apartment
		return func(c *gin.Context) {
		//get the user id from a logged-in user
		userI, exists := c.Get("user")
		if !exists {
		log.Printf("can't get user from context\n")
		response.JSON(c, "", http.StatusUnauthorized, nil, []string{"you are not logged in"})
		return
	}
		userId := userI.(*models.User).ID
		//if err := s.decode(c, &apartmentRequest); err != nil {
		//	response.JSON(c, "", http.StatusBadRequest, nil, err)
		//	return
		//}
		form, err := c.MultipartForm()
		if err != nil {
		log.Printf("error parsing multipart form: %v", err)
		response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
		return
	}

		formImages := form.File["images"]
		images := []models.Images{}
		for _, f := range formImages {
		//_ = uploadFileToS3(nil, image, "name", 12)
		log.Printf("filename: %v", f.Filename)
		url := "https://unsplash.com/photos/4ojhpgKpS68"
		img := models.Images{
		URL: url,
	}
		images = append(images, img)
	}

		price, err := strconv.Atoi(c.PostForm("price"))
		if err != nil {
		response.JSON(c, "", http.StatusBadRequest, nil, []string{err.Error()})
		return
	}

		numOfRooms, err := strconv.Atoi(c.PostForm("no_of_rooms"))
		if err != nil {
		response.JSON(c, "", http.StatusBadRequest, nil, []string{err.Error()})
		return
	}

		furnished, err := strconv.ParseBool(c.PostForm("furnished"))
		if err != nil {
		response.JSON(c, "", http.StatusBadRequest, nil, []string{err.Error()})
		return
	}

		aStatus, err := strconv.ParseBool(c.PostForm("apartment_status"))
		if err != nil {
		response.JSON(c, "", http.StatusBadRequest, nil, []string{err.Error()})
		return
	}

		exteriors := strings.Split(c.PostFormArray("exterior")[0], ",")
		interiors := strings.Split(c.PostFormArray("interior")[0], ",")

		apartment := models.Apartment{
		UserID:          userId,
		Title:           c.PostForm("title"),
		CategoryID:      c.PostForm("category"),
		Description:     c.PostForm("description"),
		Price:           price,
		NoOfRooms:       numOfRooms,
		Furnished:       furnished,
		Location:        c.PostForm("location"),
		ApartmentStatus: models.ApartmentStatus(aStatus),
		Interiors:       GetInteriors(interiors),
		Exteriors:       GetExteriors(exteriors),
		Images:          images,
	}

		err = s.DB.UpdateApartment(&apartment, apartment.ID)
		if err != nil {
		response.JSON(c, "", http.StatusBadRequest, nil, []string{err.Error()})
		return
	}
		// upload the image to aws.
		response.JSON(c, "Apartment Successfully Added", http.StatusOK, apartment, nil)

	}

}


func (s *Server) handleGetApartmentByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		apartmentID := c.Param("apartmentID")
		apartment, errs := s.DB.GetApartmentByID(apartmentID)
		if errs != nil {
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
			return
		}
		response.JSON(c, "apartment retrieved successfully", http.StatusOK, apartment, nil)
		return
	}

	//	if err := s.DB.UpdateApartment(apartment, apartmentID); err != nil {
	//		log.Printf("update apartment error : %v\n", err)
	//		response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
	//		return
	//	}
	//	response.JSON(c, "apartment updated successfully", http.StatusOK, nil, nil)
	//	return
	//}
}

func (s *Server) GetApartmentDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		apartment := c.Param("apartmentID")
		apart, err := s.DB.ApartmentDetails(apartment)
		if err != nil {
			log.Printf("error retrieving apartment: %v\n", err)
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
			return
		}
		response.JSON(c, "apartment retrieved successfully", http.StatusOK, gin.H{"apartment": apart}, nil)
    return
  }
}

func (s *Server) handleGetInteriorFeatures() gin.HandlerFunc {
	return func(c *gin.Context) {
		// fetch the interior features from database
		interiorFeatures, err := s.DB.GetAllInteriorFeatures()
		if err != nil {
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"could not retrieve all interior features"})
			return
		}
		response.JSON(c, "here are the interior features", http.StatusOK, interiorFeatures, nil)
		return
	}
}

func (s *Server) handleGetExteriorFeatures() gin.HandlerFunc {
	return func(c *gin.Context) {
		// fetch the exterior features from database
		exteriorFeatures, err := s.DB.GetAllExteriorFeatures()
		if err != nil {
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"could not retrieve all interior features"})
			return
		}
		response.JSON(c, "here are the exterior features", http.StatusOK, exteriorFeatures, nil)
		return
	}
}


func (s *Server) handleGetCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		//fetch the categories from database
		categories, err := s.DB.GetAllCategory()
		if err != nil {
			response.JSON(c, "", http.StatusInternalServerError, nil,[]string{"could not retrieve all categories"})
			return
		}
		response.JSON(c,"here are the exterior features", http.StatusOK, categories,nil)
		return
	}
}