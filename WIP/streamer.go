package main

import (
	"fmt"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net/http"
	"os"
	"encoding/json"
	"strings"
	"sort"
	"io"
	"log"
)

const (
	__ENDPOINT = "http://ewr-changes-x6.ewr.mmracks.internal:8082"
	__TOPIC    = "ewr.changes.snapshot.currency-rates.1"

	PROXY_ADDR = "127.0.0.1:8123"
)

var __TopicsURL = func() string { return __ENDPOINT + "/topics" }()
var __SubscriptionName = func() string { return "ts_" + __TOPIC }()
var __SubscriptionURL = func() string { return __ENDPOINT + "/consumers/" + __SubscriptionName }()
var __CommonURL = func() string { return __SubscriptionURL + "/instances/" + __TOPIC }()
var __DataURL = func() string { return __CommonURL + "/topics/" + __TOPIC }()
var __UnsubscribeURL = func() string { return __CommonURL }()

//var __MainContext = context.Background()

/*
STEPS
1) Create a SOCKS5 channel (run exec.Command ?)
ssh -D 8123 -C -N n2cstech &

2) Setup constants - even for topics ?
ENDPOINT="http://ewr-changes-x6.ewr.mmracks.internal:8082"
TOPIC="snapshot-campaigns-1"
SUBS_NAME="ts_campaign_consumer"

http --proxy=http:socks5://127.0.0.1:8123 ${ENDPOINT}/topics
http --proxy=http:socks5://127.0.0.1:8123 POST ${ENDPOINT}/consumers/${SUBS_NAME} Content-Type:application/vnd.kafka.v1+json name=${TOPIC} format=avro auto.offset.reset=smallest
http --proxy=http:socks5://127.0.0.1:8123 ${ENDPOINT}/consumers/${SUBS_NAME}/instances/${TOPIC}/topics/${TOPIC} Accept:application/vnd.kafka.avro.v1+json
http --proxy=http:socks5://127.0.0.1:8123 DELETE ${ENDPOINT}/consumers/${SUBS_NAME}/instances/${TOPIC}
*/

func setupHTTPClient() *http.Client {
	// create a socks5 dialer
	dialer, err := proxy.SOCKS5("tcp", PROXY_ADDR, nil, proxy.Direct)
	if err != nil {
		fmt.Println(os.Stderr, "can't connect to proxy:", err)
		os.Exit(1)
	}
	// setup a http client
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	// set our socks5 as the dialer
	httpTransport.Dial = dialer.Dial
	return httpClient
}

func postThroughProxy(httpClient *http.Client, url string, body io.Reader,
	headers map[string][]string) string {
	// create a request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't create request:", err)
		os.Exit(2)
	}
	if headers != nil {
		for k, v := range headers {
			req.Header[k] = v
		}
	}
	//fmt.Println(req)
	// use the http client to fetch the page
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't POST page:", err)
		os.Exit(3)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading body:", err)
		os.Exit(4)
	}
	return string(b)
}

func deleteThroughProxy(httpClient *http.Client, url string) string {
	// create a request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't create request:", err)
		os.Exit(2)
	}
	//fmt.Println(req)
	// use the http client to fetch the page
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't DELETE page:", err)
		os.Exit(3)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading body:", err)
		os.Exit(4)
	}
	return string(b)
}

func getThroughProxy(check bool, httpClient *http.Client, url string, headers map[string][]string) string {
	// create a request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't create request:", err)
		os.Exit(2)
	}
	if headers != nil {
		for k, v := range headers {
			req.Header[k] = v
		}
	}
	//fmt.Println(req)
	// use the http client to fetch the page
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't GET page:", err)
		os.Exit(3)
	}
	defer resp.Body.Close()
	//m := v.(map[string]interface{})["value"].(map[string]interface{})["fields"].(map[string]interface{})
	type MessageFields struct {
		Id, Version string
		Currency_Code string
		Rate string
		Date string
		//Created_on, Updated_on string
	}
	type MessageValue struct {
		//Changes json.RawMessage
		Fields MessageFields
	}
	type Message struct {
		//Key, Partition, Offset float64
		Value MessageValue
	}
	currency_map := make(map[string]string)
	if check {
		dec := json.NewDecoder(resp.Body)
		// read open bracket
		_, err := dec.Token()
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("%T: %v\n", t, t)

		array_length := 0
		// while the array contains values
		for /*i:= 0; */dec.More()/* && i < 10; i++*/ {
			array_length++
			var m Message
			// decode an array value (Message)
			err := dec.Decode(&m)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("%v\n", m)
			fields := m.Value.Fields
			currency_map[fields.Currency_Code] = fmt.Sprint(fields.Rate, " - ", fields.Date)
		}

		// read closing bracket
		_, err = dec.Token()
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("%T: %v\n", t, t)
		/*for i := 0; i < 100; i++ {
			t, err := dec.Token()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%T: %v", t, t)
			if dec.More() {
				fmt.Printf(" (more)")
			}
			fmt.Printf("\n")
		}*/
		fmt.Println("ARRAY LENGTH:", array_length)
		fmt.Println("DATA MAP LENGTH:", len(currency_map))
	} else {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error reading body:", err)
			os.Exit(4)
		}
		return string(b)
	}
	b, err := json.MarshalIndent(currency_map, "\t", "")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error marshalling:", err)
		os.Exit(4)
	}
	return string(b)
}

