package agave

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

type Client struct {
	App         string
	Channel     string
	SensorGUID  string
	SensorIp    string
	SensorPort  int
	ipCache     map[string]bool
	ipCacheLock sync.RWMutex
}

// HTTPAttack normalizes the recieved request and allows for easy marshaling into JSON.
type HTTPAttack struct {
	Protocol   string
	App        string
	Channel    string
	SensorGUID string `json:"sensor"`
	DestPort   int
	DestIp     string
	SrcPort    int
	SrcIp      string
	Signature  string
	PrevSeen   bool // True if we've seen this before
	Request    *RequestJson
}

// CredentialAttack normalizes the recieved request and allows for easy marshaling into JSON.
type CredentialAttack struct {
	Protocol   string
	App        string
	Channel    string
	SensorGUID string `json:"sensor"`
	DestPort   int
	DestIp     string
	SrcPort    int
	SrcIp      string
	Username   string `json:"agave_username"`
	Password   string `json:"agave_password"`
}

func NewClient(app string, channel string, guid string, ip string, port int) *Client {
	return &Client{App: app, Channel: channel, SensorGUID: guid, SensorIp: ip, SensorPort: port}
}

func (c *Client) NewHTTPAttack(signature string, r *http.Request) (*HTTPAttack, error) {
	ip, p, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(p)
	if err != nil {
		return nil, err
	}

	rj := TrimRequest(r)

	return &HTTPAttack{
		Protocol:   r.Proto,
		App:        c.App,
		Channel:    c.Channel,
		SensorGUID: c.SensorGUID,
		DestPort:   c.SensorPort,
		DestIp:     c.SensorIp,
		SrcPort:    port,
		SrcIp:      ip,
		PrevSeen:   c.SeenIP(ip),
		Request:    rj,
	}, nil
}

func (c *Client) NewCredentialAttack(r *http.Request, username string, password string) (*CredentialAttack, error) {
	ip, p, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(p)
	if err != nil {
		return nil, err
	}

	return &CredentialAttack{
		Protocol:   r.Proto,
		App:        c.App,
		Channel:    c.Channel,
		SensorGUID: c.SensorGUID,
		DestPort:   c.SensorPort,
		DestIp:     c.SensorIp,
		SrcPort:    port,
		SrcIp:      ip,
		Username:   username,
		Password:   password,
	}, nil
}

func (c *Client) SaveIP(ip string) {
	c.ipCacheLock.Lock()
	c.ipCache[ip] = true
	c.ipCacheLock.Unlock()
}

func (c *Client) SeenIP(ip string) bool {
	c.ipCacheLock.RLock()
	seen := c.ipCache[ip]
	c.ipCacheLock.RUnlock()
	return seen
}

func TrimRequest(r *http.Request) *RequestJson {
	body, _ := ioutil.ReadAll(r.Body)
	r.ParseForm()
	rj := &RequestJson{
		Method:           r.Method,
		URL:              r.URL,
		Proto:            r.Proto,
		ProtoMajor:       r.ProtoMajor,
		ProtoMinor:       r.ProtoMinor,
		Header:           r.Header,
		Body:             body,
		TransferEncoding: r.TransferEncoding,
		Host:             r.Host,
		PostForm:         r.PostForm,
	}

	return rj
}

type RequestJson struct {
	Method           string
	URL              *url.URL
	Proto            string
	ProtoMajor       int
	ProtoMinor       int
	Header           http.Header
	Body             []byte
	TransferEncoding []string
	Host             string
	PostForm         url.Values
}

// getHost tries its best to return the request host.
func getHost(r *http.Request) string {
	r.URL.Scheme = "http"
	r.URL.Host = r.Host

	return r.URL.String()
}
