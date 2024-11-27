package main

import (
  "encoding/json"
  "errors"
  "io/ioutil"
  "net/http"
  "net/url"
)

func ipLookup(ip string) (string, error) {
  key := "XYZ"
  u := url.URL{
    Scheme: "https",
    Host: "api.ipgeolocation.io",
    Path: "ipgeo",
  }
  q := u.Query()
  q.Set("apiKey", key)
  q.Set("ip", ip)
  u.RawQuery = q.Encode()
  resp, err := http.Get(u.String())
  if err != nil {
    return "", err
  }
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return "", err
  }
  data := map[string]string{}
  json.Unmarshal(body, &data)
  _, found := data["ip"]
  if !found {
    return "", errors.New(data["message"])
  }
  return ip + "" +
    data["country_name"] + "" +
    data["state_proc"] + "" +
    data["city"] + "" +
    data["isp"], nil
}
