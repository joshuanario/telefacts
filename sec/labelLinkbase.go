package sec

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"ecks-bee.com/telefacts/xbrl"
)

func getLabelLinkbaseFromOSfiles(files []os.FileInfo) (os.FileInfo, error) {
	var xmls []os.FileInfo
	for _, file := range files {
		s := file.Name()
		if file.IsDir() || filepath.Ext(s) == "" {
			continue
		}
		if s[len(s)-len(labExt):len(s)] == labExt {
			xmls = append(xmls, file)
		}
	}
	if len(xmls) <= 0 {
		return nil, fmt.Errorf("No label linkbase found")
	}
	return xmls[0], nil
}

func getLabelLinkbaseFromUnzipfiles(unzipFiles []*zip.File) (*zip.File, error) {
	var xmls []*zip.File
	for _, unzipFile := range unzipFiles {
		s := unzipFile.Name
		if s[len(s)-len(labExt):len(s)] == labExt {
			xmls = append(xmls, unzipFile)
		}
	}
	if len(xmls) <= 0 {
		return nil, fmt.Errorf("No label linkbase found")
	}
	return xmls[0], nil
}

func unzipLabelLinkbase(unzipFile *zip.File) (*xbrl.LabelLinkbase, error) {
	rc, err := unzipFile.Open()
	defer rc.Close()
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, rc)
	if err != nil {
		return nil, err
	}
	decoded, err := xbrl.DecodeLabelLinkbase(buffer.Bytes())
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

func commitLabelLinkbase(dest string, linkbase *xbrl.LabelLinkbase) error {
	data, _ := xml.Marshal(linkbase)
	file, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	defer file.Close()
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func getLabelLinkbaseFromFilingItems(filingItems []filingItem, ticker string) (*filingItem, error) {
	for _, f := range filingItems {
		s := f.Name
		ext := filepath.Ext(s)
		a := (ext == xmlExt && strings.Index(s, ticker) == 0)
		b := len(s) >= 8
		if b {
			longExt := s[len(s)-8:]
			b = longExt == labExt
		}
		if a && b {
			return &f, nil
		}
	}
	return nil, fmt.Errorf("Cannot identify a single label linkbase")
}

func scrapeLabelLinkbaseFromSEC(filingURL string, filingItem *filingItem) (*xbrl.LabelLinkbase, error) {
	resp, err := http.Get(filingURL + "/" + filingItem.Name)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, resp.Body)
	if err != nil {
		return nil, err
	}
	return xbrl.DecodeLabelLinkbase(buffer.Bytes())
}