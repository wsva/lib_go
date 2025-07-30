package http

import (
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strings"

	"github.com/wsva/lib_go/sets"
)

func CheckClientCertExist(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if len(r.TLS.PeerCertificates) < 1 {
		w.Write([]byte("no certificate found"))
		return
	} else if len(r.TLS.PeerCertificates) > 1 {
		w.Write([]byte("two or more certificates found"))
		return
	} else {
		next(w, r)
	}
}

/*
CheckClientCertIP verifies the ip(s) in certificate with real ip
*/
func CheckClientCertIP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	realip := parseIP(r)
	ipSet := GetClientCertIPSet(r)
	if ipSet.Has(realip) {
		next(w, r)
	} else {
		fmt.Fprintf(w, "certificate allows only %v", strings.Join(ipSet.UnsortedList(), ","))
		return
	}
}

func GetClientCertCommonName(r *http.Request) string {
	reg := regexp.MustCompile(`CN=`)
	cert := r.TLS.PeerCertificates[0]
	return reg.ReplaceAllString(cert.Subject.CommonName, "")
}

func GetClientCertIPSet(r *http.Request) sets.Set[string] {
	cert := r.TLS.PeerCertificates[0]
	set := sets.New[string]()
	for _, v := range cert.IPAddresses {
		set.Insert(v.String())
	}
	return set
}

func parseIP(r *http.Request) string {
	var remoteIP net.IP
	// the default is the originating ip. but we try to find better options because this is almost
	// never the right IP
	if parts := strings.Split(r.RemoteAddr, ":"); len(parts) == 2 {
		remoteIP = net.ParseIP(parts[0])
	}

	//change names in header to uppercase
	nh := make(http.Header)
	for k := range r.Header {
		nh[strings.ToLower(k)] = r.Header[k]
	}

	// If we have a forwarded-for header, take the address from there
	if xff := strings.Trim(nh.Get("x-forwarded-for"), ","); len(xff) > 0 {
		addrs := strings.Split(xff, ",")
		lastFwd := addrs[len(addrs)-1]
		if ip := net.ParseIP(lastFwd); ip != nil {
			remoteIP = ip
		}
		// parse X-Real-Ip header
	} else if xri := nh.Get("x-real-ip"); len(xri) > 0 {
		if ip := net.ParseIP(xri); ip != nil {
			remoteIP = ip
		}
	}

	return remoteIP.String()
}
