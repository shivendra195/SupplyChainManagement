package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/multi/qrcode"
	"image"
	"io"
	"net/http"
	"strings"

	// import gif, jpeg, png
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	// import bmp, tiff, webp
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

func (srv *Server) ScanQR(w http.ResponseWriter, r *http.Request) {

	b := new(bytes.Buffer)
	if _, err := io.Copy(b, r.Body); err != nil {
		msg := fmt.Sprintf("Failed to read request body: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	res, err := scan(b.Bytes())
	if err != "" {
		msg := fmt.Sprintf("Internal server error: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	fmt.Println([]byte(res))
	_, erro := w.Write([]byte(res))
	if erro != nil {
		msg := fmt.Sprintf("Internal server error: %v", erro)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	//json.Unmarshal([]byte(res))

	//utils.EncodeJSON200Body(w, []byte(res))
}

func scan(b []byte) (string, string) {
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		msg := fmt.Sprintf("failed to read image: %v", err)
		return "", msg
	}

	source := gozxing.NewLuminanceSourceFromImage(img)
	bin := gozxing.NewHybridBinarizer(source)
	bbm, err := gozxing.NewBinaryBitmap(bin)

	if err != nil {
		msg := fmt.Sprintf("error during processing: %v", err)
		return "", msg
	}

	qrReader := qrcode.NewQRCodeMultiReader()
	result, err := qrReader.DecodeMultiple(bbm, nil)
	if err != nil {
		msg := fmt.Sprintf("unable to decode QRCode: %v", err)
		return "", msg
	}
	strRes := []string{}
	//var strByte []byte
	for _, element := range result {
		strRes = append(strRes, element.String())
		//strByte = element.GetRawBytes()
		//strByte = append(strByte, element.GetRawBytes()...)
	}

	res := strings.Join(strRes, "")
	jsonData := make(map[string]interface{})
	err = json.Unmarshal([]byte(res), &jsonData)
	if err != nil {
		fmt.Printf("unable to unmarshal byte data: %v", err)
		//return "", msg
	}
	fmt.Println("data", jsonData)
	return res, ""
}