func main() { //http.ProxyURL() ??
	//setupBGProxyServer()
	httpClient := setupHTTPClient()

	getTopics(httpClient)

	_ = subscribeTopic(httpClient)
	fetchData(httpClient)
	unsubscribeTopic(httpClient)
}

func fetchData(httpClient *http.Client) {
	fmt.Println("EXE: Fetch Data:", __DataURL)
	var headers = map[string][]string{"Accept":{"application/vnd.kafka.avro.v1+json"}}
	currency_map := getThroughProxy(true, httpClient, __DataURL, headers)
	fmt.Println("DATA MAP:", currency_map)
	/*fetchDataReader := strings.NewReader()
	//checkType("fetchDataReader", fetchDataReader)
	dec := json.NewDecoder(fetchDataReader)
	var v interface{}
	dec.Decode(&v)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	av := v.([]interface{})
	var currency_map = make(map[string]string)
	for _, v := range av {
		m := v.(map[string]interface{})["value"].(map[string]interface{})["fields"].(map[string]interface{})
		//fmt.Println(m["id"], "(", m["version"], ") has currency_code:", m["currency_code"], "with value of:", m["rate"], "for date:", m["date"])
		currency_map[fmt.Sprint(m["currency_code"])] = fmt.Sprint(m["rate"], "->", m["date"])
	}
	fmt.Println("DATA LENGTH:", len(av))*/

	/*switch v.(type) {
	case nil:
		fmt.Println("nil value")
	case []interface{}:
		fmt.Println("array value")
	}*/
}

type sortable []interface{}

func (a sortable) Len() int           { return len(a) }
func (a sortable) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortable) Less(i, j int) bool { return a[i].(string) < a[j].(string) }

func getTopics(httpClient *http.Client) {
	fmt.Println("EXE: Fetch Topics:", __TopicsURL)
	var headers map[string][]string = nil //map[string][]string{"Content-Type":{"what"}}
	getTopicsReader := strings.NewReader(getThroughProxy(false, httpClient, __TopicsURL, headers))
	//checkType("getTopicsReader", getTopicsReader)
	dec := json.NewDecoder(getTopicsReader)
	var v interface{}
	dec.Decode(&v)
	var topics []interface{}
	switch v.(type) {
	case nil:
		fmt.Println("nil value")
	case []interface{}:
		topics = v.([]interface{})
	}
	//sort.Sort(sortable(topics))
	var strArray []string
	//var myStrArray sortable_string.StringArray
	for _, topic := range topics {
		//fmt.Println(topic)
		strArray = append(strArray, topic.(string))
		//myStrArray = append(myStrArray, topic.(string))
	}
	//sortable_topics :=

	//string.StringArray{strArray...}
	//sortable_string.Sort(topics.([]string)...)
	/*fmt.Println(myStrArray)
	sort.Sort(myStrArray)
	fmt.Println(myStrArray)*/
	if index := sort.SearchStrings(strArray, __TOPIC); index < len(strArray) {
		//fmt.Println("Found at index:", index, "::", strArray[index])
		fmt.Println("Will create Subscription:", __SubscriptionName)
	} else {
		panic(fmt.Sprintln("Could not find TOPIC:", __TOPIC))
	}

	//sortable_string.SearchStrings()
}

