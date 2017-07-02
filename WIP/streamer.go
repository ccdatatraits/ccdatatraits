package main

import (
	"fmt"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net/http"
	"os"
	"encoding/json"
	"strings"
	"io"
	"log"
	"bytes"
	"sort"
)

const (
	PROXY_ADDR = "127.0.0.1:8123"
	__ENDPOINT = "http://ewr-changes-x6.ewr.mmracks.internal:8082"
)

var __TopicsURL = func() string { return __ENDPOINT + "/topics" }()
var __SubscriptionNameBuilder = func(topic string) string { return "ts_" + topic }
var __SubscriptionURLBuilder = func(topic string) string { return __ENDPOINT + "/consumers/" + __SubscriptionNameBuilder(topic) }
var __CommonURLBuilder = func(topic string) string { return __SubscriptionURLBuilder(topic) + "/instances/" + topic }
var __DataURLBuilder = func(topic string) string { return __CommonURLBuilder(topic) + "/topics/" + topic }
var __UnsubscribeURLBuilder = func(topic string) string { return __CommonURLBuilder(topic) }

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

type MessageValue struct {
	//Changes json.RawMessage
	Fields map[string]interface{}
	Wal_Buffer_Offset json.RawMessage
}
type Message struct {
	//Key float64
	//Partition, Offset float64
	Value MessageValue
}
type SortableString []string

func (s SortableString) Len() int {
	return len(s)
}
func (s SortableString) Less(i, j int) bool {
	return s[i] < s[j]
}
func (s SortableString) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (msg Message) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%6s(%6s) -> ", msg.Value.Fields["id"], msg.Value.Fields["version"]))
	fields_array := []string{}
	for k := range msg.Value.Fields {
		if k != "id" && k != "version" {
			fields_array = append(fields_array, k)
		}
	}

	sort.Sort(SortableString(fields_array))
	for _, k := range fields_array {
		buffer.WriteString(fmt.Sprintf("(\"%v\":%v), ", k, msg.Value.Fields[k]))
	}
	return buffer.String()
}

func getThroughProxy(check bool, debuggingLength int, httpClient *http.Client, url string, headers map[string][]string) string {
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
	summarisedInfo := make(map[string]string)
	debugging_check := func(dec *json.Decoder, iteration int) bool {
		if debuggingLength <= 0 {
			return dec.More()
		} else {
			return dec.More() && iteration < debuggingLength
		}
	}
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
		for debugging_check(dec, array_length) {
			var m Message
			// decode an array value (Message)
			err := dec.Decode(&m)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%v\n", m)
			array_length++
		}

		// read closing bracket
		_, err = dec.Token()
		if err != nil {
			log.Fatal(err)
		} else {
			resp.Body.Close()
		}
		//fmt.Printf("%T: %v\n", t, t)

		fmt.Println("ARRAY LENGTH:", array_length)
		fmt.Println("DATA MAP LENGTH:", len(summarisedInfo))
	} else {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error reading body:", err)
			os.Exit(4)
		}
		return string(b)
	}
	b, err := json.MarshalIndent(summarisedInfo, "\t", "")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error marshalling:", err)
		os.Exit(4)
	}
	return string(b)
}

