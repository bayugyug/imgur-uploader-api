package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"
)

func dumper(infos ...interface{}) {
	for _, v := range infos {
		j, _ := json.MarshalIndent(v, "", "\t")
		if pShowDebug {
			log.Println(string(j))
		}
	}
}

func manageImageLocalDownloader(isReady chan bool) {
	isReady <- true
	for {
		select {
		case rec := <-pUploaderChan:
			if rec.ID != "" {
				downloadImage2LocalAndPush(rec)
			}
		}
	}
}

func downloadImage2LocalAndPush(rec *UploadRecords) {

	rec.Status = StatusInProgress

	processed := make(map[string]string)

	//save it
	pImageHistory[rec.ID] = rec
	
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
			}
		} else {
			//fail to grab
			rec.URLS[k].Status = StatusFailed
			rec.UploadedList.Failed = append(rec.UploadedList.Failed, row.Link)
			processed[row.Link] = row.Link
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
	pImageHistory[rec.ID] = rec

}

func fmtImgurResponse(body string) (bool, *ImgurResult) {
	var apires ImgurResult
	err := json.Unmarshal([]byte(body), &apires)
	if err != nil {
		dumper("IMGUR_FORMAT_RESULT:", err)
		return false, nil
	}
	return true, &apires
}
