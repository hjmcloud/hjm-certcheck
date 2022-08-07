package service

import (
	"context"
	"crypto/tls"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type sCertInfo struct {
	Host         string
	Port         string
	Issuer       string
	CommonName   string
	NotBefore    time.Time
	NotAfter     time.Time
	NotAfterUnix int64
	SANs         []string
}

func GetCertInfo(ctx context.Context, req string) (certInfo *sCertInfo, err error) {
	request := strings.ToLower(string(req))
	if !strings.HasPrefix(request, "http") {
		request = "https://" + request
	}
	u, err := url.Parse(request)
	if err != nil {
		g.Log().Error(ctx, "failed to parse url.")
		return nil, err
	}
	// if u.Scheme is not https, return error
	if u.Scheme != "https" {
		return nil, gerror.New("scheme is not https")
	}
	// if u.port is null then use 443
	if u.Port() == "" {
		u.Host = u.Host + ":443"
	}
	address := u.Hostname() + ":" + u.Port()
	ipConn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		g.Log().Errorf(ctx, "SSL/TLS not enabed on %v\nDial error: %v ", u.Hostname())
		return nil, err
	}
	defer ipConn.Close()
	conn := tls.Client(ipConn, &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         u.Hostname(),
	})
	if err = conn.Handshake(); err != nil {
		g.Log().Infof(ctx, "Invalid SSL/TLS for %v\nHandshake error: %v", address, err)
		return nil, err
	}
	defer conn.Close()
	cert := conn.ConnectionState().PeerCertificates[0]
	certInfo = &sCertInfo{
		Host:         u.Hostname(),
		Port:         u.Port(),
		Issuer:       cert.Issuer.CommonName,
		CommonName:   cert.Subject.CommonName,
		NotBefore:    cert.NotBefore,
		NotAfter:     cert.NotAfter,
		NotAfterUnix: cert.NotAfter.Unix(),
		SANs:         cert.DNSNames,
	}
	return certInfo, nil
}