func main() { //http.ProxyURL() ??
	//setupBGProxyServer() check from initial commits
	httpClient := setupHTTPClient()
	//topicsInterest := []string{"ewr.changes.history.batch-strategy-status.1","ewr.changes.history.component-creatives.1","ewr.changes.history.media-deals.1","ewr.changes.history.t1db.1","ewr.changes.history.tags.1","ewr.changes.snapshot.advertisers.1","ewr.changes.snapshot.agencies.1","ewr.changes.snapshot.atomic-creatives.1","ewr.changes.snapshot.audience-segments.1","ewr.changes.snapshot.audience-vendors.1","ewr.changes.snapshot.budget-flights.1","ewr.changes.snapshot.campaigns.1","ewr.changes.snapshot.concepts.1","ewr.changes.snapshot.currencies.1","ewr.changes.snapshot.deals.1","ewr.changes.snapshot.media-deals.1","ewr.changes.snapshot.organizations.1","ewr.changes.snapshot.pixel-bundles.1","ewr.changes.snapshot.pixel-providers.1","ewr.changes.snapshot.pixels.1","ewr.changes.snapshot.publishers.1","ewr.changes.snapshot.strategies.1","ewr.changes.snapshot.strategy-concepts.1","ewr.changes.snapshot.strategy-deals.1","ewr.changes.snapshot.supply-sources.1","ewr.changes.snapshot.target-dimensions.1","ewr.changes.snapshot.target-values.1","ewr.changes.snapshot.user-advertisers.1","ewr.changes.snapshot.user-agencies.1","ewr.changes.snapshot.user-organizations.1","ewr.changes.snapshot.users.1","ewr.changes.snapshot.vendors.1","ewr.changes.snapshot.viewability.1","ewr.kessel-run.mt-event-secondary.1","ewr.kessel-run.mt-event.1","mt_event"}
	topicsInterest := []string{"ewr.changes.snapshot.currency-rates.1", "ewr.changes.snapshot.budget-flights.1", "ewr.changes.snapshot.deals.1",}
	all_topics := getTopics(httpClient)
	DEBUGGING_LENGTH := 2
	for _, thisTopic := range all_topics {
		//fmt.Println(fmt.Sprint(" -> \"", thisTopic, "\""))
		for _, topicInterest := range topicsInterest {
			if thisTopic == topicInterest {
				fmt.Println(topicInterest)
				_ = subscribeTopic(httpClient, topicInterest)
				fetchData(httpClient, topicInterest, DEBUGGING_LENGTH)
				unsubscribeTopic(httpClient, topicInterest)
			}
		}
	}
}

func fetchData(httpClient *http.Client, topicInterest string, debuggingLength int) {
	fmt.Println("EXE: Fetch Data:", __DataURLBuilder(topicInterest))
	var headers = map[string][]string{"Accept":{"application/vnd.kafka.avro.v1+json"}}
	currency_map := getThroughProxy(true, debuggingLength, httpClient, __DataURLBuilder(topicInterest), headers)
	fmt.Println("DATA MAP:", currency_map)
}

func getTopics(httpClient *http.Client) []string {
	fmt.Println("EXE: Fetch Topics:", __TopicsURL)
	var headers map[string][]string = nil //map[string][]string{"Content-Type":{"what"}}
	getTopicsReader := strings.NewReader(getThroughProxy(false, 0, httpClient, __TopicsURL, headers))
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
	var strArray []string
	for _, topic := range topics {
		strArray = append(strArray, topic.(string))
	}
	return strArray
}

func unsubscribeTopic(httpClient *http.Client, topicInterest string) {
	fmt.Println("EXE: Delete Subscription:", __SubscriptionNameBuilder(topicInterest))
	unsubscribeTopicReader := strings.NewReader(deleteThroughProxy(httpClient, __UnsubscribeURLBuilder(topicInterest)))
	//checkType("unsubscribeTopicReader", unsubscribeTopicReader)
	dec := json.NewDecoder(unsubscribeTopicReader)
	var v interface{}
	dec.Decode(&v)
	switch vv := v.(type) {
	case nil:
		//fmt.Println("nil value")
		fmt.Println("Unsubscribed from Subscription:", __SubscriptionNameBuilder(topicInterest))
	case map[string]interface{}:
		for k_vv, v_vv := range vv {
			fmt.Println(k_vv, ":", v_vv)
		}
	}
}
func subscribeTopic(httpClient *http.Client, topicInterest string) map[string]interface{} {
	fmt.Println("EXE: Create Subscription:", __SubscriptionNameBuilder(topicInterest))
	var headers = map[string][]string{"Content-Type": {"application/vnd.kafka.v1+json"}}
	subscribeTopicBody := strings.NewReader("{\"name\": \"" + topicInterest + "\", \"format\": \"avro\", \"auto.offset.reset\": \"smallest\"}")
	//checkType("subscribeTopicBody", subscribeTopicBody)
	subscribeTopicReader := strings.NewReader(postThroughProxy(httpClient, __SubscriptionURLBuilder(topicInterest), subscribeTopicBody, headers))
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
