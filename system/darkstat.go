package system

import (
	"encoding/base64"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Host struct {
	IP        string `json:"ip"`
	HostName  string `json:"hostname"`
	Mac       string `json:"mac"`
	RX        int64  `json:"rx"`
	TX        int64  `json:"tx"`
	Total     int64  `json:"total"`
	RXStep    int64  `json:"rxStep"`
	TXStep    int64  `json:"txStep"`
	TotalStep int64  `json:"totalStep"`
	Key       string `json:"key"`
	MaskedIP  string `json:"maskedIP"`
	MaskedMAC string `json:"maskedMAC"`
}

var localoHosts []Host
var weoHosts []Host
var orangeoHosts []Host

func Publish(url string, include string, sender int) ([]Host, error) {
	var nHosts []Host
	var err error
	if sender != 2 {
		nHosts, err = getCurrentHostsReadings(url, include)
	} else {
		nHosts, err = getCurrentOrangeHostsReadings(url, include)
	}

	nHosts = calc(nHosts, sender)
	return clean(nHosts), err
}

func clean(nHosts []Host) []Host {
	clean := []Host{}
	for _, h := range nHosts {
		if h.TotalStep == 0 {
			continue
		}
		clean = append(clean, h)
	}
	return clean
}

func calc(nHosts []Host, sender int) []Host {
	var oHosts []Host
	if sender == 0 {
		oHosts = localoHosts
	} else if sender == 1 {
		oHosts = weoHosts
	} else {
		oHosts = orangeoHosts
	}
	temp := []Host{}
	if len(oHosts) > 0 {
		for _, n := range nHosts {
			for _, o := range oHosts {
				if strings.EqualFold(n.Key, o.Key) {
					if n.Total >= o.Total {
						n.RXStep = n.RX - o.RX
						n.TXStep = n.TX - o.TX
						n.TotalStep = n.Total - o.Total
					} else {
						n.RXStep = 0
						n.TXStep = 0
						n.TotalStep = 0
					}
					temp = append(temp, n)
					break
				}
			}
		}
	}
	if sender == 0 {
		localoHosts = nHosts
	} else if sender == 1 {
		weoHosts = nHosts
	} else {
		orangeoHosts = nHosts
	}
	return temp
}

func getCurrentOrangeHostsReadings(url string, include string) ([]Host, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Panic(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Panicf("status code error: %d %s", res.StatusCode, res.Status)
	}
	// Load the XML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Panic(err)
	}

	download := doc.Find("currentmonthdownload").Map(func(i int, sel *goquery.Selection) string {
		return strings.Trim(sel.Text(), "")
	})

	upload := doc.Find("currentmonthupload").Map(func(i int, sel *goquery.Selection) string {
		return strings.Trim(sel.Text(), "")
	})

	hosts := []Host{}
	var h Host
	h.IP = strings.TrimSpace("192.168.20.1")
	h.HostName = strings.TrimSpace("orange.lan")
	h.Mac = strings.ToUpper(strings.TrimSpace("d8:9e:61:d9:c8:97"))
	h.RX, _ = strconv.ParseInt(strings.Replace(strings.TrimSpace(download[0]), ",", "", -1), 10, 64)
	h.TX, _ = strconv.ParseInt(strings.Replace(strings.TrimSpace(upload[0]), ",", "", -1), 10, 64)
	h.Total = h.RX + h.TX
	h.Key = base64.StdEncoding.EncodeToString([]byte(h.Mac + "|" + h.IP))
	h.MaskedIP = base64.StdEncoding.EncodeToString(net.ParseIP(h.IP).To4())
	haddr, _ := net.ParseMAC(h.Mac)
	h.MaskedMAC = base64.StdEncoding.EncodeToString(haddr)
	hosts = append(hosts, h)

	return hosts, nil
}

func getCurrentHostsReadings(url string, include string) ([]Host, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Panic(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Panicf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Panic(err)
	}

	words := doc.Find("tr").Map(func(i int, sel *goquery.Selection) string {
		return strings.Trim(sel.Text(), "")
	})

	hosts := []Host{}
	for i, s := range words {
		if i > 0 {
			host_src := strings.Split(s, "\n")
			if strings.Contains(host_src[1], include) {
				clean := []string{}
				for i := 0; i < len(host_src); i++ {
					if strings.EqualFold(host_src[i], "") {
						continue
					}
					clean = append(clean, host_src[i])
				}
				var h Host
				h.IP = strings.TrimSpace(clean[0])
				h.HostName = strings.TrimSpace(clean[1])
				h.Mac = strings.ToUpper(strings.TrimSpace(clean[2]))
				h.RX, _ = strconv.ParseInt(strings.Replace(strings.TrimSpace(clean[3]), ",", "", -1), 10, 64)
				h.TX, _ = strconv.ParseInt(strings.Replace(strings.TrimSpace(clean[4]), ",", "", -1), 10, 64)
				h.Total, _ = strconv.ParseInt(strings.Replace(strings.TrimSpace(clean[5]), ",", "", -1), 10, 64)
				h.Key = base64.StdEncoding.EncodeToString([]byte(h.Mac + "|" + h.IP))
				h.MaskedIP = base64.StdEncoding.EncodeToString(net.ParseIP(h.IP).To4())
				haddr, _ := net.ParseMAC(h.Mac)
				h.MaskedMAC = base64.StdEncoding.EncodeToString(haddr)
				hosts = append(hosts, h)
			}
		}
	}
	return hosts, nil
}
