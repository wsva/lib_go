package net

import (
	"net"
	"net/http"
	"regexp"
	"strings"
)

/*
get public ip of this host
use IP.String() to get a string
*/
func GetHostIP() net.IP {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP
}

/*
get public ip of this host
use IP.String() to get a string

dest ip:port
*/
func GetHostIPv2(dest string) net.IP {
	conn, _ := net.Dial("udp", dest)
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP
}

// use IP.String() to get a string
func GetIPFromRequest(r *http.Request) net.IP {
	var remoteIP net.IP
	// the default is the originating ip. but we try to find better options because this is almost
	// never the right IP
	if parts := strings.Split(r.RemoteAddr, ":"); len(parts) == 2 {
		remoteIP = net.ParseIP(parts[0])
	}
	// If we have a forwarded-for header, take the address from there
	if xff := strings.Trim(r.Header.Get("X-Forwarded-For"), ","); len(xff) > 0 {
		addrs := strings.Split(xff, ",")
		lastFwd := addrs[len(addrs)-1]
		if ip := net.ParseIP(lastFwd); ip != nil {
			remoteIP = ip
		}

	} else if xri := r.Header.Get("X-Real-Ip"); len(xri) > 0 {
		// parse X-Real-Ip header
		if ip := net.ParseIP(xri); ip != nil {
			remoteIP = ip
		}
	} else if xri := r.Header.Get("X-Real-IP"); len(xri) > 0 {
		// parse X-Real-IP header
		if ip := net.ParseIP(xri); ip != nil {
			remoteIP = ip
		}
	}

	return remoteIP
}

/*
destination: ip:port
*/
func GetTCPSourceIP(destination string) (net.IP, error) {
	conn, err := net.Dial("udp", destination)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return conn.LocalAddr().(*net.UDPAddr).IP, nil
}

func GetAvailablePort() (int, error) {
	conn, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	return conn.Addr().(*net.TCPAddr).Port, nil
}

func GetFirstPathFromRequest(r *http.Request) string {
	reg := regexp.MustCompile(`\w+`)
	return reg.FindString(r.URL.Path)
}

func GetSecondPathFromRequest(r *http.Request) string {
	reg := regexp.MustCompile(`\w+`)
	pathList := reg.FindAllString(r.URL.Path, 2)
	if len(pathList) > 1 {
		return pathList[1]
	}
	return ""
}

func GetQueryValueListFromRequest(r *http.Request, key string) []string {
	query := r.URL.Query()
	return query[key]
}

func GetOneQueryValueFromRequest(r *http.Request, key string) string {
	query := r.URL.Query()
	return query[key][0]
}

func GetOneQueryMapFromRequest(r *http.Request) map[string]string {
	queryMap := make(map[string]string)
	query := r.URL.Query()
	for k := range query {
		queryMap[k] = query[k][0]
	}
	return queryMap
}

func GetSchemaAndHost(r *http.Request) string {
	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}
	return strings.Join([]string{scheme, r.Host}, "")
}

func GetFullURL(r *http.Request) string {
	return GetSchemaAndHost(r) + r.RequestURI
}
