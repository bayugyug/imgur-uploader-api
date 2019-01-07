package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"
)

//dumper dummy dumper
func dumper(infos ...interface{}) {
	for _, v := range infos {
		j, _ := json.MarshalIndent(v, "", "\t")
		if pShowDebug {
			log.Println(string(j))
		}
	}
}

//manageImageLocalDownloader watcher for the queue [pending]
func manageImageLocalDownloader(isReady chan bool) {
	isReady <- true
	for {
		select {
		case rec := <-pUploaderChan:
			if rec.ID != "" {
				//save it
				rec.Status = StatusInProgress
				pImageHistoryChan <- rec
				//forget and process
				go downloadImage2LocalAndPush(rec)
			}
		}
	}
}

//manageImageHistory watcher for the queue [history]
func manageImageHistory(isReady chan bool) {
	isReady <- true
	for {
		select {
		case rec := <-pImageHistoryChan:
			if rec.ID != "" {
				//save it
				pImageHistory[rec.ID] = rec
			}
		}
	}
}

//downloadImage2LocalAndPush grab the file/conveert to binary/push to imgur
func downloadImage2LocalAndPush(rec *UploadRecords) {

	processed := make(map[string]string)
	//parse & download and convert
	for k, row := range rec.URLS {
		resp, rcode, derr := httpGet(row.Link, map[string]string{})

		if rcode == http.StatusOK && derr == nil {
			//push to imgur
			frm := url.Values{}
			frm.Set("image", base64.StdEncoding.EncodeToString([]byte(resp)))
			headers := fmtBearerHeader()
			postresp, postcode, perr := httpPost(pImgurPostImageURL, frm.Encode(), headers)
			if postcode == http.StatusOK && perr == nil {
				dumper(postresp)
				oks, apiret := fmtImgurResponse(postresp)
				if !oks || apiret == nil {
					rec.URLS[k].Status = StatusFailed
					rec.UploadedList.Failed = append(rec.UploadedList.Failed, row.Link)
					processed[row.Link] = row.Link
					continue
				}

				//save the result
				if !apiret.Success {
					rec.URLS[k].Status = StatusFailed
					rec.UploadedList.Failed = append(rec.UploadedList.Failed, row.Link)
					processed[row.Link] = row.Link
					continue
				}

				//Imgur
				rec.URLS[k].ImgurID = apiret.Data.ID
				rec.URLS[k].ImgurLink = apiret.Data.Link
				rec.URLS[k].Status = StatusComplete
				rec.UploadedList.Complete = append(rec.UploadedList.Complete, apiret.Data.Link)
				processed[row.Link] = row.Link
			} else {
				rec.URLS[k].Status = StatusFailed
				rec.UploadedList.Failed = append(rec.UploadedList.Failed, row.Link)
				processed[row.Link] = row.Link
				dumper("FAIL_UPLOAD_IMG", postresp, postcode, perr)
			}
		} else {
			//fail to grab
			rec.URLS[k].Status = StatusFailed
			rec.UploadedList.Failed = append(rec.UploadedList.Failed, row.Link)
			processed[row.Link] = row.Link
			dumper("FAIL_GRAB_RAW_IMG", resp, rcode, derr)
		}
	}

	//remove the pending if needed
	var pending []string
	for _, v := range rec.UploadedList.Pending {
		if _, oks := processed[v]; !oks {
			pending = append(pending, v)
		}
	}

	//done
	rec.UploadedList.Pending = pending
	rec.Status = StatusComplete
	rec.Finished = time.Now().Format(time.RFC3339)
	dumper(rec)

	//maybe overwrite it
	pImageHistoryChan <- rec
}

//fmtImgurResponse format imgur response
func fmtImgurResponse(body string) (bool, *ImgurResult) {
	var apires ImgurResult
	err := json.Unmarshal([]byte(body), &apires)
	if err != nil {
		dumper("IMGUR_FORMAT_RESULT:", err)
		return false, nil
	}
	return true, &apires
}
