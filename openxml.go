package metadata

import (
  "archive/zip"
  "encoding/xml"
  "strings"
  "fmt"
)

var OfficeVersions = map[string]string{
  "16": "2016",
  "15": "2013",
  "14": "2010",
  "12": "2007",
  "11": "2003",
}

type OfficeCoreProperty struct {
  XMLName        xml.Name `xml:"coreProperties"`
  Creator        string   `xml:"creator"`
  LastModifiedBy string   `xml:"lastModifiedBy"`
}

type OfficeAppProperty struct {
  XMLName        xml.Name `xml:"Properties"`
  Application    string   `xml:"Application"`
  Company        string   `xml:"Company"`
  Version        string    `xml:"AppVersion"`
}

func (a *OfficeAppProperty) GetMajorVersion() string {
  tokens := strings.Split(a.Version, ".")

  if len(tokens) < 2 {
    fmt.Println("token was less than 2")
    return "Unknown"
  }
  v, ok := OfficeVersions [tokens[0]]
  if !ok {
    fmt.Println("office version was whack!")
    return "Unknown"
  }
  return v
}

func NewProperties(r *zip.Reader) (*OfficeCoreProperty, *OfficeAppProperty, error) {
  var coreProps OfficeCoreProperty
  var appProps OfficeAppProperty

  for _, f := range r.File {
    switch f.Name {
    case "docProps/core.xml":
      if err := process(f, &coreProps); err != nil {
        fmt.Println("processing the file failed - coreProps")
        return nil, nil, err
      }
    case "docProps/app.xml":
      if err := process(f, &appProps); err != nil {
        fmt.Println("processing the file failed = appProps")
        return nil, nil, err
      }
    default:
      continue
    }
  }
  return &coreProps, &appProps, nil
}

func process(f *zip.File, prop interface{}) error {
  rc, err := f.Open()
  if err != nil {
    fmt.Println("failed to open the file process()")
    return err
  }
  defer rc.Close()

  if err := xml.NewDecoder(rc).Decode(&prop); err != nil {
    fmt.Println("failed to decode the file - after it was opened")
    return err
  }
  return nil
}