func unsubscribeTopic(httpClient *http.Client) {
	fmt.Println("EXE: Delete Subscription:", __SubscriptionName)
	unsubscribeTopicReader := strings.NewReader(deleteThroughProxy(httpClient, __UnsubscribeURL))
	//checkType("unsubscribeTopicReader", unsubscribeTopicReader)
	dec := json.NewDecoder(unsubscribeTopicReader)
	var v interface{}
	dec.Decode(&v)
	switch vv := v.(type) {
	case nil:
		//fmt.Println("nil value")
		fmt.Println("Unsubscribed from Subscription:", __SubscriptionName)
	case map[string]interface{}:
		for k_vv, v_vv := range vv {
			fmt.Println(k_vv, ":", v_vv)
		}
	}
}
func subscribeTopic(httpClient *http.Client) map[string]interface{} {
	fmt.Println("EXE: Create Subscription:", __SubscriptionName)
	var headers = map[string][]string{"Content-Type": {"application/vnd.kafka.v1+json"}}
	subscribeTopicBody := strings.NewReader("{\"name\": \"" + __TOPIC + "\", \"format\": \"avro\", \"auto.offset.reset\": \"smallest\"}")
	//checkType("subscribeTopicBody", subscribeTopicBody)
	subscribeTopicReader := strings.NewReader(postThroughProxy(httpClient, __SubscriptionURL, subscribeTopicBody, headers))
	//checkType("subscribeTopicReader", subscribeTopicReader)
	dec := json.NewDecoder(subscribeTopicReader)
	var v interface{}
	dec.Decode(&v)
	var output map[string]interface{}
	switch vv := v.(type) {
	case nil:
		fmt.Println("nil value")
	case map[string]interface{}:
		for k_vv, v_vv := range vv {
			fmt.Println(k_vv, ":", v_vv)
		}
		output = vv
	}
	return output
}
/*func checkType(name string, obj interface{}) {
	fmt.Println("TESTING BEGIN:", name, "::")
	switch obj.(type) {
	case *strings.Reader:
		fmt.Println("*strings.Reader")
	case strings.Reader:
		fmt.Println("strings.Reader")
	case *io.Reader:
		fmt.Println("*io.Reader")
	case io.Reader:
		fmt.Println("io.Reader")
	default:
		fmt.Println("DON'T KNOW WHAT's:", obj)
	}
	fmt.Println("TESTING END:", name)
}*/

/*func setupBGProxyServer() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "ssh", "-D", "8123", "-C", "-N", "n2cstech",)
		//"-i", "/Users/sokhan/.ssh/mmth_sohail_id_rsa", "-p", "722", "cstech@ewr-cs-n2.mediamath.com",)
	//if err := cmd.Run(); err != nil {
		// This will fail after 100 milliseconds. The 5 second sleep
		// will be interrupted.
	//	fmt.Fprintln(os.Stderr, fmt.Errorf("%s failed: %v", strings.Join(cmd.Args, " "), err))
	//}
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	fmt.Println("Result: " + out.String())
	/*go func() {
		time.Sleep(2 * time.Second)
		cancel()
	}()
	if err := cmd.Run(); err != nil {
		//return fmt.Errorf("%s failed: %v", strings.Join(cmd.Args, " "), err)
		fmt.Fprintln(os.Stderr, fmt.Errorf("%s failed: %v", strings.Join(cmd.Args, " "), err))
	}
}
var (
	username         = "cstech"
	password         = "password"
	serverAddrString = "ewr-cs-n2.mediamath.com:722"
	localAddrString  = "localhost:8123"
	remoteAddrString = "localhost:8123"
)

func forward(localConn net.Conn, config *ssh.ClientConfig) {
	// Setup sshClientConn (type *ssh.ClientConn)
	sshClientConn, err := ssh.Dial("tcp", serverAddrString, config)
	if err != nil {
		log.Fatalf("ssh.Dial failed: %s", err)
	}

	// Setup sshConn (type net.Conn)
	sshConn, err := sshClientConn.Dial("tcp", remoteAddrString)

	// Copy localConn.Reader to sshConn.Writer
	go func() {
		_, err = io.Copy(sshConn, localConn)
		if err != nil {
			log.Fatalf("io.Copy failed: %v", err)
		}
	}()

	// Copy sshConn.Reader to localConn.Writer
	go func() {
		_, err = io.Copy(localConn, sshConn)
		if err != nil {
			log.Fatalf("io.Copy failed: %v", err)
		}
	}()
}

func main() {
	// Setup SSH config (type *ssh.ClientConfig)
	//ssh.NewSignerFromKey()
	//rsa.PrivateKey{}
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
	}

	// Setup localListener (type net.Listener)
	localListener, err := net.Listen("tcp", localAddrString)
	if err != nil {
		log.Fatalf("net.Listen failed: %v", err)
	}

	for {
		// Setup localConn (type net.Conn)
		localConn, err := localListener.Accept()
		if err != nil {
			log.Fatalf("listen.Accept failed: %v", err)
		}
		go forward(localConn, config)
	}
}*/
