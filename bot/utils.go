package bot

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"net/http"
)

func DownloadFile(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", response.Status)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func ConvertImageToBytes(image image.Image) (*bytes.Buffer, error) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	err := jpeg.Encode(w, image, nil)
	if err != nil {
		return nil, err
	}
	return &b, nil
}
