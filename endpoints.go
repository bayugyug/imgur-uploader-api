package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

//ApiHandlerSub skeleton interface
type ApiHandlerSub interface {
	GetOneImage(w http.ResponseWriter, r *http.Request)
	GetAllImages(w http.ResponseWriter, r *http.Request)
	SetUserCode(w http.ResponseWriter, r *http.Request)
	IndexPage(w http.ResponseWriter, r *http.Request)
	UploadImage(w http.ResponseWriter, r *http.Request)
	CheckAuthBearer(w http.ResponseWriter, r *http.Request) bool
	SanityCheckImageUpload(imgs ParamImageURLS) (int, string)
	ReplyNoContent(w http.ResponseWriter, r *http.Request)
}

//ApiHandler empty holder
type ApiHandler struct{}

//UploadImage upload images to the imgur api
func (api ApiHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	//auth-checker
	if !api.CheckAuthBearer(w, r) {
		return
	}

	var images ParamImageURLS
	err := json.NewDecoder(r.Body).Decode(&images)
	if err != nil {
		dumper("INVALID_PARAMETER_FORMAT:", err)
		api.ReplyNoContent(w, r)
		return
	}
	//just in case :-)
	defer r.Body.Close()

	//chk if url's valid and/or uniq
	valid, jobId := api.SanityCheckImageUpload(images)
	if valid <= 0 || jobId == "" {
		dumper("SANITY_CHECK_FAILED")
		api.ReplyNoContent(w, r)
		return
	}

	//give the jobId ;-)
	render.JSON(w, r, UploadJobResponse{JobID: jobId})
}

//GetAllImages get list of images from mem
func (api ApiHandler) GetAllImages(w http.ResponseWriter, r *http.Request) {

	//auth-checker
	if !api.CheckAuthBearer(w, r) {
		return
	}
	reply := getListOfImages()
	render.JSON(w, r, reply)
}

//GetOneImage get image info from mem
func (api ApiHandler) GetOneImage(w http.ResponseWriter, r *http.Request) {

	//auth-checker
	if !api.CheckAuthBearer(w, r) {
		return
	}

	rok, reply := getImageByJobId(chi.URLParam(r, "id"))
	if !rok || reply == nil {
		dumper("JOB_ID_NOT_IN_HISTORY_LOG")
		api.ReplyNoContent(w, r)
		return
	}

	//okay
	render.JSON(w, r, reply)
}

//SetUserCode set the auth-code on imgur api
func (api ApiHandler) SetUserCode(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimSpace(chi.URLParam(r, "code"))
	if code != "" {
		pParamConfig.Code = code
		if tok := getAuthBearerToken(); tok != "" {
			dumper("BEARER_CODE_GET_PERMISSION_FAILED")
			api.ReplyNoContent(w, r)
			return
		}
	} else if pParamConfig.Code == "" {
		dumper("BEARER_CODE_EMPTY")
		api.ReplyNoContent(w, r)
		return
	}
	render.JSON(w, r, APIResponse{
		Code:    http.StatusAccepted,
		Message: http.StatusText(http.StatusAccepted),
	})
}

//IndexPage root page
func (api ApiHandler) IndexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, fmt.Sprintf("Welcome to Imgur-Uploader!\n\nVersion: %s\n", pVersion))
}

//SanityCheckImageUpload check if the params format is valid
func (api ApiHandler) SanityCheckImageUpload(imgs ParamImageURLS) (int, string) {
	if len(imgs.URLS) <= 0 {
		return -1, ""
	}

	//for testing purposes, dont set limit :-)
	var validImgs []*URLInfo
	var imgsPending []string

	imgUniqHash := make(map[string]string)

	for _, v := range imgs.URLS {
		v = strings.TrimSpace(v)
		//chk old
		if !strings.HasPrefix(v, "http") {
			continue
		}
		hexmd5 := fmt.Sprintf("%x", md5.Sum([]byte(v)))
		if _, found := imgUniqHash[hexmd5]; !found {
			imgUniqHash[hexmd5] = v
			validImgs = append(validImgs, &URLInfo{
				Link:   v,
				Status: StatusPending,
			})
			imgsPending = append(imgsPending, v)
		}
	}

	//try if all ok
	validT := len(validImgs)
	if validT <= 0 {
		return -1, ""
	}

	//get uuid
	jobId := genJobId()

	//async downloader
	pUploaderChan <- &UploadRecords{
		Status:       StatusPending,
		ID:           jobId,
		Created:      time.Now().Format(time.RFC3339),
		URLS:         validImgs,
		UploadedList: &UploadedList{Pending: imgsPending},
	}

	dumper("VALID_UPLOAD_IMAGES:", validImgs, jobId)

	return validT, jobId

}

//CheckAuthBearer check if the config is set
func (api ApiHandler) CheckAuthBearer(w http.ResponseWriter, r *http.Request) bool {

	//check if set
	if pParamConfig == nil || pParamConfig.Bearer == nil {
		render.JSON(w, r, APIResponse{
			Code:    http.StatusNonAuthoritativeInfo,
			Message: http.StatusText(http.StatusNonAuthoritativeInfo)})
		return false
	}
	//chk type
	if !strings.EqualFold(pParamConfig.Bearer.TokenType, "bearer") ||
		pParamConfig.Bearer.RefreshToken == "" ||
		pParamConfig.Bearer.AccessToken == "" {
		dumper("AUTH_BEARER_EMPTY", pParamConfig)
		render.JSON(w, r, APIResponse{
			Code:    http.StatusNonAuthoritativeInfo,
			Message: http.StatusText(http.StatusNonAuthoritativeInfo)})
		return false
	}

	//good
	return true
}

//ReplyNoContent send 204 msg
func (api ApiHandler) ReplyNoContent(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, APIResponse{
		Code:    http.StatusNoContent,
		Message: http.StatusText(http.StatusNoContent),
	})
}

//genJobId format the uuid
func genJobId() string {
	return uuid.New().String() + `-` + time.Now().Format("20060102-150405")
}

//getListOfImages get list of images status since the start
func getListOfImages() ListImage {
	//get a copy :-)
	var all ListImage
	images := pImageHistory
	//iterate
	for _, row := range images {
		if row.Status == StatusComplete {
			for _, rec := range row.URLS {
				if rec.ImgurID != "" {
					all.Uploaded = append(all.Uploaded, rec.ImgurLink)
				}
			}
		}
	}

	return all
}

//getImageByJobId get image based on id
func getImageByJobId(idx string) (bool, *UploadRecords) {
	//get a copy :-)
	images, found := pImageHistory[idx]
	return found, images
}
